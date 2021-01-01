package logger

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRotateWriter(t *testing.T) {
	filePath := "/tmp"
	fileLog1 := filePath + "/1.log"
	fileLog2 := filePath + "/2.log"

	_ = os.Remove(fileLog1)
	_ = os.Remove(fileLog2)

	w := NewRotateWriter(filePath, "1")
	if w == nil {
		t.Fatalf("NewRotateWriter Fail")
	}
	// 写入内容
	_, e := w.Write([]byte("aaaa"))
	if e != nil {
		t.Fatalf("Write Fail: %s", e)
	}
	// 同文件进行 Rotate 并尝试写入
	e = w.Rotate("1")
	if e != nil {
		t.Fatalf("Rotate Fail: %s", e)
	}
	_, e = w.Write([]byte("aaaa"))
	if e != nil {
		t.Fatalf("Write Fail: %s", e)
	}
	// 更新文件名执行 Rotate 并尝试写入
	e = w.Rotate("2")
	if e != nil {
		t.Fatalf("Rotate Fail: %s", e)
	}
	_, e = w.Write([]byte("bbbb"))
	if e != nil {
		t.Fatalf("Write Fail: %s", e)
	}
	// 执行 Close
	ok := w.Close()
	if !ok {
		t.Fatalf("Close Fail: %s", e)
	}

	// 再次打开已经存在的日志斌执行写入
	w = NewRotateWriter(filePath, "2")
	_, e = w.Write([]byte("bbbb"))
	if e != nil {
		t.Fatalf("Write Fail: %s", e)
	}
	w.Close()

	// 验证结果
	log1, e := ioutil.ReadFile(fileLog1)
	log2, e := ioutil.ReadFile(fileLog2)
	if string(log1) != "aaaaaaaa" {
		t.Fatalf("Expect log1 'aaaaaaaa' but %s", string(log1))
	}
	if string(log2) != "bbbbbbbb" {
		t.Fatalf("Expect log2 'bbbbbbbb' but %s", string(log2))
	}
}
