package server

import (
	"fmt"
	"net"

	"github.com/jian-hua-he/reference-app/internal/adapter/grpc/handler"
	notev1 "github.com/jian-hua-he/reference-app/proto/note/v1"

	"google.golang.org/grpc"
)

type Server struct {
	grpcPort int
	handler  *handler.Handler
	server   *grpc.Server
}

func NewServer(grpcPort int, h *handler.Handler) *Server {
	return &Server{
		grpcPort: grpcPort,
		handler:  h,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.server = grpc.NewServer()
	notev1.RegisterNoteServiceServer(s.server, s.handler)

	return s.server.Serve(lis)
}

func (s *Server) Shutdown() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
