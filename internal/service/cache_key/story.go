package cache_key

import "fmt"

func CacheKeyStoryBySlug(slug string) string {
	return fmt.Sprintf("[STORYWEB]:story:slug:%s", slug)
}
