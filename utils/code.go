package utils

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/skip2/go-qrcode"
)

func GenerateCouponCode(length int) string {
	// 可选的字符集
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	var result strings.Builder

	// 设置随机数种子
	rand.NewSource(time.Now().UnixNano())

	// 随机生成指定长度的字符串
	for i := 0; i < length; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}

	return result.String()
}

func GenerateRandomChars(length int) string {
	// 可选的字符集
	const chars = "abcdefghijklmnopqrstuvwxyz"
	var result strings.Builder

	// 设置随机数种子
	rand.NewSource(time.Now().UnixNano())

	// 随机生成指定长度的字符串
	for i := 0; i < length; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}

	return strings.ToUpper(result.String())
}

func GenerateRandomClientID(length int) string {
	// 可选的字符集
	const chars = "abcdefghijklmnopqrstuvwxyz"
	var result strings.Builder

	// 设置随机数种子
	rand.NewSource(time.Now().UnixNano())

	// 随机生成指定长度的字符串
	for i := 0; i < length; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}

	return strings.ToUpper(result.String())
}

func GenerateRandomNumber(length int) string {
	// 可选的字符集
	const chars = "0123456789"
	var result strings.Builder

	// 设置随机数种子
	rand.NewSource(time.Now().UnixNano())

	// 随机生成指定长度的字符串
	for i := 0; i < length; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}
	return result.String()
}

func Now() *time.Time {
	curr := time.Now()
	return &curr
}
func AllLettersAreUpper(s string) bool {
	hasLetter := false
	for _, c := range s {
		if unicode.IsLetter(c) {
			if !unicode.IsUpper(c) {
				return false
			}
			hasLetter = true
		}
	}
	return hasLetter // 如果没有字母，可以根据需求决定是否返回 true
}
func ExtraAllLetters(s string) (ns string) {

	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			ns += string(c)
		}
	}
	return ns // 如果没有字母，可以根据需求决定是否返回 true
}

func HasUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func GenerateQRCodeBase64(content string, size int) (string, error) {
	qrBytes, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(qrBytes)

	// 添加 Data URL 前缀（可选，用于前端直接使用）
	dataURL := "data:image/png;base64," + base64Str
	return dataURL, nil

	// return base64Str, nil
}
