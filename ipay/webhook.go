/*
 * MIT License
 *
 * Copyright (c) 2024 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package ipay

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Payment represents the root element of the notification with an ID.
type Payment struct {
	XMLName       xml.Name      `xml:"payment"`
	ID            int64         `xml:"id,attr"`                                // Payment ID in the iPay system
	Ident         string        `xml:"ident" json:"ident"`                     // Unique payment identifier
	Status        PaymentStatus `xml:"status" json:"status"`                   // Payment status
	Amount        float64       `xml:"amount" json:"amount"`                   // Total payment amount
	Currency      string        `xml:"currency" json:"currency"`               // Currency code
	Timestamp     int64         `xml:"timestamp" json:"timestamp"`             // Date of authorization/completion in UNIX-timestamp
	CardToken     *string       `xml:"card_token" json:"card_token"`           // Card token
	CardIsPrepaid string        `xml:"card_is_prepaid" json:"card_is_prepaid"` // Whether the card is prepaid (1) or not (0), optional
	ValidTaxID    int           `xml:"valid_tax_id" json:"valid_tax_id"`       // Valid (1) or not (0) tax ID sent in one of the requests: CreateToken, CreateToken3DS, PaymentCreate, optional
	CardHolder    string        `xml:"card_holder" json:"card_holder"`         // Full name of the cardholder, optional
	PaymentType   string        `xml:"payment_type" json:"payment_type"`       // Type of payment: Manual/GooglePay/ApplePay, optional
	Transactions  Transactions  `xml:"transactions" json:"transactions"`       // Transactions element
	Salt          string        `xml:"salt"`                                   // Signature salt
	Sign          string        `xml:"sign"`                                   // Request signature
	PmtId         int           `xml:"pmt_id" json:"pmt_id"`
	CardMask      *string       `xml:"card_mask" json:"card_mask"`
	Card          *string       `xml:"card" json:"card"`
	Invoice       int           `xml:"invoice" json:"invoice"`
	Desc          *string       `xml:"desc" json:"desc"`
	BnkErrorGroup interface{}   `xml:"bnk_error_group" json:"bnk_error_group"`
	BnkErrorNote  interface{}   `xml:"bnk_error_note" json:"bnk_error_note"`
	InitDate      string        `xml:"init_date" json:"init_date"`
	ExtID         *string       `xml:"ext_id" json:"ext_id"`

	// Additional fields for 3DS
	MchID          *string `xml:"mch_id" json:"mch_id"`
	Use3DS         *bool   `xml:"use_3ds" json:"use_3ds"`
	CardType       *string `xml:"card_type" json:"card_type"`
	BankName       *string `xml:"bank_name" json:"bank_name"`
	BnkError       *string `xml:"bnk_error" json:"bnk_error"`
	RRN            *string `xml:"rrn" json:"rrn"`
	RecurrentToken *string `xml:"recurrent_token" json:"recurrent_token"`
}

// Transactions represents a collection of Transaction.
type Transactions struct {
	Transaction []Transaction `xml:"transaction" json:"transaction"` // Transaction element with transaction ID
}

func (t *Transactions) Len() int {
	return len(t.Transaction)
}

func (t *Transactions) First() *Transaction {
	if t.Len() == 0 {
		return nil
	}

	return &t.Transaction[0]
}

func (t *Transactions) Last() *Transaction {
	if t.Len() == 0 {
		return nil
	}

	return &t.Transaction[t.Len()-1]
}

// Transaction represents an individual transaction.
type Transaction struct {
	ID       int64   `xml:"id,attr" json:"id"`          // Transaction ID in the iPay system
	MchID    int     `xml:"mch_id" json:"mch_id"`       // Merchant ID
	SrvID    int     `xml:"srv_id" json:"srv_id"`       // Legal entity for which the operation is carried out
	Invoice  int     `xml:"invoice" json:"invoice"`     // Payment amount in kopecks
	Amount   int     `xml:"amount" json:"amount"`       // Amount to be paid (including commission) in kopecks
	Desc     string  `xml:"desc" json:"desc"`           // Payment description
	Info     *string `xml:"info" json:"info,omitempty"` // Information for the payment provided by the merchant
	InfoData *Info   `xml:"-"`                          // Parsed JSON object from transaction info
}

func ParsePaymentXML(data []byte) (*Payment, error) {
	var payment Payment
	err := xml.Unmarshal(data, &payment)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling PaymentURL XML: %w", err)
	}

	// Parse JSON content in the "info" field of each transaction
	for i, transaction := range payment.Transactions.Transaction {
		if transaction.Info == nil {
			continue
		}

		var infoData Info
		err := json.Unmarshal([]byte(*transaction.Info), &infoData)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling transaction info JSON: %w", err)
		}
		payment.Transactions.Transaction[i].InfoData = &infoData
	}

	return &payment, nil
}

// String returns a string representation of the Payment struct.
func (p *Payment) String() string {
	return fmt.Sprintf(
		"Payment[ID=%d, Ident=%s, Status=%d, Amount=%.2f, Currency=%s, Timestamp=%d, Transactions=%d]",
		p.ID, p.Ident, p.Status, p.Amount, p.Currency, p.Timestamp, len(p.Transactions.Transaction),
	)
}

// IsValid returns true if the Payment struct contains valid data.
func (p *Payment) IsValid() bool {
	return p.ID > 0 && p.Amount >= 0 && p.Currency != "" && len(p.Transactions.Transaction) > 0
}

// GetTransactionByID returns a transaction by its ID.
func (p *Payment) GetTransactionByID(id int64) *Transaction {
	for _, transaction := range p.Transactions.Transaction {
		if transaction.ID == id {
			return &transaction
		}
	}
	return nil
}
