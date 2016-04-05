package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)


func Sell(w http.ResponseWriter, r *http.Request){
	zero,_ := decimal.NewFromString("0");
	type sell_struct struct {
		Amount string
		Symbol string
	}

	type return_struct struct {
		Error bool
		SaleId int
		Price string
		NumShares int64
		Expiration time.Duration
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
	var strExpiration string
	strPrice, strExpiration = getStockPrice(TransId ,"true", UserId, StockId, Guid.String())
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
	//Verify Expiration Time
	ExpirationTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", strExpiration)
	if err != nil{
		writeResponse(w, http.StatusOK, "Expiration Conversion Error")
		return
	}

	//Calculate Amount to Sell
	toSell := (AmountDec.Div(quotePrice)).IntPart()
	if toSell < 1 {
		writeResponse(w, http.StatusOK, "Can't Sell less than 1 Stock")
		return
	}

	strSell := strconv.Itoa(int(toSell))

	//Create Pending Sale
	var SaleId int
	rows ,err := db.Query(addPendingSale, id, t.Symbol, strSell, strPrice, time.Now(), ExpirationTime)
	defer rows.Close()
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "Add pending Sale; " + err.Error())
		return
	}

	TimeToExpiration := (ExpirationTime.Sub(time.Now()))/time.Millisecond
	rows.Next()
	err = rows.Scan(&SaleId)

	//Build Response
	rtnStruct := return_struct{false, SaleId, strPrice, toSell, TimeToExpiration} 
	strRtnStruct, err := json.Marshal(rtnStruct)

	//success
	writeResponse(w, http.StatusOK, string(strRtnStruct))
	return    
}

func CommitSell(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	type return_struct struct {
		Error bool
		SaleId int
		Price string
		Stock string
		NumShares int
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

	//Build Response
	rtnStruct := return_struct{false, id, stock, share_price, num_shares} 
	strRtnStruct, err := json.Marshal(rtnStruct)

	//success
	writeResponse(w, http.StatusOK, string(strRtnStruct))
	return
}

func CancelSell(w http.ResponseWriter, r *http.Request){
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


func UpdateSale(w http.ResponseWriter, r *http.Request){
	zero,_ := decimal.NewFromString("0");
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := vars["PurchaseId"]


	//get User Account Information
	db, uid, found, _ := getDatabaseUserId(UserId) 
	if(found == false){
		writeResponse(w, http.StatusOK, "User Account Does Not Exist")
		return
	}

	Guid := getNewGuid()
	//Find last Sell Command
	LatestPendingrows, err := db.Query(getLatestPendingSale, uid)
	defer LatestPendingrows.Close()
	if err != nil{
		writeResponse(w, http.StatusOK, "Error Getting Last Sale: " + err.Error())
		return
	}

	var id string
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

	strOldPrice := strings.TrimPrefix(share_price, "$")
	strOldPrice = strings.Replace(strOldPrice, ",", "", -1)
	OldPrice, err := decimal.NewFromString(strOldPrice)
	if err != nil{
		writeResponse(w, http.StatusBadRequest, err.Error() + strOldPrice)
		return;
	}

	//Get and Verify Quote
	var strPrice string
	strPrice, _ = getStockPrice(TransId ,"true", UserId, stock, Guid.String())
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

	totalAmount := decimal.New(int64(num_shares),0).Mul(OldPrice)
	newShareNum := totalAmount.Div(quotePrice).IntPart()
	diffShares := int(newShareNum) - num_shares
	_, err = db.Exec(updateSale, TransId, int(newShareNum), strPrice, int(diffShares), time.Now().Add(time.Duration(60)*time.Second))
	if err != nil{
	    writeResponse(w, http.StatusBadRequest, "Unable to update Sale")
	    return
	}
			
	type return_struct struct {
		Error bool
		SaleId string
		Price string
		NumShares int64
		Expiration time.Duration
	}
	//Build Response
	rtnStruct := return_struct{false, id, strPrice, newShareNum, -1} 
	strRtnStruct, err := json.Marshal(rtnStruct)
	writeResponse(w, http.StatusOK, string(strRtnStruct))  
	return
}