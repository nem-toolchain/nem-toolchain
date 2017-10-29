// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"testing"

	"regexp"

	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestSeqMultiSelector_Pass_true(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.True(t, seqSelector{}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{TrueSelector{}}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{excludeSelector{"234567"}}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{
		prefixSelector{re: regexp.MustCompile("^TAA")},
	}}.Pass(addr))
	assert.True(t, seqSelector{[]Selector{
		excludeSelector{"BCD"},
		prefixSelector{re: regexp.MustCompile("^TA")},
	}}.Pass(addr))
}

func TestSeqMultiSelector_Pass_false(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.False(t, seqSelector{[]Selector{FalseSelector{}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{TrueSelector{}, FalseSelector{}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{excludeSelector{"ABC"}}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{
		prefixSelector{re: regexp.MustCompile("^TB")},
	}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{
		excludeSelector{"BCD"},
		prefixSelector{re: regexp.MustCompile("^TB")},
	}}.Pass(addr))
	assert.False(t, seqSelector{[]Selector{
		excludeSelector{"ABC"},
		prefixSelector{re: regexp.MustCompile("^TA")},
	}}.Pass(addr))
}

func TestSeqMultiSelector_rules(t *testing.T) {
	assert.Equal(t, []searchRule{{}},
		seqSelector{}.rules())
	assert.Equal(t, []searchRule{},
		seqSelector{[]Selector{FalseSelector{}}}.rules())
	assert.Equal(t, []searchRule{},
		seqSelector{[]Selector{FalseSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t, []searchRule{{}},
		seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{exclude: &excludeSelector{}}},
		seqSelector{[]Selector{TrueSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{exclude: &excludeSelector{}}},
		seqSelector{[]Selector{excludeSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{prefix: &prefixSelector{prefix: "TA"}}},
		seqSelector{[]Selector{
			seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}},
			prefixSelector{prefix: "TA"},
		}}.rules())
	assert.Equal(t,
		[]searchRule{{
			&excludeSelector{"BCD"},
			&prefixSelector{prefix: "TA"},
		}},
		seqSelector{[]Selector{
			excludeSelector{"BCD"}, prefixSelector{prefix: "TA"},
		}}.rules())
}

func TestParMultiSelector_Pass_true(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.True(t, parSelector{}.Pass(addr))
	assert.True(t, parSelector{[]Selector{TrueSelector{}}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{TrueSelector{}, FalseSelector{}}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{excludeSelector{"234567"}}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{
		excludeSelector{"BCD"},
		prefixSelector{re: regexp.MustCompile("^TB")},
	}}.Pass(addr))
	assert.True(t, parSelector{[]Selector{
		excludeSelector{"ABC"},
		prefixSelector{re: regexp.MustCompile("^TA")},
	}}.Pass(addr))
}

func TestParMultiSelector_Pass_false(t *testing.T) {
	addr, _ := keypair.ParseAddress("TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	assert.False(t, parSelector{[]Selector{FalseSelector{}}}.Pass(addr))
	assert.False(t, parSelector{[]Selector{excludeSelector{"ABC"}}}.Pass(addr))
	assert.False(t, parSelector{[]Selector{
		prefixSelector{re: regexp.MustCompile("^TB")},
	}}.Pass(addr))
	assert.False(t, parSelector{[]Selector{
		excludeSelector{"ABC"},
		prefixSelector{re: regexp.MustCompile("^TB")},
	}}.Pass(addr))
}

func TestParMultiSelector_rules(t *testing.T) {
	assert.Equal(t, []searchRule{{}},
		parSelector{}.rules())
	assert.Equal(t, []searchRule{},
		parSelector{[]Selector{FalseSelector{}}}.rules())
	assert.Equal(t, []searchRule{},
		parSelector{[]Selector{FalseSelector{}, FalseSelector{}}}.rules())
	assert.Equal(t, []searchRule{{}},
		parSelector{[]Selector{TrueSelector{}, FalseSelector{}}}.rules())
	assert.Equal(t, []searchRule{{}},
		parSelector{[]Selector{TrueSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{{exclude: &excludeSelector{}}},
		parSelector{[]Selector{excludeSelector{}, excludeSelector{}}}.rules())
	assert.Equal(t,
		[]searchRule{
			{exclude: &excludeSelector{"BCD"}},
			{prefix: &prefixSelector{prefix: "TA"}},
		},
		parSelector{[]Selector{
			excludeSelector{"BCD"},
			prefixSelector{prefix: "TA"},
		}}.rules())
	assert.Equal(t,
		[]searchRule{
			{&excludeSelector{"2BCD"}, &prefixSelector{prefix: "TA"}},
			{&excludeSelector{"4BCD"}, &prefixSelector{prefix: "TA"}},
			{&excludeSelector{"6BCD"}, &prefixSelector{prefix: "TA"}},
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
				prefixSelector{prefix: "TA"},
			}},
			seqSelector{[]Selector{
				seqSelector{[]Selector{TrueSelector{}, TrueSelector{}}},
				parSelector{[]Selector{TrueSelector{}, FalseSelector{}, excludeSelector{}}},
			}},
		}}.rules())
}
