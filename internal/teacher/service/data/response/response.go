package response

import (
	"github.com/google/uuid"
)

type GetTeacherResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreateAt   int64     `json:"created_at"`
	UpdateAt   int64     `json:"updated_at"`
	Subjects   []Subject `json:"subjects"`
	Classes    []Class   `json:"classes"`
}

type TeacherStatistics struct {
	TotalTeachers    int `json:"total_teachers"`
	VerifiedTeachers int `json:"verified_teachers"`
	PendingTeachers  int `json:"pending_teachers"`
	ActiveTeachers   int `json:"active_teachers"`
}

type GetListTeacherResponse []GetTeacherResponse

type GetDetailTeacherResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreateAt   int64     `json:"created_at"`
	UpdateAt   int64     `json:"updated_at"`
}

type Subject struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SchoolID  uuid.UUID `json:"school_id"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}
