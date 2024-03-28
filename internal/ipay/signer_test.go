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

import (
	"testing"
)

func TestNewSigner(t *testing.T) {
	key := "testkey"
	got := NewSigner(key)
	if got == nil {
		t.Errorf("NewSigner() returned nil, expected a Signer instance")
	}
}

func Test_hashHmacSha512(t *testing.T) {
	args := struct {
		data string
		key  string
	}{
		data: "hello",
		key:  "key",
	}
	want := "ff06ab36757777815c008d32c8e14a705b4e7bf310351a06a23b612dc4c7433e7757d20525a5593b71020ea2ee162d2311b247e9855862b270122419652c0c92"
	if got := hashHmacSha512(args.data, args.key); got != want {
		t.Errorf("hashHmacSha512() = %v, want %v", got, want)
	}
}

func Test_sha1string(t *testing.T) {
	args := struct {
		data int64
	}{
		data: 1234567890,
	}
	want := "01b307acba4f54f55aafc33bb06bbbf6ca803e9a"
	if got := sha1string(args.data); got != want {
		t.Errorf("sha1string() = %v, want %v", got, want)
	}
}

func Test_signer_Sign(t *testing.T) {
	s := NewSigner("testkey")
	got := s.Sign("testkey")
	if got == nil || got.Sign == "" {
		t.Errorf("Sign() returned an empty string, expected a non-empty signature")
	}
}
