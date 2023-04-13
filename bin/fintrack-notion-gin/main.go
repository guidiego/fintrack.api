package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	storage "github.com/guidiego/fintrack.api/packages/storage/notion"
	"github.com/guidiego/fintrack.api/ports"
)

func isAuth(authToken string, c *gin.Context) bool {
	return authToken == c.Request.Header["token"][0]
}

func main() {
	r := gin.Default()
	s := storage.New()
	authToken := os.Getenv("API_AUTH_TOKEN")

	r.GET("/budget", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		budgets, err := s.ListBudgets(&ports.Budget{MonthKey: "202304"})

		if err != nil {
			return
		}

		c.IndentedJSON(http.StatusOK, budgets)
	})

	r.GET("/account", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		accounts, err := s.ListAccounts()

		if err != nil {
			return
		}

		c.IndentedJSON(http.StatusOK, accounts)
	})

	r.POST("/transaction", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		var transaction ports.Transaction
		var err error

		if err = c.BindJSON(&transaction); err != nil {
			panic(err)
		}

		transaction, err = s.SaveTransaction(transaction)

		if err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusCreated, transaction)
	})

	r.POST("/transfer", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		var transfer ports.Transfer
		var err error

		if err = c.BindJSON(&transfer); err != nil {
			panic(err)
		}

		desc := fmt.Sprintf("🔄 Money Transfer")
		rmTransaction := ports.Transaction{
			AccountID:   &transfer.FromAccountId,
			Description: &desc,
			Value:       -transfer.Value,
		}

		addTransaction := ports.Transaction{
			AccountID:   &transfer.ToAccountId,
			Description: &desc,
			Value:       transfer.Value,
		}

		_, err = s.SaveTransaction(rmTransaction)

		if err != nil {
			panic(err)
		}

		_, err = s.SaveTransaction(addTransaction)

		if err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"ok": true})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
