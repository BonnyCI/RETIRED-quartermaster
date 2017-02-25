package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"sync"
)

type RSAKey struct {
	key *rsa.PrivateKey
}

func (m *RSAKey) Generate() error {
	var err error
	reader := rand.Reader

	m.key, err = rsa.GenerateKey(reader, KeyBits)
	if err != nil {
		return err
	}

	return nil
}

func (m *RSAKey) Key() *rsa.PrivateKey {
	return m.key
}

func (m *RSAKey) Private() ([]byte, error) {
	if m.key == nil {
		err := m.Generate()
		if err != nil {
			return nil, err
		}
	}

	block := PrivateBlock
	block.Bytes = x509.MarshalPKCS1PrivateKey(m.key)

	return pem.EncodeToMemory(&block), nil
}

func (m *RSAKey) Public() ([]byte, error) {
	if m.key == nil {
		err := m.Generate()
		if err != nil {
			return nil, err
		}
	}

	pubDer, err := x509.MarshalPKIXPublicKey(&m.key.PublicKey)
	if err != nil {
		return nil, err
	}

	block := PublicBlock
	block.Bytes = pubDer
	return pem.EncodeToMemory(&block), nil
}

func (m *RSAKey) Write(priv io.Writer, pub io.Writer) error {
	if priv != nil {
		privKey, err := m.Private()
		if err != nil {
			return err
		}

		fmt.Fprint(priv, string(privKey))
	}

	if pub != nil {
		pubKey, err := m.Public()
		if err != nil {
			return err
		}

		fmt.Fprint(pub, string(pubKey))
	}
	return nil
}

func (m *RSAKey) Read(priv []byte, pub []byte) error {
	var err error

	if priv != nil {
		var block *pem.Block
		if block, _ := pem.Decode(priv); block == nil {
			return ErrMustBePEMEncoded
		}

		var key interface{}
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return err
		}

		var ok bool
		if m.key, ok = key.(*rsa.PrivateKey); !ok {
			return ErrNotRSAPrivateKey
		}
	}

	if pub != nil {
		var block *pem.Block
		if block, _ := pem.Decode(priv); block == nil {
			return ErrMustBePEMEncoded
		}

		var key interface{}
		if key, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			return err
		}

		var pubKey *rsa.PublicKey
		var ok bool
		if pubKey, ok = key.(*rsa.PublicKey); !ok {
			return ErrNotRSAPublicKey
		}

		m.key.PublicKey = *pubKey
	}
	return nil
}

var instance *RSAKey
var keySync sync.Once

func GetRSAKeySingleton() *RSAKey {
	keySync.Do(func() {
		instance = &RSAKey{}
		instance.Generate()
	})
	return instance
}

func RSASingletonPrivate() []byte {
	key, _ := GetRSAKeySingleton().Private()
	return key
}

func RSASingletonPublic() []byte {
	key, _ := GetRSAKeySingleton().Public()
	return key
}
