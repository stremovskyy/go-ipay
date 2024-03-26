package ipay

import "encoding/json"

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

func (r *RequestWrapper) AddTransaction(amount int, currency, description, info string) {
	if r.Request.Body.Transactions == nil {
		r.Request.Body.Transactions = make([]RequestTransaction, 0)
	}

	r.Request.Body.Transactions = append(
		r.Request.Body.Transactions, RequestTransaction{
			Amount: amount,
			Desc:   description,
			Info:   []string{info, "key"},
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
	MchID   int      `xml:"mch_id" json:"mch_id"`   // Merchant ID
	SrvID   int      `xml:"srv_id" json:"srv_id"`   // Legal entity for which the operation is carried out
	Invoice int      `xml:"invoice" json:"invoice"` // Payment amount in kopecks
	Amount  int      `xml:"amount" json:"amount"`   // Amount to be paid (including commission) in kopecks
	Desc    string   `xml:"desc" json:"desc"`       // Payment description
	Info    []string `xml:"info" json:"info"`       // Information for the payment provided by the merchant
}

// Card represents the card data.
type Card struct {
	Token *string `json:"token,omitempty"` // Card token.
}

// Info holds additional information related to the payment, provided by the merchant.
type Info struct {
	OrderId     *string `json:"order_id,omitempty"` // Order ID.
	ExtId       *string `json:"ext_id,omitempty"`   // External ID.
	UserID      *string `json:"user_id,omitempty"`  // User ID.
	Cvd         *Cvd    `json:"cvd,omitempty"`      // Card Verification Data.
	Aml         *Aml    `json:"aml,omitempty"`      // Anti-Money Laundering information.
	MctsVts     bool    `json:"mcts_vts"`           // If set, creates a token of type mcts/vts along with the default tokly token.
	ExternalCVD *Cvd    `json:"external_cvd"`       // External Card Verification Data.
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
