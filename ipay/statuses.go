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

type PaymentStatus int

const (
	PaymentStatusUnknown             PaymentStatus = 0
	PaymentStatusRegistered          PaymentStatus = 1
	PaymentStatusPreAuthorized       PaymentStatus = 3
	PaymentStatusFailed              PaymentStatus = 4
	PaymentStatusSuccess             PaymentStatus = 5
	PaymentStatusCanceled            PaymentStatus = 9
	PaymentStatusManualProcessing    PaymentStatus = 11
	PaymentStatusSuccessWithoutClaim PaymentStatus = 13
	PaymentStatusSecurityRefusal     PaymentStatus = 14
)

func (s *PaymentStatus) String() string {
	if s == nil {
		return "Unknown"
	}

	switch *s {
	case PaymentStatusRegistered:
		return "Registered"
	case PaymentStatusPreAuthorized:
		return "PreAuthorized"
	case PaymentStatusFailed:
		return "Failed"
	case PaymentStatusSuccess:
		return "Success"
	case PaymentStatusCanceled:
		return "Canceled"
	case PaymentStatusManualProcessing:
		return "ManualProcessing"
	case PaymentStatusSuccessWithoutClaim:
		return "SuccessWithoutClaim"
	case PaymentStatusSecurityRefusal:
		return "SecurityRefusal"
	default:
		return "Unknown"
	}
}

func (s *PaymentStatus) Is(paymentStatus PaymentStatus) bool {
	if s == nil {
		return false
	}

	return *s == paymentStatus
}
