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

package go_ipay

import (
	"strconv"

	"github.com/stremovskyy/go-ipay/internal/ipay"
)

type Merchant struct {
	// Merchant Name
	Name string
	// Merchant ID
	MerchantID string
	// Merchant Key
	MerchantKey string
	// System Key
	SystemKey string
	// Sub Merchant ID
	SubMerchantID int

	// Login
	Login string

	// SuccessRedirect
	SuccessRedirect string

	// FailRedirect
	FailRedirect string

	signer ipay.Signer
}

func (m *Merchant) GetMerchantID() *int64 {
	id, err := strconv.ParseInt(m.MerchantID, 10, 64)

	if err != nil {
		return nil
	}

	return &id
}

func (m *Merchant) GetSign() ipay.Sign {
	if m.signer == nil {
		m.signer = ipay.NewSigner(m.SystemKey)
	}

	return *m.signer.Sign(m.MerchantKey)
}

func (m *Merchant) GetMobileSign() ipay.MobileSign {
	if m.signer == nil {
		m.signer = ipay.NewSigner(m.SystemKey)
	}

	return *m.signer.MobileSign(m.MerchantKey)
}

func (m *Merchant) GetMobileLogin() *string {
	if m.Login == "" {
		return nil
	}

	return &m.Login
}
