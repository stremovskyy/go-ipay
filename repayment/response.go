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

	var resp ResponseWrapper
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling repayment JSON response: %w", err)
	}

	return &resp.Response, nil
}
