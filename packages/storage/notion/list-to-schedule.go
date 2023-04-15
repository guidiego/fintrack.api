package storagenotion

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fintrack.api/ports"
)

func (s *NotionStorage) ListToSchedule(f *ports.ToScheduleFilterInput) ([]ports.ToSchedule, error) {
	query := notion.DatabaseQuery{}

	if f != nil {
		qs := []notion.DatabaseQueryFilter{}

		if f.AutoDebit != nil {
			qs = append(qs, notion.DatabaseQueryFilter{
				Property: "AutoDebit",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Checkbox: &notion.CheckboxDatabaseQueryFilter{
						Equals: f.AutoDebit,
					},
				},
			})
		}

		if f.FromDay != nil {
			negative := -*f.FromDay
			zero := 0
			qs = append(qs, notion.DatabaseQueryFilter{
				Or: []notion.DatabaseQueryFilter{
					{
						Property: "Day",
						DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
							Number: &notion.NumberDatabaseQueryFilter{
								GreaterThanOrEqualTo: f.FromDay,
							},
						},
					},
					{
						Property: "Day",
						DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
							Number: &notion.NumberDatabaseQueryFilter{
								GreaterThanOrEqualTo: &negative,
								LessThan:             &zero,
							},
						},
					},
				},
			})
		}

		query.Filter = &notion.DatabaseQueryFilter{
			And: qs,
		}
	}

	dbItems, cliError := s.cli.QueryDatabase(context.Background(), s.table.ToScheduleID, &query)

	if cliError != nil {
		return []ports.ToSchedule{}, cliError
	}

	toschedules := make([]ports.ToSchedule, len(dbItems.Results))
	for i, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		toschedules[i] = ports.ToSchedule{
			Day:       *props["Day"].Number,
			Ref:       *&props["Ref"].Title[0].PlainText,
			AccountID: props["Account"].Relation[0].ID,
			Value:     *props["Value"].Number,
			AutoDebit: *props["AutoDebit"].Checkbox,
		}

		if len(props["Budget"].Relation) > 0 {
			toschedules[i].BudgetID = &props["Budget"].Relation[0].ID
		}
	}

	return toschedules, nil
}
