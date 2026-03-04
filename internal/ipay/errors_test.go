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
