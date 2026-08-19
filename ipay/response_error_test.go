package ipay

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestUnmarshalJSONResponseTypedDecodeErrors(t *testing.T) {
	tests := []struct {
		name string
		raw  []byte
	}{
		{name: "nil", raw: nil},
		{name: "empty", raw: []byte{}},
		{name: "whitespace", raw: []byte(" \n\t")},
		{name: "malformed", raw: []byte(`{"response":`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalJSONResponse(tt.raw)
			if !errors.Is(err, ErrDecode) {
				t.Fatalf("UnmarshalJSONResponse() error = %v, want ErrDecode", err)
			}
			var decodeErr *DecodeError
			if !errors.As(err, &decodeErr) {
				t.Fatalf("errors.As(*DecodeError) = false; error = %v", err)
			}
			if strings.Contains(err.Error(), string(tt.raw)) && len(tt.raw) > 0 {
				t.Fatalf("error contains raw response body: %v", err)
			}
		})
	}
}

func TestResponseErrorSentinelsAndStatusCheck(t *testing.T) {
	raw := []byte("<html>sensitive</html>")
	err := &UnexpectedResponseError{
		Action:         ActionCredit,
		Operation:      "credit",
		RequestID:      "request-123",
		StatusCode:     504,
		ContentType:    "text/html",
		Body:           diagnosticBody(raw),
		OutcomeUnknown: true,
	}
	raw[0] = 'X'

	wrapped := fmt.Errorf("credit call: %w", err)
	if !errors.Is(wrapped, ErrUnexpectedResponse) {
		t.Fatalf("errors.Is(ErrUnexpectedResponse) = false; error = %v", wrapped)
	}
	var typed *UnexpectedResponseError
	if !errors.As(wrapped, &typed) {
		t.Fatalf("errors.As(*UnexpectedResponseError) = false; error = %v", wrapped)
	}
	if !RequiresA2CStatusCheck(wrapped) {
		t.Fatalf("RequiresA2CStatusCheck() = false; error = %v", wrapped)
	}
	if string(typed.Body) != "<html>sensitive</html>" {
		t.Fatalf("Body = %q, want defensive copy", typed.Body)
	}
	if strings.Contains(wrapped.Error(), "sensitive") {
		t.Fatalf("Error() contains diagnostic body: %v", wrapped)
	}

	statusErr := &TransportError{
		Action:         ActionA2CPaymentStatus,
		OutcomeUnknown: true,
	}
	if RequiresA2CStatusCheck(statusErr) {
		t.Fatal("status-check request must not require another status check")
	}

	providerErr := createIpayError(601, "declined", "business refusal")
	if RequiresA2CStatusCheck(providerErr) {
		t.Fatal("business error must not require an A2C status check")
	}
}

func TestDiagnosticBodyLimit(t *testing.T) {
	raw := []byte(strings.Repeat("x", maxDiagnosticBodyBytes+1))
	snapshot := diagnosticBody(raw)
	if len(snapshot) != maxDiagnosticBodyBytes {
		t.Fatalf("len(snapshot) = %d, want %d", len(snapshot), maxDiagnosticBodyBytes)
	}
	raw[0] = 'y'
	if snapshot[0] != 'x' {
		t.Fatal("snapshot aliases source body")
	}
}
