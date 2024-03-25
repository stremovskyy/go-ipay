package go_ipay

import (
	"net/url"

	"github.com/megakit-pro/go-ipay/ipay"
)

type Ipay interface {
	VerificationLink(request *Request) (*url.URL, error)
	Status(request *Request) (*ipay.Response, error)
}
