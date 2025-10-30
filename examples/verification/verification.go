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
	"github.com/google/uuid"

	go_ipay "github.com/stremovskyy/go-ipay"
	"github.com/stremovskyy/go-ipay/examples/internal/config"
	"github.com/stremovskyy/go-ipay/internal/utils"
	"github.com/stremovskyy/go-ipay/log"
)

func main() {
	cfg := config.MustLoad()
	client := go_ipay.NewDefaultClient()

	merchant := &go_ipay.Merchant{
		Name:            cfg.MerchantName,
		MerchantID:      cfg.MerchantID,
		SubMerchantID:   cfg.SubMerchantID,
		MerchantKey:     cfg.MerchantKey,
		SuccessRedirect: cfg.SuccessRedirect,
		FailRedirect:    cfg.FailRedirect,
	}

	uuidString := uuid.New().String()

	VerificationRequest := &go_ipay.Request{
		Merchant: merchant,
		PaymentData: &go_ipay.PaymentData{
			IpayPaymentID: utils.Ref(int64(cfg.IpayPaymentID)),
			PaymentID:     utils.Ref(uuidString),
			OrderID:       uuidString,
			Description:   "Verification payment: " + uuidString,
		},
		PersonalData: &go_ipay.PersonalData{
			UserID:    utils.Ref(123),
			FirstName: utils.Ref("John"),
			LastName:  utils.Ref("Doe"),
			TaxID:     utils.Ref("1234567890"),
		},
	}

	client.SetLogLevel(log.LevelDebug)
	VerificationRequest.SetWebhookURL(utils.Ref(cfg.WebhookURL))

	tokenURL, err := client.VerificationLink(VerificationRequest)
	if err != nil {
		panic(err)
	}

	println(tokenURL.String())
}
