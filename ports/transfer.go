package ports

type Transfer struct {
	FromRecipientID string  `json:"fromRecipientId"`
	ToRecipientID   string  `json:"toRecipientId"`
	Value           float64 `json:"value"`
}
