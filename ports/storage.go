package ports

type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Total     float64
	Allocated float64
	Free      float64
}

type Budget struct {
	ID       string `json:"id"`
	Budget   string `json:"budget"`
	Limit    float64
	Used     float64
	Free     float64
	MonthKey string
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

type Storage interface {
	SaveTransaction(Transaction) (Transaction, error)
	ListAccounts() ([]Account, error)
	ListBudgets(*Budget) ([]Budget, error)
	ListToSchedule(filter *ToScheduleFilterInput) ([]ToSchedule, error)
	// ListGoals(*Goal) ([]Goal, error)
}
