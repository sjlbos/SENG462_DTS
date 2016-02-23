package main

import (
//    "encoding/json"
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

func Quote(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    StockId := vars["symbol"]
    UserId := vars["id"]
    TransId := vars["TransNo"]

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

}
