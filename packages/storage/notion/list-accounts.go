package storagenotion

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fintrack.api/ports"
)

func (s *NotionStorage) ListAccounts() ([]ports.Account, error) {
	dbItems, cliError := s.cli.QueryDatabase(context.Background(), s.table.AccountID, nil)

	if cliError != nil {
		return []ports.Account{}, cliError
	}

	accounts := make([]ports.Account, len(dbItems.Results))
	for idx, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)

		accounts[idx] = ports.Account{
			ID:        r.ID,
			Name:      props["Name"].Title[0].PlainText,
			Total:     *props["Total"].Rollup.Number,
			Allocated: *props["Allocated"].Formula.Number,
			Free:      *props["Free"].Formula.Number,
		}
	}

	return accounts, nil
}
