package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-k8s-job/internal/biz"
)

type jobRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewJobRepo(data *Data, logger log.Logger) biz.JobRepo {
	return &jobRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *jobRepo) Save(ctx context.Context) error {

	err := CallJob(r)
	if err != nil {
		return err
	}

	return nil

}
