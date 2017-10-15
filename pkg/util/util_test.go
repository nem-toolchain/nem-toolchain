// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringsDistinct(t *testing.T) {
	assert.Empty(t, DistinctStrings(nil))

	assert.Equal(t, []string{"a"}, DistinctStrings([]string{"a"}))
	assert.Equal(t, []string{"a"}, DistinctStrings([]string{"a", "a"}))
	assert.Equal(t, []string{"a"}, DistinctStrings([]string{"a", "a", "a"}))

	assert.Equal(t, []string{"a", "b"}, DistinctStrings([]string{"a", "b"}))
	assert.Equal(t, []string{"a", "b"}, DistinctStrings([]string{"a", "a", "b", "b", "b"}))

	assert.Equal(t, []string{"a", "b", "a"}, DistinctStrings([]string{"a", "b", "a"}))
	assert.Equal(t, []string{"a", "b", "a"}, DistinctStrings([]string{"a", "a", "b", "a"}))
	assert.Equal(t, []string{"a", "b", "a"}, DistinctStrings([]string{"a", "b", "b", "a"}))
}

func TestStringsIntersection(t *testing.T) {
	assert.Empty(t, IntersectStrings(nil, nil))
	assert.Empty(t, IntersectStrings(nil, []string{}))
	assert.Empty(t, IntersectStrings([]string{"a"}, nil))

	assert.Empty(t, IntersectStrings([]string{"a", "b", "c"}, []string{}))
	assert.Empty(t, IntersectStrings([]string{"a", "b", "c"}, []string{"1", "2", "3"}))

	assert.Equal(t, []string{"a"},
		IntersectStrings([]string{"a", "bb", "ccc"}, []string{"1", "2", "a"}))
	assert.Equal(t, []string{"a", "bb", "ccc"},
		IntersectStrings([]string{"a", "bb", "ccc"}, []string{"ccc", "bb", "a"}))
}
