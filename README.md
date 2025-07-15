# go-pictrue-compress

## 项目简介

`go-pictrue-compress` 是一款基于Go语言开发的高效图片批量压缩命令行工具，支持 JPEG/PNG 格式，适用于大规模图片目录的自动化压缩与空间优化。

## 功能特性
- 支持 JPEG、PNG 图片批量压缩
- 可设置最大宽度、压缩质量、最小处理文件大小
- 仅压缩后变小才覆盖原文件，安全高效
- 多线程并发处理，充分利用多核CPU
- 实时进度条反馈
- 详细CSV日志记录，支持自动轮转
- 命令行参数灵活，适合自动化脚本集成

## 安装方法
1. 安装 Go 1.18 及以上版本
2. 克隆本项目并进入目录：
   ```sh
   git clone https://github.com/yourname/go-pictrue-compress.git
   cd go-pictrue-compress
   ```
3. 拉取依赖并编译：
   ```sh
   go mod tidy
   go build -o go-pictrue-compressor main.go
   ```

## 使用示例
```sh
./go-pictrue-compressor -d ./images -w 1920 -q 85 -s 1M -t 4 -l ./go-pictrue-compressor.csv
```

## 参数说明
| 参数 | 缩写 | 必选 | 类型 | 描述 | 默认值 |
|------|------|------|------|------|--------|
| --directory | -d | 是 | 字符串 | 要处理的目录路径 | 无 |
| --max-width | -w | 是 | 整数 | 最大宽度（像素） | 无 |
| --quality | -q | 否 | 整数(1-100) | 压缩质量 | 75 |
| --min-size | -s | 否 | 字符串 | 最小处理大小（如5M, 1G） | 0 |
| --threads | -t | 否 | 整数 | 最大并发线程数 | CPU核心数/2 |
| --log-file | -l | 否 | 字符串 | 日志文件路径 | ./go-pictrue-compressor.csv |
| --help | -h | 否 | 无 | 显示帮助信息 | false |

## 日志说明
- 日志为CSV格式，包含：时间戳、文件路径、原始大小、压缩后大小、操作类型、耗时、状态。
- 日志文件超过10MB自动轮转，保留历史日志。

## 常见问题
- 仅支持JPEG/PNG图片，其他格式自动跳过。
- 仅当压缩后文件变小时才会覆盖原文件。
- 需保证有写入权限，否则会记录失败。

## 许可证
MIT License 