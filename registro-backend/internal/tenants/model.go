package tenants

import (
	"time"
)

type TenantQuota struct {
	MaxStudents int   `json:"max_students"`
	MaxTeachers int   `json:"max_teachers"`
	MaxStorageMB int64 `json:"max_storage_mb"`
}

type Tenant struct {
	ID          string      `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Code        string      `json:"code" db:"code"`
	Status      string      `json:"status" db:"status"` // 'active', 'suspended'
	Quota       TenantQuota `json:"quota"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	MaxStudents int    `json:"max_students"`
	MaxTeachers int    `json:"max_teachers"`
	MaxStorageMB int64 `json:"max_storage_mb"`
}
