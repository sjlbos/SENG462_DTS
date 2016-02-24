package main

import (
    "encoding/json"
    "fmt"
//    "net"
    "net/http"
//    "os"
//    "log"
//    "strings"
//    "strconv"
    "time"
//    "math"
//    "io/ioutil"

    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
//    "github.com/streadway/amqp"
//    "github.com/nu7hatch/gouuid"
)

func Buy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Sell Request:")
    type buy_struct struct {
        strAmount string
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
    Amount,err := decimal.NewFromString(t.strAmount)
    if err != nil{

    }
    StockId := t.Symbol
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.strAmount)
    fmt.Fprintln(w, StockId)

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
        Command         : "BUY",
        StockSymbol     : t.Symbol,
        Funds           : t.strAmount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
 
    //Get A Quote
    quotePrice := getStockPrice(TransId ,"true", UserId, StockId, Guid.String())

    id, found, _ := getDatabaseUserId(UserId, "BUY") 
    if(found == false){
        Error := ErrorEvent{
            EventType       : "ErrorEvent",
            Guid            : Guid.String(),
            OccuredAt       : time.Now(),
            TransactionId   : TransId,
            UserId          : UserId,
            Service         : "API",
            Server          : Hostname,
            Command         : "BUY",
            StockSymbol     : "",
            Funds           : t.strAmount,
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)
    }else{
	    toBuy := (Amount.Div(quotePrice)).Floor()
        _, err = db.Exec(addPendingPurchase, id, t.Symbol, toBuy.String(), quotePrice.String(), time.Now(), time.Now().Add(time.Second*60))
        if(err != nil){
    		Error := ErrorEvent{
    		    EventType       : "ErrorEvent",
    		    Guid            : Guid.String(),
    		    OccuredAt       : time.Now(),
    		    TransactionId   : TransId,
    		    UserId          : UserId,
    		    Service         : "API",
    		    Server          : Hostname,
    		    Command         : "BUY",
    		    StockSymbol     : "",
    		    Funds           : t.strAmount,
    		    FileName        : "",
    		    ErrorMessage    : "Failed to create purchase",   
    		}
    		SendRabbitMessage(Error,Error.EventType)
                err = nil
        }
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
        Command         : "COMMIT_BUY",
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
            Command         : "COMMIT_BUY",
            StockSymbol     : "",
            Funds           : "",
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)        
    }else{
        LatestPendingrows, err := db.Query(getLatestPendingPurchase, id)
        defer LatestPendingrows.Close()
        failOnError(err, "Failed to Create Statement getLastestPendingPurchase for commitBuy")
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
                Command         : "COMMIT_BUY",
                StockSymbol     : "",
                Funds           : "",
                FileName        : "",
                ErrorMessage    : "No recent BUY commands issued",   
            }
            SendRabbitMessage(Error,Error.EventType)                  
        }else{
            _, err := db.Exec(commitPurchase, id, time.Now())
            if(err != nil){
                Error := ErrorEvent{
                    EventType       : "ErrorEvent",
                    Guid            : Guid.String(),
                    OccuredAt       : time.Now(),
                    TransactionId   : TransId,
                    UserId          : UserId,
                    Service         : "API",
                    Server          : Hostname,
                    Command         : "COMMIT_BUY",
                    StockSymbol     : stock,
                    Funds           : "",
                    FileName        : "",
                    ErrorMessage    : "Not Enough funds to buy Stocks",   
                }
                SendRabbitMessage(Error,Error.EventType)
            }
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
        Command         : "CANCEL_BUY",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
    id, found, _ := getDatabaseUserId(UserId, "CANCEL_BUY") 

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
        LatestPendingrows, err := db.Query(getLatestPendingPurchase, id)
        defer LatestPendingrows.Close()
        failOnError(err, "Failed to Create Statement: getLatestPurchase for cancelBuy")
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
                StockSymbol     : stock,
                Funds           : "",
                FileName        : "",
                ErrorMessage    : "No recent BUY commands issued",   
            }
            SendRabbitMessage(Error,Error.EventType)                  
        }else{
            Cancelrows, err := db.Query(cancelTransaction, id)
            defer Cancelrows.Close()
            failOnError(err, "Error with DB Query: cancelPurchase for cancelBuy")
        }   
    }
}
