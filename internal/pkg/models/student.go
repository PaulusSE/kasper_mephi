package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	// ID студента
	StudentID uuid.UUID `db:"students.student_id" json:"student_id,omitempty" format:"uuid"`
	// Полное имя
	FullName string `db:"students.full_name" json:"full_name,omitempty"`
	// Актуальный семестр
	ActualSemester int32 `db:"students.actual_semester" json:"actual_semester,omitempty"`
	// Количество лет обучения
	Years int32 `db:"students.years" json:"years,omitempty"`
	// Дата начала обучения
	StartDate time.Time `db:"students.start_date" json:"start_date"`
	// Статус обучения
	StudyingStatus string `db:"students.studying_status" json:"studying_status,omitempty" enums:"academic,graduated,studying,expelled"`
	// Статус проверки и подтверждения
	Status string `db:"students.status" json:"status,omitempty" enums:"todo,approved,on review,in progress,empty,failed"`
	// Специализация
	Specialization string `db:"specializations.title" json:"specialization,omitempty"`
	// Название группы
	GroupName string `db:"groups.group_name" json:"group_name,omitempty"`
	// Флаг о возможности редактировать всю информацию
	CanEdit bool   `db:"students.can_edit" json:"can_edit,omitempty"`
	Phone   string `db:"students.phone" json:"phone"`
	// Бюджетное или платное обучение
	Category string    `db:"students.category" json:"category"`
	EndDate  time.Time `db:"students.end_date" json:"end_date"`
}

type Supervisor struct {
	// ID научного руководителя
	SupervisorID uuid.UUID `db:"supervisor_id" json:"supervisor_id" format:"uuid"`
	// Полное имя руководителя
	FullName   string  `db:"full_name" json:"full_name"`
	Faculty    *string `db:"faculty" json:"faculty"`
	Department *string `db:"department" json:"department"`
	Degree     *string `db:"degree" json:"degree"`
	Phone      string  `db:"phone" json:"phone"`
	Archived   bool    `db:"archived" json:"archived"`
	Rank       *string `db:"rank" json:"rank"`
	Position   *string `db:"position" json:"position"`
}

type SupervisorProfile struct {
	// ID научного руководителя
	SupervisorID uuid.UUID `db:"supervisors.supervisor_id" json:"supervisor_id" format:"uuid"`
	// Полное имя руководителя
	FullName   string  `db:"supervisors.full_name" json:"full_name"`
	Faculty    *string `db:"supervisors.faculty" json:"faculty"`
	Department *string `db:"supervisors.department" json:"department"`
	Degree     *string `db:"supervisors.degree" json:"degree"`
	Email      string  `db:"users.email" json:"email"`
	Phone      string  `db:"supervisors.phone" json:"phone"`
	Archived   bool    `db:"archived" json:"archived"`
	Rank       *string `db:"rank" json:"rank"`
	Position   *string `db:"position" json:"position"`
}

type StudentProfile struct {
	// ID студента
	StudentID uuid.UUID `db:"students.student_id" json:"student_id,omitempty" format:"uuid"`
	// Полное имя
	FullName string `db:"students.full_name" json:"full_name,omitempty"`
	// Актуальный семестр
	ActualSemester int32 `db:"students.actual_semester" json:"actual_semester,omitempty"`
	// Количество лет обучения
	Years int32 `db:"students.years" json:"years,omitempty"`
	// Дата начала обучения
	StartDate time.Time `db:"students.start_date" json:"start_date"`
	// Статус обучения
	StudyingStatus string `db:"students.studying_status" json:"studying_status,omitempty" enums:"academic,graduated,studying,expelled"`
	// Статус проверки и подтверждения
	Status string `db:"students.status" json:"status,omitempty" enums:"todo,approved,on review,in progress,empty,failed"`
	// Специализация
	Specialization string `db:"specializations.title" json:"specialization,omitempty"`
	// Название группы
	GroupName string `db:"groups.group_name" json:"group_name,omitempty"`
	// Флаг о возможности редактировать всю информацию
	CanEdit bool   `db:"students.can_edit" json:"can_edit,omitempty"`
	Email   string `db:"users.email" json:"email"`
	Phone   string `db:"students.phone" json:"phone"`
	// Бюджетное или платное обучение
	Category string    `db:"students.category" json:"category"`
	EndDate  time.Time `db:"students.end_date" json:"end_date"`
}

type SupervisorFull struct {
	// ID научного руководителя
	SupervisorID uuid.UUID `db:"supervisor_id" json:"supervisor_id" format:"uuid"`
	// Полное имя руководителя
	FullName string `db:"full_name" json:"full_name"`
	Phone    string `db:"phone" json:"phone"`
	// Дата начала
	StartAt time.Time `db:"start_at" json:"start_at" format:"date-time"`
	// Дата окончания (пустое, если руководитель актуальный)
	EndAt *time.Time `db:"end_at" json:"end_at,omitempty" format:"date-time"`
}

type StudentSupervisorPair struct {
	// Информация о студенте в паре
	Student Student `json:"student"`
	// Информация о научном руководителе в паре
	Supervisor Supervisor `json:"supervisor"`
}

type UpdateProfile struct {
	// Полное имя
	FullName  string    `db:"students.full_name" json:"full_name,omitempty"`
	Email     string    `db:"users.email" json:"email"`
	GroupID   int32     `db:"groups.group_name" json:"group_id,omitempty"`
	Phone     string    `db:"students.phone" json:"phone"`
	Category  string    `db:"students.category" json:"category"`
	StartDate time.Time `db:"students.start_date" json:"date"`
	Years     int32     `db:"students.years" json:"years,omitempty"`
}

// Определим новую структуру для передачи данных в Python скрипт
type ReportData struct {
	CurrentSemester        int32                    `json:"current_semester"`
	FullName               string                   `json:"full_name"`
	SupervisorName         string                   `json:"supervisor_name"`
	EducationDirection     string                   `json:"education_direction"`
	EducationProfile       string                   `json:"education_profile"`
	EnrollmentDate         time.Time                `json:"enrollment_date"`
	Specialty              string                   `json:"specialty"`
	TrainingYearFGOS       string                   `json:"training_year_fgos"`
	CandidateExams         []Exam                   `json:"candidate_exams"`
	Category               string                   `json:"category"`
	Topic                  string                   `json:"topic"`
	ReportPeriodWork       string                   `json:"report_period_work"`
	ScientificObject       string                   `json:"scientific_obj"`
	ScientificSubject      string                   `json:"scientific_subj"`
	MentorRate             string                   `json:"mentor_rate"`
	ProgressPercents       []int                    `json:"progress_percents"`
	ProgressDescriptions   []string                 `json:"progress_descriptions"`
	Publications           []Publication            `json:"publications"`
	AllPublications        []Publication            `json:"all_publications"`
	PedagogicalData        []PedagogicalWork        `json:"pedagogical_data"`
	ReportOtherAchievments string                   `json:"report_other_achievments"`
	PedagogicalDataAll     []PedagogicalWorkSummary `json:"pedagogical_data_all"`
	NextSemesterPlan       []string                 `json:"next_semester_plan"`
}

// // Структура для экзаменов
// type Exam struct {
// 	ExamName string `json:"exam_name"`
// 	Mark     int    `json:"mark"`
// }

// Структура для публикаций
// type Publication struct {
// 	Name       string  `json:"name"`
// 	Status     string  `json:"status"`
// 	Impact     float64 `json:"impact"`
// 	OutputData string  `json:"output_data"`
// 	CoAuthors  string  `json:"co_authors"`
// 	Volume     int     `json:"volume"`
// }

// Структура для педагогической работы
type PedagogicalWork struct {
	Semester    int    `json:"semester"`
	WorkType    string `json:"work_type"`
	Hours       int    `json:"hours"`
	MainTeacher string `json:"main_teacher"`
	GroupName   string `json:"group_name"`
}

// Структура для обобщенной педагогической нагрузки
type PedagogicalWorkSummary struct {
	Semester              string `json:"semester"`
	ClassroomHours        int    `json:"classroom_hours"`
	OtherTeachingPractice int    `json:"other_teaching_practice"`
}
