package server

import (
	"fmt"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewApiRenderer() func(w http.ResponseWriter, r *http.Request, v interface{}) {
	return func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(http.StatusInternalServerError)
			}

			log.Error(fmt.Errorf("request processing failed: %w", err))
			render.DefaultResponder(w, r, createErrorResponse(fmt.Errorf("something went wrong")))
			return
		}

		render.DefaultResponder(w, r, v)
	}
}
