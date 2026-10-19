package ipay

import (
	"encoding/json"
	"testing"
)

func TestWithForceDiiaVerification(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []func(*RequestWrapper)
		want    string
	}{
		{"enabled", []func(*RequestWrapper){WithForceDiiaVerification(true)}, `{"force_diia_verification":true}`},
		{"disabled", []func(*RequestWrapper){WithForceDiiaVerification(false)}, ``},
		{"reset", []func(*RequestWrapper){WithForceDiiaVerification(true), WithForceDiiaVerification(false)}, `{}`},
		{"before info", []func(*RequestWrapper){WithForceDiiaVerification(true), WithPersonalData(&Info{MctsVts: true})}, `{"force_diia_verification":true,"mcts_vts":true}`},
		{"after info", []func(*RequestWrapper){WithPersonalData(&Info{MctsVts: true}), WithForceDiiaVerification(true)}, `{"force_diia_verification":true,"mcts_vts":true}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := json.Marshal(NewRequest(ActionCreateToken3DS, tt.options...))
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(requestObject(t, raw)["body"], &body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got := string(body["info"]); got != tt.want {
				t.Errorf("body.info = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestInfoMergeWithDiiaVerification(t *testing.T) {
	for _, initial := range []bool{false, true} {
		for _, incoming := range []bool{false, true} {
			info := &Info{ForceDiiaVerification: initial}
			info.MergeWith(&Info{ForceDiiaVerification: incoming})
			if got, want := info.ForceDiiaVerification, initial || incoming; got != want {
				t.Errorf("merge %v with %v = %v, want %v", initial, incoming, got, want)
			}
			info.MergeWith(nil)
			if info.ForceDiiaVerification != (initial || incoming) {
				t.Error("merging nil changed verification flag")
			}
		}
	}
}
