package repository

import (
	"github.com/coocood/freecache"
)

type FreeCache struct {
	cache *freecache.Cache
}

func NewFreeCache(cache *freecache.Cache) *FreeCache {
	return &FreeCache{cache: cache}
}

func (c *FreeCache) Get(uuid []byte) ([]byte, error) {
	got, err := c.cache.Get(uuid)
	return got, err
}

func (c *FreeCache) Set(key, val []byte, expireIn int) error {
	err := c.cache.Set(key, val, expireIn)
	if err != nil {
		return err
	}
	return nil
}

func (c *FreeCache) Del(key []byte) (affected bool) {
	return c.cache.Del(key)
}
