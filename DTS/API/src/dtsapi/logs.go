package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gorilla/mux"
)

func DisplaySummary(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Display Summary for User: ")

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
		Server          : "B134",
		Command         : "DISPLAY_SUMMARY",
		StockSymbol     : "",
		Funds           : "",
	}
	SendRabbitMessage(CommandEvent,CommandEvent.EventType);

	//Get and Return User Information
	db, uid, found, balanceStr := getDatabaseUserId(UserId) 
	if !found {
		fmt.Fprintln(w, "User Does Not Exist")
		return
	}
	fmt.Fprintln(w, UserId + string(uid) + balanceStr)
	
	//Get And Return Stocks
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Stocks:")
	rows, err := db.Query(getAllStocks, uid)
	if err != nil{
		return
	}
	defer rows.Close()
	var stock string
	var num_shares string

	for rows.Next() {
		err = rows.Scan(&stock, &num_shares)
		if err !=nil {
			return
		}
		fmt.Fprintln(w, stock + num_shares)
	}

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Triggers:")
	rows, err = db.Query(getAllTriggers, uid)
	if err != nil{
		return
	}
	defer rows.Close()
	var id string
	var uidStr string
	var trigger_type string
	var price string
	var created_time time.Time

	for rows.Next() {
		err = rows.Scan(&id, &uidStr, &stock, &trigger_type, &price, &num_shares, &created_time)
		if err != nil{
			return
		}
		fmt.Fprintln(w, id + uidStr + stock + trigger_type + price + num_shares + created_time.String())
	}
	return
}




