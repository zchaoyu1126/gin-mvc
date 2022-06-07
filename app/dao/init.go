package dao

import (
	"gin-mvc/common/db"

	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

func init() {
	mysqlDB = db.NewMySQLConnInstance().DB
}
