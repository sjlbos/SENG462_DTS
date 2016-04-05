package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"log"
	"strings"
	"strconv"
	"time"
	"bytes"
    "database/sql"
	"github.com/streadway/amqp"
	"github.com/nu7hatch/gouuid"
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

type UserCommandEvent struct{
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
}

type AccountTransactionEvent struct{
    EventType       string

    Guid            string
    OccuredAt       time.Time
    TransactionId   string
    UserId          string
    Service         string
    Server          string

    AccountAction   string
    Funds           string
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

type TriggerEvent struct{
    TriggerType     string
    UserId          string
    TransactionId   string
    UpdatedAt       time.Time
}

func getStockPrice(TransId string, getNew string, UserId string, StockId string ,guid string) string {
    //Format Request String
	strEcho :=  getNew + "," + UserId + "," + StockId + "," + TransId + "," + guid + "\n"

    //Create Connection
	addr, err := net.ResolveTCPAddr("tcp", quoteCacheConnectionString)
    if err != nil {
        println("addr Error: " + err.Error())
        return "-1"
    }
	qconn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		//error
		println("Error Connecting to Quote Cache: " + err.Error())
		return "-1"
	}
	
    //Write Request String
	_, err = qconn.Write([]byte(strEcho))
	if err != nil {
		println("Write to server Error:", err.Error())
		return "-1"
	}

    //Create Reply
	reply := make([]byte, 100)
	_, err = qconn.Read(reply)
	reply = bytes.Trim(reply, "\x00")
	return string(reply)
}

func writeResponse(w http.ResponseWriter, responseCode int, response string){
	w.WriteHeader(200)
	fmt.Fprintln(w, response)
}

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}

func getTimeStamp() int64 {
  return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func getNewGuid() (uuid.UUID){
    guid,err := uuid.NewV4()
    if err != nil{
    }
    return *guid
}

func getDatabasePointerForUser(userid string) (*sql.DB){
    var hash = 0
    for i := range userid {
        hash += int(userid[i]);
    }
    hash = hash % len(dbPointers)
    return dbPointers[hash];
} 

func msToTime(ms string) (time.Time, error) {
    msInt, err := strconv.ParseInt(ms, 10, 64)
    if err != nil {
        return time.Time{}, err
    }

    return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func SendRabbitMessage(message interface{}, EventType string){
    q := message
    var mainKey string = "TransactionEvent"


    if(EventType == "QuoteServerEvent"){
        q = message.(QuoteServerEvent)
    }else if(EventType == "UserCommandEvent"){
        q = message.(UserCommandEvent)
    }else if(EventType == "AccountTransactionEvent"){
	    q = message.(AccountTransactionEvent)
    }else if(EventType == "ErrorEvent"){
        q = message.(ErrorEvent)
    }else if(EventType == "DebugEvent"){
        q = message.(DebugEvent)
    }else if(EventType == "SystemEvent"){
        q = message.(SystemEvent)
    }else if(EventType == "buy" || EventType =="sell"){
        q = message.(TriggerEvent)
        mainKey = "Triggers"
    }else{
       panic("NOT YET IMPLEMENTED")
    }
    
    body, err := json.Marshal(q)
    
    err = ch.Publish(
        "DtsEvents",          // exchange
        mainKey + "." + EventType, // routing key
        false, // mandatory
        false, // immediate
        amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(body),
        })
    failOnError(err, "Failed to publish a message")
}

func getDatabaseUserId(userId string) (*sql.DB, int, bool, string){
    db := getDatabasePointerForUser(userId)
    rows, err := db.Query(getUserId, userId)
    failOnError(err, "Failed to Create Statement: getUserId")
    found := false
    var id int
    var userid string
    var balanceStr string
    for rows.Next() {
       found = true
       err = rows.Scan(&id, &userid, &balanceStr)
    }
    rows.Close()
    return db, id, found, balanceStr
}
