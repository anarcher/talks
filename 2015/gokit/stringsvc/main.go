// +build OMIT
package main

import (
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"net/rpc"
	"strings"
)

var ErrEmpty = errors.New("empty string")

// S:SERVICE OMIT

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

// E:SERVICE OMIT

// S:RR OMIT

type UppercaseRequest struct {
	S string `json:"s"`
}

type UppercaseResponse struct {
	V   string `json:"v"`
	Err error  `json:"err"`
}

type CountRequest struct {
	S string `json:"s"`
}

type CountResponse struct {
	V int `json:"v"`
}

// E:RR OMIT

// S:ENDPOINT OMIT

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UppercaseRequest)
		v, err := svc.Uppercase(req.S)
		return UppercaseResponse{v, err}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CountRequest)
		v := svc.Count(req.S)
		return CountResponse{v}, nil
	}
}

// E:ENDPOINT OMIT

// S:NETRPC_BINDING1 OMIT

// NetRpc's handler
type NetRpcBinding struct {
	Context           context.Context
	uppercaseEndpoint endpoint.Endpoint
	countEndpoint     endpoint.Endpoint
}

// E:NETRPC_BINDING1 OMIT

// S:NETRPC_BINDING2 OMIT

func (n NetRpcBinding) Uppercase(req UppercaseRequest, res *UppercaseResponse) error { // HL
	ctx, cancel := context.WithCancel(n.Context)
	defer cancel()
	responses := make(chan UppercaseResponse, 1)
	errs := make(chan error, 1)

	go func() {
		resp, err := n.uppercaseEndpoint(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		responses <- resp.(UppercaseResponse)
	}()

	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*res) = resp
		return nil
	}
}

// E:NETRPC_BINDING2 OMIT

// S:NETRPC_BINDING3 OMIT

func (n NetRpcBinding) Count(req CountRequest, res *CountResponse) error { // HL
	ctx, cancel := context.WithCancel(n.Context)
	defer cancel()
	responses := make(chan CountResponse, 1)
	errs := make(chan error, 1)

	go func() {
		resp, err := n.countEndpoint(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		responses <- resp.(CountResponse)
	}()

	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*res) = resp
		return nil
	}
}

// E:NETRPC_BINDING3 OMIT

// S:MAIN1 OMIT

func main() {
	ctx := context.Background()
	svc := stringService{}

	//net/rpc
	netRpcBinding := NetRpcBinding{ctx, makeUppercaseEndpoint(svc), makeCountEndpoint(svc)} // HL

	s := rpc.NewServer()
	s.RegisterName("stringsvc", netRpcBinding) // HL
	s.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		log.Fatal(err)
	}

}

// E:MAIN1 OMIT
