package container

import (
	"os"
	"strconv"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"go.uber.org/zap"

	"codebase/internal/datastore"
	codebase "codebase/internal/models"
	"codebase/internal/service/mail"
	"codebase/internal/service/stock"

	"codebase/internal/service/utils"
	"codebase/pkg/caching"
	"codebase/pkg/db"
	"codebase/pkg/limiter"
	"codebase/pkg/logger"
)

func init() {
	//nolint:errcheck
	godotenv.Load("./.env") // for production
}

func NewContainer(vs map[string]string) *do.Injector {
	injector := do.New()
	do.ProvideNamedValue(injector, "envs", vs)

	do.ProvideValue(injector, &db.PostgresConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Database:     os.Getenv("DB_DATABASE"),
		PoolMaxConns: os.Getenv("DB_POOL_MAX_CONNS"),
	})

	do.Provide(injector, ProvidePgxConnectionPool)
	do.Provide(injector, ProvideReadOnlyDatastore)
	do.Provide(injector, ProvideDatastore)
	// do.Provide(injector, ProvideCache)
	// do.Provide(injector, ProvideLimiter)
	// do.Provide(injector, ProvideMutex)
	// do.Provide(injector, ProvideRedis)
	// do.Provide(injector, ProvideReadOnlyCache)
	do.Provide(injector, ProvideLogger)

	// SERVICE
	do.Provide(injector, ProvideServiceMail)
	do.Provide(injector, ProvideServiceUtils)
	do.Provide(injector, ProvideServiceStock)

	return injector
}

func ProvidePgxConnectionPool(i *do.Injector) (*pgxpool.Pool, error) {
	config := do.MustInvoke[*db.PostgresConfig](i)
	return db.InitPGXPool(config)
}

func ProvideReadOnlyDatastore(i *do.Injector) (codebase.ReadOnlyDatastore, error) {
	dsn := os.Getenv("DB_READONLY_HOST")
	if dsn == "" {
		return do.Invoke[codebase.Datastore](i)
	}

	roPool, err := db.InitPGXPoolFromDSN(dsn)
	if err != nil {
		return nil, err
	}
	return datastore.NewPgxDatastore(roPool)
}

func ProvideDatastore(i *do.Injector) (codebase.Datastore, error) {
	pool, err := do.Invoke[*pgxpool.Pool](i)
	if err != nil {
		return nil, err
	}

	return datastore.NewPgxDatastore(pool)
}

func ProvideRedis(_ *do.Injector) (redis.UniversalClient, error) {
	clusterCacheRedisURL := os.Getenv("CLUSTER_REDIS_URL")
	if clusterCacheRedisURL != "" {
		clusterOpts, err := redis.ParseClusterURL(clusterCacheRedisURL)
		if err != nil {
			return nil, err
		}
		return redis.NewClusterClient(clusterOpts), nil
	}

	return db.InitRedis(&db.RedisConfig{
		URL: os.Getenv("REDIS_URL"),
	})
}

func ProvideCache(_ *do.Injector) (caching.Cache, error) {
	var redisClient redis.UniversalClient
	clusterCacheRedisURL := os.Getenv("CLUSTER_CACHE_REDIS_URL")
	if clusterCacheRedisURL != "" {
		clusterOpts, err := redis.ParseClusterURL(clusterCacheRedisURL)
		if err != nil {
			return nil, err
		}
		redisClient = redis.NewClusterClient(clusterOpts)
	} else {
		var err error
		redisClient, err = db.InitRedis(&db.RedisConfig{
			URL: os.Getenv("CACHE_REDIS_URL"),
		})
		if err != nil {
			return nil, err
		}
	}

	return caching.NewCacheRedis(redisClient, false)
}

func ProvideReadOnlyCache(i *do.Injector) (caching.ReadOnlyCache, error) {
	// if cluster
	clusterCacheRedisURL := os.Getenv("CLUSTER_CACHE_REDIS_URL")
	if clusterCacheRedisURL != "" {
		var clusterOpts *redis.ClusterOptions
		var err error
		clusterCacheRedisReadOnlyURL := os.Getenv("CLUSTER_CACHE_REDIS_READONLY_URL")
		if clusterCacheRedisReadOnlyURL != "" {
			clusterOpts, err = redis.ParseClusterURL(clusterCacheRedisReadOnlyURL)
		} else {
			clusterOpts, err = redis.ParseClusterURL(clusterCacheRedisURL)
		}

		if err != nil {
			return nil, err
		}
		clusterOpts.ReadOnly = true
		redisClient := redis.NewClusterClient(clusterOpts)
		return caching.NewCacheRedis(redisClient, false)
	}

	// if cluster mode is not enabled, there is no "read_only" option
	return do.Invoke[caching.Cache](i)
}

func ProvideLimiter(_ *do.Injector) (codebase.Limiter, error) {
	var redisClient redis.UniversalClient
	clusterCacheRedisURL := os.Getenv("CLUSTER_LIMITER_REDIS_URL")
	if clusterCacheRedisURL != "" {
		clusterOpts, err := redis.ParseClusterURL(clusterCacheRedisURL)
		if err != nil {
			return nil, err
		}
		redisClient = redis.NewClusterClient(clusterOpts)
	} else {
		var err error
		redisClient, err = db.InitRedis(&db.RedisConfig{
			URL: os.Getenv("LIMITER_REDIS_URL"),
		})
		if err != nil {
			return nil, err
		}
	}

	return limiter.NewLimiter(redisClient)
}

func ProvideMutex(_ *do.Injector) (*redsync.Redsync, error) {
	var redisClient redis.UniversalClient
	clusterCacheRedisURL := os.Getenv("CLUSTER_MUTEX_REDIS_URL")
	if clusterCacheRedisURL != "" {
		clusterOpts, err := redis.ParseClusterURL(clusterCacheRedisURL)
		if err != nil {
			return nil, err
		}
		redisClient = redis.NewClusterClient(clusterOpts)
	} else {
		var err error
		redisClient, err = db.InitRedis(&db.RedisConfig{
			URL: os.Getenv("MUTEX_REDIS_URL"),
		})
		if err != nil {
			return nil, err
		}
	}

	pool := goredis.NewPool(redisClient)
	return redsync.New(pool), nil
}

func ProvideLogger(_ *do.Injector) (*zap.SugaredLogger, error) {
	defaultLogLevel := 0

	logLevelEnv := os.Getenv("LOG_LEVEL")
	if logLevelEnv == "" {
		return logger.NewLogger(defaultLogLevel)
	}

	logLevel, err := strconv.Atoi(logLevelEnv)
	if err != nil {
		return logger.NewLogger(defaultLogLevel)
	}

	return logger.NewLogger(logLevel)
}

// SERVICE
func ProvideServiceMail(i *do.Injector) (*mail.ServiceMail, error) {
	return mail.NewServiceMail(i)
}

func ProvideServiceUtils(i *do.Injector) (*utils.ServiceUtils, error) {
	return utils.NewServiceUtils(i)
}

func ProvideServiceStock(i *do.Injector) (stock.ServiceStock, error) {
	return stock.NewServiceStock(i)
}
