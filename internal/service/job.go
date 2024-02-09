package service

import (
	"context"

	pb "kratos-k8s-job/api/scheduler/v1"
)

type JobService struct {
	pb.UnimplementedJobServer
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) ExecuteJob(ctx context.Context, req *pb.ExecuteJobRequest) (*pb.ExecuteJobReply, error) {
	return &pb.ExecuteJobReply{}, nil
}
