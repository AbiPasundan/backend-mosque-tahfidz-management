package routes

import (
	"encoding/json"
	"fmt"
	"time"

	"backend-mosque-tahfidz-management/internal/domain/auth"
	"backend-mosque-tahfidz-management/internal/domain/progress"
	"backend-mosque-tahfidz-management/internal/domain/student"
	"backend-mosque-tahfidz-management/internal/middleware"
	"backend-mosque-tahfidz-management/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Setup(app *fiber.App, authHandler *auth.AuthHandler, studentHandler *student.StudentHandler, progressHandler *progress.ProgressHandler, tokenMaker token.Maker) {
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	api := app.Group("/api/v1")

	api.Get("/me", authHandler.Me)

	// Auth with rate limiting
	authGroup := api.Group("/auth")
	authGroup.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
	}))
	authGroup.Get("/me", middleware.JWT(tokenMaker), authHandler.Me)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", authHandler.Logout)

	// Users (admin only)
	userGroup := api.Group("/users", middleware.JWT(tokenMaker), middleware.RBAC("admin"))
	userGroup.Get("/", authHandler.ListUsers)
	userGroup.Post("/", authHandler.CreateUser)
	userGroup.Get("/:id", authHandler.GetUser)
	userGroup.Put("/:id", authHandler.UpdateUser)
	userGroup.Delete("/:id", authHandler.DeleteUser)

	// Students
	studentGroup := api.Group("/students")
	studentGroup.Get("/", studentHandler.ListStudents)
	studentGroup.Post("/", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), studentHandler.CreateStudent)
	studentGroup.Get("/:id", studentHandler.GetStudent)
	studentGroup.Put("/:id", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"), studentHandler.UpdateStudent)
	studentGroup.Delete("/:id", middleware.JWT(tokenMaker), middleware.RBAC("admin"), studentHandler.DeleteStudent)

	// Surahs
	api.Get("/surahs", func(c *fiber.Ctx) error {
		data := getSurahs()

		if data == nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Gagal mengambil data dari API pusat",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    data,
		})
	})

	// Progress
	progressGroup := api.Group("/progress", middleware.JWT(tokenMaker), middleware.RBAC("mentor", "admin"))
	progressGroup.Get("/", progressHandler.ListProgress)
	progressGroup.Post("/", progressHandler.CreateProgress)
	progressGroup.Put("/:id", progressHandler.UpdateProgress)

	// Dashboard
	api.Get("/dashboard/summary", progressHandler.GetDashboardSummary)
}

func getSurahs() any {
	agent := fiber.Get("https://equran.id/api/v2/surat")

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Nomor       int               `json:"nomor"`
			Nama        string            `json:"nama"`
			NamaLatin   string            `json:"namaLatin"`
			JumlahAyat  int               `json:"jumlahAyat"`
			TempatTurun string            `json:"tempatTurun"`
			Arti        string            `json:"arti"`
			Deskripsi   string            `json:"deskripsi"`
			AudioFull   map[string]string `json:"audioFull"`
		} `json:"data"`
	}

	agent.Timeout(10 * time.Second)

	statusCode, body, errs := agent.Bytes()

	if len(errs) > 0 {
		fmt.Println("Error Fetching:", errs)
		return nil
	}

	if statusCode != 200 {
		fmt.Println("Status Code Not 200:", statusCode)
		return nil
	}

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Unmarshal Error:", err)
		return nil
	}

	return response.Data
}
