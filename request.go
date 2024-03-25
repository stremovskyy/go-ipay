package go_ipay

import (
	"strconv"

	"github.com/megakit-pro/go-ipay/ipay"
)

type Request struct {
	Merchant     *Merchant
	PersonalData *PersonalData
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

func (r *Request) GetPersonalData() ipay.Info {
	if r.PersonalData == nil {
		return ipay.Info{}
	}

	info := ipay.Info{}

	if r.PersonalData.UserID != nil {
		info.UserID = strconv.Itoa(*r.PersonalData.UserID)
	}

	info.Cvd = &ipay.Cvd{}

	if r.PersonalData.FirstName != nil {
		info.Cvd.Firstname = *r.PersonalData.FirstName
	}

	if r.PersonalData.LastName != nil {
		info.Cvd.Lastname = *r.PersonalData.LastName
	}

	if r.PersonalData.TaxID != nil {
		info.Cvd.TaxID = *r.PersonalData.TaxID
	}

	return info
}
