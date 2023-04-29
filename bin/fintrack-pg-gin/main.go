package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	storage "github.com/guidiego/fintrack.api/packages/storage/postgres"
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

		year, _ := strconv.ParseInt(c.Query("year"), 10, 64)
		month, _ := strconv.ParseInt(c.Query("month"), 10, 64)

		budgets, err := s.ListBudgets(ports.BudgetFilterInput{
			AccountID: "123",
			Month:     int32(month),
			Year:      int32(year),
		})

		if err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusOK, budgets)
	})

	r.GET("/recipient", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		recipients, err := s.ListRecipients("123")

		if err != nil {
			return
		}

		c.IndentedJSON(http.StatusOK, recipients)
	})

	r.GET("/transaction", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		transactions, err := s.ListTransactions("123")

		if err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusCreated, transactions)
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

		err = s.SaveTransaction(transaction)
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
			RecipientID: &transfer.FromRecipientID,
			Description: &desc,
			Value:       -transfer.Value,
		}

		addTransaction := ports.Transaction{
			RecipientID: &transfer.ToRecipientID,
			Description: &desc,
			Value:       transfer.Value,
		}

		err = s.SaveTransaction(rmTransaction)
		if err != nil {
			panic(err)
		}

		err = s.SaveTransaction(addTransaction)
		if err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"ok": true})
	})

	r.GET("/check-health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"ok": true})
	})

	r.GET("/upcomming", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		todayStr := c.Query("day")
		if todayStr == "" {
			now := time.Now()
			todayStr = now.Format("02")
		}

		today, _ := strconv.Atoi(todayStr)
		upcomming, err := s.ListUpComming("123", today)

		if err != nil {
			log.Fatalln(err)
			return
		}

		c.IndentedJSON(http.StatusCreated, upcomming)
	})

	r.GET("/goal", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		status := 2
		goals, err := s.ListGoal("123", &status)

		if err != nil {
			return
		}

		c.IndentedJSON(http.StatusOK, goals)
	})

	r.POST("/automation/schedule", func(c *gin.Context) {
		if !isAuth(authToken, c) {
			return
		}

		now := time.Now()
		todayStr := now.Format("02")
		today, _ := strconv.ParseInt(todayStr, 10, 64)
		lastDayTime := now.AddDate(0, 1, -int(today))

		schedules, err := s.ListUpComming("123", -int(today))

		if err != nil {
			return
		}

		lastDay, err := strconv.ParseInt(lastDayTime.Format("02"), 10, 64)
		if err != nil {
			return
		}

		for _, schedule := range schedules {
			scheduledDay := schedule.Day

			if scheduledDay < 0 {
				scheduledDay = lastDay + scheduledDay
			}

			if scheduledDay == today {
				t := ports.Transaction{
					Value:       schedule.Value,
					Description: &schedule.Name,
					AccountID:   schedule.AccountID,
					RecipientID: schedule.RecipientID,
					BudgetID:    schedule.BudgetID,
				}

				if schedule.BudgetID != nil {
					t.BudgetID = schedule.BudgetID
				}

				err := s.SaveTransaction(t)
				if err != nil {
					fmt.Printf("%e", err)
				}
			}
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"ok": true})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
