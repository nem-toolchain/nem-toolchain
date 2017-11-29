// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package wallet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// deriveKey performs pseudo-HMAC procedure to derive 32-byte key
// from arbitrary password in string representation.
func deriveKey(pass string) []byte {
	k := []byte(pass)
	for i := 0; i < 20; i++ {
		h := sha3.SumKeccak256(k)
		k = h[:]
	}
	return k
}

// encryptData encrypts arbitrary data with AES (Rijndael) cipher in CBC mode
// with PKCS7 padding and a 128-bit block size (AES-128/CBC/PKCS7Padding).
func encryptData(key, iv, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	src := padData(data, aes.BlockSize)
	enc := make([]byte, len(src))
	mode.CryptBlocks(enc, src)
	return enc
}

// decryptData decrypts data encrypted by `encryptData` method.
func decryptData(key, iv, enc []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(enc))
	mode.CryptBlocks(dst, enc)
	return unpadData(dst, aes.BlockSize)
}

// padData performs PKCS7 data padding for arbitrary block size before block encryption.
func padData(src []byte, blockSize int) []byte {
	p := blockSize - len(src)%blockSize
	s := bytes.Repeat([]byte{byte(p)}, p)
	return append(src, s...)
}

// unpadData performs PKCS7 data unpadding for arbitrary block size after block decryption.
func unpadData(src []byte, blockSize int) ([]byte, error) {
	l := len(src)
	if l == 0 || l%blockSize != 0 {
		return nil, &dataLengthError{l, blockSize}
	}
	p := int(src[l-1])
	if p > blockSize {
		return nil, &paddingLengthError{p, blockSize}
	}
	if !bytes.HasSuffix(src, bytes.Repeat([]byte{byte(p)}, p)) {
		return nil, &invalidPaddingError{src[l-p:]}
	}
	return src[:l-p], nil
}

type dataLengthError struct {
	len, blockSize int
}

func (e *dataLengthError) Error() string {
	return fmt.Sprintf("wallet: invalid length for unpadding, "+
		"should be a not zero multiple of the block size %v, but got %v", e.blockSize, e.len)
}

type paddingLengthError struct {
	len, blockSize int
}

func (e *paddingLengthError) Error() string {
	return fmt.Sprintf("wallet: invalid padding size, "+
		"should be maximum block size %v, but got %v", e.blockSize, e.len)
}

type invalidPaddingError struct {
	pad []byte
}

func (e *invalidPaddingError) Error() string {
	return fmt.Sprintf("wallet: invalid non-uniform padding - %v",
		hex.EncodeToString(e.pad))
}
