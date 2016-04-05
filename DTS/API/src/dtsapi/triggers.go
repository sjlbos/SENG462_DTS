package main

import (
    "encoding/json"
    "net/http"
    "time"
    "strings"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

func CreateBuyTrigger(w http.ResponseWriter, r *http.Request){
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
			Server          : Hostname,
			Command         : "SET_BUY_AMOUNT",
			StockSymbol     : Symbol,
			Funds           : "",
		}
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)
		writeResponse(w, http.StatusOK, "Request Body Is Invalid")
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
			Server          : Hostname,
			Command         : "SET_BUY_AMOUNT",
			StockSymbol     : Symbol,
			Funds           : "",
		}
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)
	
		//Validate Request Body
		if err != nil {
			writeResponse(w, http.StatusBadRequest, "Request Body Is Invalid")
			return
		}

		//Get User and Database Information
		db, uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			writeResponse(w, http.StatusOK, "User Account Does Not Exist")
			return
		}

		//Validate Amount
		_, err := decimal.NewFromString(t.Amount)
		if err != nil{
			//error
			writeResponse(w, http.StatusBadRequest, "Amount Is Not A Valid Number")
			return
		}

		//Create Buy Trigger
		_, err = db.Exec(addBuyTrigger, uid, Symbol, t.Amount, time.Now())
		if err != nil{
			writeResponse(w, http.StatusBadRequest, "addBuyTrigger: " + err.Error())
			return
		}

		//success
		writeResponse(w, http.StatusOK, "Buy Trigger Created")
		return
	}else{
		CommandEvent := UserCommandEvent{
			EventType       : "UserCommandEvent",
			Guid            : Guid.String(),
			OccuredAt       : OccuredAt,
			TransactionId   : TransId,
			UserId          : UserId,
			Service         : "Command",
			Server          : Hostname,
			Command         : "SET_BUY_TRIGGER",
			StockSymbol     : Symbol,
			Funds           : "",
		}
		SendRabbitMessage(CommandEvent,CommandEvent.EventType);
	
		//Validate Request Body
		if err != nil {
			writeResponse(w, http.StatusBadRequest, "Request Body Is Invalid")
			return
		}

		//Get User and Database Information
		db, uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			writeResponse(w, http.StatusBadRequest, "User Account Does Not Exist")
			return
		}

		//Validate Price
		_, err := decimal.NewFromString(t.Price)
		if err != nil{
			writeResponse(w, http.StatusBadRequest, "Amount Is Not A Valid Number")
			return
		}

		//Get Created Trigger
		getPendingTriggerRows, err := db.Query(getPendingTriggerId, uid, Symbol, "buy")
		defer getPendingTriggerRows.Close()
		if err != nil{
			writeResponse(w, http.StatusBadRequest, "getPendingTriggerId: " + err.Error())
			return
		}

		//Get Trigger Info
		var id int
		found = false
		for getPendingTriggerRows.Next() {
			found = true
			err = getPendingTriggerRows.Scan(&id)
		} 
		if !found {
			writeResponse(w, http.StatusBadRequest, "No Recent SET_BUY_AMOUNT commands issued")
			return               
		}

		//Create Buy Trigger
		_, err = db.Exec(setBuyTrigger, id, uid, t.Price, time.Now())
		if err != nil {
			writeResponse(w, http.StatusBadRequest, "Unable to set Buy triggger: " + err.Error())
			return
		}

		//Send Trigger Message
		Trigger := TriggerEvent{
			TriggerType     : "buy",
			UserId          : UserId,
			TransactionId   : TransId,
			UpdatedAt       : time.Now(),
		}
		SendRabbitMessage(Trigger, "buy")

		//success
		writeResponse(w, http.StatusOK, "Buy Trigger Set")
		return
	}
}

func CreateSellTrigger(w http.ResponseWriter, r *http.Request){
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

	if t.Amount == "" && t.Price == ""{
		CommandEvent := UserCommandEvent{
			EventType       : "UserCommandEvent",
			Guid            : Guid.String(),
			OccuredAt       : OccuredAt,
			TransactionId   : TransId,
			UserId          : UserId,
			Service         : "Command",
			Server          : Hostname,
			Command         : "SET_BUY_AMOUNT",
			StockSymbol     : Symbol,
			Funds           : "",
		}
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)
		writeResponse(w, http.StatusBadRequest, "Invalid Request Body")
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
			Server          : Hostname,
			Command         : "SET_SELL_AMOUNT",
			StockSymbol     : Symbol,
			Funds           : "",
		}
		SendRabbitMessage(CommandEvent,CommandEvent.EventType);
	
		if err != nil {
			writeResponse(w, http.StatusBadRequest, "Invalid Request Body")
			return
		}

		db, uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			writeResponse(w, http.StatusBadRequest, "User Account Does Not Exist")
			return
		}

		_, err := decimal.NewFromString(t.Amount)
		if err != nil{
			writeResponse(w, http.StatusBadRequest, "Amount Is Not A Valid Number")
			return
		}

		_, err = db.Exec(addSellTrigger, uid, Symbol, t.Amount, time.Now())
		if err != nil{
			writeResponse(w, http.StatusBadRequest, "addSellTrigger: " + err.Error())
			return
		}

		//success
		writeResponse(w, http.StatusOK, "Sell Trigger Created")
		return	
	}else{
		CommandEvent := UserCommandEvent{
			EventType       : "UserCommandEvent",
			Guid            : Guid.String(),
			OccuredAt       : OccuredAt,
			TransactionId   : TransId,
			UserId          : UserId,
			Service         : "Command",
			Server          : Hostname,
			Command         : "SET_SELL_TRIGGER",
			StockSymbol     : Symbol,
			Funds           : "",
		}
		SendRabbitMessage(CommandEvent,CommandEvent.EventType)
	
		if err != nil {
			writeResponse(w, http.StatusBadRequest, "Invalid Request Body")
			return
		}

		db, uid, found, _ := getDatabaseUserId(UserId)
		if !found {
			writeResponse(w, http.StatusBadRequest, "User Account Not Found")
			return
		}

		_, err := decimal.NewFromString(t.Price)
		if err != nil{
			writeResponse(w, http.StatusBadRequest, "Price Is Not A Valid Number")
			return
		}

		getPendingTriggerRows, err := db.Query(getPendingTriggerId, uid, Symbol, "sell")
		defer getPendingTriggerRows.Close()
		if err != nil{
			writeResponse(w, http.StatusBadRequest, "getPendingTriggerId: " + err.Error())
			return
		}

		var id int
		found = false
		for getPendingTriggerRows.Next() {
			found = true
			err = getPendingTriggerRows.Scan(&id)
		} 
		if(found == false){
			writeResponse(w, http.StatusBadRequest, "No recent SET_SELL_AMOUNT commands Issued")
			return             
		}

		_, err = db.Exec(setSellTrigger, id, uid, t.Price, time.Now())
		if err != nil {
			writeResponse(w, http.StatusBadRequest, "setSellTrigger: " + err.Error())
			return
		}

		//Send Trigger Message
		Trigger := TriggerEvent{
			TriggerType     : "sell",
			UserId          : UserId,
			TransactionId   : TransId,
			UpdatedAt       : time.Now(),
		}
		SendRabbitMessage(Trigger, "sell")

		//success
		writeResponse(w, http.StatusOK, "Sell Trigger Set")
		return
	}
}

func CancelBuyTrigger(w http.ResponseWriter, r *http.Request){
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
		Server          : Hostname,
		Command         : "CANCEL_SET_BUY",
		StockSymbol     : Symbol,
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);
	
	//Get User and Database Information
	db , uid, found, _ := getDatabaseUserId(UserId)
	if !found {
		writeResponse(w, http.StatusOK, "User Does Not Exist")
		return
	}

	//Get Buy Trigger to Cancel
	var id int
	getBuyTriggerIdRows, err := db.Query(getBuyTriggerId, uid, Symbol)
	defer getBuyTriggerIdRows.Close()
	if err != nil{
		writeResponse(w, http.StatusOK, "Failed To Get Trigger: " + err.Error())
		return;
	}
	found = false
	for getBuyTriggerIdRows.Next() {
		found = true
		err = getBuyTriggerIdRows.Scan(&id)
	}

	//Cancel Buy Trigger
	_, err = db.Exec(cancelBuyTrigger, id)
	if err != nil{
		writeResponse(w, http.StatusOK, "Failed To Cancel: " + err.Error())
		return
	}

	//Send Trigger Event
	Trigger := TriggerEvent{
		TriggerType     : "buy",
		UserId          : UserId,
		TransactionId   : TransId,
		UpdatedAt       : time.Now(),
	}
	SendRabbitMessage(Trigger, "buy")

	//success
	writeResponse(w, http.StatusOK, "Buy Trigger Cancelled")
	return
}

func CancelSellTrigger(w http.ResponseWriter, r *http.Request){
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
		Server          : Hostname,
		Command         : "CANCEL_SET_SELL",
		StockSymbol     : Symbol,
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);
	
	//Get Database and User Information
	db, uid, found, _ := getDatabaseUserId(UserId)
	if !found {
		writeResponse(w, http.StatusOK, "User Does Not Exist")
		return
	}

	//Get Sell Trigger to Cancel
	var id int
	getSellTriggerIdRows, err := db.Query(getSellTriggerId, uid, Symbol)
	defer getSellTriggerIdRows.Close()
	if err != nil{
		writeResponse(w, http.StatusOK, "Failed To Cancel: " + err.Error())
		return
	}
	found = false
	for getSellTriggerIdRows.Next() {
		found = true
		err = getSellTriggerIdRows.Scan(&id)
	}

	//Cancel Sell Trigger
	_, err = db.Exec(cancelSellTrigger, id)
	if err != nil{
		writeResponse(w, http.StatusOK, "Failed To Cancel: " + err.Error())
		return
	}

	//Send Trigger update
	Trigger := TriggerEvent{
		TriggerType     : "buy",
		UserId          : UserId,
		TransactionId   : TransId,
		UpdatedAt       : time.Now(),
	}
	SendRabbitMessage(Trigger, "buy")

	//success
	writeResponse(w, http.StatusOK, "Sell Trigger Cancelled")
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

	//Check If User Exists
	db, uid, found, _ := getDatabaseUserId(UserId) 
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
	defer rows.Close()
	if err != nil {
		//error
		return
	}

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

	//Check If User Exists
	db, uid, found, _ := getDatabaseUserId(UserId) 
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
	defer rows.Close()
	if err != nil {
		//error
		return
	}

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
