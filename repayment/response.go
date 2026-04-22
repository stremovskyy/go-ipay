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

package repayment

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/stremovskyy/go-ipay/internal/utils"
)

type ResponseWrapper struct {
	Response Response `json:"response"`
}

type Response struct {
	RepaymentGUID   *string `json:"repayment_guid"`
	ExtID           *string `json:"ext_id"`
	Status          *int    `json:"status"`
	Invoice         *int    `json:"invoice"`
	Amount          *int    `json:"amount"`
	MchID           *int64  `json:"mch_id"`
	MchBalance      *int    `json:"mch_balance"`
	SuccessPayments *int    `json:"success_payments"`
	FailedPayments  *int    `json:"failed_payments"`

	Error *string `json:"error"`
}

func (r Response) GetError() error {
	if r.Error == nil || *r.Error == "" {
		return nil
	}

	return &APIError{Message: *r.Error}
}

func UnmarshalJSONResponse(data []byte) (*Response, error) {
	if len(data) == 0 {
		return &Response{Error: utils.Ref("empty response data")}, nil
	}

	var raw struct {
		Response struct {
			RepaymentGUID   *string         `json:"repayment_guid"`
			ExtID           *string         `json:"ext_id"`
			Status          json.RawMessage `json:"status"`
			Invoice         json.RawMessage `json:"invoice"`
			Amount          json.RawMessage `json:"amount"`
			MchID           json.RawMessage `json:"mch_id"`
			MchBalance      json.RawMessage `json:"mch_balance"`
			SuccessPayments json.RawMessage `json:"success_payments"`
			FailedPayments  json.RawMessage `json:"failed_payments"`
			Error           *string         `json:"error"`
		} `json:"response"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response: %w", err)
	}

	resp := &Response{
		RepaymentGUID: raw.Response.RepaymentGUID,
		ExtID:         raw.Response.ExtID,
		Error:         raw.Response.Error,
	}

	var err error

	resp.Status, err = parseOptionalInt(raw.Response.Status)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response status: %w", err)
	}
	resp.Invoice, err = parseOptionalInt(raw.Response.Invoice)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response invoice: %w", err)
	}
	resp.Amount, err = parseOptionalInt(raw.Response.Amount)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response amount: %w", err)
	}
	resp.MchID, err = parseOptionalInt64(raw.Response.MchID)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response mch_id: %w", err)
	}
	resp.MchBalance, err = parseOptionalInt(raw.Response.MchBalance)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response mch_balance: %w", err)
	}
	resp.SuccessPayments, err = parseOptionalInt(raw.Response.SuccessPayments)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response success_payments: %w", err)
	}
	resp.FailedPayments, err = parseOptionalInt(raw.Response.FailedPayments)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response failed_payments: %w", err)
	}

	return resp, nil
}

func parseOptionalInt(raw json.RawMessage) (*int, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	var value int
	if err := json.Unmarshal(raw, &value); err == nil {
		return &value, nil
	}

	var str string
	if err := json.Unmarshal(raw, &str); err != nil {
		return nil, fmt.Errorf("invalid int payload %q", string(raw))
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("invalid int string %q: %w", str, err)
	}

	return &value, nil
}

func parseOptionalInt64(raw json.RawMessage) (*int64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	var value int64
	if err := json.Unmarshal(raw, &value); err == nil {
		return &value, nil
	}

	var str string
	if err := json.Unmarshal(raw, &str); err != nil {
		return nil, fmt.Errorf("invalid int64 payload %q", string(raw))
	}

	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid int64 string %q: %w", str, err)
	}

	return &value, nil
}
