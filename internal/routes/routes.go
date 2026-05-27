package routes

import (
	"backend-mosque-tahfidz-management/internal/domain/activity_log"
	"backend-mosque-tahfidz-management/internal/domain/auth"
	"backend-mosque-tahfidz-management/internal/domain/progress"
	"backend-mosque-tahfidz-management/internal/domain/student"
	"backend-mosque-tahfidz-management/internal/domain/surah"
	"backend-mosque-tahfidz-management/internal/domain/upload"
	"backend-mosque-tahfidz-management/internal/middleware"
	"backend-mosque-tahfidz-management/pkg/token"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Setup(
	app *fiber.App,
	authHandler *auth.AuthHandler,
	studentHandler *student.StudentHandler,
	progressHandler *progress.ProgressHandler,
	activityLogHandler *activity_log.ActivityLogHandler,
	surahHandler *surah.SurahHandler,
	uploadHandler *upload.UploadHandler,
	tokenMaker token.Maker,
) {
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	api := app.Group("/api/v1")

	// Auth with rate limiting
	authGroup := api.Group("/auth")
	authGroup.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
	}))
	authGroup.Get("me", middleware.JWT(tokenMaker), authHandler.Me)
	authGroup.Patch("profile", middleware.JWT(tokenMaker), authHandler.UpdateProfile)
	authGroup.Patch("password", middleware.JWT(tokenMaker), authHandler.UpdatePassword)
	authGroup.Post("login", authHandler.Login)
	authGroup.Post("logout", authHandler.Logout)

	// Users (admin only)
	userGroup := api.Group("/users", middleware.JWT(tokenMaker), middleware.RBAC("admin"))
	userGroup.Get("/", authHandler.ListUsers)
	userGroup.Post("/", authHandler.CreateUser)
	userGroup.Get("/:id", authHandler.GetUser)
	userGroup.Put("/:id", authHandler.UpdateUser)
	userGroup.Delete("/:id", authHandler.DeleteUser)

	// Mentors
	api.Get("/mentors/:id", middleware.JWT(tokenMaker), authHandler.GetMentorDetail)

	// Students

	studentGroup := api.Group("/students")
	studentGroup.Get("/", studentHandler.ListStudents)
	studentGroup.Post("/", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), studentHandler.CreateStudent)
	studentGroup.Get("/:id", studentHandler.GetStudent)
	studentGroup.Put("/:id", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), studentHandler.UpdateStudent)
	studentGroup.Delete("/:id", middleware.JWT(tokenMaker), middleware.RBAC("admin"), studentHandler.DeleteStudent)

	// Surahs
	api.Get("/surahs", surahHandler.GetSurahs)

	// Progress
	api.Post("/progress/bulk", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), progressHandler.BulkCreateProgress)
	api.Get("/progress", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), progressHandler.ListProgress)
	api.Post("/progress", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), progressHandler.CreateProgress)
	api.Put("/progress/:id", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), progressHandler.UpdateProgress)

	// Dashboard
	api.Get("/dashboard/summary", progressHandler.GetDashboardSummary)

	// Activity Logs
	api.Get("/activity-logs", middleware.JWT(tokenMaker), middleware.RBAC("admin"), activityLogHandler.ListActivityLogs)

	// File Upload (Cloudinary)
	api.Post("/upload", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), uploadHandler.UploadFile)
}
