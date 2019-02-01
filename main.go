package main

import (
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"os"
	"path/filepath"
)

// 输入：
// 1. -f 指定文件，出现顺序（默认）
// 2. -i 指定目录
//
// 处理：
// 1. -m max=按宽度最大的缩放（默认）min=按宽度最小的缩放
// 2. -s 结果缩放
//
// 输出
// 1. -od 输出目录，默认文件名
// 2. -of 输出文件

var (
	input      = flag.String("i", "", "input directory")
	mode       = flag.String("m", "max", fmt.Sprintf("merge mode, %s or %s", Max, Min))
	outputFile = flag.String("of", "", "output file.")
	outputDir  = flag.String("od", "", "output folder.")
	scale      = flag.Float64("s", 0.0, "scale the output.")
	exts       = []string{".jpg", "jpeg", "png"}
)

func main() {
	flag.Parse()
	inputFiles := flag.Args()

	if *input == "" && len(inputFiles) == 0 {
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
	fmt.Printf("output picture: %d*%d\n", w, h)
	dc := gg.NewContext(w, h)
	var currentH int
	for _, img := range imgs {
		dc.DrawImage(img, 0, currentH)
		currentH += img.Bounds().Dy()
	}

	if err := save(dc, getOutput(*outputFile, *outputDir)); err != nil {
		panic(err)
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

type picture struct {
	image.Image
	path string
}

const (
	Max   = "max"
	Min   = "min"
	Usage = `Usage:
1. imgmerge -od ~/Desktop/ -i ~/Desktop/imgs
2. imgmerge -of ~/Desktop/imgmerge.png ~/Desktop/1.jpg ~/Desktop/2.jpg
3. imgmerge -od ~/Desktop -m min -i ~/Desktop/imgs
4. imgmerge -od ~/Desktop -i ~/Desktop/imgs -s 1.5`
)
