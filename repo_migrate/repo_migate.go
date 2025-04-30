package repomigrate

import (
	"quicky-go/models/user"
	"quicky-go/repo"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {
	err := repo.Repo.AutoMigrate(&user.User{})
	panic(err)
}
