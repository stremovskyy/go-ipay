package main

import (
	go_ipay "github.com/megakit-pro/go-ipay"
	"github.com/megakit-pro/go-ipay/internal/log"
	"github.com/megakit-pro/go-ipay/internal/utils"
	"github.com/megakit-pro/go-ipay/private"
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

	request := &go_ipay.Request{
		Merchant: merchant,
		PersonalData: &go_ipay.PersonalData{
			UserID:    utils.Ref(123),
			FirstName: utils.Ref("John"),
			LastName:  utils.Ref("Doe"),
			TaxID:     utils.Ref("1234567890"),
		},
	}

	log.SetLevel(log.LevelDebug)

	tokenURL, err := client.VerificationLink(request)
	if err != nil {
		panic(err)
	}

	println(tokenURL.String())
}
