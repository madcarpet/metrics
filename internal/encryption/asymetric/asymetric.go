// Package asymetric - contains functions to work with asymetric encryption.
package asymetric

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// EncryptWithPubKey - encrypts data with the given public key.
func EncryptWithPubKey(keyPath string, data []byte) ([]byte, error) {
	pubKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("public key reading error: %s", err)
	}
	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("public key decoding error, key not found at path %s: %s", keyPath, err)
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("public key parsing failed: %s", err)
	}
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type *rsa.PublicKey")
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, fmt.Errorf("data encrypting with pub key error: %s", err)
	}
	return encryptedData, nil
}

func DecryptWithPrivateKey(keyPath string, encdata []byte) ([]byte, error) {
	privKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("private key reading error: %s", err)
	}
	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("private key decoding error, key not found at path %s: %s", keyPath, err)
	}
	privKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("private key parsing failed: %s", err)
	}
	privKey, ok := privKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type *rsa.PrivateKey")
	}
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, encdata)
	if err != nil {
		return nil, fmt.Errorf("data decrypting with priv key error: %s", err)
	}
	return decryptedData, nil
}
