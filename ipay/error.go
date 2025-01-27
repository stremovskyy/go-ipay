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

import "fmt"

type Error struct {
	Code    int
	Message string
	Details string
}

// Implement the error interface for IpayError
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("IpayError: code %d, message: %s, details: %s", e.Code, e.Message, e.Details)
	}

	return fmt.Sprintf("IpayError: code %d, message: %s", e.Code, e.Message)
}

func getErrorMessageA2CPay(code int) string {
	errorMessages := map[int]string{
		0:   "Successful transaction",
		100: "Declined, contact the issuing bank",
		101: "Card expired",
		104: "Card restriction (local or forbidden transactions)",
		106: "Card blocked",
		110: "Transaction amount exceeds allowed limit",
		111: "Incorrect card number",
		116: "Transaction amount exceeds allowed limit",
		118: "Card inactive, contact the issuing bank",
		120: "Card restriction, contact the issuing bank",
		121: "Card limits (internet transactions restrictions)",
		123: "Issuer/bank system decline due to transaction volume",
		124: "Card restriction (legally prohibited)",
		200: "Incorrect card number",
		202: "Incorrect card number",
		208: "Lost card",
		209: "Stolen card",
		907: "Issuer bank not operational",
		908: "Bank unavailable",
		909: "Technical system failure",
		600: "Digital signature invalid",
		601: "Public key not found",
		602: "Incorrect card-acceptor number (failed Luhn check)",
		603: "Only cards issued by Ukrainian banks are supported",
		604: "Transaction cannot be completed due to technical reasons",
		605: "Exceeded recipient's transfer limit",
		606: "Transaction declined by recipient’s issuing bank",
		607: "Daily top-up limit reached",
		608: "Monthly top-up limit reached",
		609: "Daily payment limit exceeded for one card",
		610: "Monthly payment limit exceeded for one card",
		611: "Attempt to process payment exceeding the limit",
		612: "Payment already processed with a different amount",
		613: "Different transaction date provided during confirmation",
		614: "Different partner system ID provided during confirmation",
		615: "Different payment amount provided during confirmation",
		616: "Different card hash number provided during confirmation",
		617: "Card listed in blacklists",
		618: "Exceeded recipient’s transfer quantity limit",
		619: "Recipient's card blocked due to debt",
		620: "Issuer bank of recipient's card unavailable",
		621: "Issuer bank of recipient's card cannot process the transaction",
		622: "Need to clarify recipient card details with issuing bank",
		623: "Recipient's card blocked by issuing bank",
		900: "Banker system error",
		901: "Banker validation error",
	}

	if message, found := errorMessages[code]; found {
		return message
	}

	return "Unknown error"
}

// BankErrorInfo contains details about a payment error
type BankErrorInfo struct {
	Code        string // Original error code
	Description string // Error description in English
	UserMessage string // User-friendly message in English
}

func (e *BankErrorInfo) Error() string {
	return fmt.Sprintf("BankErrorInfo: code %s, description: %s, user message: %s", e.Code, e.Description, e.UserMessage)
}

// GetBankErrorInfo extracts and interprets error information from either Response or Payment
func GetBankErrorInfo(data interface{}) *BankErrorInfo {
	var errorCode string

	// Extract error code based on input type
	switch v := data.(type) {
	case *Response:
		if v.ErrorCode != nil {
			errorCode = *v.ErrorCode
		} else if v.BnkErrorNote != nil {
			errorCode = string(*v.BnkErrorNote)
		}
	case *Payment:
		if v.BnkErrorNote != nil {
			errorCode = v.BnkErrorNote.(string)
		}
	default:
		return nil
	}

	if errorCode == "" {
		return nil
	}

	// Map error codes to user-friendly messages
	errorMap := map[string]BankErrorInfo{
		"41-eminent_decline": {
			Code:        "41-eminent_decline",
			Description: "Operation declined by issuing bank",
			UserMessage: "Payment declined by your bank. There might be restrictions or limits on internet transactions. Please contact your bank's support for details.",
		},
		"42-insufficient_funds": {
			Code:        "42-insufficient_funds",
			Description: "Insufficient funds",
			UserMessage: "Payment declined due to insufficient funds. Please try a smaller amount or use a different card.",
		},
		"43-limits_emitent": {
			Code:        "43-limits_emitent",
			Description: "Card transaction limits exceeded - card might not be enabled for internet payments",
			UserMessage: "Transaction failed due to card limits. Please contact your bank to enable/increase internet transaction limits.",
		},
		"44-limits_terminal": {
			Code:        "44-limits_terminal",
			Description: "Merchant limits exceeded or transactions forbidden for merchant",
			UserMessage: "Transaction declined due to merchant bank restrictions. Please try again later or use a different payment method.",
		},
		"50-verification_error_CVV": {
			Code:        "50-verification_error_CVV",
			Description: "Invalid CVV code",
			UserMessage: "Payment declined. Please verify your card details and try again.",
		},
		"51-verification_error_3d_2d": {
			Code:        "51-verification_error_3d_2d",
			Description: "Invalid 3DS confirmation code or session expired",
			UserMessage: "Payment declined. Please try again and enter the new confirmation code sent to your registered phone number.",
		},
		"52-connection_error": {
			Code:        "52-connection_error",
			Description: "Script error",
			UserMessage: "Connection error occurred. Please wait 2 minutes and try again.",
		},
		"55-unmatched_error": {
			Code:        "55-unmatched_error",
			Description: "Undefined error",
			UserMessage: "Payment declined. Please contact your bank's support for details.",
		},
		"56-expired_card": {
			Code:        "56-expired_card",
			Description: "Card expired or invalid expiration date",
			UserMessage: "Payment declined. Please check your card's expiration date and try again.",
		},
		"57-invalid_card": {
			Code:        "57-invalid_card",
			Description: "Invalid card number or card in invalid state",
			UserMessage: "Payment declined. Please verify your card details or contact your bank for assistance.",
		},
		"58-card_limits_failed": {
			Code:        "58-card_limits_failed",
			Description: "Card limits exceeded",
			UserMessage: "Payment declined due to card limits. Please contact your bank to review your transaction limits.",
		},
		"59-invalid_amount": {
			Code:        "59-invalid_amount",
			Description: "Invalid amount",
			UserMessage: "Payment declined due to invalid amount. Please contact your bank for details.",
		},
		"60-3ds_fail": {
			Code:        "60-3ds_fail",
			Description: "Unable to perform 3DS transaction",
			UserMessage: "3D Secure verification failed. Please try again in 2 minutes.",
		},
		"61-call_issuer": {
			Code:        "61-call_issuer",
			Description: "Call card issuer",
			UserMessage: "Please contact your card issuer for authorization.",
		},
		"62-card_lost_or_stolen": {
			Code:        "62-card_lost_or_stolen",
			Description: "Card lost or stolen",
			UserMessage: "Transaction declined. Please contact your bank immediately.",
		},
		"66-required_3ds": {
			Code:        "66-required_3ds",
			Description: "3DS verification required",
			UserMessage: "3D Secure verification required. Please complete the verification process.",
		},
		"67-card_country_not_allowed": {
			Code:        "67-card_country_not_allowed",
			Description: "Foreign bank card not allowed for this operation",
			UserMessage: "This card is not accepted. Please use a different payment card.",
		},
	}

	if info, exists := errorMap[errorCode]; exists {
		return &info
	}

	// Return generic error for unknown error codes
	return &BankErrorInfo{
		Code:        errorCode,
		Description: "Unknown error",
		UserMessage: "An error occurred processing your payment. Please try again or contact support.",
	}
}
