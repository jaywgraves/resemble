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

type Comparison struct {
	Score    float64
	FileName string
}

type Comparisons []Comparison

func (s Comparisons) Len() int {
	return len(s)
}
func (s Comparisons) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Comparisons) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func NewComparisons(length int) Comparisons {
	return make(Comparisons, length)
}
