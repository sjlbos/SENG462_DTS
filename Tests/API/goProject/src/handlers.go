package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "os"
    "strings"

    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
    todos := Todos{
        Todo{Name: "Write presentation"},
        Todo{Name: "Host meetup"},
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(todos); err != nil {
        panic(err)
    }
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    fmt.Fprintln(w, "Todo show:", todoId)
}



func Quote(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    StockId := vars["symbol"]
    UserId := vars["id"]

    strEcho :=  StockId + "," + UserId + "\n"
    servAddr := "quoteserve.seng.uvic.ca:4444"

    tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
    if err != nil {
        println("ResolveTCPAddr failed:", err.Error())
        os.Exit(1)
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        println("Dial failed:", err.Error())
        os.Exit(1)
    }

    _, err = conn.Write([]byte(strEcho))
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }

    reply := make([]byte, 1024)
    _, err = conn.Read(reply)

    result := strings.Split(string(reply),",")
    fmt.Fprintln(w, result[0])
    fmt.Fprintln(w, result[1])
    fmt.Fprintln(w, result[2])
    fmt.Fprintln(w, result[3])

    conn.Close()
}

func Add(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Adding Funds to account:")
    type add_struct struct {
        Amount float64
    }
    vars := mux.Vars(r)
    UserId := vars["id"]
    
    decoder := json.NewDecoder(r.Body)
    var t add_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.Amount)

//TODO database stuff!

}

func Buy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Sell Request:")
    type buy_struct struct {
        Amount float64
        Symbol string
    }
    vars := mux.Vars(r)
    UserId := vars["id"]

    decoder := json.NewDecoder(r.Body)
    var t buy_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Symbol)

//TODO database stuff!    
    
}

func CommitBuy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Buy Command Commited:") 
    vars := mux.Vars(r)
    UserId := vars["id"]

    fmt.Fprintln(w, UserId) 
//TODO database Stuff

}

func CancelBuy(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Buy Command Cancelled:") 
    vars := mux.Vars(r)
    UserId := vars["id"]

    fmt.Fprintln(w, UserId) 
//TODO database Stuff

}

func Sell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Sell Request:")
    type sell_struct struct {
        Amount float64
        Symbol string
    }
    vars := mux.Vars(r)
    UserId := vars["id"]

    decoder := json.NewDecoder(r.Body)
    var t sell_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Symbol)

//TODO database stuff!    
    
}


func CommitSell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Sell Command Commited:") 
    vars := mux.Vars(r)
    UserId := vars["id"]

    fmt.Fprintln(w, UserId) 
//TODO database Stuff

}

func CancelSell(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Last Sell Command Cancelled:") 
    vars := mux.Vars(r)
    UserId := vars["id"]

    fmt.Fprintln(w, UserId) 
//TODO database Stuff

}

func CreateBuyTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Creating Buy Trigger:") 
    type trigger_struct struct{
        Amount int
        Price float64
    }

    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]

    decoder := json.NewDecoder(r.Body)
    var t trigger_struct   
    err := decoder.Decode(&t)

    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, Symbol)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Price)

//TODO database stuff

}

func CreateSellTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Creating Sell Trigger:") 
    type trigger_struct struct{
        Amount int
        Price float64
    }

    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]

    decoder := json.NewDecoder(r.Body)
    var t trigger_struct   
    err := decoder.Decode(&t)

    if err != nil {

    }
    fmt.Fprintln(w, UserId)
    fmt.Fprintln(w, Symbol)
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Price)

//TODO database Stuff

}


func CancelBuyTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Cancelling Buy Trigger: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]    
    
    fmt.Fprintln(w, UserId)   
    fmt.Fprintln(w, Symbol) 
//TODO database Stuff

}



func CancelSellTrigger(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Cancelling Sell Trigger: ")
    vars := mux.Vars(r)
    UserId := vars["id"]
    Symbol := vars["symbol"]

    fmt.Fprintln(w, UserId)   
    fmt.Fprintln(w, Symbol) 
    
     
//TODO database Stuff

}

func DumplogUser(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Display Log for User: ")
    vars := mux.Vars(r)
    UserId := vars["id"]

    fmt.Fprintln(w, UserId)
}


func Dumplog(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Dumplog of all transactions: ")


}


func DisplaySummary(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w, "Display Summary for User: ")

    vars := mux.Vars(r)
    UserId := vars["id"]

    fmt.Fprintln(w, UserId)
}





