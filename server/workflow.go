package main

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/maki3cat/toy-temporal/api-go/workflow"
)

func StartWorkflowExecutionHandler(ctx context.Context, req *pb.StartWorkflowExecutionRequest) (*pb.StartWorkflowExecutionResponse, error) {

	runId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	taskId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	data := storeStartExecutionReq{
		ExecutionId:  req.ExecutionId,
		RunId:        runId.String(),
		StartPayload: req.Payload,
		TaskId:       taskId.String(),
	}
	err = workflowStorageInstance.storeStartExecution(data)
	if err != nil {
		return nil, err
	}
	// I don't understand why Temporal have so many parts, let's start simple
	// transactional:
	// 1) create worklfow execution state machine
	// 2) insert workflow Events ?

	// 3) update workflow task in execution
	// 4) assign workflow task queue

	// 5) assign workflow task start event ?
	res := pb.StartWorkflowExecutionResponse{
		RunId: data.RunId,
	}
	return &res, nil
}
