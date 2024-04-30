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

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/stremovskyy/go-ipay/consts"
	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/go-ipay/log"
)

type Client struct {
	client    *http.Client
	options   *Options
	logger    *log.Logger
	xmlLogger *log.Logger
}

func (c *Client) Api(apiRequest *ipay.RequestWrapper) (*ipay.Response, error) {
	requestID := uuid.New().String()

	c.logger.Debug("Request ID: %v", requestID)

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %v", err)
	}

	c.logger.Debug("Request: %v", string(jsonBody))

	req, err := http.NewRequest("POST", consts.ApiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.logger.Error("cannot create request: %v", err)
		return nil, fmt.Errorf("cannot create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GO IPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("Api-Version", consts.ApiVersion)

	tStart := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("cannot send request: %v", err)
		return nil, fmt.Errorf("cannot send request: %v", err)
	}
	c.logger.Debug("Request time: %v", time.Since(tStart))

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.logger.Error("cannot close response body: %v", err)
		}

	}()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("cannot read response: %v", err)
		return nil, fmt.Errorf("cannot read response: %v", err)
	}

	c.logger.Debug("Response: %v", string(raw))
	c.logger.Debug("Response status: %v", resp.StatusCode)

	response, err := ipay.UnmarshalJSONResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %v", err)
	}

	if response.GetError() != nil {
		return nil, fmt.Errorf("ipay error: %v", response.GetError())
	}

	return response, nil
}

func (c *Client) ApiXML(ipayXMLPayment *ipay.XmlPayment) (*ipay.PaymentResponse, error) {
	requestID := uuid.New().String()

	c.xmlLogger.Debug("Request ID: %v", requestID)

	xmlBody, err := ipayXMLPayment.Marshal()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %v", err)
	}

	c.xmlLogger.Debug("Request: %v", string(xmlBody))

	// Form-encode the XML data
	formData := url.Values{}
	formData.Set("data", string(xmlBody))

	req, err := http.NewRequest("POST", "https://tokly.ipay.ua/api302", strings.NewReader(formData.Encode()))
	if err != nil {
		c.xmlLogger.Error("cannot create request: %v", err)
		return nil, fmt.Errorf("cannot create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("User-Agent", "GO IPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("Api-Version", consts.ApiVersion)

	tStart := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		c.xmlLogger.Error("cannot send request: %v", err)
		return nil, fmt.Errorf("cannot send request: %v", err)
	}
	c.xmlLogger.Debug("Request time: %v", time.Since(tStart))

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.xmlLogger.Error("cannot close response body: %v", err)
		}
	}()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		c.xmlLogger.Error("cannot read response: %v", err)
		return nil, fmt.Errorf("cannot read response: %v", err)
	}

	c.xmlLogger.Debug("Response: %v", string(raw))
	c.xmlLogger.Debug("Response status: %v", resp.StatusCode)

	return ipay.UnmarshalXmlResponse(raw)
}

func (c *Client) SetClient(cl *http.Client) {
	c.client = cl
}

func NewClient(options *Options) *Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: options.KeepAlive,
	}

	tr := &http.Transport{
		MaxIdleConns:       options.MaxIdleConns,
		IdleConnTimeout:    options.IdleConnTimeout,
		DisableCompression: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
	}

	cl := &http.Client{
		Transport: tr,
		Timeout:   options.Timeout,
	}

	return &Client{
		client:    cl,
		options:   options,
		logger:    log.NewLogger("iPay HTTP:"),
		xmlLogger: log.NewLogger("iPay HTTP XML:"),
	}
}
