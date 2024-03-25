package ipay

type Sign struct {
	Salt string `json:"salt"`
	Sign string `json:"sign"`
}
