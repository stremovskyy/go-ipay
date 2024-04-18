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
	"time"

	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/go-ipay/private"
)

func main() {
	xmlData := []byte(private.ValidationXMLSuccess)

	payment, err := ipay.ParsePaymentXML(xmlData)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}

	// For demonstration, print out the payment ID and status
	fmt.Printf("PaymentURL ID: %d, Status: %s, Card Token: %v, Card Is Prepaid: %s\n", payment.ID, payment.Status.String(), payment.CardToken, payment.CardIsPrepaid)

	// To demonstrate timestamp conversion to a readable format
	timestamp := time.Unix(payment.Timestamp, 0)
	fmt.Println("Timestamp:", timestamp)
}
