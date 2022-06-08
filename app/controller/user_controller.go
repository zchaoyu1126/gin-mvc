package controller

import (
	"gin-mvc/app/service"
	"gin-mvc/common/xerr"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserVO struct {
	ID            int64  `json:"id"`
	UserName      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserLoginResponse struct {
	Response
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	UserVO `json:"user"`
}

// 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.Register(username, password)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: success,
		UserID:   id,
		Token:    token,
	})
}

// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.Login(username, password)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: success,
		UserID:   id,
		Token:    token,
	})
}

// 用户信息
func UserInfo(c *gin.Context) {
	// 获取Query参数 user_id, user_id解析失败时返回ErrBadRequest错误
	toUserIDTmp := c.Query("user_id")
	toUserID, err := strconv.ParseInt(toUserIDTmp, 10, 64)
	if err != nil {
		zap.S().Error("parse user_id:%v failed", toUserIDTmp)
		errorHandler(c, xerr.ErrBadRequest)
	}

	// 获取经过路由中间件UserAuthMiddlerWare解析后又写入的fromUserID信息
	fromUserIDTmp, exists := c.Get("fromUserID")
	fromUserID, ok := fromUserIDTmp.(int64)
	if !exists || !ok {
		zap.S().Errorf("parse fromUserID:%v failed", fromUserIDTmp)
		errorHandler(c, xerr.ErrBadRequest)
		return
	}

	userDTO, err := service.UserInfo(toUserID, fromUserID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	userVO := UserVO{
		ID:            userDTO.ID,
		UserName:      userDTO.UserName,
		FollowCount:   userDTO.FollowCount,
		FollowerCount: userDTO.FollowerCount,
		IsFollow:      userDTO.IsFollow,
	}

	c.JSON(http.StatusOK, UserInfoResponse{
		Response: success,
		UserVO:   userVO,
	})
}
