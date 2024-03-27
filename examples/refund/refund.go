package main

import (
	go_ipay "github.com/stremovskyy/go-ipay"
	"github.com/stremovskyy/go-ipay/internal/log"
	"github.com/stremovskyy/go-ipay/internal/utils"
	"github.com/stremovskyy/go-ipay/private"
)

func main() {
	client := go_ipay.NewDefaultClient()

	merchant := &go_ipay.Merchant{
		Name:        private.MerchantName,
		MerchantID:  private.MerchantID,
		MerchantKey: private.MerchantKey,
	}

	request := &go_ipay.Request{
		Merchant: merchant,
		PaymentData: &go_ipay.PaymentData{
			IpayPaymentID: utils.Ref(int64(private.IpayPaymentID)),
		},
	}

	log.SetLevel(log.LevelDebug)

	status, err := client.Refund(request)
	if err != nil {
		panic(err)
	}

	println(status)
}
