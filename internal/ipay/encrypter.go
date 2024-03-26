package ipay

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"github.com/megakit-pro/go-ipay/internal/log"
)

type Cipher interface {
	EncryptData(rawData string) (string, error)
}

type encrypter struct {
	key    string
	logger *log.Logger
}

func NewCipher(key string) Cipher {
	return &encrypter{key: key, logger: log.NewLogger("cipher")}
}

func (c *encrypter) EncryptData(rawData string) (string, error) {
	if rawData == "" {
		return "", fmt.Errorf("data to encrypt is empty")
	}

	if c.key == "" {
		return "", fmt.Errorf("key is empty")
	}

	// Convert the key to a SHA-512 hash to ensure it's 64 bytes and then truncate to 32 bytes for AES-256
	keyHash := sha512.Sum512([]byte(c.key))
	key := keyHash[:32]

	iv := key[:12]

	c.logger.Debug("Encrypting data: %s", rawData)
	c.logger.Debug("Key: %x", key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create new cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create new GCM: %w", err)
	}

	// Encrypt the data
	ciphertext := aesgcm.Seal(nil, iv, []byte(rawData), nil)
	// Separate the tag from the ciphertext
	tag := ciphertext[len(ciphertext)-aesgcm.Overhead():]
	encData := ciphertext[:len(ciphertext)-aesgcm.Overhead()]

	c.logger.Debug("Encrypted data: %x", encData)

	// Encode the tag to base64 for compatibility
	tagBase64 := base64.StdEncoding.EncodeToString(tag)

	c.logger.Debug("Tag: %x", tag)

	// Return the encoded data and the tag concatenated, similar to the PHP version
	return base64.StdEncoding.EncodeToString(encData) + "." + tagBase64, nil
}
