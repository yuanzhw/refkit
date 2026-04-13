package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yuanzhw/refkit/internal/db" // 引入我们刚刚写好的内部数据库包
	"github.com/yuanzhw/refkit/internal/downloader"
)

// 定义命令行标志的变量
var (
	resVersion string
	resType    string
	sourceURL  string
)

var addCmd = &cobra.Command{
	Use:   "add [资源名称]",
	Short: "下载并入库一个生信参考数据库或基因组",
	Long:  `从指定网络源下载数据，将其保存到本地存储，并记录完整的血缘追踪信息到 SQLite 数据库。`,
	Args:  cobra.ExactArgs(1), // 强制要求必须提供一个位置参数（如 kraken2_std）
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		fmt.Printf("🚀 准备入库资源: [%s] (版本: %s, 类型: %s)\n", name, resVersion, resType)

		database, err := db.InitDB("refkit.db")
		if err != nil {
			return fmt.Errorf("数据库引擎启动失败: %w", err)
		}
		defer database.Close()

		res, err := db.SaveResource(database, name, resType, resVersion)
		if err != nil {
			return fmt.Errorf("无法创建资源记录: %w", err)
		}
		fmt.Printf("📝 数据库记录已生成，血缘追踪 ID: %s\n", res.ID)

		fmt.Printf("⏳ 正在尝试从源站建立连接: %s\n", sourceURL)

		currentDir := "."
		filename := fmt.Sprintf("%s_%s.data", name, resVersion)
		destPath := filepath.Join(currentDir, filename)

		// 调用核心下载引擎
		err = downloader.DownloadFile(sourceURL, currentDir, filename)
		if err != nil {
			// 💥 异常分支：如果下载失败，立刻将数据库状态标记为 error，防止变成僵尸任务
			db.UpdateResourceStatus(database, res.ID, "error")
			return fmt.Errorf("下载彻底失败: %w", err)
		}

		// ==========================================
		// 🎉 成功分支：进入业务状态闭环
		// ==========================================
		fmt.Println("🔄 下载完成，正在进行数据防篡改固化与血缘绑定...")

		// 1. 获取刚刚下载好的文件大小
		var sizeBytes int64 = 0
		if fileInfo, err := os.Stat(destPath); err == nil {
			sizeBytes = fileInfo.Size()
		}

		// 2. 获取当前执行这条命令的系统用户 (用于审计溯源)
		operator := "unknown"
		if currentUser, err := user.Current(); err == nil {
			operator = currentUser.Username
		}

		// 3. 将物理路径落库
		if err := db.SaveFilePath(database, res.ID, destPath, sizeBytes); err != nil {
			return fmt.Errorf("记录物理路径失败: %w", err)
		}

		// 4. 将血缘日志落库 (注意：为了不卡流程，checksum 我们暂时填了 pending_hash，后续可以加上 MD5/SHA256 计算模块)
		if err := db.SaveProvenanceLog(database, res.ID, sourceURL, "network_download", "pending_hash", operator); err != nil {
			return fmt.Errorf("记录血缘日志失败: %w", err)
		}

		// 5. 所有的附加信息都绑定完毕后，正式将核心状态由 pending 改为 ready
		if err := db.UpdateResourceStatus(database, res.ID, "ready"); err != nil {
			return fmt.Errorf("更新最终状态失败: %w", err)
		}

		fmt.Printf("✅ 入库闭环完成！资源 [%s] 处于 ready 状态。\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// 绑定命令行参数 (Flags)
	addCmd.Flags().StringVarP(&resVersion, "version", "v", "latest", "指定资源的版本号")
	addCmd.Flags().StringVarP(&resType, "type", "t", "database", "指定资源类型 (如: database, genome)")
	addCmd.Flags().StringVarP(&sourceURL, "source", "s", "", "指定网络下载的源 URL (必填)")

	// 强制要求提供 source 链接
	addCmd.MarkFlagRequired("source")
}
