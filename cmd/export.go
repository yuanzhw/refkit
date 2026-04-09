package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	targetTarget string // 基因组名称或软件库名称
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出供 WDL/Nextflow 使用的路径配置",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: 从 SQLite 查询该靶标的真实相对路径，并结合环境变量拼接

		if exportFormat == "wdl" {
			// 模拟输出 WDL inputs.json 格式
			jsonOutput := fmt.Sprintf(`{
  "pipeline.ref_fasta": "/data/refdb/genomes/Homo_sapiens/%s/seq/genome.fasta",
  "pipeline.ref_dict": "/data/refdb/genomes/Homo_sapiens/%s/seq/genome.dict"
}`, targetTarget, targetTarget)
			fmt.Println(jsonOutput)

		} else if exportFormat == "nextflow" {
			// 模拟输出 Nextflow params 格式
			nfOutput := fmt.Sprintf("params.ref_fasta = '/data/refdb/genomes/Homo_sapiens/%s/seq/genome.fasta'", targetTarget)
			fmt.Println(nfOutput)

		} else {
			fmt.Println("❌ 错误：不支持的格式。仅支持 'wdl' 或 'nextflow'。")
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "wdl", "输出格式 (wdl/nextflow)")
	exportCmd.Flags().StringVarP(&targetTarget, "target", "t", "", "目标标识符 (如 GRCh38_Ensembl)")

	exportCmd.MarkFlagRequired("target")
}
