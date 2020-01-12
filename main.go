package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

// ChargeJSON incoming data for Stripe API
type ChargeJSON struct {
	Amount       int64  `json:"amount"`
	ReceiptEmail string `json:"receiptEmail"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.POST("/api/charge", func(c *gin.Context) {
		var json ChargeJSON
		c.BindJSON(&json)

		// Set Stripe API key
		apiKey := os.Getenv("SK_TEST_KEY")
		stripe.Key = apiKey

		_, err := charge.New(&stripe.ChargeParams{
			Amount:       stripe.Int64(json.Amount),
			Currency:     stripe.String(string(stripe.CurrencyUSD)),
			Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")}, // this should come from clientside
			ReceiptEmail: stripe.String(json.ReceiptEmail)})

		if err != nil {
			// handle
			c.String(http.StatusBadRequest, "Request failed")
			return
		}

		c.String(http.StatusCreated, "Successfully charged")
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
