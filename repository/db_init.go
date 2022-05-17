package repository

import (
	"bufio"
	"encoding/json"
	"os"
)

var topicID int64
var postID int64
var topicIndexMap map[int64]*Topic
var postIndexMap map[int64][]*Post
var filePath string

func Init(path string) error {
	filePath = path
	if err := initTopicIndexMap(path); err != nil {
		return err
	}
	if err := initPostsIndexMap(path); err != nil {
		return err
	}
	return nil
}

func initTopicIndexMap(filePath string) error {
	open, err := os.Open(filePath + "topic")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	topicTmpMap := make(map[int64]*Topic)
	for scanner.Scan() {
		// 使用迭代器的方式遍历数据行
		text := scanner.Text()
		var topic Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil {
			return err
		}
		topicTmpMap[topic.ID] = &topic
		if topic.ID > topicID {
			topicID = topic.ID
		}
	}
	topicIndexMap = topicTmpMap
	return nil
}

func initPostsIndexMap(filePath string) error {
	open, err := os.Open(filePath + "post")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	postsTmpMap := make(map[int64][]*Post)
	for scanner.Scan() {
		// 使用迭代器的方式遍历数据行
		text := scanner.Text()
		var post Post
		if err := json.Unmarshal([]byte(text), &post); err != nil {
			return err
		}
		if _, has := postsTmpMap[post.ParentID]; !has {
			postsTmpMap[post.ParentID] = make([]*Post, 0)
		}
		postsTmpMap[post.ParentID] = append(postsTmpMap[post.ParentID], &post)
		if post.ID > postID {
			postID = post.ID
		}
	}
	postIndexMap = postsTmpMap
	return nil
}
