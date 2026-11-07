package grpcapi

import (
	"context"
	"github.com/rinatkh/homework_go_18/internal/document"
	documentv1 "github.com/rinatkh/homework_go_18/pkg/proto/document/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	documentv1.UnimplementedDocumentServiceServer
	UseCase *document.UseCase
}

// TODO-05: validate positive id, call UseCase.Get(ctx,id), map NotFound/Internal, map DTO.
func (s *Server) GetDocument(ctx context.Context, r *documentv1.GetDocumentRequest) (*documentv1.DocumentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO-05")
}
