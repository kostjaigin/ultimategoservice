package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/kostjaigin/ultimategoservice/business/web/metrics"
	"github.com/kostjaigin/ultimategoservice/foundation/web"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		// See how we named the error - with the named variable we can change it's value in deferred call
		// without returning in the function
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

					metrics.AddPanics(ctx)
				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
