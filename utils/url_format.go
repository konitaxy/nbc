package utils

import (
	"fmt"
	"math/rand"
	"net/url"
	"path"
)

var (
	TEXTURE            = "Texture"
	HDGLOSSY           = "HD_Glossy"
	Artwork_Url_Format = "%s-%08d-%03d"
)

func ArtworkUrlTextureBuilder(artworkId uint) string {

	return ArtworkUrlBuilder(artworkId, TEXTURE)
}
func ArtworkUrlHDGlossyBuilder(artworkId uint) string {

	return ArtworkUrlBuilder(artworkId, HDGLOSSY)
}

func ArtworkUrlBuilder(artworkId uint, t string) string {
	return fmt.Sprintf(Artwork_Url_Format, t, artworkId, rand.Intn(1000))
}

func ExtGet(urlStr string) (string, error) {

	// 解析 URL
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}

	// 获取路径部分
	pathPart := u.Path

	// 获取文件名
	fileName := path.Base(pathPart)

	// 获取文件扩展名
	fileExt := path.Ext(fileName)

	return fileExt, nil
}
