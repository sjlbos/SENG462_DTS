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
//    "io/ioutil"

    "github.com/gorilla/mux"
//    "github.com/streadway/amqp"
//    "github.com/nu7hatch/gouuid"
)

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
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : time.Now(),
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : Hostname,
        CommandType     : "BUY",
        StockSymbol     : t.Symbol,
        Funds           : t.Amount,
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


    rows, err := db.Query(getUserId, UserId)
    failOnError(err, "Failed to Create Statement")
    found := false
    var id int
    var userid string
    var balanceStr string
    
    for rows.Next() {
        found = true
        err = rows.Scan(&id, &userid, &balanceStr)
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
            Command         : "ADD",
            StockSymbol     : "",
            Funds           : t.Amount,
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)
    }else{
        _, err := db.Query(addPendingPurchase, id, t.Symbol, t.Amount, QuoteEvent.Price, time.Now(), time.Now().Add(time.Second*60))
        failOnError(err, "Failed to Create Statement")
    }
    
}

func CommitBuy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Buy Command Commited:") 
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    fmt.Fprintln(w, UserId)

    //Audit UserCommand
    Guid := getNewGuid()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : time.Now(),
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : Hostname,
        CommandType     : "COMMIT_BUY",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType); 
//TODO database Stuff
    rows, err := db.Query(getUserId, UserId)
    failOnError(err, "Failed to Create Statement")
    found := false
    var id int
    var userid string
    var balanceStr string

    for rows.Next() {
        found = true
        err = rows.Scan(&id, &userid, &balanceStr)
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
            Command         : "COMMIT_BUY",
            StockSymbol     : "",
            Funds           : "",
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)        
    }else{
        rows, err := db.Query(getLatestPendingPurchase, id)
        failOnError(err, "Failed to Create Statement")
        var id int
        var uid int 
        var stock string
        var num_shares int
        var share_price string
        var requested_at time.Time 
        var expires_at time.Time   
        found = false
        for rows.Next() {
            found = true
            err = rows.Scan(&id, &uid, &stock, &num_shares, &share_price, &requested_at, &expires_at)
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
                Command         : "COMMIT_BUY",
                StockSymbol     : "",
                Funds           : "",
                FileName        : "",
                ErrorMessage    : "No recent BUY commands issued",   
            }
            SendRabbitMessage(Error,Error.EventType)                  
        }else{
            _, err := db.Query(commitPurchase, id, time.Now())
            failOnError(err, "Error with DB Query")
        }   
    }
}

func CancelBuy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Buy Command Cancelled:") 
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]

    fmt.Fprintln(w, UserId) 

    //Audit UserCommand
    Guid := getNewGuid()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : time.Now(),
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : Hostname,
        CommandType     : "CANCEL_BUY",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
//TODO database Stuff
   rows, err := db.Query(getUserId, UserId)
    failOnError(err, "Failed to Create Statement")
    found := false
    var id int
    var userid string
    var balanceStr string

    for rows.Next() {
        found = true
        err = rows.Scan(&id, &userid, &balanceStr)
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
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)        
    }else{
        rows, err := db.Query(getLatestPendingPurchase, id)
        failOnError(err, "Failed to Create Statement")
        var id int
        var uid int 
        var stock string
        var num_shares int
        var share_price string
        var requested_at time.Time 
        var expires_at time.Time   
        found = false
        for rows.Next() {
            found = true
            err = rows.Scan(&id, &uid, &stock, &num_shares, &share_price, &requested_at, &expires_at)
        } 
        if(found == false){
            Error := ErrorEvent{
                EventType       : "DebugEvent",
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
            _, err := db.Query(cancelPurchase, id)
            failOnError(err, "Error with DB Query")

        }   
    }


}
