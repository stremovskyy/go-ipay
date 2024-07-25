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
	"strconv"

	"github.com/stremovskyy/go-ipay/internal/ipay"
)

type ResponseWrapper struct {
	Response Response `json:"response"`
}

type Response struct {
	Transactions     []ResponseTransaction `json:"transactions"`
	PmtId            interface{}           `json:"pmt_id"`
	Pmt              *Payment              `json:"pmt"`
	Url              string                `json:"url"`
	Salt             string                `json:"salt"`
	Sign             string                `json:"sign"`
	Status           *PaymentStatus        `json:"status"`
	BnkErrorNote     *ipay.StatusCode      `json:"bnk_error_note"`
	ResAuthCode      int                   `json:"res_auth_code"`
	Error            *string               `json:"error"`
	ErrorCode        *string               `json:"error_code"`
	Invoice          interface{}           `json:"invoice"`
	Amount           interface{}           `json:"amount"`
	PmtStatus        *string               `json:"pmt_status"`
	CardMask         *string               `json:"card_mask"`
	BankResponse     *BankResponse         `json:"bank_response"`
	BankAcquirerName *string               `json:"bank_acquirer_name"`
}

type BankResponse struct {
	ErrorGroup int `json:"error_group"`
}

func (r Response) GetPaymentStatus() PaymentStatus {
	if r.Status != nil {
		return *r.Status
	}

	if r.PmtStatus != nil {
		return r.mobilePaymentStatus()
	}

	return PaymentStatusUnknown
}

func (r Response) mobilePaymentStatus() PaymentStatus {
	if r.PmtStatus == nil {
		return PaymentStatusUnknown
	}

	parsedIntStatus, err := strconv.Atoi(*r.PmtStatus)
	if err != nil {
		return PaymentStatusUnknown
	}

	return PaymentStatus(parsedIntStatus)
}
func (r Response) getInt64FromInterface(value interface{}) int64 {
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int, int64, float64:
		return int64(v.(float64))
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func (r Response) PmtIdInt64() int64 {
	return r.getInt64FromInterface(r.PmtId)
}

func (r Response) AmountInt64() int64 {
	return r.getInt64FromInterface(r.Amount)
}

func (r Response) InvoiceAmountInt64() int64 {
	return r.getInt64FromInterface(r.Invoice)
}

func (r Response) GetError() error {
	// Check if there's a general error message
	if r.Error != nil {
		if r.ErrorCode != nil {
			return createIpayError(
				900,
				fmt.Sprintf("ipay general error: %s", *r.Error),
				fmt.Sprintf("code: %s", *r.ErrorCode),
			)
		}
		return createIpayError(900, *r.Error, "")
	}

	// Check if there's a bank error note
	if r.BnkErrorNote != nil {
		if statusCode, found := ipay.GetStatusCode(*r.BnkErrorNote); found {
			return createIpayError(
				statusCode.ExtCode,
				fmt.Sprintf("bank error: %s", *r.BnkErrorNote),
				fmt.Sprintf("reason: %s, message: %s", statusCode.Reason, statusCode.Message),
			)
		}
		return createIpayError(900, string(*r.BnkErrorNote), "")
	}

	// Check if there's a specific authorization code error
	if r.ResAuthCode != 0 {
		message := getErrorMessageA2CPay(r.ResAuthCode)
		return createIpayError(r.ResAuthCode, message, "")
	}

	// Check for payment status errors
	if r.GetPaymentStatus() == PaymentStatusSecurityRefusal {
		return createIpayError(900, "payment status: security refusal", "")
	}

	if r.GetPaymentStatus() == PaymentStatusFailed {
		if r.BankResponse != nil {
			if r.BankResponse.ErrorGroup != 0 {
				return createIpayError(
					r.BankResponse.ErrorGroup,
					"payment status: payment failed",
					fmt.Sprintf("bank acquirer name: %v", r.BankAcquirerName),
				)
			}
		}

		return createIpayError(900, "payment status: payment failed", "")
	}

	// No errors found
	return nil
}

type ResponseTransaction struct {
	TrnId    *int    `json:"trn_id"`
	SmchRr   *int    `json:"smch_rr"`
	SmchMfo  *int    `json:"smch_mfo"`
	SmchOkpo *int    `json:"smch_okpo"`
	SmchBank *string `json:"smch_bank"`
}

func (ctr *ResponseWrapper) Debug() string {
	return fmt.Sprintf(
		"Debug Info:\nPayment ID: %d\nValidation URL: %s\nSalt: %s\nSignature: %s\n",
		ctr.Response.PmtId,
		ctr.Response.Url,
		ctr.Response.Salt,
		ctr.Response.Sign,
	)
}

func UnmarshalJSONResponse(data []byte) (*Response, error) {
	var resp ResponseWrapper

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}

	return &resp.Response, nil
}
