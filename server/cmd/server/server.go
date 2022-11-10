package server

import (
	"net/http"

	"github.com/nndergunov/auctuionApp/server/api"
)

type Repository interface {
	Start() error
}

// Server is the server.
type Server struct {
	api *api.API
}

// NewServer creates a new server.
func NewServer(api *api.API) *Server {
	return &Server{
		api: api,
	}
}

// Start starts the server on the given address.
func (s *Server) Start(address string) error {
	return http.ListenAndServe(address, s.api)
}
