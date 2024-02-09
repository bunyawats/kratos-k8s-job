package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type (
	Message struct {
		TemPlateName string `json:"templateName"`
		Version      string `json:"version"`
	}

	JobRepo interface {
		QueryMySqlDB(context.Context) ([]Message, error)
		SendMessage2RabbitMQ(context.Context, []Message) error
		ReadInfluxDB(context.Context) error
	}
)

type JobUseCase struct {
	repo JobRepo
	log  *log.Helper
}

func NewJobUseCase(repo JobRepo, logger log.Logger) *JobUseCase {
	return &JobUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *JobUseCase) ExecuteJob(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("ExecuteJob")

	messageList, err := uc.repo.QueryMySqlDB(ctx)
	if err != nil {
		return err
	}

	err = uc.repo.SendMessage2RabbitMQ(ctx, messageList)
	if err != nil {
		return err
	}

	err = uc.repo.ReadInfluxDB(ctx)
	if err != nil {
		return err
	}

	return nil
}
