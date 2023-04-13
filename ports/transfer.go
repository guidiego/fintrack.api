package ports

type Transfer struct {
	FromAccountId string  `json:"fromAccountId"`
	ToAccountId   string  `json:"toAccountId"`
	Value         float64 `json:"value"`
}
