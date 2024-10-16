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

package go_ipay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/stremovskyy/go-ipay/currency"
	"github.com/stremovskyy/go-ipay/internal/utils"
	"github.com/stremovskyy/go-ipay/ipay"
)

type Request struct {
	Merchant      *Merchant
	PersonalData  *PersonalData
	PaymentData   *PaymentData
	PaymentMethod *PaymentMethod
}

func (r *Request) GetAuth() ipay.Auth {
	if r.Merchant == nil {
		return ipay.Auth{
			Sign: "",
		}
	}

	sign := r.Merchant.GetSign()

	return ipay.Auth{
		MchID: r.Merchant.GetMerchantID(),
		Salt:  sign.Salt,
		Sign:  sign.Sign,
	}
}

func (r *Request) GetMobileAuth() ipay.Auth {
	if r.Merchant == nil {
		return ipay.Auth{
			Sign: "",
		}
	}

	sign := r.Merchant.GetMobileSign()

	return ipay.Auth{
		Login: r.Merchant.GetMobileLogin(),
		Sign:  sign.Sign,
		Time:  sign.Time,
	}
}

func (r *Request) GetRedirects() (string, string) {
	if r.Merchant == nil {
		return "", ""
	}

	return r.Merchant.SuccessRedirect, r.Merchant.FailRedirect
}

func (r *Request) GetPersonalData() *ipay.Info {
	if r.PersonalData == nil {
		return &ipay.Info{}
	}

	info := &ipay.Info{}

	if r.PersonalData.UserID != nil {
		info.UserID = utils.Ref(strconv.Itoa(*r.PersonalData.UserID))
	}

	info.Cvd = &ipay.Cvd{
		Firstname:  r.PersonalData.FirstName,
		Lastname:   r.PersonalData.LastName,
		Middlename: r.PersonalData.MiddleName,
		TaxID:      r.PersonalData.TaxID,
	}

	return info
}

func (r *Request) GetAML() *ipay.Aml {
	if r.PersonalData == nil {
		return &ipay.Aml{}
	}

	info := &ipay.Aml{
		ReceiverFirstname:            r.PersonalData.FirstName,
		ReceiverMiddlename:           r.PersonalData.MiddleName,
		ReceiverLastname:             r.PersonalData.LastName,
		ReceiverIdentificationNumber: r.PersonalData.TaxID,
		ReceiverToklyToken:           r.PersonalData.TrackingCardToken,
	}

	return info
}

func (r *Request) GetReceiver() *ipay.Receiver {
	if r.PersonalData == nil {
		return &ipay.Receiver{}
	}

	info := &ipay.Receiver{
		Lastname:             r.PersonalData.LastName,
		Firstname:            r.PersonalData.FirstName,
		Middlename:           r.PersonalData.MiddleName,
		IdentificationNumber: r.PersonalData.TaxID,
	}

	return info
}

func (r *Request) GetSender() *ipay.Sender {
	if r.PersonalData == nil {
		return &ipay.Sender{}
	}

	info := &ipay.Sender{
		Lastname:             r.PersonalData.LastName,
		Firstname:            r.PersonalData.FirstName,
		Middlename:           r.PersonalData.MiddleName,
		IdentificationNumber: r.PersonalData.TaxID,
	}

	return info
}

func (r *Request) GetIpayPaymentID() int64 {
	if r.PaymentData == nil || r.PaymentData.IpayPaymentID == nil {
		return 0
	}

	return *r.PaymentData.IpayPaymentID
}

func (r *Request) GetTransaction() (int, currency.Code, string) {
	if r.PaymentData == nil {
		return 0, "", ""
	}

	return r.PaymentData.Amount, r.PaymentData.Currency, r.PaymentData.Description
}

func (r *Request) GetCardToken() *string {
	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return nil
	}

	return r.PaymentMethod.Card.Token
}

func (r *Request) GetCardPan() *string {
	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return nil
	}

	return r.PaymentMethod.Card.Pan
}

func (r *Request) GetPaymentID() *string {
	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.PaymentID
}

func (r *Request) SetRedirects(successURL string, failURL string) {
	if r.Merchant == nil {
		r.Merchant = &Merchant{}
	}

	r.Merchant.SuccessRedirect = successURL
	r.Merchant.FailRedirect = failURL
}

func (r *Request) GetWebhookURL() *string {
	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.WebhookURL
}

func (r *Request) SetWebhookURL(webhookURL *string) {
	if r.PaymentData == nil {
		r.PaymentData = &PaymentData{}
	}

	r.PaymentData.WebhookURL = webhookURL
}

func (r *Request) GetAmount() int {
	if r.PaymentData == nil {
		return 0
	}

	return r.PaymentData.Amount

}

func (r *Request) GetDescription() string {
	if r.PaymentData == nil {
		return ""
	}

	return r.PaymentData.Description
}

func (r *Request) GetCurrency() currency.Code {
	if r.PaymentData == nil {
		return ""
	}

	return r.PaymentData.Currency

}

func (r *Request) GetSubMerchantID() *int {
	if r.Merchant == nil || r.Merchant.SubMerchantID == 0 {
		return nil
	}

	return &r.Merchant.SubMerchantID
}

func (r *Request) IsMobile() bool {
	if r.PaymentData == nil {
		return false
	}

	return r.PaymentData.IsMobile || r.PaymentMethod.AppleContainer != nil || r.PaymentMethod.GoogleToken != nil
}

func (r *Request) GetAppleContainer() (*string, error) {
	if r.PaymentMethod == nil || r.PaymentMethod.AppleContainer == nil {
		return nil, fmt.Errorf("Apple Container is not set")
	}

	decoded, err := base64.StdEncoding.DecodeString(*r.PaymentMethod.AppleContainer)
	if err != nil {
		return nil, fmt.Errorf("cannot decode Apple Container: %v", err)
	}

	var token map[string]interface{}
	if err := json.Unmarshal(decoded, &token); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}

	outputJSON, err := json.Marshal(token["token"])
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	outputBase64 := base64.StdEncoding.EncodeToString(outputJSON)
	return &outputBase64, nil
}

func (r *Request) IsApplePay() bool {
	return r.PaymentMethod != nil && r.PaymentMethod.AppleContainer != nil
}

func (r *Request) GetGoogleToken() (*string, error) {
	if r.PaymentMethod == nil || r.PaymentMethod.GoogleToken == nil {
		return nil, fmt.Errorf("Google Token is not set")
	}

	decoded, err := base64.StdEncoding.DecodeString(*r.PaymentMethod.GoogleToken)
	if err != nil {
		return nil, fmt.Errorf("cannot decode Google Token: %v", err)
	}

	var data struct {
		PaymentMethodData struct {
			TokenizationData struct {
				Token string `json:"token"`
			} `json:"tokenizationData"`
		} `json:"paymentMethodData"`
	}

	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}

	unescapedToken, err := strconv.Unquote(fmt.Sprintf("%q", data.PaymentMethodData.TokenizationData.Token))
	if err != nil {
		return nil, fmt.Errorf("unquote error: %v", err)
	}

	outputBase64 := base64.StdEncoding.EncodeToString([]byte(unescapedToken))
	return &outputBase64, nil
}

func (r *Request) GetTrackingData() *int64 {
	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.IpayPaymentID
}

func (r *Request) GetReceiverTIN() *string {
	if r.PersonalData == nil {
		return nil
	}

	return r.PersonalData.TaxID
}

func (r *Request) GetRelatedIDs() []int64 {
	if r.PaymentData == nil || r.PaymentData.RelatedIds == nil {
		return nil
	}

	return r.PaymentData.RelatedIds
}
