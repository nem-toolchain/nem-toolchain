// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/sha3"
)

type EncryptedKey struct {
	ciphertext []byte
	iv         []byte
}

func encrypt(key []byte, password []byte) (EncryptedKey, error) {
	var encrypted EncryptedKey
	block, err := aes.NewCipher(password)

	if err != nil {
		return encrypted, err
	}

	ciphertext := make([]byte, len(key))
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return encrypted, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, key)

	encrypted = EncryptedKey{
		ciphertext: ciphertext,
		iv:         iv,
	}

	return encrypted, nil
}

func decrypt(encrypted EncryptedKey, password []byte) ([]byte, error) {
	pk := make([]byte, 32)
	block, err := aes.NewCipher(password)

	if err != nil {
		return pk, err
	}

	mode := cipher.NewCBCDecrypter(block, encrypted.iv)
	mode.CryptBlocks(pk, encrypted.ciphertext)

	return pk, nil
}

func derive(pass []byte) ([]byte, error) {
	hash := sha3.NewKeccak256()
	h := pass

	for i := 0; i < 20; i++ {
		_, err := hash.Write(h)
		if err != nil {
			return h, err
		}

		h = hash.Sum(nil)
		hash.Reset()
	}

	return h, nil
}
