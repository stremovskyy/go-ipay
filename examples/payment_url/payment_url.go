package main

import (
	"github.com/google/uuid"
	go_ipay "github.com/megakit-pro/go-ipay"
	"github.com/megakit-pro/go-ipay/internal/log"
	"github.com/megakit-pro/go-ipay/private"
	"github.com/megakit-pro/go-ipay/utils"
)

func main() {
	client := go_ipay.NewDefaultClient()

	merchant := &go_ipay.Merchant{
		Name:            private.MerchantName,
		MerchantID:      private.MerchantID,
		MerchantKey:     private.MerchantKey,
		SuccessRedirect: private.SuccessRedirect,
		FailRedirect:    private.FailRedirect,
	}

	paymentRequest := &go_ipay.Request{
		Merchant: merchant,
		PaymentMethod: &go_ipay.PaymentMethod{
			Card: &go_ipay.Card{
				Name:  "Test Card",
				Token: utils.Ref(private.CardToken),
			},
		},
		PaymentData: &go_ipay.PaymentData{
			IpayPaymentID: utils.Ref(int64(private.IpayPaymentID)),
			PaymentID:     utils.Ref(uuid.New().String()),
			Amount:        100,
			Currency:      "UAH",
			OrderID:       uuid.New().String(),
			Description:   "Test payment URL",
		},
		PersonalData: &go_ipay.PersonalData{
			UserID:    utils.Ref(123),
			FirstName: utils.Ref("John"),
			LastName:  utils.Ref("Doe"),
			TaxID:     utils.Ref("1234567890"),
		},
	}

	log.SetLevel(log.LevelDebug)

	paymentResponse, err := client.PaymentURL(paymentRequest)
	if err != nil {
		panic(err)
	}

	println(paymentResponse)
}
