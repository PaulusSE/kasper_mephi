package main //nolint

import (
	"context"
	"os"
	"uir_draft/internal/pkg/service/registration"

	"uir_draft/internal/app/kasper"
	"uir_draft/internal/handlers/administator_handler"
	"uir_draft/internal/handlers/authorization_handler"
	"uir_draft/internal/handlers/student_handler"
	"uir_draft/internal/handlers/supervisor_handler"
	"uir_draft/internal/pkg/configs"
	"uir_draft/internal/pkg/repository"
	"uir_draft/internal/pkg/service/admin"
	"uir_draft/internal/pkg/service/authentication"
	"uir_draft/internal/pkg/service/email"
	"uir_draft/internal/pkg/service/enum"
	"uir_draft/internal/pkg/service/student"
	"uir_draft/internal/pkg/service/supervisor"

	"github.com/spf13/viper"
)

//	@title		Сервис Kasper
//	@version	0.1.2
// description Серверная часть Системы учета деятельности аспирантов

//	@host	localhost:8080

func main() {
	err := initConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := configs.InitPostgresDB(ctx, configs.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		panic(err)
	}

	clientRepo := repository.NewClientRepository()
	usersRepo := repository.NewUsersRepository()

	studentService := student.NewService(db)
	adminService := admin.NewService(db)
	supervisorService := supervisor.NewService(db)
	authenticationService := authentication.NewService(db)
	registrationService := registration.NewService(db)

	emailService := email.NewService("info@kasper-mephi.ru", os.Getenv("MAIL_PASSWORD"), "mail.hosting.reg.ru", db, usersRepo, clientRepo)
	enumService := enum.NewService(db)

	cache := student_handler.NewPgRecommendationCache(db)
	studentHandler := student_handler.NewHandler(studentService, authenticationService, emailService, enumService, adminService, cache)
	supervisorHandler := supervisor_handler.NewHandler(studentService, authenticationService, supervisorService, emailService, registrationService)
	adminHandler := administator_handler.NewHandler(adminService, authenticationService, enumService, supervisorService, emailService, registrationService, studentService)
	authenticationHandler := authorization_handler.NewHandler(authenticationService, studentService, supervisorService, registrationService)

	server := kasper.NewHTTPServer(studentHandler, supervisorHandler, adminHandler, authenticationHandler)
	r := server.InitRouter()

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
