package main

import (
	"log"

	"quicky-go/configs"
	_ "quicky-go/repo_migrate"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// print current directory
	log.Println("Current directory:", configs.CurrentAppPath)
	log.Println("Config file path:", configs.ConfigFilePath)
	log.Println("Starting application...")
	log.Println(configs.Info.DB.DBName)

}
