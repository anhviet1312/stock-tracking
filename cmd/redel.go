package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"

	"codebase/pkg/db"
)

func NewRedisDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "redel",
		Usage: "Delete multiple Redis keys either based on user_id or based on a specific pattern",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "user",
				Usage: "User ID",
			},
			&cli.StringFlag{
				Name:  "pattern",
				Usage: "Key pattern",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Print deleted keys",
			},
		},

		Action: redel,
	}
}

func redel(c *cli.Context) error {
	// Since this function uses Redis's KEYS (can be resource-intensive and might affect performance of Redis server),
	// this command is PROHIBITED if the environment is not clearly defined
	// More info at: https://redis.io/docs/latest/commands/keys/
	// environment := os.Getenv("MODE")
	// if environment == "" || strings.ToLower(environment) == "production" {
	// 	return errors.New("this environment is not allowed to run this command")
	// }

	// THIS COMMAND IS NOW SAFE TO BE RUN ON PRODUCTION //

	// Validate values
	keyPattern := c.String("pattern")
	verbose := c.Bool("verbose")

	if keyPattern == "" {
		return errors.New("keyPattern have specify either --user or --pattern")
	}

	// Init new redis client
	redisClient, err := InitRedisClient()
	if err != nil {
		return err
	}

	ctx := c.Context

	err = deleteKeyByPattern(ctx, redisClient, keyPattern, verbose)
	return err
}

func deleteKeyByPattern(ctx context.Context, redisClient redis.UniversalClient, pattern string, verbose bool) error {
	if pattern == "*" {
		return errors.New("* pattern is prohibited")
	}

	// Pattern input usually does not have quotes
	formattedPattern := fmt.Sprintf("*STORYWEB*%s", pattern)
	fmt.Println("Deleting keys by pattern: ", formattedPattern)

	deleted := 0
	redisClusterClient, ok := redisClient.(*redis.ClusterClient)
	if !ok {
		iter := redisClient.Scan(ctx, 0, formattedPattern, 100).Iterator()
		for iter.Next(ctx) {
			if verbose {
				fmt.Println("Deleting", iter.Val())
			}
			_ = redisClient.Del(ctx, iter.Val()).Err()
			deleted++
		}

		fmt.Printf("Redis single: %d deleted\n", deleted)
		return nil
	} else {
		err := redisClusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			iter := client.Scan(ctx, 0, formattedPattern, 100).Iterator()
			for iter.Next(ctx) {
				if verbose {
					fmt.Println("Deleting", iter.Val())
				}

				_ = client.Del(ctx, iter.Val()).Err()
				deleted++
			}
			return nil
		})

		fmt.Printf("Redis cluster: %d deleted\n", deleted)
		return err
	}
}

func InitRedisClient() (redis.UniversalClient, error) {
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

	return redisClient, nil
}
