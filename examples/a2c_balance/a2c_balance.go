/*
 * MIT License
 *
 * Copyright (c) 2026 Anton Stremovskyy
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

	go_ipay "github.com/stremovskyy/go-ipay"
	"github.com/stremovskyy/go-ipay/examples/internal/config"
)

func main() {
	cfg := config.MustLoadA2CBalance()
	client := go_ipay.NewDefaultClient()

	balance, err := client.A2CBalance(&go_ipay.Request{
		Merchant: &go_ipay.Merchant{
			Name:        cfg.MerchantNameWithdraw,
			MerchantID:  cfg.MerchantIDWithdraw,
			MerchantKey: cfg.MerchantKeyWithdraw,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "A2C balance error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("A2C balance values are in kopecks; current_balance may be negative.")
	fmt.Printf("current_balance: %d\n", balance.CurrentBalance)
	fmt.Printf("overdraft: %d\n", balance.Overdraft)
	fmt.Printf("credit: %d\n", balance.Credit)
}
