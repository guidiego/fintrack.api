package storagenotion

import (
	"context"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/google/uuid"
	"github.com/guidiego/fintrack.api/ports"
)

func (n *NotionStorage) SaveTransaction(ipt ports.Transaction) (ports.Transaction, error) {
	p, err := n.cli.CreatePage(context.Background(), TransactionToCreatePageParams(n.table.TransactionID, ipt))

	if err != nil {
		return ports.Transaction{}, err
	}

	ipt.ID = p.ID
	return ipt, err
}

func TransactionToCreatePageParams(tableId string, t ports.Transaction) notion.CreatePageParams {
	now := time.Now()
	emoji := "🟢"
	preTitle := now.Format("2006-01")
	descriptionVal := ""

	if t.Description != nil {
		descriptionVal = *t.Description
	}

	title := []notion.RichText{
		{
			Text: &notion.Text{Content: preTitle + "-" + uuid.New().String()},
		},
	}

	description := []notion.RichText{
		{
			Text: &notion.Text{Content: descriptionVal},
		},
	}

	account := []notion.Relation{}

	if t.AccountID != nil {
		account = append(account, notion.Relation{
			ID: *t.AccountID,
		})
	}

	if t.Value < 0 {
		emoji = "🔴"
	}

	database_properties := notion.DatabasePageProperties{
		"Ref": notion.DatabasePageProperty{
			Title: title,
		},
		"Desc": notion.DatabasePageProperty{
			RichText: description,
		},
		"Account": notion.DatabasePageProperty{
			Relation: account,
		},
		"Value": notion.DatabasePageProperty{
			Number: &t.Value,
		},
	}

	if t.BudgetID != nil {
		database_properties["Budget"] = notion.DatabasePageProperty{
			Relation: []notion.Relation{
				{
					ID: *t.BudgetID,
				},
			},
		}
	}
	return notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   tableId,
		Icon: &notion.Icon{
			Type:  "emoji",
			Emoji: &emoji,
		},
		DatabasePageProperties: &database_properties,
	}
}
