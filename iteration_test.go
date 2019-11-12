package main

import (
	"strconv"
	"testing"
)

var lIter = 1_000_000

var V string

// BenchmarkHashMapIter is very simple benchmark to see just how much slower
// iterating over items in hashmap is compared to slice iteration.
func BenchmarkHashMapIter(b *testing.B) {
	var hashmap = make(map[int]string)
	for i := 0; i < lIter; i++ {
		hashmap[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, v := range hashmap {
			V = v
		}
	}
}

// BenchmarkSliceRangeIter is simple benchmark for slice iteration. Mainly to
// compare against hashmap iteration and index access of slice.
func BenchmarkSliceRangeIter(b *testing.B) {
	var slice = make([]string, lIter, lIter)
	for i := 0; i < lIter; i++ {
		slice[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, v := range slice {
			V = v
		}
	}
}

// BenchmarkSliceIndexIter is benchmark of slice iteration using integer index.
// The general rule of thumb is that this should be faster than "range", but in
// reality it is not.
func BenchmarkSliceIndexIter(b *testing.B) {
	var slice = make([]string, lIter, lIter)
	for i := 0; i < lIter; i++ {
		slice[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for idx := 0; idx <= len(slice)-1; idx++ {
			V = slice[idx]
		}
	}
}
