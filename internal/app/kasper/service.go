package kasper

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type StudentHandler interface {
	GetDissertationPage(ctx *gin.Context)
	UpsertSemesterProgress(ctx *gin.Context)
	GetScientificWorks(ctx *gin.Context)
	InsertScientificWorks(ctx *gin.Context)
	UpdateScientificWorks(ctx *gin.Context)
	DeleteScientificWork(ctx *gin.Context)
	GetTeachingLoad(ctx *gin.Context)
	UpsertTeachingLoad(ctx *gin.Context)
	DeleteTeachingLoad(ctx *gin.Context)
	UploadDissertation(ctx *gin.Context)
	DownloadDissertation(ctx *gin.Context)
}

type SupervisorHandler interface {
	GetListOfStudents(ctx *gin.Context)
	GetStudentsDissertationPage(ctx *gin.Context)
}

type AuthorizationHandler interface {
	Authorize(ctx *gin.Context)
}

func InitRoutes(student StudentHandler, supervisor SupervisorHandler, authorization AuthorizationHandler) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding", "StudentID"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "StudentID"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/students/dissertation/:id", student.GetDissertationPage)
	router.POST("/students/dissertation/progress/:id", student.UpsertSemesterProgress)
	router.GET("/students/scientific_works/:id", student.GetScientificWorks)
	router.POST("/students/scientific_works/:id", student.InsertScientificWorks)
	router.PATCH("/students/scientific_works/:id", student.UpdateScientificWorks)
	router.DELETE("/students/scientific_works/:id", student.DeleteScientificWork)
	router.GET("/students/teaching_load/:id", student.GetTeachingLoad)
	router.POST("/students/teaching_load/:id", student.UpsertTeachingLoad)
	router.DELETE("/students/teaching_load/:id", student.DeleteTeachingLoad)
	router.POST("/students/dissertation/file/:id", student.UploadDissertation)
	router.PUT("/students/dissertation/file/:id", student.DownloadDissertation)

	router.GET("/supervisors/list_of_students/:id", supervisor.GetListOfStudents)
	router.PUT("/supervisors/student/:id", supervisor.GetStudentsDissertationPage)

	router.POST("authorization/authorize", authorization.Authorize)
	return router
}
