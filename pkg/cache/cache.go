package cache

import (
	"sync"

	types "gitlab.eng.vmware.com/nidhig1/goassignment-postandcomments/pkg/types"
)

var postCacheArray []*types.Post

type PostsAndCommentsCache struct {
	mu       sync.RWMutex
	cacheMap map[int]*types.Post
}

func NewPostsAndCommentsCache() *PostsAndCommentsCache {
	return &PostsAndCommentsCache{
		mu:       sync.RWMutex{},
		cacheMap: make(map[int]*types.Post),
	}
}

func (c *PostsAndCommentsCache) AddAll(posts []types.Post) { // change it post id
	c.mu.Lock()
	for index := range posts {
		c.cacheMap[posts[index].Id] = &posts[index]
	}
	c.mu.Unlock()
}
func (c *PostsAndCommentsCache) GetAll() []*types.Post {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, value := range c.cacheMap {
		postCacheArray = append(postCacheArray, value)
	}
	return postCacheArray
}
func (c *PostsAndCommentsCache) GetOnePost(postId int) *types.Post {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i := range c.cacheMap {
		if c.cacheMap[i].Id == postId {
			return c.cacheMap[i]
		}
	}
	return nil
}
