package go_ipay

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stremovskyy/go-ipay/consts"
	"github.com/stremovskyy/go-ipay/internal/teststand"
	"github.com/stremovskyy/go-ipay/ipay"
)

func TestA2CBalanceHTTPShapeAndResponse(t *testing.T) {
	attempts := 0
	roundTripper := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		attempts++

		if req.Method != http.MethodPost {
			t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
		}
		if req.URL.String() != consts.ApiUrl {
			t.Fatalf("url = %q, want %q", req.URL.String(), consts.ApiUrl)
		}
		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}
		if got := req.Header.Get("Api-Version"); got != consts.ApiVersion {
			t.Fatalf("Api-Version = %q, want %q", got, consts.ApiVersion)
		}
		if got := req.Header.Get("User-Agent"); got != "GO IPAY/1.16.0" {
			t.Fatalf("User-Agent = %q, want GO IPAY/1.16.0", got)
		}

		raw, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}

		var envelope struct {
			Request map[string]json.RawMessage `json:"request"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			t.Fatalf("unmarshal request JSON: %v\nraw=%s", err, string(raw))
		}
		if len(envelope.Request) != 2 {
			t.Fatalf("request fields = %v, want only auth and action", envelope.Request)
		}
		if _, exists := envelope.Request["body"]; exists {
			t.Fatalf("request.body must be omitted: %s", string(raw))
		}

		var action ipay.Action
		if err := json.Unmarshal(envelope.Request["action"], &action); err != nil {
			t.Fatalf("unmarshal request.action: %v", err)
		}
		if action != ipay.ActionA2CBalance {
			t.Fatalf("request.action = %q, want %q", action, ipay.ActionA2CBalance)
		}

		var auth struct {
			MchID int64  `json:"mch_id"`
			Salt  string `json:"salt"`
			Sign  string `json:"sign"`
		}
		if err := json.Unmarshal(envelope.Request["auth"], &auth); err != nil {
			t.Fatalf("unmarshal request.auth: %v", err)
		}
		if auth.MchID != 2995 {
			t.Fatalf("auth.mch_id = %d, want 2995", auth.MchID)
		}
		if auth.Salt == "" || auth.Sign == "" {
			t.Fatalf("auth signature is incomplete: salt=%q sign=%q", auth.Salt, auth.Sign)
		}

		return teststand.Response(
			http.StatusOK,
			"application/json",
			[]byte(`{"response":{"current_balance":-30,"overdraft":1000,"credit":0,"salt":"6af2952a2f1fe5e2c0385f99e84a9a05274328ef","sign":"0adb98c05edf8e86db11046ca25af4a433a03a3b776bb28f318f57852620319f09c71df06809904cda0fcb288bdfec5887fc504a691a82684a1db817afdc09c7"}}`),
		), nil
	})
	client := NewClient(WithClient(&http.Client{Transport: roundTripper}))

	response, err := client.A2CBalance(balanceRequest())
	if err != nil {
		t.Fatalf("A2CBalance() error: %v", err)
	}
	if attempts != 1 {
		t.Fatalf("HTTP attempts = %d, want 1", attempts)
	}
	if response == nil {
		t.Fatal("A2CBalance() response = nil")
	}
	if response.CurrentBalance != -30 || response.Overdraft != 1000 || response.Credit != 0 {
		t.Fatalf(
			"balance = current:%d overdraft:%d credit:%d, want -30/1000/0",
			response.CurrentBalance,
			response.Overdraft,
			response.Credit,
		)
	}
	if response.Salt != "6af2952a2f1fe5e2c0385f99e84a9a05274328ef" {
		t.Fatalf("Salt = %q", response.Salt)
	}
	if response.Sign != "0adb98c05edf8e86db11046ca25af4a433a03a3b776bb28f318f57852620319f09c71df06809904cda0fcb288bdfec5887fc504a691a82684a1db817afdc09c7" {
		t.Fatalf("Sign = %q", response.Sign)
	}
}

func TestA2CBalanceDryRun(t *testing.T) {
	attempts := 0
	roundTripper := teststand.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		return nil, errors.New("unexpected HTTP request")
	})
	client := NewClient(WithClient(&http.Client{Transport: roundTripper}))

	var endpoint string
	var payload *ipay.RequestWrapper
	response, err := client.A2CBalance(balanceRequest(), DryRun(func(gotEndpoint string, gotPayload any) {
		endpoint = gotEndpoint
		var ok bool
		payload, ok = gotPayload.(*ipay.RequestWrapper)
		if !ok {
			t.Fatalf("dry-run payload type = %T, want *ipay.RequestWrapper", gotPayload)
		}
	}))
	if err != nil {
		t.Fatalf("A2CBalance() dry-run error: %v", err)
	}
	if response != nil {
		t.Fatalf("A2CBalance() dry-run response = %#v, want nil", response)
	}
	if attempts != 0 {
		t.Fatalf("HTTP attempts = %d, want 0", attempts)
	}
	if endpoint != consts.ApiUrl {
		t.Fatalf("dry-run endpoint = %q, want %q", endpoint, consts.ApiUrl)
	}
	if payload == nil || payload.Request.Action != ipay.ActionA2CBalance {
		t.Fatalf("dry-run payload = %#v, want A2CBalance request", payload)
	}
	if payload.Operation != consts.A2CBalance {
		t.Fatalf("dry-run operation = %q, want %q", payload.Operation, consts.A2CBalance)
	}
}

func TestA2CBalanceNilRequest(t *testing.T) {
	client := NewDefaultClient()

	response, err := client.A2CBalance(nil)
	if !errors.Is(err, ErrRequestIsNil) {
		t.Fatalf("A2CBalance() error = %v, want ErrRequestIsNil", err)
	}
	if response != nil {
		t.Fatalf("A2CBalance() response = %#v, want nil", response)
	}
}

func TestA2CBalanceProviderError(t *testing.T) {
	roundTripper := teststand.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return teststand.Response(
			http.StatusOK,
			"application/json",
			[]byte(`{"response":{"error":"overall error","error_code":"U0"}}`),
		), nil
	})
	client := NewClient(WithClient(&http.Client{Transport: roundTripper}))

	response, err := client.A2CBalance(balanceRequest())
	if err == nil {
		t.Fatal("A2CBalance() error = nil")
	}
	if response == nil {
		t.Fatal("A2CBalance() response = nil, want structured provider response")
	}

	var providerErr *ipay.IpayError
	if !errors.As(err, &providerErr) {
		t.Fatalf("A2CBalance() error = %T, want *ipay.IpayError", err)
	}
	if ipay.RequiresA2CStatusCheck(err) {
		t.Fatalf("A2CBalance() error unexpectedly requires A2C status check: %v", err)
	}
}

func TestA2CBalanceRejectsIncompleteSuccessResponse(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing current balance",
			body: `{"response":{"overdraft":0,"credit":0,"salt":"salt","sign":"sign"}}`,
		},
		{
			name: "missing overdraft",
			body: `{"response":{"current_balance":0,"credit":0,"salt":"salt","sign":"sign"}}`,
		},
		{
			name: "missing credit",
			body: `{"response":{"current_balance":0,"overdraft":0,"salt":"salt","sign":"sign"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roundTripper := teststand.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
				return teststand.Response(http.StatusOK, "application/json", []byte(tt.body)), nil
			})
			client := NewClient(WithClient(&http.Client{Transport: roundTripper}))

			response, err := client.A2CBalance(balanceRequest())
			if !errors.Is(err, ipay.ErrDecode) {
				t.Fatalf("A2CBalance() error = %v, want ErrDecode", err)
			}
			if response == nil {
				t.Fatal("A2CBalance() response = nil, want decoded partial response")
			}
			if ipay.RequiresA2CStatusCheck(err) {
				t.Fatalf("A2CBalance() validation error unexpectedly requires A2C status check: %v", err)
			}

			var decodeErr *ipay.DecodeError
			if !errors.As(err, &decodeErr) {
				t.Fatalf("A2CBalance() error = %T, want *ipay.DecodeError", err)
			}
			if decodeErr.Action != ipay.ActionA2CBalance || decodeErr.Operation != consts.A2CBalance {
				t.Fatalf("decode context = action:%q operation:%q", decodeErr.Action, decodeErr.Operation)
			}
		})
	}
}

func TestA2CBalanceAcceptsZeroValues(t *testing.T) {
	roundTripper := teststand.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return teststand.Response(
			http.StatusOK,
			"application/json",
			[]byte(`{"response":{"current_balance":0,"overdraft":0,"credit":0,"salt":"salt","sign":"sign"}}`),
		), nil
	})
	client := NewClient(WithClient(&http.Client{Transport: roundTripper}))

	response, err := client.A2CBalance(balanceRequest())
	if err != nil {
		t.Fatalf("A2CBalance() error: %v", err)
	}
	if response == nil || !response.HasA2CBalanceFields() {
		t.Fatalf("A2CBalance() response = %#v, want complete zero balance", response)
	}
}

func balanceRequest() *Request {
	return &Request{
		Merchant: &Merchant{
			MerchantID:  "2995",
			MerchantKey: "merchant-key",
		},
	}
}
