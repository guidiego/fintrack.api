package storagenotion

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fintrack.api/ports"
)

func (s *NotionStorage) ListToSchedule() ([]ports.ToSchedule, error) {
	query := notion.DatabaseQuery{}
	dbItems, cliError := s.cli.QueryDatabase(context.Background(), s.table.ToScheduleID, &query)

	if cliError != nil {
		return []ports.ToSchedule{}, cliError
	}

	budgets := make([]ports.ToSchedule, len(dbItems.Results))
	for i, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		budgets[i] = ports.ToSchedule{
			Day:       *props["Day"].Number,
			Ref:       *&props["Ref"].Title[0].PlainText,
			AccountID: props["Account"].Relation[0].ID,
			Value:     *props["Value"].Number,
			AutoDebit: *props["AutoDebit"].Checkbox,
		}

		if len(props["Budget"].Relation) > 0 {
			budgets[i].BudgetID = props["Budget"].Relation[0].ID
		}
	}

	return budgets, nil
}
