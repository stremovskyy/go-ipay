package go_ipay

import (
	"fmt"
	"net/url"

	"github.com/megakit-pro/go-ipay/internal/http"
	"github.com/megakit-pro/go-ipay/internal/ipay"
)

type client struct {
	client *http.Client
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
	// TODO implement me
	panic("implement me")
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
