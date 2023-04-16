package storagenotion

import (
	"context"
	"sync"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fintrack.api/ports"
)

func (s *NotionStorage) ListBudgets(ipt *ports.BudgetFilterInput) ([]ports.Budget, error) {
	var wg sync.WaitGroup
	query := notion.DatabaseQuery{}

	if ipt.MonthKey != nil {
		query.Filter = &notion.DatabaseQueryFilter{
			Property: "MonthKey",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				RichText: &notion.TextPropertyFilter{
					Equals: *ipt.MonthKey,
				},
			},
		}
	}

	dbItems, cliError := s.cli.QueryDatabase(context.Background(), s.table.BudgetID, &query)

	if cliError != nil {
		return []ports.Budget{}, cliError
	}

	budgets := make([]ports.Budget, 0)
	wg.Add(len(dbItems.Results))
	for _, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)

		go func() {
			defer wg.Done()
			prop := props["Used"]

			budget := ports.Budget{
				ID:       r.ID,
				Budget:   props["Budget"].Title[0].PlainText,
				Limit:    *props["Limit"].Number,
				Used:     *prop.Rollup.Number,
				Free:     *props["Free"].Formula.Number,
				MonthKey: *ipt.MonthKey,
			}

			if *ipt.ExactValues {
				field, err := s.cli.FindPagePropertyByID(context.Background(), r.ID, prop.ID, nil)

				if err != nil {
					budget.Used = *field.PropertyItem.Rollup.Number
					budget.Free = budget.Limit + budget.Used
				}
			}

			budgets = append(budgets, budget)
		}()
	}

	wg.Wait()
	return budgets, nil
}
