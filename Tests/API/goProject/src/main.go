package main

import (
    "log"
    "net/http"
    "github.com/streadway/amqp"
    _ "github.com/lib/pq"
    "flag"
    "database/sql"
    "fmt"
    "os"
)

const (
    DB_USER     = "dts_user"
    DB_PASSWORD = "Group1"
    DB_NAME     = "DTS"
    DB_PORT     = "44410"
//    DB_CONNECTION = "dts_user:Group1@tcp(localhost:44410)/DTS"
)


var rabbitConnectionString string
var rabbitAudit bool
var db *sql.DB
var err error

var getUserId string = "SELECT * FROM \"get_user_account_by_char_id\"($1)"
var addUser string = "SELECT * FROM \"add_user_account\"($1::varchar, $2::money, $3::timestamptz)"
var updateBalance string = "SELECT * FROM \"update_user_account_balance\"($1, $2::money)"
var addPendingPurchase string = "SELECT * FROM \"add_pending_purchase\"($1,$2,$3,$4::money,$5, $6)"
var getLatestPendingPurchase string = "SELEECT * FROM \"get_latest_pending_purchase_for_user\"($1)"
var commitPurchase string = "SELECT * FROM \"commit_pending_purchase\"($1,$2)"
var cancelPurchase string = "SELECT * FROM \"cancel_pending_purchase\"($1)"

var Hostname string

func main() {

    var rhost string
    flag.StringVar(&rhost,"rhost","localhost","name of host for RabbitMQ")

    var rport string
    flag.StringVar(&rport,"rport","5672", "port number for RabbitMQ")

    var port string
    flag.StringVar(&port, "port","44411","port to run Router off of")
    
    flag.BoolVar(&rabbitAudit,"audit",false,"should actions be audited to RabbitMQ")

    flag.Parse()
        
    rabbitConnectionString = "amqp://dts_user:Group1@" + rhost + ":" + rport + "/"


    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
    db, err = sql.Open("postgres", dbinfo)
    failOnError(err, "Failed to connect to DTS Database")

    Hostname, err := os.Hostname()
    println("Running on :", Hostname)

    router := NewRouter()
    
    rconn, err := amqp.Dial(rabbitConnectionString)
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
    rconn.Close()

    RouterPort := ":" + port
    log.Fatal(http.ListenAndServe(RouterPort, router))


}
