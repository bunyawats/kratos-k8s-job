package service

import (
	"context"
	"kratos-k8s-job/internal/biz"

	pb "kratos-k8s-job/api/scheduler/v1"
)

type JobService struct {
	pb.UnimplementedJobServer
	uc *biz.JobUseCase
}

func NewJobService(uc *biz.JobUseCase) *JobService {
	return &JobService{uc: uc}
}

func (s *JobService) ExecuteJob(ctx context.Context, req *pb.ExecuteJobRequest) (*pb.ExecuteJobReply, error) {
	err := s.uc.ExecuteJob(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ExecuteJobReply{}, nil
}
