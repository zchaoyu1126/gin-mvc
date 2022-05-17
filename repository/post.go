package repository

import "sync"

type Post struct {
	ID         int64  `json:"id"`
	ParentID   int64  `json:"parent_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"created_time"`
}

type PostsDao struct{}

var postsDao *PostsDao
var postsOnce sync.Once

// 设计模式中的单例模式
func NewPostsDaoInstance() *PostsDao {
	postsOnce.Do(
		func() {
			postsDao = &PostsDao{}
		})
	return postsDao
}

func (*PostsDao) QueryPostsByParentID(id int64) []*Post {
	return postIndexMap[id]
}
