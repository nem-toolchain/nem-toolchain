// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"reflect"

	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

// seqSelector allows nested selectors to be combined into a sequential multi selector chain (`AND` logic)
type seqSelector struct {
	items []Selector
}

// AndSelector combines several selectors into a sequential multi selector chain (`AND` logic)
func AndSelector(selectors ...Selector) Selector {
	return seqSelector{selectors}
}

func (sel seqSelector) Pass(addr keypair.Address) bool {
	for _, it := range sel.items {
		if !it.Pass(addr) {
			return false
		}
	}
	return true
}

func (sel seqSelector) rules() []searchRule {
	res := []searchRule{{}}
	for _, it := range sel.items {
		n := make([]searchRule, 0)
		for _, r1 := range res {
			for _, r2 := range it.rules() {
				n = append(n, r1.merge(r2))
			}
		}
		res = n
	}
	return res
}

// parSelector allows nested selectors to be combined into a parallel multi selector chain (`OR` logic)
type parSelector struct {
	items []Selector
}

// OrSelector combines several selectors into a parallel multi selector chain (`OR` logic)
func OrSelector(selectors ...Selector) Selector {
	return parSelector{selectors}
}

func (sel parSelector) Pass(addr keypair.Address) bool {
	for _, it := range sel.items {
		if it.Pass(addr) {
			return true
		}
	}
	return len(sel.items) == 0
}

func (sel parSelector) rules() []searchRule {
	if len(sel.items) == 0 {
		return []searchRule{{}}
	}
	res := make([]searchRule, 0)
	for _, it := range sel.items {
	OUTER:
		for _, r := range it.rules() {
			if r == (searchRule{}) {
				return []searchRule{{}}
			}
			for _, o := range res {
				if reflect.DeepEqual(r, o) {
					continue OUTER
				}
			}
			res = append(res, r)
		}
	}
	return res
}
