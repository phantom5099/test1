package main

import (
	"github.com/yourname/neocode/config"
	"github.com/yourname/neocode/internal/repl"
	"log"
)

// 程序入口：启动 neocode 的交互式 REPL
func main() {
	cfg := config.LoadConfig()
	if err := repl.Run(cfg); err != nil {
		log.Fatalf("neocode 运行时错误: %v", err)
	}
}
