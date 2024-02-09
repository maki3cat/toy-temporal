package main

import (
	"context"

	pb "github.com/maki3cat/toy-temporal/api-go/workflow"
)

func PollWorkflowTaskHandler(ctx context.Context, req *pb.PollWorkflowTaskRequest) (*pb.PollWorkflowTaskResponse, error) {

	taskEntity, err := workflowStorageInstance.pollPendingWorkflowTask()
	if err != nil {
		return nil, err
	}

	res := &pb.PollWorkflowTaskResponse{}
	if taskEntity == nil {
		return res, nil
	}
	res.TaskId = taskEntity.TaskId
	res.RunId = taskEntity.RunId
	res.Payload = taskEntity.TaskPayload

	return res, nil
}
