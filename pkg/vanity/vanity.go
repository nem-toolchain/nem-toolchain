package vanity

import (
	"encoding/base32"
	"os"

	"errors"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

const (
	MijinId   = byte(0x60)
	MainnetId = byte(0x68)
	TestnetId = byte(0x98)
)

type KeyPair struct {
	public  []byte
	private []byte
}

func GenAddress(chainId byte) (string, error) {
	pair, err := NewKeyPair()
	if err != nil {
		os.Exit(-1)
	}
	return ToAddress(pair.public, chainId)
}

func ToAddress(pubKey []byte, chainId byte) (string, error) {
	if !IsValidChainId(chainId) {
		return "", errors.New("Invalid chain id")
	}
	h := sha3.SumKeccak256(pubKey)
	r := ripemd160.New()
	_, err := r.Write(h[:])
	if err != nil {
		return "", err
	}
	b := append([]byte{chainId}, r.Sum(nil)...)
	h = sha3.SumKeccak256(b)
	a := append(b, h[:4]...)
	return base32.StdEncoding.EncodeToString(a), nil
}

func NewKeyPair() (KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	return KeyPair{
		pub, priv[:32],
	}, err
}

func IsValidChainId(id byte) bool {
	for _, i := range []byte{MijinId, MainnetId, TestnetId} {
		if i == id {
			return true
		}
	}
	return false
}
