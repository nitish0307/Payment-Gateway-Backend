package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"io"
	"log"
	"net/http"
)

func main() {
	stripe.Key = "sk_test_51PsPjR2NXZqtTyA8YjF13ww68sD7wuE1aELmUUKL7Imjf11uHCAzSrkxiNmPXy9f5zFICDg9xhNg8QJ5bL3uTUSU00phgCWT31"

	http.HandleFunc("/create-payment-intent", handleCreatePaymentIntent)
	http.HandleFunc("/health", handlehealth)

	log.Println("Server is running on port 4242.....")
	var err error = http.ListenAndServe("localhost:4242", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handleCreatePaymentIntent(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Only POST method is supported.", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProductId string `json:"product_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address1  string `json:"address_1"`
		Address2  string `json:"address_2"`
		City      string `json:"city"`
		State     string `json:"state"`
		Zip       string `json:"zip"`
		Country   string `json:"country"`
	}

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(req.ProductId)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	fmt.Println(paymentIntent.ClientSecret)

	var response struct {
		ClientSecret string `json:"clientSecret"`
	}

	response.ClientSecret = paymentIntent.ClientSecret

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	writer.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(writer, &buf)
	if err != nil {
		fmt.Println(err)
	}

}

func handlehealth(w http.ResponseWriter, r *http.Request) {
	response := []byte("Server is running")

	_, err := w.Write(response)
	if err != nil {
		fmt.Println(err)
	}
}

func calculateOrderAmount(productId string) int64 {
	switch productId {
	case "Forever Pants":
		return 26000
	case "Forever Shirt":
		return 15500
	case "Forever Shorts":
		return 30000
	}
	return 0
}
