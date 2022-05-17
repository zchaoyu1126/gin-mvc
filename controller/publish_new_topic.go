package controller

import "goweb-sample/service"

type PublishTopicResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func PublishNewTopic(title, content string) *PublishTopicResponse {
	err := service.PublishNewTopic(title, content)
	if err != nil {
		return &PublishTopicResponse{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	return &PublishTopicResponse{
		Code: 0,
		Msg:  "success",
	}
}
