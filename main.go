package main

import (
	"log"

	"quicky-go/configs"
	"quicky-go/user_repository"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// print current directory
	log.Println("Current directory:", configs.CurrentAppPath)
	log.Println("Config file path:", configs.ConfigFilePath)
	log.Println("Starting application...")
	log.Println(configs.Info.DB.DBName)

	user := user_repository.User.FindById(1)

	log.Println(user)
}
