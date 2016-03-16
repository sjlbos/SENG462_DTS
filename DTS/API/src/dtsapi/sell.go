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



func Sell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Sell Request:")
    type sell_struct struct {
        strAmount string
        Symbol string
    }
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := r.Header.Get("TransNo")

    decoder := json.NewDecoder(r.Body)
    var t sell_struct   
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

    db := getDatabasePointerForUser(UserId)
    if err != nil{
        //error
    }

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
        Funds           : t.strAmount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    //Get A Quote
    var strPrice string
    strPrice = getStockPrice(TransId ,"true", UserId, StockId, Guid.String())

    var quotePrice decimal.Decimal
    quotePrice, err = decimal.NewFromString(strPrice)
    if err != nil || quotePrice == decimal.NewFromFloat(0){
        //error
        return;
    }

    id, found, _ := getDatabaseUserId(UserId) 
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
            Funds           : t.strAmount,
            FileName        : "",
            ErrorMessage    : "User Account Does Not Exist",   
        }
        SendRabbitMessage(Error,Error.EventType)
    }else{
        toSell := (Amount.Div(quotePrice)).Floor()
        _,err = db.Exec(addPendingSale, id, t.Symbol, toSell.String(), strPrice, time.Now(), time.Now().Add(time.Second*60))
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
                Funds           : t.strAmount,
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
    TransId := r.Header.Get("TransNo")

    fmt.Fprintln(w, UserId) 

    db := getDatabasePointerForUser(UserId)
    if err != nil{
        //error
    }

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

    id, found, _ := getDatabaseUserId(UserId) 

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
    TransId := r.Header.Get("TransNo")

    fmt.Fprintln(w, UserId) 

    db := getDatabasePointerForUser(UserId)
    if err != nil{
        //error
    }

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

    id, found,_ := getDatabaseUserId(UserId) 

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
