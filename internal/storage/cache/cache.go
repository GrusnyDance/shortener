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
	c.RLock()
	defer c.RUnlock()
	link := c.Map[hash]
	if link == "" {
		return link, errors.New("link not found")
	}
	return link, nil
}

func (c *Cache) CheckIfHashedExists(ctx context.Context, hash string) error {
	c.RLock()
	defer c.RUnlock()
	if _, exist := c.Map[hash]; !exist {
		return errors.New("link not found")
	}
	return nil
}

func (c *Cache) CreateLink(ctx context.Context, hashed string, original string) error {
	c.Lock()
	defer c.Unlock()
	c.Map[hashed] = original
	return nil
}

func (c *Cache) Close() {
	// Gives pointer to previous map to GC
	c.Map = make(map[string]string)
}
