package ports

type Transaction struct {
	ID          string  `json:"id"`
	AccountID   *string `json:"accountId,omitempty"`
	BudgetID    *string `json:"budgetId,omitempty"`
	Description *string `json:"description,omitempty"`
	Value       float64 `json:"value"`
}
