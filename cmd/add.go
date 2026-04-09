package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 定义当前命令专用的局部变量
var (
	dataType   string
	dataName   string
	dataVer    string
	dataSource string
	checksum   string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加并溯源新的参考数据",
	Long:  `下载、校验并规范化存放数据，同时将数据血缘信息写入 SQLite 数据库。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 这里的逻辑会在用户输入 refkit add ... 后执行
		fmt.Println("🚀 正在初始化入库流程...")
		fmt.Printf("📂 使用数据库: %s\n", dbPath) // 可以直接读取 root.go 中定义的全局变量
		fmt.Printf("📦 数据信息: [%s] %s (版本: %s)\n", dataType, dataName, dataVer)

		if dataSource != "" {
			fmt.Printf("⬇️  下载来源: %s\n", dataSource)
		}
		if checksum != "" {
			fmt.Printf("🛡️  预期校验: %s\n", checksum)
		}

		// TODO: 1. 执行下载
		// TODO: 2. 流式计算 MD5 并对比
		// TODO: 3. 解压/移动到物理目录
		// TODO: 4. Insert into SQLite

		fmt.Println("✅ 数据已成功添加并记录溯源信息！")
	},
}

func init() {
	// 将 add 命令注册为 root 命令的子命令
	rootCmd.AddCommand(addCmd)

	// 绑定 Flags
	addCmd.Flags().StringVarP(&dataType, "type", "t", "", "数据类别: genomes 或 software_dbs (必填)")
	addCmd.Flags().StringVarP(&dataName, "name", "n", "", "数据名称，如 kraken2 (必填)")
	addCmd.Flags().StringVarP(&dataVer, "version", "v", "", "版本号，如 standard_202403 (必填)")
	addCmd.Flags().StringVarP(&dataSource, "source", "s", "", "下载链接或本地绝对路径")
	addCmd.Flags().StringVarP(&checksum, "checksum", "c", "", "MD5校验码，格式为 md5:xxx")

	// 标记必填参数，Cobra 会在入口处自动拦截遗漏操作
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("version")
}
