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
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"github.com/stremovskyy/go-ipay/log"
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
