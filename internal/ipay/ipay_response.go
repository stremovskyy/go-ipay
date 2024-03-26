package ipay

import (
	"encoding/json"
	"fmt"
)

type IpayResponseWrapper struct {
	Response Response `json:"response"`
}

type Response struct {
	Transactions []ResponseTransaction `json:"transactions"`
	PmtId        int                   `json:"pmt_id"`
	Pmt          *Payment              `json:"pmt"`
	Url          string                `json:"url"`
	Salt         string                `json:"salt"`
	Sign         string                `json:"sign"`
	Status       PaymentStatus         `json:"status"`
	BnkErrorNote *StatusCode           `json:"bnk_error_note"`
	Error        *string               `json:"error"`
	ErrorCode    *string               `json:"error_code"`
}

func (r Response) GetError() error {
	if r.BnkErrorNote != nil {
		if statusCode, found := GetStatusCode(*r.BnkErrorNote); found {
			return fmt.Errorf(fmt.Sprintf("ipay error: %s, reason: %s, message: %s", *r.BnkErrorNote, statusCode.Reason, statusCode.Message))
		} else {
			return fmt.Errorf("ipay error: %s", *r.BnkErrorNote)
		}
	}

	if r.Status == PaymentStatusSecurityRefusal {
		return fmt.Errorf("ipay error: security refusal")
	}

	if r.Status == PaymentStatusCanceled {
		return fmt.Errorf("ipay error: payment canceled")
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
