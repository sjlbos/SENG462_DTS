package main

import (
    "encoding/json"
    "net/http"
    "time"
    "strings"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
   	"fmt"
)

func CreateBuyTrigger(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Creating Buy Trigger");
	type trigger_struct struct{
		Amount string
		Price  string
	}

	vars := mux.Vars(r)
	UserId := vars["id"]
	Symbol := vars["symbol"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	db := getDatabasePointerForUser(UserId)

	decoder := json.NewDecoder(r.Body)
	var t trigger_struct   
	err = decoder.Decode(&t)

	//Audit UserCommand
	Guid := getNewGuid()
	OccuredAt := time.Now()

	if t.Amount == "" && t.Price == ""{
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
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)
		//error
		return
	}

	if t.Amount != "" {
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
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)
		if err != nil {
			//error
			return
		}

		uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			//error
			return
		}

		_, err := decimal.NewFromString(t.Amount)
		if err != nil{
			//error
			return
		}

		_, err = db.Exec(addBuyTrigger, uid, Symbol, t.Amount, time.Now())
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
				Funds           : t.Amount,
				FileName        : "",
				ErrorMessage    : "Unable to Create Buy Trigger",   
			}
			SendRabbitMessage(Error,Error.EventType)
			//error
			return
		}
		//success
		//writeResponse(w, http.StatusOK, "Buy Trigger Created")
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

		if err != nil {
			//error
			return
		}

		uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			//error
			return
		}

		AmountDec, err := decimal.NewFromString(t.Price)
		if err != nil{
			//error
			return
		}

		getPendingTriggerRows, err := db.Query(getPendingTriggerId, uid, Symbol, "buy")
		if err != nil{
			//error
			return
		}
		defer getPendingTriggerRows.Close()

		var id int
		found = false

		for getPendingTriggerRows.Next() {
			found = true
			err = getPendingTriggerRows.Scan(&id)
		} 
		if !found {
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
				Funds           : AmountDec.String(),
				FileName        : "",
				ErrorMessage    : "No recent SET_BUY_AMOUNT commands issued",   
			}
			SendRabbitMessage(Error,Error.EventType)    
			//error
			return               
		}
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
		//success
		//writeResponse(w, http.StatusOK, "Buy Trigger Set")
		return
	}
}

func CreateSellTrigger(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Creating Sell Trigger");
	type trigger_struct struct{
		Amount string
		Price  string
	}

	vars := mux.Vars(r)
	UserId := vars["id"]
	Symbol := vars["symbol"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	decoder := json.NewDecoder(r.Body)
	var t trigger_struct   
	err := decoder.Decode(&t)

	//Audit UserCommand
	Guid := getNewGuid()
	OccuredAt := time.Now()

	db := getDatabasePointerForUser(UserId)

	if t.Amount == "" && t.Price == ""{
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
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)
		//error
		return
	}


	if(t.Amount != ""){
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
		if err != nil {
			//error
			return
		}

		uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			//error
			return
		}

		_, err := decimal.NewFromString(t.Amount)
		if err != nil{
			//error
			return
		}

		_, err = db.Exec(addSellTrigger, uid, Symbol, t.Amount, time.Now())
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
				Funds           : t.Amount,
				FileName        : "",
				ErrorMessage    : "Unable to Create Buy Trigger",   
			}
			SendRabbitMessage(Error,Error.EventType)
			//error
			return
		}
		//success
		//writeResponse(w, http.StatusOK, "Sell Trigger Created")
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
			Command         : "SET_SELL_TRIGGER",
			StockSymbol     : Symbol,
			Funds           : "",
		}
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)

		if err != nil {
			//error
			return
		}

		uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			//error
			return
		}

		_, err := decimal.NewFromString(t.Price)
		if err != nil{
			//error
			return
		}

		getPendingTriggerRows, err := db.Query(getPendingTriggerId, uid, Symbol, "sell")
		if err != nil{
			//error
			return
		}
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
				Funds           : t.Amount,
				FileName        : "",
				ErrorMessage    : "No recent SET_SELL_AMOUNT commands issued",   
			}
			SendRabbitMessage(Error,Error.EventType)     
			//error
			return             
		}
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
		//success
		//writeResponse(w, http.StatusOK, "Sell Trigger Set")
		return
	}
}

func CancelBuyTrigger(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Cancelling Buy Trigger");
	vars := mux.Vars(r)
	UserId := vars["id"]
	Symbol := vars["symbol"]    
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
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
		Server          : "B134",
		Command         : "CANCEL_SET_BUY",
		StockSymbol     : Symbol,
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);

	uid, found, _ := getDatabaseUserId(UserId)

	if !found {
		//error
		return
	}

	db := getDatabasePointerForUser(UserId)
	var id int
	getBuyTriggerIdRows, err := db.Query(getBuyTriggerId, uid, Symbol)
	if err != nil{
		//error
		return;
	}
	defer getBuyTriggerIdRows.Close()
	found = false
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
	//success
	//writeResponse(w, http.StatusOK, "Buy Trigger Cancelled")
	return
}

func CancelSellTrigger(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Cancelling Sell Request");
	vars := mux.Vars(r)
	UserId := vars["id"]
	Symbol := vars["symbol"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
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
		Server          : "B134",
		Command         : "CANCEL_SET_SELL",
		StockSymbol     : Symbol,
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);
	uid, found, _ := getDatabaseUserId(UserId)

	if !found {
		//error
		return
	}

	db := getDatabasePointerForUser(UserId)
	var id int
	getSellTriggerIdRows, err := db.Query(getSellTriggerId, uid, Symbol)
	if err != nil{
		//error
		return
	}
	defer getSellTriggerIdRows.Close()
	found = false

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
	//success
	//writeResponse(w, http.StatusOK, "Sell Trigger Cancelled")
	return
}

func PerformSellTrigger(w http.ResponseWriter, r *http.Request){
	zero,_ := decimal.NewFromString("0");
	vars := mux.Vars(r)
	UserId := vars["id"]
	StockId := vars["symbol"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	type trigger_struct struct {
		TriggerId string
	}

	Guid := getNewGuid()
	decoder := json.NewDecoder(r.Body)
	var t trigger_struct   
	err := decoder.Decode(&t)
	if err != nil{

	}
	
	db := getDatabasePointerForUser(UserId)

	//Check If User Exists
	uid, found, _ := getDatabaseUserId(UserId) 
	
	if !found {
		//error
		return
	}

	//Get A Quote
	var strPrice string
	strPrice = getStockPrice(TransId ,"false", UserId, StockId, Guid.String())
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
			Funds           : strPrice,
			FileName        : "",
			ErrorMessage    : "Quote is not greater than 0",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//writeResponse(w, http.StatusBadRequest, "Amount to buy is not a valid number")
		return
	}

	//Get Trigger Information
	rows, err := db.Query(getTriggerById, t.TriggerId)
	if err != nil {
		//error
		return
	}
	defer rows.Close()

	var id int
	var stock string
	var trigger_type string 
	var trigger_price string 
	var num_shares int
	var created_at time.Time
	found = false
	for rows.Next() {
		found = true
		err = rows.Scan(&id, &uid, &stock, &trigger_type, &trigger_price, &num_shares, &created_at)
	}	
	
	if !found {
		//error
		return
	}

	trigger_price = strings.Trim(trigger_price, "$")
	trigger_price = strings.Replace(trigger_price, ",", "", -1)
	trigger_priceDec, err := decimal.NewFromString(trigger_price)

	if trigger_priceDec.Cmp(quotePrice) != -1{
		//error
		return
	}
	
	//Commit trigger at price
	_,err = db.Exec(performSellTrigger, id, strPrice)
	if err != nil{
		//error
		return
	}	
}

func PerformBuyTrigger(w http.ResponseWriter, r *http.Request){
	zero,_ := decimal.NewFromString("0");
	vars := mux.Vars(r)
	UserId := vars["id"]
	StockId := vars["symbol"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	type trigger_struct struct {
		TriggerId string
	}
	decoder := json.NewDecoder(r.Body)
	var t trigger_struct   
	err := decoder.Decode(&t)
	if err != nil{

	}
	Guid := getNewGuid()
	db := getDatabasePointerForUser(UserId)

	//Check If User Exists
	uid, found, _ := getDatabaseUserId(UserId) 
	if !found {
		//error
		return
	}

	//Get A Quote
	var strPrice string
	strPrice = getStockPrice(TransId ,"false", UserId, StockId, Guid.String())
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
			Funds           : strPrice,
			FileName        : "",
			ErrorMessage    : "Quote is not greater than 0",   
		}
		SendRabbitMessage(Error,Error.EventType)
		//writeResponse(w, http.StatusBadRequest, "Amount to buy is not a valid number")
		return
	}

	//Get Trigger Information
	rows, err := db.Query(getTriggerById, t.TriggerId)
	if err != nil {
		//error
		return
	}
	defer rows.Close()

	var id int 
	var stock string
	var trigger_type string 
	var trigger_price string 
	var num_shares int
	var created_at time.Time
	found = false
	for rows.Next() {
		found = true
		err = rows.Scan(&id, &uid, &stock, &trigger_type, &trigger_price, &num_shares, &created_at)
	}	
	if err != nil {
		//error 
		return
	}
	
	if !found {
		//error
		return
	}

	trigger_price = strings.Trim(trigger_price, "$")
	trigger_price = strings.Replace(trigger_price, ",", "", -1)
	trigger_priceDec, err := decimal.NewFromString(trigger_price)

	if quotePrice.Cmp(trigger_priceDec) != -1{
		//error
		return
	}
	
	//Commit trigger at price
	_,err = db.Exec(performBuyTrigger, id, strPrice)
	if err != nil{
		//error
		return
	}
}
