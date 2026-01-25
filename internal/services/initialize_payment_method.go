package services

// import (
// 	"fmt"
// 	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
// 	"github.com/stripe/stripe-go/v74"
// 	"github.com/stripe/stripe-go/v74/checkout/session"
// 	"log"
// 	"os"
// )

// func InitializeStripePayment(payment *gen.Payment) (string, error) {
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
// 	if stripe.Key == "" {
// 		log.Fatal("STRIPE_SECRET_KEY is required")
// 	}
// 	log.Printf("Initializing Stripe payment: payment_id=%s, transaction_id=%s, amount=%f, currency=%s",
// 		payment.ID, payment.TransactionID.String(), payment.Amount, payment.Currency)

// 	// Basic validations
// 	if payment.Amount <= 0 {
// 		err := fmt.Errorf("payment amount must be greater than 0")
// 		log.Println("Error:", err)
// 		return "", err
// 	}
// 	if payment.Currency == "" {
// 		err := fmt.Errorf("payment currency is required")
// 		log.Println("Error:", err)
// 		return "", err
// 	}
// 	successURLBase := os.Getenv("PAYMENT_SUCCESS_URL")
// 	cancelURL := os.Getenv("PAYMENT_CANCEL_URL")
// 	successRedirect := os.Getenv("PAYMENT_SUCCESS_REDIRECT")
// 	cancelRedirect := os.Getenv("PAYMENT_CANCEL_REDIRECT")

// 	if successURLBase == "" || cancelURL == "" || successRedirect == "" || cancelRedirect == "" {
// 		return "", fmt.Errorf("missing Stripe URL env vars")
// 	}

// 	params := &stripe.CheckoutSessionParams{
// 		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
// 		LineItems: []*stripe.CheckoutSessionLineItemParams{
// 			{
// 				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
// 					Currency: stripe.String(payment.Currency),
// 					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
// 						Name: stripe.String("Subscription Plan"),
// 					},
// 					UnitAmount: stripe.Int64(int64(payment.Amount * 100)),
// 				},
// 				Quantity: stripe.Int64(1),
// 			},
// 		},
// 		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
// 		// Stripe will replace {CHECKOUT_SESSION_ID} automatically
// 		SuccessURL: stripe.String(successURLBase + "?session_id={CHECKOUT_SESSION_ID}"),
// 		CancelURL:  stripe.String(cancelURL + "?session_id={CHECKOUT_SESSION_ID}"),
// 	}
// 	params.AddMetadata("transaction_id", payment.TransactionID.String())
// 	params.AddMetadata("payment_id", payment.ID.String())
// 	params.AddMetadata("user_id", payment.UserID.String()) // safe way for pgtype.UUID
// 	s, err := session.New(params)
// 	if err != nil {
// 		log.Printf("[ERROR] Failed to create Stripe session: %v", err)
// 		return "", fmt.Errorf("failed to create Stripe session: %w", err)
// 	}

// 	log.Printf("Stripe Checkout session created – ID: %s, URL: %s", s.ID, s.URL)
// 	return s.URL, nil
// }
