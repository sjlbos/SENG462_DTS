package main

import (
	"encoding/json"
	"net/http"
	//"strings"
	"time"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	//"fmt"
)

func Add(w http.ResponseWriter, r *http.Request){
	zero,_ := decimal.NewFromString("0");

	type add_struct struct {
		Amount string
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

	//Success
	writeResponse(w, http.StatusOK, "Account Updated with Funds")
	return
}
