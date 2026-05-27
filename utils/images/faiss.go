package images

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"gitlab.com/ucard/global"
)

type GenerateVectorsRequest struct {
	Directory string `json:"directory"`
}

type FindSimilarImagesRequest struct {
	QueryImageName string `json:"query_image_name"`
	N              int    `json:"N,omitempty"`
}

type Response struct {
	Message       string   `json:"message,omitempty"`
	Error         string   `json:"error,omitempty"`
	SimilarImages []Result `json:"similar_images,omitempty"`
}

func callPythonEndpoint(endpoint string, body *bytes.Buffer, contentType string) (Response, error) {
	// proxyURL, _ := url.Parse("http://127.0.0.1:7890")

	// // 创建自定义 Transport
	// transport := &http.Transport{
	// 	Proxy: func(req *http.Request) (*url.URL, error) {
	// 		return proxyURL, nil
	// 	},
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }

	// // 创建 client 并使用代理 transport
	// client := &http.Client{
	// 	Transport: transport,
	// }
	resp, err := http.Post(fmt.Sprintf("%s%s", global.GVA_CONFIG.MachineLearning.MLUrl, endpoint), contentType, body)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func GenerateFaissVectors() (string, error) {
	req := GenerateVectorsRequest{
		Directory: global.GVA_CONFIG.MachineLearning.ImageDir,
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	response, err := callPythonEndpoint("/generate_vectors", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return "", err
	}

	return response.Message, err
}

func FindSimilarImagesByClip(key string, topN int) ([]Result, error) {
	var req = FindSimilarImagesRequest{
		N:              topN,
		QueryImageName: key,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := callPythonEndpoint("/find_similar_images", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}

	return response.SimilarImages, err
}
func AddNewImageToFaiss(buf []byte, key string) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 添加第一个图片文件

	part, err := writer.CreateFormFile("image", fmt.Sprintf("%s.png", key))
	if err != nil {
		return err
	}
	filepath := fmt.Sprintf("%s/%s.png", global.GVA_CONFIG.MachineLearning.ImageDir, key)
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	// 关闭writer
	err = writer.Close()
	if err != nil {
		return err
	}
	_, err = callPythonEndpoint("/add_new_image", body, writer.FormDataContentType())
	if err != nil {
		return err
	}

	return nil
}
