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

	"github.com/stremovskyy/go-ipay/consts"
	"github.com/stremovskyy/go-ipay/internal/http"
	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/go-ipay/log"
	"github.com/stremovskyy/recorder"
)

type client struct {
	ipayClient *http.Client
}

func (c *client) SetLogLevel(levelDebug log.Level) {
	log.SetLevel(levelDebug)
}

func NewDefaultClient() Ipay {
	return &client{
		ipayClient: http.NewClient(http.DefaultOptions()),
	}
}

func NewClientWithRecorder(rec recorder.Recorder) Ipay {
	return &client{
		ipayClient: http.NewClient(http.DefaultOptions()).WithRecorder(rec),
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
		ipay.WithInvoiceInTransactions(request.GetAmount(), request.GetSubMerchantID()),
		ipay.WithRedirects(request.GetRedirects()),
		ipay.WithPersonalData(request.GetPersonalData()),
		ipay.WithPaymentID(request.GetPaymentID()),
		ipay.WithAmount(0),
		ipay.WithWebhookURL(request.GetWebhookURL()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithOutAmount(true),
		ipay.WithAML(request.GetAML()),
		ipay.WithMetadata(request.GetMetadata()),
		ipay.WithOperationOperation(consts.VerificationLink),
	)

	apiResponse, err := c.ipayClient.Api(createTokenRequest)
	if err != nil {
		return nil, fmt.Errorf("verification link API call: %w", err)
	}

	if apiResponse == nil || apiResponse.Url == "" {
		return nil, fmt.Errorf("verification link: empty URL in API response")
	}

	u, err := url.Parse(apiResponse.Url)
	if err != nil {
		return nil, fmt.Errorf("verification link URL parsing: %w", err)
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
		ipay.WithOperationOperation(consts.Status),
	)

	return c.ipayClient.Api(statusRequest)
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
		return nil, fmt.Errorf("payment URL API call: %w", err)
	}

	return apiResponse, nil
}

func (c *client) Payment(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if request.IsMobile() {
		return c.handleMobilePayment(request, false)
	}

	return c.handleStandardPayment(request, false)
}

func (c *client) Hold(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if request.IsMobile() {
		return c.handleMobilePayment(request, true)
	}

	return c.handleStandardPayment(request, true)
}

func (c *client) handleMobilePayment(request *Request, isPreauth bool) (*ipay.Response, error) {
	var paymentRequest *ipay.RequestWrapper
	var apiFunc func(*ipay.RequestWrapper) (*ipay.Response, error)

	operationKind := consts.Payment

	if isPreauth {
		operationKind = consts.Hold
	}

	common := []func(*ipay.RequestWrapper){
		ipay.WithAuth(request.GetMobileAuth()),
		ipay.WithInvoiceAmount(request.GetAmount()),
		ipay.WithInvoiceInTransactions(request.GetAmount(), request.GetSubMerchantID()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithPaymentID(request.GetPaymentID()),
		ipay.WithPersonalData(request.GetPersonalData()),
		ipay.WithMetadata(request.GetMetadata()),
	}

	if isPreauth {
		common = append(common, ipay.WithPreauth(true))
	}

	if request.IsApplePay() {
		container, err := request.GetAppleContainer()
		if err != nil {
			return nil, fmt.Errorf("cannot get Apple Container: %w", err)
		}

		operationKind += consts.ApplePaySuffix

		common = append(common, ipay.WithAppleContainer(container))
		common = append(common, ipay.WithOperationOperation(operationKind))

		paymentRequest = ipay.NewRequest(
			ipay.MobilePaymentCreate, ipay.LangUk, common...,
		)
		apiFunc = c.ipayClient.ApplePayApi
	} else {
		token, err := request.GetGoogleToken()
		if err != nil {
			return nil, fmt.Errorf("cannot get Google Token: %w", err)
		}

		operationKind += consts.GooglePaySuffix

		common = append(common, ipay.WithGoogleContainer(token))
		common = append(common, ipay.WithOperationOperation(operationKind))

		paymentRequest = ipay.NewRequest(
			ipay.MobilePaymentCreate, ipay.LangUk, common...,
		)
		apiFunc = c.ipayClient.GooglePayApi
	}

	apiResponse, err := apiFunc(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("mobile payment API call: %w", err)
	}
	return apiResponse, nil
}

func (c *client) handleStandardPayment(request *Request, preauth bool) (*ipay.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("standard payment: %w", ErrRequestIsNil)
	}

	options := []func(*ipay.RequestWrapper){
		ipay.WithAmount(request.GetAmount()),
		ipay.WithCurrency(request.GetCurrency()),
		ipay.WithAuth(request.GetAuth()),
		ipay.WithPersonalData(request.GetPersonalData()),
		ipay.WithInvoiceInTransactions(request.GetAmount(), request.GetSubMerchantID()),
		ipay.WithCardToken(request.GetCardToken()),
		ipay.WithPaymentID(request.GetPaymentID()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
		ipay.WithMetadata(request.GetMetadata()),
		ipay.WithAML(request.GetAML()),
		ipay.WithOperationOperation(consts.Payment),
	}

	if preauth {
		options = append(options, ipay.WithPreauth(true))
	}

	holdRequest := ipay.NewRequest(ipay.ActionDebiting, ipay.LangUk, options...)

	return c.ipayClient.Api(holdRequest)
}

func (c *client) Capture(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("capture: %w", ErrRequestIsNil)
	}

	captureRequest := ipay.NewRequest(
		ipay.ActionCompletion, ipay.LangUk,
		ipay.WithAuth(request.GetAuth()),
		ipay.WithAmountInTransactions(request.GetAmount(), request.GetSubMerchantID()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithIpayPaymentID(request.GetIpayPaymentID()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
		ipay.WithRelatedIDs(request.GetRelatedIDs()),
		ipay.WithMetadata(request.GetMetadata()),
		ipay.WithAML(request.GetAML()),
		ipay.WithOperationOperation(consts.Capture),
	)

	return c.ipayClient.Api(captureRequest)
}

func (c *client) Refund(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("refund: %w", ErrRequestIsNil)
	}

	refundRequest := ipay.NewRequest(
		ipay.ActionReversal, ipay.LangUk,
		ipay.WithAuth(request.GetAuth()),
		ipay.WithIpayPaymentID(request.GetIpayPaymentID()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
		ipay.WithMetadata(request.GetMetadata()),
		ipay.WithOperationOperation(consts.Refund),
	)

	return c.ipayClient.Api(refundRequest)
}

func (c *client) Credit(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("credit: %w", ErrRequestIsNil)
	}

	options := []func(*ipay.RequestWrapper){
		ipay.WithAuth(request.GetAuth()),
		ipay.WithInvoiceAmount(request.GetAmount()),
		ipay.WithPaymentID(request.GetPaymentID()),
		ipay.WithWebhookURL(request.GetWebhookURL()),
		ipay.WithTrackingData(request.GetTrackingData()),
		ipay.WithDescription(request.GetDescription()),
		ipay.WithReceiver(request.GetReceiver()),
		ipay.WithMetadata(request.GetMetadata()),
		ipay.WithOperationOperation(consts.Credit),
		ipay.WithRelatedIDs(request.GetRelatedIDs()),
	}

	if request.GetCardToken() != nil {
		options = append(options, ipay.WithCardToken(request.GetCardToken()))
	} else if request.GetCardPan() != nil {
		options = append(options, ipay.WithCardPan(request.GetCardPan()))
	} else {
		return nil, fmt.Errorf("credit: neither CardToken nor CardPan provided")
	}

	creditRequest := ipay.NewRequest(
		ipay.ActionCredit, ipay.LangUk,
		options...,
	)

	response, err := c.ipayClient.Api(creditRequest)
	if err != nil {
		return nil, fmt.Errorf("credit API call: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("credit: empty response from API")
	}

	return response, nil
}

func (c *client) A2CPaymentStatus(request *Request) (*ipay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	extID := request.GetPaymentID()
	pmtID := request.GetIpayPaymentID()

	if (extID == nil || *extID == "") && pmtID == 0 {
		return nil, fmt.Errorf("A2CPaymentStatus: either ext_id or pmt_id must be provided")
	}
	if extID != nil && *extID != "" && pmtID != 0 {
		return nil, fmt.Errorf("A2CPaymentStatus: only one of ext_id or pmt_id must be provided")
	}

	opts := []func(*ipay.RequestWrapper){
		ipay.WithAuth(request.GetAuth()),
		ipay.WithOperationOperation(consts.A2CPaymentStatus),
	}

	if extID != nil && *extID != "" {
		opts = append(opts, ipay.WithExtID(extID))
	} else if pmtID != 0 {
		opts = append(opts, ipay.WithIpayPaymentID(pmtID))
	}

	statusRequest := ipay.NewRequest(
		ipay.ActionA2CPaymentStatus, ipay.LangUk,
		opts...,
	)

	return c.ipayClient.Api(statusRequest)
}
