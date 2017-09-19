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
	KeySize uint = 32

	MijinId   = byte(0x60)
	MainnetId = byte(0x68)
	TestnetId = byte(0x98)
)

type KeyPair struct {
	public  []byte
	private []byte
}

func IsValidChainId(id byte) bool {
	for _, i := range []byte{MainnetId, TestnetId, MijinId} {
		if i == id {
			return true
		}
	}
	return false
}

func ToChainId(ch string) (byte, error) {
	switch ch {
	case "mainnet", "main", "0x68", "68":
		return MainnetId, nil
	case "testnet", "test", "0x98", "98":
		return TestnetId, nil
	case "mijin", "0x60", "60":
		return MijinId, nil
	}

	return 0, errors.New("invalid chain id")
}

func NewKeyPair() (KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(nil)

	return KeyPair{
		pub, priv[:KeySize],
	}, err
}

func ToAccount(pub []byte, chainId byte) (string, error) {
	if !IsValidChainId(chainId) {
		return "", errors.New("invalid chain id")
	}

	h := sha3.SumKeccak256(pub)
	//fmt.Printf("SHA3 %x\n", h)

	md := ripemd160.New()
	_, err := md.Write(h[:])
	if err != nil {
		return "", err
	}

	s := md.Sum(nil)
	//fmt.Printf("Ripemd %x\n", s)

	s = append([]byte{chainId}, s...)
	h = sha3.SumKeccak256(s)
	address := append(s, h[:4]...)
	//fmt.Printf("Address %x\n", address)

	return base32.StdEncoding.EncodeToString(address), nil
}

func GenerateAccount(chainId byte) (string, error) {
	keyPair, err := NewKeyPair()
	if err != nil {
		os.Exit(-1)
	}
	//fmt.Printf("Public: %x, Private: %x\n", keyPair.public, keyPair.private)
	return ToAccount(keyPair.public, chainId)
}
