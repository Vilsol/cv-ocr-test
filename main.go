package main

import (
	"bytes"
	_ "embed"
	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
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

	src := imaging.Grayscale(img)
	src = imaging.Sharpen(src, 4)

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, src); err != nil {
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
}
