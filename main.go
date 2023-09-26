package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"golang.org/x/sync/errgroup"
)

func main() {
	var imagePath string
	var saveDirPath string

	var width float64
	var height float64
	var quality int

	flag.StringVar(&imagePath, "imdir", "image/", "the directory or file path for images")
	flag.StringVar(&saveDirPath, "svdir", "thumbnail/", "the directory path to save")
	flag.Float64Var(&width, "wid", 0, "Specify what the denominator of the original image width.\nThe height will also be adjusted to match the aspect ratio of the original image")
	flag.Float64Var(&height, "hei", 0, "Specify what the denominator of the original image height.\nThe width will also be adjusted to match the aspect ratio of the original image")
	flag.IntVar(&quality, "q", 80, "the quality to resize")
	flag.Parse()

	if err := run(imagePath, saveDirPath, width, height, quality); err != nil {
		log.Fatalln(err)
	}
}

func run(imagePath, saveDirPath string, width, height float64, quality int) error {
	// ファイル情報取得
	fileinfo, err := os.Stat(imagePath)
	if err != nil {
		return err
	}

	// ファイルかフォルダかを判別しファイルだったら1回リサイズ
	if !fileinfo.IsDir() {
		return resizeImage(imagePath, saveDirPath, width, height, quality)
	}

	fmt.Println("directory")
	// フォルダだったら拡張子jpeg,jpg,pngのものを抜き出して繰り返し実行
	files, err := os.ReadDir(imagePath)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)

	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".jpeg" ||
			strings.ToLower(filepath.Ext(file.Name())) == ".jpg" ||
			strings.ToLower(filepath.Ext(file.Name())) == ".png" {
			imagePath2 := filepath.Join(imagePath, file.Name())
			fmt.Println(imagePath2)
			g.Go(
				func() error {
					return resizeImage(imagePath2, saveDirPath, width, height, quality)
				})
		}
	}

	err = g.Wait()
	if err != nil {
		return err
	}

	fmt.Println("正常終了")
	return nil
}

func resizeImage(imagePath, saveDirPath string, width, height float64, quality int) error {
	// 元の画像の読み込み
	file := imagePath
	fileData, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("resizeImage: os.Open: %v", err)
	}

	// 画像をimage.Image型にdecodeします
	img, data, err := image.Decode(fileData)
	if err != nil {
		return fmt.Errorf("resizeImage: image.Decode: %v", err)
	}
	fmt.Printf("%s: data: %s\n", filepath.Base(imagePath), data)
	fmt.Printf("%s: image width: %d\n", filepath.Base(imagePath), img.Bounds().Dx())
	fmt.Printf("%s: image height: %d\n", filepath.Base(imagePath), img.Bounds().Dy())
	fileData.Close()

	widthInt := uint(width)
	heightInt := uint(height)
	if width != 0 {
		widthInt = uint(math.Round(float64(img.Bounds().Dx()) / width))
	}
	if height != 0 {
		heightInt = uint(math.Round(float64(img.Bounds().Dy()) / height))
	}

	// ここでリサイズします
	// 片方のサイズを0にするとアスペクト比固定してくれます
	resizedImg := resize.Resize(widthInt, heightInt, img, resize.NearestNeighbor)
	fmt.Printf("%s: resized image width: %d\n", filepath.Base(imagePath), resizedImg.Bounds().Dx())
	fmt.Printf("%s: resized image height: %d\n", filepath.Base(imagePath), resizedImg.Bounds().Dy())

	// 書き出すファイル名を指定します
	createFilePath := path.Join(saveDirPath, filepath.Base(imagePath)) // + "." + data
	output, err := os.Create(createFilePath)
	if err != nil {
		return fmt.Errorf("resizeImage: os.Create: %v", err)
	}
	fmt.Println("save as:", createFilePath)
	// 最後にファイルを閉じる
	defer output.Close()

	// 画像のエンコード(書き込み)
	switch data {
	case "png":
		if err := png.Encode(output, resizedImg); err != nil {
			return fmt.Errorf("resizeImage: png.Encode: %v", err)
		}
	case "jpeg", "jpg":
		opts := &jpeg.Options{Quality: int(quality)}
		if err := jpeg.Encode(output, resizedImg, opts); err != nil {
			return fmt.Errorf("resizeImage: jpeg.Encode: %v", err)
		}
	default:
		return errors.New("unknown file format")
	}

	return nil
}
