package file

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

func makeFile(orginal, want string) error {
	ext := strings.Replace(filepath.Ext(want), ".", "", -1)
	format, ok := Format[ext]
	if !ok {
		return ErrorNotValidFile
	}
	if format != Image {
		return errors.New("not an image to create resizer")
	}
	w, h := getSizes(want)
	if w == 0 && h == 0 {
		return errors.New("not valid sizes for resize")
	}
	return resizer(orginal, want, w, h)
}

func getSizes(path string) (uint, uint) {
	rw := regexp.MustCompile("[w,W]([0-9]+)")
	rh := regexp.MustCompile("[h,H]([0-9]+)")
	wMatches := rw.FindStringSubmatch(path)
	hMatches := rh.FindStringSubmatch(path)
	var (
		w, h uint
	)
	if len(wMatches) == 2 {
		if width, err := strconv.Atoi(wMatches[1]); err == nil {
			w = uint(width)
		}
	}
	if len(hMatches) == 2 {
		if height, err := strconv.Atoi(hMatches[1]); err == nil {
			h = uint(height)
		}
	}
	return w, h
}

func resizer(srcPath, destPath string, width, height uint) error {
	var (
		white      = color.RGBA{255, 255, 255, 255}
		point      = image.Pt(0, 0)
		imgW, imgH int
	)
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()
	imgConfig, _, err := image.DecodeConfig(src)
	if err != nil {
		return err
	}
	src.Seek(0, 0)
	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}
	w, h := getValidWH(width, height, uint(imgConfig.Width), uint(imgConfig.Height))
	imgResized := resize.Resize(w, h, img, resize.Bicubic)
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()
	imgW = int(width)
	imgH = int(height)
	if imgW == 0 {
		ratio := float32(imgH) / float32(imgConfig.Height)
		imgW = round(float32(imgConfig.Width) * ratio)
	}
	if imgH == 0 {
		ratio := float32(imgW) / float32(imgConfig.Width)
		imgH = round(float32(imgConfig.Height) * ratio)
	}

	m := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	b := m.Bounds()
	draw.Draw(m, b, &image.Uniform{white}, image.ZP, draw.Src)
	// draw image center if resized image scaled ratio
	if height != 0 && width != 0 {
		if imgConfig.Width > imgConfig.Height {
			y := int((height - h) / 2)
			point = image.Pt(0, y)
		} else {
			x := int((width - w) / 2)
			point = image.Pt(x, 0)
		}
	}
	draw.Draw(m, b, imgResized, b.Min.Sub(point), draw.Src)

	return jpeg.Encode(dest, m, nil)
}

func getValidWH(width, height, orginalW, orginalH uint) (uint, uint) {
	if width == 0 || height == 0 {
		return width, height
	}
	ratio := float32(width) / float32(orginalW)
	toHeight := float32(orginalH) * ratio
	return width, uint(round(toHeight))
}

func round(num float32) int {
	if num-float32(int(num)) >= 0.5 {
		num++
	}
	return int(num)
}
