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

package repayment

type Action string

const (
	ActionCreateRepayment            Action = "CreateRepayment"
	ActionCancelRepayment            Action = "CancelRepayment"
	ActionGetRepaymentStatus         Action = "GetRepaymentStatus"
	ActionGetRepaymentProcessingFile Action = "GetRepaymentProcessingFile"
)

// RequestWrapper is the outer payload for Repayment API requests.
// For CreateRepayment it is JSON-marshaled and sent as multipart/form-data field named "request".
// For other actions it is sent as an application/json payload.
type RequestWrapper struct {
	Request Request `json:"request"`

	// Operation is not part of the API payload; it is used for logging/recording.
	Operation string `json:"-"`
}

type Request struct {
	Auth   Auth   `json:"auth"`
	Action Action `json:"action"`
	Body   Body   `json:"body"`
}

type Body struct {
	MchID         *int64  `json:"mch_id,omitempty"`
	ExtID         *string `json:"ext_id,omitempty"`
	SmchID        *int64  `json:"smch_id,omitempty"`
	RepaymentGUID *string `json:"repayment_guid,omitempty"`
}
