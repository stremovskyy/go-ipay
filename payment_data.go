package go_ipay

type PaymentData struct {
	IpayPaymentID *int64
	PaymentID     *string
	Amount        int
	Currency      string
	OrderID       string
	Description   string
}
