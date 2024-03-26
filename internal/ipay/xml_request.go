package ipay

import (
	"encoding/xml"
	"fmt"
)

// XmlPayment represents the root of the payment request structure.
type XmlPayment struct {
	XMLName      xml.Name        `xml:"payment"`
	Auth         Auth            `xml:"auth"`
	Urls         *XmlUrls        `xml:"urls"`
	Card         *XmlCard        `xml:"card,omitempty"`
	Transactions XmlTransactions `xml:"transactions"`
	Lifetime     float64         `xml:"lifetime"`
	Lang         Lang            `xml:"lang"`
	Environment  string          `xml:"environment,omitempty"`
}

// XmlUrls contains the URLs for redirection after payment attempts.
type XmlUrls struct {
	Good string `xml:"good"`
	Bad  string `xml:"bad"`
}

// XmlCard contains the token or encoded PAN of the card.
type XmlCard struct {
	TokenType *string `xml:"token_type,omitempty"`
	Token     *string `xml:"token,omitempty"`
	Cdata     *string `xml:"cdata,omitempty"`
}

// XmlTransactions contains a slice of transactions.
type XmlTransactions struct {
	Transaction []XmlTransaction `xml:"transaction"`
}

// XmlTransaction represents a single transaction.
type XmlTransaction struct {
	Amount           int                  `xml:"amount"`
	Currency         string               `xml:"currency"`
	Desc             string               `xml:"desc"`
	Info             string               `xml:"info"` // This could be more complex depending on the structure of the info
	SmchID           *int                 `xml:"smch_id,omitempty"`
	AdditionalTokens *XmlAdditionalTokens `xml:"additional_tokens,omitempty"`
}

// XmlAdditionalTokens contains additional token information.
type XmlAdditionalTokens struct {
	MctsVts bool `xml:"mcts_vts"`
}

// Custom marshal function to ensure proper XML output
func (p XmlPayment) Marshal() ([]byte, error) {
	output, err := xml.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal XML: %w", err)
	}
	header := []byte(xml.Header)
	output = append(header, output...)
	return output, nil
}

// AddTransaction is a method to conveniently add a transaction to the XmlPayment.
func (p *XmlPayment) AddTransaction(amount int, currency, description, info string) {
	transaction := XmlTransaction{
		Amount:   amount,
		Currency: currency,
		Desc:     description,
		Info:     info,
	}
	p.Transactions.Transaction = append(p.Transactions.Transaction, transaction)
}

// AddCardToken adds card information to the XmlPayment.
func (p *XmlPayment) AddCardToken(token *string) {
	p.Card = &XmlCard{
		Token: token,
	}
}

func (p *XmlPayment) SetAuth(auth Auth) {
	p.Auth = auth
}

func (p *XmlPayment) SetRedirects(successUrl string, failUrl string) {
	if p.Urls == nil {
		p.Urls = &XmlUrls{}
	}

	p.Urls.Good = successUrl
	p.Urls.Bad = failUrl
}

func (p *XmlPayment) SetPersonalData(personalData *Info) {
	if personalData == nil {
		return
	}
	if p.Transactions.Transaction == nil {
		p.Transactions.Transaction = make([]XmlTransaction, 0)
	}

	for i := range p.Transactions.Transaction {
		p.Transactions.Transaction[i].Info = personalData.JsonString()
	}
}

func CreatePaymentCreateRequest() *XmlPayment {
	return &XmlPayment{
		Lang:     LangUk,
		Lifetime: 24,
	}
}
