/*
 * MIT License
 *
 * Copyright (c) 2025 Anton Stremovskyy
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

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stremovskyy/go-ipay/internal/ipay"
	"github.com/stremovskyy/go-ipay/internal/utils"
)

const (
	ErrorTypeValidation = "validation"
	ErrorTypeBank       = "bank"
	ErrorTypeSystem     = "system"
	ErrorTypeTransport  = "transport"
)

// IpayError represents a structured error for the iPay system.
type IpayError struct {
	Code        int         `json:"code"`
	Type        string      `json:"type"`
	Message     string      `json:"message"`
	Details     string      `json:"details,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
	Context     interface{} `json:"context,omitempty"`
	UserMessage string      `json:"user_message,omitempty"`
}

// Error satisfies the error interface.
func (e *IpayError) Error() string {
	return fmt.Sprintf("IpayError [Code: %d]: %s - %s", e.Code, e.Message, e.Details)
}

// IsTransient returns true if the error is transient and can be retried.
func (e *IpayError) IsTransient() bool {
	switch e.Code {
	case 907, 908, 909, 52: // TODO: more transient error codes as needed
		return true
	default:
		return false
	}
}

func inferErrorType(code int) string {
	switch {
	case code >= 600 && code <= 699:
		return ErrorTypeValidation
	case code >= 100 && code <= 299:
		return ErrorTypeBank
	case code >= 900:
		return ErrorTypeSystem
	default:
		return ErrorTypeSystem
	}
}

// IsValidationError checks if the error is a validation error
func (e *IpayError) IsValidationError() bool {
	return e.Type == ErrorTypeValidation
}

// IsBankError checks if the error is a bank-related error
func (e *IpayError) IsBankError() bool {
	return e.Type == ErrorTypeBank
}

// IsSystemError checks if the error is a system error
func (e *IpayError) IsSystemError() bool {
	return e.Type == ErrorTypeSystem
}

// LogJSON logs the error in JSON format for structured logging systems.
func (e *IpayError) LogJSON() {
	jsonError, err := json.Marshal(e)
	if err != nil {
		log.Printf("❌ Failed to marshal IpayError: %v\n", err)
		return
	}

	log.Println(string(jsonError))
}

// LogPlain logs the error in plain text format.
func (e *IpayError) LogPlain() {
	log.Printf("❌ Error [Code: %d]: %s\nDetails: %s\n", e.Code, e.Message, e.Details)
}

// createIpayError creates a new instance of IpayError with all details.
func createIpayError(code int, message, details string) *IpayError {
	return &IpayError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
		Type:      inferErrorType(code),
	}
}

func (r Response) GetError() error {
	// Check for general error messages
	if r.Error != nil {
		errorCode := 900
		if r.ErrorCode != nil {
			errorCode, _ = strconv.Atoi(*r.ErrorCode)
		}
		return createIpayError(
			errorCode,
			"iPay General Error",
			fmt.Sprintf("Error: %s, Code: %s", *r.Error, utils.SafeString(r.ErrorCode)),
		)
	}

	// Bank error note handling
	if r.BnkErrorNote != nil {
		if statusCode, found := ipay.GetStatusCode(*r.BnkErrorNote); found {
			return createIpayError(
				statusCode.ExtCode,
				"Bank Error",
				fmt.Sprintf("Reason: %s, Message: %s", statusCode.Reason, statusCode.Message),
			)
		}
		return createIpayError(
			900,
			"Bank Error Note",
			fmt.Sprintf("Note: %s", *r.BnkErrorNote),
		)
	}

	// Authorization code errors
	if r.ResAuthCode != 0 {
		message := getErrorMessageA2CPay(r.ResAuthCode)
		return createIpayError(
			r.ResAuthCode,
			"Authorization Code Error",
			message,
		)
	}

	switch r.GetPaymentStatus() {
	case PaymentStatusSecurityRefusal:
		return createIpayError(900, "Payment Status: Security Refusal", "")
	case PaymentStatusFailed:
		details := "Payment failed for unknown reasons"
		if r.Pmt != nil && r.Pmt.BnkErrorGroup != nil && r.Pmt.BnkErrorNote != nil {
			details = fmt.Sprintf("Bank Error: %s, Group: %v", r.Pmt.BnkErrorNote, r.Pmt.BnkErrorGroup)
		} else if r.BankResponse != nil && r.BankResponse.ErrorGroup != 0 {
			details = fmt.Sprintf("Bank Response Error Group: %d", r.BankResponse.ErrorGroup)
		}
		return createIpayError(900, "Payment Failed", details)
	}

	return nil
}
