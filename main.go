package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"golang.org/x/sync/errgroup"
)

// ここに元データのファイル名を書きます。
var imageName = "sakura.jpg"

// 変更後のファイル名を書きます。
var thumbnailName = "sakura-thumbnail"

func main() {
	// ログ設定
	log.SetFlags(log.Lshortfile)

	var imagePath string
	var saveDirPath string

	var width uint
	var height uint
	var quality uint

	flag.StringVar(&imagePath, "imdir", "image/", "the directory or file path for images")
	flag.StringVar(&saveDirPath, "svdir", "thumbnail/", "the directory path to save")
	flag.UintVar(&width, "wid", 0, "the width to resize")
	flag.UintVar(&height, "hei", 0, "th height to resize")
	flag.UintVar(&quality, "q", 80, "the quality to resize")
	flag.Parse()

	if err := run(imagePath, saveDirPath, width, height, quality); err != nil {
		log.Fatalln(err)
	}
}

func run(imagePath, saveDirPath string, width, height, quality uint) error {
	// ファイル情報取得
	fileinfo, err := os.Stat(imagePath)
	if err != nil {
		return err
	}

	// ファイルかフォルダかを判別しファイルだったら1回リサイズ
	if !fileinfo.IsDir() {
		return resizeImage(imagePath, saveDirPath, width, height, quality)
	}

	fmt.Println("フォルダだった")
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

func resizeImage(imagePath, saveDirPath string, width, height, quality uint) error {
	// 元の画像の読み込み
	file := imagePath
	fileData, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("resizeImage: os.Open: %v", err)
	}

	// 画像をimage.Image型にdecodeします
	img, data, err := image.Decode(fileData)
	fmt.Println("data:", data)
	if err != nil {
		return fmt.Errorf("resizeImage: image.Decode: %v", err)
	}
	fileData.Close()

	// ここでリサイズします
	// 片方のサイズを0にするとアスペクト比固定してくれます
	resizedImg := resize.Resize(width, height, img, resize.NearestNeighbor)

	// 書き出すファイル名を指定します
	createFilePath := saveDirPath + filepath.Base(imagePath) // + "." + data
	output, err := os.Create(createFilePath)
	if err != nil {
		return fmt.Errorf("resizeImage: os.Create: %v", err)
	}
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
