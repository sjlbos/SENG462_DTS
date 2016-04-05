package main

import (
    //"fmt"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

func Quote(w http.ResponseWriter, r *http.Request){
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
		writeResponse(w, http.StatusBadRequest, "Symbol is Not Valid")
		return
	}

	//Get Stock Price
	var strPrice string
	strPrice = getStockPrice(TransId ,"false", UserId, StockId, Guid.String())

	//Verify Return Price
	var price decimal.Decimal
	price, err := decimal.NewFromString(strPrice)
	if err != nil{
		writeResponse(w, http.StatusBadRequest, "Quote Return: " + err.Error())
		return
	}

	//Success
	var Output string = "The Quote For UserId " + UserId + " and StockId " + StockId + " returned " + price.String()
	writeResponse(w, http.StatusOK, Output)
}
