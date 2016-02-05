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
    println("write to server = ", strEcho)

    reply := make([]byte, 1024)
    _, err = conn.Read(reply)
    println("reply from server=", string(reply))

    result := strings.Split(string(reply),",")
    fmt.Fprintln(w, result[0])
    fmt.Fprintln(w, result[1])
    fmt.Fprintln(w, result[2])
    fmt.Fprintln(w, result[3])

    conn.Close()
}

func Add(w http.ResponseWriter, r *http.Request){

    type add_struct struct {
        Amount float64
    }
    decoder := json.NewDecoder(r.Body)
    var t add_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, t.Amount)

//TODO database stuff!

}

func Buy(w http.ResponseWriter, r *http.Request){

    type buy_struct struct {
        Amount float64
        Symbol string
    }

    decoder := json.NewDecoder(r.Body)
    var t buy_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Symbol)

//TODO database stuff!    
    
}

func Sell(w http.ResponseWriter, r *http.Request){

    type sell_struct struct {
        Amount float64
        Symbol string
    }

    decoder := json.NewDecoder(r.Body)
    var t sell_struct   
    err := decoder.Decode(&t)
    if err != nil {

    }
    fmt.Fprintln(w, t.Amount)
    fmt.Fprintln(w, t.Symbol)

//TODO database stuff!    
    
}





