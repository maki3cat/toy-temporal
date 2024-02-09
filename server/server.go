package main

import (
	"context"
	"flag"

	pb "github.com/maki3cat/toy-temporal/api-go/workflow"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type workflowServer struct {
	pb.UnimplementedWorkflowServer
}

func (workflowServer) StartWorkflowExecution(ctx context.Context, req *pb.StartWorkflowExecutionRequest) (*pb.StartWorkflowExecutionResponse, error) {
	return StartWorkflowExecutionHandler(ctx, req)
}

func (workflowServer) PollWorkflowTask(ctx context.Context, req *pb.PollWorkflowTaskRequest) (*pb.PollWorkflowTaskResponse, error) {
	return PollWorkflowTaskHandler(ctx, req)
}
