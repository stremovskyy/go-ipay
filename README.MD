# GO-iPay

[![GoDoc](https://godoc.org/github.com/stremovskyy/go-ipay?status.svg)](https://godoc.org/github.com/stremovskyy/go-ipay)
[![Go Report Card](https://goreportcard.com/badge/github.com/stremovskyy/go-ipay)](https://goreportcard.com/report/github.com/stremovskyy/go-ipay)
[![License](https://img.shields.io/github/license/stremovskyy/go-ipay)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/stremovskyy/go-ipay)](go.mod)

A Go client library for the iPay.ua Payment Gateway API. This library provides a simple and reliable way to integrate iPay.ua payment processing into your Go applications.

## Features

- Simple and intuitive API
- Comprehensive error handling
- Full iPay.ua API coverage
- Thread-safe client implementation
- Configurable HTTP client and timeouts
- Detailed logging support

## Installation

```bash
go get github.com/stremovskyy/go-ipay
```

## Documentation

 **Complete documentation is available in:**
- [Usage Guide](docs/Usage.md) - Comprehensive guide to using the library
- [Code Examples](docs/Examples.md) - Real-world examples and best practices

For API reference and package documentation:
- [GoDoc Documentation](https://godoc.org/github.com/stremovskyy/go-ipay)

## Quick Start

```go
package main

import (
    "github.com/stremovskyy/go-ipay"
)

func main() {
    // Create a new client with default options
    client := go_ipay.NewDefaultClient()

    // Configure merchant details
    merchant := &go_ipay.Merchant{
        Name:        "Your Merchant Name",
        MerchantID:  "your-merchant-id",
        MerchantKey: "your-merchant-key",
    }

    // Use the client...
}
```

## Configuration Options

The client can be configured with various options:

```go
client := go_ipay.NewClient(
    go_ipay.WithClient(&http.Client{Timeout: 30 * time.Second}),
)
```

## Error Handling

The `ipay` package exposes typed transport, HTTP-response, decode, and provider
business errors. Use `errors.Is` for error kinds and `errors.As` for diagnostic
metadata. Diagnostic response bodies are bounded and are deliberately excluded
from `Error()`.

```go
import (
    "errors"

    go_ipay "github.com/stremovskyy/go-ipay"
    "github.com/stremovskyy/go-ipay/ipay"
)

if err != nil {
    var providerErr *ipay.IpayError
    var responseErr *ipay.UnexpectedResponseError

    switch {
    case errors.As(err, &providerErr):
        // A structured business response from iPay.
    case errors.As(err, &responseErr):
        // Inspect status/content type/request ID; do not log responseErr.Body.
    case errors.Is(err, ipay.ErrTransport):
        // The HTTP exchange failed.
    case errors.Is(err, ipay.ErrDecode):
        // The response looked like JSON but was malformed.
    }
}
```

An `A2CPay` transport or response failure may leave the payout outcome unknown.
Never retry it automatically. Reconcile the original `ext_id` with
`A2CPaymentStatus` first:

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please make sure to update tests as appropriate.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

For security issues, please email security@yourdomain.com instead of using the issue tracker.

## Acknowledgments

- iPay.ua for providing the payment gateway service
- Contributors who have helped improve this library

## Support

If you are having problems, please let us know by [raising a new issue](https://github.com/stremovskyy/go-ipay/issues/new).
