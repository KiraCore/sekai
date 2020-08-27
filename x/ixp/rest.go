package ixp

import (
	"strings"

	"github.com/KiraCore/sekai/x/ixp/client/rest"
	"github.com/KiraCore/sekai/x/ixp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

// RegisterRESTRoutes register rest routes
func RegisterRESTRoutes(cliContext client.Context, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, types.CreateOrderBookTransaction}, "/"), rest.CreateOrderbookRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, types.CreateOrderTransaction}, "/"), rest.CreateOrderRequestHandler(cliContext)).Methods("POST")
	// TODO: should add cancel order rest endpoint

	router.HandleFunc(strings.Join([]string{"", QuerierRoute, types.ListOrderBooksQuery}, "/"), rest.GetOrderBooks(cliContext)).Methods("GET")
	router.HandleFunc(strings.Join([]string{"", QuerierRoute, types.ListOrderBooksQueryByTradingPair}, "/"), rest.GetOrderBooksByTradingPair(cliContext)).Methods("GET")
	router.HandleFunc(strings.Join([]string{"", QuerierRoute, types.ListOrders}, "/"), rest.GetOrders(cliContext)).Methods("GET")
}
