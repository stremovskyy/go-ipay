package go_ipay

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stremovskyy/go-ipay/internal/teststand"
	"github.com/stremovskyy/go-ipay/ipay"
)

func TestCreditUnknownOutcomeRequiresStatusCheckWithoutRetry(t *testing.T) {
	const extID = "withdrawal-123"
	cardToken := "card-token"
	attempts := 0
	roundTripper := teststand.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		return teststand.Response(
			http.StatusGatewayTimeout,
			"text/html",
			[]byte("<html>upstream timeout</html>"),
		), nil
	})
	client := NewClient(WithClient(&http.Client{Transport: roundTripper}))

	_, err := client.Credit(&Request{
		Merchant: &Merchant{
			MerchantID:  "1234",
			MerchantKey: "merchant-key",
		},
		PaymentData: &PaymentData{
			PaymentID: refString(extID),
			Amount:    100,
		},
		PaymentMethod: &PaymentMethod{
			Card: &Card{Token: &cardToken},
		},
	})
	if err == nil {
		t.Fatal("Credit() error = nil")
	}
	if attempts != 1 {
		t.Fatalf("HTTP attempts = %d, want 1", attempts)
	}
	if !errors.Is(err, ipay.ErrUnexpectedResponse) {
		t.Fatalf("Credit() error = %v, want ErrUnexpectedResponse", err)
	}
	if !ipay.RequiresA2CStatusCheck(err) {
		t.Fatalf("RequiresA2CStatusCheck() = false; error = %v", err)
	}
	var responseErr *ipay.UnexpectedResponseError
	if !errors.As(err, &responseErr) {
		t.Fatalf("errors.As(*UnexpectedResponseError) = false; error = %v", err)
	}
	if responseErr.Action != ipay.ActionCredit || responseErr.Operation != "Credit" {
		t.Fatalf("response context = action %q, operation %q", responseErr.Action, responseErr.Operation)
	}

	_, statusErr := client.A2CPaymentStatus(&Request{
		Merchant: &Merchant{
			MerchantID:  "1234",
			MerchantKey: "merchant-key",
		},
		PaymentData: &PaymentData{PaymentID: refString(extID)},
	})
	if statusErr == nil {
		t.Fatal("A2CPaymentStatus() error = nil")
	}
	if attempts != 2 {
		t.Fatalf("HTTP attempts after explicit status check = %d, want 2", attempts)
	}
	if ipay.RequiresA2CStatusCheck(statusErr) {
		t.Fatalf("status-check failure requires another status check: %v", statusErr)
	}
}
