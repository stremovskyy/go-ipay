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
)

type RequestWrapper struct {
	Request Request `json:"request"`
}

// Request represents the main structure of a payment request.
type Request struct {
	Auth   Auth   `json:"auth"`           // Authentication details for the payment request.
	Action Action `json:"action"`         // Specifies the action to be performed.
	Body   Body   `json:"body"`           // Contains the core data of the payment request.
	Lang   Lang   `json:"lang,omitempty"` // Optional language setting for web pages (ua - Ukrainian, en - English).
}

// Body encompasses the main content of the payment request.
type Body struct {
	Cdata        *string              `json:"cdata,omitempty"`        // Encoded card PAN.
	UrlGood      *string              `json:"url_good,omitempty"`     // Merchant's success URL.
	UrlBad       *string              `json:"url_bad,omitempty"`      // Merchant's failure URL.
	Info         *Info                `json:"info,omitempty"`         // Additional payment information.
	VerifyType   *string              `json:"verify_type,omitempty"`  // Verification type.
	PmtId        *int64               `json:"pmt_id,omitempty"`       // Payment ID.
	Transactions []RequestTransaction `json:"transactions,omitempty"` // List of transactions.
	Card         Card                 `json:"card,omitempty"`         // Card data.
	ExtId        *string              `json:"ext_id,omitempty"`       // External ID.
	Invoice      *int                 `json:"invoice,omitempty"`      // Payment amount in kopecks.
}

// Transaction represents an individual transaction.
type RequestTransaction struct {
	MchID    int           `xml:"mch_id" json:"mch_id,omitempty"`     // Merchant ID
	SrvID    int           `xml:"srv_id" json:"srv_id,omitempty"`     // Legal entity for which the operation is carried out
	Invoice  int           `xml:"invoice" json:"invoice,omitempty"`   // Payment amount in kopecks
	Amount   int           `xml:"amount" json:"amount,omitempty"`     // Amount to be paid (including commission) in kopecks
	Desc     string        `xml:"desc" json:"desc,omitempty"`         // Payment description
	Info     Info          `xml:"info" json:"info,omitempty"`         // Information for the payment provided by the merchant
	Currency currency.Code `xml:"currency" json:"currency,omitempty"` // Currency code
}

// Card represents the card data.
type Card struct {
	Token *string `json:"token,omitempty"` // Card token.
}

// Info holds additional information related to the payment, provided by the merchant.
type Info struct {
	OrderId     *string     `json:"order_id,omitempty"`     // Order ID.
	ExtId       *string     `json:"ext_id,omitempty"`       // External ID.
	UserID      *string     `json:"user_id,omitempty"`      // User ID.
	Cvd         interface{} `json:"cvd,omitempty"`          // Card Verification Data.
	Aml         *Aml        `json:"aml,omitempty"`          // Anti-Money Laundering information.
	MctsVts     bool        `json:"mcts_vts,omitempty"`     // If set, creates a token of type mcts/vts along with the default tokly token.
	ExternalCVD *Cvd        `json:"external_cvd,omitempty"` // External Card Verification Data.
	Preauth     *int        `json:"preauth,omitempty"`      // Preauthorization flag.
	NotifyUrl   *string     `json:"notify_url,omitempty"`
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
	SenderFirstname            *string `json:"sender_firstname,omitempty"`             // Sender's first name.
	SenderMiddlename           *string `json:"sender_middlename,omitempty"`            // Sender's middle name.
	SenderLastname             *string `json:"sender_lastname,omitempty"`              // Sender's last name.
	SenderIdentificationNumber *string `json:"sender_identification_number,omitempty"` // Sender's Identification Number.
	SenderDocument             *string `json:"sender_document,omitempty"`              // Sender's document number.
	SenderAddress              *string `json:"sender_address,omitempty"`               // Sender's address.
	SenderAccountNumber        *string `json:"sender_account_number,omitempty"`        // Sender's account number.
}
