package db

import (
	"database/sql"
	_ "embed" // 用于将 schema.sql 嵌入二进制
	"fmt"

	"github.com/google/uuid" // 用于生成资源 ID
	_ "modernc.org/sqlite"   // 纯 Go 实现的 SQLite 驱动
)

//go:embed schema.sql
var ddlSchema string

// ==========================================
// 1. 数据模型定义
// ==========================================

// Resource 代表数据库中的一条核心资产记录
type Resource struct {
	ID      string
	Name    string
	Type    string
	Version string
	Status  string
}

// ==========================================
// 2. 数据库初始化与操作
// ==========================================

// InitDB 初始化数据库并创建表
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 强制开启外键支持（SQLite 默认关闭级联删除等功能）
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	// 执行嵌入的 DDL 语句
	if _, err := db.Exec(ddlSchema); err != nil {
		return nil, err
	}

	return db, nil
}

// SaveResource 创建并保存一条全新的资源记录
func SaveResource(db *sql.DB, name string, resType string, version string) (*Resource, error) {
	// 生成全局唯一的 ID
	id := uuid.New().String()
	initialStatus := "pending"

	query := `
		INSERT INTO resources (id, name, type, version, status)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, id, name, resType, version, initialStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to insert resource: %w", err)
	}

	return &Resource{
		ID:      id,
		Name:    name,
		Type:    resType,
		Version: version,
		Status:  initialStatus,
	}, nil
}

// ==========================================
// 3. 状态更新与闭环记录
// ==========================================

// UpdateResourceStatus 修改资源的核心状态 (如: ready, error)
func UpdateResourceStatus(db *sql.DB, id string, status string) error {
	query := `UPDATE resources SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, status, id)
	return err
}

// SaveFilePath 记录物理文件在硬盘上的位置和大小
func SaveFilePath(db *sql.DB, resourceID string, relativePath string, sizeBytes int64) error {
	query := `INSERT INTO file_paths (resource_id, relative_path, size_bytes) VALUES (?, ?, ?)`
	_, err := db.Exec(query, resourceID, relativePath, sizeBytes)
	return err
}

// SaveProvenanceLog 写入极其关键的“血缘溯源”记录
func SaveProvenanceLog(db *sql.DB, resourceID string, source string, method string, checksum string, operator string) error {
	query := `INSERT INTO provenance_log (resource_id, original_source, ingest_method, checksum, operator) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, resourceID, source, method, checksum, operator)
	return err
}
