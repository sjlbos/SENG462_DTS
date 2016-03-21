package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"fmt"
)

func Buy(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Buying Stock");
	zero,_ := decimal.NewFromString("0");
	type buy_struct struct {
		Amount string
		Symbol string
	}
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	decoder := json.NewDecoder(r.Body)
	var t buy_struct   
	err := decoder.Decode(&t)

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
		Command         : "BUY",
		StockSymbol     : t.Symbol,
		Funds           : t.Amount,
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);
	if err != nil {
		//error decoding
		return
	}

	//Error Checking
	AmountDec,err := decimal.NewFromString(t.Amount)
	if err != nil{
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
	    //writeResponse(w, http.StatusBadRequest, "Amount to buy is not a valid number")
	    return
	}
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
	    //writeResponse(w, http.StatusBadRequest, "Symbol is Not Valid")
	    return
	}

	//Open DB connection
	db := getDatabasePointerForUser(UserId)

	//Get A Quote
	var strPrice string
	strPrice = getStockPrice(TransId ,"true", UserId, StockId, Guid.String())
	var quotePrice decimal.Decimal
	quotePrice, err = decimal.NewFromString(strPrice)
	if err != nil {
		//writeResponse(w, http.StatusInternalServerError, "Quote Return is not Valid")
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
			ErrorMessage    : "Quote is not greater than 0",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//writeResponse(w, http.StatusBadRequest, "Amount to buy is not a valid number")
		return
	}
	//Check If User Exists
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
			Command         : "BUY",
			StockSymbol     : "",
			Funds           : t.Amount,
			FileName        : "",
			ErrorMessage    : "User Does not Exist",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//writeResponse(w, http.StatusBadRequest, "User Does Not Exist")
		return
	}
	toBuy := (AmountDec.Div(quotePrice)).Floor()
	
	//Check to make sure amount is valid
	if(toBuy.Cmp(zero) != 1){
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
			ErrorMessage    : "Cannot Buy less than 1 stock",   
	    }
	    SendRabbitMessage(Error,Error.EventType)
	    //writeResponse(w, http.StatusBadRequest, "Cannot Buy " + toBuy.String() + " stock")
	    return
	}

	//commit buy 
	_, err = db.Exec(addPendingPurchase, id, t.Symbol, toBuy.String(), strPrice, time.Now(), time.Now().Add(time.Second*60))
	if(err != nil){
		Error := ErrorEvent{
			EventType       : "ErrorEvent",
			Guid            : Guid.String(),
			OccuredAt       : time.Now(),
			TransactionId   : TransId,
			UserId          : UserId,
			Service         : "API",
			Server          : Hostname,
			Command         : "BUY",
			StockSymbol     : "",
			Funds           : t.Amount,
			FileName        : "",
			ErrorMessage    : "Failed to create purchase",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//writeResponse(w, http.StatusInternalServerError, "Failed to Create Purchase")
	    	return
	}
	//success
	//writeResponse(w, http.StatusOK, "Purchase Request has been Created")
	return
}

func CommitBuy(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Commiting Buy Request");
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	db := getDatabasePointerForUser(UserId)

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
		Command         : "COMMIT_BUY",
		StockSymbol     : "",
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType); 


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
			Command         : "COMMIT_BUY",
			StockSymbol     : "",
			Funds           : "",
			FileName        : "",
			ErrorMessage    : "User Account Does Not Exist",   
		}
		SendRabbitMessage(Error,Error.EventType) 
		//error
		//writeResponse(w, http.StatusBadRequest, "User Account Does Not Exist")
		return       
	}

	LatestPendingrows, err := db.Query(getLatestPendingPurchase, uid)
	if err != nil{
		//error
		//writeResponse(w, http.StatusBadRequest, "No Recent BUY Commands to Commit")
		return
	}
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
			Command         : "COMMIT_BUY",
			StockSymbol     : "",
			Funds           : "",
			FileName        : "",
			ErrorMessage    : "No recent BUY commands issued",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//error
		return                  
	}
	if expires_at.Before(time.Now()){
		//success (Kinda)
		//writeResponse(w, http.StatusBadRequest, "Purchase Request has Timed Out")
		_, err = db.Exec(cancelTransaction, id)
		if err != nil{
			//error
			return
		}
		return
	}
	_, err = db.Exec(commitPurchase, id, time.Now())
	if err != nil {
		Error := ErrorEvent{
			EventType       : "ErrorEvent",
			Guid            : Guid.String(),
			OccuredAt       : time.Now(),
			TransactionId   : TransId,
			UserId          : UserId,
			Service         : "API",
			Server          : Hostname,
			Command         : "COMMIT_BUY",
			StockSymbol     : stock,
			Funds           : "",
			FileName        : "",
			ErrorMessage    : "Not Enough funds to buy Stocks",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//error
		return
	}
	//success
	//writeResponse(w, http.StatusOK, "Purchase Request has been Commited")
	return
}

func CancelBuy(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Cancelling Buy Request");
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}
	db := getDatabasePointerForUser(UserId)

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
		Command         : "CANCEL_BUY",
		StockSymbol     : "",
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);

	//Get Database User
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
			Command         : "CANCEL_BUY",
			StockSymbol     : "",
			Funds           : "",
			FileName        : "",
			ErrorMessage    : "User Account Does Not Exist",   
		}
		SendRabbitMessage(Error,Error.EventType)    
		//error
		//writeResponse(w, http.StatusBadRequest, "User Account Does Not Exist")
		return    
	}

	LatestPendingrows, err := db.Query(getLatestPendingPurchase, uid)
	defer LatestPendingrows.Close()
	if err != nil {
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
			StockSymbol     : stock,
			Funds           : "",
			FileName        : "",
			ErrorMessage    : "No recent BUY commands issued",   
		}
		SendRabbitMessage(Error,Error.EventType)       
		//error
		//writeResponse(w, http.StatusBadRequest, "No Recent BUY commands to be Cancelled")
		return           
	}

	_, err = db.Exec(cancelTransaction, id)
	if err != nil {
		//error
		//writeResponse(w, http.StatusInternalServerError, "Failed To Cancel Transaction")
		return
	}	
	//writeResponse(w, http.StatusOK, "Purchase Request has been Cancelled")
	return 
}
