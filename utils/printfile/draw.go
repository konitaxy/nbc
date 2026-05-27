package printfile

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"

	"math/rand/v2"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
	"golang.org/x/image/tiff"
)

func DrawBackground(text string, x, y float64, fontSize float64, file string) bytes.Buffer {

	img, err := gg.LoadImage(file)
	if err != nil {
		panic(err)
	}
	dc := gg.NewContextForImage(img)

	dc.SetColor(color.RGBA{158, 53, 49, 255})

	err = dc.LoadFontFace("./resource/printfile/arialbd.ttf", fontSize)
	if err != nil {
		panic(err)
	}

	// dc.SetColor(color.White)
	dc.DrawStringAnchored(text, float64(x), float64(y), 0, 0)
	var b bytes.Buffer
	jpeg.Encode(&b, dc.Image(), &jpeg.Options{Quality: 180})
	return b
}

func Resize(buf []byte, targetX, targetY int, withWatermark bool) ([]byte, error) {
	// img :=bimg.NewImage(buf)
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	if targetX == 0 {
		targetX = targetY * img.Bounds().Dx() / img.Bounds().Dy()
	}
	if targetY == 0 {
		targetY = targetX * img.Bounds().Dy() / img.Bounds().Dx()
	}
	thumbnail := image.NewNRGBA(image.Rect(0, 0, int(targetX), int(targetY)))
	draw.Draw(thumbnail, img.Bounds(), &image.Uniform{C: color.Alpha{A: 255}}, image.Point{}, draw.Src)
	draw.CatmullRom.Scale(thumbnail, thumbnail.Bounds(), img, img.Bounds(), draw.Over, &draw.Options{})
	// thumbnail := resize.Resize(targetX, targetY, img, resize.Lanczos3)

	b := bytes.Buffer{}

	if withWatermark {
		wmWidth, wmHeight := 200, 200 // 水印大小
		wmSpacing := 300              // 水印之间的间隔
		startX, startY := 0, 0        // 第一行水印起始位置
		lineOffsetY := 150            // 行与行之间的垂直偏移

		maxX, maxY := thumbnail.Bounds().Max.X, thumbnail.Bounds().Max.Y
		wmImg, _ := imaging.Open("./resource/printfile/watermark.png")
		maxWatermarksPerLine := (maxX - startX) / wmSpacing
		if maxWatermarksPerLine > 0 {

			for row := 0; ; row++ {

				offsetY := startY + row*lineOffsetY
				if offsetY+wmHeight > maxY {
					break // 如果超出源图片高度，则停止添加水印
				}
				xOffset := startX
				if row%2 == 1 {
					xOffset = 150 // 第二行从250px开始
				}
				for col := 0; ; col++ {
					if xOffset+wmWidth > maxX {
						break // 如果超出源图片宽度，则停止添加水印
					}
					thumbnail = imaging.Overlay(thumbnail, wmImg, image.Point{xOffset, offsetY}, 0.2)
					xOffset += wmSpacing
				}
			}

		}
	}
	err = jpeg.Encode(&b, thumbnail, &jpeg.Options{Quality: 90})
	return b.Bytes(), err
}

func ResizeWithFactor(buf []byte, factor int) ([]byte, error) {

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	thumbnail := resize.Resize(uint(img.Bounds().Dx()*factor), uint(img.Bounds().Dy()*factor), img, resize.Lanczos3)
	b := bytes.Buffer{}
	err = jpeg.Encode(&b, thumbnail, &jpeg.Options{Quality: 100})
	return b.Bytes(), err
}

func ImageToCmyk(buf []byte, width, height, span, lineWidth, totalWidth int, text string) ([]byte, error) { //单位毫米

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}
	if width > totalWidth {
		if img.Bounds().Dx() > img.Bounds().Dy() {
			img = imaging.Rotate90(img)
		}
		tmp := width
		width = height
		height = tmp
	} else {
		if height > totalWidth {
			if img.Bounds().Dx() > img.Bounds().Dy() {
				img = imaging.Rotate90(img)
			}
		} else {
			if img.Bounds().Dx() < img.Bounds().Dy() {
				img = imaging.Rotate90(img)
				tmp := width
				width = height
				height = tmp
			}
		}
	}
	isHorizal := img.Bounds().Dx() > img.Bounds().Dy()
	textHeight := 0
	if isHorizal && text != "" {
		textHeight = 120
	}

	left := img.Bounds().Dx()*span/width + lineWidth
	top := img.Bounds().Dy()*span/height + lineWidth
	maxX := img.Bounds().Dx() * totalWidth / width //白色区域包括黑框大小
	rtp := img.Bounds()

	canvasWidth := rtp.Dx() + 2*left //黑框内的图片宽度
	rtp.Max.X = maxX
	space := (maxX - canvasWidth) / 2
	rtp.Max.Y = rtp.Dy() + 2*top + textHeight
	newImg := image.NewNRGBA(rtp)
	dx := rtp.Dx()
	dy := rtp.Dy()
	for x := space; x <= space+canvasWidth; x++ {
		for y := textHeight; y < dy; y++ {
			if x < left+space || x > dx-left-space || y < textHeight+top || y > dy-top {
				if x < space+lineWidth || x > dx-lineWidth-space || y < textHeight+lineWidth || y > dy-lineWidth {
					newImg.Set(x, y, color.Black)
				} else if y == textHeight+top-1 || y == dy-top+1 || y == textHeight+top-2 || y == dy-top+2 {
					newImg.Set(x, y, color.Black)
				}
			} else {
				newImg.Set(x, y, img.At(x-left-space, y-top-textHeight))
			}
		}
	}

	if text != "" {
		var lastImg image.Image
		if isHorizal {
			dc := gg.NewContextForImage(newImg)
			dc.SetColor(color.Black)
			err = dc.LoadFontFace("./resource/printfile/arialbd.ttf", 80)
			if err != nil {
				panic(err)
			}
			dc.DrawStringAnchored(text, float64(space+left), float64(textHeight), 0, 0)
			lastImg = dc.Image()
			// newImg = &lastImg
		} else {

			tempImg := imaging.Rotate270(newImg)
			dc := gg.NewContextForImage(tempImg)
			dc.SetColor(color.Black)
			err = dc.LoadFontFace("./resource/printfile/arialbd.ttf", 80)
			if err != nil {
				panic(err)
			}
			dc.DrawStringAnchored(text, float64(left), float64(space), 0, 0)
			lastImg = imaging.Rotate90(dc.Image())
		}
		var b bytes.Buffer
		// jpeg.Encode(&b, newImg, &jpeg.Options{Quality: 100})
		tiff.Encode(&b, lastImg, &tiff.Options{Compression: tiff.Deflate})

		// imaging.Encode(&b, newImg, imaging.TIFF)
		return b.Bytes(), err
	}
	// b := bytes.Buffer{}
	var b bytes.Buffer
	// jpeg.Encode(&b, newImg, &jpeg.Options{Quality: 100})
	tiff.Encode(&b, newImg, &tiff.Options{Compression: tiff.Deflate})

	// imaging.Encode(&b, newImg, imaging.TIFF)
	return b.Bytes(), err
}

func ImageToCmykLarge(buf []byte, width, height, span, lineWidth, totalHeight int, text string) ([]byte, error) {

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}
	if img.Bounds().Dx() > img.Bounds().Dy() {
		img = imaging.Rotate90(img)
	}
	left := img.Bounds().Dx()*span/width + lineWidth
	top := img.Bounds().Dy()*span/height + lineWidth
	maxY := img.Bounds().Dy() * totalHeight / height
	rtp := img.Bounds()

	canvasHeight := rtp.Dy() + 2*top
	rtp.Max.Y = maxY
	space := (maxY - canvasHeight) / 2
	rtp.Max.X = rtp.Dx() + 2*left
	newImg := image.NewCMYK(rtp)
	dx := rtp.Dx()
	dy := rtp.Dy()
	for x := 0; x <= dx; x++ {
		for y := space; y < space+canvasHeight; y++ {
			// if x < space || x > space+canvasWidth {
			// 	continue
			// }
			if x < left || x > dx-left || y < space+top || y > dy-space-top {
				if x < lineWidth || x > dx-lineWidth || y < space+lineWidth || y > dy-space-lineWidth {
					newImg.Set(x, y, color.Black)
				}
			} else {
				newImg.Set(x, y, img.At(x-left, y-top-space))
			}
		}
	}
	dc := gg.NewContextForImage(newImg)

	dc.SetColor(color.Black)

	err = dc.LoadFontFace("./resource/printfile/arialbd.ttf", 120)
	if err != nil {
		panic(err)
	}

	// dc.SetColor(color.White)
	dc.DrawStringAnchored(text, float64(dx/2), float64(space/2), 0, 0)

	lastImg := imaging.Rotate90(dc.Image())
	// b := bytes.Buffer{}
	var b bytes.Buffer
	// jpeg.Encode(&b, newImg, &jpeg.Options{Quality: 100})
	tiff.Encode(&b, lastImg, &tiff.Options{Compression: tiff.Deflate})

	// imaging.Encode(&b, newImg, imaging.TIFF)
	return b.Bytes(), err
}

func DrawScene(buf []byte) ([]bytes.Buffer, error) {
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	//195/314 810/579 271/88 650/911
	x, y := img.Bounds().Max.X, img.Bounds().Max.Y
	ox, oy := 272, 90
	imgX, imgY := 658, 921
	var layout = "v"
	if x > y {
		layout = "h"
		ox, oy = 190, 315
		imgX, imgY = 822, 587
	}
	var bs []bytes.Buffer
	var idx = []int{1, 2, 3}
	rand.Shuffle(len(idx), func(i, j int) {
		idx[i], idx[j] = idx[j], idx[i]
	})
	img = imaging.Resize(img, imgX, imgY, imaging.Box)
	for _, i := range idx {
		fileName := fmt.Sprintf("resource/printfile/scene_%s_%d.png", layout, i)

		sceneImg, err := gg.LoadImage(fileName)
		if err != nil {
			return nil, err
		}
		canvas := image.NewNRGBA(image.Rect(0, 0, sceneImg.Bounds().Dx(), sceneImg.Bounds().Dy()))
		canvas = imaging.Overlay(canvas, img, image.Point{X: ox, Y: oy}, 1)
		canvas = imaging.Overlay(canvas, sceneImg, image.Point{X: 0, Y: 0}, 1)
		var b bytes.Buffer
		jpeg.Encode(&b, canvas, &jpeg.Options{Quality: 90})
		bs = append(bs, b)
	}
	return bs, err
}
