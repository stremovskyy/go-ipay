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
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/google/uuid"

	"github.com/stremovskyy/go-ipay/consts"
	"github.com/stremovskyy/go-ipay/log"
	"github.com/stremovskyy/go-ipay/repayment"
)

const loggerTypeRepayment loggerType = "iPay Repayment:"

// RepaymentJSONApi sends a Repayment API request as JSON (application/json).
func (c *Client) RepaymentJSONApi(apiRequest *repayment.RequestWrapper) (*repayment.Response, error) {
	return c.sendRepaymentJSONRequest(consts.RepaymentUrl, apiRequest, c.loggerFor(loggerTypeRepayment))
}

// RepaymentProcessingFileApi sends a Repayment API request and returns raw bytes (typically a CSV file).
// The API may respond with JSON errors, so callers should treat a non-nil error as authoritative even
// when raw bytes are returned.
func (c *Client) RepaymentProcessingFileApi(apiRequest *repayment.RequestWrapper) ([]byte, error) {
	return c.sendRepaymentProcessingFileRequest(consts.RepaymentUrl, apiRequest, c.loggerFor(loggerTypeRepayment))
}

// RepaymentApi sends a Repayment API request with a CSV file in multipart/form-data.
func (c *Client) RepaymentApi(apiRequest *repayment.RequestWrapper, fileName string, file io.Reader) (*repayment.Response, error) {
	return c.sendRepaymentMultipartRequest(consts.RepaymentUrl, apiRequest, fileName, file, c.loggerFor(loggerTypeRepayment))
}

func (c *Client) sendRepaymentJSONRequest(apiURL string, apiRequest *repayment.RequestWrapper, logger *log.Logger) (*repayment.Response, error) {
	requestID := uuid.New().String()
	logger.Debug("Request ID: %v", requestID)

	if apiRequest == nil {
		return nil, c.logAndReturnError("repayment request is nil", fmt.Errorf("request is nil"), logger, requestID, nil)
	}

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, c.logAndReturnError("cannot marshal repayment request", err, logger, requestID, tagsRetrieverRepayment(apiRequest))
	}

	logger.Debug("Request: %v", string(jsonBody))

	ctx := context.WithValue(context.Background(), CtxKeyRequestID, requestID)
	tags := tagsRetrieverRepayment(apiRequest)

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, c.logAndReturnError("cannot create repayment request", err, logger, requestID, tags)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GO IPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)

	if c.recorder != nil {
		if errr := c.recorder.RecordRequest(ctx, nil, requestID, jsonBody, tags); errr != nil {
			logger.Error("%s: cannot record request: %v", "error", errr)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.logAndReturnError("cannot send repayment request", err, logger, requestID, tags)
	}
	defer c.safeClose(resp.Body, logger)

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, c.logAndReturnError("cannot read repayment response", err, logger, requestID, tags)
	}

	logger.Debug("Response: %v", string(raw))
	logger.Debug("Response status: %v", resp.StatusCode)

	if c.recorder != nil {
		if errr := c.recorder.RecordResponse(ctx, nil, requestID, raw, tags); errr != nil {
			logger.Error("%s: cannot record response %v", "error", errr)
		}
	}

	if !isLikelyJSONResponse(resp, raw) {
		apiErr := nonJSONRepaymentAPIError(resp, raw)
		return nil, c.logAndReturnError("repayment API returned non-JSON response", apiErr, logger, requestID, tags)
	}

	response, err := repayment.UnmarshalJSONResponse(raw)
	if err != nil {
		return nil, c.logAndReturnError("cannot unmarshal repayment response", err, logger, requestID, tags)
	}

	return response, response.GetError()
}

func (c *Client) sendRepaymentProcessingFileRequest(apiURL string, apiRequest *repayment.RequestWrapper, logger *log.Logger) ([]byte, error) {
	requestID := uuid.New().String()
	logger.Debug("Request ID: %v", requestID)

	if apiRequest == nil {
		return nil, c.logAndReturnError("repayment request is nil", fmt.Errorf("request is nil"), logger, requestID, nil)
	}

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, c.logAndReturnError("cannot marshal repayment request", err, logger, requestID, tagsRetrieverRepayment(apiRequest))
	}

	logger.Debug("Request: %v", string(jsonBody))

	ctx := context.WithValue(context.Background(), CtxKeyRequestID, requestID)
	tags := tagsRetrieverRepayment(apiRequest)

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, c.logAndReturnError("cannot create repayment request", err, logger, requestID, tags)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/csv, application/json")
	req.Header.Set("User-Agent", "GO IPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)

	if c.recorder != nil {
		if errr := c.recorder.RecordRequest(ctx, nil, requestID, jsonBody, tags); errr != nil {
			logger.Error("%s: cannot record request: %v", "error", errr)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.logAndReturnError("cannot send repayment request", err, logger, requestID, tags)
	}
	defer c.safeClose(resp.Body, logger)

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, c.logAndReturnError("cannot read repayment response", err, logger, requestID, tags)
	}

	logger.Debug("Response: %v", string(raw))
	logger.Debug("Response status: %v", resp.StatusCode)

	if c.recorder != nil {
		if errr := c.recorder.RecordResponse(ctx, nil, requestID, raw, tags); errr != nil {
			logger.Error("%s: cannot record response %v", "error", errr)
		}
	}

	isJSON := false
	if ct := resp.Header.Get("Content-Type"); ct != "" && (ct == "application/json" || bytes.Contains([]byte(ct), []byte("application/json"))) {
		isJSON = true
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		isJSON = true
	}

	// If the server returned JSON, attempt to parse and surface API errors.
	if isJSON {
		parsed, parseErr := repayment.UnmarshalJSONResponse(raw)
		if parseErr != nil {
			return raw, c.logAndReturnError("cannot unmarshal repayment response", parseErr, logger, requestID, tags)
		}
		if apiErr := parsed.GetError(); apiErr != nil {
			return raw, apiErr
		}
	}

	// Treat non-2xx statuses as errors even when the payload is not JSON.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return raw, fmt.Errorf("repayment processing file: unexpected HTTP status %d", resp.StatusCode)
	}

	return raw, nil
}

func (c *Client) sendRepaymentMultipartRequest(apiURL string, apiRequest *repayment.RequestWrapper, fileName string, file io.Reader, logger *log.Logger) (*repayment.Response, error) {
	requestID := uuid.New().String()
	logger.Debug("Request ID: %v", requestID)

	if apiRequest == nil {
		return nil, c.logAndReturnError("repayment request is nil", fmt.Errorf("request is nil"), logger, requestID, nil)
	}
	if file == nil {
		return nil, c.logAndReturnError("repayment file is nil", fmt.Errorf("file is nil"), logger, requestID, tagsRetrieverRepayment(apiRequest))
	}
	if fileName == "" {
		return nil, c.logAndReturnError("repayment file name is empty", fmt.Errorf("file name is empty"), logger, requestID, tagsRetrieverRepayment(apiRequest))
	}

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, c.logAndReturnError("cannot marshal repayment request", err, logger, requestID, tagsRetrieverRepayment(apiRequest))
	}

	logger.Debug("Request: %v", string(jsonBody))
	logger.Debug("File: %s", fileName)

	ctx := context.WithValue(context.Background(), CtxKeyRequestID, requestID)
	tags := tagsRetrieverRepayment(apiRequest)

	bodyReader, contentType := buildRepaymentMultipartBody(jsonBody, fileName, file)

	req, err := http.NewRequest("POST", apiURL, bodyReader)
	if err != nil {
		return nil, c.logAndReturnError("cannot create repayment request", err, logger, requestID, tags)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GO IPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)

	if c.recorder != nil {
		if errr := c.recorder.RecordRequest(ctx, nil, requestID, jsonBody, tags); errr != nil {
			logger.Error("%s: cannot record request: %v", "error", errr)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.logAndReturnError("cannot send repayment request", err, logger, requestID, tags)
	}
	defer c.safeClose(resp.Body, logger)

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, c.logAndReturnError("cannot read repayment response", err, logger, requestID, tags)
	}

	logger.Debug("Response: %v", string(raw))
	logger.Debug("Response status: %v", resp.StatusCode)

	if c.recorder != nil {
		if errr := c.recorder.RecordResponse(ctx, nil, requestID, raw, tags); errr != nil {
			logger.Error("%s: cannot record response %v", "error", errr)
		}
	}

	if !isLikelyJSONResponse(resp, raw) {
		apiErr := nonJSONRepaymentAPIError(resp, raw)
		return nil, c.logAndReturnError("repayment API returned non-JSON response", apiErr, logger, requestID, tags)
	}

	response, err := repayment.UnmarshalJSONResponse(raw)
	if err != nil {
		return nil, c.logAndReturnError("cannot unmarshal repayment response", err, logger, requestID, tags)
	}

	return response, response.GetError()
}

func isLikelyJSONResponse(resp *http.Response, raw []byte) bool {
	if resp != nil {
		ct := strings.ToLower(resp.Header.Get("Content-Type"))
		if ct != "" && strings.Contains(ct, "application/json") {
			return true
		}
	}

	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[')
}

func nonJSONRepaymentAPIError(resp *http.Response, raw []byte) error {
	const maxBodyLen = 4096

	status := 0
	contentType := ""
	if resp != nil {
		status = resp.StatusCode
		contentType = resp.Header.Get("Content-Type")
	}

	trimmed := bytes.TrimSpace(raw)
	body := strings.TrimSpace(string(trimmed))

	if body == "" {
		body = "<empty body>"
	}
	if len(body) > maxBodyLen {
		body = body[:maxBodyLen] + "...(truncated)"
	}

	// Include HTTP metadata to help debug gateways/proxies returning plain text.
	msg := fmt.Sprintf("unexpected non-JSON response (status=%d, content-type=%q): %s", status, contentType, body)
	return &repayment.APIError{Message: msg}
}

func buildRepaymentMultipartBody(jsonBody []byte, fileName string, file io.Reader) (io.Reader, string) {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)
	contentType := w.FormDataContentType()

	go func() {
		reqHeader := make(textproto.MIMEHeader)
		reqHeader.Set("Content-Disposition", `form-data; name="request"`)
		reqHeader.Set("Content-Type", "application/json; charset=utf-8")

		reqPart, err := w.CreatePart(reqHeader)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("repayment multipart: create request field: %w", err))
			return
		}

		if _, err := reqPart.Write(jsonBody); err != nil {
			_ = pw.CloseWithError(fmt.Errorf("repayment multipart: write request field: %w", err))
			return
		}

		escapedFileName := escapeQuotes(fileName)

		fileHeader := make(textproto.MIMEHeader)
		fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapedFileName))
		fileHeader.Set("Content-Type", "text/csv")

		part, err := w.CreatePart(fileHeader)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("repayment multipart: create file field: %w", err))
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			_ = pw.CloseWithError(fmt.Errorf("repayment multipart: copy file: %w", err))
			return
		}

		// Close multipart writer first to ensure the closing boundary is written.
		if err := w.Close(); err != nil {
			_ = pw.CloseWithError(fmt.Errorf("repayment multipart: close writer: %w", err))
			return
		}

		_ = pw.Close()
	}()

	return pr, contentType
}

func tagsRetrieverRepayment(request *repayment.RequestWrapper) map[string]string {
	tags := make(map[string]string)

	if request == nil {
		return tags
	}

	if request.Request.Body.ExtID != nil && *request.Request.Body.ExtID != "" {
		tags["ext_id"] = *request.Request.Body.ExtID
	}

	if request.Request.Body.RepaymentGUID != nil && *request.Request.Body.RepaymentGUID != "" {
		tags["repayment_guid"] = *request.Request.Body.RepaymentGUID
	}

	if request.Request.Body.MchID != nil && *request.Request.Body.MchID != 0 {
		tags["mch_id"] = fmt.Sprintf("%d", *request.Request.Body.MchID)
	}

	if request.Operation != "" {
		tags["operation"] = request.Operation
	}

	return tags
}

var repaymentQuoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func escapeQuotes(s string) string {
	return repaymentQuoteEscaper.Replace(s)
}
