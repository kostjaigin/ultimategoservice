package handlers

import (
	"net/http"
	"os"

	"github.com/kostjaigin/ultimategoservice/app/services/sales-api/handlers/v1/testgrp"
	"github.com/kostjaigin/ultimategoservice/business/web/v1/mid"
	"github.com/kostjaigin/ultimategoservice/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs a http.Handler with all applicatoin routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	// httptreemux context mux uses same fct signature as http package for its handlers //
	// we take test, wrap error around it, then wrap logger around it //
	// from right to left // (with panics coming as close as possible to the handler call)
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	// bounding a route to mux //
	// when request comes in to /test, we handle it with h //
	app.Handle(http.MethodGet, "/test", testgrp.Test)

	return app
}
