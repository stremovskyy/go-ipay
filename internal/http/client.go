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
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/stremovskyy/go-ipay/consts"
	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/go-ipay/log"
	"github.com/stremovskyy/recorder"
)

type loggerType string

const (
	loggerTypeHTTP      loggerType = "iPay HTTP:"
	loggerTypeHTTPXML   loggerType = "iPay HTTP XML:"
	loggerTypeApplePay  loggerType = "iPay ApplePay:"
	loggerTypeGooglePay loggerType = "iPay GooglePay:"
)

type Client struct {
	client    *http.Client
	options   *Options
	loggers   map[loggerType]*log.Logger
	loggersMu sync.RWMutex
	recorder  recorder.Recorder
}

// Api handles the standard iPay API request.
func (c *Client) Api(apiRequest *ipay.RequestWrapper) (*ipay.Response, error) {
	return c.sendRequest(consts.ApiUrl, apiRequest, c.loggerFor(loggerTypeHTTP))
}

// ApplePayApi handles the Apple Pay-specific API request.
func (c *Client) ApplePayApi(apiRequest *ipay.RequestWrapper) (*ipay.Response, error) {
	return c.sendRequest(consts.ApplePayUrl, apiRequest, c.loggerFor(loggerTypeApplePay))
}

// GooglePayApi handles the Google Pay-specific API request.
func (c *Client) GooglePayApi(apiRequest *ipay.RequestWrapper) (*ipay.Response, error) {
	return c.sendRequest(consts.GooglePayUrl, apiRequest, c.loggerFor(loggerTypeGooglePay))
}

func (c *Client) loggerFor(category loggerType) *log.Logger {
	c.loggersMu.RLock()
	logger, ok := c.loggers[category]
	c.loggersMu.RUnlock()
	if ok {
		return logger
	}

	c.loggersMu.Lock()
	defer c.loggersMu.Unlock()

	// Double-check to avoid races when creating loggers lazily.
	if logger, ok = c.loggers[category]; ok {
		return logger
	}

	if c.loggers == nil {
		c.loggers = make(map[loggerType]*log.Logger)
	}

	logger = log.NewLogger(string(category))
	c.loggers[category] = logger

	return logger
}

// WithRecorder attaches a recorder to the client.
func (c *Client) WithRecorder(rec recorder.Recorder) *Client {
	c.recorder = rec

	return c
}

// sendRequest handles sending an HTTP request and processing the response.
func (c *Client) sendRequest(apiURL string, apiRequest *ipay.RequestWrapper, logger *log.Logger) (*ipay.Response, error) {
	requestID := uuid.New().String()
	logger.Debug("Request ID: %v", requestID)

	if apiRequest == nil {
		return nil, c.logAndReturnError(
			"request is nil",
			fmt.Errorf("ipay request is nil"),
			logger,
			requestID,
			nil,
		)
	}

	action := apiRequest.Request.Action
	operation := apiRequest.Operation
	const method = http.MethodPost

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, c.logAndReturnError("cannot marshal request", err, logger, requestID, nil)
	}

	logger.Debug("Request: %v", string(jsonBody))

	ctx := context.WithValue(context.Background(), CtxKeyRequestID, requestID)
	tags := tagsRetriever(apiRequest)

	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, c.logAndReturnError("cannot create request", err, logger, requestID, tags)
	}

	c.setHeaders(req, requestID)

	if c.recorder != nil {
		if errr := c.recorder.RecordRequest(ctx, nil, requestID, jsonBody, tags); errr != nil {
			logger.Error("%s: cannot record request: %v", "error", errr)
		}
	}

	if c.client == nil {
		transportErr := &ipay.TransportError{
			Op:        "http.client",
			Action:    action,
			Operation: operation,
			Method:    method,
			Endpoint:  apiURL,
			RequestID: requestID,
			Cause:     fmt.Errorf("http client is nil"),
		}

		return nil, c.logAndReturnError(
			"cannot send request",
			transportErr,
			logger,
			requestID,
			errorTags(tags, "transport"),
		)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		transportErr := &ipay.TransportError{
			Op:             "http.do",
			Action:         action,
			Operation:      operation,
			Method:         method,
			Endpoint:       apiURL,
			RequestID:      requestID,
			OutcomeUnknown: true,
			Cause:          err,
		}

		return nil, c.logAndReturnError(
			"cannot send request",
			transportErr,
			logger,
			requestID,
			errorTags(tags, "transport"),
		)
	}
	if resp == nil {
		responseErr := &ipay.UnexpectedResponseError{
			Op:             "http.do",
			Action:         action,
			Operation:      operation,
			Method:         method,
			Endpoint:       apiURL,
			RequestID:      requestID,
			Message:        "response is nil",
			OutcomeUnknown: true,
		}

		return nil, c.logAndReturnError(
			"invalid response",
			responseErr,
			logger,
			requestID,
			errorTags(tags, "unexpected_response"),
		)
	}
	if resp.Body == nil {
		responseErr := &ipay.UnexpectedResponseError{
			Op:             "response.validate",
			Action:         action,
			Operation:      operation,
			Method:         method,
			Endpoint:       apiURL,
			RequestID:      requestID,
			StatusCode:     resp.StatusCode,
			ContentType:    responseContentType(resp),
			Message:        "response body is nil",
			OutcomeUnknown: true,
		}
		tagsWithResponse := responseTags(tags, resp)

		return nil, c.logAndReturnError(
			"invalid response",
			responseErr,
			logger,
			requestID,
			errorTags(tagsWithResponse, "unexpected_response"),
		)
	}
	defer c.safeClose(resp.Body, logger)

	tagsWithResponse := responseTags(tags, resp)
	raw, tooLarge, err := readBoundedResponseBody(resp.Body)
	if err != nil {
		if c.recorder != nil {
			if errr := c.recorder.RecordResponse(ctx, nil, requestID, raw, tagsWithResponse); errr != nil {
				logger.Error("%s: cannot record response %v", "error", errr)
			}
		}

		transportErr := &ipay.TransportError{
			Op:             "response.read",
			Action:         action,
			Operation:      operation,
			Method:         method,
			Endpoint:       apiURL,
			RequestID:      requestID,
			StatusCode:     resp.StatusCode,
			ContentType:    responseContentType(resp),
			Body:           diagnosticResponseBody(raw),
			OutcomeUnknown: true,
			Cause:          err,
		}

		return nil, c.logAndReturnError(
			"cannot read response",
			transportErr,
			logger,
			requestID,
			errorTags(responseTags(tags, resp), "transport"),
		)
	}

	logger.Debug("Response body: %d bytes", len(raw))
	logger.Debug("Response status: %v", resp.StatusCode)
	if tooLarge {
		tagsWithResponse["response_truncated"] = "true"
	}

	if c.recorder != nil {
		if errr := c.recorder.RecordResponse(ctx, nil, requestID, raw, tagsWithResponse); errr != nil {
			logger.Error("%s: cannot record response %v", "error", errr)
		}
	}

	responseContext := ipay.UnexpectedResponseError{
		Action:         action,
		Operation:      operation,
		Method:         method,
		Endpoint:       apiURL,
		RequestID:      requestID,
		StatusCode:     resp.StatusCode,
		ContentType:    responseContentType(resp),
		Body:           diagnosticResponseBody(raw),
		OutcomeUnknown: true,
	}

	if tooLarge {
		responseContext.Op = "response.read"
		responseContext.Message = fmt.Sprintf("response exceeds %d bytes", maxResponseBodyBytes)

		return nil, c.logAndReturnError(
			"response too large",
			&responseContext,
			logger,
			requestID,
			errorTags(tagsWithResponse, "unexpected_response"),
		)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		responseContext.Op = "response.status"
		responseContext.Message = "unexpected HTTP status"
		if isLikelyIPayJSONResponse(raw) {
			if parsed, parseErr := ipay.UnmarshalJSONResponse(raw); parseErr == nil && parsed != nil {
				if providerErr := parsed.GetError(); providerErr != nil {
					responseContext.Cause = providerErr
					responseContext.OutcomeUnknown = false
				}
			}
		}

		return nil, c.logAndReturnError(
			"unexpected response status",
			&responseContext,
			logger,
			requestID,
			errorTags(tagsWithResponse, "unexpected_response"),
		)
	}

	if len(bytes.TrimSpace(raw)) == 0 {
		responseContext.Op = "response.validate"
		responseContext.Message = "empty response body"

		return nil, c.logAndReturnError(
			"empty response",
			&responseContext,
			logger,
			requestID,
			errorTags(tagsWithResponse, "unexpected_response"),
		)
	}

	if !isLikelyIPayJSONResponse(raw) {
		responseContext.Op = "response.validate"
		responseContext.Message = "expected JSON response"

		return nil, c.logAndReturnError(
			"unexpected response format",
			&responseContext,
			logger,
			requestID,
			errorTags(tagsWithResponse, "unexpected_response"),
		)
	}

	response, err := ipay.UnmarshalJSONResponse(raw)
	if err != nil {
		decodeErr := &ipay.DecodeError{
			Op:             "response.decode",
			Action:         action,
			Operation:      operation,
			Method:         method,
			Endpoint:       apiURL,
			RequestID:      requestID,
			StatusCode:     resp.StatusCode,
			ContentType:    responseContentType(resp),
			Message:        "cannot decode JSON response",
			Body:           diagnosticResponseBody(raw),
			OutcomeUnknown: true,
			Cause:          err,
		}

		return nil, c.logAndReturnError(
			"cannot unmarshal response",
			decodeErr,
			logger,
			requestID,
			errorTags(tagsWithResponse, "decode"),
		)
	}

	return response, response.GetError()
}

// logAndReturnError logs an error and optionally records it.
func (c *Client) logAndReturnError(msg string, err error, logger *log.Logger, requestID string, tags map[string]string) error {
	// Logger is printf-style, not structured; include the error in the formatted message.
	logger.Error("%s: %v", msg, err)
	if c.recorder != nil {
		ctx := context.WithValue(context.Background(), CtxKeyRequestID, requestID)
		if recordErr := c.recorder.RecordError(ctx, nil, requestID, err, tags); recordErr != nil {
			logger.Error("%s: cannot record error %v", "error", recordErr)
		}
	}

	return err
}

// setHeaders sets common headers for all requests.
func (c *Client) setHeaders(req *http.Request, requestID string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GO IPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("Api-Version", consts.ApiVersion)
}

// safeClose ensures the body is closed properly and logs any error.
func (c *Client) safeClose(body io.ReadCloser, logger *log.Logger) {
	if err := body.Close(); err != nil {
		logger.Error("%s: cannot close response body, %v", "error", err)
	}
}

// tagsRetriever extracts tags from the request for logging or recording purposes.
func tagsRetriever(request *ipay.RequestWrapper) map[string]string {
	tags := make(map[string]string)

	if request.Request.Body.PmtId != nil {
		tags["payment_id"] = fmt.Sprintf("%v", *request.Request.Body.PmtId)
	}

	if request.Request.Body.ExtId != nil {
		tags["invoice_id"] = fmt.Sprintf("%v", *request.Request.Body.ExtId)
	}

	if request.Operation != "" {
		tags["operation"] = request.Operation
	}

	return tags
}

// ApiXML handles XML API requests.
func (c *Client) ApiXML(ipayXMLPayment *ipay.XmlPayment) (*ipay.PaymentResponse, error) {
	logger := c.loggerFor(loggerTypeHTTPXML)
	requestID := uuid.New().String()

	logger.Debug("Request ID: %v", requestID)

	xmlBody, err := ipayXMLPayment.Marshal()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	logger.Debug("Request: %v", string(xmlBody))

	formData := url.Values{}
	formData.Set("data", string(xmlBody))

	req, err := http.NewRequest("POST", consts.ApiXMLUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, c.logAndReturnError("cannot create XML request", err, logger, requestID, nil)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("User-Agent", "GO IPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("Api-Version", consts.ApiVersion)

	tStart := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.logAndReturnError("cannot send XML request", err, logger, requestID, nil)
	}
	logger.Debug("Request time: %v", time.Since(tStart))

	defer c.safeClose(resp.Body, logger)

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, c.logAndReturnError("cannot read XML response", err, logger, requestID, nil)
	}

	logger.Debug("Response: %v", string(raw))
	logger.Debug("Response status: %v", resp.StatusCode)

	return ipay.UnmarshalXmlResponse(raw)
}

// SetClient allows for replacing the default HTTP client.
func (c *Client) SetClient(cl *http.Client) {
	c.client = cl
}

// SetRecorder allows for attaching a new recorder.
func (c *Client) SetRecorder(r recorder.Recorder) {
	c.recorder = r
}

// NewClient initializes a new HTTP client with options.
func NewClient(options *Options) *Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: options.KeepAlive,
	}

	tr := &http.Transport{
		MaxIdleConns:       options.MaxIdleConns,
		IdleConnTimeout:    options.IdleConnTimeout,
		DisableCompression: true,
		DialContext:        dialer.DialContext,
	}

	cl := &http.Client{
		Transport: tr,
		Timeout:   options.Timeout,
	}

	return &Client{
		client:  cl,
		options: options,
		loggers: map[loggerType]*log.Logger{
			loggerTypeHTTP:      log.NewLogger(string(loggerTypeHTTP)),
			loggerTypeApplePay:  log.NewLogger(string(loggerTypeApplePay)),
			loggerTypeGooglePay: log.NewLogger(string(loggerTypeGooglePay)),
			loggerTypeHTTPXML:   log.NewLogger(string(loggerTypeHTTPXML)),
		},
	}
}
