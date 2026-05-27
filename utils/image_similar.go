package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/utils/images"
	"go.uber.org/zap"
)

var Bins = []int{60, 4, 4}
var HISTOGRAMS = make(map[string][]float64)
var PHASH = make(map[string]*goimagehash.ImageHash)

// rgbToHSV converts an RGB color to HSV.
func rgbToHSV(r, g, b float64) (float64, float64, float64) {
	rf := r / 255.0
	gf := g / 255.0
	bf := b / 255.0

	maxc := math.Max(math.Max(rf, gf), bf)
	minc := math.Min(math.Min(rf, gf), bf)
	v := maxc

	delta := maxc - minc
	if maxc != 0 {
		s := delta / maxc
		if rf == maxc {
			return math.Mod(60*(gf-bf)/delta+360, 360), s, v
		} else if gf == maxc {
			return math.Mod(60*(bf-rf)/delta+120, 360), s, v
		} else {
			return math.Mod(60*(rf-gf)/delta+240, 360), s, v
		}
	}

	return 0, 0, v
}

// extractHSVHistogram extracts HSV histogram from an image.
func extractHSVHistogram(img image.Image, bins []int) []float64 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	histogram := make([]float64, bins[0]*bins[1]*bins[2])

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			hr, sr, vr := rgbToHSV(float64(r>>8), float64(g>>8), float64(b>>8))

			hBin := int(math.Min(float64(bins[0]-1), hr/360*float64(bins[0])))
			sBin := int(math.Min(float64(bins[1]-1), sr*float64(bins[1])))
			vBin := int(math.Min(float64(bins[2]-1), vr*float64(bins[2])))

			index := hBin*bins[1]*bins[2] + sBin*bins[2] + vBin
			histogram[index]++
		}
	}

	// Normalize the histogram
	totalPixels := float64(width * height)
	for i := range histogram {
		histogram[i] /= totalPixels
	}

	return histogram
}

var mu sync.Mutex

// loadImagesFromDir loads all images from a directory and computes their HSV histograms.
func LoadImagesFromDir(dir string) (map[string][]float64, error) {
	histograms := make(map[string][]float64)
	i := 0
	var wg sync.WaitGroup
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png") {
			img, err := imaging.Open(path)
			if err != nil {
				return err
			}
			wg.Add(1)
			go func(info os.FileInfo) {
				defer wg.Done()
				hist := extractHSVHistogram(img, Bins)
				mu.Lock()
				histograms[strings.Split(info.Name(), ".")[0]] = hist
				mu.Unlock()
				i++
				fmt.Println(i)
			}(info)

		}
		return nil
	})
	wg.Wait()
	global.GVA_LOG.Sugar().Infof("success extract [%d] histograms done", i)
	return histograms, err
}

func LoadImagesFromDirUsePHash(dir string) (map[string]string, error) {
	hashDict := make(map[string]string)
	i := 0
	var wg sync.WaitGroup
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png") {
			img, err := imaging.Open(path)
			if err != nil {
				return err
			}
			wg.Add(1)
			go func(info os.FileInfo) {
				defer wg.Done()
				phash, _ := goimagehash.PerceptionHash(img)
				hashDict[strings.Split(info.Name(), ".")[0]] = phash.ToString()
				i++
			}(info)

		}
		return nil
	})
	wg.Wait()
	global.GVA_LOG.Sugar().Infof("success extract [%d] phash done", i)

	return hashDict, err
}

// LoadHashFromFile 从文件中加载哈希值到map[string]*goimagehash.ImageHash
func LoadHashFromFile(filename string) (map[string]*goimagehash.ImageHash, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hashDict := make(map[string]*goimagehash.ImageHash)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			continue
		}
		hashValue := parts[0]
		imagePath := parts[1]

		// 将字符串形式的哈希值转换回*goimagehash.ImageHash类型
		hash, err := goimagehash.ImageHashFromString(hashValue)
		if err != nil {
			continue
		}

		hashDict[imagePath] = hash
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	PHASH = hashDict
	global.GVA_LOG.Info(fmt.Sprintf("成功读取文件并解析PHASH数据 %d条数据", len(hashDict)))

	return hashDict, nil
}

// 保存 hash 到文件
func SaveHashToFile(hashMap map[string]string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for path, hash := range hashMap {
		fmt.Fprintf(file, "%s %s\n", hash, path)
	}
	return err
}
func AddHashToFile(buff []byte, key string) error {
	// Step 1: 解码图片
	if _, e := PHASH[key]; e {
		global.GVA_LOG.Info("exit key:" + key)
	}
	img, _, err := image.Decode(bytes.NewReader(buff))
	if err != nil {
		return fmt.Errorf("解码图片失败: %v", err)
	}

	hash, _ := goimagehash.PerceptionHash(img)

	file, err := os.OpenFile("log/phash.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Fprintf(file, "%s %s\n", hash.ToString(), key)
	global.GVA_LOG.Info("写入文件成功", zap.String("file", "log/phash.csv"), zap.String("key", key), zap.String("hash", hash.ToString()))
	PHASH[key] = hash
	return err
}

func SaveHistogramsToFile(histograms map[string][]float64, outputPath string) error {
	os.Remove(outputPath)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	for path, hist := range histograms {
		// 在每一行开始处添加文件名
		_, err := fmt.Fprintf(file, "%s:", path)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}

		// 写入直方图数据
		for _, val := range hist {
			_, err := fmt.Fprintf(file, " %.6f", val)
			if err != nil {
				return fmt.Errorf("failed to write to file: %v", err)
			}
		}
		// 换行以便下一个文件的数据从新行开始
		_, err = fmt.Fprintln(file)
		if err != nil {
			return fmt.Errorf("failed to write newline to file: %v", err)
		}
	}

	return nil
}
func AppendHistogramsToFile(key string, hist []float64, outputPath string) error {
	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// 在每一行开始处添加文件名
	_, err = fmt.Fprintf(file, "%s:", key)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	// 写入直方图数据
	for _, val := range hist {
		_, err := fmt.Fprintf(file, " %.6f", val)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	// 换行以便下一个文件的数据从新行开始
	_, err = fmt.Fprintln(file)
	if err != nil {
		return fmt.Errorf("failed to write newline to file: %v", err)
	}

	return nil
}

// loadHistogramsFromFile 从文件中加载 HSV 直方图数据
func LoadHistogramsFromFile(filePath string) (map[string][]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	histograms := make(map[string][]float64)
	scanner := bufio.NewScanner(file)

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// 按冒号分割文件名和特征值
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("第 %d 行格式错误: 缺少冒号分隔符", lineNum)
		}

		filename := strings.Split(parts[0], ".")[0]
		valuesStr := strings.Fields(strings.TrimSpace(parts[1]))

		// 将字符串切片转换为 float64 切片
		features := make([]float64, len(valuesStr))
		for i, s := range valuesStr {
			val, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, fmt.Errorf("第 %d 行解析浮点数失败: %v", lineNum, err)
			}
			features[i] = val
		}

		histograms[filename] = features
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时出错: %v", err)
	}

	HISTOGRAMS = histograms
	global.GVA_LOG.Info(fmt.Sprintf("成功读取文件并解析直方图数据 %d条数据", len(histograms)))
	return histograms, nil
}

// euclideanDistance 计算两个直方图之间的欧氏距离
func euclideanDistance(a, b []float64) float64 {
	var sum float64
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// findTopNSimilarImages 根据输入的图片字节流查找最相似的 Top N 图片

func FindTopNSimilarImagesUsePHash(
	key string,
	phash map[string]*goimagehash.ImageHash,
	distince int,
	topN int,
) ([]images.Result, error) {

	// // Step 1: 解码图片
	// img, _, err := image.Decode(bytes.NewReader(imgBytes))
	// if err != nil {
	// 	return nil, fmt.Errorf("解码图片失败: %v", err)
	// }

	// // Step 2: 缩放 + 提取直方图
	// img = imaging.Resize(img, 512, 512, imaging.Lanczos)

	// Step 3: 找出最相似的 Top N 张图片

	var results []images.Result
	if phash == nil {
		phash = PHASH
	}
	targetHist := phash[key]
	if targetHist == nil {
		return results, nil
	}
	for path, hash := range phash {
		d, _ := targetHist.Distance(hash)
		if d > distince {
			continue
		}
		results = append(results, images.Result{
			Path: path,
			Dist: float64(d),
		})
	}

	// 排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Dist < results[j].Dist
	})

	return results, nil
}

func FindTopNSimilarImages(
	key string,
	dbHistograms map[string][]float64,
	topN int,
	dist float64,
	bins []int,
) ([]images.Result, error) {

	targetHist := HISTOGRAMS[key]

	// Step 3: 找出最相似的 Top N 张图片

	var results []images.Result
	if dbHistograms == nil {
		dbHistograms = HISTOGRAMS
	}
	for path, hist := range dbHistograms {
		d := euclideanDistance(targetHist, hist)
		if d > dist {
			continue
		}
		results = append(results, images.Result{
			Path: path,
			Dist: d,
		})
	}

	// 排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Dist < results[j].Dist
	})

	return results, nil
}

func AddHistograms(buf []byte, key string) {
	img, _ := imaging.Decode(bytes.NewReader(buf))
	hist := extractHSVHistogram(img, Bins)
	if _, e := HISTOGRAMS[key]; !e {
		HISTOGRAMS[key] = hist
		AppendHistogramsToFile(key, hist, "log/histograms.csv")
		global.GVA_LOG.Info(fmt.Sprintf("Add Histogram [%s]", key))
	}
}
