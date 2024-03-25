package ipay

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/megakit-pro/go-ipay/internal/log"
)

type Signer interface {
	Sign(key string) *Sign
}

type signer struct {
	logger *log.Logger
}

func (s *signer) Sign(key string) *Sign {
	timeNow := time.Now().UnixNano()
	salt := sha1string(timeNow)

	s.logger.Debug("Signing data: %d", timeNow)

	sign := hashHmacSha512(salt, key)

	return &Sign{
		Salt: salt,
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

	return fmt.Sprintf("%x", sha1.Sum([]byte(dataStr)))
}
