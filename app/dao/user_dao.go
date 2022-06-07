package dao

import (
	"gin-mvc/app/entity"
	"gin-mvc/common/logger"
)

func UserAdd(user *entity.User) error {
	err := mysqlDB.AutoMigrate(&entity.User{})
	if err != nil {
		logger.Errorf("table users automigrate failed:%w", err)
		return err
	}

	if err = mysqlDB.Create(user).Error; err != nil {
		logger.Errorf("table users create new user:%v failed:%w", user, err)
	}
	return err
}

func UserGetByName(user *entity.User) error {
	err := mysqlDB.AutoMigrate(&entity.User{})
	if err != nil {
		logger.Errorf("table users automigrate failed:%w", err)
	}

	if err = mysqlDB.Where("username=?", user.UserName).First(user).Error; err != nil {
		logger.Errorf("table users query user:%v by username failed:%w", user, err)
	}
	return err
}

func UserGetByUID(user *entity.User) error {
	err := mysqlDB.AutoMigrate(&entity.User{})
	if err != nil {
		logger.Errorf("table users automigrate failed:%w", err)
	}

	if err = mysqlDB.Where("user_id=?", user.UserID).First(user).Error; err != nil {
		logger.Errorf("table users query user:%v by ID failed:%w", user, err)
	}
	return err
}

func UserUpdate(user *entity.User) error {
	err := mysqlDB.AutoMigrate(&entity.User{})
	if err != nil {
		logger.Errorf("table users automigrate failed:%w", err)
	}
	if err = mysqlDB.Save(user).Error; err != nil {
		logger.Errorf("table users update user:%v failed:%w", user, err)
	}
	return err
}
