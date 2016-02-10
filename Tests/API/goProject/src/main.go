package main

import (
    "log"
    "net/http"
    "github.com/streadway/amqp"
)

func main() {

    router := NewRouter()

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
    rconn.Close()

    log.Fatal(http.ListenAndServe(":44410", router))


}
