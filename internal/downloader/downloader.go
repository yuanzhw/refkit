package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// DownloadFile 带有防屏蔽和进度条的核心下载引擎
// 参数: sourceURL (目标网址), destDir (保存目录), filename (保存的文件名)
func DownloadFile(sourceURL string, destDir string, filename string) error {
	// 1. 确保目标目录存在 (比如 ~/.refkit/data/)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("无法创建目录 %s: %w", destDir, err)
	}

	// 拼装完整的本地文件路径
	destPath := filepath.Join(destDir, filename)

	// 2. 伪装准备：不要直接 GET，而是构建一个自定义 Request
	req, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		return fmt.Errorf("构建请求失败: %w", err)
	}

	// ⚠️ 核心防屏蔽魔法：把自己伪装成最新版的 Chrome 浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	// 告诉服务器我们接受断点续传（为以后的功能打底）
	req.Header.Set("Accept", "*/*")

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("网络请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码是否正常
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误状态码: %d %s", resp.StatusCode, resp.Status)
	}

	// 4. 创建本地文件准备接收数据
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("无法创建本地文件: %w", err)
	}
	defer out.Close()

	// 5. 渲染极其优雅的进度条
	// resp.ContentLength 会告诉我们文件总大小，如果服务器没给，它就是 -1
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"⬇️ 正在下载",
	)

	// 6. 数据流转魔法 (io.MultiWriter)
	// 将网络数据流 (resp.Body) 一边写入硬盘文件 (out)，一边写入进度条更新器 (bar)
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件时发生中断: %w", err)
	}

	fmt.Printf("\n✅ 文件已成功保存至: %s\n", destPath)
	return nil
}
