package utils

import (
	"bytes"
)

func GetImageType(buf []byte) string {
	// 检测图片格式
	var format string
	switch {
	case bytes.HasPrefix(buf, []byte{0xFF, 0xD8, 0xFF}): // JPEG
		format = "jpeg"
	case bytes.HasPrefix(buf, []byte{0x89, 0x50, 0x4E, 0x47}): // PNG
		format = "png"
	case bytes.HasPrefix(buf, []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}): // GIF
		format = "gif"
	case bytes.HasPrefix(buf, []byte{0x42, 0x4D}): // BMP
		format = "bmp"
	default:
		format = "unknown"
	}
	return format
}
