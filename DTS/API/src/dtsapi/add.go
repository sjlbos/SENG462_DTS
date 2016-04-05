package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

func Add(w http.ResponseWriter, r *http.Request){
	zero,_ := decimal.NewFromString("0");

	type add_struct struct {
		Amount string
	}

	type return_struct struct {
		Error bool
		Amount string
		UserId string
	}
	vars := mux.Vars(r)
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" { 
		TransId = "0"
	}
	decoder := json.NewDecoder(r.Body)
	var t add_struct   
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
		Command         : "ADD",
		StockSymbol     : "",
		Funds           : t.Amount,
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType)
	if err != nil{
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	AmountDec,err := decimal.NewFromString(t.Amount)
	if err != nil {	
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//Get user id from Database
	db := getDatabasePointerForUser(UserId)

	//Amount to add is invalid
	if(AmountDec.Cmp(zero) != 1){
	    writeResponse(w, http.StatusBadRequest, "Amount to add is not a valid number")
	    return
	}

	_, err = db.Exec(addOrCreateUser, UserId, t.Amount, time.Now())
	//Failed to Create Account
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//Build Response
	rtnStruct := return_struct{false, t.Amount, UserId} 
	strRtnStruct, err := json.Marshal(rtnStruct)

	//Success
	writeResponse(w, http.StatusOK, string(strRtnStruct))
	return
}
