package ipay

import "testing"

func TestGetStatusCode_Includes66Required3DS(t *testing.T) {
	code, found := GetStatusCode(StatusCode("66-required_3ds"))
	if !found {
		t.Fatalf("expected status code 66-required_3ds to be found")
	}

	if code.ExtCode != 66 {
		t.Fatalf("expected ExtCode 66, got %d", code.ExtCode)
	}
}

func TestGetStatusCode_Includes67CardCountryNotAllowed(t *testing.T) {
	code, found := GetStatusCode(StatusCode("67-card_country_not_allowed"))
	if !found {
		t.Fatalf("expected status code 67-card_country_not_allowed to be found")
	}

	if code.ExtCode != 67 {
		t.Fatalf("expected ExtCode 67, got %d", code.ExtCode)
	}
}

func TestGetStatusCodeByExtCode(t *testing.T) {
	code, found := GetStatusCodeByExtCode(42)
	if !found {
		t.Fatalf("expected ext code 42 to be found")
	}

	if code.Code != "42-insufficient_funds" {
		t.Fatalf("unexpected code mapping %q", code.Code)
	}
}

func TestGetStatusCodeByExtCode_NotFound(t *testing.T) {
	_, found := GetStatusCodeByExtCode(999)
	if found {
		t.Fatalf("did not expect ext code 999 to be found")
	}
}
