package main

import (
    "log"
    "net/http"
    "net"
    "github.com/streadway/amqp"
    "github.com/rs/cors"
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
var Hostname string

func main() {

    type Configuration struct {
        RabbitHost     string
    	RabbitPort     string
    	HostPort       string
    	QuoteRunnerHost string
    	QuoteRunnerPort string
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

    var qhost string = configuration.QuoteRunnerHost
    var qport string = configuration.QuoteRunnerPort
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
        dbPointers[i].SetMaxOpenConns(200)
        dbPointers[i].SetMaxIdleConns(200)
    }

    Hostname, err = os.Hostname()
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

	c := cors.New(cors.Options{
	    AllowedOrigins: []string{"*"},
	    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	    AllowCredentials: true,
	    Debug : false,
	})

	// Insert the middleware
	handler := c.Handler(router)

    log.Fatal(http.ListenAndServe(RouterPort, handler))
}
