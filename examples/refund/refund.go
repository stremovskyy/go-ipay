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

	go_ipay "github.com/stremovskyy/go-ipay"
	"github.com/stremovskyy/go-ipay/internal/utils"
	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/go-ipay/log"
	"github.com/stremovskyy/go-ipay/private"
)

func main() {
	client := go_ipay.NewDefaultClient()

	merchant := &go_ipay.Merchant{
		Name:        private.MerchantName,
		MerchantID:  private.MerchantID,
		MerchantKey: private.MerchantKey,
	}

	refundRequest := &go_ipay.Request{
		Merchant: merchant,
		PaymentData: &go_ipay.PaymentData{
			IpayPaymentID: utils.Ref(int64(376470216)),
		},
	}

	client.SetLogLevel(log.LevelDebug)
	refundRequest.SetWebhookURL(utils.Ref(private.WebhookURL))

	refundResponse, err := client.Refund(refundRequest)
	if err != nil {
		panic(err)
	}

	if refundResponse.Status != ipay.PaymentStatusCanceled {
		panic(refundResponse.GetError())
	}

	fmt.Printf("Refund success: %s\n", refundResponse.Status)
}
