package listOrderBooks

type QueryListOrderBooks struct {
	By    string `json:"by"`
	Value string `json:"value"`
}

type QueryListOrderBooksByTP struct {
	Base    string `json:"base"`
	Quote   string `json:"quote"`
}

