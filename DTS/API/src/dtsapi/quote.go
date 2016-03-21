package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

func Quote(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Getting Quote");
	vars := mux.Vars(r)
	StockId := vars["symbol"]
	UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}

	//Audit UserCommand
	Guid := getNewGuid()
	CommandEvent := UserCommandEvent{
		EventType       : "UserCommandEvent",
		Guid            : Guid.String(),
		OccuredAt       : time.Now(),
		TransactionId   : TransId,
		UserId          : UserId,
		Service         : "Command",
		Server          : "B134",
		Command         : "QUOTE",
		StockSymbol     : StockId,
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);

	//Check Stock Symbol
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
			StockSymbol     : StockId,
			Funds           : "",
			FileName        : "",
			ErrorMessage    : "Symbol is Not Valid",   
		}
		SendRabbitMessage(Error,Error.EventType)
		writeResponse(w, http.StatusBadRequest, "Symbol is Not Valid")
		return
	}

	var strPrice string
	strPrice = getStockPrice(TransId ,"false", UserId, StockId, Guid.String())


	var price decimal.Decimal
	price, err := decimal.NewFromString(strPrice)
	if err != nil{
		//error
		return
	}
	var Output string = "The Quote For UserId " + UserId + " and StockId " + StockId + " returned " + price.String()
	fmt.Fprintln(w, Output)
}
