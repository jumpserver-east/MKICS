package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

func DecryptPassword(cipherText string, privateKey string) (string, error) {
	parts := strings.Split(cipherText, ":")
	if len(parts) != 3 {
		return "", errors.New("cipherText format invalid")
	}

	aesKeyBase64, err := rsaDecrypt(parts[0], privateKey)
	if err != nil {
		return "", err
	}

	aesKey, err := base64.StdEncoding.DecodeString(aesKeyBase64)
	if err != nil {
		return "", err
	}

	iv, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	cipherData, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return "", err
	}

	return aesDecryptCBC(cipherData, []byte(aesKey), iv)
}

func GenerateRSAKeyPair() (privateKeyPEM string, publicKeyPEM string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	privASN1 := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privASN1,
	})
	privateKeyPEM = string(privPEM)

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	publicKeyPEM = string(pubPEM)

	return
}

func rsaDecrypt(cipherTextBase64, privateKeyPEM string) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("private key decode error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherData)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func aesDecryptCBC(cipherText []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipherText is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	plainText = pkcs7Unpad(plainText)
	return string(plainText), nil
}

func pkcs7Unpad(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return src
	}
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
