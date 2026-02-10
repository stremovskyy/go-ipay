package go_ipay

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stremovskyy/go-ipay/consts"
	repayinternal "github.com/stremovskyy/go-ipay/internal/repayment"
	"github.com/stremovskyy/go-ipay/internal/teststand"
	"github.com/stremovskyy/go-ipay/repayment"
)

func TestRepayment_CreateRepayment_HTTPShape(t *testing.T) {
	const (
		login = "test-login"
		key   = "test-key"
		mchID = int64(2023)
		extID = "8ae61d49-9d31-4390-9a12-532590f00422"
	)

	tx := RepaymentTransaction{
		PmtID: 867305735,
		ExtID: "4207d62e-d7e4-4fe5-b210-49bffc4cd5ce",
	}

	rt := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
		}
		if req.URL.Scheme != "https" || req.URL.Host != "api-repayment.ipay.ua" {
			t.Fatalf("url = %q, want https://api-repayment.ipay.ua", req.URL.String())
		}

		ct := req.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(ct)
		if err != nil {
			t.Fatalf("parse Content-Type %q: %v", ct, err)
		}
		if mediaType != "multipart/form-data" {
			t.Fatalf("Content-Type media = %q, want %q", mediaType, "multipart/form-data")
		}
		boundary := params["boundary"]
		if boundary == "" {
			t.Fatalf("multipart boundary is empty (Content-Type=%q)", ct)
		}

		raw, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		_ = req.Body.Close()

		r := multipart.NewReader(bytes.NewReader(raw), boundary)

		var (
			gotRequestJSON []byte
			gotFileName    string
			gotFileCT      string
			gotFile        []byte
		)

		for {
			part, err := r.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				t.Fatalf("NextPart: %v", err)
			}

			b, err := io.ReadAll(part)
			if err != nil {
				t.Fatalf("read part %q: %v", part.FormName(), err)
			}

			switch part.FormName() {
			case "request":
				gotRequestJSON = b
				gotCT := part.Header.Get("Content-Type")
				if gotCT != "application/json; charset=utf-8" {
					t.Fatalf("request part Content-Type = %q, want %q", gotCT, "application/json; charset=utf-8")
				}
			case "file":
				gotFileName = part.FileName()
				gotFileCT = part.Header.Get("Content-Type")
				gotFile = b
			default:
				t.Fatalf("unexpected multipart field: %q", part.FormName())
			}
		}

		if len(gotRequestJSON) == 0 {
			t.Fatalf("missing request part")
		}
		if gotFileName == "" {
			t.Fatalf("missing file part")
		}
		if gotFileName != "transactions.csv" {
			t.Fatalf("file name = %q, want %q", gotFileName, "transactions.csv")
		}
		if gotFileCT != "text/csv" {
			t.Fatalf("file Content-Type = %q, want %q", gotFileCT, "text/csv")
		}

		var wrapper repayment.RequestWrapper
		if err := json.Unmarshal(gotRequestJSON, &wrapper); err != nil {
			t.Fatalf("unmarshal request JSON: %v\nraw=%s", err, string(gotRequestJSON))
		}

		if wrapper.Request.Action != repayment.ActionCreateRepayment {
			t.Fatalf("action = %q, want %q", wrapper.Request.Action, repayment.ActionCreateRepayment)
		}
		if wrapper.Request.Auth.Login != login {
			t.Fatalf("auth.login = %q, want %q", wrapper.Request.Auth.Login, login)
		}
		if _, err := time.Parse("2006-01-02 15:04:05", wrapper.Request.Auth.Time); err != nil {
			t.Fatalf("auth.time parse error: %v (time=%q)", err, wrapper.Request.Auth.Time)
		}
		wantSign := repayinternal.SignSHA256Hex(wrapper.Request.Auth.Time, key)
		if wrapper.Request.Auth.Sign != wantSign {
			t.Fatalf("auth.sign = %q, want %q", wrapper.Request.Auth.Sign, wantSign)
		}

		if wrapper.Request.Body.MchID == nil || *wrapper.Request.Body.MchID != mchID {
			t.Fatalf("body.mch_id = %v, want %d", wrapper.Request.Body.MchID, mchID)
		}
		if wrapper.Request.Body.ExtID == nil || *wrapper.Request.Body.ExtID != extID {
			t.Fatalf("body.ext_id = %v, want %q", wrapper.Request.Body.ExtID, extID)
		}
		if wrapper.Request.Body.SmchID != nil {
			t.Fatalf("body.smch_id = %v, want nil", *wrapper.Request.Body.SmchID)
		}

		wantCSV := "867305735;4207d62e-d7e4-4fe5-b210-49bffc4cd5ce\n"
		if string(gotFile) != wantCSV {
			t.Fatalf("CSV = %q, want %q", string(gotFile), wantCSV)
		}

		responseBody := []byte(`{"response":{"repayment_guid":"68D1D550-0BC9-4BE7-9A44-964A0E2AE3A2","ext_id":"` + extID + `","status":5,"invoice":13800,"amount":13800,"mch_id":2624,"mch_balance":3485750,"success_payments":1,"failed_payments":0}}`)
		return teststand.Response(200, "application/json", responseBody), nil
	})

	httpClient := &http.Client{Transport: rt}
	cl := NewClient(WithClient(httpClient))

	resp, err := cl.CreateRepayment(&CreateRepaymentRequest{
		Merchant: &Merchant{
			MerchantID:   "2023",
			Login:        login,
			RepaymentKey: key,
		},
		MchID:        mchID,
		ExtID:        extID,
		Transactions: []RepaymentTransaction{tx},
	})
	if err != nil {
		t.Fatalf("CreateRepayment() error: %v", err)
	}
	if resp == nil || resp.RepaymentGUID == nil || *resp.RepaymentGUID == "" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestRepayment_CreateRepayment_NonJSONResponse(t *testing.T) {
	const (
		login = "test-login"
		key   = "test-key"
		mchID = int64(2023)
		extID = "8ae61d49-9d31-4390-9a12-532590f00422"
	)

	rt := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		body := []byte("An error occurred during the request, please try again later.")
		return teststand.Response(200, "text/plain", body), nil
	})

	httpClient := &http.Client{Transport: rt}
	cl := NewClient(WithClient(httpClient))

	_, err := cl.CreateRepayment(&CreateRepaymentRequest{
		Merchant: &Merchant{
			MerchantID:   "2023",
			Login:        login,
			RepaymentKey: key,
		},
		MchID: mchID,
		ExtID: extID,
		Transactions: []RepaymentTransaction{
			{
				PmtID: 867305735,
				ExtID: "4207d62e-d7e4-4fe5-b210-49bffc4cd5ce",
			},
		},
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "An error occurred during the request") {
		t.Fatalf("error = %q, want it to include the response body", err.Error())
	}
	if strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("error should not be a JSON unmarshal error, got %q", err.Error())
	}

	var apiErr *repayment.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want repayment.APIError in chain", err)
	}
}

func TestRepayment_CancelRepayment_HTTPShape(t *testing.T) {
	const (
		login = "test-login"
		key   = "test-key"
	)
	guid := "68D1D550-0BC9-4BE7-9A44-964A0E2AE3A2"

	rt := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
		}
		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want %q", got, "application/json")
		}

		raw, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		_ = req.Body.Close()

		var wrapper repayment.RequestWrapper
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			t.Fatalf("unmarshal JSON: %v (raw=%s)", err, string(raw))
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

		if wrapper.Request.Auth.Login != login {
			t.Fatalf("auth.login = %q, want %q", wrapper.Request.Auth.Login, login)
		}
		if _, err := time.Parse("2006-01-02 15:04:05", wrapper.Request.Auth.Time); err != nil {
			t.Fatalf("auth.time parse error: %v (time=%q)", err, wrapper.Request.Auth.Time)
		}
		wantSign := repayinternal.SignSHA256Hex(wrapper.Request.Auth.Time, key)
		if wrapper.Request.Auth.Sign != wantSign {
			t.Fatalf("auth.sign = %q, want %q", wrapper.Request.Auth.Sign, wantSign)
		}

		responseBody := []byte(`{"response":{"repayment_guid":"` + guid + `","ext_id":"ext","status":9,"invoice":13800,"amount":13800,"mch_id":2624,"mch_balance":3485750}}`)
		return teststand.Response(200, "application/json", responseBody), nil
	})

	httpClient := &http.Client{Transport: rt}
	cl := NewClient(WithClient(httpClient))

	resp, err := cl.CancelRepayment(&CancelRepaymentRequest{
		Merchant: &Merchant{
			Login:        login,
			RepaymentKey: key,
		},
		RepaymentGUID: &guid,
	})
	if err != nil {
		t.Fatalf("CancelRepayment() error: %v", err)
	}
	if resp == nil || resp.Status == nil || *resp.Status != 9 {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestRepayment_GetRepaymentStatus_HTTPShape(t *testing.T) {
	const (
		login = "test-login"
		key   = "test-key"
	)
	extID := "8ae61d49-9d31-4390-9a12-532590f00429"

	rt := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want %q", got, "application/json")
		}

		raw, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		_ = req.Body.Close()

		var wrapper repayment.RequestWrapper
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			t.Fatalf("unmarshal JSON: %v (raw=%s)", err, string(raw))
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

		wantSign := repayinternal.SignSHA256Hex(wrapper.Request.Auth.Time, key)
		if wrapper.Request.Auth.Sign != wantSign {
			t.Fatalf("auth.sign = %q, want %q", wrapper.Request.Auth.Sign, wantSign)
		}

		responseBody := []byte(`{"response":{"repayment_guid":"guid","ext_id":"` + extID + `","status":5,"invoice":13800,"amount":13800,"mch_id":2624,"mch_balance":3485750,"success_payments":3,"failed_payments":0}}`)
		return teststand.Response(200, "application/json", responseBody), nil
	})

	httpClient := &http.Client{Transport: rt}
	cl := NewClient(WithClient(httpClient))

	resp, err := cl.GetRepaymentStatus(&GetRepaymentStatusRequest{
		Merchant: &Merchant{
			Login:        login,
			RepaymentKey: key,
		},
		ExtID: &extID,
	})
	if err != nil {
		t.Fatalf("GetRepaymentStatus() error: %v", err)
	}
	if resp == nil || resp.Status == nil || *resp.Status != 5 {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestRepayment_GetRepaymentProcessingFile_HTTPShape_SuccessCSV(t *testing.T) {
	const (
		login = "test-login"
		key   = "test-key"
	)
	extID := "8ae61d49-9d31-4390-9a12-532590f00429"

	csvBody := "repayment_guid;pmt_id;ext_id;status;fail_reason\n"

	rt := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want %q", got, "application/json")
		}

		accept := req.Header.Get("Accept")
		if accept == "" || accept == "application/json" {
			t.Fatalf("Accept = %q, want csv-capable Accept header", accept)
		}

		raw, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		_ = req.Body.Close()

		var wrapper repayment.RequestWrapper
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			t.Fatalf("unmarshal JSON: %v (raw=%s)", err, string(raw))
		}

		if wrapper.Request.Action != repayment.ActionGetRepaymentProcessingFile {
			t.Fatalf("action = %q, want %q", wrapper.Request.Action, repayment.ActionGetRepaymentProcessingFile)
		}
		if wrapper.Request.Body.ExtID == nil || *wrapper.Request.Body.ExtID != extID {
			t.Fatalf("ext_id = %v, want %q", wrapper.Request.Body.ExtID, extID)
		}

		wantSign := repayinternal.SignSHA256Hex(wrapper.Request.Auth.Time, key)
		if wrapper.Request.Auth.Sign != wantSign {
			t.Fatalf("auth.sign = %q, want %q", wrapper.Request.Auth.Sign, wantSign)
		}

		return teststand.Response(200, "text/csv", []byte(csvBody)), nil
	})

	httpClient := &http.Client{Transport: rt}
	cl := NewClient(WithClient(httpClient))

	raw, err := cl.GetRepaymentProcessingFile(&GetRepaymentProcessingFileRequest{
		Merchant: &Merchant{
			Login:        login,
			RepaymentKey: key,
		},
		ExtID: &extID,
	})
	if err != nil {
		t.Fatalf("GetRepaymentProcessingFile() error: %v", err)
	}
	if string(raw) != csvBody {
		t.Fatalf("csv = %q, want %q", string(raw), csvBody)
	}
}

func TestRepayment_GetRepaymentProcessingFile_HTTPShape_JSONError(t *testing.T) {
	const (
		login = "test-login"
		key   = "test-key"
	)
	guid := "68D1D550-0BC9-4BE7-9A44-964A0E2AE3A2"

	rt := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		responseBody := []byte(`{"response":{"error":"repayment not found"}}`)
		return teststand.Response(200, "application/json", responseBody), nil
	})

	httpClient := &http.Client{Transport: rt}
	cl := NewClient(WithClient(httpClient))

	raw, err := cl.GetRepaymentProcessingFile(&GetRepaymentProcessingFileRequest{
		Merchant: &Merchant{
			Login:        login,
			RepaymentKey: key,
		},
		RepaymentGUID: &guid,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if len(raw) == 0 {
		t.Fatalf("expected raw error response, got empty")
	}

	var apiErr *repayment.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want repayment.APIError in chain", err)
	}
}

func TestRepayment_ConstOperations(t *testing.T) {
	// Ensure the strings we use for logging/recording match the action names.
	if consts.CreateRepayment != string(repayment.ActionCreateRepayment) {
		t.Fatalf("CreateRepayment const mismatch: %q vs %q", consts.CreateRepayment, repayment.ActionCreateRepayment)
	}
	if consts.CancelRepayment != string(repayment.ActionCancelRepayment) {
		t.Fatalf("CancelRepayment const mismatch: %q vs %q", consts.CancelRepayment, repayment.ActionCancelRepayment)
	}
	if consts.GetRepaymentStatus != string(repayment.ActionGetRepaymentStatus) {
		t.Fatalf("GetRepaymentStatus const mismatch: %q vs %q", consts.GetRepaymentStatus, repayment.ActionGetRepaymentStatus)
	}
	if consts.GetRepaymentProcessingFile != string(repayment.ActionGetRepaymentProcessingFile) {
		t.Fatalf("GetRepaymentProcessingFile const mismatch: %q vs %q", consts.GetRepaymentProcessingFile, repayment.ActionGetRepaymentProcessingFile)
	}
}
