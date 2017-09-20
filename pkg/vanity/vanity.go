// Package vanity implements a bundle of vanity address generators
package vanity

import (
	"encoding/base32"
	"os"

	"errors"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Supported chains
const (
	MijinId   = byte(0x60)
	MainnetId = byte(0x68)
	TestnetId = byte(0x98)
)

// Account keypair
type KeyPair struct {
	public  []byte
	private []byte
}

// ErrInvalidChain indicates invalid chain id.
var ErrInvalidChain = errors.New("invalid chain id")

// GenAddress generates a new address for required chain on crypto random basis.
// It’s a run-time error for unknown chain.
func GenAddress(chainId byte) (string, error) {
	pair, err := NewKeyPair()
	if err != nil {
		os.Exit(-1)
	}
	return ToAddress(pair.public, chainId)
}

// NewKeyPair generates a public/private key pair using entropy from crypto rand.
func NewKeyPair() (KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	return KeyPair{
		pub, priv[:32],
	}, err
}

// ToAddress converts public key to public account address.
// It’s a run-time error for unknown chain.
func ToAddress(pubKey []byte, chainId byte) (string, error) {
	if !IsValidChainId(chainId) {
		return "", ErrInvalidChain
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

// IsValidChainId checks chain id for existence
func IsValidChainId(id byte) bool {
	for _, i := range []byte{MijinId, MainnetId, TestnetId} {
		if i == id {
			return true
		}
	}
	return false
}
