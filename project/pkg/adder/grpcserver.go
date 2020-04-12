package adder

import (
	"golang.org/x/net/context"
)

// GRPCServer ..
type GRPCServer struct {

}

// Add Запускаем...
func (s *GRPCServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	return
}