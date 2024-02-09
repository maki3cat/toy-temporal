package main

import (
	"context"

	"github.com/gogo/status"
	pb "github.com/maki3cat/toy-temporal/api-go/workflow"
	"google.golang.org/grpc/codes"
)

func PollWorkflowTaskHandler(ctx context.Context, pb *pb.PollWorkflowTaskRequest) (*pb.PollWorkflowTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PollWorkflowTask not implemented")
}
