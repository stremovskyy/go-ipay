package ipay

type Action string

// Action types for RequestWrapper
const (
	ActionCreateToken    Action = "CreateToken"
	ActionCreateToken3DS Action = "CreateToken3DS"
)

func CreateCreateToken3DSRequest(withAmount bool) RequestWrapper {
	amountString := "no_amount"

	if withAmount {
		amountString = "with_amount"
	}

	return RequestWrapper{
		Request: Request{
			Auth:   Auth{},
			Action: ActionCreateToken3DS,
			Body: Body{
				VerifyType: amountString,
			},
		},
	}
}

func CreateCreateTokenRequest() *RequestWrapper {
	return &RequestWrapper{
		Request: Request{
			Action: ActionCreateToken,
			Lang:   "uk",
		},
	}
}

type RequestWrapper struct {
	Request Request `json:"request"`
}

func (r *RequestWrapper) SetPersonalData(info Info) {
	r.Request.Body.Info = info
}

func (r *RequestWrapper) SetAuth(auth Auth) {
	r.Request.Auth = auth
}

func (r *RequestWrapper) SetRedirects(success string, fail string) {
	r.Request.Body.UrlGood = success
	r.Request.Body.UrlBad = fail
}

type Request struct {
	Auth   Auth   `json:"auth"`
	Action Action `json:"action"`
	Body   Body   `json:"body"`
	Lang   string `json:"lang,omitempty"`
}

type Auth struct {
	MchID int64  `json:"mch_id"`
	Salt  string `json:"salt"`
	Sign  string `json:"sign"`
}

type Body struct {
	Cdata      string `json:"cdata,omitempty"`
	UrlGood    string `json:"url_good,omitempty"`
	UrlBad     string `json:"url_bad,omitempty"`
	Info       Info   `json:"info,omitempty"`
	VerifyType string `json:"verify_type,omitempty"`
}

type Info struct {
	UserID  string `json:"user_id,omitempty"`
	Cvd     *Cvd   `json:"cvd,omitempty"`
	Aml     *Aml   `json:"aml,omitempty"`
	MctsVts bool   `json:"mcts_vts,omitempty"`
}

type Cvd struct {
	TaxID       string `json:"tax_id,omitempty"`
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	Middlename  string `json:"middlename,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type Aml struct {
	SenderFirstname            string `json:"sender_firstname,omitempty"`
	SenderMiddlename           string `json:"sender_middlename,omitempty"`
	SenderLastname             string `json:"sender_lastname,omitempty"`
	SenderIdentificationNumber string `json:"sender_identification_number,omitempty"`
	SenderDocument             string `json:"sender_document,omitempty"`
	SenderAddress              string `json:"sender_address,omitempty"`
	SenderAccountNumber        string `json:"sender_account_number,omitempty"`
}
