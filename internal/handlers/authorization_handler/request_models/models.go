package request_models

import (
	"github.com/google/uuid"
)

type FirstStudentRegistry struct {
	//Department       string     `json:"department,omitempty"`
	Email            string     `json:"email"            binding:"required,email"`
	FullName         string     `json:"full_name" binding:"required"`
	Password         string     `json:"password"         binding:"required,min=8"`
	SpecializationID int32      `json:"specialization_id" binding:"required"`
	ActualSemester   int32      `json:"actual_semester" binding:"required"`
	NumberOfYears    int32      `json:"number_of_years" binding:"required"`
	StartDate        string     `json:"start_date" binding:"required"`
	GroupID          int32      `json:"group_number" binding:"required"`
	SupervisorID     *uuid.UUID `json:"supervisor_id" binding:"required"`
	Phone            string     `json:"phone" binding:"required"`
	Category         string     `json:"category" binding:"required"` // Бюджетное или платное обучение
}

type ChangePasswordRequest struct {
	OldPassword string
	NewPassword string
}

type FirstSupervisorRegistry struct {
	Email      *string `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=8"`
	FullName   string  `json:"full_name" binding:"required"`
	Phone      string  `json:"phone" binding:"required"`
	Faculty    string  `json:"faculty,omitempty"`
	Department string  `json:"department,omitempty"`
	Degree     string  `json:"degree,omitempty"`
	Rank       *string `json:"rank,omitempty"`
	Position   *string `json:"position,omitempty"`
}
