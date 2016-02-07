package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "os"
    "log"
    "strings"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/streadway/amqp"
)


//topic exchange DtsEvents, durable
//TransactionEvent.*
/*
   <xsd:element name="userCommand" type="UserCommandType"/>
   <xsd:element name="quoteServer" type="QuoteServerType"/>
   <xsd:element name="accountTransaction" type="AccountTransactionType"/>
   <xsd:element name="systemEvent" type="SystemEventType"/>
   <xsd:element name="errorEvent" type="ErrorEventType"/>
   <xsd:element name="debugEvent" type="DebugType"/>
*/

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}



func Quote(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    StockId := vars["symbol"]
    UserId := vars["id"]

    type QuoteCommand struct{
        Price   float64
        StockSymbol string
        QuoteServerTime string
        Cryptokey string
    }

    strEcho :=  StockId + "," + UserId + "\n"
    servAddr := "quoteserve.seng.uvic.ca:4444"

    tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
    if err != nil {
        println("ResolveTCPAddr failed:", err.Error())
        os.Exit(1)
    }

    qconn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        println("Dial failed:", err.Error())
        os.Exit(1)
    }

    _, err = qconn.Write([]byte(strEcho))
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }

    reply := make([]byte, 1024)
    _, err = qconn.Read(reply)
    println(reply)
    result := strings.Split(string(reply),",")
    tmpResult0, err := strconv.ParseFloat(result[0], 64)
    fmt.Fprintln(w, tmpResult0)
    fmt.Fprintln(w, result[1])
    fmt.Fprintln(w, result[2])
    fmt.Fprintln(w, result[3])
    tmpResult4 := stripCtlAndExtFromUTF8(result[4])
    fmt.Fprintln(w, tmpResult4)

    //Audit Quote

    rconn, err := amqp.Dial("amqp://dts_user:Group1@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer rconn.Close()

    ch, err := rconn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    err = ch.ExchangeDeclare(
            "DtsEvents", // name
            "topic",      // type
            true,         // durable
            true,        // auto-deleted
            false,        // internal
            false,        // no-wait
            nil,          // arguments
    )
    failOnError(err, "Failed to declare an exchange")

    q := QuoteCommand{tmpResult0,result[1],result[3],tmpResult4}
    body, err := json.Marshal(q)

    
    err = ch.Publish(
        "DtsEvents",          // exchange
        "TransactionEvent.quoteServer", // routing key
        false, // mandatory
        false, // immediate
        amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(body),
        })
    failOnError(err, "Failed to publish a message")

    log.Printf(" [x] Sent %s", body)

    rconn.Close()


    qconn.Close()
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





