package repayment

// Auth holds the authentication details required for a Repayment API request.
type Auth struct {
	Login string `json:"login"`
	Time  string `json:"time"`
	Sign  string `json:"sign"`
}
