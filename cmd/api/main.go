package main

import (
	"log"

	"backend-mosque-tahfidz-management/container"
	"backend-mosque-tahfidz-management/internal/config"
	"backend-mosque-tahfidz-management/internal/middleware"
	"backend-mosque-tahfidz-management/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	c := container.New(cfg, db)

	app := middleware.NewApp()
	app.Use(middleware.CORS())

	routes.Setup(app, c.AuthHandler, c.StudentHandler, c.ProgressHandler, c.ActivityLogHandler, c.SurahHandler, c.UploadHandler, c.TokenMaker)

	log.Fatal(app.Listen(":3010"))
}
