package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)



func Sell(w http.ResponseWriter, r *http.Request){
	zero,_ := decimal.NewFromString("0");
	type sell_struct struct {
		Amount string
		Symbol string
	}
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	//Decode Body
	decoder := json.NewDecoder(r.Body)
	var t sell_struct   
	err := decoder.Decode(&t)

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
		Funds           : t.Amount,
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);
	if err != nil {
		//error
		return
	}

	//Check amount
	AmountDec,err := decimal.NewFromString(t.Amount)
	if err != nil{
		//error
		return
	}
	if(AmountDec.Cmp(zero) != 1){
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
	    writeResponse(w, http.StatusBadRequest, "Amount to buy is not a valid number")
	    return
	}

	//Check Stock Symbol
	StockId := t.Symbol
	if(len(StockId) == 0 || len(StockId) > 3){
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
		ErrorMessage    : "Symbol is Not Valid",   
	    }
	    SendRabbitMessage(Error,Error.EventType)
	    writeResponse(w, http.StatusBadRequest, "Symbol is Not Valid")
	    return
	}

	//get a database pointer
	db := getDatabasePointerForUser(UserId)

	//Get A Quote
	var strPrice string
	strPrice = getStockPrice(TransId ,"true", UserId, StockId, Guid.String())

	var quotePrice decimal.Decimal
	quotePrice, err = decimal.NewFromString(strPrice)
	if err != nil{
		//error
		return;
	}
	if(quotePrice.Cmp(zero) != 1){
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
	    writeResponse(w, http.StatusBadRequest, "Amount to buy is not a valid number")
	    return
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
			Funds           : t.Amount,
			FileName        : "",
			ErrorMessage    : "User Account Does Not Exist",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//error
		return
	}


	toSell := (AmountDec.Div(quotePrice)).Floor()
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
			Funds           : t.Amount,
			FileName        : "",
			ErrorMessage    : "Failed to create sale",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//error
		return
	}else{
		//success
		return
	}      
}


func CommitSell(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	//get a db pointer
	db := getDatabasePointerForUser(UserId)
	if err != nil{
		//error
		return
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

	//Find user in database
	uid, found, _ := getDatabaseUserId(UserId) 
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
	}

	//find last sell
	LatestPendingrows, err := db.Query(getLatestPendingSale, uid)
	defer LatestPendingrows.Close()
	var id int
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
		//error
		return                 
	}
	_, err = db.Exec(commitSale, id, time.Now())
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
		//error
		return
	}else{
		//success
		return
	} 
}

func CancelSell(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	//get DB Pointer
	db := getDatabasePointerForUser(UserId)
	if err != nil{
		//error
		return
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

	//Find user in DB
	uid, found,_ := getDatabaseUserId(UserId) 
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
	}

	//Find last Sell Command
	LatestPendingrows, err := db.Query(getLatestPendingSale, uid)
	defer LatestPendingrows.Close()
	if err != nil{
		//error
		return
	}
	var id int
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
		//error
		return               
	}
	_, err = db.Exec(cancelTransaction, id)
	if err != nil{
		//error
		return
	}else{
		//success
		return
	}
}
