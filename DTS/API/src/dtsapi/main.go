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
    "netpool"
    //"sync"
)

const (
    DB_USER     = "dts_user"
    DB_PASSWORD = "Group1"
    DB_NAME     = "dts"
    DB_PORT     = "44410"
    DB_HOST     = "B133.seng.uvic.ca"
)


var rabbitConnectionString string
var quoteCacheConnectionString string
var rabbitAudit bool
var db *sql.DB
var err error
var rconn *amqp.Connection
var ch *amqp.Channel
var quoteCache net.Conn
var QuoteNetpool *netpool.Netpool
//var NetpoolMutex *sync.Mutex

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

var Hostname string

func main() {

    type Configuration struct {
        RabbitHost     string
	RabbitPort     string
	HostPort       string
	QuoteCacheHost string
	QuoteCachePort string
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


    dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
    db, err = sql.Open("postgres", dbinfo)
    failOnError(err, "Failed to connect to DTS Database")

    db.SetMaxOpenConns(50)
    db.SetMaxIdleConns(50)


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

    RouterPort := ":" + port
    log.Fatal(http.ListenAndServe(RouterPort, router))


}
