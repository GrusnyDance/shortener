package cache

import (
	"context"
	"errors"
	"sync"
)

type Cache struct {
	Map map[string]string
	sync.RWMutex
}

func Init() *Cache {
	m := make(map[string]string)
	return &Cache{
		Map: m,
	}
}

func (c *Cache) ReturnLink(ctx context.Context, hash string) (string, error) {
	link := c.Map[hash]
	return link, nil
}

func (c *Cache) CheckIfHashedExists(ctx context.Context, hash string) error {
	if _, exist := c.Map[hash]; !exist {
		return errors.New("link not found")
	}
	return nil
}

func (c *Cache) CreateLink(ctx context.Context, hashed string, original string) error {
	c.Lock()
	c.Map[hashed] = original
	c.Unlock()
	return nil
}
