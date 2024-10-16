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
	"github.com/stremovskyy/go-ipay/currency"
)

func NewRequest(action Action, lang Lang, options ...func(*RequestWrapper)) *RequestWrapper {
	rw := &RequestWrapper{
		Request: Request{
			Action: action,
			Lang:   lang,
			Body:   Body{},
		},
	}

	for _, option := range options {
		option(rw)
	}

	return rw
}

func WithAmount(amountString int) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Amount: amountString,
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			rw.Request.Body.Transactions[i].Amount = amountString
		}
	}
}

func WithAmountInTransactions(amountString int, subMerchantId *int) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Amount: amountString,
					SmchId: subMerchantId,
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			rw.Request.Body.Transactions[i].Amount = amountString
			rw.Request.Body.Transactions[i].SmchId = subMerchantId
		}
	}
}

func WithInvoiceInTransactions(amountString int, subMerchantId *int) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Invoice: amountString,
					SmchId:  subMerchantId,
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			rw.Request.Body.Transactions[i].Invoice = amountString
			rw.Request.Body.Transactions[i].SmchId = subMerchantId
		}
	}
}

func WithInvoiceAmount(amount int) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.Invoice = &amount
	}
}

func WithCardToken(cardToken *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Card == nil {
			rw.Request.Body.Card = &Card{
				Token: cardToken,
			}

			return
		}

		rw.Request.Body.Card.Token = cardToken
	}
}

func WithPreauth(preauth bool) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		preauthInt := 0

		if preauth {
			preauthInt = 1
		}

		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Info: &Info{
						Preauth: &preauthInt,
					},
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			if rw.Request.Body.Transactions[i].Info == nil {
				rw.Request.Body.Transactions[i].Info = &Info{}
			}

			rw.Request.Body.Transactions[i].Info.Preauth = &preauthInt
		}
	}
}

func WithNotifyURL(url *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Info: &Info{
						NotifyUrl: url,
					},
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			if rw.Request.Body.Transactions[i].Info == nil {
				rw.Request.Body.Transactions[i].Info = &Info{}
			}

			rw.Request.Body.Transactions[i].Info.NotifyUrl = url
		}
	}
}

func WithRedirects(success string, fail string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.UrlGood = &success
		rw.Request.Body.UrlBad = &fail
	}
}

func WithIpayPaymentID(ipayPaymentID int64) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.PmtId = &ipayPaymentID
	}
}

func WithPaymentID(paymentID *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.ExtId = paymentID

		if rw.Request.Body.Info == nil {
			rw.Request.Body.Info = &Info{
				OrderId: paymentID,
				ExtId:   paymentID,
			}

			return
		}

		rw.Request.Body.Info.OrderId = paymentID
		rw.Request.Body.ExtId = paymentID

		if rw.Request.Body.Transactions != nil && len(rw.Request.Body.Transactions) != 0 {
			for i := range rw.Request.Body.Transactions {
				if rw.Request.Body.Transactions[i].Info == nil {
					rw.Request.Body.Transactions[i].Info = &Info{}
				}

				rw.Request.Body.Transactions[i].Info.OrderId = paymentID
			}

			rw.Request.Body.Info.OrderId = paymentID
			rw.Request.Body.Info.ExtId = paymentID

		}
	}
}

func WithAuth(auth Auth) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Auth = auth
	}
}

func WithPersonalData(data *Info) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.Info = data
	}
}

func WithWebhookURL(url *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Action == ActionDebiting {
			if rw.Request.Body.Transactions == nil {
				rw.Request.Body.Transactions = []RequestTransaction{
					{
						Info: &Info{
							NotifyUrl: url,
						},
					},
				}

				return
			}

			for i := range rw.Request.Body.Transactions {
				if rw.Request.Body.Transactions[i].Info == nil {
					rw.Request.Body.Transactions[i].Info = &Info{}
				}

				rw.Request.Body.Transactions[i].Info.NotifyUrl = url
			}

			return
		}

		if rw.Request.Body.Info == nil {
			rw.Request.Body.Info = &Info{
				NotifyUrl: url,
			}

			return
		}

		rw.Request.Body.Info.NotifyUrl = url
	}
}

func WithDescription(description string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Desc: description,
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			rw.Request.Body.Transactions[i].Desc = description
		}

	}
}

func WithCurrency(currency currency.Code) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Currency: currency,
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			rw.Request.Body.Transactions[i].Currency = currency
		}
	}
}

func WithOutAmount(withAmount bool) func(*RequestWrapper) {
	amountString := "no_amount"

	if !withAmount {
		amountString = "with_amount"
	}

	return func(rw *RequestWrapper) {
		rw.Request.Body.VerifyType = &amountString
	}
}

func WithAppleContainer(base64Object *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.AppleData = base64Object
	}
}

func WithGoogleContainer(token *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.Token = token
	}
}

func WithTrackingData(data *int64) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.Info.PmtIdIn = data
	}
}

func WithReceiverTIN(tin *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		var c = Cvd{
			TaxID: tin,
		}

		rw.Request.Body.Info.Cvd = c

		if rw.Request.Body.Aml == nil {
			rw.Request.Body.Aml = &Aml{}
		}

		rw.Request.Body.Aml.ReceiverIdentificationNumber = tin
	}
}

func WithTrackingToken(token *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.Info.ReceiverAccountNumber = token

		if rw.Request.Body.Aml == nil {
			rw.Request.Body.Aml = &Aml{}
		}

		rw.Request.Body.Aml.ReceiverAccountNumber = token
	}
}

func WithCardPan(pan *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Card == nil {
			rw.Request.Body.Card = &Card{
				Pan: pan,
			}

			return
		}

		rw.Request.Body.Card.Pan = pan
	}
}

func WithAML(aml *Aml) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.Aml = aml
	}
}

func WithReceiver(receiver *Receiver) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Info == nil {
			rw.Request.Body.Info = &Info{}
		}

		rw.Request.Body.Receiver = receiver
	}
}

func WithRelatedIDs(relatedIDs []int64) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		// rw.Request.Body.RelatedIDs = relatedIDs
	}
}
