package ipay

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Payment represents the root element of the notification with an ID.
type Payment struct {
	XMLName       xml.Name      `xml:"payment"`
	ID            int64         `xml:"id,attr"`         // Payment ID in the iPay system
	Ident         string        `xml:"ident"`           // Unique payment identifier
	Status        PaymentStatus `xml:"status"`          // Payment status
	Amount        float64       `xml:"amount"`          // Total payment amount
	Currency      string        `xml:"currency"`        // Currency code
	Timestamp     int64         `xml:"timestamp"`       // Date of authorization/completion in UNIX-timestamp
	CardToken     string        `xml:"card_token"`      // Card token
	CardIsPrepaid string        `xml:"card_is_prepaid"` // Whether the card is prepaid (1) or not (0), optional
	ValidTaxID    int           `xml:"valid_tax_id"`    // Valid (1) or not (0) tax ID sent in one of the requests: CreateToken, CreateToken3DS, PaymentCreate, optional
	CardHolder    string        `xml:"card_holder"`     // Full name of the cardholder, optional
	PaymentType   string        `xml:"payment_type"`    // Type of payment: Manual/GooglePay/ApplePay, optional
	Transactions  Transactions  `xml:"transactions"`    // Transactions element
	Salt          string        `xml:"salt"`            // Signature salt
	Sign          string        `xml:"sign"`            // Request signature
}

// Transactions represents a collection of Transaction.
type Transactions struct {
	Transaction []Transaction `xml:"transaction"` // Transaction element with transaction ID
}

// Transaction represents an individual transaction.
type Transaction struct {
	ID       int64  `xml:"id,attr"` // Transaction ID in the iPay system
	MchID    int    `xml:"mch_id"`  // Merchant ID
	SrvID    int    `xml:"srv_id"`  // Legal entity for which the operation is carried out
	Invoice  int    `xml:"invoice"` // Payment amount in kopecks
	Amount   int    `xml:"amount"`  // Amount to be paid (including commission) in kopecks
	Desc     string `xml:"desc"`    // Payment description
	Info     string `xml:"info"`    // Information for the payment provided by the merchant
	InfoData *Info  `xml:"-"`       // Parsed JSON object from transaction info
}

func ParsePaymentXML(data []byte) (*Payment, error) {
	var payment Payment
	err := xml.Unmarshal(data, &payment)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling Payment XML: %w", err)
	}

	// Parse JSON content in the "info" field of each transaction
	for i, transaction := range payment.Transactions.Transaction {
		var infoData Info
		err := json.Unmarshal([]byte(transaction.Info), &infoData)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling transaction info JSON: %w", err)
		}
		payment.Transactions.Transaction[i].InfoData = &infoData
	}

	return &payment, nil
}
