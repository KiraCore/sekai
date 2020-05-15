package listOrderBooks

type QueryListOrderBooks struct {
	index int `json:"owner"`
}

func (r QueryListOrderBooks) String() int {
	return r.index
}