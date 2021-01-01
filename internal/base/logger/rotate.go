package logger

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
)

type RotateWriter struct {
	lock     sync.Mutex
	path     string // should be set to the actual filename
	writer   *bufio.Writer
	file     *os.File
	filename string
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func NewRotateWriter(path string, filename string) *RotateWriter {
	w := &RotateWriter{path: path, lock: sync.Mutex{}}
	err := w.Rotate(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return w
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.writer.Write(output)
}

func (w *RotateWriter) WriteString(output string) (int, error) {
	return w.Write([]byte(output))
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) Rotate(name string) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// 如果传入的日期是一致的，不需要执行操作
	if w.filename == name {
		// 执行一下 Flush 把日志刷盘
		err = w.writer.Flush()
		return err
	}

	filename := path.Join(w.path, name+".log")
	_, err = os.Stat(filename)

	// 如果日期为空或者 writer 不存在，表明是新创建
	if w.writer != nil {
		w.Close()
	}
	if err != nil {
		// 文件不存在
		w.file, err = os.Create(filename)
	} else {
		w.file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}
	w.filename = name
	w.writer = bufio.NewWriter(w.file)
	return err
}

func (w *RotateWriter) Close() bool {
	// 关闭文件流
	var err error
	if w.writer != nil && w.file != nil {
		err = w.writer.Flush()
		err = w.file.Close()
		w.writer = nil
		w.file = nil

	}
	if err != nil {
		return false
	}
	return true
}
