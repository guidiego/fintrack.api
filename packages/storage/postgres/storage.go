package storagenotion

import (
	"os"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/guidiego/fintrack.api/ports"
)

type PostgresStorage struct {
	cli *sql.DB
}

func New() *PostgresStorage {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	return &PostgresStorage{
		cli: db,
	}
}

func (p *PostgresStorage) SaveBudget(b ports.Budget) {
	q := `
	INSERT INTO public.budget("account_id", "name", "limit", "year", "month")
	VALUES($1, $2, $3, $4, $5)`

	_, err := p.cli.Exec(q, b.AccountID, b.Name, b.Limit, b.Year, b.Month)

	if err != nil {
		panic(err)
	}
}

func (p *PostgresStorage) SaveGoal(g ports.Goal) {
	q := `
	INSERT INTO public.goal("account_id", "recipient_id", "name", "desired", "status")
	VALUES($1, $2, $3, $4, $5)`

	_, err := p.cli.Exec(q, g.AccountID, g.RecipientId, g.Name, g.Desired, g.Status)

	if err != nil {
		panic(err)
	}
}

func (p *PostgresStorage) SaveTransaction(t ports.Transaction) error {
	q := `
	INSERT INTO public.transaction("value", "description", "budget_id", "recipient_id", "goal_id", "account_id")
	VALUES($1, $2, $3, $4, $5, $6)`

	_, err := p.cli.Exec(q, t.Value, t.Description, t.BudgetID, t.RecipientID, t.GoalID, "123")

	return err
}

func (p *PostgresStorage) SaveUpComming(u ports.UpComming) error {
	q := `
	INSERT INTO public.upcomming("name", "day", "value", "account_id", "budget_id", "recipient_id", "auto_debit")
	VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.cli.Exec(q, u.Name, u.Day, u.Value, u.AccountID, u.BudgetID, u.RecipientID, u.AutoDebit)

	return err
}

func (p *PostgresStorage) ListBudgets(ipt ports.BudgetFilterInput) ([]ports.Budget, error) {
	q := `
		SELECT b.id, b.name, b.account_id, b.year, b.month, b.limit, SUM(t.value) as used
		FROM public.budget as b
		LEFT JOIN public.transaction as t ON b.id = t.budget_id
		WHERE b.account_id = $1 AND b.year = $2 AND b.month = $3
		GROUP BY b.id
	`

	r, err := p.cli.Query(q, ipt.AccountID, ipt.Year, ipt.Month)

	if err != nil {
		return []ports.Budget{}, err
	}

	budgets := []ports.Budget{}
	for r.Next() {
		b := ports.Budget{}
		r.Scan(&b.ID, &b.Name, &b.AccountID, &b.Year, &b.Month, &b.Limit, &b.Used)
		budgets = append(budgets, b)
	}

	return budgets, nil
}

func (p *PostgresStorage) ListRecipients(accountID string) ([]ports.Recipient, error) {
	q := `
		SELECT r.id, r.name, r.account_id, SUM(t.value) as total_value
		FROM public.recipient as r
		LEFT JOIN public.transaction as t ON r.id = t.recipient_id
		WHERE r.account_id = $1
		GROUP BY r.id
	`

	res, err := p.cli.Query(q, accountID)

	if err != nil {
		return []ports.Recipient{}, err
	}

	recipients := []ports.Recipient{}
	for res.Next() {
		r := ports.Recipient{}
		res.Scan(&r.ID, &r.Name, &r.AccountID, &r.TotalValue)
		recipients = append(recipients, r)
	}

	return recipients, nil
}

func (p *PostgresStorage) ListUpComming(accountID string, day int) ([]ports.UpComming, error) {
	q := `
		SELECT u.id, u.name, u.day, u.account_id, u.value, b.id as budget_id, b.name as budget_name, r.id as recipient_id, r.name as recipient_name FROM upcomming as u
		INNER JOIN budget as b ON b.id = u.budget_id
		INNER JOIN recipient as r ON r.id = u.recipient_id
		WHERE u.account_id = $1 and (u.day >= $2 OR u.day BETWEEN $3 AND 0)
	`

	res, err := p.cli.Query(q, accountID, day, -day)

	if err != nil {
		return []ports.UpComming{}, err
	}

	upcommings := []ports.UpComming{}
	for res.Next() {
		u := ports.UpComming{}
		res.Scan(&u.ID, &u.Name, &u.Day, &u.AccountID, &u.Value, &u.BudgetID, &u.BudgetName, &u.RecipientID, &u.RecipientName)
		upcommings = append(upcommings, u)
	}

	return upcommings, nil
}

func (p *PostgresStorage) ListTransactions(accountID string) ([]ports.Transaction, error) {
	q := `
		SELECT
			t.id,
			t.description,
			t.value,
			t.created_at,
			t.updated_at,
			t.recipient_id,
			t.budget_id,
			b.name as bname,
			r.name as rname
		FROM transaction as t
		INNER JOIN budget as b ON b.id = t.budget_id
		INNER JOIN recipient as r ON r.id = t.recipient_id
		WHERE t.account_id = $1
		ORDER BY t.created_at DESC
		LIMIT 10
	`

	res, err := p.cli.Query(q, accountID)

	if err != nil {
		return []ports.Transaction{}, err
	}

	transactions := []ports.Transaction{}
	for res.Next() {
		t := ports.Transaction{
			AccountID: "123",
		}

		res.Scan(
			&t.ID,
			&t.Description,
			&t.Value,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.RecipientID,
			&t.BudgetID,
			&t.BudgetName,
			&t.RecipientName,
		)

		transactions = append(transactions, t)
	}

	return transactions, nil
}
