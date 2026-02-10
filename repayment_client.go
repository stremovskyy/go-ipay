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
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/stremovskyy/go-ipay/consts"
	repayinternal "github.com/stremovskyy/go-ipay/internal/repayment"
	"github.com/stremovskyy/go-ipay/repayment"
)

// RepaymentTransaction describes a single transaction row in the CreateRepayment CSV file.
// The Repayment API expects a semicolon-separated CSV with two columns: pmt_id;ext_id
type RepaymentTransaction struct {
	PmtID int64
	ExtID string
}

// CreateRepaymentRequest is the high-level request for Repayment API CreateRepayment action.
// Provide either TransactionsFilePath or Transactions.
type CreateRepaymentRequest struct {
	Merchant *Merchant

	// MchID is the merchant ID to debit for the repayment. If 0, Merchant.MerchantID is used.
	MchID int64
	// ExtID is the repayment request ID on merchant side (max length: 50).
	ExtID string
	// SmchID is optional.
	SmchID *int64

	// TransactionsFilePath is a path to the CSV file (multipart field name: file).
	TransactionsFilePath string
	// Transactions, when provided, are encoded into a CSV file and uploaded as "transactions.csv".
	Transactions []RepaymentTransaction
}

// CancelRepaymentRequest is the high-level request for Repayment API CancelRepayment action.
// Provide either RepaymentGUID or ExtID.
type CancelRepaymentRequest struct {
	Merchant *Merchant

	RepaymentGUID *string
	ExtID         *string
}

// GetRepaymentStatusRequest is the high-level request for Repayment API GetRepaymentStatus action.
// Provide either RepaymentGUID or ExtID.
type GetRepaymentStatusRequest struct {
	Merchant *Merchant

	RepaymentGUID *string
	ExtID         *string
}

// GetRepaymentProcessingFileRequest is the high-level request for Repayment API GetRepaymentProcessingFile action.
// Provide either RepaymentGUID or ExtID.
type GetRepaymentProcessingFileRequest struct {
	Merchant *Merchant

	RepaymentGUID *string
	ExtID         *string
}

// CreateRepayment creates a repayment request via Repayment API.
func (c *client) CreateRepayment(request *CreateRepaymentRequest, runOpts ...RunOption) (*repayment.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}
	if request.Merchant == nil {
		return nil, ErrMerchantIsNil
	}

	opts := collectRunOptions(runOpts)

	if request.Merchant.Login == "" {
		return nil, fmt.Errorf("create repayment: merchant login is empty")
	}
	if request.Merchant.RepaymentKey == "" {
		return nil, fmt.Errorf("create repayment: merchant repayment key is empty")
	}

	mchID := request.MchID
	if mchID == 0 {
		mid := request.Merchant.GetMerchantID()
		if mid == nil {
			return nil, fmt.Errorf("create repayment: invalid merchant ID")
		}
		mchID = *mid
	}

	if request.ExtID == "" {
		return nil, fmt.Errorf("create repayment: ext_id is empty")
	}
	if len(request.ExtID) > 50 {
		return nil, fmt.Errorf("create repayment: ext_id is too long (max 50)")
	}

	if request.TransactionsFilePath != "" && len(request.Transactions) > 0 {
		return nil, fmt.Errorf("create repayment: provide either TransactionsFilePath or Transactions, not both")
	}

	var (
		fileName string
		file     *os.File
		reader   = bytes.NewReader([]byte(nil))
		txCount  int
		filePath string
	)

	if len(request.Transactions) > 0 {
		raw, err := encodeRepaymentTransactionsCSV(request.Transactions)
		if err != nil {
			return nil, fmt.Errorf("create repayment: encode transactions CSV: %w", err)
		}
		reader = bytes.NewReader(raw)
		fileName = "transactions.csv"
		txCount = len(request.Transactions)
	} else {
		if request.TransactionsFilePath == "" {
			return nil, fmt.Errorf("create repayment: transactions file path is empty")
		}
		f, err := os.Open(request.TransactionsFilePath)
		if err != nil {
			return nil, fmt.Errorf("create repayment: open transactions file: %w", err)
		}
		defer func() { _ = f.Close() }()

		file = f
		fileName = filepath.Base(request.TransactionsFilePath)
		filePath = request.TransactionsFilePath
	}

	timeString := time.Now().Format("2006-01-02 15:04:05")
	sign := repayinternal.SignSHA256Hex(timeString, request.Merchant.RepaymentKey)
	extID := request.ExtID

	repaymentRequest := &repayment.RequestWrapper{
		Request: repayment.Request{
			Auth: repayment.Auth{
				Login: request.Merchant.Login,
				Time:  timeString,
				Sign:  sign,
			},
			Action: repayment.ActionCreateRepayment,
			Body: repayment.Body{
				MchID:  &mchID,
				ExtID:  &extID,
				SmchID: request.SmchID,
			},
		},
		Operation: consts.CreateRepayment,
	}

	if opts.isDryRun() {
		payload := struct {
			Operation         string            `json:"operation"`
			Request           repayment.Request `json:"request"`
			TransactionsFile  string            `json:"transactions_file,omitempty"`
			TransactionsCount int               `json:"transactions_count,omitempty"`
			FileName          string            `json:"file_name"`
			MchID             int64             `json:"mch_id"`
			ExtID             string            `json:"ext_id"`
			SmchID            *int64            `json:"smch_id,omitempty"`
			AuthLogin         string            `json:"auth_login"`
			AuthTime          string            `json:"auth_time"`
			AuthSign          string            `json:"auth_sign"`
		}{
			Operation:         repaymentRequest.Operation,
			Request:           repaymentRequest.Request,
			TransactionsFile:  filePath,
			TransactionsCount: txCount,
			FileName:          fileName,
			MchID:             mchID,
			ExtID:             request.ExtID,
			SmchID:            request.SmchID,
			AuthLogin:         request.Merchant.Login,
			AuthTime:          timeString,
			AuthSign:          sign,
		}

		opts.handleDryRun(consts.RepaymentUrl, payload)
		return nil, nil
	}

	var apiFileReader io.Reader = reader
	if file != nil {
		apiFileReader = file
	}

	resp, err := c.ipayClient.RepaymentApi(repaymentRequest, fileName, apiFileReader)
	if err != nil {
		return resp, fmt.Errorf("create repayment API call: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("create repayment: empty response from API")
	}

	return resp, nil
}

// CancelRepayment cancels a repayment (allowed only on the creation day, per API docs).
func (c *client) CancelRepayment(request *CancelRepaymentRequest, runOpts ...RunOption) (*repayment.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}
	if request.Merchant == nil {
		return nil, ErrMerchantIsNil
	}

	opts := collectRunOptions(runOpts)

	auth, err := repaymentAuth(request.Merchant)
	if err != nil {
		return nil, fmt.Errorf("cancel repayment: %w", err)
	}

	repaymentGUID, extID, err := validateRepaymentLookup(request.RepaymentGUID, request.ExtID)
	if err != nil {
		return nil, fmt.Errorf("cancel repayment: %w", err)
	}

	wrapper := &repayment.RequestWrapper{
		Request: repayment.Request{
			Auth:   auth,
			Action: repayment.ActionCancelRepayment,
			Body: repayment.Body{
				RepaymentGUID: repaymentGUID,
				ExtID:         extID,
			},
		},
		Operation: consts.CancelRepayment,
	}

	if opts.isDryRun() {
		opts.handleDryRun(consts.RepaymentUrl, wrapper)
		return nil, nil
	}

	resp, err := c.ipayClient.RepaymentJSONApi(wrapper)
	if err != nil {
		return resp, fmt.Errorf("cancel repayment API call: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("cancel repayment: empty response from API")
	}

	return resp, nil
}

// GetRepaymentStatus returns the current repayment status.
func (c *client) GetRepaymentStatus(request *GetRepaymentStatusRequest, runOpts ...RunOption) (*repayment.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}
	if request.Merchant == nil {
		return nil, ErrMerchantIsNil
	}

	opts := collectRunOptions(runOpts)

	auth, err := repaymentAuth(request.Merchant)
	if err != nil {
		return nil, fmt.Errorf("get repayment status: %w", err)
	}

	repaymentGUID, extID, err := validateRepaymentLookup(request.RepaymentGUID, request.ExtID)
	if err != nil {
		return nil, fmt.Errorf("get repayment status: %w", err)
	}

	wrapper := &repayment.RequestWrapper{
		Request: repayment.Request{
			Auth:   auth,
			Action: repayment.ActionGetRepaymentStatus,
			Body: repayment.Body{
				RepaymentGUID: repaymentGUID,
				ExtID:         extID,
			},
		},
		Operation: consts.GetRepaymentStatus,
	}

	if opts.isDryRun() {
		opts.handleDryRun(consts.RepaymentUrl, wrapper)
		return nil, nil
	}

	resp, err := c.ipayClient.RepaymentJSONApi(wrapper)
	if err != nil {
		return resp, fmt.Errorf("get repayment status API call: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("get repayment status: empty response from API")
	}

	return resp, nil
}

// GetRepaymentProcessingFile returns the CSV processing results for a repayment.
func (c *client) GetRepaymentProcessingFile(request *GetRepaymentProcessingFileRequest, runOpts ...RunOption) ([]byte, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}
	if request.Merchant == nil {
		return nil, ErrMerchantIsNil
	}

	opts := collectRunOptions(runOpts)

	auth, err := repaymentAuth(request.Merchant)
	if err != nil {
		return nil, fmt.Errorf("get repayment processing file: %w", err)
	}

	repaymentGUID, extID, err := validateRepaymentLookup(request.RepaymentGUID, request.ExtID)
	if err != nil {
		return nil, fmt.Errorf("get repayment processing file: %w", err)
	}

	wrapper := &repayment.RequestWrapper{
		Request: repayment.Request{
			Auth:   auth,
			Action: repayment.ActionGetRepaymentProcessingFile,
			Body: repayment.Body{
				RepaymentGUID: repaymentGUID,
				ExtID:         extID,
			},
		},
		Operation: consts.GetRepaymentProcessingFile,
	}

	if opts.isDryRun() {
		opts.handleDryRun(consts.RepaymentUrl, wrapper)
		return nil, nil
	}

	raw, err := c.ipayClient.RepaymentProcessingFileApi(wrapper)
	if err != nil {
		return raw, fmt.Errorf("get repayment processing file API call: %w", err)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("get repayment processing file: empty response from API")
	}

	return raw, nil
}

func repaymentAuth(merchant *Merchant) (repayment.Auth, error) {
	if merchant == nil {
		return repayment.Auth{}, ErrMerchantIsNil
	}

	if merchant.Login == "" {
		return repayment.Auth{}, fmt.Errorf("merchant login is empty")
	}
	if merchant.RepaymentKey == "" {
		return repayment.Auth{}, fmt.Errorf("merchant repayment key is empty")
	}

	timeString := time.Now().Format("2006-01-02 15:04:05")
	sign := repayinternal.SignSHA256Hex(timeString, merchant.RepaymentKey)

	return repayment.Auth{
		Login: merchant.Login,
		Time:  timeString,
		Sign:  sign,
	}, nil
}

func validateRepaymentLookup(guid, extID *string) (*string, *string, error) {
	hasGUID := guid != nil && strings.TrimSpace(*guid) != ""
	hasExtID := extID != nil && strings.TrimSpace(*extID) != ""

	switch {
	case !hasGUID && !hasExtID:
		return nil, nil, fmt.Errorf("either repayment_guid or ext_id must be provided")
	case hasGUID && hasExtID:
		return nil, nil, fmt.Errorf("only one of repayment_guid or ext_id must be provided")
	}

	if hasExtID && len(*extID) > 50 {
		return nil, nil, fmt.Errorf("ext_id is too long (max 50)")
	}

	if hasGUID {
		return guid, nil, nil
	}

	return nil, extID, nil
}

func encodeRepaymentTransactionsCSV(transactions []RepaymentTransaction) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	w.Comma = ';'

	for i, tx := range transactions {
		if tx.PmtID == 0 {
			return nil, fmt.Errorf("transaction[%d]: pmt_id is empty", i)
		}
		if tx.ExtID == "" {
			return nil, fmt.Errorf("transaction[%d]: ext_id is empty", i)
		}
		if len(tx.ExtID) > 50 {
			return nil, fmt.Errorf("transaction[%d]: ext_id is too long (max 50)", i)
		}

		if err := w.Write(
			[]string{
				strconv.FormatInt(tx.PmtID, 10),
				tx.ExtID,
			},
		); err != nil {
			return nil, fmt.Errorf("transaction[%d]: write CSV: %w", i, err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("flush CSV: %w", err)
	}

	return buf.Bytes(), nil
}
