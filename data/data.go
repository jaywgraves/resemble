package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"resemble/models"
	"resemble/phash"
	"strings"
	"time"
)

const CorpusFileName = ".resemble-cache"

func getConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		// uhhhh.
		log.Panic("can't get working directory")
	}
	p := filepath.Join(wd, CorpusFileName)
	return p
}

func LoadCorpus() (models.ImageCorpus, error) {
	cp := getConfigPath()
	//fmt.Println("LoadCorpus", cp)
	corpus := models.NewImageCorpus()
	cj, err := ioutil.ReadFile(cp)
	//fmt.Printf("empty file %v %v\n", cj, err)
	if err != nil {
	} else {
		err = json.Unmarshal(cj, &corpus)
		if err != nil {
			log.Fatal(err)
			return corpus, err
		}
	}
	return corpus, nil
}

func SaveCorpus(corpus models.ImageCorpus) error {
	cp := getConfigPath()
	//fmt.Println("SaveCorpus", cp)
	cj, err := json.MarshalIndent(corpus, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cp, cj, 0644)
	if err != nil {
		return err
	}
	return nil
}

func isImage(filename string) bool {
	ext := strings.ToLower((path.Ext(filename)))
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		return true
	} else {
		return false
	}
}

func RefreshCorpus(corpus models.ImageCorpus) (models.ImageCorpus, bool, error) {
	var updated bool
	newcorpus := models.NewImageCorpus()
	newcorpus.Version = corpus.Version
	newcorpus.LastRefreshTime = corpus.LastRefreshTime

	// walk image files
	err := filepath.Walk(".", func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			// skip directories
			return nil
		}
		if !isImage(p) {
			// skip non-images (by extension)
			return nil
		}
		// check corpus by filename
		img, fnd := corpus.Images[p]
		if fnd {
			if info.Size() == img.SizeBytes && info.ModTime() == img.ModTime {
				// still exists and didn't change
				newcorpus.Images[p] = img
			} else {
				// still exists but changed.  re-calc hash
				updated = true
				img.PHash = phash.GetPHash(p)
			}
		} else {
			// file exists but is new
			updated = true
			img = models.Image{
				SizeBytes: info.Size(),
				ModTime:   info.ModTime(),
				PHash:     phash.GetPHash(p),
			}
		}
		newcorpus.Images[p] = img
		return nil
	})
	if updated {
		newcorpus.LastRefreshTime = time.Now()
	}
	return newcorpus, updated, err
}
