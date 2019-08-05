package http

import (
	"context"
	"net/http"

	transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type Handler struct {
	Address string
	Router  *mux.Router
}

const (
	RouteHealth = "/health"
)

func NewHandler(address string) *Handler {
	health := transport.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			return nil, nil
		},
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
			w.WriteHeader(http.StatusOK)
			return nil
		},
	)

	router := mux.NewRouter()
	router.Methods("GET").Path(RouteHealth).Handler(health)

	return &Handler{
		Address: address,
		Router:  router,
	}
}
