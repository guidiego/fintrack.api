package storagenotion

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fintrack.api/ports"
)

func formatFilter(b ports.Budget) *notion.DatabaseQueryFilter {
	filter := notion.DatabaseQueryFilter{}

	if b.MonthKey != "" {
		filter.Property = "MonthKey"
		filter.RichText = &notion.TextPropertyFilter{
			Equals: b.MonthKey,
		}
	}

	return &filter
}

func (s *NotionStorage) ListBudgets(b *ports.Budget) ([]ports.Budget, error) {
	query := notion.DatabaseQuery{}

	if b != nil {
		query.Filter = formatFilter(*b)
	}

	dbItems, cliError := s.cli.QueryDatabase(context.Background(), s.table.BudgetID, &query)

	if cliError != nil {
		return []ports.Budget{}, cliError
	}

	budgets := make([]ports.Budget, len(dbItems.Results))
	for i, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		budgets[i] = ports.Budget{
			ID:     r.ID,
			Budget: props["Budget"].Title[0].PlainText,
			Limit:  *props["Limit"].Number,
			Used:   *props["Used"].Rollup.Number,
			Free:   *props["Free"].Formula.Number,
		}
	}

	return budgets, nil
}
