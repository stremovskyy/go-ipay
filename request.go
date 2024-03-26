package go_ipay

import (
	"strconv"

	"github.com/megakit-pro/go-ipay/ipay"
	"github.com/megakit-pro/go-ipay/utils"
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
			MchID: 0,
			Salt:  "EMPTY",
			Sign:  "",
		}
	}

	sign := r.Merchant.GetSign()

	return ipay.Auth{
		MchID: r.Merchant.GetMerchantID(),
		Salt:  sign.Salt,
		Sign:  sign.Sign,
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
		Firstname: r.PersonalData.FirstName,
		Lastname:  r.PersonalData.LastName,
		TaxID:     r.PersonalData.TaxID,
	}

	return info
}

func (r *Request) GetIpayPaymentID() int64 {
	if r.PaymentData == nil || r.PaymentData.IpayPaymentID == nil {
		return 0
	}

	return *r.PaymentData.IpayPaymentID
}

func (r *Request) GetTransaction() (int, string, string, string) {
	if r.PaymentData == nil {
		return 0, "", "", ""
	}

	return r.PaymentData.Amount, r.PaymentData.Currency, r.PaymentData.OrderID, r.PaymentData.Description
}

func (r *Request) GetCardToken() *string {
	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return nil
	}

	return r.PaymentMethod.Card.Token
}

func (r *Request) GetPaymentID() *string {
	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.PaymentID
}
