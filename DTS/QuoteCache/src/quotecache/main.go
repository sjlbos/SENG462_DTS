package main 

import(
	"net"
	"bytes"
	"time"
	"fmt"
	"os"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"
	"github.com/nu7hatch/gouuid"
	"log"
	"math/rand"
)

type QuoteServerEvent struct{
    EventType       string

    Guid            string
    OccuredAt       time.Time
    TransactionId   string
    UserId          string
    Service         string
    Server          string

    Price           string
    StockSymbol     string
    QuoteServerTime time.Time
    Cryptokey       string
}
type ErrorEvent struct{
    EventType       string

    Guid            string
    OccuredAt       time.Time
    TransactionId   string
    UserId          string
    Service         string
    Server          string

    Command         string
    StockSymbol     string
    Funds           string
    FileName        string
    ErrorMessage    string
}

type DebugEvent struct{
    EventType       string

    Guid            string
    OccuredAt       time.Time
    TransactionId   string
    UserId          string
    Service         string
    Server          string

    Command         string
    StockSymbol     string
    Funds           string
    FileName        string
    DebugMessage    string
}

type SystemEvent struct{
    EventType       string

    Guid            string
    OccuredAt       time.Time
    TransactionId   string
    UserId          string
    Service         string
    Server          string

    Command         string
    StockSymbol     string
    Funds           string
    FileName        string
}

var num_threads int

var rabbitConnectionString string
var rconn *amqp.Connection
var ch *amqp.Channel

var err error
var quotePort string

type QuoteCacheItem struct{
	Expiration      time.Time
	Value		string
}

func spawnQuoteThreads(messages chan string, timeout <-chan bool, stockSymbol string, APIUserId string){
	go getQuote(messages, timeout, stockSymbol, APIUserId)
	select{
		case message, ok := <-timeout:
			println(message)
			println(ok)
			return
		default:
			time.Sleep(time.Duration(10) * time.Millisecond)
			go getQuote(messages, timeout, stockSymbol, APIUserId)
			spawnQuoteThreads(messages, timeout, stockSymbol, APIUserId)
	}
}

func getQuote(messages chan string, timeout <-chan bool, stockSymbol string, APIUserId string) {
	sendString := stockSymbol + "," + APIUserId + "\n"
	addr, err := net.ResolveTCPAddr("tcp", "quoteserve.seng.uvic.ca:" + quotePort)
	if err != nil {
		return
	}
	qconn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		//error
		println("ERROR qconn: " + err.Error())
		return
	}
	_, err =fmt.Fprintf(qconn, sendString)
	if err != nil {
		failOnError(err, "Error with fprintf")
	}
	response := make([]byte, 100)
	_, err = qconn.Read(response)	
	qconn.Close()
	select { 
		case _, ok := <- timeout:
		if ok {
			messages <- string(response)
		}else{
			return
		}
	}
}

func getNewGuid() (uuid.UUID){
    guid,err := uuid.NewV4()
    if err != nil{
    }
    return *guid
}

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func SendRabbitMessage(message interface{}, EventType string){
    q := message

    if(EventType == "QuoteServerEvent"){
        q = message.(QuoteServerEvent)
    }else if(EventType == "ErrorEvent"){
        q = message.(ErrorEvent)
    }else if(EventType == "DebugEvent"){
        q = message.(DebugEvent)
    }else if(EventType == "SystemEvent"){
        q = message.(SystemEvent)
    }else{
       panic("NOT YET IMPLEMENTED")
    }
    
    body, err := json.Marshal(q)
    
    err = ch.Publish(
        "DtsEvents",          // exchange
        "TransactionEvent." + EventType, // routing key
        false, // mandatory
        false, // immediate
        amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(body),
        })
    failOnError(err, "Failed to publish a message")
}

func msToTime(ms string) (time.Time, error) {
    msInt, err := strconv.ParseInt(ms, 10, 64)
    if err != nil {
        return time.Time{}, err
    }

    return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

var memCache map[string]QuoteCacheItem
var miss float64 = 0
var hit  float64 = 0

func EfficiencyCalc(){
	for {
		var HitToMiss float64 = hit/miss

		var PercentHit float64 = hit/(hit+miss) 

		fmt.Printf("Hit To Miss: %f\n", HitToMiss)
		fmt.Printf("Percent Hit: %f\n\n", PercentHit) 

		time.Sleep(time.Duration(10) * time.Second)
	}
}


func handleConnection(conn net.Conn){
	var found bool
	var QuoteItem QuoteCacheItem
	if err != nil{
		// do stuff
		return
	}
	status := make([]byte, 100)
	_,err = conn.Read(status)
	if err != nil{
		// do stuff
		println("ERROR READ: " + err.Error())
		return
	}
	
	status = bytes.Trim(status, "\x00")

	inputs := strings.Split(string(status), ",")

	var price decimal.Decimal

	_, err 	:= strconv.ParseBool(strings.ToLower(inputs[0]))
	APIUserId	:= inputs[1]
	stockSymbol 	:= inputs[2]
	TransId := "1"
	GGuid := getNewGuid()
	Guid := GGuid.String()
	if len(inputs) > 3 {
		TransId 	= inputs[3]
		Guid		= inputs[4]
	}

	QuoteItem, found = memCache[stockSymbol]
	if found {
		if QuoteItem.Expiration.After(time.Now()){
			fmt.Fprintf(conn, string(QuoteItem.Value))
			conn.Close()
			hit = hit + 1
			return
		}else{
			found = false;
		}
	}
	miss = miss + 1
	if !found {
		messages := make(chan string)
		timeout := make(chan bool)
		startTime := time.Now()

		go spawnQuoteThreads(messages, timeout, stockSymbol, APIUserId)

		QuoteReturn := <-messages
		timeout <- true;


		diffTime := time.Since(startTime)
		println(diffTime.String())
		ParsedQuoteReturn := strings.Split(QuoteReturn,",")
		price, err = decimal.NewFromString(ParsedQuoteReturn[0])
		if err != nil{
			//error

		}
		_, err = conn.Write([]byte(price.String()))
		conn.Close()
		stockSymbol = ParsedQuoteReturn[1]
		ReturnUserId := ParsedQuoteReturn[2]
		msTimeStamp,err := msToTime(ParsedQuoteReturn[3])
		if err != nil{
			//error

		}
		cryptoKey := stripCtlAndExtFromUTF8(ParsedQuoteReturn[4])

		if ReturnUserId != APIUserId {
			// system error

		}

		QuoteEvent := QuoteServerEvent{
			EventType       : "QuoteServerEvent",
			Guid            : Guid,
			OccuredAt       : time.Now(),
			TransactionId   : TransId,
			UserId          : APIUserId,
			Service         : "QUOTE",
			Server          : "QuoteCache",
			Price           : price.String(),
			StockSymbol     : stockSymbol,
			QuoteServerTime : msTimeStamp,
			Cryptokey       : cryptoKey,
		}
		SendRabbitMessage(QuoteEvent,QuoteEvent.EventType)

		backoff := 35 + rand.Intn(10)

		tmpQuoteItem := QuoteCacheItem{
			Expiration : time.Now().Add(time.Duration(backoff)*time.Second),
			Value : price.String(),
		}
		memCache[stockSymbol] = tmpQuoteItem
		if err != nil {
		    // system error
		}
		for i := 0; i < num_threads-1; i++ {
			<- messages
		}
		close(messages)
		close(timeout)
		return
	}		
}

func main(){
    memCache = map[string]QuoteCacheItem{}


    type Configuration struct {
        	RabbitHost     string
		RabbitPort     string
		HostPort       string
		QuotePort      string
		NumThreads     string
    }
    file, err := os.Open("conf.json")
    if err != nil {
        fmt.Println("error:", err)
    }
    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    err = decoder.Decode(&configuration)
    if err != nil {
        fmt.Println("error:", err)
    }

    var rhost = configuration.RabbitHost
    var rport = configuration.RabbitPort
    quotePort = configuration.QuotePort
    num_threads,err = strconv.Atoi(configuration.NumThreads)
    if err != nil {
    	//error
    }

    var port string = configuration.HostPort

    println("Connected to Rabbit: " + rhost + ":" + rport)
    println("Connected to QuoteServer: " + "quoteserve.seng.uvic.ca:" + quotePort)
    println("Running locally on port " + port)

    rabbitConnectionString = "amqp://dts_user:Group1@" + rhost + ":" + rport + "/"
    rconn, err = amqp.Dial(rabbitConnectionString)
    failOnError(err, "Failed to connect to RabbitMQ")
    defer rconn.Close()

    ch, err = rconn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    err = ch.ExchangeDeclare(
            "DtsEvents", // name
            "topic",      // type
            true,         // durable
            true,        // auto-deleted
            false,        // internal
            false,        // no-wait
            nil,          // arguments
    )
    failOnError(err, "Failed to declare an exchange")

	//go EfficiencyCalc()

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			println("ERROR: " + err.Error())	
			// handle error
		}else{
			go handleConnection(conn)
		}
	}
}
