package service

import (
	"context"
	"kratos-k8s-job/internal/biz"

	pb "kratos-k8s-job/api/scheduler/v1"
)

type JobService struct {
	pb.UnimplementedJobServer
	uc *biz.GreeterUsecase
}

func NewJobService(uc *biz.GreeterUsecase) *JobService {
	return &JobService{uc: uc}
}

func (s *JobService) ExecuteJob(ctx context.Context, req *pb.ExecuteJobRequest) (*pb.ExecuteJobReply, error) {
	_, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: "Job Message"})
	if err != nil {
		return nil, err
	}
	return &pb.ExecuteJobReply{}, nil
}
