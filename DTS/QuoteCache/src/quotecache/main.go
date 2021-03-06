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

var thread_wait int

var rabbitConnectionString string
var rconn *amqp.Connection
var ch *amqp.Channel

var err error
var quotePort string

type QuoteCacheItem struct{
	Expiration      time.Time
	Value		string
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
		var PercentHit float64 = hit/(hit+miss) 
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
			fmt.Fprintf(conn, string(QuoteItem.Value) + "," + QuoteItem.Expiration.String())
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
		var returned bool = false
		num_threads := 0

		go func(){
			for returned == false {
				num_threads = num_threads + 1
				go func(){
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
					defer qconn.Close()
					_, err =fmt.Fprintf(qconn, sendString)
					if err != nil {
						failOnError(err, "Error with fprintf")
					}
					response := make([]byte, 100)
					_, err = qconn.Read(response)	
					returned = true
					messages <- string(response)
				}()
				time.Sleep(time.Duration(thread_wait) * time.Millisecond)
			}
		}()


		QuoteReturn := <-messages
		ParsedQuoteReturn := strings.Split(QuoteReturn,",")
		price, err = decimal.NewFromString(ParsedQuoteReturn[0])
		if err != nil{
			//error
			println("ERROR PARSING")
		}

		backoff := 50 + rand.Intn(5)
		QuoteExpiration := time.Now().Add(time.Duration(backoff)*time.Second)
		
		_, err = conn.Write([]byte(price.String() + "," + QuoteExpiration.String()))
		conn.Close()
		stockSymbol = ParsedQuoteReturn[1]
		ReturnUserId := ParsedQuoteReturn[2]
		msTimeStamp,err := msToTime(ParsedQuoteReturn[3])
		if err != nil{
			//error

		}
		cryptoKey := stripCtlAndExtFromUTF8(ParsedQuoteReturn[4])
		cryptoKey = strings.TrimSpace(cryptoKey)

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


		tmpQuoteItem := QuoteCacheItem{
			Expiration : QuoteExpiration,
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
		ThreadWait     string
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
    thread_wait,err = strconv.Atoi(configuration.ThreadWait)
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

	go EfficiencyCalc()

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
