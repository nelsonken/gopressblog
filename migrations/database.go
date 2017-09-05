package main

import (
	"blog/config"
	"blog/models"
	"blog/services"
)

func main() {

	opts := &config.Options{}
	opts.Database = &services.DBOptions{}
	config.GetConfig("config/config.yaml", opts)

	dbService := services.NewDBService(opts.Database.DBType, opts.Database)
	defer dbService.ORM.Close()

	// 自动迁移模式
	dbService.ORM.AutoMigrate(
		&models.Account{},
		&models.Post{},
		&models.Comment{},
		&models.Message{},
		&models.User{},
	)

	u := new(models.User)
	dbService.ORM.Where("name = ?", "system").First(u)
	if u.ID <= 0 {
		dbService.ORM.Create(&models.User{Name: "system", Password: ""})
	}

}
