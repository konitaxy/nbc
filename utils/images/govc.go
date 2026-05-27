package images

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"strings"
// 	"sync"

// 	"github.com/r0busta/go-shopify-graphql-model/v4/graph/model"
// 	"gitlab.com/ucard/global"
// 	"gocv.io/x/gocv"
// )

// var (
// 	detector            = gocv.NewSIFTWithParams(model.NewInt(200), nil, model.NewFloat64(0.08), model.NewFloat64(15), nil)
// 	wg                  sync.WaitGroup
// 	matcher             = gocv.NewFlannBasedMatcher()
// 	rows, cols, matType = 1, 128, gocv.MatTypeCV32F
// 	descMap             = make(map[string]gocv.Mat)
// 	ImgDir              = "log/images"
// 	descriptorPath      = "log/descriptors"
// )

// func extractAndSaveDescriptors(imagePath, descriptorPath string) ([]gocv.KeyPoint, gocv.Mat) {

// 	img := gocv.IMRead(imagePath, gocv.IMReadGrayScale)
// 	defer img.Close()
// 	keypoints, descriptors := detector.DetectAndCompute(img, gocv.NewMat())

// 	// 使用自定义方式保存描述符为二进制文件
// 	data := descriptors.ToBytes()
// 	err := os.WriteFile(descriptorPath, data, 0644)
// 	if err != nil {
// 		fmt.Println("Error saving descriptors:", err)
// 	}
// 	return keypoints, descriptors
// }
// func loadDescriptors(descriptorPath string, rows, cols int, matType gocv.MatType) (gocv.Mat, error) {
// 	data, _ := os.ReadFile(descriptorPath)

// 	return gocv.NewMatFromBytes(rows, cols, matType, data)
// }
// func matchFeaturesFLANN(queryDescriptors, trainDescriptors gocv.Mat) [][]gocv.DMatch {
// 	result := matcher.KnnMatch(queryDescriptors, trainDescriptors, 2)
// 	return result
// }
// func extractAndSaveDescriptorsFromBuff(buff []byte, descriptorPath string, save bool) ([]gocv.KeyPoint, gocv.Mat) {
// 	img, _ := gocv.IMDecode(buff, gocv.IMReadGrayScale)
// 	defer img.Close()
// 	keypoints, descriptors := detector.DetectAndCompute(img, gocv.NewMat())

// 	// 使用自定义方式保存描述符为二进制文件
// 	data := descriptors.ToBytes()
// 	if !save {
// 		return keypoints, descriptors
// 	}
// 	os.WriteFile(descriptorPath, data, 0644)
// 	return keypoints, descriptors
// }

// func ExtraAllDescriptors() error {
// 	i := 0
// 	os.MkdirAll(descriptorPath, os.ModePerm)
// 	// var wg sync.WaitGroup
// 	err := filepath.Walk(ImgDir, func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png") {
// 			key := strings.Split(info.Name(), ".")[0]
// 			if _, err := os.Stat(fmt.Sprintf("%s/%s.npy", descriptorPath, key)); err == nil {
// 				return nil
// 			}
// 			_, _ = extractAndSaveDescriptors(path, fmt.Sprintf("%s/%s.npy", descriptorPath, key))
// 			i++
// 			fmt.Println(i)
// 			// }(info)

// 		}
// 		return nil
// 	})
// 	return err
// }

// func LoadAllDescriptors(keys []string) error {
// 	var i = 0
// 	if len(descMap) > 1000 {
// 		return nil
// 	}
// 	err := filepath.Walk(descriptorPath, func(path string, info os.FileInfo, err error) error {

// 		// defer wg.Done()
// 		if !info.IsDir() && (filepath.Ext(path) == ".npy") {
// 			exit := false
// 			for _, k := range keys {
// 				if strings.Contains(info.Name(), k) {
// 					exit = true
// 					break
// 				}
// 			}
// 			if exit {
// 				key := strings.Split(strings.Split(info.Name(), ".")[0], "_")[0]
// 				mat, err := loadDescriptors(path, rows, cols, matType)
// 				if err == nil {
// 					descMap[key] = mat
// 					i++
// 					fmt.Println(i)
// 				}

// 			}

// 			// fmt.Println(i)
// 		}

// 		return nil
// 	})
// 	// wg.Wait()
// 	global.GVA_LOG.Sugar().Infof("加载描述符成功%d 特征", len(descMap))

// 	return err
// }

// type Result struct {
// 	Path  string
// 	Dist  float64
// 	Dist2 float64
// }

// func FindTopNSimilarImages(buff []byte, key string, distKey []string, topN int) ([]Result, error) {
// 	var results []Result
// 	dist, err := loadDescriptors(fmt.Sprintf("%s/%s.npy", descriptorPath, key), rows, cols, matType)
// 	if err != nil {
// 		global.GVA_LOG.Sugar().Infof("Load描述符%s不存在", key)
// 	}
// 	for _, k := range distKey {

// 		v, err := loadDescriptors(fmt.Sprintf("%s/%s.npy", descriptorPath, k), rows, cols, matType)
// 		if err != nil {
// 			global.GVA_LOG.Sugar().Infof("Load描述符%s不存在", k)
// 		}
// 		if v.Empty() || dist.Empty() {
// 			global.GVA_LOG.Sugar().Infof("描述符%s不存在", k)
// 			continue
// 		}
// 		result := matchFeaturesFLANN(dist, v)
// 		for _, matches := range result {

// 			if len(matches) >= 2 && matches[0].Distance < 0.9*matches[1].Distance {
// 				results = append(results, Result{
// 					Path:  k,
// 					Dist:  matches[0].Distance,
// 					Dist2: 0.7 * matches[1].Distance,
// 				})
// 				// }

// 			}

// 		}

// 	}
// 	sort.Slice(results, func(i, j int) bool {
// 		return results[i].Dist < results[j].Dist
// 	})
// 	if len(results) > topN {
// 		results = results[:topN]
// 	}
// 	return results, nil
// }
