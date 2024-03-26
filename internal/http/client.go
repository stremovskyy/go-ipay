package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	go_ipay "github.com/megakit-pro/go-ipay/internal/consts"
	"github.com/megakit-pro/go-ipay/internal/log"
	"github.com/megakit-pro/go-ipay/ipay"
)

type Client struct {
	client  *http.Client
	options *Options
	logger  *log.Logger
}

func (c *Client) Api(apiRequest *ipay.RequestWrapper) (*ipay.Response, error) {
	requestID := uuid.New().String()

	c.logger.Debug("Request ID: %v", requestID)

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %v", err)
	}

	c.logger.Debug("Request: %v", string(jsonBody))

	req, err := http.NewRequest("POST", ipay.ApiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.logger.Error("cannot create request: %v", err)
		return nil, fmt.Errorf("cannot create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GO IPAY/"+go_ipay.Version)
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("Api-Version", ipay.ApiVersion)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("cannot send request: %v", err)
		return nil, fmt.Errorf("cannot send request: %v", err)
	}

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

	return ipay.UnmarshalCreateTokenResponse(raw)
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
		client:  cl,
		options: options,
		logger:  log.NewLogger("iPay HTTP:"),
	}
}