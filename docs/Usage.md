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
    var providerErr *ipay.IpayError
    if errors.As(err, &providerErr) {
        fmt.Printf("iPay rejected the request with code %d\n", providerErr.Code)
    } else {
        fmt.Printf("Payment call failed: %v\n", err)
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

GO-iPay separates valid iPay business errors from transport and invalid-response
failures. The public sentinels support `errors.Is`, and the typed errors support
`errors.As`:

```go
import (
    "errors"
    "fmt"

    "github.com/stremovskyy/go-ipay/ipay"
)

func classifyError(err error) error {
    var providerErr *ipay.IpayError
    var transportErr *ipay.TransportError
    var responseErr *ipay.UnexpectedResponseError
    var decodeErr *ipay.DecodeError

    switch {
    case errors.As(err, &providerErr):
        return fmt.Errorf("iPay business error %d: %w", providerErr.Code, err)
    case errors.As(err, &transportErr):
        return fmt.Errorf("iPay transport failed (request %s): %w", transportErr.RequestID, err)
    case errors.As(err, &responseErr):
        // Body is a bounded diagnostic snapshot. Keep it out of ordinary logs.
        return fmt.Errorf("unexpected iPay HTTP status %d (request %s): %w", responseErr.StatusCode, responseErr.RequestID, err)
    case errors.As(err, &decodeErr):
        return fmt.Errorf("invalid iPay JSON (request %s): %w", decodeErr.RequestID, err)
    default:
        return err
    }
}
```

### A2C outcome reconciliation

Do not automatically retry `Credit` (`A2CPay`) when the request outcome is
unknown. Reusing the same payout as a new request could duplicate the transfer.
Check the original external ID with `A2CPaymentStatus` first:

```go
response, err := client.Credit(request)
if ipay.RequiresA2CStatusCheck(err) {
    extID := request.PaymentData.PaymentID
    response, err = client.A2CPaymentStatus(&go_ipay.Request{
        Merchant: request.Merchant,
        PaymentData: &go_ipay.PaymentData{PaymentID: extID},
    })
}
```

`RequiresA2CStatusCheck` is false for structured iPay business refusals and for
errors returned by the status-check request itself.

## Best Practices

1. **Error Handling**
   - Always check for errors and handle them appropriately
   - Reconcile outcome-unknown A2C payouts by `ext_id`; never retry them automatically
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

4. **Timeouts and A2C Safety**
   ```go
   client := go_ipay.NewClient(
       go_ipay.WithClient(&http.Client{Timeout: 30 * time.Second}),
   )
   ```
   Configure the HTTP timeout explicitly, but do not retry an outcome-unknown
   A2C payout. Run `A2CPaymentStatus` with the original `ext_id` first.

5. **Testing**
   - Use test credentials in development
   - Implement proper error handling
   - Test various payment scenarios
   - Verify webhook handling
