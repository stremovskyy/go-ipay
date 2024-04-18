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
	"fmt"
	"net/url"

	"github.com/stremovskyy/go-ipay/internal/http"
	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/go-ipay/log"
)

type client struct {
	ipayClient *http.Client
}

func (c *client) SetLogLevel(levelDebug log.Level) {
	log.SetLevel(log.Level(levelDebug))
}

func NewDefaultClient() Ipay {
	return &client{
		ipayClient: http.NewClient(http.DefaultOptions()),
	}
}
func NewClient(options ...Option) Ipay {
	c := &client{
		ipayClient: http.NewClient(http.DefaultOptions()),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *client) VerificationLink(request *Request) (*url.URL, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	createTokenRequest := ipay.NewRequest(
		ipay.ActionCreateToken3DS, ipay.LangUk,
		ipay.WithAuth(request.GetAuth()),
		ipay.WithRedirects(request.GetRedirects()),
		ipay.WithPersonalData(request.GetPersonalData()),
		ipay.WithPaymentID(request.GetPaymentID()),
		ipay.WithAmount(0),
		ipay.WithWebhookURL(request.GetWebhookURL()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithOutAmount(true),
	)

	apiResponse, err := c.ipayClient.Api(createTokenRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	u, err := url.Parse(apiResponse.Url)
	if err != nil {
		return nil, fmt.Errorf("cannot parse URL: %v", err)
	}

	return u, nil
}

func (c *client) Status(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	statusRequest := ipay.NewRequest(
		ipay.ActionGetPaymentStatus, ipay.LangUk,
		ipay.WithAuth(request.GetAuth()),
		ipay.WithIpayPaymentID(request.GetIpayPaymentID()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
	)

	apiResponse, err := c.ipayClient.Api(statusRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) PaymentURL(request *Request) (*ipay.PaymentResponse, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	XMLPaymentURLRequest := ipay.CreateXMLPaymentCreateRequest()
	XMLPaymentURLRequest.SetAuth(request.GetAuth())
	XMLPaymentURLRequest.SetRedirects(request.GetRedirects())
	XMLPaymentURLRequest.AddTransaction(request.GetTransaction())
	XMLPaymentURLRequest.SetPersonalData(request.GetPersonalData())
	XMLPaymentURLRequest.AddCardToken(request.GetCardToken())

	apiResponse, err := c.ipayClient.ApiXML(XMLPaymentURLRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Payment(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	paymentRequest := ipay.NewRequest(
		ipay.ActionDebiting, ipay.LangUk,
		ipay.WithPreauth(false),
		ipay.WithAmount(request.GetAmount()),
		ipay.WIthCurrency(request.GetCurrency()),
		ipay.WithAuth(request.GetAuth()),
		ipay.WithRedirects(request.GetRedirects()),
		ipay.WithPersonalData(request.GetPersonalData()),
		ipay.WithCardToken(request.GetCardToken()),
		ipay.WithPaymentID(request.GetPaymentID()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
	)

	apiResponse, err := c.ipayClient.Api(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Hold(request *Request) (*ipay.Response, error) {
	holdRequest := ipay.NewRequest(
		ipay.ActionDebiting, ipay.LangUk,
		ipay.WithPreauth(true),
		ipay.WithAmount(request.GetAmount()),
		ipay.WIthCurrency(request.GetCurrency()),
		ipay.WithAuth(request.GetAuth()),
		ipay.WithRedirects(request.GetRedirects()),
		ipay.WithPersonalData(request.GetPersonalData()),
		ipay.WithCardToken(request.GetCardToken()),
		ipay.WithPaymentID(request.GetPaymentID()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
	)

	apiResponse, err := c.ipayClient.Api(holdRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Capture(invoiceRequest *Request) (*ipay.Response, error) {
	captureRequest := ipay.NewRequest(
		ipay.ActionCompletion, ipay.LangUk,
		ipay.WithAuth(invoiceRequest.GetAuth()),
		ipay.WithIpayPaymentID(invoiceRequest.GetIpayPaymentID()),
		ipay.WithWebhookURL(invoiceRequest.GetWebhookURL()),
	)

	apiResponse, err := c.ipayClient.Api(captureRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Refund(invoiceRequest *Request) (*ipay.Response, error) {
	refundRequest := ipay.NewRequest(
		ipay.ActionReversal, ipay.LangUk,
		ipay.WithAuth(invoiceRequest.GetAuth()),
		ipay.WithIpayPaymentID(invoiceRequest.GetIpayPaymentID()),
		ipay.WithWebhookURL(invoiceRequest.GetWebhookURL()),
	)

	apiResponse, err := c.ipayClient.Api(refundRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Credit(invoiceRequest *Request) (*ipay.Response, error) {
	// TODO implement me
	panic("implement me")
}
