package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func absPath(path string) string {
	path = strings.TrimSuffix(path, "/")
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME")+"/", 1)
	}
	return path
}

func isImage(path string) bool {
	for _, ext := range exts {
		if endWith(path, ext) {
			return true
		}
	}
	return false;
}

func endWith(path, extension string) bool {
	return strings.HasSuffix(strings.ToLower(path), extension)
}

func getOutput(outputFile, outputDir string) string {
	if outputFile != "" {
		return outputFile
	}

	p := absPath(outputDir)
	if _, err := os.Stat(p); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/imgmerge_%s.png", p, time.Now().Format("20060102150405"))
}

