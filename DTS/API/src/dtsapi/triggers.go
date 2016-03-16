package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

func CreateBuyTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Creating Buy Trigger:") 
    type trigger_struct struct{
        strAmount string
        Price  string
    }

    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]
    TransId := r.Header.Get("TransNo")

    db := getDatabasePointerForUser(UserId)

    decoder := json.NewDecoder(r.Body)
    var t trigger_struct   
    err = decoder.Decode(&t)
    if err != nil {
        //error
        return
    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, Symbol)
    fmt.Fprintln(w, t.strAmount)
    fmt.Fprintln(w, t.Price)

    uid, found, _ := getDatabaseUserId(UserId)
    Amount, err := decimal.NewFromString(t.strAmount)
    if err != nil{
        //error
        return
    }
    if !found {
        //error
        return
    }

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()
    if(t.strAmount != ""){
	    CommandEvent := UserCommandEvent{
    		EventType       : "UserCommandEvent",
    		Guid            : Guid.String(),
    		OccuredAt       : OccuredAt,
    		TransactionId   : TransId,
    		UserId          : UserId,
    		Service         : "Command",
    		Server          : "B134",
    		Command         : "SET_BUY_AMOUNT",
    		StockSymbol     : Symbol,
    		Funds           : "",
	    }
	    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
        _, err = db.Exec(addBuyTrigger, uid, Symbol, t.strAmount, time.Now())
        if err != nil{
            Error := ErrorEvent{
                EventType       : "ErrorEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "API",
                Server          : Hostname,
                Command         : "BUY",
                StockSymbol     : Symbol,
                Funds           : t.strAmount,
                FileName        : "",
                ErrorMessage    : "Unable to Create Buy Trigger",   
            }
            SendRabbitMessage(Error,Error.EventType)
        }
        return

    }else{
	    CommandEvent := UserCommandEvent{
    		EventType       : "UserCommandEvent",
    		Guid            : Guid.String(),
    		OccuredAt       : OccuredAt,
    		TransactionId   : TransId,
    		UserId          : UserId,
    		Service         : "Command",
    		Server          : "B134",
    		Command         : "SET_BUY_TRIGGER",
    		StockSymbol     : Symbol,
    		Funds           : "",
	    }
	    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
        
        getPendingTriggerRows, err := db.Query(getPendingTriggerId, uid, Symbol, "buy")
        defer getPendingTriggerRows.Close()

        var id int
        found = false
        for getPendingTriggerRows.Next() {
            found = true
            err = getPendingTriggerRows.Scan(&id)
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
                StockSymbol     : Symbol,
                Funds           : Amount.String(),
                FileName        : "",
                ErrorMessage    : "No recent SET_BUY_AMOUNT commands issued",   
            }
            SendRabbitMessage(Error,Error.EventType)                  
        }else{
            _, err = db.Exec(setBuyTrigger, id, uid, t.Price, time.Now())
            if err != nil {
                //error
                return
            }
            Trigger := TriggerEvent{
                TriggerType     : "buy",
                UserId          : UserId,
                TransactionId   : TransId,
                UpdatedAt       : time.Now(),
            }
            SendRabbitMessage(Trigger, "buy")
        }
    }
}

func CreateSellTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Creating Sell Trigger:") 
    type trigger_struct struct{
        strAmount string
        Price  string
    }

    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]
    TransId := r.Header.Get("TransNo")

    decoder := json.NewDecoder(r.Body)
    var t trigger_struct   
    err := decoder.Decode(&t)

    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, Symbol)
    fmt.Fprintln(w, t.strAmount)
    fmt.Fprintln(w, t.Price)

    //Audit UserCommand
    Guid := getNewGuid()
    OccuredAt := time.Now()

    uid, found, _ := getDatabaseUserId(UserId)
    if !found {
        //error
        return
    }

    db := getDatabasePointerForUser(UserId)

    if(t.strAmount != ""){
	    CommandEvent := UserCommandEvent{
		EventType       : "UserCommandEvent",
		Guid            : Guid.String(),
		OccuredAt       : OccuredAt,
		TransactionId   : TransId,
		UserId          : UserId,
		Service         : "Command",
		Server          : "B134",
		Command         : "SET_SELL_AMOUNT",
		StockSymbol     : Symbol,
		Funds           : "",
	    }
	    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

        _, err = db.Exec(addSellTrigger, uid, Symbol, t.strAmount, time.Now())
        if err != nil{
            Error := ErrorEvent{
                EventType       : "ErrorEvent",
                Guid            : Guid.String(),
                OccuredAt       : time.Now(),
                TransactionId   : TransId,
                UserId          : UserId,
                Service         : "API",
                Server          : Hostname,
                Command         : "BUY",
                StockSymbol     : Symbol,
                Funds           : t.strAmount,
                FileName        : "",
                ErrorMessage    : "Unable to Create Buy Trigger",   
            }
            SendRabbitMessage(Error,Error.EventType)
        }
    }else{
	    CommandEvent := UserCommandEvent{
		EventType       : "UserCommandEvent",
		Guid            : Guid.String(),
		OccuredAt       : OccuredAt,
		TransactionId   : TransId,
		UserId          : UserId,
		Service         : "Command",
		Server          : "B134",
		Command         : "SET_SELL_TRIGGER",
		StockSymbol     : Symbol,
		Funds           : "",
	    }
	    SendRabbitMessage(CommandEvent,CommandEvent.EventType)

        getPendingTriggerRows, err := db.Query(getPendingTriggerId, uid, Symbol, "sell")
        defer getPendingTriggerRows.Close()

        var id int
        found = false
        for getPendingTriggerRows.Next() {
            found = true
            err = getPendingTriggerRows.Scan(&id)
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
                StockSymbol     : Symbol,
                Funds           : t.strAmount,
                FileName        : "",
                ErrorMessage    : "No recent SET_SELL_AMOUNT commands issued",   
            }
            SendRabbitMessage(Error,Error.EventType)                  
        }else{
            _, err = db.Exec(setSellTrigger, id, uid, t.Price, time.Now())
            if err != nil {
                //error
                return
            }
            Trigger := TriggerEvent{
                TriggerType     : "sell",
                UserId          : UserId,
                TransactionId   : TransId,
                UpdatedAt       : time.Now(),
            }
            SendRabbitMessage(Trigger, "sell")
        }
    }
}

func CancelBuyTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Cancelling Buy Trigger: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]    
    TransId := r.Header.Get("TransNo")
    
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
        Command         : "CANCEL_SET_BUY",
        StockSymbol     : Symbol,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

    uid, found, _ := getDatabaseUserId(UserId)

    if found {
        db := getDatabasePointerForUser(UserId)
        var id int
        getBuyTriggerIdRows, err := db.Query(getBuyTriggerId, uid, Symbol)
        defer getBuyTriggerIdRows.Close()
        found := false
        for getBuyTriggerIdRows.Next() {
            found = true
            err = getBuyTriggerIdRows.Scan(&id)
        }
        _, err = db.Exec(cancelBuyTrigger, id)
        if err != nil{
            //error
            return
        }
        Trigger := TriggerEvent{
            TriggerType     : "buy",
            UserId          : UserId,
            TransactionId   : TransId,
            UpdatedAt       : time.Now(),
        }
        SendRabbitMessage(Trigger, "buy")
    }else {
        //error
    }
}

func CancelSellTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Cancelling Sell Trigger: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]
    TransId := r.Header.Get("TransNo")

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
        Command         : "CANCEL_SET_SELL",
        StockSymbol     : Symbol,
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
    uid, found, _ := getDatabaseUserId(UserId)

    if found {
        db := getDatabasePointerForUser(UserId)
        var id int
        getSellTriggerIdRows, err := db.Query(getSellTriggerId, uid, Symbol)
        defer getSellTriggerIdRows.Close()
        found := false
        for getSellTriggerIdRows.Next() {
            found = true
            err = getSellTriggerIdRows.Scan(&id)
        }
        _, err = db.Exec(cancelSellTrigger, id)
        if err != nil{
            //error
            return
        }
        Trigger := TriggerEvent{
            TriggerType     : "buy",
            UserId          : UserId,
            TransactionId   : TransId,
            UpdatedAt       : time.Now(),
        }
        SendRabbitMessage(Trigger, "buy")
    }else {
        //error
    }
}
