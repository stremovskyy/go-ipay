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
)

const sampleWebhookResponse = `
<payment id="123456">
	<ident>demo-transaction</ident>
	<status>4</status>
	<amount>100.00</amount>
	<currency>UAH</currency>
	<timestamp>1700000000</timestamp>
	<transactions>
		<transaction id="1699368014">
			<mch_id>13581</mch_id>
			<srv_id>4767</srv_id>
			<invoice>100</invoice>
			<amount>100</amount>
			<desc>Demo payment</desc>
			<info>{"preauth":1,"notify_url":"https://example.test"}</info>
		</transaction>
	</transactions>
	<salt>demo-salt</salt>
	<sign>demo-signature</sign>
	<pmt_id>123</pmt_id>
</payment>
`

func main() {
	xmlData := []byte(sampleWebhookResponse)

	payment, err := ipay.ParsePaymentXML(xmlData)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}

	// For demonstration, print out the payment ID and status
	fmt.Printf("Payment: %s\n", payment.String())

	// To demonstrate timestamp conversion to a readable format
	timestamp := time.Unix(payment.Timestamp, 0)
	fmt.Println("Timestamp:", timestamp)
}
