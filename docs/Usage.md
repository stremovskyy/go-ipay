# GO-iPay Usage Guide

This guide provides comprehensive examples for using the GO-iPay library to integrate with the iPay.ua payment gateway.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Basic Payment Flow](#basic-payment-flow)
- [Advanced Features](#advanced-features)
  - [Card Payments](#card-payments)
  - [Apple Pay](#apple-pay)
  - [Google Pay](#google-pay)
  - [Run Options](#run-options)
  - [Payment Status](#payment-status)
  - [Refunds](#refunds)
  - [Webhooks](#webhooks)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## Installation

```bash
go get github.com/stremovskyy/go-ipay
```

## Quick Start

Here's a minimal example to get started with GO-iPay:

```go
package main

import (
    "fmt"
    go_ipay "github.com/stremovskyy/go-ipay"
    "github.com/stremovskyy/go-ipay/currency"
)

func main() {
    // Create a new client
    client := go_ipay.NewDefaultClient()

    // Configure merchant details
    merchant := &go_ipay.Merchant{
        Name:            "Your Store Name",
        MerchantID:      "your-merchant-id",
        MerchantKey:     "your-merchant-key",
        SuccessRedirect: "https://your-store.com/success",
        FailRedirect:    "https://your-store.com/fail",
    }

    // Create a payment request
    request := &go_ipay.Request{
        Merchant: merchant,
        PaymentData: &go_ipay.PaymentData{
            Amount:      100.00,
            Currency:    currency.UAH,
            OrderID:     "order-123",
            Description: "Test payment",
        },
    }

    // Process payment
    response, err := client.Payment(request)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Payment URL: %s\n", response.PaymentURL)
}
```

## Basic Payment Flow

### 1. Create a Payment

```go
// Initialize payment data
paymentData := &go_ipay.PaymentData{
    Amount:      1000.00,
    Currency:    currency.UAH,
    OrderID:     "order-123",
    Description: "Product purchase",
}

// Add customer information
personalData := &go_ipay.PersonalData{
    FirstName: utils.Ref("John"),
    LastName:  utils.Ref("Doe"),
    Email:     utils.Ref("john@example.com"),
    Phone:     utils.Ref("+380501234567"),
}

// Create the request
request := &go_ipay.Request{
    Merchant:     merchant,
    PaymentData:  paymentData,
    PersonalData: personalData,
}

// Process payment
response, err := client.Payment(request)
```

### 2. Handle the Response

```go
if err != nil {
    switch e := err.(type) {
    case *go_ipay.ValidationError:
        fmt.Printf("Validation error: %v\n", e)
    case *go_ipay.APIError:
        fmt.Printf("API error: %v\n", e)
    default:
        fmt.Printf("Unknown error: %v\n", e)
    }
    return
}

// Use the payment URL or token
fmt.Printf("Payment URL: %s\n", response.PaymentURL)
fmt.Printf("Payment Token: %s\n", response.PaymentToken)
```

## Advanced Features

### Card Payments

Process payment with saved card:

```go
paymentMethod := &go_ipay.PaymentMethod{
    Card: &go_ipay.Card{
        Name:  "My Saved Card",
        Token: utils.Ref("saved-card-token"),
    },
}

request := &go_ipay.Request{
    Merchant:      merchant,
    PaymentMethod: paymentMethod,
    PaymentData:   paymentData,
}

response, err := client.Payment(request)
```

### Payment Status

Check payment status:

```go
status, err := client.Status(&go_ipay.StatusRequest{
    Merchant:  merchant,
    PaymentID: "payment-123",
})

if err != nil {
    fmt.Printf("Error checking status: %v\n", err)
    return
}

fmt.Printf("Payment Status: %s\n", status.Status)
fmt.Printf("Amount: %.2f %s\n", status.Amount, status.Currency)
```

### Run Options

All client calls accept optional run options. Use them to adjust behaviour per request, for example to perform a dry run without contacting the API while inspecting the payload that would be sent.

> Tip: when you call `options.DryRun()` without a custom handler, the library pretty-prints the request payload to the logger. Make sure to enable logging (`client.SetLogLevel(log.LevelInfo)`) to see it.

```go
import (
    "fmt"

    go_ipay "github.com/stremovskyy/go-ipay"
    "github.com/stremovskyy/go-ipay/options"
    "github.com/stremovskyy/go-ipay/ipay"
)

client := go_ipay.NewDefaultClient()

var wrapper *ipay.RequestWrapper

_, err := client.Hold(request, options.DryRun(func(endpoint string, payload any) {
    fmt.Println("Skipping request to:", endpoint)

    if v, ok := payload.(*ipay.RequestWrapper); ok {
        wrapper = v
    }
}))

if err != nil {
    fmt.Printf("Hold error: %v\n", err)
    return
}

fmt.Printf("Operation: %s\n", wrapper.Operation)
```

### Refunds

Process a refund:

```go
refundRequest := &go_ipay.RefundRequest{
    Merchant:  merchant,
    PaymentID: "payment-123",
    Amount:    100.00, // Partial refund amount
    Reason:    "Customer request",
}

refund, err := client.Refund(refundRequest)
if err != nil {
    fmt.Printf("Refund error: %v\n", err)
    return
}

fmt.Printf("Refund ID: %s\n", refund.RefundID)
```

### Webhooks

Handle payment notifications:

```go
func handleWebhook(w http.ResponseWriter, r *http.Request) {
    webhook, err := go_ipay.ParseWebhook(r)
    if err != nil {
        http.Error(w, "Invalid webhook", http.StatusBadRequest)
        return
    }

    switch webhook.Event {
    case "payment.success":
        // Handle successful payment
        processSuccessfulPayment(webhook.PaymentID)
    case "payment.failed":
        // Handle failed payment
        processFailedPayment(webhook.PaymentID)
    }

    w.WriteHeader(http.StatusOK)
}
```

## Error Handling

GO-iPay provides detailed error types:

```go
func processPayment(request *go_ipay.Request) error {
    response, err := client.Payment(request)
    if err != nil {
        switch e := err.(type) {
        case *go_ipay.ValidationError:
            // Handle validation errors (e.g., invalid amount)
            log.Printf("Validation error: %v", e)
            return fmt.Errorf("invalid payment data: %w", e)
            
        case *go_ipay.APIError:
            // Handle API errors (e.g., authentication failed)
            log.Printf("API error: %v", e)
            return fmt.Errorf("payment processing failed: %w", e)
            
        default:
            // Handle unexpected errors
            log.Printf("Unexpected error: %v", e)
            return fmt.Errorf("payment system error: %w", e)
        }
    }
    
    return nil
}
```

## Best Practices

1. **Error Handling**
   - Always check for errors and handle them appropriately
   - Log errors with sufficient context
   - Provide meaningful error messages to users

2. **Security**
   - Never store merchant keys in code
   - Use environment variables or secure configuration management
   - Validate all input data before sending to the API

3. **Logging**
   ```go
   // Enable debug logging
   client.SetLogLevel(log.LevelDebug)
   
   // Available log levels
   log.LevelNone    // Disables logging
   log.LevelError   // Logs errors only
   log.LevelWarning // Logs warnings and errors
   log.LevelInfo    // Logs info, warnings, and errors
   log.LevelDebug   // Logs everything including detailed IO
   ```

4. **Timeouts and Retries**
   ```go
   client := go_ipay.NewClient(
       go_ipay.WithTimeout(30 * time.Second),
       go_ipay.WithRetryAttempts(3),
   )
   ```

5. **Testing**
   - Use test credentials in development
   - Implement proper error handling
   - Test various payment scenarios
   - Verify webhook handling
