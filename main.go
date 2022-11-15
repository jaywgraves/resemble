package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"resemble/data"
	"resemble/models"
	"sort"

	"github.com/ajdnik/imghash/similarity"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

var Usage = func() {
	fmt.Fprint(os.Stderr, "Given an image, find similar images in the same directory")
	fmt.Fprintf(os.Stderr, "Usage:  %s filename \n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	cntflg := flag.Int("cnt", 10, "how many results to return")

	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "No filename given as argument.\n")
		Usage()
		return
	}

	filename := flag.Args()[0]
	if !FileExists(filename) {
		fmt.Fprintf(os.Stderr, "Missing filename given as argument.\n")
		Usage()
		return
	}

	images, _ := data.LoadCorpus()
	images, updated, err := data.RefreshCorpus(images)
	if err != nil {
		panic(err)
	}
	comparisons := calcSimilarity(filename, images)
	fmt.Println("Score\tFilename")
	cnt := 0
	for _, comp := range comparisons {
		fmt.Printf("%d\t%s\n", int(comp.Score), comp.FileName)
		cnt += 1
		if cnt >= *cntflg {
			break
		}
	}
	if updated {
		err = data.SaveCorpus(images)
		if err != nil {
			panic(err)
		}
	}
}

func calcSimilarity(filename string, corpus models.ImageCorpus) models.Comparisons {
	count := len(corpus.Images)
	comparisons := models.NewComparisons(count - 1)
	hash := corpus.Images[filename].PHash
	i := 0
	for k, v := range corpus.Images {
		// no reason to check similarity vs self
		if filename == k {
			continue
		}
		d := float64(similarity.Hamming(hash, v.PHash))
		comparisons[i] = models.Comparison{Score: d, FileName: k}
		i += 1
	}
	// reverse sort
	sort.Sort(comparisons)
	return comparisons
}
