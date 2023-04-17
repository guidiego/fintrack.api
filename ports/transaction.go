package ports

import "time"

type Transaction struct {
	ID          string     `json:"id"`
	AccountID   string     `json:"accountId"`
	RecipientID *string    `json:"recipientId"`
	BudgetID    *string    `json:"budgetId,omitempty"`
	GoalID      *string    `json:"goalId,omitempty"`
	Description *string    `json:"description,omitempty"`
	Value       float64    `json:"value"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`

	// Aggregated
	BudgetName    *string `json:"budgetName"`
	RecipientName *string `json:"recipientName"`
}
