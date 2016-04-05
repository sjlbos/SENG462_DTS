package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"fmt"
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
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//get User Account Information
	db, id, found, _ := getDatabaseUserId(UserId) 
	if(found == false){
		writeResponse(w, http.StatusOK, "User Account Does Not Exist")
		return
	}

	//Check amount
	AmountDec,err := decimal.NewFromString(t.Amount)
	if err != nil{
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if(AmountDec.Cmp(zero) != 1){
	    writeResponse(w, http.StatusBadRequest, "Amount to sell is not a valid number")
	    return
	}

	//Check Stock Symbol
	StockId := t.Symbol
	if(len(StockId) == 0 || len(StockId) > 3){
	    writeResponse(w, http.StatusBadRequest, "Symbol is Not Valid")
	    return
	}

	//Get and Verify Quote
	var strPrice string
	strPrice = getStockPrice(TransId ,"true", UserId, StockId, Guid.String())
	var quotePrice decimal.Decimal
	quotePrice, err = decimal.NewFromString(strPrice)
	if err != nil{
		writeResponse(w, http.StatusBadRequest, err.Error())
		return;
	}
	if(quotePrice.Cmp(zero) != 1){
	    writeResponse(w, http.StatusBadRequest, "Quote is not a valid number")
	    return
	}

	//Calculate Amount to Sell
	toSell := (AmountDec.Div(quotePrice)).Floor()
	if toSell.Cmp(zero) != 1 {
		writeResponse(w, http.StatusOK, "Can't Sell less than 1 Stock")
		return
	}

	//Create Pending Sale
	_,err = db.Exec(addPendingSale, id, t.Symbol, toSell.String(), strPrice, time.Now(), time.Now().Add(time.Second*60))
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "Add pending Sale; " + err.Error())
		return
	}

	//success
	writeResponse(w, http.StatusOK, "Sale Request Has Been Created")
	return    
}

func CommitSell(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Commiting Sell Request");
	vars := mux.Vars(r)
	UserId := vars["id"]
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
		Command         : "COMMIT_SELL",
		StockSymbol     : "",
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);

	//Find User and Database Information
	db, uid, found, _ := getDatabaseUserId(UserId) 
	if(found == false){
		writeResponse(w, http.StatusOK, "User Does not Exist")
		return       
	}

	//Get Last Pending Sale
	LatestPendingrows, err := db.Query(getLatestPendingSale, uid)
	defer LatestPendingrows.Close()
	if err != nil {
		writeResponse(w, http.StatusOK, "Error Getting Recent Transaction: " + err.Error())
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

	//Verify Pending Sale
	if !found {
		writeResponse(w, http.StatusOK, "No Recent Sell Commands")
		return                 
	}
	if expires_at.Before(time.Now()){
		_, err = db.Exec(cancelTransaction, id)
		if err != nil{
			//error
			writeResponse(w, http.StatusOK, "Error Cancelling Expired Transaction: " + err.Error())
			return
		}
		writeResponse(w, http.StatusOK, "Transaction Has Expired")
		return
	}
	_, err = db.Exec(commitSale, id, time.Now())
	if(err != nil){
		writeResponse(w, http.StatusOK, "Error Commiting: " + err.Error())
		return
	}

	//success
	writeResponse(w, http.StatusOK, "Sale Request Has Been Commited")
	return
}

func CancelSell(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Cancelling Sell Request");
	vars := mux.Vars(r)
	UserId := vars["id"]
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
		Command         : "CANCEL_SELL",
		StockSymbol     : "",
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);

	//Find user in DB
	db, uid, found,_ := getDatabaseUserId(UserId) 
	if(found == false){
		writeResponse(w, http.StatusOK, "User Does Not Exist")
		return        
	}

	//Find last Sell Command
	LatestPendingrows, err := db.Query(getLatestPendingSale, uid)
	defer LatestPendingrows.Close()
	if err != nil{
		writeResponse(w, http.StatusOK, "Error Getting Last Sale: " + err.Error())
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
		writeResponse(w, http.StatusOK, "No Recent Sell Commands")  
		return               
	}
	_, err = db.Exec(cancelTransaction, id)
	if err != nil{
		writeResponse(w, http.StatusOK, "Error Cancelling Sale: " + err.Error())
		return
	}
	//success
	writeResponse(w, http.StatusOK, "Sale Request Has Been Cancelled")
	return
}
