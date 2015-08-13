// +build OMIT

package main

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/hashicorp/golang-lru"
	"golang.org/x/net/context"
	"strings"
)

// S:SERVICE OMIT
var ErrEmpty = errors.New("empty string")

type UppercaseRequest struct {
	S string `json:"s"`
}

type UppercaseResponse struct {
	V   string `json:"v"`
	Err error  `json:"err"`
}

type StringService interface { // HL
	Uppercase(string) (string, error) // HL
	Count(string) int                 // HL
} // HL

type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UppercaseRequest)
		v, err := svc.Uppercase(req.S)
		return UppercaseResponse{v, err}, nil
	}
}

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
	svc := stringService{}
	// S:MAIN OMIT
	cacheKeyFunc := func(request interface{}) (interface{}, bool) {
		if req, ok := request.(UppercaseRequest); ok {
			return req.S, true
		}
		return nil, false
	}

	cache, _ := lru.New(10)
	e := makeUppercaseEndpoint(svc)
	e = NewLRUCache(cache, cacheKeyFunc)(e)

	req := UppercaseRequest{"gophercon!"}
	resp, err := e(context.Background(), req)
	fmt.Println("resp", resp.(UppercaseResponse).V, "err", err)

	resp, err = e(context.Background(), req)
	fmt.Println("resp", resp.(UppercaseResponse).V, "err", err)

	// E:MAIN OMIT

}
