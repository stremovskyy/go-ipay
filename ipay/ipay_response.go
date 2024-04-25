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

package ipay

import (
	"encoding/json"
	"fmt"

	"github.com/stremovskyy/go-ipay/internal/ipay"
)

type IpayResponseWrapper struct {
	Response Response `json:"response"`
}

type Response struct {
	Transactions []ResponseTransaction `json:"transactions"`
	PmtId        int64                 `json:"pmt_id"`
	Pmt          *Payment              `json:"pmt"`
	Url          string                `json:"url"`
	Salt         string                `json:"salt"`
	Sign         string                `json:"sign"`
	Status       PaymentStatus         `json:"status"`
	BnkErrorNote *ipay.StatusCode      `json:"bnk_error_note"`
	ResAuthCode  int                   `json:"res_auth_code"`
	Error        *string               `json:"error"`
	ErrorCode    *string               `json:"error_code"`
}

func (r Response) GetError() error {
	if r.BnkErrorNote != nil {
		if statusCode, found := ipay.GetStatusCode(*r.BnkErrorNote); found {
			return fmt.Errorf(fmt.Sprintf("bank error: %s, reason: %s, message: %s", *r.BnkErrorNote, statusCode.Reason, statusCode.Message))
		} else {
			return fmt.Errorf("general error: %s", *r.BnkErrorNote)
		}
	}

	if r.Status == PaymentStatusSecurityRefusal {
		return fmt.Errorf("payment status: security refusal")
	}

	if r.Status == PaymentStatusFailed {
		return fmt.Errorf("payment status: payment failed")
	}

	if r.Error != nil {
		if r.ErrorCode != nil {
			return fmt.Errorf("ipay general error: %s, code: %s", *r.Error, *r.ErrorCode)
		} else {
			return fmt.Errorf("ipay general error: %s", *r.Error)
		}
	}

	return nil
}

type ResponseTransaction struct {
	TrnId    *int    `json:"trn_id"`
	SmchRr   *int    `json:"smch_rr"`
	SmchMfo  *int    `json:"smch_mfo"`
	SmchOkpo *int    `json:"smch_okpo"`
	SmchBank *string `json:"smch_bank"`
}

func (ctr *IpayResponseWrapper) Debug() string {
	return fmt.Sprintf(
		"Debug Info:\nPayment ID: %d\nValidation URL: %s\nSalt: %s\nSignature: %s\n",
		ctr.Response.PmtId,
		ctr.Response.Url,
		ctr.Response.Salt,
		ctr.Response.Sign,
	)
}

func UnmarshalJSONResponse(data []byte) (*Response, error) {
	var resp IpayResponseWrapper
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}
	return &resp.Response, nil
}
