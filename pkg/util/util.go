// Copyright 2017 The nem-toolchain project authors. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package util contains common util methods.
package util

// DistinctStrings skip all duplicates in an ordered slice of strings.
func DistinctStrings(a []string) []string {
	res := make([]string, 0)
	var s string
	for i := 0; i < len(a); i++ {
		if s != a[i] {
			res = append(res, a[i])
			s = a[i]
		}
	}
	return res
}

// IntersectStrings intersects two sets of string slices.
func IntersectStrings(a []string, b []string) []string {
	res := make([]string, 0)
	for _, c := range a {
		for _, d := range b {
			if c == d {
				res = append(res, c)
			}
		}
	}
	return res
}
