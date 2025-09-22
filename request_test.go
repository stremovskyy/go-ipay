package go_ipay

import "testing"

func TestRequestIsMobile(t *testing.T) {
	applePayload := "apple"
	googlePayload := "google"

	tests := []struct {
		name string
		req  Request
		want bool
	}{
		{
			name: "no payment data",
			req:  Request{},
			want: false,
		},
		{
			name: "mobile flag true",
			req: Request{
				PaymentData: &PaymentData{IsMobile: true},
			},
			want: true,
		},
		{
			name: "payment method nil",
			req: Request{
				PaymentData: &PaymentData{},
			},
			want: false,
		},
		{
			name: "apple pay container",
			req: Request{
				PaymentData:   &PaymentData{},
				PaymentMethod: &PaymentMethod{AppleContainer: &applePayload},
			},
			want: true,
		},
		{
			name: "google pay token",
			req: Request{
				PaymentData:   &PaymentData{},
				PaymentMethod: &PaymentMethod{GoogleToken: &googlePayload},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.IsMobile(); got != tt.want {
				t.Errorf("IsMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}
