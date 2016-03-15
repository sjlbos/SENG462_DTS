package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "time"
    "github.com/gorilla/mux"
)

func Add(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Adding Funds to account:")
    type add_struct struct {
        Amount float64
    }
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := vars["TransNo"]
    
    decoder := json.NewDecoder(r.Body)
    var t add_struct   
    err := decoder.Decode(&t)

    strAmount := strconv.FormatFloat(t.Amount, 'f', -1, 64)

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
        Command         : "ADD",
        StockSymbol     : "",
        Funds           : strAmount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType)

    var balanceFloat float64
    id, found, balanceStr := getDatabaseUserId(UserId, "ADD") 
    if(found == false){
    	/*Debug := DebugEvent{
            EventType       : "DebugEvent",
    	    Guid            : Guid.String(),
    	    OccuredAt       : time.Now(),
    	    TransactionId   : TransId,
    	    UserId          : UserId,
    	    Service         : "API",
    	    Server          : Hostname,
    	    Command         : "ADD",
    	    StockSymbol     : "",
           	Funds           : strAmount,
    	    FileName        : "",
    	    DebugMessage    : "Created User Account",   
        }*/
        //SendRabbitMessage(Debug,Debug.EventType)
        Addrows, _ := db.Query(addUser, UserId, strAmount, time.Now())
        defer Addrows.Close()
    }else{
        if(t.Amount < 0){
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
                Funds           : strAmount,
                FileName        : "",
                ErrorMessage    : "Amount to add is not a valid number",   
            }
            SendRabbitMessage(Error,Error.EventType)
        }else{
            balanceStr = strings.TrimLeft(balanceStr, "$")
            balanceFloat, err = strconv.ParseFloat(balanceStr, 64)
            newBalance := balanceFloat + t.Amount
            AccountEvent := AccountTransactionEvent{
                EventType       : "AccountTransactionEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "Account",
                Server          : Hostname,
                AccountAction   : "Add",
                Funds           : strAmount,
            }
            SendRabbitMessage(AccountEvent,AccountEvent.EventType)

            Updaterows, _ := db.Query(updateBalance, id, newBalance)
            defer Updaterows.Close()
        }
    }
}
