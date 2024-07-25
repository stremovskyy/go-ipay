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

package ipay

// BankErrorStatusCode holds the machine-readable error and the user-friendly message.
type BankErrorStatusCode struct {
	Code    string
	Reason  string
	Message string
	ExtCode int
}

// StatusCode is a type alias for string to represent the status codes.
type StatusCode string

// statusCodes holds the map of possible API response codes and messages.
var statusCodes = map[StatusCode]BankErrorStatusCode{
	"41-eminent_decline": {
		Code:    "41-eminent_decline",
		Reason:  "Operation declined by the issuing bank",
		Message: "Payment declined by the bank. Possible restrictions or limits on internet operations. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
		ExtCode: 41,
	},
	"42-insufficient_funds": {
		Code:    "42-insufficient_funds",
		Reason:  "Insufficient funds",
		Message: "Payment declined by the bank due to insufficient funds. Recommend specifying a lower amount or using another card.",
		ExtCode: 42,
	},
	"43-limits_emitent": {
		Code:    "43-limits_emitent",
		Reason:  "Exceeded the card's limit for transactions - possibly not open for online payments",
		Message: "Unsuccessful payment, your card has hit the limits for internet operations. It's recommended to contact the bank's hotline to have the operator increase the limits.",
		ExtCode: 43,
	},
	"44-limits_terminal": {
		Code:    "44-limits_terminal",
		Reason:  "Exceeded the merchant's limit or transactions prohibited to the merchant",
		Message: "Payment declined by the bank used for payment as our partner bank has set its restrictions.",
		ExtCode: 44,
	},
	"50-verification_error_CVV": {
		Code:    "50-verification_error_CVV",
		Reason:  "Incorrect CVV code",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
		ExtCode: 50,
	},
	"51-verification_error_3d_2d": {
		Code:    "51-verification_error_3d_2d",
		Reason:  "Incorrect 3DS confirmation code or session expired",
		Message: "Payment declined because you did not enter the confirmation code from your bank. Recommend trying the payment again, entering the new code sent to the phone linked to your card.",
		ExtCode: 51,
	},
	"52-connection_error": {
		Code:    "52-connection_error",
		Reason:  "Script error",
		Message: "Payment declined as there was no connection with the bank at the time of payment. Recommend trying the payment again in 2 minutes.",
		ExtCode: 52,
	},
	"55-unmatched_error": {
		Code:    "55-unmatched_error",
		Reason:  "Undefined error",
		Message: "Payment declined by the bank. It's recommended to contact the hotline to clarify the reason for the refusal.",
		ExtCode: 55,
	},
	"56-expired_card": {
		Code:    "56-expired_card",
		Reason:  "Card expired or incorrectly specified validity period",
		Message: "Payment declined by the bank. The card might be expired or the validity period incorrectly specified. Check the card's expiry date and try again.",
		ExtCode: 56,
	},
	"57-invalid_card": {
		Code:    "57-invalid_card",
		Reason:  "Incorrect card number entered, or card in an unacceptable state",
		Message: "Payment declined by the bank. Possible restrictions or limits on internet operations. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
		ExtCode: 57,
	},
	"58-card_limits_failed": {
		Code:    "58-card_limits_failed",
		Reason:  "Exceeded card limit",
		Message: "Payment declined by the bank. Your card has hit the limits for internet operations. It's recommended to contact the bank's hotline to clarify the reason for the refusal.",
		ExtCode: 58,
	},
	"59-invalid_amount": {
		Code:    "59-invalid_amount",
		Reason:  "Incorrect amount",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
		ExtCode: 59,
	},
	"60-3ds_fail": {
		Code:    "60-3ds_fail",
		Reason:  "Unable to perform 3DS transaction",
		Message: "Payment declined as there was no connection with the bank or the one-time password was incorrectly specified. Recommend trying the payment again in 2 minutes.",
		ExtCode: 60,
	},
	"61-call_issuer": {
		Code:    "61-call_issuer",
		Reason:  "Call the card issuer",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
		ExtCode: 61,
	},
	"62-card_lost_or_stolen": {
		Code:    "62-card_lost_or_stolen",
		Reason:  "Card lost or stolen",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
		ExtCode: 62,
	},
}

// GetStatusCode retrieves status code information by its code.
func GetStatusCode(code StatusCode) (BankErrorStatusCode, bool) {
	statusCode, found := statusCodes[code]
	return statusCode, found
}
