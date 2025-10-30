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

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	go_ipay "github.com/stremovskyy/go-ipay"
	"github.com/stremovskyy/go-ipay/examples/internal/config"
	"github.com/stremovskyy/go-ipay/internal/utils"
	"github.com/stremovskyy/go-ipay/log"
	"github.com/stremovskyy/recorder"
	"github.com/stremovskyy/recorder/gorm_recorder"
)

func main() {
	sqlDb, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to open sqlite DB:", err)
		return
	}

	err = sqlDb.AutoMigrate(&IpayRequestRecord{}, &IpayRequestTag{})
	if err != nil {
		fmt.Println("Failed to migrate DB:", err)
		return
	}

	opts := gorm_recorder.NewOptions(func() *IpayRequestRecord { return &IpayRequestRecord{} }, func() *IpayRequestTag { return &IpayRequestTag{} }).
		WithRecordTable("custom_records").
		WithRecordColumns("id", "kind", "correlation_id").
		WithTagTable("custom_tags").
		WithTagColumns("record_ref", "tag_key", "tag_value")

	scrub := recorder.NewScrubber(
		recorder.WithDefaultReplacement("<hidden>"),
	)

	rec, err := gorm_recorder.NewRecorderWithModels(sqlDb, opts, recorder.WithScrubber(scrub))
	if err != nil {
		fmt.Println("Failed to create recorder:", err)
		return
	}

	cfg := config.MustLoad()
	client := go_ipay.NewClientWithRecorder(rec)

	merchant := &go_ipay.Merchant{
		Name:        cfg.MerchantName,
		MerchantID:  cfg.MerchantID,
		MerchantKey: cfg.MerchantKey,
	}

	statusRequest := &go_ipay.Request{
		Merchant: merchant,
		PaymentData: &go_ipay.PaymentData{
			IpayPaymentID: utils.Ref(int64(65465465)),
		},
	}

	client.SetLogLevel(log.LevelDebug)
	// statusRequest.SetWebhookURL(utils.Ref(cfg.WebhookURL))

	statusResponse, err := client.Status(statusRequest)
	if err != nil {
		if statusResponse != nil {
			statusResponse.PrettyPrint()

			return
		}

		fmt.Println(err)
		return
	}

	statusResponse.PrettyPrint()
}
