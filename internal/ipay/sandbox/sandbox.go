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

package sandbox

import (
	"errors"
	"fmt"
)

type CardType int

const (
	SuccessRegardlessOfAmount CardType = iota
	SuccessIfUnder100
	PreAuthorizationPossible
	FailureRegardlessOfAmount
	PreAuthorizationRegardlessOfAmount
	SuccessA2CPay
	FailureRandomErrorA2CPay
	FailureInsufficientBalanceA2CPay
	InvalidSandboxPan
)

type Card struct {
	Number string
	Type   CardType
}

type PaymentOutcome int

const (
	PaymentSuccess PaymentOutcome = iota
	PaymentFailure
	PaymentPreAuthorized
	PaymentInvalid
)

type Sandbox interface {
	SimulatePayment(cardNumber string, amount float64) (PaymentOutcome, error)
}

type SandboxSimulator struct {
	Cards []Card
}

func NewSandboxSimulator() Sandbox {
	return &SandboxSimulator{
		Cards: []Card{
			// Always success cards
			{Number: "3333333333333331", Type: SuccessRegardlessOfAmount},
			{Number: "3333333333332705", Type: SuccessRegardlessOfAmount},
			{Number: "3333333333333000", Type: SuccessRegardlessOfAmount},
			{Number: "3333333333331640", Type: SuccessRegardlessOfAmount},
			{Number: "3333333333334909", Type: SuccessRegardlessOfAmount},
			{Number: "3333333333333703", Type: SuccessRegardlessOfAmount},
			{Number: "3333333333332820", Type: SuccessRegardlessOfAmount},
			{Number: "5375913862726080", Type: SuccessRegardlessOfAmount},

			// Success cards if amount is under 100
			{Number: "3333333333333430", Type: SuccessIfUnder100},
			{Number: "3333333333331509", Type: SuccessIfUnder100},
			{Number: "5375912476960515", Type: SuccessIfUnder100},
			{Number: "5117962099480048", Type: SuccessIfUnder100},
			{Number: "4005520000000129", Type: SuccessIfUnder100},
			{Number: "4164978855760477", Type: SuccessIfUnder100},

			{Number: "3333333333479407", Type: PreAuthorizationPossible},
			{Number: "3333333333334503", Type: PreAuthorizationPossible},
			{Number: "5117963095204135", Type: PreAuthorizationPossible},
			{Number: "4341500505113562", Type: PreAuthorizationPossible},

			// Failure cards
			{Number: "3333333333333349", Type: FailureRegardlessOfAmount},
			{Number: "3333333333336409", Type: FailureRegardlessOfAmount},
			{Number: "3333333333339205", Type: FailureRegardlessOfAmount},
			{Number: "3333333333338710", Type: FailureRegardlessOfAmount},
			{Number: "3333333333337605", Type: FailureRegardlessOfAmount},
			{Number: "3333333333337340", Type: FailureRegardlessOfAmount},
			{Number: "3333333333339403", Type: FailureRegardlessOfAmount},
			{Number: "3333333333337720", Type: FailureRegardlessOfAmount},
			{Number: "3333333333335120", Type: FailureRegardlessOfAmount},
			{Number: "3333333333335930", Type: FailureRegardlessOfAmount},
			{Number: "4280596505234682", Type: FailureRegardlessOfAmount},
			{Number: "5218572540397762", Type: FailureRegardlessOfAmount},
		},
	}
}

func (s *SandboxSimulator) SimulatePayment(cardNumber string, amount float64) (PaymentOutcome, error) {
	for _, card := range s.Cards {
		if card.Number == cardNumber {
			switch card.Type {
			case SuccessRegardlessOfAmount, SuccessA2CPay:
				return PaymentSuccess, nil
			case SuccessIfUnder100:
				if amount <= 100 {
					return PaymentSuccess, nil
				}
				return PaymentFailure, errors.New("payment amount exceeds limit for success")
			case FailureRegardlessOfAmount, FailureRandomErrorA2CPay:
				return PaymentFailure, errors.New("simulated failure")
			case FailureInsufficientBalanceA2CPay:
				return PaymentFailure, errors.New("insufficient_balance")
			case InvalidSandboxPan:
				return PaymentInvalid, errors.New("invalid sandbox pan")
			}
		}
	}
	return PaymentInvalid, fmt.Errorf("card number %s not recognized", cardNumber)
}
