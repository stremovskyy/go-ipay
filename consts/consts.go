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

package consts

const (
	Version    = "1.0.0"
	ApiVersion = "1.28"

	baseUrl = "https://tokly.ipay.ua"

	ApiUrl       = baseUrl + "/api"
	ApplePayUrl  = "https://api-applepay.ipay.ua"
	GooglePayUrl = "https://api-googlepay.ipay.ua"
	ApiXMLUrl    = baseUrl + "/api302"

	RepaymentUrl = "https://api-repayment.ipay.ua"
)

const (
	VerificationLink           = "VerificationLink"
	Status                     = "Status"
	Payment                    = "Payment"
	Hold                       = "Hold"
	Capture                    = "Capture"
	Refund                     = "Refund"
	Credit                     = "Credit"
	ApplePaySuffix             = "ApplePay"
	GooglePaySuffix            = "GooglePay"
	A2CPaymentStatus           = "A2CPaymentStatus"
	CreateRepayment            = "CreateRepayment"
	CancelRepayment            = "CancelRepayment"
	GetRepaymentStatus         = "GetRepaymentStatus"
	GetRepaymentProcessingFile = "GetRepaymentProcessingFile"
)
