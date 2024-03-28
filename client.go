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
	"github.com/stremovskyy/go-ipay/internal/log"
	"github.com/stremovskyy/go-ipay/ipay"
)

type client struct {
	client *http.Client
}

func (c *client) SetLogLevel(levelDebug int) {
	log.SetLevel(log.Level(levelDebug))
}

func NewDefaultClient() Ipay {
	return &client{
		client: http.NewClient(http.DefaultOptions()),
	}
}

func (c *client) VerificationLink(request *Request) (*url.URL, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	createTokenRequest := ipay.CreateCreateToken3DSRequest(false)
	createTokenRequest.SetAuth(request.GetAuth())
	createTokenRequest.SetRedirects(request.GetRedirects())
	createTokenRequest.SetPersonalData(request.GetPersonalData())
	createTokenRequest.SetPaymentID(request.GetPaymentID())

	apiResponse, err := c.client.Api(createTokenRequest)
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

	statusRequest := ipay.CreateStatusRequest()
	statusRequest.SetAuth(request.GetAuth())
	statusRequest.SetIpayPaymentID(request.GetIpayPaymentID())

	apiResponse, err := c.client.Api(statusRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) PaymentURL(request *Request) (*ipay.PaymentResponse, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	paymentURLRequest := ipay.CreatePaymentCreateRequest()
	paymentURLRequest.SetAuth(request.GetAuth())
	paymentURLRequest.SetRedirects(request.GetRedirects())
	paymentURLRequest.AddTransaction(request.GetTransaction())
	paymentURLRequest.SetPersonalData(request.GetPersonalData())
	paymentURLRequest.AddCardToken(request.GetCardToken())

	apiResponse, err := c.client.ApiXML(paymentURLRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Payment(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	paymentRequest := ipay.CreatePaymentRequest()
	paymentRequest.SetAuth(request.GetAuth())
	paymentRequest.SetRedirects(request.GetRedirects())
	paymentRequest.AddTransaction(request.GetTransaction())
	paymentRequest.SetPersonalData(request.GetPersonalData())
	paymentRequest.AddCardToken(request.GetCardToken())
	paymentRequest.SetPaymentID(request.GetPaymentID())

	apiResponse, err := c.client.Api(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Hold(invoiceRequest *Request) (*ipay.Response, error) {
	holdRequest := ipay.CreateHoldRequest()
	holdRequest.SetAuth(invoiceRequest.GetAuth())
	holdRequest.SetIpayPaymentID(invoiceRequest.GetIpayPaymentID())

	apiResponse, err := c.client.Api(holdRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Capture(invoiceRequest *Request) (*ipay.Response, error) {
	captureRequest := ipay.CreateCaptureRequest()
	captureRequest.SetAuth(invoiceRequest.GetAuth())
	captureRequest.SetIpayPaymentID(invoiceRequest.GetIpayPaymentID())

	apiResponse, err := c.client.Api(captureRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Refund(invoiceRequest *Request) (*ipay.Response, error) {
	refundRequest := ipay.CreateRefundRequest()
	refundRequest.SetAuth(invoiceRequest.GetAuth())
	refundRequest.SetIpayPaymentID(invoiceRequest.GetIpayPaymentID())

	apiResponse, err := c.client.Api(refundRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Credit(invoiceRequest *Request) (*ipay.Response, error) {
	// TODO implement me
	panic("implement me")
}
