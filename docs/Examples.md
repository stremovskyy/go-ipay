# GO-iPay Examples

This document contains real-world examples for common payment scenarios using GO-iPay.

## Table of Contents

- [Simple Payment](#simple-payment)
- [Card Payment](#card-payment)
- [Apple Pay](#apple-pay)
- [Google Pay](#google-pay)
- [Payment with Hold](#payment-with-hold)
- [Payment Status Check](#payment-status-check)
- [Refund Processing](#refund-processing)
- [Webhook Handler](#webhook-handler)

## Simple Payment

Complete example of processing a simple payment:

```go
package main

import (
    "fmt"
    "log"
    
    go_ipay "github.com/stremovskyy/go-ipay"
    "github.com/stremovskyy/go-ipay/currency"
)

func main() {
    // Create client
    client := go_ipay.NewDefaultClient()
    
    // Configure merchant
    merchant := &go_ipay.Merchant{
        Name:            "My Store",
        MerchantID:      "merchant-id",
        MerchantKey:     "merchant-key",
        SuccessRedirect: "https://mystore.com/success",
        FailRedirect:    "https://mystore.com/fail",
    }
    
    // Create order ID
    orderID := fmt.Sprintf("order-%d", time.Now().Unix())
    
    // Create payment request
    request := &go_ipay.Request{
        Merchant: merchant,
        PaymentData: &go_ipay.PaymentData{
            Amount:      999.99,
            Currency:    currency.UAH,
            OrderID:     orderID,
            Description: "Order #" + orderID,
        },
        PersonalData: &go_ipay.PersonalData{
            FirstName: utils.Ref("John"),
            LastName:  utils.Ref("Doe"),
            Email:     utils.Ref("john@example.com"),
            Phone:     utils.Ref("+380501234567"),
        },
    }
    
    // Enable debug logging
    client.SetLogLevel(log.LevelDebug)
    
    // Process payment
    response, err := client.Payment(request)
    if err != nil {
        log.Fatalf("Payment error: %v", err)
    }
    
    // Use payment URL
    fmt.Printf("Please complete payment at: %s\n", response.PaymentURL)
}
```

## Card Payment

Example of processing a payment with a saved card:

```go
package main

import (
    "fmt"
    "log"
    
    go_ipay "github.com/stremovskyy/go-ipay"
    "github.com/stremovskyy/go-ipay/currency"
)

func main() {
    client := go_ipay.NewDefaultClient()
    
    // Configure payment method
    paymentMethod := &go_ipay.PaymentMethod{
        Card: &go_ipay.Card{
            Name:  "My Visa Card",
            Token: utils.Ref("saved-card-token"),
        },
    }
    
    // Create request
    request := &go_ipay.Request{
        Merchant:      merchant,
        PaymentMethod: paymentMethod,
        PaymentData: &go_ipay.PaymentData{
            Amount:      500.00,
            Currency:    currency.UAH,
            OrderID:     "card-payment-123",
            Description: "Payment with saved card",
        },
    }
    
    // Process payment
    response, err := client.Payment(request)
    if err != nil {
        log.Fatalf("Card payment error: %v", err)
    }
    
    fmt.Printf("Payment processed: %s\n", response.PaymentID)
}
```

## Payment with Hold

Example of creating a payment with hold (authorization without capture):

```go
package main

import (
    "fmt"
    "log"
    
    go_ipay "github.com/stremovskyy/go-ipay"
    "github.com/stremovskyy/go-ipay/currency"
)

func main() {
    client := go_ipay.NewDefaultClient()
    
    // Create hold request
    request := &go_ipay.Request{
        Merchant: merchant,
        PaymentData: &go_ipay.PaymentData{
            Amount:      1000.00,
            Currency:    currency.UAH,
            OrderID:     "hold-123",
            Description: "Payment with hold",
            Hold:        utils.Ref(true), // Enable hold
        },
    }
    
    // Create hold
    response, err := client.Hold(request)
    if err != nil {
        log.Fatalf("Hold creation error: %v", err)
    }
    
    // Later, capture the held payment
    captureRequest := &go_ipay.CaptureRequest{
        Merchant:  merchant,
        PaymentID: response.PaymentID,
        Amount:    1000.00, // Can be less than original amount
    }
    
    capture, err := client.Capture(captureRequest)
    if err != nil {
        log.Fatalf("Capture error: %v", err)
    }
    
    fmt.Printf("Payment captured: %s\n", capture.PaymentID)
}
```

## Webhook Handler

Example of handling payment webhooks:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    
    go_ipay "github.com/stremovskyy/go-ipay"
)

func main() {
    http.HandleFunc("/webhook", handleWebhook)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    // Parse webhook
    webhook, err := go_ipay.ParseWebhook(r)
    if err != nil {
        log.Printf("Webhook parse error: %v", err)
        http.Error(w, "Invalid webhook", http.StatusBadRequest)
        return
    }
    
    // Process based on event type
    switch webhook.Event {
    case "payment.success":
        handleSuccessfulPayment(webhook)
    case "payment.failed":
        handleFailedPayment(webhook)
    case "payment.refunded":
        handleRefundedPayment(webhook)
    default:
        log.Printf("Unknown webhook event: %s", webhook.Event)
    }
    
    w.WriteHeader(http.StatusOK)
}

func handleSuccessfulPayment(webhook *go_ipay.Webhook) {
    // Update order status
    err := updateOrderStatus(webhook.OrderID, "paid")
    if err != nil {
        log.Printf("Error updating order %s: %v", webhook.OrderID, err)
        return
    }
    
    // Send confirmation email
    err = sendConfirmationEmail(webhook.OrderID)
    if err != nil {
        log.Printf("Error sending confirmation for order %s: %v", webhook.OrderID, err)
    }
}

func handleFailedPayment(webhook *go_ipay.Webhook) {
    // Update order status
    err := updateOrderStatus(webhook.OrderID, "failed")
    if err != nil {
        log.Printf("Error updating failed order %s: %v", webhook.OrderID, err)
        return
    }
    
    // Notify customer
    err = sendPaymentFailureNotification(webhook.OrderID)
    if err != nil {
        log.Printf("Error sending failure notification for order %s: %v", webhook.OrderID, err)
    }
}
```

These examples demonstrate common use cases and best practices for integrating GO-iPay into your application. Remember to:

1. Always handle errors appropriately
2. Log important events and errors
3. Validate input data
4. Secure sensitive information
5. Implement proper webhook handling
6. Test thoroughly with test credentials before going live
