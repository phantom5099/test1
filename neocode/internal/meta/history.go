package meta

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// HistoryFilePath 记录历史的文件路径（当前在仓库根目录下的隐藏历史文件）
func HistoryFilePath() string {
	// 使用工作目录下的隐藏文件，便于沙盒环境下的简单持久化
	return filepath.Join(".", ".neocode_history.json")
}

// LoadHistory 读取历史条目
func LoadHistory() ([]string, error) {
	path := HistoryFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []string{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var hist []string
	if err := json.Unmarshal(data, &hist); err != nil {
		return nil, err
	}
	return hist, nil
}

// SaveHistory 将历史条目写入磁盘
func SaveHistory(hist []string) error {
	path := HistoryFilePath()
	data, err := json.MarshalIndent(hist, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// AppendHistory 将一条历史记录追加到历史中
func AppendHistory(entry string) error {
	hist, err := LoadHistory()
	if err != nil {
		return err
	}
	hist = append(hist, entry)
	return SaveHistory(hist)
}
