package repository

import (
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Topic struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type TopicDao struct{}

var topicDao *TopicDao
var topicOnce sync.Once
var mutex sync.Mutex

// 设计模式中的单例模式
func NewTopicDaoInstance() *TopicDao {
	topicOnce.Do(
		func() {
			topicDao = &TopicDao{}
		})
	return topicDao
}

func (*TopicDao) QueryTopicByID(id int64) *Topic {
	return topicIndexMap[id]
}

func (*TopicDao) AddTopic(title, content string) error {
	atomic.AddInt64(&topicID, 1)
	timeStamp := time.Now().Unix()
	newTopic := &Topic{ID: topicID, Title: title, Content: content, CreateTime: timeStamp}
	// json序列化
	bytes, err := json.Marshal(newTopic)
	if err != nil {
		return err
	}

	mutex.Lock()
	topicIndexMap[topicID] = newTopic
	// 添加换行符
	bytes = append(bytes, '\n')
	// 追加的方式打开文件
	f, err := os.OpenFile(filePath+"post", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(bytes)
	mutex.Unlock()
	if err != nil {
		return err
	}

	return nil
}
