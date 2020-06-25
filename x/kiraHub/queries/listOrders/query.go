package listOrders

type QueryListOrders struct {
	ID    string `json:"id"`
	Max_Orders string `json:"max_orders"`
	Min_Amount string `json:"min_amount"`
}
