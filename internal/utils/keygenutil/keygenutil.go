package keygenutil 

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func GenerateKeyPair(bitSize int) ([]byte, []byte, error) {
    privateKey, err := generatePrivateKey(bitSize)
    if err != nil {
        return nil, nil, err
    }

    publicKeyByte, err := generatePublicKey(privateKey)
    if err != nil {
        return nil, nil, err
    }

    privateKeyByte := encodePrivateKeyToPEM(privateKey)

    return privateKeyByte, publicKeyByte, nil
}

func GenerateRandomKey(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func generatePublicKey(privatekey *rsa.PrivateKey) ([]byte, error) {
	// Get ASN.1 DER format
	pubDER, err := x509.MarshalPKIXPublicKey(&privatekey.PublicKey)
    if err != nil {
        return nil, err
    }

	// pem.Block
	pubBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDER,
	}

	// Private key in PEM format
	publicPEM := pem.EncodeToMemory(&pubBlock)

	return publicPEM, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}
