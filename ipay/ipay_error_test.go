package ipay

import (
	"errors"
	"strings"
	"testing"

	internalipay "github.com/stremovskyy/go-ipay/internal/ipay"
)

func TestResponseGetError_BankErrorNote66(t *testing.T) {
	note := internalipay.StatusCode("66-required_3ds")
	resp := Response{
		BnkErrorNote: &note,
	}

	err := resp.GetError()
	if err == nil {
		t.Fatalf("expected error for bnk_error_note")
	}

	var ipayErr *IpayError
	if !errors.As(err, &ipayErr) {
		t.Fatalf("expected *IpayError, got %T", err)
	}

	if ipayErr.Code != 66 {
		t.Fatalf("expected code 66, got %d", ipayErr.Code)
	}
}

func TestResponseGetError_BankErrorNote67(t *testing.T) {
	note := internalipay.StatusCode("67-card_country_not_allowed")
	resp := Response{
		BnkErrorNote: &note,
	}

	err := resp.GetError()
	if err == nil {
		t.Fatalf("expected error for bnk_error_note")
	}

	var ipayErr *IpayError
	if !errors.As(err, &ipayErr) {
		t.Fatalf("expected *IpayError, got %T", err)
	}

	if ipayErr.Code != 67 {
		t.Fatalf("expected code 67, got %d", ipayErr.Code)
	}
}

func TestResponseGetError_FailedWithBankResponseGroup42(t *testing.T) {
	status := "4"
	resp := Response{
		PmtStatus: &status,
		BankResponse: &BankResponse{
			ErrorGroup: 42,
		},
	}

	err := resp.GetError()
	if err == nil {
		t.Fatalf("expected error for failed payment status with error group")
	}

	var ipayErr *IpayError
	if !errors.As(err, &ipayErr) {
		t.Fatalf("expected *IpayError, got %T", err)
	}

	if ipayErr.Code != 42 {
		t.Fatalf("expected code 42, got %d", ipayErr.Code)
	}
	if ipayErr.Message != "Bank Error" {
		t.Fatalf("expected message %q, got %q", "Bank Error", ipayErr.Message)
	}
	if !strings.Contains(ipayErr.Details, "Insufficient funds") {
		t.Fatalf("expected details to include reason, got %q", ipayErr.Details)
	}
}

func TestResponseGetError_FailedWithUnknownBankResponseGroup(t *testing.T) {
	status := "4"
	resp := Response{
		PmtStatus: &status,
		BankResponse: &BankResponse{
			ErrorGroup: 777,
		},
	}

	err := resp.GetError()
	if err == nil {
		t.Fatalf("expected error for failed payment status with error group")
	}

	var ipayErr *IpayError
	if !errors.As(err, &ipayErr) {
		t.Fatalf("expected *IpayError, got %T", err)
	}

	if ipayErr.Code != 777 {
		t.Fatalf("expected code 777, got %d", ipayErr.Code)
	}
	if ipayErr.Message != "Payment Failed" {
		t.Fatalf("expected message %q, got %q", "Payment Failed", ipayErr.Message)
	}
}
