// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package account

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeriveKey(t *testing.T) {
	for k, v := range map[string]string{
		"":      "c63750c62c700d20b5364524cfde0a7c5077b17092edfdffe13705e646ed8aae",
		" ":     "3931c4b8958f4af4b855321333bf01603a888a31123de1a5b7824d7c0444a4de",
		"12345": "ad40ae7b162ebe4744ecba9e509648c0fad29fce4efe3a625dcce969e253722a",
	} {
		t.Run(k, func(t *testing.T) {
			d, _ := hex.DecodeString(v)
			assert.Equal(t, d, derive(k))
		})
	}
}

func TestEncryptDecryptData(t *testing.T) {
	for i, s := range []struct {
		key, iv, data, enc string
	}{
		{
			"c63750c62c700d20b5364524cfde0a7c5077b17092edfdffe13705e646ed8aae",
			"00000000000000000000000000000000",
			"0000000000000000000000000000000000000000000000000000000000000000",
			"90dcb3aabfe1df18dca68f94a2e2d38a16cec6c272b135da009c33e35c63dda43872a1e5eb5bf1d82ffde96df7196f62",
		},
		{
			"ad40ae7b162ebe4744ecba9e509648c0fad29fce4efe3a625dcce969e253722a",
			"846e971d8deb3d02f80876119cc30f43",
			"6945c1c5db2aba903a18d12d9c5401fdbdb6eec8d3807455856a1f98f83b5880",
			"14e980dcd411d27aed5431c7a4f9afbf036df2bfc3d9fdaeb23281c6d08902017c71c34cd5653defc38bc889be1b473a",
		},
	} {
		t.Run(string(i), func(t *testing.T) {
			key, _ := hex.DecodeString(s.key)
			iv, _ := hex.DecodeString(s.iv)
			data, _ := hex.DecodeString(s.data)
			enc, _ := hex.DecodeString(s.enc)

			assert.Equal(t, enc, encryptData(key, iv, data))

			act, err := decryptData(key, iv, enc)
			assert.NoError(t, err)
			assert.Equal(t, data, act)
		})
	}
}

func TestPadUnpadData(t *testing.T) {
	for i, s := range []struct {
		l       int
		in, out []byte
	}{
		{
			8,
			[]byte{},
			[]byte{8, 8, 8, 8, 8, 8, 8, 8},
		},
		{
			8,
			[]byte{1, 2, 3},
			[]byte{1, 2, 3, 5, 5, 5, 5, 5},
		},
		{
			8,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8,
				8, 8, 8, 8, 8, 8, 8, 8},
		},
		{
			16,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 8},
		},
		{
			16,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 4, 4, 4, 4},
		},
		{
			16,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
				16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16},
		},
	} {
		t.Run(string(i), func(t *testing.T) {
			assert.Equal(t, s.out, padData(s.in, s.l))

			act, err := unpadData(s.out, s.l)
			assert.NoError(t, err)
			assert.Equal(t, s.in, act)
		})
	}
}

func TestUnpadData_dataLengthError(t *testing.T) {
	for i, s := range []struct {
		l int
		s []byte
		e *dataLengthError
	}{
		{
			8,
			[]byte{},
			&dataLengthError{0, 8},
		},
		{
			8,
			[]byte{1},
			&dataLengthError{1, 8},
		},
		{
			8,
			[]byte{1, 2, 3, 4, 5, 6, 7},
			&dataLengthError{7, 8},
		},
		{
			8,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			&dataLengthError{12, 8},
		},
		{
			16,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			&dataLengthError{8, 16},
		},
		{
			16,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			&dataLengthError{12, 16},
		},
		{
			16,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			&dataLengthError{18, 16},
		},
	} {
		t.Run(string(i), func(t *testing.T) {
			_, err := unpadData(s.s, s.l)
			assert.Equal(t, err, s.e)
		})
	}
}

func TestUnpadData_paddingLengthError(t *testing.T) {
	for i, s := range []struct {
		l int
		s []byte
		e *paddingLengthError
	}{
		{
			8,
			[]byte{1, 2, 3, 4, 5, 6, 7, 9},
			&paddingLengthError{9, 8},
		},
		{
			8,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			&paddingLengthError{16, 8},
		},
		{
			16,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17},
			&paddingLengthError{17, 16},
		},
	} {
		t.Run(string(i), func(t *testing.T) {
			_, err := unpadData(s.s, s.l)
			assert.Equal(t, err, s.e)
		})
	}
}

func TestUnpadData_invalidPaddingError(t *testing.T) {
	for i, s := range []struct {
		l int
		s []byte
		e *invalidPaddingError
	}{
		{
			8,
			[]byte{1, 8, 8, 8, 8, 8, 8, 8},
			&invalidPaddingError{[]byte{1, 8, 8, 8, 8, 8, 8, 8}},
		},
		{
			8,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 1, 2, 3, 4},
			&invalidPaddingError{[]byte{1, 2, 3, 4}},
		},
	} {
		t.Run(string(i), func(t *testing.T) {
			_, err := unpadData(s.s, s.l)
			assert.Equal(t, err, s.e)
		})
	}
}
