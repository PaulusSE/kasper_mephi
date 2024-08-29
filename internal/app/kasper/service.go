package kasper

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type (
	StudentHandler interface {
		AllToReview(ctx *gin.Context)
		GetStudentStatus(ctx *gin.Context)

		GetDissertationPage(ctx *gin.Context)
		UpsertSemesterProgress(ctx *gin.Context)
		DownloadDissertation(ctx *gin.Context)
		UploadDissertation(ctx *gin.Context)
		// DownloadPresentation(ctx *gin.Context)
		DissertationTitleToReview(ctx *gin.Context)
		UpsertDissertationTitle(ctx *gin.Context)
		DissertationToReview(ctx *gin.Context)

		GetTeachingLoad(ctx *gin.Context)
		UpsertAdditionalLoads(ctx *gin.Context)
		UpsertClassroomLoads(ctx *gin.Context)
		UpsertIndividualLoads(ctx *gin.Context)
		DeleteAdditionalLoads(ctx *gin.Context)
		DeleteClassroomLoads(ctx *gin.Context)
		DeleteIndividualLoads(ctx *gin.Context)
		TeachingLoadToReview(ctx *gin.Context)

		GetScientificWorks(ctx *gin.Context)
		UpsertResearchProjects(ctx *gin.Context)
		UpsertPublications(ctx *gin.Context)
		UpsertConferences(ctx *gin.Context)
		DeleteProjects(ctx *gin.Context)
		DeletePublications(ctx *gin.Context)
		DeleteConferences(ctx *gin.Context)
		ScientificWorksToReview(ctx *gin.Context)

		GetSpecializations(ctx *gin.Context)
		GetGroups(ctx *gin.Context)

		GetSupervisors(ctx *gin.Context)
		UpdateProgressiveness(ctx *gin.Context)
		GetStudentProfile(ctx *gin.Context)

		GetReportComments(ctx *gin.Context)
		UpsertReportComments(ctx *gin.Context)

		GetAllMarks(ctx *gin.Context)
		UpsertExamResults(ctx *gin.Context)

		GetSemesterAmounts(ctx *gin.Context)

		UpsertPatents(ctx *gin.Context)
		DeletePatents(ctx *gin.Context)

		UpdateProfile(ctx *gin.Context)
		DeleteMarks(ctx *gin.Context)
	}

	SupervisorHandler interface {
		GetStudentsList(ctx *gin.Context)

		AllToStatus(ctx *gin.Context)

		GetDissertationPage(ctx *gin.Context)
		UpsertFeedback(ctx *gin.Context)
		DownloadDissertation(ctx *gin.Context)

		GetTeachingLoad(ctx *gin.Context)
		GetScientificWorks(ctx *gin.Context)

		GetStudentStatus(ctx *gin.Context)
		GetSupervisorProfile(ctx *gin.Context)
		GetStudentProfile(ctx *gin.Context)

		GetStudentsReportComments(ctx *gin.Context)

		GetAllMarks(ctx *gin.Context)
	}

	AdministratorHandler interface {
		ChangeSupervisor(ctx *gin.Context)
		GetPairs(ctx *gin.Context)
		SetStudentStudyingStatus(ctx *gin.Context)
		GetSupervisors(ctx *gin.Context)

		GetSpecializations(ctx *gin.Context)
		GetGroups(ctx *gin.Context)

		AddSpecializations(ctx *gin.Context)
		AddGroups(ctx *gin.Context)

		GetStudents(ctx *gin.Context)
		GetSupervisorsStudents(ctx *gin.Context)

		DeleteGroups(ctx *gin.Context)
		DeleteSpecializations(ctx *gin.Context)

		GetSupervisorProfile(ctx *gin.Context)

		UpsertAttestationMarks(ctx *gin.Context)

		AddStudents(ctx *gin.Context)
		AddSupervisors(ctx *gin.Context)

		DeleteSemesterAmounts(ctx *gin.Context)
		GetSemesterAmounts(ctx *gin.Context)
		AddAmounts(ctx *gin.Context)

		ArchiveSupervisor(ctx *gin.Context)

		GetNotRegisteredUsers(ctx *gin.Context)
		DeleteNotRegisteredUsers(ctx *gin.Context)
	}

	AuthenticationHandler interface {
		Authorize(ctx *gin.Context)
		FirstStudentRegistry(ctx *gin.Context)
		FirstSupervisorRegistry(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		TokenCheck(ctx *gin.Context)
	}
)

type HTTPServer struct {
	student        StudentHandler
	supervisor     SupervisorHandler
	administrator  AdministratorHandler
	authentication AuthenticationHandler
}

func NewHTTPServer(studentHandler StudentHandler, supervisorHandler SupervisorHandler, adminHandler AdministratorHandler, authenticationHandler AuthenticationHandler) *HTTPServer {
	return &HTTPServer{
		student:        studentHandler,
		supervisor:     supervisorHandler,
		administrator:  adminHandler,
		authentication: authenticationHandler,
	}
}

func (h *HTTPServer) InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding", "StudentID"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Content-Disposition"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	r := router.Group("/")
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// StudentHandlers init
	r.POST("/students/review/:token", h.student.AllToReview)
	r.GET("/students/info/:token", h.student.GetStudentStatus)

	r.GET("/students/dissertation/:token", h.student.GetDissertationPage)
	r.POST("/students/dissertation/progress/:token", h.student.UpsertSemesterProgress)
	r.PUT("/students/dissertation/file/:token", h.student.DownloadDissertation)
	r.POST("/students/dissertation/file/:token", h.student.UploadDissertation)
	//r.POST("/students/dissertation_title/review/:token", h.student.DissertationTitleToReview)
	r.POST("/students/dissertation_title/:token", h.student.UpsertDissertationTitle)
	//r.POST("/students/dissertation/review/:token", h.student.DissertationToReview)

	r.GET("/students/load/:token", h.student.GetTeachingLoad)
	r.POST("/students/load/classroom/:token", h.student.UpsertClassroomLoads)
	r.PUT("/students/load/classroom/:token", h.student.DeleteClassroomLoads)
	r.POST("/students/load/individual/:token", h.student.UpsertIndividualLoads)
	r.PUT("/students/load/individual/:token", h.student.DeleteIndividualLoads)
	r.POST("/students/load/additional/:token", h.student.UpsertAdditionalLoads)
	r.PUT("/students/load/additional/:token", h.student.DeleteAdditionalLoads)
	//r.POST("/student/load/review/:token", h.student.TeachingLoadToReview)

	r.GET("/students/works/:token", h.student.GetScientificWorks)
	r.POST("/students/works/publications/:token", h.student.UpsertPublications)
	r.PUT("/students/works/publications/:token", h.student.DeletePublications)
	r.POST("/students/works/conferences/:token", h.student.UpsertConferences)
	r.PUT("/students/works/conferences/:token", h.student.DeleteConferences)
	r.POST("/students/works/projects/:token", h.student.UpsertResearchProjects)
	r.PUT("/students/works/projects/:token", h.student.DeleteProjects)
	//r.POST("/students/works/review/:token", h.student.ScientificWorksToReview)
	r.POST("/students/works/patents/:token", h.student.UpsertPatents)
	r.PUT("/students/works/patents/:token", h.student.DeletePatents)

	r.GET("/student/enum/specializations/:token", h.student.GetSpecializations)
	r.GET("/student/enum/groups/:token", h.student.GetGroups)
	r.GET("/student/supervisors/list/:token", h.student.GetSupervisors)

	r.POST("/students/dissertation/progress/percent/:token", h.student.UpdateProgressiveness)

	r.GET("/student/profile/:token", h.student.GetStudentProfile)

	r.GET("/students/report/comments/:token", h.student.GetReportComments)
	r.POST("/students/report/comments/:token", h.student.UpsertReportComments)

	r.GET("/students/marks/:token", h.student.GetAllMarks)
	r.POST("/students/exams/:token", h.student.UpsertExamResults)

	r.GET("/students/enum/amounts/:token", h.student.GetSemesterAmounts)

	r.POST("/students/profile/:token", h.student.UpdateProfile)
	r.PUT("/students/marks/:token", h.student.DeleteMarks)

	// r.PUT("/students/dissertation/file/:token", h.student.DownloadPresentation)

	// SupervisorHandler init
	r.GET("/supervisors/student/list/:token", h.supervisor.GetStudentsList)

	r.POST("/supervisors/student/review/:token", h.supervisor.AllToStatus)

	r.PUT("/supervisors/student/dissertation/:token", h.supervisor.GetDissertationPage)
	r.POST("/supervisors/student/feedback/:token", h.supervisor.UpsertFeedback)
	r.PUT("/supervisors/student/dissertation/file/:token", h.supervisor.DownloadDissertation)

	r.PUT("/supervisors/student/load/:token", h.supervisor.GetTeachingLoad)
	r.PUT("/supervisors/student/works/:token", h.supervisor.GetScientificWorks)

	r.PUT("/supervisors/student/info/:token", h.supervisor.GetStudentStatus)
	r.GET("/supervisors/profile/:token", h.supervisor.GetSupervisorProfile)

	r.PUT("/supervisor/student/profile/:token", h.supervisor.GetStudentProfile)

	r.PUT("/supervisors/report/comments/:token", h.supervisor.GetStudentsReportComments)

	r.PUT("/supervisors/student/marks/:token", h.supervisor.GetAllMarks)
	//r.POST("/supervisors/student/marks/:token", h.supervisor.UpsertSupervisorMark)

	// AdministratorHandler init
	r.POST("/administrator/student/change/:token", h.administrator.ChangeSupervisor)

	r.GET("/administrator/pairs/:token", h.administrator.GetPairs)
	r.POST("/administrator/student/status/:token", h.administrator.SetStudentStudyingStatus)

	r.GET("/administrator/supervisors/list/:token", h.administrator.GetSupervisors)

	r.GET("/administrator/enum/specializations/:token", h.administrator.GetSpecializations)
	r.GET("/administrator/enum/groups/:token", h.administrator.GetGroups)

	r.POST("/administrator/enum/specializations/:token", h.administrator.AddSpecializations)
	r.POST("/administrator/enum/groups/:token", h.administrator.AddGroups)

	r.GET("/administrator/students/list/:token", h.administrator.GetStudents)

	r.PUT("/administrator/supervisor/students/:token", h.administrator.GetSupervisorsStudents)

	r.PUT("/administrator/enum/specializations/:token", h.administrator.DeleteSpecializations)
	r.PUT("/administrator/enum/groups/:token", h.administrator.DeleteGroups)

	r.PUT("/administrator/supervisors/profile/:token", h.administrator.GetSupervisorProfile)

	r.POST("/administrator/student/attestation/marks/:token", h.administrator.UpsertAttestationMarks)

	r.POST("/administrator/users/students/:token", h.administrator.AddStudents)
	r.POST("/administrator/users/supervisors/:token", h.administrator.AddSupervisors)
	r.GET("/administrator/enum/amounts/:token", h.administrator.GetSemesterAmounts)
	r.PUT("/administrator/enum/amounts/:token", h.administrator.DeleteSemesterAmounts)

	r.PUT("/administrator/supervisor/:token", h.administrator.ArchiveSupervisor)
	r.POST("/administrator/enum/amounts/:token", h.administrator.AddAmounts)

	r.GET("/administrator/users/not_registered/:token", h.administrator.GetNotRegisteredUsers)
	r.PUT("/administrator/users/not_registered/:token", h.administrator.DeleteNotRegisteredUsers)

	// AuthenticationHandler init
	r.POST("/authorize", h.authentication.Authorize)

	r.POST("/authorize/registration/student/:token", h.authentication.FirstStudentRegistry)
	r.POST("/authorize/registration/supervisor/:token", h.authentication.FirstSupervisorRegistry)

	r.POST("/authorize/password/change/:token", h.authentication.ChangePassword)
	r.GET("/authorize/token/check/:token", h.authentication.TokenCheck)

	return router
}
