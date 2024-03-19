package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Dealer struct {
	key    [32]byte
	aesgcm cipher.AEAD
}

func NewDealer(k string) (*Dealer, error) {
	key := sha256.Sum256([]byte(k))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}
	return &Dealer{
		key:    key,
		aesgcm: aesgcm,
	}, nil
}

func (d Dealer) Encrypt(msg string) (string, error) {
	nonce := d.key[len(d.key)-d.aesgcm.NonceSize():]

	dst := d.aesgcm.Seal(nil, nonce, []byte(msg), nil) // зашифровываем
	return hex.EncodeToString(dst), nil
}

func (d Dealer) Decrypt(msg string) (string, error) {
	nonce := d.key[len(d.key)-d.aesgcm.NonceSize():]

	encrypted, err := hex.DecodeString(msg)
	if err != nil {
		return "", fmt.Errorf("hex.DecodeString: %w", err)
	}

	decrypted, err := d.aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", fmt.Errorf("aesgcm.Open: %w", err)
	}
	return string(decrypted), nil
}
