package ipay

import (
	"math"
	"reflect"
	"testing"
)

func TestUnmarshalJSONResponseApplePayNumericPmtStatus(t *testing.T) {
	raw := []byte(`{
		"response": {
			"pmt_id": 1234567,
			"invoice": 100,
			"amount": 100,
			"pmt_status": 1,
			"security_rate": "3D",
			"security_data": {
				"redirect_url": "https://example.com/acs"
			},
			"expected_aml_fields": ["sender.firstname", "receiver.lastname"],
			"recurrent_token": "recurrent-token",
			"rrn": "rrn-123",
			"terminal_id": "terminal-123",
			"terminal_merchant_id": "terminal-merchant-123",
			"auth_code": "auth-123",
			"mch_amount": [
				{"smch_id": "112233", "amount": "100"}
			],
			"transactions": [
				{"trn_id": 7654321, "smch_id": 112233, "invoice": 100, "amount": 100}
			]
		}
	}`)

	resp, err := UnmarshalJSONResponse(raw)
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}
	if got := resp.GetPaymentStatus(); got != PaymentStatusRegistered {
		t.Fatalf("GetPaymentStatus() = %v, want %v", got, PaymentStatusRegistered)
	}
	if got := resp.PmtIdInt64(); got != 1234567 {
		t.Fatalf("PmtIdInt64() = %d, want 1234567", got)
	}

	assertStringPointerField(t, resp, "SecurityRate", "3D")
	assertNestedStringPointerField(t, resp, "SecurityData", "RedirectURL", "https://example.com/acs")
	assertStringSliceField(t, resp, "ExpectedAMLFields", []string{"sender.firstname", "receiver.lastname"})
	assertStringPointerField(t, resp, "RecurrentToken", "recurrent-token")
	assertStringPointerField(t, resp, "RRN", "rrn-123")
	assertStringPointerField(t, resp, "TerminalID", "terminal-123")
	assertStringPointerField(t, resp, "TerminalMerchantID", "terminal-merchant-123")
	assertStringPointerField(t, resp, "AuthCode", "auth-123")
	assertMchAmountField(t, resp)
	assertTransactionSmchIDField(t, resp)
}

func TestUnmarshalJSONResponseApplePayStringPmtStatus(t *testing.T) {
	resp, err := UnmarshalJSONResponse([]byte(`{"response":{"pmt_id":"1234567","pmt_status":"5"}}`))
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}
	if got := resp.GetPaymentStatus(); got != PaymentStatusSuccess {
		t.Fatalf("GetPaymentStatus() = %v, want %v", got, PaymentStatusSuccess)
	}
}

func TestUnmarshalJSONResponseExtIDCompatibilityShapes(t *testing.T) {
	tests := []struct {
		name  string
		extID string
		want  *string
	}{
		{
			name:  "legacy string ext_id",
			extID: `"56bff39a-e594-4f17-b7be-6096bb5c2b50"`,
			want:  refString("56bff39a-e594-4f17-b7be-6096bb5c2b50"),
		},
		{
			name:  "numeric ext_id",
			extID: `1108174424`,
			want:  refString("1108174424"),
		},
		{
			name:  "structured ext_id",
			extID: `{"ext_id":"56bff39a-e594-4f17-b7be-6096bb5c2b50","mch_id":"4767","pmt_id":"1108174424"}`,
			want:  refString("56bff39a-e594-4f17-b7be-6096bb5c2b50"),
		},
		{
			name:  "structured numeric ext_id",
			extID: `{"ext_id":1108174424,"mch_id":"4767","pmt_id":"1108174424"}`,
			want:  refString("1108174424"),
		},
		{
			name:  "null ext_id",
			extID: `null`,
			want:  nil,
		},
		{
			name:  "unsupported structured ext_id",
			extID: `{"mch_id":"4767","pmt_id":"1108174424"}`,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := UnmarshalJSONResponse([]byte(`{
				"response": {
					"ext_id": ` + tt.extID + `,
					"pmt_id": "1108174424",
					"pmt_status": "3",
					"invoice": "100",
					"amount": "100"
				}
			}`))
			if err != nil {
				t.Fatalf("UnmarshalJSONResponse() error: %v", err)
			}
			if !sameStringPointer(resp.ExtId, tt.want) {
				t.Fatalf("ExtId = %v, want %v", resp.ExtId, tt.want)
			}
			if got := resp.PmtIdInt64(); got != 1108174424 {
				t.Fatalf("PmtIdInt64() = %d, want 1108174424", got)
			}
			if got := resp.GetPaymentStatus(); got != PaymentStatusPreAuthorized {
				t.Fatalf("GetPaymentStatus() = %v, want %v", got, PaymentStatusPreAuthorized)
			}
			if got := resp.InvoiceAmountInt64(); got != 100 {
				t.Fatalf("InvoiceAmountInt64() = %d, want 100", got)
			}
			if got := resp.AmountInt64(); got != 100 {
				t.Fatalf("AmountInt64() = %d, want 100", got)
			}
		})
	}
}

func TestUnmarshalJSONResponseMchAmountCompatibilityShapes(t *testing.T) {
	tests := []struct {
		name       string
		mchAmount  string
		wantSmchID *string
		wantAmount *string
	}{
		{
			name:       "sentry mobile payment numeric amount",
			mchAmount:  `[{"smch_id":"13581","amount":97}]`,
			wantSmchID: refString("13581"),
			wantAmount: refString("97"),
		},
		{
			name:       "legacy string amount",
			mchAmount:  `[{"smch_id":"112233","amount":"100"}]`,
			wantSmchID: refString("112233"),
			wantAmount: refString("100"),
		},
		{
			name:       "numeric smch id and amount",
			mchAmount:  `[{"smch_id":13581,"amount":97}]`,
			wantSmchID: refString("13581"),
			wantAmount: refString("97"),
		},
		{
			name:       "null fields",
			mchAmount:  `[{"smch_id":null,"amount":null}]`,
			wantSmchID: nil,
			wantAmount: nil,
		},
		{
			name:       "unsupported object fields",
			mchAmount:  `[{"smch_id":{"id":"13581"},"amount":{"value":97}}]`,
			wantSmchID: nil,
			wantAmount: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := UnmarshalJSONResponse([]byte(`{
				"response": {
					"pmt_id": "1108644127",
					"pmt_status": "5",
					"invoice": "97",
					"amount": "97",
					"mch_amount": ` + tt.mchAmount + `
				}
			}`))
			if err != nil {
				t.Fatalf("UnmarshalJSONResponse() error: %v", err)
			}
			if len(resp.MchAmount) != 1 {
				t.Fatalf("len(MchAmount) = %d, want 1", len(resp.MchAmount))
			}
			if !sameStringPointer(resp.MchAmount[0].SmchID, tt.wantSmchID) {
				t.Fatalf("MchAmount[0].SmchID = %v, want %v", resp.MchAmount[0].SmchID, tt.wantSmchID)
			}
			if !sameStringPointer(resp.MchAmount[0].Amount, tt.wantAmount) {
				t.Fatalf("MchAmount[0].Amount = %v, want %v", resp.MchAmount[0].Amount, tt.wantAmount)
			}
			if got := resp.GetPaymentStatus(); got != PaymentStatusSuccess {
				t.Fatalf("GetPaymentStatus() = %v, want %v", got, PaymentStatusSuccess)
			}
		})
	}
}

func TestResponseNumericHelpersAcceptManualIntegerTypes(t *testing.T) {
	resp := Response{
		PmtId:   int64(1234567),
		Invoice: int(100),
		Amount:  int64(125),
	}

	if got := resp.PmtIdInt64(); got != 1234567 {
		t.Fatalf("PmtIdInt64() = %d, want 1234567", got)
	}
	if got := resp.InvoiceAmountInt64(); got != 100 {
		t.Fatalf("InvoiceAmountInt64() = %d, want 100", got)
	}
	if got := resp.AmountInt64(); got != 125 {
		t.Fatalf("AmountInt64() = %d, want 125", got)
	}
}

func TestResponseNumericHelpersRejectOverflowingUint(t *testing.T) {
	if uint64(^uint(0)) <= math.MaxInt64 {
		t.Skip("uint cannot exceed int64 on this architecture")
	}

	resp := Response{PmtId: ^uint(0)}
	if got := resp.PmtIdInt64(); got != 0 {
		t.Fatalf("PmtIdInt64() = %d, want 0 for overflowing uint", got)
	}
}

func assertStringPointerField(t *testing.T, resp *Response, fieldName string, want string) {
	t.Helper()

	field := responseField(t, resp, fieldName)
	if field.IsNil() {
		t.Fatalf("%s is nil, want %q", fieldName, want)
	}
	if got := field.Elem().String(); got != want {
		t.Fatalf("%s = %q, want %q", fieldName, got, want)
	}
}

func assertNestedStringPointerField(t *testing.T, resp *Response, fieldName string, nestedFieldName string, want string) {
	t.Helper()

	field := responseField(t, resp, fieldName)
	if field.IsNil() {
		t.Fatalf("%s is nil, want nested %s=%q", fieldName, nestedFieldName, want)
	}

	nested := field.Elem().FieldByName(nestedFieldName)
	if !nested.IsValid() {
		t.Fatalf("%s.%s field is missing", fieldName, nestedFieldName)
	}
	if nested.IsNil() {
		t.Fatalf("%s.%s is nil, want %q", fieldName, nestedFieldName, want)
	}
	if got := nested.Elem().String(); got != want {
		t.Fatalf("%s.%s = %q, want %q", fieldName, nestedFieldName, got, want)
	}
}

func assertStringSliceField(t *testing.T, resp *Response, fieldName string, want []string) {
	t.Helper()

	field := responseField(t, resp, fieldName)
	if field.Len() != len(want) {
		t.Fatalf("len(%s) = %d, want %d", fieldName, field.Len(), len(want))
	}
	for i := range want {
		if got := field.Index(i).String(); got != want[i] {
			t.Fatalf("%s[%d] = %q, want %q", fieldName, i, got, want[i])
		}
	}
}

func assertMchAmountField(t *testing.T, resp *Response) {
	t.Helper()

	field := responseField(t, resp, "MchAmount")
	if field.Len() != 1 {
		t.Fatalf("len(MchAmount) = %d, want 1", field.Len())
	}

	row := field.Index(0)
	smchID := row.FieldByName("SmchID")
	amount := row.FieldByName("Amount")
	if !smchID.IsValid() || !amount.IsValid() {
		t.Fatalf("MchAmount row fields are missing: %#v", row)
	}
	if smchID.IsNil() || smchID.Elem().String() != "112233" {
		t.Fatalf("MchAmount[0].SmchID = %#v, want 112233", smchID)
	}
	if amount.IsNil() || amount.Elem().String() != "100" {
		t.Fatalf("MchAmount[0].Amount = %#v, want 100", amount)
	}
}

func assertTransactionSmchIDField(t *testing.T, resp *Response) {
	t.Helper()

	if len(resp.Transactions) != 1 {
		t.Fatalf("len(Transactions) = %d, want 1", len(resp.Transactions))
	}

	field := reflect.ValueOf(resp.Transactions[0]).FieldByName("SmchID")
	if !field.IsValid() {
		t.Fatalf("ResponseTransaction.SmchID field is missing")
	}
	if field.IsNil() || int(field.Elem().Int()) != 112233 {
		t.Fatalf("Transactions[0].SmchID = %#v, want 112233", field)
	}
}

func responseField(t *testing.T, resp *Response, fieldName string) reflect.Value {
	t.Helper()

	field := reflect.ValueOf(resp).Elem().FieldByName(fieldName)
	if !field.IsValid() {
		t.Fatalf("Response.%s field is missing", fieldName)
	}

	return field
}

func refString(value string) *string {
	return &value
}

func sameStringPointer(got *string, want *string) bool {
	if got == nil || want == nil {
		return got == want
	}

	return *got == *want
}
