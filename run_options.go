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
	"encoding/json"
	"fmt"

	"github.com/stremovskyy/go-ipay/ipay"
	"github.com/stremovskyy/go-ipay/log"
)

// RunOption controls the behavior of a single API call.
type RunOption func(*runOptions)

// DryRunHandler receives information about a skipped request.
type DryRunHandler func(endpoint string, payload any)

type runOptions struct {
	dryRun       bool
	dryRunHandle DryRunHandler
}

var dryRunLogger = log.NewLogger("iPay DryRun:")

// DryRun skips the underlying HTTP call. An optional handler can be provided to inspect the request payload.
func DryRun(handler ...DryRunHandler) RunOption {
	return func(o *runOptions) {
		o.dryRun = true

		if len(handler) > 0 && handler[0] != nil {
			o.dryRunHandle = handler[0]
			return
		}

		o.dryRunHandle = defaultDryRunHandler
	}
}

func collectRunOptions(opts []RunOption) *runOptions {
	if len(opts) == 0 {
		return nil
	}

	r := &runOptions{}

	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}

	return r
}

func (o *runOptions) isDryRun() bool {
	return o != nil && o.dryRun
}

func (o *runOptions) handleDryRun(endpoint string, payload any) {
	if o == nil || !o.dryRun {
		return
	}

	if o.dryRunHandle != nil {
		o.dryRunHandle(endpoint, payload)
	}
}

func defaultDryRunHandler(endpoint string, payload any) {
	dryRunLogger.Info("Dry run: skipping request to %s", endpoint)

	if payload == nil {
		dryRunLogger.Info("Dry run payload: <nil>")
		return
	}

	switch req := payload.(type) {
	case *ipay.RequestWrapper:
		body := struct {
			Operation string       `json:"operation"`
			Request   ipay.Request `json:"request"`
		}{
			Operation: req.Operation,
			Request:   req.Request,
		}

		dryRunLogger.Info("Dry run payload:\n%s", marshalIndent(body))
	case interface{ Marshal() ([]byte, error) }:
		raw, err := req.Marshal()
		if err != nil {
			dryRunLogger.Info("Dry run payload marshal error for %T: %v", req, err)
			return
		}
		dryRunLogger.Info("Dry run payload:\n%s", string(raw))
	default:
		dryRunLogger.Info("Dry run payload:\n%s", marshalIndent(req))
	}
}

func marshalIndent(v any) string {
	if v == nil {
		return "<nil>"
	}

	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("unable to marshal %T: %v", v, err)
	}

	return string(out)
}
