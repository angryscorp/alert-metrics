package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type PrivateKeyDecrypter struct {
	privateKey *rsa.PrivateKey
}

func NewPrivateKeyDecrypter(keyPath string) (*PrivateKeyDecrypter, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var privateKey *rsa.PrivateKey

	if block.Type == "PRIVATE KEY" {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not RSA private key")
		}
	} else if block.Type == "RSA PRIVATE KEY" {
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
		}
		privateKey = key
	} else {
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	return &PrivateKeyDecrypter{privateKey: privateKey}, nil
}

func (d *PrivateKeyDecrypter) Decrypt(data []byte) ([]byte, error) {
	hash := sha256.New()
	decrypted, err := rsa.DecryptOAEP(hash, rand.Reader, d.privateKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	return decrypted, nil
}

type NoOpDecrypter struct{}

func NewNoOpDecrypter() *NoOpDecrypter {
	return &NoOpDecrypter{}
}

func (n *NoOpDecrypter) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}
