package repayment

import "testing"

func TestUnmarshalJSONResponse_WithStringNumerics(t *testing.T) {
	resp, err := UnmarshalJSONResponse([]byte(`{"response":{"repayment_guid":"69CFC09C-12DF-4BB1-9702-6D1E4BB22554","ext_id":"ext-1","status":"5","invoice":"0","amount":"0","mch_id":"4769","mch_balance":"362165220","success_payments":"3","failed_payments":"0"}}`))
	if err != nil {
		t.Fatalf("expected string numeric payload to parse, got error: %v", err)
	}
	if resp == nil || resp.Status == nil || *resp.Status != 5 {
		t.Fatalf("expected status=5, got %#v", resp)
	}
	if resp.MchID == nil || *resp.MchID != 4769 {
		t.Fatalf("expected mch_id=4769, got %#v", resp.MchID)
	}
	if resp.SuccessPayments == nil || *resp.SuccessPayments != 3 {
		t.Fatalf("expected success_payments=3, got %#v", resp.SuccessPayments)
	}
}
