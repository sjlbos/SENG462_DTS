package main

import (
//    "encoding/json"
    "fmt"
//    "net"
    "net/http"
//    "os"
//    "log"
//    "strings"
//    "strconv"
    "time"
//    "io/ioutil"

    "github.com/gorilla/mux"
//    "github.com/streadway/amqp"
//   "github.com/nu7hatch/gouuid"
)


func DumplogUser(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Display Log for User: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
	TransId := r.Header.Get("X-TransNo")
	if TransId == "" {
		TransId = "0"
	}
    fmt.Fprintln(w, UserId)

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
        Command         : "DUMPLOG",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);
}


func Dumplog(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Dumplog of all transactions: ")
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
        UserId          : "",
        Service         : "Command",
        Server          : "B134",
        Command         : "DUMPLOG",
        StockSymbol     : "",
        Funds           : "",
    }
    SendRabbitMessage(CommandEvent,CommandEvent.EventType);

}


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

	db, uid, found, balanceStr := getDatabaseUserId(UserId) 

	if !found {
		fmt.Fprintln(w, "User Does Not Exist")
	}

	fmt.Fprintln(w, UserId + string(uid) + balanceStr)
	
	rows, err := db.Query(getAllTriggers, uid)
	if err != nil{
		return
	}

	var id string
	var uidStr string
	var stock string
	var trigger_type string
	var price string
	var num_shares string
	var created_time time.Time

	for rows.Next() {
		err = rows.Scan(&id, &uidStr, &stock, &trigger_type, &price, &num_shares, &created_time)
		if err != nil{
			return
		}
		fmt.Fprintln(w, id + uidStr + stock + trigger_type + price + num_shares + created_time.String())
	}
	rows.Close()
	return
}




