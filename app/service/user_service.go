package service

import (
	"gin-mvc/app/dao"
	"gin-mvc/app/entity"
	"gin-mvc/common/auth"
	"gin-mvc/common/db"
	"gin-mvc/common/logger"
	"gin-mvc/common/utils"
	"gin-mvc/common/xerr"
	"regexp"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID            int64
	UserName      string
	FollowCount   int64
	FollowerCount int64
	IsFollow      bool
}

// 去数据库中查询用户名是否存在
func CheckUserNameExist(username string) (bool, error) {
	user := &entity.User{UserName: username}
	err := dao.UserGetByName(user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, xerr.ErrDatabase
	}
	return true, nil
}

// Before register, controller has checked wheather username exist.
// If registation success, return userID, token, nil, otherwise return 0, "", err.
func Register(username, password string) (int64, string, error) {
	// check username and password
	// 用户名由字母数据下划线英文句号组成，长度要求4-16之间
	usernameReg, err := regexp.Compile(`^[a-zA-Z0-9_\.]{4,16}$`)
	if err != nil {
		logger.Errorf("init regexp:%v failed", usernameReg)
		return 0, "", xerr.ErrInternalServer
	}
	if !usernameReg.MatchString(username) {
		return 0, "", xerr.ErrUsernameValidation
	}

	// 密码匹配6-16位英文数据大部分英文标点
	passwordReg, err := regexp.Compile(`^([A-Za-z0-9\-=\[\];,\./~!@#\$%^\*\(\)_\+}{:\?]){6,16}$`)
	if err != nil {
		return 0, "", xerr.ErrInternalServer
	}
	if !passwordReg.MatchString(password) {
		return 0, "", xerr.ErrPasswordValidation
	}

	// use redis check wheather this user has already existed
	exist, err := db.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		logger.Errorf("redis query username:%v exist failed:%w", username, err)
		return 0, "", xerr.ErrDatabase
	}
	if exist {
		return 0, "", xerr.ErrUserExist
	}

	// encrypt password and then store into mysql
	encodePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("encrypt password failed:%w", err)
		return 0, "", xerr.ErrInternalServer
	}
	user := &entity.User{UserName: username, Password: string(encodePassword)}
	snowFlake, err := utils.NewSnowFlake(0, 0)
	if err != nil {
		logger.Errorf("init snow flake failed:%w", err)
		return 0, "", xerr.ErrInternalServer
	}

	user.UserID, err = snowFlake.NextId()
	if err != nil {
		logger.Errorf("generate uid failed:%w", err)
		return 0, "", xerr.ErrInternalServer
	}

	if err := dao.UserAdd(user); err != nil {
		return 0, "", xerr.ErrDatabase
	}

	// store namelist username into redis
	if err := db.NewRedisDaoInstance().AddToNameList(user.UserName); err != nil {
		logger.Errorf("redis username failed:%w", err)
		return 0, "", xerr.ErrDatabase
	}

	// store (token, userID) into redis
	token := auth.GenerateToken(user.UserID)
	if err := db.NewRedisDaoInstance().SetToken(token, user.UserID); err != nil {
		logger.Errorf("cache token failed:%w", err)
		return 0, "", xerr.ErrDatabase
	}
	return user.UserID, token, nil
}

// Before login, controller has checked wheather username exist.
// If login success, return userID, token, nil, otherwise return 0, "", err.
func Login(username, password string) (int64, string, error) {
	// check username and password and sql injection attack

	// use redis check wheather this user has already existed
	exist, err := db.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		return 0, "", xerr.ErrDatabase
	}
	if !exist {
		return 0, "", xerr.ErrUserNotFound
	}

	user := &entity.User{UserName: username}
	if err := dao.UserGetByName(user); err != nil {
		return 0, "", xerr.ErrDatabase
	}

	// decrypt the string stored in mysql
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return -1, "", xerr.ErrPasswordIncorrect
	}

	token := auth.GenerateToken(user.UserID)
	db.NewRedisDaoInstance().SetToken(token, user.UserID)
	return user.UserID, token, nil
}

// fromUserID wants to look over toUserID's user information.
func UserInfo(toUserID, fromUserID int64) (*UserDTO, error) {
	// check toUserID and fromUserID

	user := &entity.User{UserID: toUserID}
	if err := dao.UserGetByUID(user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.ErrUserNotFound
		} else {
			return nil, xerr.ErrDatabase
		}
	}

	// query like table, check fromUserID is follow toUserID or not
	// follow := &entity.Follow{FromUserID: fromUserID, ToUserID: toUserID}
	// if fromUserID == toUserID {
	// 	follow.IsFollow = false
	// } else {
	// 	if err := dao.FollowGetByIDs(follow); err != nil {
	// 		return nil, errors.Warp(errors.ErrDatabase, err.Error())
	// 	}
	// }
	// userDTO := &UserDTO{user.UserID, user.UserName, user.FollowCount, user.FollowerCount, follow.IsFollow}
	userDTO := &UserDTO{user.UserID, user.UserName, user.FollowCount, user.FollowerCount, false}
	return userDTO, nil
}
