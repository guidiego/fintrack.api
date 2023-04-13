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
	Ref       string
	Day       float64
	AccountID string
	Value     float64
	BudgetID  string
	AutoDebit bool
}

type Goal struct {
}

type Storage interface {
	SaveTransaction(Transaction) (Transaction, error)
	ListAccounts() ([]Account, error)
	ListBudgets(*Budget) ([]Budget, error)
	ListToSchedule() ([]ToSchedule, error)
	// ListGoals(*Goal) ([]Goal, error)
}
