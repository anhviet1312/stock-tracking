package utils

//
//import (
//	"context"
//	"database/sql"
//	"errors"
//
//	"github.com/google/uuid"
//	storyW "story-web/internal/models"
//	"story-web/internal/service/cache_key"
//	"story-web/pkg/caching"
//)
//
//func (service *ServiceUtils) FindUserByID(ctx context.Context, id uuid.UUID) (*storyW.User, error) {
//	return caching.UseCacheWithRO(
//		ctx,
//		service.ReadOnlyCache,
//		service.Cache,
//		cache_key.CacheKeyUserByID(id),
//		storyW.DEFAULT_CACHE_TTL,
//		func() (*storyW.User, error) {
//			user, err := service.ReadOnlyDatastore.FindUserByConditions(ctx, map[string]interface{}{
//				"id": id,
//			})
//			if err != nil && !errors.Is(err, sql.ErrNoRows) {
//				return nil, err
//			}
//
//			return user, nil
//		})
//}
//
//func (service *ServiceUtils) FindUserByUsername(ctx context.Context, username string) (*storyW.User, error) {
//	return caching.UseCacheWithRO(
//		ctx,
//		service.ReadOnlyCache,
//		service.Cache,
//		cache_key.CacheKeyUserByUsername(username),
//		storyW.DEFAULT_CACHE_TTL,
//		func() (*storyW.User, error) {
//			return service.ReadOnlyDatastore.FindUserByConditions(ctx, map[string]interface{}{
//				"username": username,
//			})
//		})
//}
//
//func (service *ServiceUtils) FindUserByEmail(ctx context.Context, email string) (*storyW.User, error) {
//	return caching.UseCacheWithRO(
//		ctx,
//		service.ReadOnlyCache,
//		service.Cache,
//		cache_key.CacheKeyUserByEmail(email),
//		storyW.DEFAULT_CACHE_TTL,
//		func() (*storyW.User, error) {
//			return service.ReadOnlyDatastore.FindUserByConditions(ctx, map[string]interface{}{
//				"email": email,
//			})
//		})
//}
