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

func TestVerificationLinkDiiaHTTPContract(t *testing.T) {
	for _, force := range []bool{false, true} {
		name := "default"
		if force {
			name = "forced"
		}
		t.Run(name, func(t *testing.T) {
			attempts := 0
			transport := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				attempts++
				if req.Method != http.MethodPost || req.URL.String() != consts.ApiUrl {
					t.Fatalf("request = %s %s, want POST %s", req.Method, req.URL, consts.ApiUrl)
				}
				if got := req.Header.Get("Content-Type"); got != "application/json" {
					t.Fatalf("Content-Type = %q, want application/json", got)
				}
				raw, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("read request: %v", err)
				}
				var envelope struct {
					Request struct {
						Action ipay.Action `json:"action"`
						Body   struct {
							Info map[string]json.RawMessage `json:"info"`
						} `json:"body"`
					} `json:"request"`
				}
				if err := json.Unmarshal(raw, &envelope); err != nil {
					t.Fatalf("decode request: %v", err)
				}
				if envelope.Request.Action != ipay.ActionCreateToken3DS {
					t.Fatalf("action = %q, want CreateToken3DS", envelope.Request.Action)
				}
				info := envelope.Request.Body.Info
				flag, exists := info["force_diia_verification"]
				if exists != force || (force && string(flag) != "true") {
					t.Fatalf("force_diia_verification = %s (present %v), want force %v", flag, exists, force)
				}
				for key, want := range map[string]string{
					"order_id": `"verification-order"`,
					"ext_id":   `"verification-order"`,
					"user_id":  `"123"`,
					"metadata": `"source:test"`,
				} {
					if got := string(info[key]); got != want {
						t.Errorf("info.%s = %s, want %s", key, got, want)
					}
				}
				var cvd map[string]string
				if err := json.Unmarshal(info["cvd"], &cvd); err != nil {
					t.Fatalf("decode CVD: %v", err)
				}
				if cvd["firstname"] != "Test" {
					t.Errorf("cvd.firstname = %q, want Test", cvd["firstname"])
				}
				return teststand.Response(http.StatusOK, "application/json",
					[]byte(`{"response":{"url":"https://sandbox.example/verify"}}`)), nil
			})
			client := NewClient(WithClient(&http.Client{Transport: transport}))
			userID := 123
			link, err := client.VerificationLink(&Request{
				Merchant: &Merchant{MerchantID: "123", MerchantKey: "test-key"},
				PaymentData: &PaymentData{
					PaymentID: refString("verification-order"),
					Metadata:  map[string]string{"source": "test"},
				},
				PersonalData:          &PersonalData{UserID: &userID, FirstName: refString("Test")},
				ForceDiiaVerification: force,
			})
			if err != nil {
				t.Fatalf("VerificationLink() error: %v", err)
			}
			if link == nil || link.String() != "https://sandbox.example/verify" {
				t.Fatalf("link = %v, want sandbox verification URL", link)
			}
			if attempts != 1 {
				t.Fatalf("HTTP attempts = %d, want 1", attempts)
			}
		})
	}
}

func TestVerificationLinkDiiaDryRun(t *testing.T) {
	attempts := 0
	transport := teststand.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		return nil, errors.New("unexpected HTTP request")
	})
	client := NewClient(WithClient(&http.Client{Transport: transport}))
	var payload *ipay.RequestWrapper
	link, err := client.VerificationLink(&Request{ForceDiiaVerification: true}, DryRun(func(endpoint string, value any) {
		if endpoint != consts.ApiUrl {
			t.Errorf("endpoint = %q, want %q", endpoint, consts.ApiUrl)
		}
		var ok bool
		payload, ok = value.(*ipay.RequestWrapper)
		if !ok {
			t.Fatalf("payload type = %T, want *ipay.RequestWrapper", value)
		}
	}))
	if err != nil || link != nil || attempts != 0 {
		t.Fatalf("dry run: link=%v err=%v attempts=%d", link, err, attempts)
	}
	if payload == nil || payload.Request.Body.Info == nil || !payload.Request.Body.Info.ForceDiiaVerification {
		t.Fatalf("dry run lost Diia verification flag: %#v", payload)
	}
}
