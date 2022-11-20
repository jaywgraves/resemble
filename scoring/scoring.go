package scoring

import (
	"resemble/models"
	"sort"

	"github.com/ajdnik/imghash/similarity"
)

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

func CalcBinarySimilarity(filename string, corpus models.ImageCorpus) BinaryComparisons {
	count := len(corpus.Images)
	comparisons := NewBinaryComparisons(count - 1)
	hash := corpus.Images[filename].PHash
	i := 0
	for k, v := range corpus.Images {
		// no reason to check similarity vs self
		if filename == k {
			continue
		}
		d := float64(similarity.Hamming(hash, v.PHash))
		comparisons[i] = BinaryComparison{Score: d, FileName: k}
		i += 1
	}
	sort.Sort(comparisons)
	return comparisons
}
