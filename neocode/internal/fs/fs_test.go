package fs

import (
	"os"
	"testing"
)

func TestEnsureDirAndWriteAtomic(t *testing.T) {
	base := t.TempDir()
	path := base + string(os.PathSeparator) + "sub" + string(os.PathSeparator) + "file.txt"
	if err := EnsureDir(path); err != nil {
		t.Fatalf("EnsureDir 失败: %v", err)
	}
	data := []byte("content")
	if err := WriteFileAtomic(path, data); err != nil {
		t.Fatalf("WriteFileAtomic 失败: %v", err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile 失败: %v", err)
	}
	if string(got) != string(data) {
		t.Fatalf("内容不匹配，期望 %q，得到 %q", data, got)
	}
	// 测试备份功能
	if _, err := BackupFile(path); err != nil {
		t.Fatalf("BackupFile 失败: %v", err)
	}
}
