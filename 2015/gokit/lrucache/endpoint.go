// +build OMIT

package main

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/hashicorp/golang-lru"
	"golang.org/x/net/context"
	"log"
)

// START1 OMIT
type CacheKeyFunc func(request interface{}) (interface{}, bool)

func NewLRUCache(cache *lru.Cache, cacheKey CacheKeyFunc) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			key, ok := cacheKey(request)
			if !ok {
				return next(ctx, request)
			}

			val, ok := cache.Key(key)
			if !ok {
				return next(ctx, request)
			}
			return val
		}
	}
}

// STOP1 OMIT
