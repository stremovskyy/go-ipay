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
	"os"

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
		MerchantID:   cfg.MerchantIDWithdraw,
		Login:        cfg.Login,
		RepaymentKey: cfg.RepaymentKey,
	}

	repaymentExtID := uuid.New().String()

	req := &go_ipay.CreateRepaymentRequest{
		Merchant: merchant,
		ExtID:    repaymentExtID,
		Transactions: []go_ipay.RepaymentTransaction{
			{
				PmtID: int64(cfg.IpayPaymentID),
				ExtID: uuid.New().String(),
			},
		},
	}

	client.SetLogLevel(log.LevelDebug)

	resp, err := client.CreateRepayment(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateRepayment error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Repayment GUID: %s\n", utils.SafeString(resp.RepaymentGUID))
	fmt.Printf("Repayment ext_id: %s\n", utils.SafeString(resp.ExtID))
	fmt.Printf("Status: %d\n", utils.SafeInt(resp.Status))
	fmt.Printf("Invoice: %d\n", utils.SafeInt(resp.Invoice))
	fmt.Printf("Amount: %d\n", utils.SafeInt(resp.Amount))
	if resp.MchID != nil {
		fmt.Printf("mch_id: %d\n", *resp.MchID)
	}
	fmt.Printf("mch_balance: %d\n", utils.SafeInt(resp.MchBalance))
	fmt.Printf("success_payments: %d\n", utils.SafeInt(resp.SuccessPayments))
	fmt.Printf("failed_payments: %d\n", utils.SafeInt(resp.FailedPayments))
}
