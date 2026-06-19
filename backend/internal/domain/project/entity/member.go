package entity

import "time"

type ProjectRole string

const (
	RoleProjectManager ProjectRole = "PROJECT_MANAGER"
	RoleMember         ProjectRole = "MEMBER"
	RoleObserver       ProjectRole = "OBSERVER"
)

type MemberStatus string

const (
	MemberActive    MemberStatus = "ACTIVE"
	MemberSuspended MemberStatus = "SUSPENDED"
	MemberLeft      MemberStatus = "LEFT"
	MemberRemoved   MemberStatus = "REMOVED"
)

type ProjectMember struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	ProjectID uint64 `gorm:"not null" json:"project_id"`
	UserID    uint64 `gorm:"not null" json:"user_id"`

	ProjectRole ProjectRole  `gorm:"not null" json:"project_role"`
	Status      MemberStatus `gorm:"not null;default:ACTIVE" json:"status"`

	JoinedAt time.Time  `json:"joined_at"`
	LeftAt   *time.Time `json:"left_at,omitempty"`

	StatusChangedBy *uint64    `json:"status_changed_by,omitempty"`
	StatusChangedAt *time.Time `json:"status_changed_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}
