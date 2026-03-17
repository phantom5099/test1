package fs

import (
	"io"
	"os"
	"path/filepath"
)

// EnsureDir 确保给定路径的目录存在
func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

// PathExists 用于判断给定路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// BackupFile 为给定文件创建一个简单备份并返回备份路径
func BackupFile(path string) (string, error) {
	if !PathExists(path) {
		return "", nil
	}
	backup := path + ".bak"
	if err := copyFile(path, backup); err != nil {
		return "", err
	}
	return backup, nil
}

// WriteFileAtomic 使用临时文件进行原子写入
func WriteFileAtomic(path string, data []byte) error {
	if err := EnsureDir(path); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		// Fallback: try to overwrite directly
		if err2 := os.WriteFile(path, data, 0644); err2 != nil {
			return err2
		}
		_ = os.Remove(tmp)
	}
	return nil
}

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
