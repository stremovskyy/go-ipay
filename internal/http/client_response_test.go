package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	"github.com/stremovskyy/go-ipay/internal/teststand"
	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/recorder"
)

func TestSendRequestResponseClassification(t *testing.T) {
	tests := []struct {
		name                string
		action              ipay.Action
		response            *stdhttp.Response
		transportErr        error
		wantSentinel        error
		wantStatus          int
		wantContentType     string
		wantBody            string
		wantStatusCheck     bool
		wantIpayErrorCode   int
		wantResponseSuccess bool
	}{
		{
			name:                "normal JSON response",
			action:              ipay.ActionCredit,
			response:            teststand.Response(stdhttp.StatusOK, "application/json", []byte(`{"response":{"pmt_id":12345,"pmt_status":"5"}}`)),
			wantResponseSuccess: true,
		},
		{
			name:            "empty body",
			action:          ipay.ActionCredit,
			response:        teststand.Response(stdhttp.StatusOK, "application/json", nil),
			wantSentinel:    ipay.ErrUnexpectedResponse,
			wantStatus:      stdhttp.StatusOK,
			wantContentType: "application/json",
			wantStatusCheck: true,
		},
		{
			name:            "whitespace body",
			action:          ipay.ActionCredit,
			response:        teststand.Response(stdhttp.StatusOK, "application/json; charset=utf-8", []byte(" \n\t")),
			wantSentinel:    ipay.ErrUnexpectedResponse,
			wantStatus:      stdhttp.StatusOK,
			wantContentType: "application/json",
			wantBody:        " \n\t",
			wantStatusCheck: true,
		},
		{
			name:            "HTML gateway timeout",
			action:          ipay.ActionCredit,
			response:        teststand.Response(stdhttp.StatusGatewayTimeout, "text/html; charset=utf-8", []byte("<html>secret-gateway-page</html>")),
			wantSentinel:    ipay.ErrUnexpectedResponse,
			wantStatus:      stdhttp.StatusGatewayTimeout,
			wantContentType: "text/html",
			wantBody:        "<html>secret-gateway-page</html>",
			wantStatusCheck: true,
		},
		{
			name:            "HTML with success status",
			action:          ipay.ActionCredit,
			response:        teststand.Response(stdhttp.StatusOK, "text/html", []byte("<html>upstream login</html>")),
			wantSentinel:    ipay.ErrUnexpectedResponse,
			wantStatus:      stdhttp.StatusOK,
			wantContentType: "text/html",
			wantBody:        "<html>upstream login</html>",
			wantStatusCheck: true,
		},
		{
			name:            "HTML despite JSON content type",
			action:          ipay.ActionCredit,
			response:        teststand.Response(stdhttp.StatusOK, "application/json", []byte("<html>proxy failure</html>")),
			wantSentinel:    ipay.ErrUnexpectedResponse,
			wantStatus:      stdhttp.StatusOK,
			wantContentType: "application/json",
			wantBody:        "<html>proxy failure</html>",
			wantStatusCheck: true,
		},
		{
			name:            "malformed JSON",
			action:          ipay.ActionCredit,
			response:        teststand.Response(stdhttp.StatusOK, "application/json", []byte(`{"response":`)),
			wantSentinel:    ipay.ErrDecode,
			wantStatus:      stdhttp.StatusOK,
			wantContentType: "application/json",
			wantBody:        `{"response":`,
			wantStatusCheck: true,
		},
		{
			name:            "network timeout",
			action:          ipay.ActionCredit,
			transportErr:    timeoutError{},
			wantSentinel:    ipay.ErrTransport,
			wantStatusCheck: true,
		},
		{
			name:            "status check timeout is not A2C pay outcome",
			action:          ipay.ActionA2CPaymentStatus,
			transportErr:    timeoutError{},
			wantSentinel:    ipay.ErrTransport,
			wantStatusCheck: false,
		},
		{
			name:              "non-2xx JSON business refusal",
			action:            ipay.ActionCredit,
			response:          teststand.Response(stdhttp.StatusBadRequest, "application/json", []byte(`{"response":{"error":"declined","error_code":"601"}}`)),
			wantSentinel:      ipay.ErrUnexpectedResponse,
			wantStatus:        stdhttp.StatusBadRequest,
			wantContentType:   "application/json",
			wantBody:          `{"response":{"error":"declined","error_code":"601"}}`,
			wantIpayErrorCode: 601,
			wantStatusCheck:   false,
		},
		{
			name:              "2xx JSON business refusal",
			action:            ipay.ActionCredit,
			response:          teststand.Response(stdhttp.StatusOK, "application/json", []byte(`{"response":{"error":"declined","error_code":"601"}}`)),
			wantIpayErrorCode: 601,
			wantStatusCheck:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := clientWithTransport(func(*stdhttp.Request) (*stdhttp.Response, error) {
				return tt.response, tt.transportErr
			})
			request := ipay.NewRequest(tt.action, ipay.WithOperationOperation("test-operation"))

			response, err := client.sendRequest("https://api.example.test", request, client.loggerFor(loggerTypeHTTP))
			if tt.wantResponseSuccess {
				if err != nil {
					t.Fatalf("sendRequest() error = %v", err)
				}
				if response == nil || response.PmtIdInt64() != 12345 {
					t.Fatalf("sendRequest() response = %#v, want pmt_id 12345", response)
				}
				return
			}

			if err == nil {
				t.Fatal("sendRequest() error = nil")
			}
			if tt.wantSentinel != nil && !errors.Is(err, tt.wantSentinel) {
				t.Fatalf("errors.Is(%v) = false; error = %v", tt.wantSentinel, err)
			}
			if got := ipay.RequiresA2CStatusCheck(err); got != tt.wantStatusCheck {
				t.Fatalf("RequiresA2CStatusCheck() = %v, want %v; error = %v", got, tt.wantStatusCheck, err)
			}

			if tt.wantIpayErrorCode != 0 {
				var providerErr *ipay.IpayError
				if !errors.As(err, &providerErr) {
					t.Fatalf("errors.As(*IpayError) = false; error = %v", err)
				}
				if providerErr.Code != tt.wantIpayErrorCode {
					t.Fatalf("IpayError.Code = %d, want %d", providerErr.Code, tt.wantIpayErrorCode)
				}
			}

			var unexpectedErr *ipay.UnexpectedResponseError
			if errors.As(err, &unexpectedErr) {
				assertResponseErrorContext(t, unexpectedErr.StatusCode, unexpectedErr.ContentType, unexpectedErr.Body, tt.wantStatus, tt.wantContentType, tt.wantBody)
				if strings.Contains(err.Error(), "secret-gateway-page") ||
					strings.Contains(err.Error(), "upstream login") ||
					strings.Contains(err.Error(), "proxy failure") {
					t.Fatalf("error contains raw response body: %v", err)
				}
			}

			var decodeErr *ipay.DecodeError
			if errors.As(err, &decodeErr) {
				assertResponseErrorContext(t, decodeErr.StatusCode, decodeErr.ContentType, decodeErr.Body, tt.wantStatus, tt.wantContentType, tt.wantBody)
			}
		})
	}
}

func TestSendRequestOversizedResponse(t *testing.T) {
	body := bytes.Repeat([]byte("x"), maxResponseBodyBytes+1)
	client := clientWithTransport(func(*stdhttp.Request) (*stdhttp.Response, error) {
		return teststand.Response(stdhttp.StatusOK, "application/json", body), nil
	})

	_, err := client.sendRequest(
		"https://api.example.test",
		ipay.NewRequest(ipay.ActionCredit),
		client.loggerFor(loggerTypeHTTP),
	)
	if !errors.Is(err, ipay.ErrUnexpectedResponse) {
		t.Fatalf("sendRequest() error = %v, want ErrUnexpectedResponse", err)
	}
	if !ipay.RequiresA2CStatusCheck(err) {
		t.Fatalf("RequiresA2CStatusCheck() = false; error = %v", err)
	}

	var responseErr *ipay.UnexpectedResponseError
	if !errors.As(err, &responseErr) {
		t.Fatalf("errors.As(*UnexpectedResponseError) = false; error = %v", err)
	}
	if len(responseErr.Body) != maxErrorBodyBytes {
		t.Fatalf("len(Body) = %d, want %d", len(responseErr.Body), maxErrorBodyBytes)
	}
	body[0] = 'y'
	if responseErr.Body[0] != 'x' {
		t.Fatal("diagnostic Body aliases the source response")
	}
	if strings.Contains(err.Error(), strings.Repeat("x", 32)) {
		t.Fatalf("error contains raw response body: %v", err)
	}
}

func TestSendRequestNilBody(t *testing.T) {
	client := clientWithTransport(func(*stdhttp.Request) (*stdhttp.Response, error) {
		return &stdhttp.Response{
			StatusCode:    stdhttp.StatusOK,
			Status:        "200 OK",
			Header:        stdhttp.Header{"Content-Type": []string{"application/json"}},
			Body:          nil,
			ContentLength: 1,
		}, nil
	})

	_, err := client.sendRequest(
		"https://api.example.test",
		ipay.NewRequest(ipay.ActionCredit),
		client.loggerFor(loggerTypeHTTP),
	)
	if !errors.Is(err, ipay.ErrTransport) && !errors.Is(err, ipay.ErrUnexpectedResponse) {
		t.Fatalf("sendRequest() error = %v, want typed transport or unexpected-response error", err)
	}
	if !ipay.RequiresA2CStatusCheck(err) {
		t.Fatalf("RequiresA2CStatusCheck() = false; error = %v", err)
	}
}

func TestSendRequestResponseReadError(t *testing.T) {
	rec := &recordingRecorder{}
	client := clientWithTransport(func(*stdhttp.Request) (*stdhttp.Response, error) {
		return &stdhttp.Response{
			StatusCode: stdhttp.StatusOK,
			Status:     "200 OK",
			Header:     stdhttp.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			Body:       &failingReadCloser{data: []byte("partial-sensitive-response")},
		}, nil
	})
	client.WithRecorder(rec)

	_, err := client.sendRequest(
		"https://api.example.test",
		ipay.NewRequest(ipay.ActionCredit),
		client.loggerFor(loggerTypeHTTP),
	)
	if !errors.Is(err, ipay.ErrTransport) {
		t.Fatalf("sendRequest() error = %v, want ErrTransport", err)
	}
	if !ipay.RequiresA2CStatusCheck(err) {
		t.Fatalf("RequiresA2CStatusCheck() = false; error = %v", err)
	}
	var transportErr *ipay.TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("errors.As(*TransportError) = false; error = %v", err)
	}
	if transportErr.StatusCode != stdhttp.StatusOK || transportErr.ContentType != "application/json" {
		t.Fatalf("transport response context = status %d, content type %q", transportErr.StatusCode, transportErr.ContentType)
	}
	if string(transportErr.Body) != "partial-sensitive-response" {
		t.Fatalf("TransportError.Body = %q", transportErr.Body)
	}
	if strings.Contains(err.Error(), "partial-sensitive-response") {
		t.Fatalf("error contains raw partial response: %v", err)
	}
	if rec.responseCount != 1 || rec.errorCount != 1 {
		t.Fatalf("recorded responses/errors = %d/%d, want 1/1", rec.responseCount, rec.errorCount)
	}
	assertTag(t, rec.errorTags, "http_status", "200")
	assertTag(t, rec.errorTags, "response_content_type", "application/json")
	assertTag(t, rec.errorTags, "response_error_kind", "transport")
}

func TestSendRequestRecorderTagsErrors(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		contentType string
		body        []byte
	}{
		{
			name:        "gateway timeout",
			status:      stdhttp.StatusGatewayTimeout,
			contentType: "text/html; charset=utf-8",
			body:        []byte("<html>gateway timeout</html>"),
		},
		{
			name:        "empty success response",
			status:      stdhttp.StatusOK,
			contentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &recordingRecorder{}
			var sentRequestID string
			client := clientWithTransport(func(request *stdhttp.Request) (*stdhttp.Response, error) {
				sentRequestID = request.Header.Get("X-Request-ID")
				return teststand.Response(tt.status, tt.contentType, tt.body), nil
			})
			client.WithRecorder(rec)

			_, err := client.sendRequest(
				"https://api.example.test",
				ipay.NewRequest(ipay.ActionCredit, ipay.WithOperationOperation("credit")),
				client.loggerFor(loggerTypeHTTP),
			)
			if !errors.Is(err, ipay.ErrUnexpectedResponse) {
				t.Fatalf("sendRequest() error = %v, want ErrUnexpectedResponse", err)
			}
			if rec.responseCount != 1 || rec.errorCount != 1 {
				t.Fatalf("recorded responses/errors = %d/%d, want 1/1", rec.responseCount, rec.errorCount)
			}
			if rec.responseRequestID != sentRequestID || rec.errorRequestID != sentRequestID {
				t.Fatalf("recorded request IDs = %q/%q, sent %q", rec.responseRequestID, rec.errorRequestID, sentRequestID)
			}

			var responseErr *ipay.UnexpectedResponseError
			if !errors.As(err, &responseErr) {
				t.Fatalf("errors.As(*UnexpectedResponseError) = false; error = %v", err)
			}
			if responseErr.RequestID != sentRequestID {
				t.Fatalf("UnexpectedResponseError.RequestID = %q, want %q", responseErr.RequestID, sentRequestID)
			}
			if !strings.Contains(err.Error(), "status=") || !strings.Contains(err.Error(), "request_id="+sentRequestID) {
				t.Fatalf("error lacks status/request ID context: %v", err)
			}

			wantContentType := strings.Split(tt.contentType, ";")[0]
			assertTag(t, rec.responseTags, "http_status", fmt.Sprint(tt.status))
			assertTag(t, rec.responseTags, "response_content_type", wantContentType)
			assertTag(t, rec.errorTags, "http_status", fmt.Sprint(tt.status))
			assertTag(t, rec.errorTags, "response_content_type", wantContentType)
			assertTag(t, rec.errorTags, "response_error_kind", "unexpected_response")
		})
	}
}

func TestSendRequestRecorderErrorKinds(t *testing.T) {
	tests := []struct {
		name              string
		response          *stdhttp.Response
		transportErr      error
		wantKind          string
		wantResponseCount int
		wantStatus        string
	}{
		{
			name:         "transport",
			transportErr: timeoutError{},
			wantKind:     "transport",
		},
		{
			name:              "decode",
			response:          teststand.Response(stdhttp.StatusOK, "application/json", []byte(`{"response":`)),
			wantKind:          "decode",
			wantResponseCount: 1,
			wantStatus:        "200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &recordingRecorder{}
			client := clientWithTransport(func(*stdhttp.Request) (*stdhttp.Response, error) {
				return tt.response, tt.transportErr
			})
			client.WithRecorder(rec)

			_, err := client.sendRequest(
				"https://api.example.test",
				ipay.NewRequest(ipay.ActionCredit),
				client.loggerFor(loggerTypeHTTP),
			)
			if err == nil {
				t.Fatal("sendRequest() error = nil")
			}
			if rec.errorCount != 1 || rec.responseCount != tt.wantResponseCount {
				t.Fatalf("recorded responses/errors = %d/%d, want %d/1", rec.responseCount, rec.errorCount, tt.wantResponseCount)
			}
			assertTag(t, rec.errorTags, "response_error_kind", tt.wantKind)
			if tt.wantStatus != "" {
				assertTag(t, rec.errorTags, "http_status", tt.wantStatus)
			}
		})
	}
}

func TestSendRequestNilClientIsNotOutcomeUnknown(t *testing.T) {
	client := NewClient(DefaultOptions())
	client.SetClient(nil)

	_, err := client.sendRequest(
		"https://api.example.test",
		ipay.NewRequest(ipay.ActionCredit),
		client.loggerFor(loggerTypeHTTP),
	)
	if !errors.Is(err, ipay.ErrTransport) {
		t.Fatalf("sendRequest() error = %v, want ErrTransport", err)
	}
	if ipay.RequiresA2CStatusCheck(err) {
		t.Fatalf("RequiresA2CStatusCheck() = true for a request that was not sent: %v", err)
	}
}

func clientWithTransport(roundTrip func(*stdhttp.Request) (*stdhttp.Response, error)) *Client {
	client := NewClient(DefaultOptions())
	client.SetClient(&stdhttp.Client{
		Transport: teststand.RoundTripperFunc(roundTrip),
		Timeout:   time.Second,
	})
	return client
}

func assertResponseErrorContext(
	t *testing.T,
	status int,
	contentType string,
	body []byte,
	wantStatus int,
	wantContentType string,
	wantBody string,
) {
	t.Helper()
	if status != wantStatus {
		t.Fatalf("StatusCode = %d, want %d", status, wantStatus)
	}
	if contentType != wantContentType {
		t.Fatalf("ContentType = %q, want %q", contentType, wantContentType)
	}
	if string(body) != wantBody {
		t.Fatalf("Body = %q, want %q", string(body), wantBody)
	}
}

func assertTag(t *testing.T, tags map[string]string, key, want string) {
	t.Helper()
	if got := tags[key]; got != want {
		t.Fatalf("tag %q = %q, want %q; tags = %#v", key, got, want, tags)
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "request timed out" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

type failingReadCloser struct {
	data []byte
}

func (r *failingReadCloser) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, errors.New("response stream failed")
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}

func (*failingReadCloser) Close() error { return nil }

type recordingRecorder struct {
	responseCount     int
	errorCount        int
	responseRequestID string
	errorRequestID    string
	responseTags      map[string]string
	errorTags         map[string]string
}

func (r *recordingRecorder) RecordRequest(context.Context, *string, string, []byte, map[string]string) error {
	return nil
}

func (r *recordingRecorder) RecordResponse(_ context.Context, _ *string, requestID string, _ []byte, tags map[string]string) error {
	r.responseCount++
	r.responseRequestID = requestID
	r.responseTags = cloneTags(tags)
	return nil
}

func (r *recordingRecorder) RecordError(_ context.Context, _ *string, requestID string, _ error, tags map[string]string) error {
	r.errorCount++
	r.errorRequestID = requestID
	r.errorTags = cloneTags(tags)
	return nil
}

func (r *recordingRecorder) RecordMetrics(context.Context, *string, string, map[string]string, map[string]string) error {
	return nil
}

func (r *recordingRecorder) GetRequest(context.Context, string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (r *recordingRecorder) GetResponse(context.Context, string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (r *recordingRecorder) FindByTag(context.Context, string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (r *recordingRecorder) Async() recorder.AsyncRecorder { return nil }

var _ recorder.Recorder = (*recordingRecorder)(nil)
