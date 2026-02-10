package teststand

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func Response(status int, contentType string, body []byte) *http.Response {
	h := make(http.Header)
	if contentType != "" {
		h.Set("Content-Type", contentType)
	}

	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     h,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}
