package models

import "time"

type ImageCorpus struct {
	LastRefreshTime time.Time
	Version         int
	Images          map[string]Image
}

type Image struct {
	SizeBytes int64
	ModTime   time.Time
	PHash     []byte
}

func NewImageCorpus() ImageCorpus {
	return ImageCorpus{Images: make(map[string]Image)}
}

type BinaryComparison struct {
	Score    float64
	FileName string
}

type BinaryComparisons []BinaryComparison

func (s BinaryComparisons) Len() int {
	return len(s)
}
func (s BinaryComparisons) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s BinaryComparisons) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func NewBinaryComparisons(length int) BinaryComparisons {
	return make(BinaryComparisons, length)
}
