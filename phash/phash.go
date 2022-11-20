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
	imagePtr.AHash = getAHash(img)
	imagePtr.DHash = getDHash(img)
	imagePtr.MHash = getMHash(img)
	imagePtr.CMHash = getCMHash(img)
	imagePtr.MHHash = getMHHash(img)
	imagePtr.BMHash = getBMHash(img)
	imagePtr.RVHash = getRVHash(img)

}

func getPHash(img image.Image) hashtype.Binary {
	phash := imghash.NewPHash()
	h := phash.Calculate(img)
	return h
}

func getAHash(img image.Image) hashtype.Binary {
	ahash := imghash.NewAverage()
	h := ahash.Calculate(img)
	return h
}

func getDHash(img image.Image) hashtype.Binary {
	dhash := imghash.NewDifference()
	h := dhash.Calculate(img)
	return h
}

func getMHash(img image.Image) hashtype.Binary {
	mhash := imghash.NewMedian()
	h := mhash.Calculate(img)
	return h
}

func getCMHash(img image.Image) hashtype.Float64 {
	cmhash := imghash.NewColorMoment()
	h := cmhash.Calculate(img)
	return h
}

func getMHHash(img image.Image) hashtype.Binary {
	mhhash := imghash.NewMedian()
	h := mhhash.Calculate(img)
	return h
}

func getBMHash(img image.Image) hashtype.Binary {
	bmhash := imghash.NewBlockMean()
	h := bmhash.Calculate(img)
	return h
}

func getRVHash(img image.Image) hashtype.UInt8 {
	rvhash := imghash.NewRadialVariance()
	h := rvhash.Calculate(img)
	return h
}
