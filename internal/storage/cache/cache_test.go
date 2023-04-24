package cache_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"shortener/internal/storage/cache"
	"testing"
)

func TestInit(t *testing.T) {
	c := cache.Init()
	assert.NotEqual(t, c.Map, nil, "expected cache.Map to be initialized, got nil")
}

func TestReturnLink(t *testing.T) {
	c := cache.Init()
	c.Map["hash1"] = "link1"
	c.Map["hash2"] = "link2"

	link, err := c.ReturnLink(context.Background(), "hash1")
	assert.NoError(t, err)
	assert.Equal(t, link, "link1")

	link, err = c.ReturnLink(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, link, "")
}

func TestCheckIfHashedExists(t *testing.T) {
	c := cache.Init()
	c.Map["hash1"] = "link1"

	err := c.CheckIfHashedExists(context.Background(), "hash1")
	assert.NoError(t, err)

	err = c.CheckIfHashedExists(context.Background(), "nonexistent")
	assert.Error(t, err)
}

func TestCreateLink(t *testing.T) {
	c := cache.Init()

	err := c.CreateLink(context.Background(), "hash1", "link1")
	assert.NoError(t, err)
	assert.Equal(t, c.Map["hash1"], "link1")
}

func TestClose(t *testing.T) {
	c := cache.Init()
	c.Map["hash1"] = "link1"
	c.Close()

	assert.Len(t, c.Map, 0)
}
