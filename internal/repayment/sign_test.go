package repayment

import "testing"

func TestSignSHA256Hex(t *testing.T) {
	got := SignSHA256Hex("2024-01-01 10:00:00", "secret")
	want := "4783f576026b50a9fa7e20ba7e10abc55c2f0d705c45d1189b37405856f35b6f"

	if got != want {
		t.Fatalf("SignSHA256Hex() = %q, want %q", got, want)
	}
}
