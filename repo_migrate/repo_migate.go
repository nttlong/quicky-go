package repomigrate

import (
	"quicky-go/models/account"
	"quicky-go/models/personal"

	"quicky-go/repo"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {
	models := []interface{}{
		&account.Account{},
		&personal.PersonalInfo{},
		// Thêm các model khác của bạn vào đây
		// &your_other_package.YourOtherModel{},
	}

	for _, model := range models {
		err := repo.Repo.AutoMigrate(model)
		if err != nil {
			panic(err)
		}
	}
}
