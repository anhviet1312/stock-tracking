package utils

import (
	storyW "codebase/internal/models"

	"github.com/samber/do"
	"go.uber.org/zap"
)

type ServiceUtils struct {
	Datastore         storyW.Datastore
	ReadOnlyDatastore storyW.ReadOnlyDatastore

	// Cache         caching.Cache
	// ReadOnlyCache caching.ReadOnlyCache
	// Limiter       storyW.Limiter
	// RedisClient   redis.UniversalClient
	// Rs *redsync.Redsync

	Logger *zap.SugaredLogger
}

func NewServiceUtils(container *do.Injector) (*ServiceUtils, error) {
	datastore, err := do.Invoke[storyW.Datastore](container)
	if err != nil {
		return nil, err
	}

	readOnlyDatastore, err := do.Invoke[storyW.ReadOnlyDatastore](container)
	if err != nil {
		return nil, err
	}

	//readOnlyCache, err := do.Invoke[caching.ReadOnlyCache](container)
	//if err != nil {
	//	return nil, err
	//}
	//
	//cacher, err := do.Invoke[caching.Cache](container)
	//if err != nil {
	//	return nil, err
	//}
	//
	//redisClient, err := do.Invoke[redis.UniversalClient](container)
	//if err != nil {
	//	return nil, err
	//}
	//
	//limiter, err := do.Invoke[storyW.Limiter](container)
	//if err != nil {
	//	return nil, err
	//}
	//
	//rs, err := do.Invoke[*redsync.Redsync](container)
	//if err != nil {
	//	return nil, err
	//}

	logger, err := do.Invoke[*zap.SugaredLogger](container)
	if err != nil {
		return nil, err
	}

	return &ServiceUtils{
		Datastore:         datastore,
		ReadOnlyDatastore: readOnlyDatastore,
		// Limiter:           limiter,
		// Cache:             cacher,
		// ReadOnlyCache:     readOnlyCache,
		// RedisClient:       redisClient,
		// Rs:                rs,
		Logger: logger,
	}, nil
}

//func (service *ServiceUtils) TryDeleteCache(key string) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	_ = service.Cache.Delete(ctx, key)
//}
