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

	MySqlAdapter interface {
		QueryMySqlDB(context.Context) ([]Message, error)
	}

	RabbitMqAdapter interface {
		SendMessage2RabbitMQ(context.Context, []Message) error
	}
)

type JobUseCase struct {
	MySqlAdapter
	RabbitMqAdapter
	log *log.Helper
}

func NewJobUseCase(
	m MySqlAdapter,
	r RabbitMqAdapter,
	logger log.Logger) *JobUseCase {

	return &JobUseCase{
		MySqlAdapter:    m,
		RabbitMqAdapter: r,
		log:             log.NewHelper(logger),
	}
}

func (uc *JobUseCase) ExecuteJob(ctx context.Context) error {
	uc.log.WithContext(ctx).Info("ExecuteJob")

	messageList, err := uc.QueryMySqlDB(ctx)
	if err != nil {
		return err
	}

	err = uc.SendMessage2RabbitMQ(ctx, messageList)
	if err != nil {
		return err
	}

	return nil
}
