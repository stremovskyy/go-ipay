package go_ipay

import (
	"net/url"

	"github.com/stremovskyy/go-ipay/internal/ipay"
)

type Ipay interface {
	VerificationLink(request *Request) (*url.URL, error)
	Status(request *Request) (*ipay.Response, error)
	PaymentURL(invoiceRequest *Request) (*ipay.PaymentResponse, error)
	Payment(invoiceRequest *Request) (*ipay.Response, error)
	Hold(invoiceRequest *Request) (*ipay.Response, error)
	Capture(invoiceRequest *Request) (*ipay.Response, error)
	Refund(invoiceRequest *Request) (*ipay.Response, error)
	Credit(invoiceRequest *Request) (*ipay.Response, error)
}
