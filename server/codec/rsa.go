package codec

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
)

type RSACodec struct {
	key *rsa.PrivateKey
}

func NewRSACodec(keyPath string) (*RSACodec, error) {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &RSACodec{key: privateKey}, nil
}

func (c *RSACodec) Decode(data []byte) ([]byte, error) {
	cipher, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	plain := make([]byte, 0)
	i := 0

	for i < len(cipher) {
		j := i + 256
		if j > len(cipher) {
			j = len(cipher)
		}

		p, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, c.key, cipher[i:j], []byte(""))
		if err != nil {
			return nil, err
		}
		i += 256
		plain = append(plain, p...)
	}
	return plain, nil
}

func (c *RSACodec) Encode(data []byte) ([]byte, error) {
	data, err := EncryptWithPrivateKey(c.key, data)
	if err != nil {
		return nil, err
	}

	return []byte(base64.StdEncoding.EncodeToString(data)), nil
}
