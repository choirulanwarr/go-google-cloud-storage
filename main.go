package main

import (
	"fmt"
	"go-google-cloud-storage/app/config"
	"go-google-cloud-storage/app/helper"
	"log"
)

func main() {
	viper := config.NewViper()
	db := config.InitDatabase(viper)
	server := config.NewServer()
	validator := config.NewValidator()
	logger := helper.Logger
	defer func() {
		helper.CloseLoggerFile()
	}()

	autoMigrate := viper.GetBool("AUTO_MIGRATION_SWITCH")
	if autoMigrate {
		log.Println("» Auto migrate = true")
		log.Println("» Database migration is running")
		config.MigrateDatabase(db)
		config.SeedDatabase(db, viper)
	} else {
		log.Println("» Auto migrate = false")
		log.Println("» No database migration")
	}
	defer config.DisconnectDB(db)

	config.InitConfig(&config.AppConfig{
		Config:    viper,
		DB:        db,
		Server:    server,
		Validator: validator,
		Logger:    logger,
	})

	appPort := viper.GetInt("APP_PORT")
	err := server.Run(fmt.Sprintf(":%d", appPort))
	if err != nil {
		log.Fatalf("» Server error : %v", err)
	}
}
