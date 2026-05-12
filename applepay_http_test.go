package go_ipay

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stremovskyy/go-ipay/consts"
	"github.com/stremovskyy/go-ipay/currency"
	"github.com/stremovskyy/go-ipay/internal/teststand"
)

func TestApplePayPaymentCreateHTTPShape(t *testing.T) {
	const (
		login         = "merchant-login"
		key           = "merchant-sign-key"
		extID         = "AD68E7675FE111E79A65005056B960DF"
		subMerchantID = 112233
	)

	appleToken := map[string]any{
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
	appleContainer := base64.StdEncoding.EncodeToString(mustJSON(t, appleToken))

	rt := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
		}
		if req.URL.Scheme != "https" || req.URL.Host != "api-applepay.ipay.ua" {
			t.Fatalf("url = %q, want https://api-applepay.ipay.ua", req.URL.String())
		}
		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}

		raw, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		_ = req.Body.Close()

		var payload map[string]any
		if err := json.Unmarshal(raw, &payload); err != nil {
			t.Fatalf("unmarshal request JSON: %v\nraw=%s", err, string(raw))
		}

		requestObj := payload["request"].(map[string]any)
		if requestObj["version"] != 1.11 {
			t.Fatalf("request.version = %#v, want 1.11", requestObj["version"])
		}
		if requestObj["action"] != "PaymentCreate" {
			t.Fatalf("request.action = %#v, want PaymentCreate", requestObj["action"])
		}

		auth := requestObj["auth"].(map[string]any)
		if auth["login"] != login {
			t.Fatalf("auth.login = %#v, want %q", auth["login"], login)
		}
		timeString, ok := auth["time"].(string)
		if !ok {
			t.Fatalf("auth.time = %#v, want string", auth["time"])
		}
		if _, err := time.Parse("2006-01-02 15:04:05", timeString); err != nil {
			t.Fatalf("auth.time parse error: %v", err)
		}
		sign, ok := auth["sign"].(string)
		if !ok {
			t.Fatalf("auth.sign = %#v, want string", auth["sign"])
		}
		if len(sign) != 128 {
			t.Fatalf("len(auth.sign) = %d, want 128", len(sign))
		}
		if _, err := hex.DecodeString(sign); err != nil {
			t.Fatalf("auth.sign is not hex: %v", err)
		}

		body := requestObj["body"].(map[string]any)
		if _, exists := body["info"]; exists {
			t.Fatalf("body.info must not be sent for Apple Pay PaymentCreate: %#v", body["info"])
		}
		if body["invoice"] != float64(12345) {
			t.Fatalf("body.invoice = %#v, want 12345", body["invoice"])
		}
		if body["ext_id"] != extID {
			t.Fatalf("body.ext_id = %#v, want %q", body["ext_id"], extID)
		}
		if body["pmt_desc"] != "Apple Pay ride" {
			t.Fatalf("body.pmt_desc = %#v, want Apple Pay ride", body["pmt_desc"])
		}

		pmtInfo := body["pmt_info"].(map[string]any)
		if pmtInfo["order_id"] != extID {
			t.Fatalf("pmt_info.order_id = %#v, want %q", pmtInfo["order_id"], extID)
		}
		if pmtInfo["ext_id"] != extID {
			t.Fatalf("pmt_info.ext_id = %#v, want %q", pmtInfo["ext_id"], extID)
		}
		if pmtInfo["metadata"] != "auth:apple-pay-auth" {
			t.Fatalf("pmt_info.metadata = %#v, want auth:apple-pay-auth", pmtInfo["metadata"])
		}

		transactions := body["transactions"].([]any)
		if len(transactions) != 1 {
			t.Fatalf("len(transactions) = %d, want 1", len(transactions))
		}
		tx := transactions[0].(map[string]any)
		if tx["invoice"] != float64(12345) {
			t.Fatalf("transactions[0].invoice = %#v, want 12345", tx["invoice"])
		}
		if tx["smch_id"] != float64(subMerchantID) {
			t.Fatalf("transactions[0].smch_id = %#v, want %d", tx["smch_id"], subMerchantID)
		}
		if tx["desc"] != "Apple Pay ride" {
			t.Fatalf("transactions[0].desc = %#v, want Apple Pay ride", tx["desc"])
		}

		appleData, ok := body["apple_data"].(string)
		if !ok || appleData == "" {
			t.Fatalf("body.apple_data = %#v, want non-empty string", body["apple_data"])
		}
		decodedAppleData := decodeBase64JSON(t, appleData)
		if _, ok := decodedAppleData["paymentData"]; !ok {
			t.Fatalf("decoded apple_data missing paymentData: %#v", decodedAppleData)
		}

		responseBody := []byte(`{"response":{"pmt_id":"1234567","invoice":"12345","amount":"12345","pmt_status":"5","card_mask":"411111******1111"}}`)
		return teststand.Response(http.StatusOK, "application/json", responseBody), nil
	})

	client := NewClient(WithClient(&http.Client{Transport: rt}))
	resp, err := client.Payment(&Request{
		Merchant: &Merchant{
			Login:         login,
			SystemKey:     key,
			SubMerchantID: subMerchantID,
		},
		PaymentMethod: &PaymentMethod{
			AppleContainer: &appleContainer,
		},
		PaymentData: &PaymentData{
			PaymentID:   refString(extID),
			Amount:      12345,
			Currency:    currency.UAH,
			Description: "Apple Pay ride",
			IsMobile:    true,
			Metadata:    map[string]string{"auth": "apple-pay-auth"},
		},
	})
	if err != nil {
		t.Fatalf("Payment() error: %v", err)
	}
	if resp == nil || resp.GetPaymentStatus() != 5 {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestApplePayPaymentCreateDryRunUsesAppleEndpoint(t *testing.T) {
	var (
		gotEndpoint string
		gotPayload  any
	)

	container := base64.StdEncoding.EncodeToString(mustJSON(t, map[string]any{
		"paymentData":           map[string]any{"version": "EC_v1"},
		"paymentMethod":         map[string]any{"network": "Visa"},
		"transactionIdentifier": "TRANSACTION-ID",
	}))

	client := NewDefaultClient()
	_, err := client.Payment(&Request{
		Merchant:      &Merchant{Login: "login", SystemKey: "key"},
		PaymentMethod: &PaymentMethod{AppleContainer: &container},
		PaymentData: &PaymentData{
			PaymentID: refString("apple-dry-run"),
			Amount:    100,
			IsMobile:  true,
		},
	}, DryRun(func(endpoint string, payload any) {
		gotEndpoint = endpoint
		gotPayload = payload
	}))
	if err != nil {
		t.Fatalf("Payment(DryRun) error: %v", err)
	}
	if gotEndpoint != consts.ApplePayUrl {
		t.Fatalf("endpoint = %q, want %q", gotEndpoint, consts.ApplePayUrl)
	}
	if gotPayload == nil {
		t.Fatalf("dry-run payload is nil")
	}
}

func refString(value string) *string {
	return &value
}
