package ipay

import "testing"

func TestGetErrorMessageA2CPay_SanctionsCodes(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{
			code: 671,
			want: "Recipient or payer matches sanctions list of persons",
		},
		{
			code: 672,
			want: "Recipient or payer matches sanctions list of companies",
		},
		{
			code: 673,
			want: "Recipient or payer matches international terrorists list",
		},
	}

	for _, tc := range tests {
		got := getErrorMessageA2CPay(tc.code)
		if got != tc.want {
			t.Fatalf("unexpected message for code %d: got %q, want %q", tc.code, got, tc.want)
		}
	}
}
