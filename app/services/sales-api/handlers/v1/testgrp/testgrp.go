package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	v1 "github.com/kostjaigin/ultimategoservice/business/web/v1"
	"github.com/kostjaigin/ultimategoservice/foundation/web"
)

// Test is our example route.
func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		// can also try with untrusted errors: return errors.New("UNTRUSTED ERROR")
		// or a panic to see how we transform it into an error: panic("SOME PANIC")
		return v1.NewRequestError(errors.New("TRUSTED ERROR"), http.StatusBadRequest)
	}

	// Validate the data
	// Call into the business layer

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
