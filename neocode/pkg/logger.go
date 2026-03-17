package logger

import (
	"log"
)

// 初始化日志设置
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
