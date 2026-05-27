package utils

import (
	"fmt"
	"math/rand"
)

func RandomTitle() string {

	// 定义一组关键词
	keywords := []string{"科技", "创新", "未来", "人工智能", "编程", "开源", "互联网", "大数据", "区块链"}

	// 定义多个标题模板
	titles := []string{
		"%s %s %s",
		"探索 %s 与 %s",
		"%s 技术的发展",
		"%s 在 %s 的应用",
	}

	// 随机选择一个模板
	templateIndex := rand.Intn(len(titles))
	titleTemplate := titles[templateIndex]

	// 随机选择几个关键词
	numKeywords := rand.Intn(len(keywords)) + 3 // 至少选择3个关键词
	selectedKeywords := make([]interface{}, numKeywords)

	for i := 0; i < numKeywords; i++ {
		index := rand.Intn(len(keywords))
		selectedKeywords[i] = keywords[index]
	}

	// 组合成标题
	return fmt.Sprintf(titleTemplate, selectedKeywords...)
}
