package ports

type Recipient struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AccountID string `json:"accountID"`

	// Parsed
	TotalValue *float64 `json:"totalValue"`
}

type Budget struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AccountID string  `json:"accountID"`
	Limit     float64 `json:"limit"`
	Year      int16   `json:"year"`
	Month     int16   `json:"month"`

	// Value From Aggregated
	Used *float64 `json:"used"`
}

type UpComming struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Day         int64   `json:"day"`
	AccountID   string  `json:"accountId"`
	Value       float64 `json:"value"`
	BudgetID    *string `json:"budgetId,omitempty"`
	RecipientID *string `json:"recipientId,omitempty"`
	AutoDebit   bool    `json:"autoDebit"`

	// Extra
	BudgetName    *string `json:"budgetName"`
	RecipientName *string `json:"recipientName"`
}

type Goal struct {
	ID          string  `json:"id"`
	AccountID   string  `json:"accountId"`
	RecipientId string  `json:"recipientId"`
	Name        string  `json:"name"`
	Desired     float64 `json:"desired"`
	Status      int     `json:"status"`
}

type ToScheduleFilterInput struct {
	AutoDebit *bool
	FromDay   *int
}

type BudgetFilterInput struct {
	AccountID string
	Month     int32
	Year      int32
}

type Storage interface {
	SaveTransaction(Transaction) (Transaction, error)
	ListRecipients() ([]Recipient, error)
	ListBudgets(*BudgetFilterInput) ([]Budget, error)
}
