/*
 * MIT License
 *
 * Copyright (c) 2026 Anton Stremovskyy
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
	"errors"
	"fmt"
	"strings"
)

const maxDiagnosticBodyBytes = 4 << 10

var (
	// ErrTransport identifies a failure while executing or reading an HTTP request.
	ErrTransport = errors.New("ipay: transport error")
	// ErrUnexpectedResponse identifies an HTTP response that cannot be treated as an iPay JSON response.
	ErrUnexpectedResponse = errors.New("ipay: unexpected response")
	// ErrDecode identifies a response that looked like JSON but could not be decoded.
	ErrDecode = errors.New("ipay: decode error")
)

// TransportError describes a network or response-stream failure.
// Body is a bounded defensive copy when a partial response was available;
// Error deliberately does not print it. OutcomeUnknown is set only when an
// HTTP attempt may have reached iPay.
type TransportError struct {
	Op             string
	Action         Action
	Operation      string
	Method         string
	Endpoint       string
	RequestID      string
	StatusCode     int
	ContentType    string
	Body           []byte
	OutcomeUnknown bool
	Cause          error
}

func (e *TransportError) Error() string {
	if e == nil {
		return ErrTransport.Error()
	}

	parts := []string{ErrTransport.Error()}
	parts = appendResponseContext(parts, e.Op, e.Action, e.Operation, e.Method, e.Endpoint, e.RequestID)
	if e.StatusCode != 0 {
		parts = append(parts, fmt.Sprintf("status=%d", e.StatusCode))
	}
	if e.ContentType != "" {
		parts = append(parts, fmt.Sprintf("content-type=%q", e.ContentType))
	}
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}

	return strings.Join(parts, ": ")
}

func (e *TransportError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func (e *TransportError) Is(target error) bool {
	return target == ErrTransport
}

func (e *TransportError) requiresA2CStatusCheck() bool {
	return e != nil && e.Action == ActionCredit && e.OutcomeUnknown
}

// UnexpectedResponseError describes a non-2xx, empty, non-JSON, nil, or oversized response.
// Body is a bounded defensive copy intended for diagnostics; Error deliberately does not print it.
type UnexpectedResponseError struct {
	Op             string
	Action         Action
	Operation      string
	Method         string
	Endpoint       string
	RequestID      string
	StatusCode     int
	ContentType    string
	Message        string
	Body           []byte
	OutcomeUnknown bool
	Cause          error
}

func (e *UnexpectedResponseError) Error() string {
	if e == nil {
		return ErrUnexpectedResponse.Error()
	}

	parts := []string{ErrUnexpectedResponse.Error()}
	parts = appendResponseContext(parts, e.Op, e.Action, e.Operation, e.Method, e.Endpoint, e.RequestID)
	if e.StatusCode != 0 {
		parts = append(parts, fmt.Sprintf("status=%d", e.StatusCode))
	}
	if e.ContentType != "" {
		parts = append(parts, fmt.Sprintf("content-type=%q", e.ContentType))
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}

	return strings.Join(parts, ": ")
}

func (e *UnexpectedResponseError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func (e *UnexpectedResponseError) Is(target error) bool {
	return target == ErrUnexpectedResponse
}

func (e *UnexpectedResponseError) requiresA2CStatusCheck() bool {
	return e != nil && e.Action == ActionCredit && e.OutcomeUnknown
}

// DecodeError describes a malformed JSON response from iPay.
// Body is a bounded defensive copy intended for diagnostics; Error deliberately does not print it.
type DecodeError struct {
	Op             string
	Action         Action
	Operation      string
	Method         string
	Endpoint       string
	RequestID      string
	StatusCode     int
	ContentType    string
	Message        string
	Body           []byte
	OutcomeUnknown bool
	Cause          error
}

func (e *DecodeError) Error() string {
	if e == nil {
		return ErrDecode.Error()
	}

	parts := []string{ErrDecode.Error()}
	parts = appendResponseContext(parts, e.Op, e.Action, e.Operation, e.Method, e.Endpoint, e.RequestID)
	if e.StatusCode != 0 {
		parts = append(parts, fmt.Sprintf("status=%d", e.StatusCode))
	}
	if e.ContentType != "" {
		parts = append(parts, fmt.Sprintf("content-type=%q", e.ContentType))
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}

	return strings.Join(parts, ": ")
}

func (e *DecodeError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func (e *DecodeError) Is(target error) bool {
	return target == ErrDecode
}

func (e *DecodeError) requiresA2CStatusCheck() bool {
	return e != nil && e.Action == ActionCredit && e.OutcomeUnknown
}

// RequiresA2CStatusCheck reports whether an A2CPay result is unknown and must
// be reconciled with A2CPaymentStatus by ext_id before any retry.
func RequiresA2CStatusCheck(err error) bool {
	if err == nil {
		return false
	}

	var candidate interface {
		requiresA2CStatusCheck() bool
	}
	if !errors.As(err, &candidate) {
		return false
	}

	return candidate.requiresA2CStatusCheck()
}

func diagnosticBody(data []byte) []byte {
	if len(data) > maxDiagnosticBodyBytes {
		data = data[:maxDiagnosticBodyBytes]
	}

	return append([]byte(nil), data...)
}

func appendResponseContext(
	parts []string,
	op string,
	action Action,
	operation string,
	method string,
	endpoint string,
	requestID string,
) []string {
	if op != "" {
		parts = append(parts, op)
	}
	if action != "" {
		parts = append(parts, "action="+string(action))
	}
	if operation != "" {
		parts = append(parts, "operation="+operation)
	}
	if method != "" && endpoint != "" {
		parts = append(parts, method+" "+endpoint)
	} else if endpoint != "" {
		parts = append(parts, endpoint)
	}
	if requestID != "" {
		parts = append(parts, "request_id="+requestID)
	}

	return parts
}
