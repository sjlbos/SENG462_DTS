package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
//    "strconv"
    "time"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

func Add(w http.ResponseWriter, r *http.Request){
    zero,_ := decimal.NewFromString("0");
    fmt.Fprintln(w, "Adding Funds to account:")
    type add_struct struct {
        strAmount string
    }
    vars := mux.Vars(r)
    UserId := vars["id"]
    TransId := r.Header.Get("TransNo")
    
    decoder := json.NewDecoder(r.Body)
    var t add_struct   
    err := decoder.Decode(&t)

    Amount,err := decimal.NewFromString(t.strAmount)

    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.strAmount)

    db := getDatabasePointerForUser(UserId)
    if err != nil{
        //error
    }

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
        Funds           : t.strAmount,
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType)

    var balanceStr string
    var balance decimal.Decimal

    id, found, balanceStr := getDatabaseUserId(UserId) 
    if(found == false){
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
           	Funds           : t.strAmount,
    	    FileName        : "",
    	    DebugMessage    : "Created User Account",   
        }
        SendRabbitMessage(Debug,Debug.EventType)
        Addrows, _ := db.Query(addUser, UserId, t.strAmount, time.Now())
        defer Addrows.Close()
    }else{
        if(Amount.Cmp(zero) == -1){
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
                Funds           : t.strAmount,
                FileName        : "",
                ErrorMessage    : "Amount to add is not a valid number",   
            }
            SendRabbitMessage(Error,Error.EventType)
        }else{
            balanceStr = strings.TrimLeft(balanceStr, "$")
            balance, err = decimal.NewFromString(balanceStr)
            newBalance := balance.Add(Amount)
            AccountEvent := AccountTransactionEvent{
                EventType       : "AccountTransactionEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "Account",
                Server          : Hostname,
                AccountAction   : "Add",
                Funds           : t.strAmount,
            }
            SendRabbitMessage(AccountEvent,AccountEvent.EventType)

            Updaterows, _ := db.Query(updateBalance, id, newBalance)
            defer Updaterows.Close()
        }
    }
}
