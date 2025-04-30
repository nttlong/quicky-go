package main

import (
	"log"

	"quicky-go/configs"

	"quicky-go/repo"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// print current directory
	log.Println("Current directory:", configs.CurrentAppPath)
	log.Println("Config file path:", configs.ConfigFilePath)
	log.Println("Starting application...")
	log.Println(configs.Info.DB.DBName)
	cnn, err := repo.GetRepo("test-001")
	if err != nil {
		log.Println(err)
	}
	log.Println(cnn)

}
