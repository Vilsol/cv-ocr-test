package main

import (
	"bytes"
	_ "embed"
	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	"image"
	"image/png"
)

//go:embed ocr-sample.png
var SamplePNG []byte

func main() {
	client := gosseract.NewClient()
	if err := client.SetLanguage("eng"); err != nil {
		panic(err)
	}
	println("Tesseract Version", client.Version())

	img, err := png.Decode(bytes.NewReader(SamplePNG))
	if err != nil {
		panic(err)
	}

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, img); err != nil {
		panic(err)
	}

	if err := client.SetImageFromBytes(buff.Bytes()); err != nil {
		panic(err)
	}

	text, err := client.Text()
	if err != nil {
		panic(err)
	}

	println("Result:")
	println(text)

	cropped := imaging.Crop(img, image.Rect(50, 50, 250, 250))
	find, err := Find(img, cropped)
	if err != nil {
		panic(err)
	}
	println("Found:")
	println(find.String())
}

func Find(static image.Image, dynamic image.Image) (image.Point, error) {
	matStatic, err := gocv.ImageToMatRGB(static)
	if err != nil {
		return image.Point{}, err
	}

	grayStatic := gocv.NewMat()
	gocv.CvtColor(matStatic, &grayStatic, gocv.ColorRGBToGray)

	matDynamic, err := gocv.ImageToMatRGB(dynamic)
	if err != nil {
		return image.Point{}, err
	}

	grayDynamic := gocv.NewMat()
	gocv.CvtColor(matDynamic, &grayDynamic, gocv.ColorRGBToGray)

	m := gocv.NewMat()
	result := gocv.NewMat()
	gocv.MatchTemplate(grayStatic, grayDynamic, &result, gocv.TmCcoeffNormed, m)

	_, _, _, location := gocv.MinMaxLoc(result)
	_ = result.Close()
	_ = m.Close()
	_ = matDynamic.Close()
	_ = grayDynamic.Close()

	return location, nil
}
