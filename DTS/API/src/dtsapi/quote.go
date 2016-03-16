package main

import (
//    "encoding/json"
    "fmt"
//    "net"
    "net/http"
//   "os"
//    "log"
//    "strings"
//    "strconv"
    "time"
//    "io/ioutil"

    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
//    "github.com/streadway/amqp"
//    "github.com/nu7hatch/gouuid"
)

func Quote(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    StockId := vars["symbol"]
    UserId := vars["id"]
    TransId := r.Header.Get("TransNo")

    //Audit UserCommand
    Guid := getNewGuid()
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : time.Now(),
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : "B134",
        Command         : "QUOTE",
        StockSymbol     : StockId,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    var strPrice string
    strPrice = getStockPrice(TransId ,"false", UserId, StockId, Guid.String())
    

    var price decimal.Decimal
    price, err := decimal.NewFromString(strPrice)
    if err != nil{
        //error
    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, StockId)
    fmt.Fprintln(w, price)

}
