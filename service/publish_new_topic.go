package service

import "goweb-sample/repository"

func PublishNewTopic(title, content string) error {
	return repository.NewTopicDaoInstance().AddTopic(title, content)
}
