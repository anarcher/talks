package main

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
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

type AddService interface {
	Add(context.Context int64, int64) int64
}

type addService struct{}
func (addService) add(_ context.Context, a, b int64) int64 { return a + b } // HL

// E:SERVICE OMIT
