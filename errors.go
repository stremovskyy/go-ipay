package go_ipay

import "errors"

var ErrRequestIsNil = errors.New("request is nil")
var ErrMerchantIsNil = errors.New("merchant is nil")
var ErrPersonalDataIsNil = errors.New("personal data is nil")
