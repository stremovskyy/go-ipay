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

	"github.com/stremovskyy/go-ipay/currency"
)

type Action string

// Action types for RequestWrapper
const (
	ActionCreateToken      Action = "CreateToken"
	ActionCreateToken3DS   Action = "CreateToken3DS"
	ActionGetPaymentStatus Action = "GetPaymentStatus"
	ActionDebiting         Action = "Debiting"
	ActionCompletion       Action = "Completion"
	ActionReversal         Action = "Reversal"
	ActionCredit           Action = "A2CPay"
	MobilePaymentCreate    Action = "PaymentCreate"
	ActionA2CPaymentStatus Action = "A2CPaymenStatus"
)

type RequestWrapper struct {
	Request Request `json:"request"`

	Operation string `json:"-"`
}

// Request represents the main structure of a payment request.
type Request struct {
	Auth   Auth   `json:"auth"`           // Authentication details for the payment request.
	Action Action `json:"action"`         // Specifies the action to be performed.
	Body   Body   `json:"body"`           // Contains the core data of the payment request.
	Lang   *Lang  `json:"lang,omitempty"` // Optional language setting for web pages (ua - Ukrainian, en - English).
}

// Body encompasses the main content of the payment request.
type Body struct {
	Cdata          *string             `json:"cdata,omitempty"`           // Encoded card PAN.
	UrlGood        *string             `json:"url_good,omitempty"`        // Merchant's success URL.
	UrlBad         *string             `json:"url_bad,omitempty"`         // Merchant's failure URL.
	Info           *Info               `json:"info,omitempty"`            // Additional payment information.
	VerifyType     *string             `json:"verify_type,omitempty"`     // Verification type.
	PmtId          *int64              `json:"pmt_id,omitempty"`          // Payment ID.
	Transactions   RequestTransactions `json:"transactions,omitempty"`    // List of transactions.
	Card           *Card               `json:"card,omitempty"`            // Card data.
	ExtId          *string             `json:"ext_id,omitempty"`          // External ID.
	Invoice        *int                `json:"invoice,omitempty"`         // Payment amount in kopecks.
	AppleData      *string             `json:"apple_data,omitempty"`      // Apple Pay data.
	PmtDesc        *string             `json:"pmt_desc,omitempty"`        // Payment description.
	Token          *string             `json:"token,omitempty"`           // Token for Google Pay.
	Recurrent      *string             `json:"recurrent,omitempty"`       // Recurrent payment (true or false).
	RecurrentToken *string             `json:"recurrent_token,omitempty"` // Recurrent payment token.
	Aml            *Aml                `json:"aml,omitempty"`             // Anti-Money Laundering information.
	Sender         *Sender             `json:"sender,omitempty"`          // Sender details.
	Receiver       *Receiver           `json:"receiver,omitempty"`
}

type RequestTransactions []RequestTransaction

func (r RequestTransactions) Len() int {
	return len(r)
}

func (r RequestTransactions) First() *RequestTransaction {
	if r.Len() == 0 {
		return nil
	}

	return &r[0]
}

func (r RequestTransactions) Last() *RequestTransaction {
	if r.Len() == 0 {
		return nil
	}

	return &r[r.Len()-1]
}

// Transaction represents an individual transaction.
type RequestTransaction struct {
	MchID    int           `xml:"mch_id" json:"mch_id,omitempty"`     // Merchant ID
	SrvID    int           `xml:"srv_id" json:"srv_id,omitempty"`     // Legal entity for which the operation is carried out
	SmchId   *int          `xml:"smch_id" json:"smch_id,omitempty"`   // Submerchant ID
	Invoice  int           `xml:"invoice" json:"invoice,omitempty"`   // Payment amount in kopecks
	Amount   int           `xml:"amount" json:"amount,omitempty"`     // Amount to be paid (including commission) in kopecks
	Desc     string        `xml:"desc" json:"desc,omitempty"`         // Payment description
	Info     *Info         `xml:"info" json:"info"`                   // Information for the payment provided by the merchant
	Currency currency.Code `xml:"currency" json:"currency,omitempty"` // Currency code
}

// Card represents the card data.
type Card struct {
	Token     *string `json:"token,omitempty"`      // Card token.
	Pan       *string `json:"pan,omitempty"`        // Card PAN.
	TokenType *string `json:"token_type,omitempty"` // Token type.
}

// Info holds additional information related to the payment, provided by the merchant.
type Info struct {
	OrderId               *string     `json:"order_id,omitempty"`                // Order ID.
	ExtId                 *string     `json:"ext_id,omitempty"`                  // External ID.
	UserID                *string     `json:"user_id,omitempty"`                 // User ID.
	Cvd                   interface{} `json:"cvd,omitempty"`                     // Card Verification Data.
	Aml                   *Aml        `json:"aml,omitempty"`                     // Anti-Money Laundering information.
	MctsVts               bool        `json:"mcts_vts,omitempty"`                // If set, creates a token of type mcts/vts along with the default tokly token.
	ExternalCVD           *Cvd        `json:"external_cvd,omitempty"`            // External Card Verification Data.
	Preauth               *int        `json:"preauth,omitempty"`                 // Preauthorization flag.
	NotifyUrl             *string     `json:"notify_url,omitempty"`              // Notification URL.
	PmtIdIn               []int64     `json:"pmt_id_in,omitempty"`               // Payment IDs in.
	ReceiverAccountNumber *string     `json:"receiver_account_number,omitempty"` // Receiver's account number.
	Metadata              *string     `json:"metadata,omitempty"`                // Metadata.
}

func (i *Info) JsonString() string {
	jsonString, _ := json.Marshal(i)
	return string(jsonString)
}

// Cvd represents Card Verification Data.
type Cvd struct {
	TaxID       *string `json:"tax_id,omitempty"`       // Tax Identification Number.
	Firstname   *string `json:"firstname,omitempty"`    // First name.
	Lastname    *string `json:"lastname,omitempty"`     // Last name.
	Middlename  *string `json:"middlename,omitempty"`   // Middle name.
	PhoneNumber *string `json:"phone_number,omitempty"` // Phone number.
}

// Aml contains Anti-Money Laundering information for financial monitoring.
type Aml struct {
	SenderFirstname              *string `json:"sender_firstname,omitempty"`               // Sender's first name.
	SenderMiddlename             *string `json:"sender_middlename,omitempty"`              // Sender's middle name.
	SenderLastname               *string `json:"sender_lastname,omitempty"`                // Sender's last name.
	SenderIdentificationNumber   *string `json:"sender_identification_number,omitempty"`   // Sender's Identification Number.
	SenderDocument               *string `json:"sender_document,omitempty"`                // Sender's document number.
	SenderAddress                *string `json:"sender_address,omitempty"`                 // Sender's address.
	SenderAccountNumber          *string `json:"sender_account_number,omitempty"`          // Sender's account number.
	ReceiverFirstname            *string `json:"receiver_firstname,omitempty"`             // Receiver's first name.
	ReceiverMiddlename           *string `json:"receiver_middlename,omitempty"`            // Receiver's middle name.
	ReceiverLastname             *string `json:"receiver_lastname,omitempty"`              // Receiver's last name.
	ReceiverIdentificationNumber *string `json:"receiver_identification_number,omitempty"` // Receiver's Identification Number.
	ReceiverDocument             *string `json:"receiver_document,omitempty"`              // Receiver's document number.
	ReceiverAddress              *string `json:"receiver_address,omitempty"`               // Receiver's address.
	ReceiverAccountNumber        *string `json:"receiver_account_number,omitempty"`        // Receiver's account number.
	ReceiverToklyToken           *string `json:"receiver_tokly_token,omitempty"`           // Tokly token associated with the receiver's card number.
}

// Sender represents the sender's details.
type Sender struct {
	Lastname             *string `json:"lastname,omitempty"`             // Sender's last name.
	Firstname            *string `json:"firstname,omitempty"`            // Sender's first name.
	Middlename           *string `json:"middlename,omitempty"`           // Sender's middle name.
	Document             *string `json:"document,omitempty"`             // Sender's document number.
	Address              *string `json:"address,omitempty"`              // Sender's address.
	IdentificationNumber *string `json:"identificationNumber,omitempty"` // Sender's tax identification number (optional).
	AccountNumber        *string `json:"accountNumber,omitempty"`        // Sender's account number (optional).
}

// Receiver represents the receiver's details.
type Receiver struct {
	Lastname             *string `json:"lastname,omitempty"`             // Receiver's last name.
	Firstname            *string `json:"firstname,omitempty"`            // Receiver's first name.
	Middlename           *string `json:"middlename,omitempty"`           // Receiver's middle name.
	Document             *string `json:"document,omitempty"`             // Receiver's document number (optional).
	Address              *string `json:"address,omitempty"`              // Receiver's address.
	IdentificationNumber *string `json:"identificationNumber,omitempty"` // Receiver's tax identification number (optional).
	AccountNumber        *string `json:"accountNumber,omitempty"`        // Receiver's account number (optional).
}
