package caching

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type ReadOnlyCache interface {
	Get(ctx context.Context, key string, target any) error
	GetInt64(ctx context.Context, key string) (int64, error)
}

type Cache interface {
	ReadOnlyCache
	Set(ctx context.Context, key string, value any, ttl time.Duration, tags ...string) error
	SetInt64(ctx context.Context, key string, value int64, ttl time.Duration) error
	Incr(ctx context.Context, key string) error
	Decr(ctx context.Context, key string) error
	Delete(ctx context.Context, key string) error
	DeleteByTags(ctx context.Context, tags ...string) error

	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
}

func UseCache[T any](ctx context.Context, cash Cache, key string, ttl time.Duration, callback func() (T, error)) (T, error) {
	var v T
	err := cash.Get(ctx, key, &v)
	if !errors.Is(err, cache.ErrCacheMiss) {
		return v, err
	}

	v, err = callback()
	if err != nil {
		return v, err
	}

	// fire and forget
	//nolint:errcheck
	cash.Set(ctx, key, v, ttl)
	return v, nil
}

func UseCacheWithTags[T any](ctx context.Context, cash Cache, key string, ttl time.Duration, callback func() (T, []string, error)) (T, error) {
	var v T
	err := cash.Get(ctx, key, &v)
	if !errors.Is(err, cache.ErrCacheMiss) {
		return v, err
	}

	v, tags, err := callback()
	if err != nil {
		return v, err
	}

	// fire and forget
	//nolint:errcheck
	cash.Set(ctx, key, v, ttl, tags...)
	return v, nil
}

func UseCacheWithRO[T any](ctx context.Context, roCash ReadOnlyCache, cash Cache, key string, ttl time.Duration, callback func() (T, error)) (T, error) {
	var v T
	err := roCash.Get(ctx, key, &v)
	if !errors.Is(err, cache.ErrCacheMiss) {
		return v, err
	}

	v, err = callback()
	if err != nil {
		return v, err
	}

	// fire and forget
	//nolint:errcheck
	cash.Set(ctx, key, v, ttl)
	return v, nil
}

func UseCacheGroup[T any](ctx context.Context, roCash ReadOnlyCache, cash Cache, groupKey, key string, ttl time.Duration, callback func() (T, error)) (T, error) {
	var v T
	err := roCash.Get(ctx, key, &v)
	if !errors.Is(err, cache.ErrCacheMiss) {
		return v, err
	}

	v, err = callback()
	if err != nil {
		return v, err
	}

	// fire and forget
	//nolint:errcheck
	err = cash.Set(ctx, key, v, ttl)
	if err != nil {
		return v, nil
	}

	// add key to group keys
	err = cash.SAdd(ctx, groupKey, key)
	if err != nil {
		return v, nil
	}

	return v, nil
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func DeleteCacheGroup(ctx context.Context, cash Cache, groupKey string) error {
	keys, err := cash.SMembers(ctx, groupKey)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		maxRoutines := 1000
		chunks := chunkBy(keys, 1+len(keys)/maxRoutines)
		wg := &sync.WaitGroup{}
		for _, chunk := range chunks {
			if len(chunk) == 0 {
				continue
			}
			wg.Add(1)
			go func(ctx context.Context, wg *sync.WaitGroup, keys ...string) {
				defer wg.Done()
				for _, key := range keys {
					_ = cash.Delete(ctx, key)
				}
			}(ctx, wg, chunk...)
		}

		wg.Wait()
	}

	// return cash.Delete(ctx, groupKey)
	return nil
}

type CacheRedis struct {
	client   redis.UniversalClient
	instance *cache.Cache
}

func (c *CacheRedis) GetInt64(ctx context.Context, key string) (int64, error) {
	return c.client.Get(ctx, key).Int64()
}

func (c *CacheRedis) Get(ctx context.Context, key string, target any) error {
	return c.instance.Get(ctx, key, target)
}

func (c *CacheRedis) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SAdd(ctx, key, members).Err()
}

func (c *CacheRedis) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

func (c *CacheRedis) SetInt64(ctx context.Context, key string, value int64, ttl time.Duration) error {
	_, err := c.client.Set(ctx, key, fmt.Sprintf("%d", value), ttl).Result()
	return err
}

func (c *CacheRedis) Incr(ctx context.Context, key string) error {
	_, err := c.client.Incr(ctx, key).Result()
	return err
}

func (c *CacheRedis) Decr(ctx context.Context, key string) error {
	_, err := c.client.Decr(ctx, key).Result()
	return err
}

func (c *CacheRedis) Set(ctx context.Context, key string, value any, ttl time.Duration, tags ...string) error {
	if len(tags) == 0 {
		return c.instance.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: value,
			TTL:   ttl,
		})
	}

	_, err := c.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, tag := range tags {
			err := p.SAdd(ctx, fmt.Sprintf("cache:tag:%s", tag), key).Err()
			if err != nil {
				return err
			}
		}

		return c.instance.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: value,
			TTL:   ttl,
		})
	})

	return err
}

func (c *CacheRedis) Delete(ctx context.Context, key string) error {
	return c.instance.Delete(ctx, key)
}

func (c *CacheRedis) DeleteByTags(ctx context.Context, tags ...string) error {
	if len(tags) == 0 {
		return nil
	}

	redisClusterClient, ok := c.client.(*redis.ClusterClient)
	if !ok {
		_, err := c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
			for _, tag := range tags {
				key := fmt.Sprintf("cache:tag:%s", tag)
				members, err := c.client.SMembers(ctx, key).Result()
				if err != nil {
					return err
				}

				if len(members) > 0 {
					err = c.client.Del(ctx, members...).Err()
					if err != nil {
						return err
					}
				}

				err = c.client.Del(ctx, key).Err()
				if err != nil {
					return err
				}
			}

			return nil
		})

		return err
	}

	err := redisClusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
		for _, tag := range tags {
			key := fmt.Sprintf("cache:tag:%s", tag)
			members, err := c.client.SMembers(ctx, key).Result()
			if err != nil {
				return err
			}

			if len(members) > 0 {
				for _, member := range members {
					//nolint:errcheck
					client.Del(ctx, member)
				}
			}

			//nolint:errcheck
			client.Del(ctx, key)
		}

		return nil
	})

	return err
}

func NewCacheRedis(client redis.UniversalClient, withLocalCache bool) (*CacheRedis, error) {
	var localCache cache.LocalCache
	if withLocalCache {
		localCache = cache.NewTinyLFU(10000, time.Minute)
	}

	return &CacheRedis{
		client: client,
		instance: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: localCache,
		}),
	}, nil
}
