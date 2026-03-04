package ipay

import (
	"errors"
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
