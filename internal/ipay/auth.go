package ipay

// Auth holds the authentication details required for a payment request.
type Auth struct {
	MchID int64  `json:"mch_id" xml:"mch_id"` // Merchant ID.
	Salt  string `json:"salt" xml:"salt"`     // Salt for signature.
	Sign  string `json:"sign" xml:"sign"`     // Request signature.
}
