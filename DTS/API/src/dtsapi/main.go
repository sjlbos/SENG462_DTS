package main

import (
    "log"
    "net/http"
    "net"
    "github.com/streadway/amqp"
    _ "github.com/lib/pq"
    "flag"
    "database/sql"
    "fmt"
    "os"
    "runtime/pprof"
    "encoding/json"
)


var rabbitConnectionString string
var quoteCacheConnectionString string
var rabbitAudit bool
var dbPointers []*sql.DB
var err error
var rconn *amqp.Connection
var ch *amqp.Channel
var quoteCache net.Conn

var getUserId string = "SELECT * FROM \"get_user_account_by_char_id\"($1)"
var addUser string = "SELECT * FROM \"add_user_account\"($1::varchar, $2::money, $3::timestamptz)"
var updateBalance string = "SELECT * FROM \"update_user_account_balance\"($1, $2::money)"
var addPendingPurchase string = "SELECT * FROM \"add_pending_purchase\"($1,$2,$3::int,$4::money,$5, $6)"
var getLatestPendingPurchase string = "SELECT * FROM \"get_latest_pending_purchase_for_user\"($1)"
var commitPurchase string = "SELECT * FROM \"commit_pending_purchase\"($1,$2)"
var addPendingSale string = "SELECT * FROM \"add_pending_sale\"($1,$2,$3::int,$4::money,$5, $6)"
var getLatestPendingSale string = "SELECT * FROM \"get_latest_pending_sale_for_user\"($1::int)"
var commitSale string = "SELECT * FROM \"commit_pending_sale\"($1,$2)"
var cancelTransaction string = "SELECT * FROM \"cancel_pending_transaction\"($1)"
var addBuyTrigger string = "SELECT * FROM \"add_buy_trigger\"($1::int,$2::varchar,$3::money,$4::timestamptz)"
var addSellTrigger string = "SELECT * FROM \"add_buy_trigger\"($1::int,$2::varchar,$3::money,$4::timestamptz)"
var setBuyTrigger string = "SELECT * FROM \"commit_buy_trigger\"($1::int,$2::int,$3::money,$4::timestamptz)"
var setSellTrigger string = "SELECT * FROM \"commit_sell_trigger\"($1::int,$2::int,$3::money,$4::timestamptz)"
var cancelBuyTrigger string = "SELECT * FROM \"cancel_buy_trigger\"($1::int)"
var cancelSellTrigger string = "SELECT * FROM \"cancel_sell_trigger\"($1::int)"
var getBuyTriggerId string = "SELECT * FROM \"get_buy_trigger_id_for_user_and_stock\"($1::int, $2::varchar)"
var getSellTriggerId string = "SELECT * FROM \"get_sell_trigger_id_for_user_and_stock\"($1::int, $2::varchar)"
var getPendingTriggerId string = "SELECT * FROM \"get_pending_trigger_id_for_user_and_stock\"($1::int, $2::varchar, $3::trigger_type)"
var getTriggerById string = "SELECT * FROM \"get_trigger_by_id\"($1)"
var performBuyTrigger string = "SELECT * FROM \"perform_buy_trigger\"($1, $2::money)"
var performSellTrigger string = "SELECT * FROM \"perform_sell_trigger\"($1, $2::money)"


var Hostname string

func main() {

    type Configuration struct {
        RabbitHost     string
    	RabbitPort     string
    	HostPort       string
    	QuoteCacheHost string
    	QuoteCachePort string
        DBConnectionString []string
    }
    file, err := os.Open("conf.json")
    if err != nil {
        fmt.Println("error:", err)
    }
    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    err = decoder.Decode(&configuration)
    if err != nil {
        fmt.Println("error:", err)
    }

    var qhost string = configuration.QuoteCacheHost
    var qport string = configuration.QuoteCachePort
    var rhost string = configuration.RabbitHost
    var rport string = configuration.RabbitPort
    var port string = configuration.HostPort

    println("Connected to Rabbit: " + rhost + ":" + rport)
    println("Connected to QuoteCache: " + qhost + ":" + qport)
    println("Running locally on port " + port)

    var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    quoteCacheConnectionString = qhost + ":" + qport
        
    rabbitConnectionString = "amqp://dts_user:Group1@" + rhost + ":" + rport + "/"

    for i := range configuration.DBConnectionString{
        var db *sql.DB
        dbPointers = append(dbPointers, db)
        dbinfo := configuration.DBConnectionString[i]
        dbPointers[i], err = sql.Open("postgres", dbinfo)
        failOnError(err, "Failed to connect to DTS Database")

        dbPointers[i].SetMaxOpenConns(50)
        dbPointers[i].SetMaxIdleConns(50)

    }

    Hostname, err := os.Hostname()
    println("Running on :", Hostname)

    router := NewRouter()
    
    rconn, err = amqp.Dial(rabbitConnectionString)
    failOnError(err, "Failed to connect to RabbitMQ")
    defer rconn.Close()

    ch, err = rconn.Channel()
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

    err = ch.ExchangeDeclare(
            "Dts", // name
            "topic",      // type
            true,         // durable
            true,        // auto-deleted
            false,        // internal
            false,        // no-wait
            nil,          // arguments
    )
    failOnError(err, "Failed to declare an exchange")

    RouterPort := ":" + port
    log.Fatal(http.ListenAndServe(RouterPort, router))


}
