package cache

import (
	"ai-calls/internal/config"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Cache struct {
	Config *config.Config
	db     *ristretto.Cache[string, any]
	prefix string
}

func New(cfg *config.Config) *Cache {
	cache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
		Cost: func(value any) int64 {
			switch v := value.(type) {
			case string:
				return int64(len(v))
			case int:
				return 8
			default:
				return 1
			}
		},
	})
	if err != nil {
		panic(err)
	}
	return &Cache{db: cache, Config: cfg}
}

func (c *Cache) Prefix(prefix string) *Cache {
	return &Cache{
		Config: c.Config,
		db:     c.db,
		prefix: prefix,
	}
}

func (c *Cache) getKeyWithPrefixOrRaw(key string) string {
	if c.prefix != "" {
		return fmt.Sprintf("%s:%s", c.prefix, key)
	}
	return key
}

func (c *Cache) Set(key string, val any) {
	c.db.Set(c.getKeyWithPrefixOrRaw(key), val, 1)
	c.db.Wait()
}

func (c *Cache) SetWithTTL(key string, val any, ttl time.Duration) {
	c.db.SetWithTTL(c.getKeyWithPrefixOrRaw(key), val, 1, ttl)
	c.db.Wait()
}

func (c *Cache) Get(key string) any {
	value, found := c.db.Get(c.getKeyWithPrefixOrRaw(key))
	if !found {
		return nil
	}
	return value
}

func (c *Cache) GetValueAndDelete(key string) any {
	value := c.Get(key)
	c.Delete(key)
	return value
}

func (c *Cache) Delete(key string) {
	c.db.Del(c.getKeyWithPrefixOrRaw(key))
}
