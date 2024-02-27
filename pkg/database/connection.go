package database

import (
	"fmt"
	"log"
	"time"

	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
)

// Connection instance
var (
	DB    *gorm.DB
	SqlDB *sql.DB
)

func CheckConnection() {
	log.Println("Checking database connection")
	if err := SqlDB.Ping(); err != nil {
		log.Println("Failed to Ping database")
	}
}

func CloseConnection() {
	log.Println("Closing database connection")
	SqlDB.Close()
}

// Connect to database
func Init() {
	// Open connection to database
	db, err := gorm.Open(getDialector(), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	log.Println("Connected to database")

	// Set database connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to set database connection pool")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Println("Set database connection pool")

	// Migrate the schema
	db.AutoMigrate(&models.Advertisement{})

	// Set database connection instance
	DB = db
	SqlDB = sqlDB

	// Enable Logger, show detailed log
	if config.Settings.App.Env == "dev" {
		DB = db.Debug()
	}
}

// get database dialector
func getDialector() gorm.Dialector {
	// Get database configuration
	switch config.Settings.Database.Kind {
	case "mysql":
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Settings.Database.User,
			config.Settings.Database.Password,
			config.Settings.Database.Host,
			config.Settings.Database.Port,
			config.Settings.Database.Name,
		))
	case "postgres":
		return postgres.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Taipei",
			config.Settings.Database.Host,
			config.Settings.Database.Port,
			config.Settings.Database.User,
			config.Settings.Database.Name,
			config.Settings.Database.Password,
		))
	case "sqlite":
		return sqlite.Open(config.Settings.Database.Name + ".db")
	default:
		panic("Database kind not supported")
	}
}
