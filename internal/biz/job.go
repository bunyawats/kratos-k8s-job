package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type JobRepo interface {
	Save(context.Context) error
}

type JobUseCase struct {
	repo JobRepo
	log  *log.Helper
}

func NewJobUseCase(repo JobRepo, logger log.Logger) *JobUseCase {
	return &JobUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *JobUseCase) ExecuteJob(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("ExecuteJob")
	return uc.repo.Save(ctx)
}
