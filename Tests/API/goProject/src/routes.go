package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },
    Route{
        "TodoIndex",
        "GET",
        "/todos",
        TodoIndex,
    },
    Route{
        "Add",
        "PUT",
        "/api/users/{id}",
        Add,
    },
    Route{
	"Quote",
	"GET",
	"/api/users/{id}/stocks/quote/{symbol}",
	Quote,
    },
    Route{
	"Buy",
	"POST",
	"/api/users/{id}/pending-purchases",
	Buy,
    },
    Route{
	"Sell",
	"POST",
	"/api/users/{id}/pending-sales",
	Sell,
    },
}
