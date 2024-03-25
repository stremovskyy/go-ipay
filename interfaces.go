package go_ipay

import "net/url"

type Ipay interface {
	VerificationLink(*Request) (*url.URL, error)
}
