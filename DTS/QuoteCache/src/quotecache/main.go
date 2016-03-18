package main 

import(
	"net"
	"bytes"
	"time"
	//"sync"
	//"github.com/shopspring/decimal"
	//"github.com/streadway/amqp"
	"fmt"
	//"flag"
	"os"
	//"log"
	"encoding/json"
	//"strconv"
	"strings"
)

type QuoteCacheItem struct{
	Expiration      time.Time
	Value		string
}

var memCache map[string]QuoteCacheItem

var err error

func handleConnection(conn net.Conn){
	var found bool
	var QuoteItem QuoteCacheItem

	status := make([]byte, 100)
	_,err := conn.Read(status)
	if err != nil{
		// do stuff
	}
	
	status = bytes.Trim(status, "\x00")

	inputs := strings.Split(string(status), ",")
	found = false;

	if inputs[0] == "GET" {
		stockSymbol := inputs[1]
		QuoteItem, found = memCache[stockSymbol]
		if found {
			if QuoteItem.Expiration.After(time.Now()){
				fmt.Fprintf(conn, string(QuoteItem.Value))
			}else{
				found = false;
			}
		}
		if !found {
			fmt.Fprintf(conn, "-1")
			status = make([]byte,100)
			_,err := conn.Read(status)
			if err != nil {
			}
			status = bytes.Trim(status, "\x00")

			inputs := strings.Split(string(status), ",")
			if inputs[0] == "SET"{
				stockSymbol := inputs[1]
				price 	    := inputs[2]

				tmpQuoteItem := QuoteCacheItem{
			    	    Expiration	: time.Now().Add(time.Duration(time.Second*45)),
			       	    Value	: price,
				}
				memCache[stockSymbol] = tmpQuoteItem;
			}
		}
	} else if inputs[0] == "SET"{
		stockSymbol := inputs[1]
		price 		:= inputs[2]

		tmpQuoteItem := QuoteCacheItem{
	    	    Expiration	: time.Now().Add(time.Duration(time.Second*45)),
	       	    Value		: price,
		}
		memCache[stockSymbol] = tmpQuoteItem;
	}
	return
}

func main(){
    memCache = map[string]QuoteCacheItem{}

    type Configuration struct {
		HostPort       string
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
    var port string = configuration.HostPort


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
