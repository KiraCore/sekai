package kiraHub

import (
	"strings"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/sekai/x/kiraHub/client/rest"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/gorilla/mux"
)

func RegisterRESTRoutes(cliContext client.Context, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, types.CreateOrderBookTransaction}, "/"), rest.RestCreateOrderRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, types.CreateOrderTransaction}, "/"), rest.RestCreateOrderRequestHandler(cliContext)).Methods("POST")

	router.HandleFunc(strings.Join([]string{"", QuerierRoute, types.ListOrderBooksQuery}, "/"), rest.GetOrderBooks(cliContext)).Methods("GET")
	router.HandleFunc(strings.Join([]string{"", QuerierRoute, types.ListOrderBooksQueryByTP}, "/"), rest.GetOrderBooksByTP(cliContext)).Methods("GET")
	router.HandleFunc(strings.Join([]string{"", QuerierRoute, types.ListOrders}, "/"), rest.GetOrders(cliContext)).Methods("GET")
}
