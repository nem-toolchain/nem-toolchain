// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package vanity

import (
	"sort"
	"strings"

	"github.com/nem-toolchain/nem-toolchain/pkg/util"
)

// searchRule is one row declarative search rule used for example to calculate probability.
type searchRule struct {
	exclude *excludeSelector
	prefix  *prefixSelector
}

func (rule searchRule) merge(other searchRule) searchRule {
	if rule.exclude == nil {
		rule.exclude = other.exclude
	} else if other.exclude != nil {
		s := rule.exclude.merge(*other.exclude)
		rule.exclude = &s
	}
	if rule.prefix == nil {
		rule.prefix = other.prefix
	} else if other.prefix != nil && rule.prefix != other.prefix {
		panic("an attempt to merge two different prefix selectors")
	}
	return rule
}

func (sel excludeSelector) merge(other excludeSelector) excludeSelector {
	arr := append(
		strings.Split(sel.chars, ""),
		strings.Split(other.chars, "")...)
	sort.Strings(arr)
	return excludeSelector{
		strings.Join(util.DistinctStrings(arr), "")}
}
