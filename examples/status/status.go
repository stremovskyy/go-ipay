package main

import (
	go_ipay "github.com/megakit-pro/go-ipay"
	"github.com/megakit-pro/go-ipay/internal/log"
	"github.com/megakit-pro/go-ipay/private"
)

func main() {
	client := go_ipay.NewDefaultClient()

	merchant := &go_ipay.Merchant{
		Name:            private.MerchantName,
		MerchantID:      private.MerchantID,
		MerchantKey:     private.MerchantKey,
		SystemKey:       private.SystemKey,
		SuccessRedirect: private.SuccessRedirect,
		FailRedirect:    private.FailRedirect,
	}

	request := &go_ipay.Request{
		Merchant: merchant,
		PaymentData: &go_ipay.PaymentData{
			IpayPaymentID: private.IpayPaymentID,
		},
	}

	log.SetLevel(log.LevelDebug)

	status, err := client.Status(request)
	if err != nil {
		panic(err)
	}

	println(status)
}
