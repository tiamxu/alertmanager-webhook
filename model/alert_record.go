package model

import (
	"encoding/json"
	"time"

	"github.com/tiamxu/kit/sql"
)

// Alert 告警记录

type AlertRecord struct {
	Alertname   string            `db:"alertname" json:"alertname"`
	Level       string            `db:"level" json:"level"`
	Status      string            `db:"status" json:"status"`
	Labels      map[string]string `db:"labels" json:"labels"`
	Annotations map[string]string `db:"annotations" json:"annotations"`
	Instance    string            `db:"instance" json:"instance"`
	StartsAt    string            `db:"startsAt" json:"startsAt"`
	EndsAt      string            `db:"endsAt" json:"endsAt"`
	Summary     string            `db:"Summary" json:"Summary"`
	Description string            `db:"description" json:"description"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
}

// 告警状态常量
const (
	AlertStatusFiring    = "firing"    // 告警触发
	AlertStatusHandling  = "handling"  // 处理中
	AlertStatusResolved  = "resolved"  // 已解决
	AlertStatusEscalated = "escalated" // 已升级
)

// AlertRecordRepo 告警记录仓库
type AlertRecordRepo struct {
	db *sql.DB
}

// NewAlertRecordRepo 创建告警记录仓库实例
func NewAlertRecordRepo() *AlertRecordRepo {
	return &AlertRecordRepo{db: DB}
}

// Create 创建告警记录
func (r *AlertRecordRepo) Create(alert *AlertRecord) error {
	query := `INSERT INTO alert_records (
		alertname, level, status, labels, annotations, 
		instance, startsAt, endsAt, Summary, description,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	labels, _ := json.Marshal(alert.Labels)
	annotations, _ := json.Marshal(alert.Annotations)

	_, err := r.db.Exec(query,
		alert.Alertname,
		alert.Level,
		alert.Status,
		string(labels),
		string(annotations),
		alert.Instance,
		alert.StartsAt,
		alert.EndsAt,
		alert.Summary,
		alert.Description,
	)
	return err
}
