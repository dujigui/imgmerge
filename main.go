package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gwpp/tinify-go/tinify"
	"github.com/nfnt/resize"
	"image"
	"os"
	"path/filepath"
)

// 输入：
// 1. -f 指定输入文件，按检索顺序（默认）
// 2. -i 指定输入目录
//
// 处理：
// 1. -m max=按宽度最大缩放（默认）min=按宽度最小缩放
// 2. -s 缩放图片, 0 < i < 100
// 3. -c 压缩图片
// 4. -k tinypng.com 提供的 API key
//
// 输出
// 1. -od 指定输出到目录，使用默认文件名 (imgmerge_<timestam>.png)
// 2. -of 指定输出到文件。
//
// 其他
// 1. -v 版本信息

var (
	input         = flag.String("i", "", "input directory")
	mode          = flag.String("m", "max", fmt.Sprintf("merge mode, %s or %s", Max, Min))
	outputFile    = flag.String("of", "", "output file.")
	outputDir     = flag.String("od", "", "output folder.")
	scale         = flag.Float64("s", 0.0, "scale the output.")
	doCompress    = flag.Bool("c", false, "compress using tinypng.com")
	tinypngAPIKey = flag.String("k", "", "api key from tinypng.com")
	version       = flag.Bool("v", false, "print version info")
	exts          = []string{".jpg", "jpeg", "png"}
)

func main() {
	flag.Parse()
	inputFiles := flag.Args()

	if Debug {
		fmt.Printf(
			"input=%s\nmode=%s\noutputFile=%s\noutputDir=%s\nscale=%f\ndoCompress=%t\ntinypngAPIKey=%s\n",
			*input, *mode, *outputFile, *outputDir, *scale, *doCompress, *tinypngAPIKey)
	}

	if *version {
		fmt.Printf("Version: %s\nAuthor: github.com/dujigui/imgmerge\n", Version)
		return
	}

	if *input == "" && len(inputFiles) == 0 {
		fmt.Println("Please supply input")
		fmt.Println(Usage)
		flag.PrintDefaults()
		return
	}

	if *doCompress && *tinypngAPIKey == "" {
		fmt.Println("Please apply for api key from tinypng.com")
		fmt.Println(Usage)
		flag.PrintDefaults()
		return
	}

	imgs := load(*input, inputFiles)
	fmt.Printf("found %d images: \n", len(imgs))
	for i, img := range imgs {
		fmt.Printf("%d: %s(%d*%d)\n", i, img.path, img.Bounds().Dx(), img.Bounds().Dy())
	}

	w, h, imgs := scaleImages(imgs, *mode)
	pw := w
	ph := h
	if *scale != 0 {
		pw = int(float64(w) * *scale)
		ph = int(float64(h) * *scale)
	}
	fmt.Printf("output picture: %d*%d\n", pw, ph)
	dc := gg.NewContext(w, h)
	var currentH int
	for _, img := range imgs {
		dc.DrawImage(img, 0, currentH)
		currentH += img.Bounds().Dy()
	}

	o := getOutput(*outputFile, *outputDir)
	if err := save(dc, o); err != nil {
		panic(err)
	}

	if *doCompress {
		if err := compress(o, o); err != nil {
			panic(err)
		}
	}
}

func load(inputDir string, inputFiles []string) []picture {
	if inputDir != "" {
		inputDir = absPath(inputDir)
		var files []string
		f := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}
			if !info.IsDir() && isImage(path) {
				files = append(files, path)
			}
			return nil
		}
		if err := filepath.Walk(inputDir, f); err != nil {
			panic(err)
		}
		return loadImages(files)
	}
	return loadImages(inputFiles)
}

func loadImages(files []string) (result []picture) {
	if len(files) == 0 {
		panic(errors.New("empty input"))
		return
	}
	for _, file := range files {
		img, err := gg.LoadImage(file)
		if err != nil {
			panic(err)
		}
		result = append(result, picture{img, file})
	}
	return result
}

func scaleImages(imgs []picture, scaleMode string) (w, h int, output []picture) {
	var f func(int, int) int
	if scaleMode == Max {
		f = max
	} else {
		f = min
	}
	for _, img := range imgs {
		if w == 0 {
			w = img.Bounds().Dx()
		} else {
			w = f(img.Bounds().Dx(), w)
		}
	}
	for i, img := range imgs {
		if img.Bounds().Dx() != w {
			h := float64(w) / float64(img.Bounds().Dx()) * float64(img.Bounds().Dy())
			resized := resize.Resize(uint(w), uint(h), img.Image, resize.Lanczos3)
			imgs[i] = picture{resized, imgs[i].path}
			fmt.Printf("scaling picture %d from %d*%d to %d*%d\n",
				i, img.Bounds().Dx(), img.Bounds().Dy(), imgs[i].Bounds().Dx(), imgs[i].Bounds().Dy())
		}
		h += imgs[i].Bounds().Dy()
	}
	return w, h, imgs
}

func save(c *gg.Context, output string) error {
	b := c.Image().Bounds()
	w := float64(b.Dx()) * *scale
	h := float64(b.Dy()) * *scale
	resized := resize.Resize(uint(w), uint(h), c.Image(), resize.Lanczos3)
	err := gg.NewContextForImage(resized).SavePNG(output)
	if err == nil {
		fmt.Printf("saved to: %s\n", output)
	}
	return err
}

func compress(input, output string) error {
	fs := fileSize(input)
	fmt.Printf("before compress: %d bytes\n", fs)

	Tinify.SetKey(*tinypngAPIKey)
	source, err := Tinify.FromFile(input)
	if err != nil {
		return errors.New("compress fail")
	}

	fmt.Println("compressing...")
	err = source.ToFile(output)
	if err != nil {
		return errors.New("output fail")
	}

	fs = fileSize(input)
	fmt.Printf("after compress: %d bytes\n", fs)

	return nil
}

func fileSize(input string) int64 {
	fi, err := os.Stat(input)
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

type picture struct {
	image.Image
	path string
}

const (
	Debug   = false
	Version = "1.0.1"
	Max     = "max"
	Min     = "min"
	Usage   = `Usage:
1. imgmerge -od ~/Desktop/ -i ~/Desktop/imgs
2. imgmerge -of ~/Desktop/imgmerge.png ~/Desktop/1.jpg ~/Desktop/2.jpg
3. imgmerge -od ~/Desktop -m min -i ~/Desktop/imgs
4. imgmerge -od ~/Desktop -i ~/Desktop/imgs -s 1.5
5. imgmerge -od ~/Desktop/ -i ~/Desktop/imgs -c -k yourAPIkey`
)
