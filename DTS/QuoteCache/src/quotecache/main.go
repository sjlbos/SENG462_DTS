package main 

import(
	"net"
	"bytes"
	"time"
	"sync"
	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"
	"fmt"
	"flag"
	"log"
	"encoding/json"
	"strconv"
	"strings"
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

type QuoteCacheItem struct{
	Expiration      time.Time
	Value		string
}

var memCache map[string]QuoteCacheItem
var memMutex sync.Mutex
var readMutex sync.Mutex

var maxconns int

var rabbitConnectionString string
var rconn *amqp.Connection
var ch *amqp.Channel

var err error


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

func handleConnection(conn net.Conn){
	status := make([]byte, 100)

	readMutex.Lock()
    	_, err = conn.Read(status)
	readMutex.Unlock()
	status = bytes.Trim(status, "\x00")

	inputs := strings.Split(string(status), ",")

	if len(inputs) != 4 {
		//invalid input
	}
	var price decimal.Decimal
	var found bool
	var QuoteItem QuoteCacheItem

	TransId 		:= inputs[0]
	getNew, err 	:= strconv.ParseBool(inputs[1])
	APIUserId	    := inputs[2]
	stockSymbol 	:= inputs[3]
	Guid			:= inputs[4]

	found = false
	num_threads := 2

	if !getNew {
		memMutex.Lock()
		QuoteItem, found = memCache[stockSymbol]
		memMutex.Unlock()
		if found{
			if QuoteItem.Expiration.Before(time.Now()){
				found = false
			}
		}
	}

	
	if !found {
                messages := make(chan string)
		for num_threads > 0 {
			go func() {
				sendString := stockSymbol + "," + APIUserId + "\n"
				qconn, err := net.Dial("tcp", "quoteserve.seng.uvic.ca:4441")
				if err != nil {
					//
				}
				_, err =fmt.Fprintf(qconn, sendString)
				if err!= nil {
					failOnError(err, "Error with fprintf")
				}
				response := make([]byte, 100)
    				_, err = qconn.Read(response)	
				messages <- string(response)
				qconn.Close()
			}()
			num_threads -= 1
		}
		QuoteReturn := <-messages
		ParsedQuoteReturn := strings.Split(QuoteReturn,",")
		price, err = decimal.NewFromString(ParsedQuoteReturn[0])
		if err != nil{
			//error
		
		}
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
		tmpQuoteItem := QuoteCacheItem{
	    	    Expiration 		    : time.Now().Add(time.Duration(time.Second*60)),
	       	    Value		    : price.String(),
		}
		// update map
		memMutex.Lock()
		memCache[stockSymbol] = tmpQuoteItem
		QuoteItem = tmpQuoteItem
		memMutex.Unlock()
		_, err = conn.Write([]byte(QuoteItem.Value))
		if err != nil {
		    // system error
		}
		for i := 0; i < 1; i++ {
			<- messages
		}
		close(messages)
	}else{
		//log system event returned Quote
		SystemEvent := SystemEvent{
			EventType       : "SystemEvent",
			Guid            : Guid,
			OccuredAt       : time.Now(),
			TransactionId   : TransId,
			UserId          : APIUserId,
			Service         : "QUOTE",
			Server          : "QuoteCache",
			Command         : price.String(),
			StockSymbol     : stockSymbol,
			Funds           : "",
			FileName        : "",
	        }
		SendRabbitMessage(SystemEvent,SystemEvent.EventType)
		_, err = conn.Write([]byte(QuoteItem.Value))
		if err != nil {
		    // system error
		}
	}
	return
}

func main(){
    memCache = map[string]QuoteCacheItem{}
    var rhost string

    flag.StringVar(&rhost,"rhost","b134.seng.uvic.ca","name of host for RabbitMQ")

    var rport string
    flag.StringVar(&rport,"rport","44410", "port number for RabbitMQ")

    var port string
    flag.StringVar(&port,"port", "44410", "port number for Cache")

    rabbitConnectionString = "amqp://dts_user:Group1@" + rhost + ":" + rport + "/"
    println(rabbitConnectionString)
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


	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}
