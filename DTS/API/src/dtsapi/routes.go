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
    	"CommitBuy",
    	"POST",
    	"/api/users/{id}/pending-purchases/commit",
    	CommitBuy,
    },
    Route{
    	"CancelBuy",
    	"DELETE",
    	"/api/users/{id}/pending-purchases",
    	CancelBuy,
    },
    Route{
    	"Sell",
    	"POST",
    	"/api/users/{id}/pending-sales",
    	Sell,
    },
    Route{
    	"CommitSell",
    	"POST",
    	"/api/users/{id}/pending-sales/commit",
	CommitSell,
    },
    Route{
    	"CancelSell",
    	"DELETE",
    	"/api/users/{id}/pending-sales",
    	CancelSell,
    },
    Route{
    	"CreateBuyTrigger",
    	"PUT",
    	"/api/users/{id}/buy-triggers/{symbol}",
    	CreateBuyTrigger,
    },
    Route{
    	"CreateSellTrigger",
    	"PUT",
    	"/api/users/{id}/sell-triggers/{symbol}",
    	CreateSellTrigger,
    },
    Route{
    	"CancelBuyTrigger",
    	"DELETE",
    	"/api/users/{id}/buy-triggers/{symbol}",
    	CancelBuyTrigger,
    },
    Route{
    	"CancelSellTrigger",
    	"DELETE",
    	"/api/users/{id}/sell-triggers/{symbol}",
    	CancelSellTrigger,
    },
    Route{
    	"PerformBuyTrigger",
    	"POST",
    	"/api/users/{id}/buy-triggers/{symbol}/commit",
    	PerformBuyTrigger,
    },
    Route{
    	"PerformSellTrigger",
    	"POST",
    	"/api/users/{id}/sell-triggers/{symbol}/commit",
    	PerformSellTrigger,
    },
    Route{
    	"DumplogUser",
    	"GET",
    	"/api/users/{id}/transactions",
    	DumplogUser,
    },
    Route{
    	"Dumplog",
    	"GET",
    	"/api/users/transactions",
    	Dumplog,
    },
    Route{
    	"DisplaySummary",
    	"GET",
    	"/api/users/{id}/summary",
    	DisplaySummary,
    },
}
