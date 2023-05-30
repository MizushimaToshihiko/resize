package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

const imagePath = "image/"
const thumbnailPath = "thumbnail/"
const width = 256
const height = 0

// ここに元データのファイル名を書きます。
var imageName = "sakura.jpg"

// 変更後のファイル名を書きます。
var thumbnailName = "sakura-thumbnail"

func main() {
	resizeImage()
}

func resizeImage() {
	// 元の画像の読み込み
	file := imagePath + imageName
	fileData, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	// 画像をimage.Image型にdecodeします
	img, data, err := image.Decode(fileData)
	if err != nil {
		log.Fatal(err)
	}
	fileData.Close()

	// ここでリサイズします
	// 片方のサイズを0にするとアスペクト比固定してくれます
	resizedImg := resize.Resize(width, height, img, resize.NearestNeighbor)

	// 書き出すファイル名を指定します
	createFilePath := thumbnailPath + thumbnailName + "." + data
	output, err := os.Create(createFilePath)
	if err != nil {
		log.Fatal(err)
	}
	// 最後にファイルを閉じる
	defer output.Close()

	// 画像のエンコード(書き込み)
	switch data {
	case "jpeg", "jpg":
		opts := &jpeg.Options{Quality: 100}
		if err := jpeg.Encode(output, resizedImg, opts); err != nil {
			log.Fatal(err)
		}
	default: // "png"
		if err := png.Encode(output, resizedImg); err != nil {
			log.Fatal(err)
		}
	}
}
