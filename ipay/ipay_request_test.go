package ipay

import (
	"encoding/json"
	"testing"
)

func TestRequestMarshalA2CBalanceOmitsBody(t *testing.T) {
	raw, err := json.Marshal(NewRequest(ActionA2CBalance))
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	request := requestObject(t, raw)
	if _, exists := request["body"]; exists {
		t.Fatalf("A2CBalance body must be omitted: %s", string(raw))
	}
}

func TestRequestMarshalKeepsEmptyBodyForOtherActions(t *testing.T) {
	raw, err := json.Marshal(NewRequest(ActionGetPaymentStatus))
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	request := requestObject(t, raw)
	body, exists := request["body"]
	if !exists || string(body) != `{}` {
		t.Fatalf("request.body = %s, want {}", string(body))
	}
}

func TestRequestMarshalKeepsExplicitZeroBodyValue(t *testing.T) {
	raw, err := json.Marshal(NewRequest(ActionCredit, WithInvoiceAmount(0)))
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	request := requestObject(t, raw)
	var body map[string]json.RawMessage
	if err := json.Unmarshal(request["body"], &body); err != nil {
		t.Fatalf("unmarshal request.body: %v", err)
	}

	var invoice int
	if err := json.Unmarshal(body["invoice"], &invoice); err != nil {
		t.Fatalf("unmarshal body.invoice: %v", err)
	}
	if invoice != 0 {
		t.Fatalf("body.invoice = %d, want 0", invoice)
	}
}

func requestObject(t *testing.T, raw []byte) map[string]json.RawMessage {
	t.Helper()

	var envelope struct {
		Request map[string]json.RawMessage `json:"request"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("unmarshal request envelope: %v", err)
	}

	return envelope.Request
}
