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
)

func CreateCreateToken3DSRequest(withAmount bool) *RequestWrapper {
	amountString := "no_amount"

	if withAmount {
		amountString = "with_amount"
	}

	return &RequestWrapper{
		Request: Request{
			Auth:   Auth{},
			Action: ActionCreateToken3DS,
			Body: Body{
				VerifyType: &amountString,
			},
			Lang: LangUk,
		},
	}
}

func CreateCreateTokenRequest() *RequestWrapper {
	return &RequestWrapper{
		Request: Request{
			Action: ActionCreateToken,
			Lang:   LangUk,
		},
	}
}

func CreateStatusRequest() *RequestWrapper {
	return &RequestWrapper{
		Request: Request{
			Action: ActionGetPaymentStatus,
			Lang:   LangUk,
		},
	}
}
func CreatePaymentRequest() *RequestWrapper {
	return &RequestWrapper{
		Request: Request{
			Action: ActionDebiting,
			Lang:   LangUk,
			Body: Body{
				Info: &Info{
					Preauth: 0,
				},
			},
		},
	}
}
func CreateCaptureRequest() *RequestWrapper {
	return &RequestWrapper{
		Request: Request{
			Action: ActionCompletion,
			Lang:   LangUk,
		},
	}
}

func CreateHoldRequest() *RequestWrapper {
	return &RequestWrapper{
		Request: Request{
			Action: ActionDebiting,
			Lang:   LangUk,
			Body: Body{
				Info: &Info{
					Preauth: 1,
				},
			},
		},
	}
}

func CreateRefundRequest() *RequestWrapper {
	return &RequestWrapper{
		Request: Request{
			Action: ActionReversal,
			Lang:   LangUk,
		},
	}
}

type RequestWrapper struct {
	Request Request `json:"request"`
}

func (r *RequestWrapper) SetPersonalData(info *Info) {
	if r.Request.Body.Info != nil {
		info.NotifyUrl = r.Request.Body.Info.NotifyUrl
		info.Preauth = r.Request.Body.Info.Preauth
	}

	r.Request.Body.Info = info
}

func (r *RequestWrapper) SetAuth(auth Auth) {
	r.Request.Auth = auth
}

func (r *RequestWrapper) SetRedirects(success string, fail string) {
	r.Request.Body.UrlGood = &success
	r.Request.Body.UrlBad = &fail
}

func (r *RequestWrapper) SetIpayPaymentID(ipayPaymentID int64) {
	r.Request.Body.PmtId = &ipayPaymentID
}

func (r *RequestWrapper) AddTransaction(amount int, currency currency.Code, description string) {
	if r.Request.Body.Transactions == nil {
		r.Request.Body.Transactions = make([]RequestTransaction, 0)
	}

	r.Request.Body.Transactions = append(
		r.Request.Body.Transactions, RequestTransaction{
			Amount:   amount,
			Currency: currency,
			Desc:     description,
			Info: Info{
				NotifyUrl: r.Request.Body.Info.NotifyUrl,
				Preauth:   r.Request.Body.Info.Preauth,
			},
		},
	)
}

func (r *RequestWrapper) AddCardToken(cardToken *string) {
	if cardToken != nil {
		r.Request.Body.Card.Token = cardToken
	}
}

func (r *RequestWrapper) SetPaymentID(paymentID *string) {
	if paymentID != nil {
		r.Request.Body.Info.OrderId = paymentID
		r.Request.Body.ExtId = paymentID
	}

	if r.Request.Body.Transactions != nil && len(r.Request.Body.Transactions) != 0 {
		for i := range r.Request.Body.Transactions {
			r.Request.Body.Transactions[i].Info.OrderId = paymentID
		}
	}

	if r.Request.Body.Info == nil {
		r.Request.Body.Info = &Info{
			OrderId: paymentID,
			ExtId:   paymentID,
		}
	}
}

func (r *RequestWrapper) SetWebhookURL(url *string) {
	if url == nil {
		return
	}

	if r.Request.Body.Transactions != nil && len(r.Request.Body.Transactions) != 0 {
		for i := range r.Request.Body.Transactions {
			r.Request.Body.Transactions[i].Info.NotifyUrl = url
		}
	}

	if r.Request.Body.Info == nil {
		r.Request.Body.Info = &Info{
			NotifyUrl: url,
		}
	}
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
	Preauth     int8        `json:"preauth"`                // Preauthorization flag.
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
