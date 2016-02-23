package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "os"
//    "log"
    "strings"
    "strconv"
    "time"
    "math"
//    "io/ioutil"

    "github.com/gorilla/mux"
//    "github.com/streadway/amqp"
//    "github.com/nu7hatch/gouuid"
)



func Sell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Sell Request:")
    type sell_struct struct {
        Amount float64
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
    strAmount := strconv.FormatFloat(t.Amount, 'f', -1, 64)
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
        Server          : Hostname,
        Command         : "SELL",
        StockSymbol     : t.Symbol,
        Funds           : strAmount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    //Get A Quote
    strEcho :=  t.Symbol + "," + UserId + "\n"
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
        OccuredAt       : time.Now(),
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
    
    id, found, _ := getDatabaseUserId(UserId, "SELL") 
    if(found == false){
        Error := ErrorEvent{
            EventType       : "ErrorEvent",
            Guid            : Guid.String(),
            OccuredAt       : time.Now(),
            TransactionId   : TransId,
            UserId          : UserId,
            Service         : "API",
            Server          : Hostname,
            Command         : "SELL",
            StockSymbol     : "",
            Funds           : strAmount,
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)
    }else{

        quotePrice, err := strconv.ParseFloat(QuoteEvent.Price, 64)
        toSell := math.Floor(t.Amount/ quotePrice)
        _,err = db.Exec(addPendingSale, id, t.Symbol, toSell, QuoteEvent.Price, time.Now(), time.Now().Add(time.Second*60))
        if(err != nil){
            Error := ErrorEvent{
                EventType       : "ErrorEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "API",
                Server          : Hostname,
                Command         : "SELL",
                StockSymbol     : "",
                Funds           : strAmount,
                FileName        : "",
                ErrorMessage    : "Failed to create sale",   
            }
            SendRabbitMessage(Error,Error.EventType)
        }
    }      
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
        Server          : Hostname,
        Command         : "COMMIT_SELL",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    id, found, _ := getDatabaseUserId(UserId, "COMMIT_BUY") 

    if(found == false){
        Error := ErrorEvent{
            EventType       : "ErrorEvent",
            Guid            : Guid.String(),
            OccuredAt       : time.Now(),
            TransactionId   : TransId,
            UserId          : UserId,
            Service         : "API",
            Server          : Hostname,
            Command         : "COMMIT_SELL",
            StockSymbol     : "",
            Funds           : "",
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)        
    }else{
        LatestPendingrows, err := db.Query(getLatestPendingSale, id)
	defer LatestPendingrows.Close()
        failOnError(err, "Failed to Create Statement getLatestPendingSale for commitSell")
        var id int
        var uid int 
        var stock string
        var num_shares int
        var share_price string
        var requested_at time.Time 
        var expires_at time.Time   
        found = false
        for LatestPendingrows.Next() {
            found = true
            err = LatestPendingrows.Scan(&id, &uid, &stock, &num_shares, &share_price, &requested_at, &expires_at)
        } 
        if(found == false){
            Error := ErrorEvent{
                EventType       : "ErrorEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "API",
                Server          : Hostname,
                Command         : "COMMIT_SELL",
                StockSymbol     : "",
                Funds           : "",
                FileName        : "",
                ErrorMessage    : "No recent SELL commands issued",   
            }
            SendRabbitMessage(Error,Error.EventType)                  
        }else{
            _, err := db.Exec(commitSale, id, time.Now())
	    if(err != nil){
		Error := ErrorEvent{
		        EventType       : "ErrorEvent",
		        Guid            : Guid.String(),
		        OccuredAt       : time.Now(),
		        TransactionId   : TransId,
		        UserId          : UserId,
		        Service         : "API",
		        Server          : Hostname,
		        Command         : "COMMIT_SELL",
		        StockSymbol     : stock,
		        Funds           : "",
		        FileName        : "",
		        ErrorMessage    : "Can not Sell stocks",   
		}
		SendRabbitMessage(Error,Error.EventType)  
	    }
        } 
    }
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
        Server          : Hostname,
        Command         : "CANCEL_SELL",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    id, found,_ := getDatabaseUserId(UserId, "CANCEL_BUY") 

    if(found == false){
        Error := ErrorEvent{
            EventType       : "ErrorEvent",
            Guid            : Guid.String(),
            OccuredAt       : time.Now(),
            TransactionId   : TransId,
            UserId          : UserId,
            Service         : "API",
            Server          : Hostname,
            Command         : "CANCEL_BUY",
            StockSymbol     : "",
            Funds           : "",
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)        
    }else{
        LatestPendingrows, err := db.Query(getLatestPendingSale, id)
        defer LatestPendingrows.Close()
        failOnError(err, "Failed to Create Statement: getLatestSale for cancelSell")
        var id int
        var uid int 
        var stock string
        var num_shares int
        var share_price string
        var requested_at time.Time 
        var expires_at time.Time   
        found = false
        for LatestPendingrows.Next() {
            found = true
            err = LatestPendingrows.Scan(&id, &uid, &stock, &num_shares, &share_price, &requested_at, &expires_at)
        }
        if(found == false){
            Error := ErrorEvent{
                EventType       : "ErrorEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "API",
                Server          : Hostname,
                Command         : "CANCEL_BUY",
                StockSymbol     : "",
                Funds           : "",
                FileName        : "",
                ErrorMessage    : "No recent BUY commands issued",   
            }
            SendRabbitMessage(Error,Error.EventType)                  
        }else{
            Cancelrows, err := db.Query(cancelTransaction, id)
            defer Cancelrows.Close()
            failOnError(err, "Error with DB Query: cancelSale for cancelBuy")

        }   
    }

}
