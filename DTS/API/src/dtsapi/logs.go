package main

import (
//    "encoding/json"
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
//   "github.com/nu7hatch/gouuid"
)


func DumplogUser(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Display Log for User: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := r.Header.Get("TransNo")
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
        Command         : "DUMPLOG",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
}


func Dumplog(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Dumplog of all transactions: ")
    TransId := r.Header.Get("TransNo")

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
        Command         : "DUMPLOG",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

}


func DisplaySummary(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Display Summary for User: ")

    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := r.Header.Get("TransNo")

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
        Command         : "DISPLAY_SUMMARY",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    fmt.Fprintln(w, UserId)
}




