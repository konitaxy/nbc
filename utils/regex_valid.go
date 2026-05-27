package utils

import (
	"regexp"
	"strconv"
	"unicode"
)

func IsValidEmail(email string) bool {
	// 正则表达式模式
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 检查是否匹配
	return re.MatchString(email)
}
func IsValidPassword(password string) bool {
	// 正则表达式模式
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	var hasLetter, hasNumber bool

	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsSpace(char):
			return false
		}
	}

	return hasLetter && hasNumber
}
func IsNumeric(str string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(str)
}
func IsValidName(name string) bool {
	// 正则表达式模式
	pattern := `^[a-zA-Z0-9._ ]*$`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 检查是否匹配
	return re.MatchString(name)
}

func IsValidatePassword(password string) bool {
	// 定义密码校验的正则表达式
	passwordPattern := `^(?=.*[A-Za-z])(?=.*\d)[^\s]{8,20}$`

	// 编译正则表达式
	re := regexp.MustCompile(passwordPattern)

	// 检查密码是否符合正则表达式
	return re.MatchString(password)
}

type CodeParts struct {
	ArtworkID   uint
	ProductType string
	Size        string
	Layout      string
}

func ParseSkuAndValidateCode(code string) (CodeParts, bool) {
	// 定义正则表达式模式
	pattern := `^MP-(\d+)-([a-zA-Z]+)-([a-zA-Z]+)-([a-zA-Z]+)$`
	re := regexp.MustCompile(pattern)

	// 尝试匹配输入字符串
	matches := re.FindStringSubmatch(code)

	if matches == nil {
		return CodeParts{}, false
	}

	u, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return CodeParts{}, false
	}
	// 如果匹配成功，则提取各部分
	var parts CodeParts
	if len(matches) == 5 {
		parts = CodeParts{
			ArtworkID:   uint(u),
			ProductType: matches[2],
			Size:        matches[3],
			Layout:      matches[4],
		}
	}

	return parts, true
}
