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

package go_ipay

import "github.com/stremovskyy/go-ipay/currency"

// PaymentData represents the data related to a payment transaction.
type PaymentData struct {
	// IpayPaymentID is the unique identifier for the iPay payment.
	IpayPaymentID *int64
	// PaymentID is the unique identifier for the payment.
	PaymentID *string
	// Amount is the amount of the payment in the smallest unit of the currency.
	Amount int
	// Currency is the currency code of the payment.
	Currency currency.Code
	// OrderID is the unique identifier for the order.
	OrderID string
	// Description is a brief description of the payment.
	Description string
	// WebhookURL is the URL to which payment notifications will be sent.
	WebhookURL *string
	// IsMobile indicates whether the payment was made from a mobile device.
	IsMobile bool
	// RelatedIds is a list of related payment IDs.
	RelatedIds []int64
}
