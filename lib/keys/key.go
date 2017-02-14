package keys

import (
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"io"
)

var (
	KeyBits = 2048

	PrivateBlock = pem.Block{
		Headers: nil,
		Type:    "RSA PRIVATE KEY",
	}

	PublicBlock = pem.Block{
		Headers: nil,
		Type:    "PUBLIC KEY",
	}

	ErrMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM Encoded PKCS1.")
	ErrNotRSAPrivateKey = errors.New("Key is not a valid RSA private key.")
	ErrNotRSAPublicKey  = errors.New("Key is not a valid RSA public key.")
)

type Key interface {
	Generate() error
	Key() *rsa.PrivateKey
	Private() ([]byte, error)
	Public() ([]byte, error)
	Write(io.Writer, io.Writer) error
	Read([]byte, []byte) error
}
