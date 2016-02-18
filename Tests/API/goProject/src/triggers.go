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
//    "io/ioutil"

    "github.com/gorilla/mux"
//    "github.com/streadway/amqp"
//    "github.com/nu7hatch/gouuid"
)

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