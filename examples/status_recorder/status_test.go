package main

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	go_ipay "github.com/stremovskyy/go-ipay"
	"github.com/stremovskyy/go-ipay/internal/teststand"
	"github.com/stremovskyy/recorder"
	"github.com/stremovskyy/recorder/gorm_recorder"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestStatusPersistsRequestAndResponse(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "records.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get SQL database: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close database: %v", err)
		}
	})

	opts := gorm_recorder.NewOptions(func() *IpayRequestRecord { return &IpayRequestRecord{} }, func() *IpayRequestTag { return &IpayRequestTag{} }).
		WithRecordTable("custom_records").
		WithRecordColumns("id", "kind", "correlation_id").
		WithTagTable("custom_tags").
		WithTagColumns("record_ref", "tag_key", "tag_value")
	rec, err := gorm_recorder.NewRecorderWithModels(db, opts,
		recorder.WithScrubber(recorder.NewScrubber(recorder.WithDefaultReplacement("<hidden>"))))
	if err != nil {
		t.Fatalf("create recorder: %v", err)
	}

	var requestID string
	responseBody := []byte(`{"response":{"pmt_id":123,"status":5}}`)
	transport := teststand.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		requestID = req.Header.Get("X-Request-ID")
		return teststand.Response(http.StatusOK, "application/json", responseBody), nil
	})
	client := go_ipay.NewClient(go_ipay.WithRecorder(rec), go_ipay.WithClient(&http.Client{Transport: transport}))
	paymentID := int64(123)
	response, err := client.Status(&go_ipay.Request{
		Merchant:    &go_ipay.Merchant{MerchantID: "123", MerchantKey: "test-key"},
		PaymentData: &go_ipay.PaymentData{IpayPaymentID: &paymentID},
	})
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if response == nil || response.GetPaymentStatus() != 5 || requestID == "" {
		t.Fatalf("unexpected response or missing request ID: response=%+v requestID=%q", response, requestID)
	}

	request, err := rec.GetRequest(context.Background(), requestID)
	if err != nil {
		t.Fatalf("read stored request: %v", err)
	}
	var envelope struct {
		Request struct {
			Action string `json:"action"`
			Auth   string `json:"auth"`
		} `json:"request"`
	}
	if err := json.Unmarshal(request, &envelope); err != nil {
		t.Fatalf("decode stored request: %v", err)
	}
	if envelope.Request.Action != "GetPaymentStatus" || envelope.Request.Auth != "<hidden>" {
		t.Fatalf("stored request lost its action or signature scrubbing: %+v", envelope)
	}
	storedResponse, err := rec.GetResponse(context.Background(), requestID)
	if err != nil {
		t.Fatalf("read stored response: %v", err)
	}
	if string(storedResponse) != string(responseBody) {
		t.Fatalf("stored response = %s, want %s", storedResponse, responseBody)
	}
	var count int64
	if err := db.Model(&IpayRequestTag{}).Where("tag_key = ? AND tag_value = ?", "payment_id", "123").Count(&count).Error; err != nil {
		t.Fatalf("count stored payment tags: %v", err)
	}
	if count != 2 {
		t.Fatalf("payment tag count = %d, want request and response tags", count)
	}
}
