package utils

import (
	"regexp"
	"strings"
)

// ParseUserAgent 解析 User-Agent 并返回操作系统、浏览器、设备类型等信息
func ParseUserAgent(ua string) map[string]string {
	result := map[string]string{
		"os":       "Unknown",
		"os_name":  "Unknown",
		"os_ver":   "",
		"browser":  "Unknown",
		"device":   "Unknown",
		"platform": "Unknown",
	}

	// 小写化 UA 便于匹配
	lowerUA := strings.ToLower(ua)

	// === 操作系统识别 ===
	if strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") || strings.Contains(lowerUA, "ios") {
		// iOS 设备
		result["os_name"] = "iOS"
		re := regexp.MustCompile(`OS (\d+_?\d*_?\d*)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			osVer := strings.ReplaceAll(matches[1], "_", ".")
			result["os_ver"] = osVer
			result["os"] = "iOS " + osVer
		} else {
			result["os"] = "iOS"
		}
		result["device"] = "Mobile"
		result["platform"] = "iPhone/iPad"

	} else if strings.Contains(lowerUA, "android") {
		// Android
		result["os_name"] = "Android"
		re := regexp.MustCompile(`Android[ /]([\d.]+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			result["os_ver"] = matches[1]
			result["os"] = "Android " + matches[1]
		} else {
			result["os"] = "Android"
		}
		result["device"] = "Mobile"
		result["platform"] = "Android"

	} else if strings.Contains(lowerUA, "mac os x") {
		// macOS (非 iOS)
		result["os_name"] = "macOS"
		re := regexp.MustCompile(`Mac OS X (\d+_\d+_\d+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			osVer := strings.ReplaceAll(matches[1], "_", ".")
			result["os_ver"] = osVer
			result["os"] = "macOS " + osVer
		} else {
			result["os"] = "macOS"
		}
		result["device"] = "Desktop"
		result["platform"] = "Macintosh"

	} else if strings.Contains(lowerUA, "windows") {
		result["os_name"] = "Windows"
		re := regexp.MustCompile(`Windows NT (\d+\.\d+)`)
		if matches := re.FindStringSubmatch(lowerUA); len(matches) > 1 {
			ntVer := matches[1]
			switch ntVer {
			case "10.0":
				result["os_ver"] = "10"
				result["os"] = "Windows 10"
			case "6.3":
				result["os_ver"] = "8.1"
				result["os"] = "Windows 8.1"
			case "6.2":
				result["os"] = "Windows 8"
			case "6.1":
				result["os"] = "Windows 7"
			default:
				result["os"] = "Windows (NT " + ntVer + ")"
			}
		} else {
			result["os"] = "Windows"
		}
		result["device"] = "Desktop"
		result["platform"] = "Windows PC"

	} else if strings.Contains(lowerUA, "linux") {
		result["os_name"] = "Linux"
		result["os"] = "Linux"
		result["device"] = "Desktop"
		if strings.Contains(lowerUA, "mobile") || strings.Contains(lowerUA, "android") {
			result["device"] = "Mobile"
		}
		result["platform"] = "Linux"
	}

	// === 浏览器识别 ===
	if strings.Contains(lowerUA, "wechat") {
		re := regexp.MustCompile(`MicroMessenger/([\d.]+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			result["browser"] = "WeChat " + matches[1]
		} else {
			result["browser"] = "WeChat"
		}
	} else if strings.Contains(lowerUA, "qqbrowser") {
		re := regexp.MustCompile(`QQBrowser/([\d.]+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			result["browser"] = "QQ Browser " + matches[1]
		} else {
			result["browser"] = "QQ Browser"
		}
	} else if strings.Contains(lowerUA, "edg/") {
		re := regexp.MustCompile(`Edg/([\d.]+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			result["browser"] = "Edge " + matches[1]
		} else {
			result["browser"] = "Edge"
		}
	} else if strings.Contains(lowerUA, "chrome/") {
		re := regexp.MustCompile(`Chrome/([\d.]+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			result["browser"] = "Chrome " + matches[1]
		} else {
			result["browser"] = "Chrome"
		}
	} else if strings.Contains(lowerUA, "firefox/") {
		re := regexp.MustCompile(`Firefox/([\d.]+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			result["browser"] = "Firefox " + matches[1]
		} else {
			result["browser"] = "Firefox"
		}
	} else if strings.Contains(lowerUA, "safari") && !strings.Contains(lowerUA, "chrome") {
		// 注意：Safari 的 UA 往往也包含 Chrome，所以要排除
		if strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") {
			re := regexp.MustCompile(`Version/([\d.]+)`)
			if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
				result["browser"] = "Safari " + matches[1]
			} else {
				result["browser"] = "Safari"
			}
		} else {
			re := regexp.MustCompile(`Safari/([\d.]+)`)
			if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
				result["browser"] = "Safari " + matches[1]
			} else {
				result["browser"] = "Safari"
			}
		}
	} else if strings.Contains(lowerUA, "msie") || strings.Contains(lowerUA, "trident") {
		if strings.Contains(lowerUA, "trident") {
			result["browser"] = "Internet Explorer"
		} else {
			re := regexp.MustCompile(`MSIE (\d+\.\d+)`)
			if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
				result["browser"] = "Internet Explorer " + matches[1]
			} else {
				result["browser"] = "Internet Explorer"
			}
		}
	}

	return result
}

// func main() {
// 	userAgents := []string{
// 		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36",
// 		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
// 		"Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
// 		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36",
// 		"Mozilla/5.0 (Linux; U; Android 10; zh-cn; MI 9 Build/QKQ1.190825.002) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/79.0.3945.147 Mobile Safari/537.36 XiaoMi/MiuiBrowser/13.5.13",
// 		"Mozilla/5.0 (iPad; CPU OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/141.0.6167.82 Mobile/15E148 Safari/604.1",
// 		"Mozilla/5.0 (Linux; Android 14; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Mobile Safari/537.36 EdgA/141.0.2272.118",
// 		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
// 	}

// 	for i, ua := range userAgents {
// 		fmt.Printf("\n--- 示例 %d ---\n", i+1)
// 		info := ParseUserAgent(ua)
// 		fmt.Printf("UA: %s\n", ua)
// 		fmt.Printf("操作系统: %s\n", info["os"])
// 		fmt.Printf("浏览器: %s\n", info["browser"])
// 		fmt.Printf("设备类型: %s\n", info["device"])
// 		fmt.Printf("平台: %s\n", info["platform"])
// 	}
// }
