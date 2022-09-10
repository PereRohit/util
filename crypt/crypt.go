package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
)

func RsaOaepEncrypt(secretMessage, encryptKey string) (string, error) {
	key, err := ImportPubKeyFromPEMStr(encryptKey)
	if err != nil {
		return "", fmt.Errorf("import public key failed: %s", err.Error())
	}
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, key, []byte(secretMessage), label)
	if err != nil {
		return "", fmt.Errorf("encrypt failed: %s", err.Error())
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RsaOaepDecrypt(cipherText string, privKey rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	return string(plaintext), err
}

func ImportPubKeyFromPEMStr(pubKey string) (*rsa.PublicKey, error) {
	key, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	spkiBlock, _ := pem.Decode(key)
	if spkiBlock == nil || spkiBlock.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	return x509.ParsePKCS1PublicKey(spkiBlock.Bytes)
}

func ImportPriKeyFromPEMStr(priKey string) (*rsa.PrivateKey, error) {
	key, err := base64.StdEncoding.DecodeString(priKey)
	if err != nil {
		return nil, err
	}
	spkiBlock, _ := pem.Decode(key)
	if spkiBlock == nil || spkiBlock.Type != "RSA PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(spkiBlock.Bytes)
}

func ExportPubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return base64.StdEncoding.EncodeToString([]byte(pubKeyPem))
}

func ExportPriKeyAsPEMStr(privKey *rsa.PrivateKey) string {
	priKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	))
	return base64.StdEncoding.EncodeToString([]byte(priKeyPem))
}
