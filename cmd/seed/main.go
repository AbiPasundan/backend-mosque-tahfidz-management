package main

import (
	"flag"
	"log"
	"os"

	"backend-mosque-tahfidz-management/internal/config"
	"backend-mosque-tahfidz-management/internal/database/seeds"
)

func main() {
	fresh := flag.Bool("fresh", false, "Clean all tables before seeding")
	clean := flag.Bool("clean", false, "Only clean all tables (no seeding)")
	admin := flag.Bool("admin", false, "Seed only a single admin account")
	flag.Parse()

	cfg := config.LoadConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("📦 Database connected successfully")

	switch {
	case *clean:
		log.Println("🧹 Running clean only...")
		if err := seeds.CleanAll(db); err != nil {
			log.Fatalf("❌ Clean failed: %v", err)
		}
	case *fresh:
		log.Println("🔄 Running fresh seed (clean + seed)...")
		if err := seeds.FreshSeed(db); err != nil {
			log.Fatalf("❌ Fresh seed failed: %v", err)
		}
	case *admin:
		log.Println("👤 Seeding admin account...")
		seeder := &seeds.AdminSeeder{}
		if err := seeder.Seed(db); err != nil {
			log.Fatalf("❌ Admin seed failed: %v", err)
		}
	default:
		log.Println("🌱 Running seeders...")
		if err := seeds.RunAll(db); err != nil {
			log.Fatalf("❌ Seed failed: %v", err)
		}
	}

	os.Exit(0)
}
