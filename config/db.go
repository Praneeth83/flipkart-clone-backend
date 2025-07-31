package config

import (
	"fmt"
	"log"
	"test1/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "postgresql://flipkart_clone_user:DdAKsuDfNeN52QCeudB76EJeL40kbsM0@dpg-d25f723e5dus73a1bh30-a.oregon-postgres.render.com/flipkart_clone"

	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL on Render")

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Auto migration failed:", err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Failed to get sql.DB from GORM:", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
