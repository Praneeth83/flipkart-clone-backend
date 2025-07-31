package config

import (
	"fmt"
	"log"
	"test1/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "postgresql://flipkart_clone_user:DdAKsuDfNeN52QCeudB76EJeL40kbsM0@dpg-d25f723e5dus73a1bh30-a.oregon-postgres.render.com/flipkart_clone"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL on Render")

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Auto migration failed:", err)
	}
}
