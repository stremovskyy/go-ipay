/*
 * MIT License
 *
 * Copyright (c) 2026 Anton Stremovskyy
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
	"io"
	"mime"
	stdhttp "net/http"
	"strconv"
	"strings"
)

const (
	maxResponseBodyBytes = 4 << 20
	maxErrorBodyBytes    = 4 << 10
)

func readBoundedResponseBody(body io.Reader) ([]byte, bool, error) {
	raw, err := io.ReadAll(io.LimitReader(body, maxResponseBodyBytes+1))
	tooLarge := len(raw) > maxResponseBodyBytes
	if tooLarge {
		raw = raw[:maxResponseBodyBytes]
	}

	return raw, tooLarge, err
}

func diagnosticResponseBody(raw []byte) []byte {
	if len(raw) > maxErrorBodyBytes {
		raw = raw[:maxErrorBodyBytes]
	}

	return append([]byte(nil), raw...)
}

func isLikelyIPayJSONResponse(raw []byte) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[')
}

func responseContentType(resp *stdhttp.Response) string {
	if resp == nil {
		return ""
	}

	raw := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if raw == "" {
		return ""
	}

	mediaType, _, err := mime.ParseMediaType(raw)
	if err == nil {
		return strings.ToLower(mediaType)
	}

	return strings.ToLower(raw)
}

func cloneTags(tags map[string]string) map[string]string {
	cloned := make(map[string]string, len(tags)+3)
	for key, value := range tags {
		cloned[key] = value
	}

	return cloned
}

func responseTags(tags map[string]string, resp *stdhttp.Response) map[string]string {
	cloned := cloneTags(tags)
	if resp == nil {
		return cloned
	}

	cloned["http_status"] = strconv.Itoa(resp.StatusCode)
	if contentType := responseContentType(resp); contentType != "" {
		cloned["response_content_type"] = contentType
	}

	return cloned
}

func errorTags(tags map[string]string, kind string) map[string]string {
	cloned := cloneTags(tags)
	if kind != "" {
		cloned["response_error_kind"] = kind
	}

	return cloned
}
