package scoring

import (
	"errors"
	"resemble/models"
	"sort"
	"strings"

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

func CalcBinarySimilarity(hashType string, filename string, corpus models.ImageCorpus) BinaryComparisons {
	count := len(corpus.Images)
	comparisons := NewBinaryComparisons(count - 1)
	hash := getRequestedHashValue(hashType, corpus.Images[filename])
	i := 0
	for k, v := range corpus.Images {
		// no reason to check similarity vs self
		if filename == k {
			continue
		}
		d := float64(similarity.Hamming(hash, getRequestedHashValue(hashType, v)))
		comparisons[i] = BinaryComparison{Score: d, FileName: k}
		i += 1
	}
	sort.Sort(comparisons)
	return comparisons
}

var ErrInvalidHashType = errors.New("invalid hashtype requested")

func ValidateRequestedHashType(hashType string) error {
	hashType = strings.ToUpper(hashType)
	switch hashType {
	case "P", "A", "D", "M", "MH", "BM":
		return nil
	default:
		return ErrInvalidHashType
	}
}

func getRequestedHashValue(hashType string, img models.Image) []byte {
	hashType = strings.ToUpper(hashType)
	switch hashType {
	case "P":
		return img.PHash
	case "A":
		return img.AHash
	case "D":
		return img.DHash
	case "M":
		return img.MHash
	case "MH":
		return img.MHHash
	case "BM":
		return img.BMHash
	default:
		return img.PHash
	}
}
