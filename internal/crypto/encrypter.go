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

type PublicKeyEncrypter struct {
	publicKey *rsa.PublicKey
}

func NewPublicKeyEncrypter(keyPath string) (*PublicKeyEncrypter, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var publicKey *rsa.PublicKey

	if block.Type == "PUBLIC KEY" {
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}
		var ok bool
		publicKey, ok = pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("key is not RSA public key")
		}
	} else if block.Type == "RSA PUBLIC KEY" {
		pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
		}
		publicKey = pub
	} else {
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	return &PublicKeyEncrypter{publicKey: publicKey}, nil
}

func (e *PublicKeyEncrypter) Encrypt(data []byte) ([]byte, error) {
	hash := sha256.New()
	encrypted, err := rsa.EncryptOAEP(hash, rand.Reader, e.publicKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}
	return encrypted, nil
}
