package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	storage "github.com/guidiego/fintrack.api/packages/storage/notion"
	"github.com/guidiego/fintrack.api/ports"
)

func isAuth(authToken string, c *gin.Context) bool {
	token := c.Request.Header.Get("token")

	if token == "" {
		token = c.Request.Header.Get("Token")
	}

	return authToken == token
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

	r.GET("/to-schedule", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		now := time.Now()
		todayStr := now.Format("02")
		today, _ := strconv.Atoi(todayStr)

		schedules, err := s.ListToSchedule(&ports.ToScheduleFilterInput{
			FromDay: &today,
		})

		if err != nil {
			log.Fatalln(err)
			return
		}

		c.IndentedJSON(http.StatusCreated, schedules)
	})

	r.POST("/automation/schedule", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		now := time.Now()
		todayStr := now.Format("02")
		today, _ := strconv.Atoi(todayStr)
		lastDayTime := now.AddDate(0, 1, -today)

		onlyAutoDebit := true
		schedules, err := s.ListToSchedule(&ports.ToScheduleFilterInput{
			AutoDebit: &onlyAutoDebit,
			FromDay:   &today,
		})

		if err != nil {
			return
		}

		lastDay, err := strconv.ParseFloat(lastDayTime.Format("02"), 0)
		if err != nil {
			return
		}

		for _, schedule := range schedules {
			scheduledDay := schedule.Day

			if scheduledDay < 0 {
				scheduledDay = lastDay + scheduledDay
			}

			if scheduledDay == float64(today) {
				t := ports.Transaction{
					Value:       schedule.Value,
					Description: &schedule.Ref,
					AccountID:   &schedule.AccountID,
				}

				if schedule.BudgetID != nil {
					t.BudgetID = schedule.BudgetID
				}

				_, err := s.SaveTransaction(t)

				if err != nil {
					fmt.Printf("%e", err)
				}
			}
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"ok": true})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
