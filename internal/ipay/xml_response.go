package ipay

import (
	"encoding/xml"
	"fmt"
)

type PaymentResponse struct {
	XMLName xml.Name `xml:"payment"`
	PID     string   `xml:"pid"`
	Status  int      `xml:"status"`
	Salt    string   `xml:"salt"`
	Sign    string   `xml:"sign"`
	URL     string   `xml:"url"`
}

func UnmarshalXmlResponse(data []byte) (*PaymentResponse, error) {
	var resp PaymentResponse
	err := xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal XML: %v", err)
	}
	return &resp, nil
}
