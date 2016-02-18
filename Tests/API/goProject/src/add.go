package main

import (
    "encoding/json"
    "fmt"
//    "net"
    "net/http"
//    "os"
//    "log"
    "strings"
    "strconv"
    "time"
//    "io/ioutil"

    "github.com/gorilla/mux"
//    "github.com/streadway/amqp"
//    "github.com/nu7hatch/gouuid"
)

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
    CommandEvent := UserCommandEvent{
        EventType       : "UserCommandEvent",
        Guid            : Guid.String(),
        OccuredAt       : time.Now(),
        TransactionId   : TransId,
        UserId          : UserId,
        Service         : "Command",
        Server          : Hostname,
        CommandType     : "ADD",
        StockSymbol     : "",
        Funds           : t.Amount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType)

    rows, err := db.Query(getUserId, UserId)
    failOnError(err, "Failed to Create Statement")
    //found := false
    var id int = -1
    var userid string
    var balanceStr string
    var balanceFloat float64
    var amountFloat float64

    for rows.Next() {
	//found = true
	err = rows.Scan(&id, &userid, &balanceStr)
    }
    if(id == -1){
	Debug := DebugEvent{
            EventType       : "DebugEvent",
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
	    DebugMessage    : "Created User Account",   
        }
        SendRabbitMessage(Debug,Debug.EventType)
        db.Query(addUser, UserId, t.Amount, time.Now())
    }else{
        amountFloat, err = strconv.ParseFloat(t.Amount, 64)
        failOnError(err, "Amount is not a vaid number, Transaction" + TransId)
        if(amountFloat < 0){
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
                ErrorMessage    : "Amount to add is not a valid number",   
            }
            SendRabbitMessage(Error,Error.EventType)
        }else{
            balanceStr = strings.TrimLeft(balanceStr, "$")
            balanceFloat, err = strconv.ParseFloat(balanceStr, 64)

            newBalance := balanceFloat + amountFloat
            AccountEvent := AccountTransactionEvent{
                EventType       : "AccountTransactionEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "Account",
                Server          : Hostname,
                AccountAction   : "Add",
                Funds           : t.Amount,
            }
            SendRabbitMessage(AccountEvent,AccountEvent.EventType)

            db.Query(updateBalance, id, newBalance)
        }
    }
}
