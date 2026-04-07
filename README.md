# RefKit: Enterprise-Grade Bioinformatics Reference Data Manager

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuanzhw/refkit)](https://goreportcard.com/report/github.com/yuanzhw/refkit)
[![Release](https://img.shields.io/github/v/release/yuanzhw/refkit)](https://github.com/yuanzhw/refkit/releases)

RefKit 是一个轻量级、零依赖的命令行工具，专为现代生物信息学工作流（WDL, Nextflow）设计，用于管理底层参考基因组和生信软件数据库。它通过严格的写权限控制和本地 SQLite 数据血缘追踪，确保了计算产线中“单一数据源”的绝对可靠性与跨服务器的无缝迁移。

## ✨ 核心特性 (Key Features)

* **零依赖部署 (Zero-Dependency):** 基于 Go 静态编译，单一二进制文件，无需配置 Python/Conda 环境，极速部署至各类计算节点或容器中。
* **数据血缘追踪 (Data Provenance):** 后端集成 SQLite，每一次数据的下载、校验和建库操作均被持久化记录。精确追溯每一个索引文件的构建命令、MD5 校验和与时间戳。
* **工作流引擎无缝对接 (Workflow Integration):** 告别硬编码。一键导出标准化的 WDL `inputs.json` 或 Nextflow `params.config`，自动处理相对路径与绝对路径的动态映射。
* **面向迁移设计 (Migration-Friendly):** 数据实体与元数据在物理层面强绑定。迁移集群或上云时，只需通过 Rsync 或 S3 Sync 同步单一根目录，环境瞬间恢复。

## 🚀 快速开始 (Quick Start)

### 1. 安装

直接下载预编译的二进制文件：

```bash
wget https://github.com/yuanzhw/refkit/releases/latest/download/refkit-linux-amd64 -O refkit
chmod +x refkit
sudo mv refkit /usr/local/bin/
```

### 2. 初始化环境

设置统一的参考数据根目录（建议使用专用账号管理权限）：

```bash
export REFKIT_ROOT="/data/refdb"
refkit init
```

### 3. 添加受控数据

所有进入系统的数据必须通过 `add` 命令并提供校验和来源：

```bash
refkit add \
  --type software_dbs \
  --name kraken2 \
  --version standard_202403 \
  --source "https://genome-idx.s3.amazonaws.com/kraken/k2_standard_202403.tar.gz" \
  --checksum "md5:88e...a1b"
```

### 4. 导出工作流配置

为 WDL 工作流动态生成输入 JSON：

```bash
refkit export --format wdl --genome GRCh38_Ensembl > inputs.json
```

## 📂 架构设计 (Architecture)

RefKit 强制采用物理隔离的目录树规范，将静态基因组序列与动态软件数据库分离：

```text
/data/refdb/
├── refkit_metadata.db        # SQLite 元数据库：管理所有路径与构建日志
├── genomes/                  # 极少变动：参考基因组及比对软件索引
│   └── Homo_sapiens/
│       └── GRCh38_Ensembl/
└── software_dbs/             # 频繁迭代：生信软件运行所需数据库
    └── kraken2/
        └── standard_202403/
```

## 🛠 开发指南 (Development)

本项目使用 Go 编写，命令行交互基于 Cobra 框架构建。欢迎提交 PR。

```bash
git clone https://github.com/yuanzhw/refkit.git
cd refkit
go build -o refkit main.go
```

## 📄 协议 (License)

本项目采用 [Apache License 2.0](LICENSE) 协议开源。详细信息请参阅 `LICENSE` 文件。
在遵循企业级合规标准的前提下，欢迎在生产和商业管线中自由集成与使用。