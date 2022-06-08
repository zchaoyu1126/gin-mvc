package dao

import (
	"gin-mvc/app/entity"
	"gin-mvc/common/db"

	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

func init() {
	mysqlDB = db.NewMySQLConnInstance().DB
	mysqlDB.AutoMigrate(&entity.User{})
}
