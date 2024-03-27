package go_ipay

import (
	"strconv"

	"github.com/stremovskyy/go-ipay/internal/ipay"
)

type Merchant struct {
	// Merchant Name
	Name string
	// Merchant ID
	MerchantID string
	// Merchant Key
	MerchantKey string
	// System Key
	SystemKey string

	// SuccessRedirect
	SuccessRedirect string

	// FailRedirect
	FailRedirect string

	signer ipay.Signer
}

func (m *Merchant) GetMerchantID() int64 {
	id, err := strconv.ParseInt(m.MerchantID, 10, 64)

	if err != nil {
		return 0
	}

	return id
}

func (m *Merchant) GetSign() ipay.Sign {
	if m.signer == nil {
		m.signer = ipay.NewSigner(m.SystemKey)
	}

	return *m.signer.Sign(m.MerchantKey)
}
