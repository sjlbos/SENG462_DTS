package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "os"
    "log"
    "strings"
    "strconv"
    "time"
    "github.com/shopspring/decimal"
//    "io/ioutil"

//    "github.com/gorilla/mux"
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
/* 
q := QuoteServerEvent{
    EventType       : "AccountTransaction",
    Guid            : "",
    OccuredAt       : "",
    TransactionId   : "",
    UserId          : "",
    Service         : "",
    Server          : "",
    AccountAction   : "",
    Funds           : "",
}
*/

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

/*q := ErrorEvent{
    EventType       : "ErrorEvent",
    Guid            : "",
    OccuredAt       : "",
    TransactionId   : "",
    UserId          : "",
    Service         : "",
    Server          : "",
    Command         : "",
    StockSymbol     : "",
    Funds           : "",
    FileName        : "",
    ErrorMessage    : "",   
}*/

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
/*q := DebugEvent{
    EventType       : "DebugEvent",
    Guid            : "",
    OccuredAt       : "",
    TransactionId   : "",
    UserId          : "",
    Service         : "",
    Server          : "",
    Command         : "",
    StockSymbol     : "",
    Funds           : "",
    FileName        : "",
    DebugMessage    : "",   
}*/

func getStockPrice(TransId string, getNew string, UserId string, StockId string ,guid string) decimal.Decimal {
    strEcho :=  TransId + "," + getNew + "," + StockId + "," + UserId + "," + guid + "\n"

    tcpAddr, err := net.ResolveTCPAddr("tcp", quoteCacheConnectionString)
    if err != nil {
        println("ResolveTCPAddr failed:", err.Error())
        os.Exit(1)
    }

    qconn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        println("Dial failed:", err.Error())
        os.Exit(1)
    }

    _, err = qconn.Write([]byte(strEcho))
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }
    
    reply := make([]byte, 1024)
    _, err = qconn.Read(reply)
    qconn.Close()
    result, err := decimal.NewFromString(string(reply))
    if err != nil{
        //error
    }
    return result
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



func getDatabaseUserId(userId string, commandStr string) (int, bool, string){
    rows, err := db.Query(getUserId, userId)
    failOnError(err, "Failed to Create Statement: getUserId for add.go")
    found := false
    var id int
    var userid string
    var balanceStr string


    for rows.Next() {
       found = true
       err = rows.Scan(&id, &userid, &balanceStr)
    }
    rows.Close()
    return id, found, balanceStr
}
