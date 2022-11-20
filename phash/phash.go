package phash

import (
	"image"
	"os"
	"resemble/models"

	"github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
)

func CalcHashes(filename string, imagePtr *models.Image) {
	imgf, err := os.Open(filename)
	if err != nil {
		// shouldn't happen because we already checked the metadata
		panic(err)
	}
	defer imgf.Close()

	img, _, err := image.Decode(imgf)

	if err != nil {
		// should we just skip this?
		// probably return nil and handle upstream
		panic(err)
	}
	imagePtr.PHash = getPHash(img)

}

func getPHash(img image.Image) hashtype.Binary {
	phash := imghash.NewPHash()
	h := phash.Calculate(img)
	return h
}
