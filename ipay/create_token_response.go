package ipay

import (
	"encoding/json"
	"fmt"
)

type IpayResponseWrapper struct {
	Response Response `json:"response"`
}

type Response struct {
	PmtId int    `json:"pmt_id"`
	Url   string `json:"url"`
	Salt  string `json:"salt"`
	Sign  string `json:"sign"`
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

func UnmarshalCreateTokenResponse(data []byte) (*Response, error) {
	var resp IpayResponseWrapper
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling CreateToken response: %w", err)
	}
	return &resp.Response, nil
}
