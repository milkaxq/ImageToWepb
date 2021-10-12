package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

var (
	imageName = "Hello"
)

func main() {
	someBase64Image := ""
	wantPutWatermark := "true"
	imageType, image := decBase64ToImage(someBase64Image)
	switch imageType {
	//Use Case for Png
	case "png":
		if wantPutWatermark == "true" {
			putWaterMark(image)
			imageToWebp()
		}
	//Use Default for Jpeg and Jpg Images it have some errors in Type
	default:
		ownImage := createOwnJpegImage(image)
		putWaterMark(ownImage)
		imageToWebp()
	}
}

// Decode Base64 To Image To Get Type (Jpeg/Png)
func decBase64ToImage(base64Image string) (string, image.Image) {
	//Decode base64 To []Byte
	imgByte, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		fmt.Println(err.Error())
	}
	//Returns image/png
	filetype := http.DetectContentType(imgByte)
	//Returns imageString[1] png
	imageString := strings.SplitAfter(filetype, "/")
	//[]Byte to Image.image
	img, _, _ := image.Decode(bytes.NewReader(imgByte))
	return imageString[1], img
}

//How To Put WaterMark
func putWaterMark(img image.Image) {
	//Open Watermark Image
	waterMark, err := os.Open("watermark.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	//Decode Png To Image.image
	decWatermark, err := png.Decode(waterMark)
	if err != nil {
		fmt.Printf("failed to decode: %s", err)
	}
	//Resize WaterMark By Given Image
	resizedDecWatermark := resize.Resize(uint(img.Bounds().Max.X), uint(img.Bounds().Max.Y), decWatermark, resize.Lanczos3)
	//Get Size of Given Image
	b := img.Bounds()
	//Create RGBA Image of given Image
	imageForWatermark := image.NewRGBA(b)
	//Draw To RGBA given Image
	draw.Draw(imageForWatermark, b, img, image.Point{}, draw.Src)
	//Draw Over WaterMark
	draw.Draw(imageForWatermark, resizedDecWatermark.Bounds(), resizedDecWatermark, image.Point{}, draw.Over)
	//Create File for ImageWith Watermark
	imageWithWatermark, err := os.Create(fmt.Sprintf("mark%s.png", imageName))
	if err != nil {
		fmt.Printf("failed to create: %s", err)
	}
	//Encode Png To Given File
	png.Encode(imageWithWatermark, imageForWatermark)
}

// How To Convert Given Image to Webp
func imageToWebp() {
	var buf bytes.Buffer
	//Open encoded Png Image
	imageWithWatermarkRaw, err := os.Open(fmt.Sprintf("mark%s.png", imageName))
	if err != nil {
		fmt.Println(err)
	}
	//Decode Png Image To Image.image
	imageWithWaterMark, err := png.Decode(imageWithWatermarkRaw)
	if err != nil {
		fmt.Printf("failed to create: %s", err)
	}
	//Give Image.image to Convert it To Webp
	if err := webp.Encode(&buf, imageWithWaterMark, &webp.Options{Quality: 75}); err != nil {
		fmt.Println(err)
	}
	//Create And Write Webp to File
	if err := ioutil.WriteFile(fmt.Sprintf("%s.webp", imageName), buf.Bytes(), 0666); err != nil {
		fmt.Println(err)
	}
}

//Create Own Jpeg Image To Work With
func createOwnJpegImage(img image.Image) image.Image {
	//Create Temp Jpeg File
	out, err := os.Create("temp.jpeg")
	if err != nil {
		fmt.Println(err)
	}
	//Write Jpeg To That File
	err = jpeg.Encode(out, img, &jpeg.Options{Quality: 100})
	if err != nil {
		fmt.Println(err)
	}
	//Decode and Get image.Image file
	mainImage, _ := os.Open("temp.jpeg")
	decMainImage, _ := jpeg.Decode(mainImage)
	os.Remove("temp.jpeg")
	return decMainImage
}
