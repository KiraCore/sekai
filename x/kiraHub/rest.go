package kiraHub

import (
	"strings"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/gorilla/mux"

	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
	"github.com/KiraCore/sekai/x/kiraHub/queries/listOrderBooks"
	"github.com/KiraCore/sekai/x/kiraHub/queries/listOrders"
)

func RegisterRESTRoutes(cliContext client.Context, router *mux.Router) {
	//router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.CreateOrderBookTransaction}, "/"), createOrderBook.RestRequestHandler(cliContext)).Methods("POST")
	//router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.CreateOrderTransaction}, "/"), createOrder.RestRequestHandler(cliContext)).Methods("POST")

	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.ListOrderBooksQuery}, "/"), listOrderBooks.GetOrderBooks(cliContext)).Methods("GET")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.ListOrderBooksQueryByTP}, "/"), listOrderBooks.GetOrderBooksByTP(cliContext)).Methods("GET")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.ListOrders}, "/"), listOrders.GetOrders(cliContext)).Methods("GET")
}
