package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go-google-cloud-storage/app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDatabase(viper *viper.Viper) *gorm.DB {
	username := viper.GetString("DB_USERNAME")
	password := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	database := viper.GetString("DB_NAME")
	sslMode := viper.GetString("DB_SSLMODE")

	connMaxIdleTime := viper.GetDuration("DB_CONN_MAX_IDLE_TIME")
	if connMaxIdleTime == 0 {
		log.Fatalf("» Failed to get setting: DB_CONN_MAX_IDLE_TIME")
	}
	connMaxLifetime := viper.GetDuration("DB_CONN_MAX_LIFE_TIME")
	if connMaxLifetime == 0 {
		log.Fatalf("» Failed to get setting: DB_CONN_MAX_LIFE_TIME")
	}
	maxOpenConn := viper.GetInt("DB_MAX_OPEN_CONN")
	if maxOpenConn == 0 {
		log.Fatalf("» Failed to get setting: DB_MAX_OPEN_CONN")
	}
	maxIdleConn := viper.GetInt("DB_MAX_IDLE_CONN")
	if maxIdleConn == 0 {
		log.Fatalf("» Failed to get setting: DB_MAX_IDLE_CONN")
	}

	log.Println("» Trying to connect database")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, username, database, password, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("» Failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("» Failed to get database instance: %v", err)
	}

	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetMaxIdleConns(maxIdleConn)

	log.Println("» Success connect to database")

	return db
}

func DisconnectDB(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to kill connection from database")
	}
	err = dbSQL.Close()
	if err != nil {
		panic("Failed to kill connection from database")
	}
}

func SeedDatabase(db *gorm.DB, config *viper.Viper) {
	
}

func MigrateDatabase(db *gorm.DB) {
	var err error
	InitUUID(db)

	// Config
	err = db.AutoMigrate(&model.Config{})
	if err != nil {
		log.Println("[ERROR] Error migrate Config model: " + err.Error())
	}

}

func InitUUID(db *gorm.DB) {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
}
