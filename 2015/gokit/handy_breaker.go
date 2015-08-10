// +build OMIT

package circuitbreaker

import (
	"time"

	"github.com/streadway/handy/breaker"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

// HandyBreaker returns an endpoint.Middleware that implements the circuit
// breaker pattern using the streadway/handy/breaker package. Only errors
// returned by the wrapped endpoint count against the circuit breaker's error
// count.
//
// See http://godoc.org/github.com/streadway/handy/breaker for more
// information.
// START OMIT
func HandyBreaker(cb breaker.Breaker) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !cb.Allow() { // HL
				return nil, breaker.ErrCircuitOpen // HL
			} // HL

			defer func(begin time.Time) {
				if err == nil {
					cb.Success(time.Since(begin))
				} else {
					cb.Failure(time.Since(begin))
				}
			}(time.Now())

			response, err = next(ctx, request)
			return
		}
	}
}

// STOP OMIT
