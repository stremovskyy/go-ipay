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
}

// StatusCode is a type alias for string to represent the status codes.
type StatusCode string

// statusCodes holds the map of possible API response codes and messages.
var statusCodes = map[StatusCode]BankErrorStatusCode{
	"41-eminent_decline": {
		Code:    "41-eminent_decline",
		Reason:  "Operation declined by the issuing bank",
		Message: "Payment declined by the bank. Possible restrictions or limits on internet operations. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
	},
	"42-insufficient_funds": {
		Code:    "42-insufficient_funds",
		Reason:  "Insufficient funds",
		Message: "Payment declined by the bank due to insufficient funds. Recommend specifying a lower amount or using another card.",
	},
	"43-limits_emitent": {
		Code:    "43-limits_emitent",
		Reason:  "Exceeded the card's limit for transactions - possibly not open for online payments",
		Message: "Unsuccessful payment, your card has hit the limits for internet operations. It's recommended to contact the bank's hotline to have the operator increase the limits.",
	},
	"44-limits_terminal": {
		Code:    "44-limits_terminal",
		Reason:  "Exceeded the merchant's limit or transactions prohibited to the merchant",
		Message: "Payment declined by the bank used for payment as our partner bank has set its restrictions.",
	},
	"50-verification_error_CVV": {
		Code:    "50-verification_error_CVV",
		Reason:  "Incorrect CVV code",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
	},
	"51-verification_error_3d_2d": {
		Code:    "51-verification_error_3d_2d",
		Reason:  "Incorrect 3DS confirmation code or session expired",
		Message: "Payment declined because you did not enter the confirmation code from your bank. Recommend trying the payment again, entering the new code sent to the phone linked to your card.",
	},
	"52-connection_error": {
		Code:    "52-connection_error",
		Reason:  "Script error",
		Message: "Payment declined as there was no connection with the bank at the time of payment. Recommend trying the payment again in 2 minutes.",
	},
	"55-unmatched_error": {
		Code:    "55-unmatched_error",
		Reason:  "Undefined error",
		Message: "Payment declined by the bank. It's recommended to contact the hotline to clarify the reason for the refusal.",
	},
	"56-expired_card": {
		Code:    "56-expired_card",
		Reason:  "Card expired or incorrectly specified validity period",
		Message: "Payment declined by the bank. The card might be expired or the validity period incorrectly specified. Check the card's expiry date and try again.",
	},
	"57-invalid_card": {
		Code:    "57-invalid_card",
		Reason:  "Incorrect card number entered, or card in an unacceptable state",
		Message: "Payment declined by the bank. Possible restrictions or limits on internet operations. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
	},
	"58-card_limits_failed": {
		Code:    "58-card_limits_failed",
		Reason:  "Exceeded card limit",
		Message: "Payment declined by the bank. Your card has hit the limits for internet operations. It's recommended to contact the bank's hotline to clarify the reason for the refusal.",
	},
	"59-invalid_amount": {
		Code:    "59-invalid_amount",
		Reason:  "Incorrect amount",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
	},
	"60-3ds_fail": {
		Code:    "60-3ds_fail",
		Reason:  "Unable to perform 3DS transaction",
		Message: "Payment declined as there was no connection with the bank or the one-time password was incorrectly specified. Recommend trying the payment again in 2 minutes.",
	},
	"61-call_issuer": {
		Code:    "61-call_issuer",
		Reason:  "Call the card issuer",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
	},
	"62-card_lost_or_stolen": {
		Code:    "62-card_lost_or_stolen",
		Reason:  "Card lost or stolen",
		Message: "Payment declined by the bank. Recommended to contact your bank's hotline to clarify the reason for the refusal.",
	},
}

// GetStatusCode retrieves status code information by its code.
func GetStatusCode(code StatusCode) (BankErrorStatusCode, bool) {
	statusCode, found := statusCodes[code]
	return statusCode, found
}

// A2CPayStatusCode holds the machine-readable error and the user-friendly message
type A2CPayStatusCode struct {
	Error   string
	Message string
}

// A2CPayStatusCodes maps numeric status codes for A2CPay operations to their descriptions
var A2CPayStatusCodes = map[int]A2CPayStatusCode{
	0:   {"successful_accreditation", "Successful accreditation."},
	100: {"issuer_bank_decline", "Declined, please contact the issuing bank."},
	101: {"card_expired", "The card's expiration date has passed."},
	104: {"card_restriction", "Restrictions on the card (local or prohibited operations)."},
	106: {"card_blocked", "The card is blocked."},
	110: {"transaction_amount_exceeded", "The transaction amount exceeds the allowed limit."},
	111: {"invalid_card_number", "Invalid card number."},
	116: {"transaction_limit_exceeded", "The transaction amount exceeds the allowed limit."},
	118: {"card_not_active", "The card is not active, please contact the issuing bank."},
	120: {"card_restriction_contact_issuer", "Restrictions on the card, please contact the issuing bank."},
	121: {"card_limits", "Limits on the card (restrictions on internet transactions)."},
	123: {"issuer_or_payment_system_decline", "Decline by the issuing bank/payment system due to the number of operations."},
	124: {"card_restriction_legal", "Restrictions on the card (prohibited by law)."},
	200: {"invalid_card_number", "Invalid card number."},
	202: {"invalid_card_number", "Invalid card number."},
	208: {"lost_card", "The card is reported lost."},
	209: {"stolen_card", "The card is reported stolen."},
	907: {"issuer_bank_offline", "The issuing bank is currently offline."},
	908: {"bank_unavailable", "The bank is unavailable."},
	909: {"system_malfunction", "System malfunction."},
	600: {"invalid_digital_signature", "The digital signature is not valid."},
	601: {"public_key_not_found", "The public key was not found."},
	602: {"invalid_receiver_card_number", "Invalid card number for the receiver (failed Luhn check)."},
	603: {"service_only_for_domestic_banks", "This service is available only for cards issued by domestic banks."},
	604: {"technical_reasons", "The operation could not be performed due to technical reasons."},
	605: {"receiver_transfer_limit_exceeded", "The transfer limit for the receiver has been exceeded."},
	606: {"declined_by_receiver_issuer", "The operation was declined by the receiver's card issuing bank."},
	607: {"daily_replenishment_limit_reached", "The daily replenishment limit has been reached."},
	608: {"monthly_replenishment_limit_reached", "The monthly replenishment limit has been reached."},
	609: {"daily_payment_limit_per_card_reached", "The limit for the number of payments per card per day has been reached."},
	610: {"monthly_payment_limit_per_card_reached", "The limit for the number of payments per card per month has been reached."},
	611: {"payment_attempt_exceeds_limit", "The payment attempt exceeds the limit."},
	612: {"payment_already_made_with_different_amount", "A payment for the request has already been made with a different amount."},
	613: {"different_payment_date_provided", "A different payment date was provided upon confirmation."},
	614: {"different_partner_system_id_provided", "A different partner system ID was provided upon confirmation."},
	615: {"different_payment_amount_provided", "A different payment amount was provided upon confirmation."},
	616: {"different_card_hash_provided", "A different card hash was provided upon confirmation."},
	617: {"card_in_gray_list", "The card is in the gray list."},
	618: {"recipient_transfer_count_limit_exceeded", "The recipient's transfer count limit has been exceeded."},
	619: {"recipient_card_blocked_due_to_debt", "The recipient's card is blocked due to debt."},
	620: {"recipient_issuer_bank_unavailable", "The recipient's card issuing bank is unavailable."},
	621: {"recipient_issuer_unable_to_process", "The recipient's card issuing bank is unable to process the operation."},
	622: {"verify_recipient_card_details", "Need to verify the recipient card details with the issuing bank."},
	623: {"recipient_card_blocked_by_issuer", "The recipient's card is blocked by the issuing bank."},
}

// Function to retrieve A2CPay status code information by its numeric code
func GetA2CPayStatusCode(code int) (A2CPayStatusCode, bool) {
	statusCode, exists := A2CPayStatusCodes[code]
	return statusCode, exists
}
