package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// 定义全局变量，用于接收全局 flag 的值
var dbPath string

// rootCmd 代表基础命令
var rootCmd = &cobra.Command{
	Use:   "refkit",
	Short: "生信参考数据库与软件数据库管理工具",
	Long: `RefKit 是一个企业级的本地参考数据治理工具。
它可以帮助你规范化下载、校验、并追踪参考基因组及生信软件数据库的版本记录，
同时支持动态导出供 WDL 和 Nextflow 使用的配置文件。`,
	// 如果没有子命令，默认打印帮助信息
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute 提供给 main.go 调用
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init() 是 Go 的特殊函数，会在包加载时自动执行
func init() {
	// 定义全局 Flag。这里的 dbPath 之前要加 & 取地址，是因为要把命令行接收到的值塞进这个变量的内存里
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "/data/refdb/refkit_metadata.db", "SQLite 数据库的绝对路径")
}
