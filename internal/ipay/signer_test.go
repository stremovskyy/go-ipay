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
	"strings"
	"testing"
)

func TestSHA3512Compatibility(t *testing.T) {
	// SHA3-512 vectors also cover UTF-8 keys and input longer than one hash block.
	for _, tt := range []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"},
		{"abc", "abc", "b751850b1a57168a5693cd924b6b096e08f621827444f70d884f5d0240d2712e10e116e9192af3c91a7ec57647e3934057340b4cf408d5a56592f8274eec53f0"},
		{"unicode", "2026-10-19 18:00:00ключ-Дія", "f4af4089c731b5976224f4999939fbcad2963ba6e592a42a3c5af59c07945b28398b04e3ffe433527196844192b84003bac421192e04afd891f4600cb7c2806e"},
		{"multiple blocks", strings.Repeat("a", 200), "eae6c85c6904f11075de9f9d5e1064371d000510fa3d2d79d40cf9be34892fb01859d0a0234e138bcb0ad5c84f6c0dca226a414b0c9a2897cb695f5185fe36ec"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := sha3512(tt.input); got != tt.want {
				t.Errorf("SHA3-512 = %q, want %q", got, tt.want)
			}
		})
	}
}

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
