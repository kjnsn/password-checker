package filter

import "github.com/willf/bloom"

func NewFilter(numWords uint) *bloom.BloomFilter {
	return bloom.NewWithEstimates(numWords, 0.00000001)
}
