package storagenotion

import (
	"os"

	"github.com/dstotijn/go-notion"
)

type NotionTableIds struct {
	TransactionID string
	ToScheduleID  string
	BudgetID      string
	AccountID     string
}

type NotionStorage struct {
	cli   *notion.Client
	table NotionTableIds
}

func New() *NotionStorage {
	tableIds := NotionTableIds{
		os.Getenv("NOTION_TRANSACTION_TABLE_ID"),
		os.Getenv("NOTION_TO_SCHEDULE_TABLE_ID"),
		os.Getenv("NOTION_BUDGETS_TABLE_ID"),
		os.Getenv("NOTION_ACCOUNTS_TABLE_ID"),
	}

	return &NotionStorage{
		cli:   notion.NewClient(os.Getenv("NOTION_API_TOKEN")),
		table: tableIds,
	}
}
