package compress

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// 处理结果
type Result struct {
	OriginalPath string
	OriginalSize int64
	NewSize      int64
	Action       string // compressed, skipped, replaced, failed
	Err          error
}

// 压缩图片主函数
func CompressImage(path string, maxWidth, quality int) Result {
	file, err := os.Open(path)
	if err != nil {
		return Result{OriginalPath: path, Action: "failed", Err: err}
	}
	defer file.Close()

	info, _ := file.Stat()
	origSize := info.Size()

	ext := strings.ToLower(filepath.Ext(path))
	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return Result{OriginalPath: path, Action: "skipped", Err: errors.New("不支持的图片格式")}
	}
	if err != nil {
		return Result{OriginalPath: path, Action: "failed", Err: err}
	}

	// 等比缩放
	w := img.Bounds().Dx()
	if w > maxWidth {
		img = resize.Resize(uint(maxWidth), 0, img, resize.Lanczos3)
	}

	buf := &bytes.Buffer{}
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	case ".png":
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		err = enc.Encode(buf, img)
	}
	if err != nil {
		return Result{OriginalPath: path, Action: "failed", Err: err}
	}

	newBytes := buf.Bytes()
	if int64(len(newBytes)) >= origSize {
		return Result{OriginalPath: path, OriginalSize: origSize, NewSize: int64(len(newBytes)), Action: "skipped", Err: nil}
	}

	err = ioutil.WriteFile(path, newBytes, info.Mode())
	if err != nil {
		return Result{OriginalPath: path, Action: "failed", Err: err}
	}

	return Result{OriginalPath: path, OriginalSize: origSize, NewSize: int64(len(newBytes)), Action: "replaced", Err: nil}
}
