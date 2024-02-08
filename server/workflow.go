package main

import (
	"context"
	"flag"

	"github.com/gogo/status"
	pb "github.com/maki3cat/toy-temporal/api-go/workflow"

	"google.golang.org/grpc/codes"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type workflowExecutionServer struct {
	pb.UnimplementedWorkflowExecutionServer
}

func (workflowExecutionServer) StartWorkflowExecution(context.Context, *pb.StartWorkflowExecutionRequest) (*pb.StartWorkflowExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartWorkflowExecution not implemented")
}
