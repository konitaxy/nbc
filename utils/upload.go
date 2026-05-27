package utils

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ucard/global"
)

// func GenS3Url(url string) string {
// 	return fmt.Sprintf("https://%s.%s", global.GVA_CONFIG.AwsS3.Bucket, url)
// }

// func GetImageFileType(path string) string {

// }

// func IsImageFile(f file) (string, err) {
// 	cfg, err := jpeg.DecodeConfig(file)
// }

func GetS3Key(str string) string {
	u, _ := url.Parse(str)
	return u.Path[1:]
}
func GetS3KeyWithoutPrefix(str, prefix string) string {
	if prefix == "" {
		prefix = global.GVA_CONFIG.AwsS3.PathPrefix
	}

	u, _ := url.Parse(str)
	key := u.Path[1:]
	for {
		if strings.HasPrefix(key, prefix) {
			key = strings.Replace(key, prefix, "", 1)[1:]
			continue
		}
		break
	}
	return key
}

func UrlAddTs(str string) string {
	u, _ := url.Parse(str)
	u.RawQuery = "ts=" + strconv.FormatInt(time.Now().Unix(), 10)
	return u.String()
}
