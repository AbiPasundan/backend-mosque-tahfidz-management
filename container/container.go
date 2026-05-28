package container

import (
	"backend-mosque-tahfidz-management/internal/config"
	"backend-mosque-tahfidz-management/internal/domain/activity_log"
	"backend-mosque-tahfidz-management/internal/domain/auth"
	"backend-mosque-tahfidz-management/internal/domain/memorize"
	"backend-mosque-tahfidz-management/internal/domain/progress"
	"backend-mosque-tahfidz-management/internal/domain/student"
	"backend-mosque-tahfidz-management/internal/domain/surah"
	"backend-mosque-tahfidz-management/internal/domain/upload"
	"backend-mosque-tahfidz-management/pkg/token"

	"github.com/jmoiron/sqlx"
)

// Container holds all application dependencies (handlers, services, token maker).
type Container struct {
	AuthHandler        *auth.AuthHandler
	StudentHandler     *student.StudentHandler
	ProgressHandler    *progress.ProgressHandler
	MemorizeHandler    *memorize.MemorizeHandler
	ActivityLogHandler *activity_log.ActivityLogHandler
	SurahHandler       *surah.SurahHandler
	UploadHandler      *upload.UploadHandler
	TokenMaker         token.Maker
}

// New creates a Container with all repositories, services, and handlers wired up.
func New(cfg *config.Config, db *sqlx.DB) *Container {
	tokenMaker := token.NewJWTMaker(cfg.JWTSecret)

	// Repositories
	userRepo := auth.NewUserRepository(db)
	studentRepo := student.NewStudentRepository(db)
	progressRepo := progress.NewProgressRepository(db)
	activityLogRepo := activity_log.NewActivityLogRepository(db)

	// Services
	authService := auth.NewAuthService(userRepo, tokenMaker)
	studentService := student.NewStudentService(studentRepo)
	progressService := progress.NewProgressService(progressRepo)
	activityLogService := activity_log.NewActivityLogService(activityLogRepo)
	cloudinaryService := upload.NewCloudinaryService(cfg)

	// Handlers
	authHandler := auth.NewAuthHandler(authService, activityLogService)
	studentHandler := student.NewStudentHandler(studentService, authService, activityLogService)
	progressHandler := progress.NewProgressHandler(progressService, authService, studentService, activityLogService)
	memorizeRepo := memorize.NewMemorizeRepository(db)
	memorizeService := memorize.NewMemorizeService(memorizeRepo)
	memorizeHandler := memorize.NewMemorizeHandler(memorizeService, authService, activityLogService)
	activityLogHandler := activity_log.NewActivityLogHandler(activityLogService)
	surahHandler := surah.NewSurahHandler()
	uploadHandler := upload.NewUploadHandler(cloudinaryService)

	return &Container{
		AuthHandler:        authHandler,
		StudentHandler:     studentHandler,
		ProgressHandler:    progressHandler,
		MemorizeHandler:    memorizeHandler,
		ActivityLogHandler: activityLogHandler,
		SurahHandler:       surahHandler,
		UploadHandler:      uploadHandler,
		TokenMaker:         tokenMaker,
	}
}
