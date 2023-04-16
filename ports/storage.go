package ports

type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Total     float64
	Allocated float64
	Free      float64
}

type Budget struct {
	ID       string  `json:"id"`
	Budget   string  `json:"budget"`
	Limit    float64 `json:"limit"`
	Used     float64 `json:"used"`
	Free     float64 `json:"free"`
	MonthKey string  `json:"monthKey"`
}

type ToSchedule struct {
	Ref       string  `json:"ref"`
	Day       float64 `json:"day"`
	AccountID string
	Value     float64 `json:"value"`
	BudgetID  *string
	AutoDebit bool `json:"autoDebit"`
}

type Goal struct {
}

type ToScheduleFilterInput struct {
	AutoDebit *bool
	FromDay   *int
}

type BudgetFilterInput struct {
	MonthKey    *string
	ExactValues *bool
}

type Storage interface {
	SaveTransaction(Transaction) (Transaction, error)
	ListAccounts() ([]Account, error)
	ListBudgets(*BudgetFilterInput) ([]Budget, error)
	ListToSchedule(*ToScheduleFilterInput) ([]ToSchedule, error)
	// ListGoals(*Goal) ([]Goal, error)
}
