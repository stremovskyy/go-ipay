/*
 * MIT License
 *
 * Copyright (c) 2025 Anton Stremovskyy
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

package main

import (
	"gorm.io/gorm"
)

type IpayRequestRecord struct {
	gorm.Model
	Kind          string  `gorm:"column:kind;size:32;not null"`
	CorrelationID string  `gorm:"column:correlation_id;size:255;not null"`
	ReferenceID   *string `gorm:"column:ref_id"`
	Body          []byte  `gorm:"column:body;type:blob;not null"`
}

func (IpayRequestRecord) TableName() string { return "custom_records" }

func (r *IpayRequestRecord) GetID() uint           { return r.ID }
func (r *IpayRequestRecord) GetType() string       { return r.Kind }
func (r *IpayRequestRecord) SetType(v string)      { r.Kind = v }
func (r *IpayRequestRecord) GetRequestID() string  { return r.CorrelationID }
func (r *IpayRequestRecord) SetRequestID(v string) { r.CorrelationID = v }
func (r *IpayRequestRecord) SetPrimaryID(v *string) {
	if v == nil {
		r.ReferenceID = nil
		return
	}
	tmp := *v
	r.ReferenceID = &tmp
}
func (r *IpayRequestRecord) SetPayload(data []byte) { r.Body = append(r.Body[:0], data...) }
func (r *IpayRequestRecord) GetPayload() []byte     { return r.Body }

type IpayRequestTag struct {
	ID       uint   `gorm:"primaryKey"`
	RecordID uint   `gorm:"column:record_ref;index"`
	Key      string `gorm:"column:tag_key"`
	Value    string `gorm:"column:tag_value"`
}

func (IpayRequestTag) TableName() string { return "custom_tags" }

func (t *IpayRequestTag) SetRecordID(id uint) { t.RecordID = id }
func (t *IpayRequestTag) SetKey(k string)     { t.Key = k }
func (t *IpayRequestTag) SetValue(v string)   { t.Value = v }
