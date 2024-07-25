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

// Define the IpayError type
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

// Function to create an IpayError
func createIpayError(code int, message, details string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
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
