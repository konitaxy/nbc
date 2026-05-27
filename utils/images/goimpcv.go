package images

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sort"

	"gitlab.com/ucard/global"
)

type Result struct {
	Path  string
	Dist  float64
	Dist2 float64
}

func FindTopNSimilarImages(buff []byte, key string, distKey []string, topN int) ([]Result, error) {
	srcPath := fmt.Sprintf("%s/%s.png", global.GVA_CONFIG.MachineLearning.ImageDir, key)
	var results []Result
	// var wg sync.WaitGroup
	for _, k := range distKey {
		// defer wg.Done()
		// wg.Add(1)
		distPath := fmt.Sprintf("%s/%s.png", global.GVA_CONFIG.MachineLearning.ImageDir, k)
		response, err := sendImagesToAPI(global.GVA_CONFIG.MachineLearning.MLUrl, srcPath, distPath)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return nil, err
		}

		// 解析响应
		var result map[string]interface{}
		err = json.Unmarshal(response, &result)
		if err != nil {
			fmt.Println("Error parsing response:", err)
			return nil, err
		}

		// 打印相似度分数
		similarity, ok := result["similarity"].(float64)
		if !ok {
			fmt.Println("Invalid similarity value in response")
			continue
		}
		if similarity > global.GVA_CONFIG.MachineLearning.MLSimilar {
			results = append(results, Result{
				Path: k,
				Dist: similarity,
			})
		}

	}
	// wg.Wait()
	sort.Slice(results, func(i, j int) bool {
		return results[i].Dist > results[j].Dist
	})
	if len(results) > topN {
		results = results[:topN]
	}
	return results, nil
}

func sendImagesToAPI(url string, image1Path string, image2Path string) ([]byte, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 添加第一个图片文件
	part1, err := writer.CreateFormFile("image1", image1Path)
	if err != nil {
		return nil, err
	}
	file1, err := os.Open(image1Path)
	if err != nil {
		return nil, err
	}
	defer file1.Close()
	_, err = io.Copy(part1, file1)
	if err != nil {
		return nil, err
	}

	// 添加第二个图片文件
	part2, err := writer.CreateFormFile("image2", image2Path)
	if err != nil {
		return nil, err
	}
	file2, err := os.Open(image2Path)
	if err != nil {
		return nil, err
	}
	defer file2.Close()
	_, err = io.Copy(part2, file2)
	if err != nil {
		return nil, err
	}

	// 关闭writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
