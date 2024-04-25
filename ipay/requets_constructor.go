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

//
// func CreateCreateToken3DSRequest(withAmount bool) *RequestWrapper {
// 	amountString := "no_amount"
//
// 	if withAmount {
// 		amountString = "with_amount"
// 	}
//
// 	return &RequestWrapper{
// 		Request: Request{
// 			Auth:   Auth{},
// 			Action: ActionCreateToken3DS,
// 			Body: Body{
// 				VerifyType: &amountString,
// 			},
// 			Lang: LangUk,
// 		},
// 	}
// }
//
// func CreateCreateTokenRequest() *RequestWrapper {
// 	return &RequestWrapper{
// 		Request: Request{
// 			Action: ActionCreateToken,
// 			Lang:   LangUk,
// 		},
// 	}
// }
//
// func CreateStatusRequest() *RequestWrapper {
// 	return &RequestWrapper{
// 		Request: Request{
// 			Action: ActionGetPaymentStatus,
// 			Lang:   LangUk,
// 		},
// 	}
// }
// func CreatePaymentRequest() *RequestWrapper {
// 	return &RequestWrapper{
// 		Request: Request{
// 			Action: ActionDebiting,
// 			Lang:   LangUk,
// 			Body: Body{
// 				Info: &Info{
// 					Preauth: utils.Ref(0),
// 				},
// 			},
// 		},
// 	}
// }
// func CreateCaptureRequest() *RequestWrapper {
// 	return &RequestWrapper{
// 		Request: Request{
// 			Action: ActionCompletion,
// 			Lang:   LangUk,
// 		},
// 	}
// }
//
// func CreateHoldRequest() *RequestWrapper {
// 	return &RequestWrapper{
// 		Request: Request{
// 			Action: ActionDebiting,
// 			Lang:   LangUk,
// 			Body: Body{
// 				Info: &Info{
// 					Preauth: utils.Ref(1),
// 				},
// 			},
// 		},
// 	}
// }
//
// func CreateRefundRequest() *RequestWrapper {
// 	return &RequestWrapper{
// 		Request: Request{
// 			Action: ActionReversal,
// 			Lang:   LangUk,
// 		},
// 	}
// }

//
// func (r *RequestWrapper) SetPersonalData(info *Info) {
// 	if r.Request.Body.Info != nil {
// 		info.NotifyUrl = r.Request.Body.Info.NotifyUrl
// 		info.Preauth = r.Request.Body.Info.Preauth
// 	}
//
// 	r.Request.Body.Info = info
// }
//
// func (r *RequestWrapper) SetAuth(auth Auth) {
// 	r.Request.Auth = auth
// }
//
// func (r *RequestWrapper) SetRedirects(success string, fail string) {
// 	r.Request.Body.UrlGood = &success
// 	r.Request.Body.UrlBad = &fail
// }
//
// func (r *RequestWrapper) SetIpayPaymentID(ipayPaymentID int64) {
// 	r.Request.Body.PmtId = &ipayPaymentID
// }
//
// func (r *RequestWrapper) AddTransaction(amount int, currency currency.Code, description string) {
// 	if r.Request.Body.Transactions == nil {
// 		r.Request.Body.Transactions = make([]RequestTransaction, 0)
// 	}
//
// 	r.Request.Body.Transactions = append(
// 		r.Request.Body.Transactions, RequestTransaction{
// 			Amount:   amount,
// 			Currency: currency,
// 			Desc:     description,
// 			Info: Info{
// 				NotifyUrl: r.Request.Body.Info.NotifyUrl,
// 				Preauth:   r.Request.Body.Info.Preauth,
// 			},
// 		},
// 	)
//
// 	r.Request.Body.Info.Preauth = nil
// 	r.Request.Body.Info.NotifyUrl = nil
// }
//
// func (r *RequestWrapper) AddCardToken(cardToken *string) {
// 	if cardToken != nil {
// 		r.Request.Body.Card.Token = cardToken
// 	}
// }
//
// func (r *RequestWrapper) SetPaymentID(paymentID *string) {
// 	if paymentID != nil {
// 		r.Request.Body.Info.OrderId = paymentID
// 		r.Request.Body.ExtId = paymentID
// 	}
//
// 	if r.Request.Body.Transactions != nil && len(r.Request.Body.Transactions) != 0 {
// 		for i := range r.Request.Body.Transactions {
// 			r.Request.Body.Transactions[i].Info.OrderId = paymentID
// 		}
// 	}
//
// 	if r.Request.Body.Info == nil {
// 		r.Request.Body.Info = &Info{
// 			OrderId: paymentID,
// 			ExtId:   paymentID,
// 		}
// 	}
// }
//
// func (r *RequestWrapper) SetWebhookURL(url *string) {
// 	if url == nil {
// 		return
// 	}
//
// 	if r.Request.Body.Transactions != nil && len(r.Request.Body.Transactions) != 0 {
// 		for i := range r.Request.Body.Transactions {
// 			r.Request.Body.Transactions[i].Info.NotifyUrl = url
// 		}
// 	}
//
// 	if r.Request.Body.Info == nil {
// 		r.Request.Body.Info = &Info{
// 			NotifyUrl: url,
// 		}
// 	}
// }

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

func WithInvoiceAmount(amount int) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		rw.Request.Body.Invoice = &amount
	}
}

func WithCardToken(cardToken *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
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
					Info: Info{
						Preauth: &preauthInt,
					},
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
			rw.Request.Body.Transactions[i].Info.Preauth = &preauthInt
		}
	}
}

func WithNotifyURL(url *string) func(*RequestWrapper) {
	return func(rw *RequestWrapper) {
		if rw.Request.Body.Transactions == nil {
			rw.Request.Body.Transactions = []RequestTransaction{
				{
					Info: Info{
						NotifyUrl: url,
					},
				},
			}

			return
		}

		for i := range rw.Request.Body.Transactions {
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
						Info: Info{
							NotifyUrl: url,
						},
					},
				}

				return
			}

			for i := range rw.Request.Body.Transactions {
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

func WIthCurrency(currency currency.Code) func(*RequestWrapper) {
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
