package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"fmt"
)

func Add(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Adding Funds Too Account");
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
		//error
		return
	}

	AmountDec,err := decimal.NewFromString(t.Amount)
	if err != nil {	
		//error
		return
	}

	//Get user id from Databasse
	db := getDatabasePointerForUser(UserId)

	var balanceStr string
	var balance decimal.Decimal

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
	    writeResponse(w, http.StatusBadRequest, "Amount to add is not a valid number")
	    return
	}

	id, found, balanceStr := getDatabaseUserId(UserId)
	//User Account Does not Exist, Create Account 
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
			Funds           : t.Amount,
			FileName        : "",
			DebugMessage    : "Created User Account",   
		}
		SendRabbitMessage(Debug,Debug.EventType)
		_, err := db.Exec(addUser, UserId, t.Amount, time.Now())
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
				StockSymbol     : "",
				Funds           : t.Amount,
				FileName        : "",
				ErrorMessage    : "Failed To Create User",   
			}
			SendRabbitMessage(Error,Error.EventType)
			writeResponse(w, http.StatusInternalServerError, "Failed To Create User Account")
			return
		}
		//success
		writeResponse(w, http.StatusOK, "Account Created with Funds")
		return
	}

	//User Account Exists, Add Funds
	balanceStr = strings.Trim(balanceStr, "$")
	balanceStr = strings.Replace(balanceStr, ",", "", -1)
	balance, err = decimal.NewFromString(balanceStr)
	newBalance := balance.Add(AmountDec)
	AccountEvent := AccountTransactionEvent{
		EventType       : "AccountTransactionEvent",
		Guid            : Guid.String(),
		OccuredAt       : time.Now(),
		TransactionId   : TransId,
		UserId          : UserId,
		Service         : "Account",
		Server          : Hostname,
		AccountAction   : "Add",
		Funds           : t.Amount,
	}
	SendRabbitMessage(AccountEvent,AccountEvent.EventType)
	_, err = db.Exec(updateBalance, id, newBalance)
	if err != nil {
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
			ErrorMessage    : "Failed to add Funds to Account",   
		}
		SendRabbitMessage(Error,Error.EventType)
		writeResponse(w, http.StatusInternalServerError, "Failed To Add funds to Account")
		return
	}
	//success
	writeResponse(w, http.StatusOK, "Funds Have been Added to Account")
	return
}
