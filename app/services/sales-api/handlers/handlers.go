package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) http.Handler {
	// httptreemux context mux uses same fct signature as http package for its handlers //
	mux := httptreemux.NewContextMux()

	h := func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status string
		}{
			Status: "OK",
		}

		json.NewEncoder(w).Encode(status)
	}
	// bounding a route to mux //
	// when request comes in to /test, we handle it with h //
	mux.Handle(http.MethodGet, "/test", h)

	return mux
}
