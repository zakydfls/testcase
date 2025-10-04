package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentStatus string
type DocumentAction string

const (
	StatusPending      DocumentStatus = "pending"
	StatusApproved     DocumentStatus = "approved"
	StatusRejected     DocumentStatus = "rejected"
	StatusNeedRevision DocumentStatus = "need_revision"
)

const (
	ActionApprove DocumentAction = "approve"
	ActionReject  DocumentAction = "reject"
)

type Document struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title           string         `gorm:"not null" json:"title"`
	Status          DocumentStatus `gorm:"default:'pending'" json:"status"`
	CurrentApprover int            `gorm:"default:1" json:"current_approver"`

	Approver1Action  *DocumentAction `json:"approver1_action"`
	Approver1Comment *string         `gorm:"type:text" json:"approver1_comment"`
	Approver1Date    *time.Time      `json:"approver1_date"`

	Approver2Action  *DocumentAction `json:"approver2_action"`
	Approver2Comment *string         `gorm:"type:text" json:"approver2_comment"`
	Approver2Date    *time.Time      `json:"approver2_date"`

	Approver3Action  *DocumentAction `json:"approver3_action"`
	Approver3Comment *string         `gorm:"type:text" json:"approver3_comment"`
	Approver3Date    *time.Time      `json:"approver3_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return
}

func (d *Document) TableName() string {
	return "documents"
}
