package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "bufio"

    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
    conn, err := net.Dial("udp", "http://quoteserve.seng.uvic.ca:4441")
    if err != nil {
        // handle error
    }
    fmt.Fprintf(conn, "randomString")
    status, err := bufio.NewReader(conn).ReadString('\n')
    fmt.Fprintln(w, status)
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