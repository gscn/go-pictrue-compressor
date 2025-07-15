package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	file    *os.File
	writer  *csv.Writer
	mutex   sync.Mutex
	maxSize int64
	logPath string
}

// 新建日志，maxSize单位字节，超过则轮转
func NewLogger(path string, maxSize int64) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)
	return &Logger{file: file, writer: writer, maxSize: maxSize, logPath: path}, nil
}

func (l *Logger) RotateIfNeeded() error {
	info, err := l.file.Stat()
	if err != nil {
		return err
	}
	if info.Size() < l.maxSize {
		return nil
	}
	l.writer.Flush()
	l.file.Close()
	backup := fmt.Sprintf("%s.%s", l.logPath, time.Now().Format("20060102_150405"))
	os.Rename(l.logPath, backup)
	file, err := os.OpenFile(l.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	l.file = file
	l.writer = csv.NewWriter(file)
	return nil
}

func (l *Logger) Write(record []string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if err := l.RotateIfNeeded(); err != nil {
		return err
	}
	if err := l.writer.Write(record); err != nil {
		return err
	}
	l.writer.Flush()
	return nil
}

func (l *Logger) Close() {
	l.writer.Flush()
	l.file.Close()
}

// 日志头部
func (l *Logger) WriteHeader() {
	header := []string{"timestamp", "file_path", "original_size", "processed_size", "action_taken", "processing_time", "status"}
	l.Write(header)
}
