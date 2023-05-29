package handlers

import (
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/kostjaigin/ultimategoservice/app/services/sales-api/handlers/v1/testgrp"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs a http.Handler with all applicatoin routes defined.
func APIMux(cfg APIMuxConfig) *httptreemux.ContextMux {
	// httptreemux context mux uses same fct signature as http package for its handlers //
	mux := httptreemux.NewContextMux()

	// bounding a route to mux //
	// when request comes in to /test, we handle it with h //
	mux.Handle(http.MethodGet, "/test", testgrp.Test)

	return mux
}
