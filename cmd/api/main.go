package main

import (
	"log"

	"backend-mosque-tahfidz-management/internal/config"
	"backend-mosque-tahfidz-management/internal/domain/activity_log"
	"backend-mosque-tahfidz-management/internal/domain/auth"
	"backend-mosque-tahfidz-management/internal/domain/progress"
	"backend-mosque-tahfidz-management/internal/domain/student"
	"backend-mosque-tahfidz-management/internal/domain/surah"
	"backend-mosque-tahfidz-management/internal/routes"
	"backend-mosque-tahfidz-management/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

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

	// Handlers
	authHandler := auth.NewAuthHandler(authService, activityLogService)
	studentHandler := student.NewStudentHandler(studentService, authService, activityLogService)
	progressHandler := progress.NewProgressHandler(progressService, authService, studentService, activityLogService)
	activityLogHandler := activity_log.NewActivityLogHandler(activityLogService)
	surahHandler := surah.NewSurahHandler()

	// Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Setup routes
	routes.Setup(app, authHandler, studentHandler, progressHandler, activityLogHandler, surahHandler, tokenMaker)

	log.Fatal(app.Listen(":3010"))
}
