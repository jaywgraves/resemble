package phash

import (
	"image"
	"os"

	"github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
)

func GetPHash(filename string) hashtype.Binary {
	phash := imghash.NewPHash()
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

	h := phash.Calculate(img)
	return h
}
