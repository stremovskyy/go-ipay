package go_ipay

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestRequestIsMobile(t *testing.T) {
	applePayload := "apple"
	googlePayload := "google"

	tests := []struct {
		name string
		req  Request
		want bool
	}{
		{
			name: "no payment data",
			req:  Request{},
			want: false,
		},
		{
			name: "mobile flag true",
			req: Request{
				PaymentData: &PaymentData{IsMobile: true},
			},
			want: true,
		},
		{
			name: "payment method nil",
			req: Request{
				PaymentData: &PaymentData{},
			},
			want: false,
		},
		{
			name: "apple pay container",
			req: Request{
				PaymentData:   &PaymentData{},
				PaymentMethod: &PaymentMethod{AppleContainer: &applePayload},
			},
			want: true,
		},
		{
			name: "google pay token",
			req: Request{
				PaymentData:   &PaymentData{},
				PaymentMethod: &PaymentMethod{GoogleToken: &googlePayload},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.IsMobile(); got != tt.want {
				t.Errorf("IsMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestGetAppleContainerUsesApplePaymentToken(t *testing.T) {
	token := map[string]any{
		"paymentData": map[string]any{
			"version":   "EC_v1",
			"data":      "encrypted-data",
			"signature": "payment-signature",
			"header": map[string]any{
				"ephemeralPublicKey": "public-key",
				"publicKeyHash":      "public-key-hash",
				"transactionId":      "transaction-id",
			},
		},
		"paymentMethod": map[string]any{
			"displayName": "Visa 4655",
			"network":     "Visa",
			"type":        "debit",
		},
		"transactionIdentifier": "TRANSACTION-ID",
	}

	tests := []struct {
		name      string
		container any
	}{
		{
			name:      "direct apple token",
			container: token,
		},
		{
			name: "legacy event payment wrapper",
			container: map[string]any{
				"billingContact": map[string]any{"countryCode": "UA"},
				"token":          token,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := base64.StdEncoding.EncodeToString(mustJSON(t, tt.container))
			req := &Request{
				PaymentMethod: &PaymentMethod{AppleContainer: &container},
			}

			got, err := req.GetAppleContainer()
			if err != nil {
				t.Fatalf("GetAppleContainer() error: %v", err)
			}
			if got == nil || strings.TrimSpace(*got) == "" {
				t.Fatalf("GetAppleContainer() returned empty token")
			}

			decoded := decodeBase64JSON(t, *got)
			if _, ok := decoded["paymentData"].(map[string]any); !ok {
				t.Fatalf("decoded apple_data missing paymentData: %#v", decoded)
			}
			if _, ok := decoded["paymentMethod"].(map[string]any); !ok {
				t.Fatalf("decoded apple_data missing paymentMethod: %#v", decoded)
			}
			if decoded["transactionIdentifier"] != "TRANSACTION-ID" {
				t.Fatalf("transactionIdentifier = %#v, want TRANSACTION-ID", decoded["transactionIdentifier"])
			}
			if _, ok := decoded["billingContact"]; ok {
				t.Fatalf("decoded apple_data should contain only the Apple token, got wrapper: %#v", decoded)
			}
		})
	}
}

func TestRequestGetAppleContainerRejectsPayloadWithoutToken(t *testing.T) {
	container := base64.StdEncoding.EncodeToString(mustJSON(t, map[string]any{
		"billingContact": map[string]any{"countryCode": "UA"},
	}))
	req := &Request{
		PaymentMethod: &PaymentMethod{AppleContainer: &container},
	}

	got, err := req.GetAppleContainer()
	if err == nil {
		t.Fatalf("GetAppleContainer() error = nil, token = %v", got)
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()

	out, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}

	return out
}

func decodeBase64JSON(t *testing.T, value string) map[string]any {
	t.Helper()

	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("json unmarshal %q: %v", string(raw), err)
	}

	return out
}
