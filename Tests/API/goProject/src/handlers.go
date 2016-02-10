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
//    "io/ioutil"

    "github.com/gorilla/mux"
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
/* 
q := QuoteServerEvent{
    EventType       : "QuoteServerEvent",
    Guid            : "",
    OccuredAt       : "",
    Transactionid   : "",
    UserId          : "",
    Service         : "",
    Server          : "",
    Price           : result[0],
    StockSymbol     : result[1],
    QuoteServerTime : tmpResult3,
    Cryptokey       : tmpResult4
}
*/
type UserCommandEvent struct{
    EventType       string

    Guid            string
    OccuredAt       time.Time
    TransactionId   string
    UserId          string
    Service         string
    Server          string

    CommandType     string
    StockSymbol     string
    Funds           string
}
/* 
q := UserCommandEvent{
    EventType       : "QuoteServerEvent",
    Guid            : "",
    OccuredAt       : "",
    Transactionid   : "",
    UserId          : "",
    Service         : "",
    Server          : "",
    CommandType     : "",
    StockSymbol     : result[1],
    Funds           : ""
}
*/
type AccountTransactionEvent struct{
    EventType       string

    Guid            int64
    OccuredAt       time.Time
    TransactionId   int64
    UserId          string
    Service         string
    Server          string

    AccountAction   string
    Funds           string
}

type SystemEvent struct{
    EventType       string

    Guid            int64
    OccuredAt       time.Time
    TransactionId   int64
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

    Guid            int64
    OccuredAt       time.Time
    TransactionId   int64
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

    Guid            int64
    OccuredAt       time.Time
    TransactionId   int64
    UserId          string
    Service         string
    Server          string

    Command         string
    StockSymbol     string
    Funds           string
    FileName        string
    DebugMessage    string
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
    rconn, err := amqp.Dial("amqp://dts_user:Group1@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer rconn.Close()

    ch, err := rconn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()
    
    q := message

    if(EventType == "QuoteServerEvent"){
        q = message.(QuoteServerEvent)
    }else if(EventType == "UserCommandEvent"){
        q = message.(UserCommandEvent)
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

    log.Printf(" [x] Sent %s", body)

    rconn.Close()
}



func Quote(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    StockId := vars["symbol"]
    UserId := vars["id"]
    TransId := vars["TransNo"]

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "QUOTE",
        StockSymbol     : StockId,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    strEcho :=  StockId + "," + UserId + "\n"
    servAddr := "quoteserve.seng.uvic.ca:4444"

    tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
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
    result := strings.Split(string(reply),",")
    qconn.Close()

    tmpResult0, err := strconv.ParseFloat(result[0], 64)
    fmt.Fprintln(w, tmpResult0)
    fmt.Fprintln(w, result[1])
    fmt.Fprintln(w, result[2])
    fmt.Fprintln(w, result[3])
    tmpResult3,err := msToTime(result[3])
    tmpResult4 := stripCtlAndExtFromUTF8(result[4])

    QuoteEvent := QuoteServerEvent{
        EventType       : "QuoteServerEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "QUOTE",
        Server          : "quoteserve",
        Price           : result[0],
        StockSymbol     : result[1],
        QuoteServerTime : tmpResult3,
        Cryptokey       : tmpResult4,
    }
    SendRabbitMessage(QuoteEvent,QuoteEvent.EventType)

}

func Add(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Adding Funds to account:")
    type add_struct struct {
        Amount string
    }
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]
    
    decoder := json.NewDecoder(r.Body)
    var t add_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.Amount)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "ADD",
        StockSymbol     : "",
        Funds           : t.Amount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);


//TODO database stuff!

}

func Buy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Sell Request:")
    type buy_struct struct {
        Amount string
        Symbol string
    }
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    decoder := json.NewDecoder(r.Body)
    var t buy_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Symbol)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "BUY",
        StockSymbol     : t.Symbol,
        Funds           : t.Amount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

//TODO database stuff!    
    
}

func CommitBuy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Buy Command Commited:") 
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    fmt.Fprintln(w, UserId)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "COMMIT_BUY",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType); 
//TODO database Stuff

}

func CancelBuy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Buy Command Cancelled:") 
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    fmt.Fprintln(w, UserId) 

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "CANCEL_BUY",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
//TODO database Stuff

}

func Sell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Sell Request:")
    type sell_struct struct {
        Amount string
        Symbol string
    }
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    decoder := json.NewDecoder(r.Body)
    var t sell_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Symbol)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "SELL",
        StockSymbol     : t.Symbol,
        Funds           : t.Amount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

//TODO database stuff!    
    
}


func CommitSell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Sell Command Commited:") 
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    fmt.Fprintln(w, UserId) 

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "COMMIT_SELL",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
//TODO database Stuff

}

func CancelSell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Sell Command Cancelled:") 
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    fmt.Fprintln(w, UserId) 

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "CANCEL_SELL",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
//TODO database Stuff

}

func CreateBuyTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Creating Buy Trigger:") 
    type trigger_struct struct{
        Amount string
        Price  string
    }

    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]
    TransId := vars["TransNo"]

    decoder := json.NewDecoder(r.Body)
    var t trigger_struct   
    err := decoder.Decode(&t)

    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, Symbol)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Price)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "CREATE_BUY_TRIGGER",
        StockSymbol     : Symbol,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

//TODO database stuff

}

func CreateSellTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Creating Sell Trigger:") 
    type trigger_struct struct{
        Amount string
        Price  string
    }

    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]
    TransId := vars["TransNo"]

    decoder := json.NewDecoder(r.Body)
    var t trigger_struct   
    err := decoder.Decode(&t)

    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, Symbol)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Price)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "CREATE_SELL_TRIGGER",
        StockSymbol     : Symbol,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

//TODO database Stuff

}


func CancelBuyTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Cancelling Buy Trigger: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]    
    TransId := vars["TransNo"]
    
    fmt.Fprintln(w, UserId)   
    fmt.Fprintln(w, Symbol) 

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "CANCEL_BUY_TRIGGER",
        StockSymbol     : Symbol,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
//TODO database Stuff

}



func CancelSellTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Cancelling Sell Trigger: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]
    TransId := vars["TransNo"]

    fmt.Fprintln(w, UserId)   
    fmt.Fprintln(w, Symbol) 
    
    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "CANCEL_SELL_TRIGGER",
        StockSymbol     : Symbol,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
     
//TODO database Stuff

}

func DumplogUser(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Display Log for User: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]
    fmt.Fprintln(w, UserId)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        CommandType     : "DUMPLOG",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
}


func Dumplog(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Dumplog of all transactions: ")
    vars := mux.Vars(r)
    TransId := vars["TransNo"]

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : "",
        Service         : "Command",
        Server          : "B134",
        CommandType     : "DUMPLOG",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

}


func DisplaySummary(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Display Summary for User: ")

    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : OccuredAt,
        TransactionId   : TransId,
        UserId          : "",
        Service         : "Command",
        Server          : "B134",
        CommandType     : "DISPLAY_SUMMARY",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    fmt.Fprintln(w, UserId)
}





