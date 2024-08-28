package models

// import (
// 	"time"
// )

// // ReportData описывает структуру данных для генерации отчета,
// // интегрируя существующие модели из диссертационной работы
// type ReportData struct {
// 	FullName                string                       `json:"full_name"`
// 	SupervisorName          string                       `json:"supervisor_name"`
// 	EducationDirection      string                       `json:"education_direction"`
// 	EducationProfile        string                       `json:"education_profile"`
// 	EnrollmentDate          time.Time                    `json:"enrollment_date"`
// 	Specialty               string                       `json:"specialty"`
// 	TrainingYearFGOS        int                          `json:"training_year_fgos"`
// 	CandidateExams          []string                     `json:"candidate_exams"`
// 	Category                string                       `json:"category"`
// 	Topic                   string                       `json:"topic"`
// 	ReportPeriodWork        string                       `json:"report_period_work"`
// 	ScientificObj           string                       `json:"scientific_obj"`
// 	ScientificSubj          string                       `json:"scientific_subj"`
// 	MentorRate              int                          `json:"mentor_rate"`
// 	ProgressPercents        float64                      `json:"progress_percents"`
// 	ProgressDescriptions    []string                     `json:"progress_descriptions"`
// 	Publications            []Publication                `json:"publications"`
// 	AllPublications         []Publication                `json:"all_publications"`
// 	PedagogicalData         []ClassroomLoad              `json:"pedagogical_data"`
// 	ReportOtherAchievements []string                     `json:"report_other_achievements"`
// 	PedagogicalDataAll      []TeachingLoad               `json:"pedagogical_data_all"`
// 	NextSemesterPlan        string                       `json:"next_semester_plan"`
// 	SemesterProgress        []SemesterProgressResponse   `json:"semester_progress"`
// 	Supervisors             []SupervisorFull             `json:"supervisors"`
// 	StudentsComments        []StudentComment             `json:"students_comments"`
// 	DissertationTitles      []DissertationTitlesResponse `json:"dissertation_titles"`
// 	Feedback                []FeedbackResponse           `json:"feedback"`
// }

// Добавлены следующие поля для отображения более сложной структуры данных:
// - SemesterProgress - Прогресс по семестрам.
// - Supervisors - Информация о научных руководителях.
// - StudentsComments - Комментарии студентов.
// - DissertationTitles - Названия диссертаций.
// - Feedback - Обратная связь.
