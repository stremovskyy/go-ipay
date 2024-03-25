package go_ipay

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/megakit-pro/go-ipay/internal/http"
	"github.com/megakit-pro/go-ipay/ipay"
)

type Client struct {
	client *http.Client
}

func NewDefaultClient() *Client {
	return &Client{
		client: http.NewClient(http.DefaultOptions()),
	}
}

func (c *Client) VerificationLink(request *Request) (*url.URL, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	createTokenRequest := ipay.CreateCreateTokenRequest()
	createTokenRequest.SetAuth(request.GetAuth())
	createTokenRequest.SetRedirects(request.GetRedirects())
	createTokenRequest.SetPersonalData(request.GetPersonalData())
	createTokenRequest.Request.Body.Info.OrderId = uuid.New().String()

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

func (c *Client) Status(request *Request) (*ipay.Response, error) {
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
