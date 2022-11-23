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
	version    string = ".02"
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
	var ht, filename string
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
	versionCmd.Usage = Usage
	matchCmd := flag.NewFlagSet("match", flag.ExitOnError)
	matchCmd.Usage = Usage
	cntFlg := matchCmd.Int("cnt", 10, "how many results to return")
	htFlg1 := matchCmd.String("hashtype", "P", "which hash type to use for similarity check")

	duplicatesCmd := flag.NewFlagSet("duplicates", flag.ExitOnError)
	duplicatesCmd.Usage = Usage
	htFlg2 := duplicatesCmd.String("hashtype", "P", "which hash type to use for similarity check")

	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		versionCmd.Parse(os.Args[2:])
		fmt.Fprintf(os.Stderr, "%s version %s-%s\n", os.Args[0], version, versionSHA)
		return
	case "match":
		matchCmd.Parse(os.Args[2:])
		if matchCmd.NArg() != 1 {
			fmt.Fprintf(os.Stderr, "No filename given as argument.\n")
			Usage()
			return
		}
		filename = matchCmd.Args()[0]
		if !FileExists(filename) {
			fmt.Fprintln(os.Stderr, "Invalid filename given as argument.")
			Usage()
			return
		}
		ht = *htFlg1
		err := scoring.ValidateRequestedHashType(ht)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			Usage()
			return
		}
	case "duplicates":
		ht := *htFlg2
		err := scoring.ValidateRequestedHashType(ht)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			Usage()
			return
		}
		fmt.Fprintln(os.Stderr, "not yet implementd")
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
		if cnt >= *cntFlg {
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
