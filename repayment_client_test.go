package go_ipay

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stremovskyy/go-ipay/consts"
	"github.com/stremovskyy/go-ipay/repayment"
)

func TestEncodeRepaymentTransactionsCSV(t *testing.T) {
	got, err := encodeRepaymentTransactionsCSV([]RepaymentTransaction{
		{PmtID: 1, ExtID: "a"},
		{PmtID: 2, ExtID: "b"},
	})
	if err != nil {
		t.Fatalf("encodeRepaymentTransactionsCSV() error: %v", err)
	}

	want := "1;a\n2;b\n"
	if string(got) != want {
		t.Fatalf("encodeRepaymentTransactionsCSV() = %q, want %q", string(got), want)
	}
}

func TestCreateRepayment_DryRun(t *testing.T) {
	cl := NewDefaultClient()

	var gotEndpoint string
	var gotPayload any

	_, err := cl.CreateRepayment(&CreateRepaymentRequest{
		Merchant: &Merchant{
			MerchantID:   "123",
			Login:        "test-login",
			RepaymentKey: "test-key",
		},
		ExtID: "8ae61d49-9d31-4390-9a12-532590f00429",
		Transactions: []RepaymentTransaction{
			{PmtID: 867305735, ExtID: "4207d62e-d7e4-4fe5-b210-49bffc4cd5ce"},
		},
	}, DryRun(func(endpoint string, payload any) {
		gotEndpoint = endpoint
		gotPayload = payload
	}))
	if err != nil {
		t.Fatalf("CreateRepayment() error: %v", err)
	}

	if gotEndpoint != consts.RepaymentUrl {
		t.Fatalf("endpoint = %q, want %q", gotEndpoint, consts.RepaymentUrl)
	}

	raw, err := json.Marshal(gotPayload)
	if err != nil {
		t.Fatalf("marshal dry-run payload: %v", err)
	}

	var payload struct {
		Operation string            `json:"operation"`
		Request   repayment.Request `json:"request"`
		FileName  string            `json:"file_name"`
		AuthTime  string            `json:"auth_time"`
		AuthSign  string            `json:"auth_sign"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal dry-run payload: %v", err)
	}

	if payload.Operation != consts.CreateRepayment {
		t.Fatalf("operation = %q, want %q", payload.Operation, consts.CreateRepayment)
	}
	if payload.Request.Action != repayment.ActionCreateRepayment {
		t.Fatalf("action = %q, want %q", payload.Request.Action, repayment.ActionCreateRepayment)
	}
	if payload.FileName != "transactions.csv" {
		t.Fatalf("file_name = %q, want %q", payload.FileName, "transactions.csv")
	}
	if len(payload.AuthSign) != 64 {
		t.Fatalf("auth_sign len = %d, want 64", len(payload.AuthSign))
	}
	if _, err := time.Parse("2006-01-02 15:04:05", payload.AuthTime); err != nil {
		t.Fatalf("auth_time parse error: %v", err)
	}
}

func TestCreateRepayment_RejectsFileAndTransactions(t *testing.T) {
	cl := NewDefaultClient()

	_, err := cl.CreateRepayment(&CreateRepaymentRequest{
		Merchant: &Merchant{
			MerchantID:   "123",
			Login:        "test-login",
			RepaymentKey: "test-key",
		},
		ExtID:                "ext",
		TransactionsFilePath: "/tmp/does-not-matter.csv",
		Transactions: []RepaymentTransaction{
			{PmtID: 1, ExtID: "a"},
		},
	}, DryRun())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCancelRepayment_DryRun(t *testing.T) {
	cl := NewDefaultClient()

	guid := "68D1D550-0BC9-4BE7-9A44-964A0E2AE3A2"

	var gotEndpoint string
	var gotPayload any

	_, err := cl.CancelRepayment(&CancelRepaymentRequest{
		Merchant: &Merchant{
			Login:        "test-login",
			RepaymentKey: "test-key",
		},
		RepaymentGUID: &guid,
	}, DryRun(func(endpoint string, payload any) {
		gotEndpoint = endpoint
		gotPayload = payload
	}))
	if err != nil {
		t.Fatalf("CancelRepayment() error: %v", err)
	}

	if gotEndpoint != consts.RepaymentUrl {
		t.Fatalf("endpoint = %q, want %q", gotEndpoint, consts.RepaymentUrl)
	}

	wrapper, ok := gotPayload.(*repayment.RequestWrapper)
	if !ok {
		t.Fatalf("payload type = %T, want *repayment.RequestWrapper", gotPayload)
	}

	if wrapper.Operation != consts.CancelRepayment {
		t.Fatalf("operation = %q, want %q", wrapper.Operation, consts.CancelRepayment)
	}
	if wrapper.Request.Action != repayment.ActionCancelRepayment {
		t.Fatalf("action = %q, want %q", wrapper.Request.Action, repayment.ActionCancelRepayment)
	}

	if wrapper.Request.Body.RepaymentGUID == nil || *wrapper.Request.Body.RepaymentGUID != guid {
		t.Fatalf("repayment_guid = %v, want %q", wrapper.Request.Body.RepaymentGUID, guid)
	}
	if wrapper.Request.Body.ExtID != nil {
		t.Fatalf("ext_id should be nil, got %v", *wrapper.Request.Body.ExtID)
	}

	if wrapper.Request.Auth.Login != "test-login" {
		t.Fatalf("auth.login = %q, want %q", wrapper.Request.Auth.Login, "test-login")
	}
	if len(wrapper.Request.Auth.Sign) != 64 {
		t.Fatalf("auth.sign len = %d, want 64", len(wrapper.Request.Auth.Sign))
	}
	if _, err := time.Parse("2006-01-02 15:04:05", wrapper.Request.Auth.Time); err != nil {
		t.Fatalf("auth.time parse error: %v", err)
	}
}

func TestGetRepaymentStatus_DryRun(t *testing.T) {
	cl := NewDefaultClient()

	extID := "8ae61d49-9d31-4390-9a12-532590f00429"

	var gotEndpoint string
	var gotPayload any

	_, err := cl.GetRepaymentStatus(&GetRepaymentStatusRequest{
		Merchant: &Merchant{
			Login:        "test-login",
			RepaymentKey: "test-key",
		},
		ExtID: &extID,
	}, DryRun(func(endpoint string, payload any) {
		gotEndpoint = endpoint
		gotPayload = payload
	}))
	if err != nil {
		t.Fatalf("GetRepaymentStatus() error: %v", err)
	}

	if gotEndpoint != consts.RepaymentUrl {
		t.Fatalf("endpoint = %q, want %q", gotEndpoint, consts.RepaymentUrl)
	}

	wrapper, ok := gotPayload.(*repayment.RequestWrapper)
	if !ok {
		t.Fatalf("payload type = %T, want *repayment.RequestWrapper", gotPayload)
	}

	if wrapper.Operation != consts.GetRepaymentStatus {
		t.Fatalf("operation = %q, want %q", wrapper.Operation, consts.GetRepaymentStatus)
	}
	if wrapper.Request.Action != repayment.ActionGetRepaymentStatus {
		t.Fatalf("action = %q, want %q", wrapper.Request.Action, repayment.ActionGetRepaymentStatus)
	}

	if wrapper.Request.Body.ExtID == nil || *wrapper.Request.Body.ExtID != extID {
		t.Fatalf("ext_id = %v, want %q", wrapper.Request.Body.ExtID, extID)
	}
	if wrapper.Request.Body.RepaymentGUID != nil {
		t.Fatalf("repayment_guid should be nil, got %v", *wrapper.Request.Body.RepaymentGUID)
	}

	if wrapper.Request.Auth.Login != "test-login" {
		t.Fatalf("auth.login = %q, want %q", wrapper.Request.Auth.Login, "test-login")
	}
	if len(wrapper.Request.Auth.Sign) != 64 {
		t.Fatalf("auth.sign len = %d, want 64", len(wrapper.Request.Auth.Sign))
	}
	if _, err := time.Parse("2006-01-02 15:04:05", wrapper.Request.Auth.Time); err != nil {
		t.Fatalf("auth.time parse error: %v", err)
	}
}

func TestRepaymentLookupValidation(t *testing.T) {
	cl := NewDefaultClient()

	guid := "guid"
	extID := "ext"

	_, err := cl.GetRepaymentStatus(&GetRepaymentStatusRequest{
		Merchant: &Merchant{Login: "l", RepaymentKey: "k"},
	}, DryRun())
	if err == nil {
		t.Fatalf("expected error for missing lookup fields, got nil")
	}

	_, err = cl.GetRepaymentStatus(&GetRepaymentStatusRequest{
		Merchant:      &Merchant{Login: "l", RepaymentKey: "k"},
		RepaymentGUID: &guid,
		ExtID:         &extID,
	}, DryRun())
	if err == nil {
		t.Fatalf("expected error for both lookup fields, got nil")
	}
}

func TestGetRepaymentProcessingFile_DryRun(t *testing.T) {
	cl := NewDefaultClient()

	extID := "8ae61d49-9d31-4390-9a12-532590f00429"

	var gotEndpoint string
	var gotPayload any

	_, err := cl.GetRepaymentProcessingFile(&GetRepaymentProcessingFileRequest{
		Merchant: &Merchant{
			Login:        "test-login",
			RepaymentKey: "test-key",
		},
		ExtID: &extID,
	}, DryRun(func(endpoint string, payload any) {
		gotEndpoint = endpoint
		gotPayload = payload
	}))
	if err != nil {
		t.Fatalf("GetRepaymentProcessingFile() error: %v", err)
	}

	if gotEndpoint != consts.RepaymentUrl {
		t.Fatalf("endpoint = %q, want %q", gotEndpoint, consts.RepaymentUrl)
	}

	wrapper, ok := gotPayload.(*repayment.RequestWrapper)
	if !ok {
		t.Fatalf("payload type = %T, want *repayment.RequestWrapper", gotPayload)
	}

	if wrapper.Operation != consts.GetRepaymentProcessingFile {
		t.Fatalf("operation = %q, want %q", wrapper.Operation, consts.GetRepaymentProcessingFile)
	}
	if wrapper.Request.Action != repayment.ActionGetRepaymentProcessingFile {
		t.Fatalf("action = %q, want %q", wrapper.Request.Action, repayment.ActionGetRepaymentProcessingFile)
	}
	if wrapper.Request.Body.ExtID == nil || *wrapper.Request.Body.ExtID != extID {
		t.Fatalf("ext_id = %v, want %q", wrapper.Request.Body.ExtID, extID)
	}
}
