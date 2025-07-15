package cmd

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"go-pictrue-compress/internal/compress"
	"go-pictrue-compress/internal/logger"
	"go-pictrue-compress/internal/scanner"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	directory string
	maxWidth  int
	quality   int
	minSize   string
	threads   int
	logFile   string
)

var rootCmd = &cobra.Command{
	Use:   "go-pictrue-compressor",
	Short: "高效的图片批量压缩工具",
	Long:  `一个支持多格式、可配置参数的高效图片批量压缩命令行工具。`,
	Run: func(cmd *cobra.Command, args []string) {
		if directory == "" || maxWidth == 0 {
			fmt.Println("参数错误：--directory 和 --max-width 为必填项！")
			_ = cmd.Help()
			os.Exit(1)
		}
		if quality < 1 || quality > 100 {
			fmt.Println("参数错误：--quality 必须在1-100之间！")
			os.Exit(1)
		}

		minSizeBytes := scanner.ParseSize(minSize)
		fmt.Println("正在扫描目录，请稍候...")
		images, err := scanner.ScanImages(directory, minSizeBytes)
		if err != nil {
			fmt.Printf("扫描目录出错: %v\n", err)
			os.Exit(1)
		}
		totalSize := int64(0)
		for _, img := range images {
			totalSize += img.Size
		}
		fmt.Printf("找到图片文件: %d 个，总计 %.2f MB\n", len(images), float64(totalSize)/1024.0/1024.0)
		if len(images) == 0 {
			fmt.Println("没有找到符合条件的图片文件，程序结束。")
			return
		}

		if threads <= 0 {
			threads = runtime.NumCPU() / 2
			if threads < 1 {
				threads = 1
			}
		}
		fmt.Printf("开始处理（使用%d线程）...\n", threads)

		bar := progressbar.NewOptions(len(images),
			progressbar.OptionSetDescription("压缩进度"),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(20),
			progressbar.OptionSetPredictTime(true),
		)

		logMaxSize := int64(10 * 1024 * 1024) // 10MB
		logr, err := logger.NewLogger(logFile, logMaxSize)
		if err != nil {
			fmt.Printf("日志文件创建失败: %v\n", err)
			os.Exit(1)
		}
		defer logr.Close()
		logr.WriteHeader()

		results := make([]compress.Result, len(images))
		wg := sync.WaitGroup{}
		sem := make(chan struct{}, threads)

		start := time.Now()
		for i, img := range images {
			wg.Add(1)
			go func(i int, img scanner.ImageFile) {
				defer wg.Done()
				sem <- struct{}{}
				imgStart := time.Now()
				res := compress.CompressImage(img.Path, maxWidth, quality)
				res.OriginalPath = img.Path
				results[i] = res
				procTime := time.Since(imgStart).Milliseconds()
				status := "success"
				if res.Action == "failed" && res.Err != nil {
					status = res.Err.Error()
				}
				logr.Write([]string{
					time.Now().Format(time.RFC3339),
					img.Path,
					fmt.Sprintf("%d", res.OriginalSize),
					fmt.Sprintf("%d", res.NewSize),
					res.Action,
					fmt.Sprintf("%d", procTime),
					status,
				})
				bar.Add(1)
				<-sem
			}(i, img)
		}
		wg.Wait()
		bar.Finish()
		used := time.Since(start)

		// 统计结果
		var replaced, skipped, failed int
		var saved int64
		for _, r := range results {
			if r.Action == "replaced" {
				replaced++
				saved += r.OriginalSize - r.NewSize
			} else if r.Action == "skipped" {
				skipped++
			} else if r.Action == "failed" {
				failed++
			}
		}
		fmt.Printf("处理完成: 替换原文件: %d, 跳过: %d, 失败: %d\n", replaced, skipped, failed)
		fmt.Printf("节省空间: %.2f MB, 总耗时: %.2fs\n", float64(saved)/1024.0/1024.0, used.Seconds())
		fmt.Printf("详细日志见: %s\n", logFile)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "要处理的目录路径 (必填)")
	rootCmd.PersistentFlags().IntVarP(&maxWidth, "max-width", "w", 0, "最大宽度（像素，必填）")
	rootCmd.PersistentFlags().IntVarP(&quality, "quality", "q", 75, "压缩质量 (1-100)")
	rootCmd.PersistentFlags().StringVarP(&minSize, "min-size", "s", "0", "最小处理大小（如5M, 1G）")
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 0, "最大并发线程数 (默认CPU核心数/2)")
	rootCmd.PersistentFlags().StringVarP(&logFile, "log-file", "l", "./go-pictrue-compressor.csv", "日志文件路径")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
