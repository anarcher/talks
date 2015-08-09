// +build OMIT

package main

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/hashicorp/golang-lru"
	"golang.org/x/net/context"
)

// S:SERVICE OMIT
type AddRequest struct {
	A int64 `json:"a"`
	B int64 `json:"b"`
}

type AddResponse struct {
	V int64 `json:"v"`
}

type Add func(context.Context, int64, int64) int64

func add(ctx context.Context, a, b int64) int64 { return a + b }

// E:SERVICE OMIT

// START1 OMIT
type CacheKeyFunc func(request interface{}) (interface{}, bool)

func NewLRUCache(cache *lru.Cache, cacheKey CacheKeyFunc) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			key, ok := cacheKey(request)
			if !ok {
				return next(ctx, request)
			}
			val, ok := cache.Get(key)
			if ok {
				fmt.Println("Return from cache", request, val)
				return val, nil
			}
			val, err := next(ctx, request)
			if err != nil {
				return val, err
			}
			cache.Add(key, val)
			fmt.Println("Return from endpoint", request, val)
			return val, err
		}
	}
}

// STOP1 OMIT

func main() {
	// S:MAIN OMIT
	cacheKeyFunc := func(request interface{}) (interface{}, bool) {
		if req, ok := request.(AddRequest); ok {
			return fmt.Sprintf("%d+%d", req.A, req.B), true
		}
		return nil, false
	}

	cache, _ := lru.New(10)
	e := makeEndpoint(add)
	e = NewLRUCache(cache, cacheKeyFunc)(e)

	req := AddRequest{1, 2}
	resp, err := e(context.Background(), req)
	fmt.Println("resp", resp, "err", err)

	resp, err = e(context.Background(), req)
	fmt.Println("resp", resp, "err", err)

	// E:MAIN OMIT

}

// S:makeEndpoint OMIT
func makeEndpoint(add Add) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		addReq := request.(AddRequest)
		v := add(ctx, addReq.A, addReq.B)
		addRes := AddResponse{V: v}
		return addRes, nil
	}
}

// E:makeEndpoint OMIT
