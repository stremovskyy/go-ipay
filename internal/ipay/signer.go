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
	"crypto/hmac"
	"crypto/sha1" // #nosec G505 -- iPay requires legacy SHA-1 signatures for salts.
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/stremovskyy/go-ipay/log"
)

type Signer interface {
	Sign(key string) *Sign
	MobileSign(key string) *MobileSign
}

type signer struct {
	logger *log.Logger
}

func (s *signer) MobileSign(key string) *MobileSign {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	dataString := fmt.Sprintf("%s%s", timeNow, key)

	sign := sha3512(dataString)

	return &MobileSign{
		Time: &timeNow,
		Sign: sign,
	}
}

func sha3512(data string) string {
	hasher := sha3.New512()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *signer) Sign(key string) *Sign {
	timeNow := time.Now().UnixNano()
	salt := sha1string(timeNow)

	s.logger.Debug("Signing data: %d", timeNow)

	sign := hashHmacSha512(salt, key)

	return &Sign{
		Salt: &salt,
		Sign: sign,
	}
}

func NewSigner(key string) Signer {
	return &signer{}
}

func hashHmacSha512(data string, key string) string {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(data))

	return fmt.Sprintf("%x", mac.Sum(nil))
}

func sha1string(data int64) string {
	dataStr := fmt.Sprintf("%d", data)

	// #nosec G401 -- iPay requires SHA-1 hashing for backward-compatible signatures.
	return fmt.Sprintf("%x", sha1.Sum([]byte(dataStr)))
}
