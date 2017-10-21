// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchRule_merge(t *testing.T) {
	assert.Equal(t, searchRule{}, searchRule{}.merge(searchRule{}))
	assert.Equal(t,
		searchRule{exclude: &excludeSelector{}},
		searchRule{}.merge(searchRule{exclude: &excludeSelector{}}))
	assert.Equal(t,
		searchRule{prefix: &prefixSelector{}},
		searchRule{}.merge(searchRule{prefix: &prefixSelector{}}))
	assert.Equal(t,
		searchRule{exclude: &excludeSelector{}, prefix: &prefixSelector{}},
		searchRule{exclude: &excludeSelector{}}.merge(searchRule{prefix: &prefixSelector{}}))
	assert.Equal(t,
		searchRule{exclude: &excludeSelector{"246"}, prefix: &prefixSelector{"TABC"}},
		searchRule{exclude: &excludeSelector{"246"}}.merge(searchRule{prefix: &prefixSelector{"TABC"}}))
	assert.Equal(t,
		searchRule{exclude: &excludeSelector{"246BC"}, prefix: &prefixSelector{"TA"}},
		searchRule{exclude: &excludeSelector{"24"}, prefix: &prefixSelector{"TA"}}.
			merge(searchRule{exclude: &excludeSelector{"6BC"}}))
}

func TestSearchRule_fail(t *testing.T) {
	assert.Panics(t, func() { searchRule{prefix: &prefixSelector{}}.merge(searchRule{prefix: &prefixSelector{}}) })
}

func TestNewExcludeSelector_merge(t *testing.T) {
	assert.Equal(t, excludeSelector{}, excludeSelector{}.merge(excludeSelector{}))
	assert.Equal(t, excludeSelector{"A"}, excludeSelector{}.merge(excludeSelector{"A"}))
	assert.Equal(t, excludeSelector{"A"}, excludeSelector{"A"}.merge(excludeSelector{}))
	assert.Equal(t, excludeSelector{"AB"}, excludeSelector{"A"}.merge(excludeSelector{"B"}))
	assert.Equal(t, excludeSelector{"123AB"}, excludeSelector{"AB"}.merge(excludeSelector{"123"}))
}
