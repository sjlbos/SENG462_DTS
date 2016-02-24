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
        "/api/users/{id}/{TransNo}",
        Add,
    },
    Route{
	 "Quote",
    	"GET",
    	"/api/users/{id}/stocks/quote/{symbol}/{TransNo}",
    	Quote,
    },
    Route{
    	"Buy",
    	"POST",
    	"/api/users/{id}/pending-purchases/{TransNo}",
    	Buy,
    },
    Route{
    	"CommitBuy",
    	"POST",
    	"/api/users/{id}/pending-purchases/commit/{TransNo}",
    	CommitBuy,
    },
    Route{
    	"CancelBuy",
    	"DELETE",
    	"/api/users/{id}/pending-purchases/{TransNo}",
    	CancelBuy,
    },
    Route{
    	"Sell",
    	"POST",
    	"/api/users/{id}/pending-sales/{TransNo}",
    	Sell,
    },
    Route{
    	"CommitSell",
    	"POST",
    	"/api/users/{id}/pending-sales/commit/{TransNo}",
	   CommitSell,
    },
    Route{
    	"CancelSell",
    	"DELETE",
    	"/api/users/{id}/pending-sales/{TransNo}",
    	CancelSell,
    },
    Route{
    	"CreateBuyTrigger",
    	"PUT",
    	"/api/users/{id}/buy-triggers/{symbol}/{TransNo}",
    	CreateBuyTrigger,
    },
    Route{
    	"CreateSellTrigger",
    	"PUT",
    	"/api/users/{id}/sell-triggers/{symbol}/{TransNo}",
    	CreateSellTrigger,
    },
    Route{
    	"CancelBuyTrigger",
    	"DELETE",
    	"/api/users/{id}/buy-triggers/{symbol}/{TransNo}",
    	CreateBuyTrigger,
    },
    Route{
    	"CancelSellTrigger",
    	"DELETE",
    	"/api/users/{id}/sell-triggers/{symbol}/{TransNo}",
    	CreateSellTrigger,
    },
    Route{
    	"DumplogUser",
    	"GET",
    	"/api/users/{id}/transactions/{TransNo}",
    	DumplogUser,
    },
    Route{
    	"Dumplog",
    	"GET",
    	"/api/users/transactions/{TransNo}",
    	Dumplog,
    },
    Route{
    	"DisplaySummary",
    	"GET",
    	"/api/users/{id}/summary/{TransNo}",
    	DisplaySummary,
    },
}
