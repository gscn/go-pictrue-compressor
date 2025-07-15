package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 支持的图片格式
var supportedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// 图片文件信息结构体
// 便于后续处理
type ImageFile struct {
	Path string
	Size int64
}

// 解析最小大小字符串（如1M, 5K, 2G）为字节数
func ParseSize(sizeStr string) int64 {
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))
	if sizeStr == "" || sizeStr == "0" {
		return 0
	}
	var mul int64 = 1
	if strings.HasSuffix(sizeStr, "K") {
		mul = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "K")
	} else if strings.HasSuffix(sizeStr, "M") {
		mul = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "M")
	} else if strings.HasSuffix(sizeStr, "G") {
		mul = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "G")
	}
	var val int64
	_, err := fmt.Sscanf(sizeStr, "%d", &val)
	if err != nil {
		return 0
	}
	return val * mul
}

// 递归扫描目录，返回符合条件的图片文件列表
func ScanImages(root string, minSize int64) ([]ImageFile, error) {
	var files []ImageFile
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无法访问的文件
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if !supportedExt[ext] {
			return nil
		}
		if info.Size() < minSize {
			return nil
		}
		files = append(files, ImageFile{
			Path: path,
			Size: info.Size(),
		})
		return nil
	})
	return files, err
}
