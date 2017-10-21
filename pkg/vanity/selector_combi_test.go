// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestSeqMultiSelector_Pass(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.True(t, seqSelector{}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{TrueSelector{}}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{excludeSelector{"234567"}}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{prefixSelector{"TAA"}}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{excludeSelector{"BCD"}, prefixSelector{"TA"}}}.Pass(addr))
}

func TestSeqMultiSelector_Pass_fail(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.False(t, seqSelector{[]Selector{FalseSelector{}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{TrueSelector{}, FalseSelector{}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{excludeSelector{"ABC"}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{prefixSelector{"TB"}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{excludeSelector{"BCD"}, prefixSelector{"TB"}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{excludeSelector{"ABC"}, prefixSelector{"TA"}}}.Pass(addr))
}

func TestSeqMultiSelector_rules(t *testing.T) {
	assert.Equal(t, []searchRule{{}}, seqSelector{}.rules())
	assert.Equal(t, []searchRule{}, seqSelector{[]Selector{FalseSelector{}}}.rules())
	assert.Equal(t, []searchRule{}, seqSelector{[]Selector{FalseSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t, []searchRule{{}}, seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{exclude: &excludeSelector{}}},
		seqSelector{[]Selector{TrueSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{exclude: &excludeSelector{}}},
		seqSelector{[]Selector{excludeSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{prefix: &prefixSelector{"TA"}}},
		seqSelector{[]Selector{
			seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}}, prefixSelector{"TA"},
		}}.rules())
	assert.Equal(t,
		[]searchRule{{&excludeSelector{"BCD"}, &prefixSelector{"TA"}}},
		seqSelector{[]Selector{excludeSelector{"BCD"}, prefixSelector{"TA"}}}.rules())
}

func TestParMultiSelector_Pass(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.True(t, parSelector{}.Pass(addr))
	assert.True(t, parSelector{[]Selector{TrueSelector{}}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{TrueSelector{}, FalseSelector{}}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{excludeSelector{"234567"}}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{excludeSelector{"BCD"}, prefixSelector{"TB"}}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{excludeSelector{"ABC"}, prefixSelector{"TA"}}}.Pass(addr))
}

func TestParMultiSelector_Pass_fail(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.False(t, parSelector{[]Selector{FalseSelector{}}}.Pass(addr))
	assert.False(t, parSelector{[]Selector{excludeSelector{"ABC"}}}.Pass(addr))
	assert.False(t, parSelector{[]Selector{prefixSelector{"TB"}}}.Pass(addr))
	assert.False(t, parSelector{[]Selector{excludeSelector{"ABC"}, prefixSelector{"TB"}}}.Pass(addr))
}

func TestParMultiSelector_rules(t *testing.T) {
	assert.Equal(t, []searchRule{{}}, parSelector{}.rules())
	assert.Equal(t, []searchRule{}, parSelector{[]Selector{FalseSelector{}}}.rules())
	assert.Equal(t, []searchRule{{}}, parSelector{[]Selector{TrueSelector{}, FalseSelector{}}}.rules())
	assert.Equal(t, []searchRule{{}}, parSelector{[]Selector{TrueSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{exclude: &excludeSelector{}}},
		parSelector{[]Selector{excludeSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{exclude: &excludeSelector{"BCD"}}, {prefix: &prefixSelector{"TA"}}},
		parSelector{[]Selector{excludeSelector{"BCD"}, prefixSelector{"TA"}}}.rules())
	assert.Equal(t,
		[]searchRule{
			{&excludeSelector{"2BCD"}, &prefixSelector{"TA"}},
			{&excludeSelector{"4BCD"}, &prefixSelector{"TA"}},
			{&excludeSelector{"6BCD"}, &prefixSelector{"TA"}},
		},
		seqSelector{[]Selector{
			parSelector{[]Selector{
				excludeSelector{"2"},
				excludeSelector{"4"},
				excludeSelector{"6"},
			}},
			seqSelector{[]Selector{
				excludeSelector{"B"},
				excludeSelector{"C"},
				excludeSelector{"D"},
				prefixSelector{"TA"},
			}},
			seqSelector{[]Selector{
				seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}},
				parSelector{[]Selector{TrueSelector{}, FalseSelector{}, excludeSelector{}}},
			}},
		}}.rules())
}
