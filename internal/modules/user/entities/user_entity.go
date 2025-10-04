package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleEnum string

const (
	RoleAdmin1 RoleEnum = "admin1"
	RoleAdmin2 RoleEnum = "admin2"
	RoleAdmin3 RoleEnum = "admin3"
	RoleUser   RoleEnum = "user"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=2,max=255"`
	Username  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username" validate:"required,alphanum,min=3,max=100"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-" validate:"required,min=8"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Role      RoleEnum       `gorm:"type:varchar(50);not null;default:'user'" json:"role" validate:"required,oneof=admin user manager"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
