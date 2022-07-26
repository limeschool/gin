package gin

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

type ExtRsaConfig struct {
	Enable bool   `json:"enable" mapstructure:"enable"`
	Name   string `json:"name"  mapstructure:"name"`
	Path   string `json:"path" mapstructure:"path"`
}

type ExtRsa struct {
	key string
}

func initRsa() {
	var confList []ExtRsaConfig
	_ = globalConfig.UnmarshalKey("rsa", &confList)
	rsa := make(map[string]*ExtRsa)
	for _, item := range confList {
		if !item.Enable {
			return
		}
		file, err := os.Open(item.Path)
		if err != nil {
			panic("rsa init fail:" + err.Error())
		}
		key, _ := ioutil.ReadAll(file)
		rsa[item.Name] = &ExtRsa{
			key: string(key),
		}
	}
	globalRsa = rsa
}

func (e ExtRsa) Encode(plainText string) (string, error) {
	block, _ := pem.Decode([]byte(e.key))
	if block == nil {
		return "", errors.New("rsa public key error")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	key := publicKeyInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(plainText))
	return base64.StdEncoding.EncodeToString(cipherText), err
}

func (e ExtRsa) Decode(cipherText string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(e.key))
	if block == nil {
		return "", errors.New("rsa private key error")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, key, text)
	return string(plainText), err
}
