package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"resemble/data"
	"resemble/scoring"
)

var (
	version    string = ".01"
	versionSHA string
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

var Usage = func() {
	fmt.Fprintln(os.Stderr, "Given an image, find similar images in the same directory")
	fmt.Fprintf(os.Stderr, "Usage:  %s filename \n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	cntflg := flag.Int("cnt", 10, "how many results to return")
	htflg := flag.String("hashtype", "P", "which hash type to use for similarity check")
	versionflg := flag.Bool("version", false, "show version")

	flag.Usage = Usage
	flag.Parse()

	if *versionflg == true {
		fmt.Fprintf(os.Stderr, "%s version %s-%s\n", os.Args[0], version, versionSHA)
		return
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "No filename given as argument.\n")
		Usage()
		return
	}

	filename := flag.Args()[0]
	if !FileExists(filename) {
		fmt.Fprintln(os.Stderr, "Missing filename given as argument.")
		Usage()
		return
	}
	ht := *htflg
	err := scoring.ValidateRequestedHashType(ht)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Usage()
		return
	}

	images, _ := data.LoadCorpus()
	images, updated, err := data.RefreshCorpus(images)
	if err != nil {
		panic(err)
	}
	comparisons := scoring.CalcBinarySimilarity(ht, filename, images)
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
