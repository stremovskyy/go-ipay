/*
 * MIT License
 *
 * Copyright (c) 2024 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"

	"github.com/google/uuid"

	go_ipay "github.com/stremovskyy/go-ipay"
	"github.com/stremovskyy/go-ipay/currency"
	"github.com/stremovskyy/go-ipay/examples/internal/config"
	"github.com/stremovskyy/go-ipay/internal/utils"
	"github.com/stremovskyy/go-ipay/log"
)

func main() {
	cfg := config.MustLoad()
	client := go_ipay.NewDefaultClient()

	merchant := &go_ipay.Merchant{
		Name:            cfg.MerchantName,
		SubMerchantID:   cfg.SubMerchantID,
		MerchantID:      cfg.MerchantID,
		MerchantKey:     cfg.MerchantKey,
		SuccessRedirect: cfg.SuccessRedirect,
		FailRedirect:    cfg.FailRedirect,
	}

	uuidString := uuid.New().String()
	metaData := map[string]string{
		"key1": "value1",
	}

	holdRequest := &go_ipay.Request{
		Merchant: merchant,
		PaymentMethod: &go_ipay.PaymentMethod{
			Card: &go_ipay.Card{
				Name:  "Test Card",
				Token: utils.Ref(cfg.CardToken),
			},
		},
		PaymentData: &go_ipay.PaymentData{
			IpayPaymentID: utils.Ref(int64(cfg.IpayPaymentID)),
			PaymentID:     utils.Ref(uuidString),
			Amount:        100,
			Currency:      currency.UAH,
			OrderID:       uuidString,
			Description:   "Test payment: " + uuidString,
			Metadata:      metaData,
		},
		PersonalData: &go_ipay.PersonalData{
			UserID:     utils.Ref(123),
			FirstName:  utils.Ref("John"),
			LastName:   utils.Ref("Doe"),
			MiddleName: utils.Ref("Middle"),
			TaxID:      utils.Ref("1234567890"),
		},
	}

	client.SetLogLevel(log.LevelDebug)

	holdRequest.SetWebhookURL(utils.Ref(cfg.WebhookURL))

	holdResponse, err := client.Hold(holdRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	holdResponse.PrettyPrint()
}
