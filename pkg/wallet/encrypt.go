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

func deriveKey(pass string) []byte {
	k := []byte(pass)
	for i := 0; i < 20; i++ {
		h := sha3.SumKeccak256(k)
		k = h[:]
	}
	return k
}

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

func padData(src []byte, blockSize int) []byte {
	p := blockSize - len(src)%blockSize
	s := bytes.Repeat([]byte{byte(p)}, p)
	return append(src, s...)
}

func unpadData(src []byte, blockSize int) ([]byte, error) {
	l := len(src)
	if l < blockSize {
		return nil,
			fmt.Errorf("wallet: insufficient slice length for unpadding, "+
				"should be minimal %v, but got %v", blockSize, l)
	}
	p := int(src[l-1])
	if p > blockSize {
		return nil, fmt.Errorf("wallet: invalid padding size, "+
			"should be maximum %v, but got %v", blockSize, p)
	}
	if !bytes.HasSuffix(src, bytes.Repeat([]byte{byte(p)}, p)) {
		return nil,
			fmt.Errorf("wallet: non-uniform padding - %v",
				hex.EncodeToString(src[l-p:]))
	}
	return src[:l-p], nil
}
