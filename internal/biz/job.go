package biz

import (
	"context"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-k8s-job/internal/utility"
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

	InfluxDbAdapter interface {
		ReadInfluxDB(context.Context) error
		WriteMetric2InfluxDB([]*influxdb3.Point) error
	}
)

type JobUseCase struct {
	MySqlAdapter
	RabbitMqAdapter
	InfluxDbAdapter
	log *log.Helper
}

func NewJobUseCase(m MySqlAdapter, r RabbitMqAdapter, i InfluxDbAdapter, logger log.Logger) *JobUseCase {
	return &JobUseCase{
		MySqlAdapter:    m,
		RabbitMqAdapter: r,
		InfluxDbAdapter: i,
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

	pts, err := utility.GetGoRuntimeMetrics()
	if err != nil {
		log.Warnf("Fail on get runtime matix: %v", err)
	}
	err = uc.WriteMetric2InfluxDB(pts)
	if err != nil {
		log.Warnf("Fail on write runtime matix to influxdb: %v", err)
	}

	return nil
}
