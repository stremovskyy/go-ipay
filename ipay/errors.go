package ipay

// StatusCode holds the machine-readable error and the user-friendly message
type StatusCode struct {
	Error   string
	Message string
}

// StatusCodes maps numeric status codes to their descriptions
var StatusCodes = map[int]StatusCode{
	41: {"eminent_decline", "The payment was declined by the bank possibly due to restrictions or limits on online operations. Please contact your bank's support to inquire about the refusal."},
	42: {"insufficient_funds", "The payment was declined by the bank due to insufficient funds. Please consider using a different card or checking for the correct amount."},
	43: {"limits_emitent", "Payment unsuccessful due to limits on internet operations set on your card. Please contact your bank to modify these limits."},
	44: {"limits_terminal", "The payment was declined due to restrictions set by our banking partner. Please try a different payment method."},
	50: {"verification_error_CVV", "The payment was declined by the bank. Please verify the CVV code and try again or contact your bank for assistance."},
	51: {
		"verification_error_3d_2d",
		"The payment was declined because the 3DS verification failed or the session expired. Please attempt the payment again and ensure you enter the confirmation code sent by your bank.",
	},
	52: {"connection_error", "The payment was declined due to a connection error at the time of the transaction. Please try again in a few minutes."},
	55: {"unmatched_error", "The payment was declined by the bank for an unspecified reason. Please contact your bank's support for more information."},
	56: {"expired_card", "The payment was declined because the card is expired or the validity period was incorrectly specified. Please check the card's expiry date and try again."},
	57: {"invalid_card", "The payment was declined by the bank due to an invalid card number or because the card is in an unacceptable state. Please verify the card details and try again."},
	58: {"card_limits_failed", "The payment was declined by the bank because the transaction exceeded the card's limits. Please contact your bank to inquire about modifying these limits."},
	59: {"invalid_amount", "The payment was declined due to an incorrect amount being specified. Please verify the amount and try again."},
	60: {"3ds_fail", "The payment was declined because the 3DS verification could not be completed. Please try the payment again later."},
	61: {"call_issuer", "The payment was declined. Please contact your card issuer for further assistance."},
	62: {"card_lost_or_stolen", "The payment was declined because the card was reported lost or stolen. Please contact your bank for further assistance."},
}

// Function to retrieve status code information by its numeric code
func GetStatusCode(code int) (StatusCode, bool) {
	statusCode, exists := StatusCodes[code]
	return statusCode, exists
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
