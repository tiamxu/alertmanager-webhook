package database

import (
	"fmt"

	"github.com/tiamxu/kit/sql"
)

// DB 全局数据库连接
var DB *sql.DB

// Config 数据库配置
type Config = sql.Config

// Init 初始化数据库连接和表结构
func Init(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("database config is nil")
	}

	db := sql.NewPreDB()
	if err := db.Init(cfg); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	DB = db.DB

	// 创建告警记录表
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS alert_records (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		alertname VARCHAR(255) NOT NULL,
		level VARCHAR(50) NOT NULL,
		status VARCHAR(50) NOT NULL,
		labels TEXT NOT NULL,
		annotations TEXT NOT NULL,
		instance VARCHAR(255),
		startsAt VARCHAR(50),
		endsAt VARCHAR(50),
		summary TEXT,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	if err != nil {
		return fmt.Errorf("failed to create alert_records table: %w", err)
	}

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
