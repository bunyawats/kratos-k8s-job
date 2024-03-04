package server

import (
	"context"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"kratos-k8s-job/internal/conf"
	"kratos-k8s-job/internal/utility"
)

type (
	influxDbMiddleware struct {
		InfluxDBClient *influxdb3.Client
		Bucket         string
		log            *log.Helper
	}
)

func NewInfluxDbMiddleware(c *conf.Data, logger log.Logger) (*influxDbMiddleware, func(), error) {

	l := log.NewHelper(logger)

	influxCf := c.Influxdb
	l.Debug("influxdb address: ", influxCf.GetAddr())
	l.Debug("influxdb token: ", influxCf.GetToken())
	influxClient, err := influxdb3.New(influxdb3.ClientConfig{
		Host:  influxCf.Addr,
		Token: influxCf.Token,
	})
	if err != nil {
		l.Error("Fail on connect to InfluxDB")
		return nil, nil, err
	}

	cleanup := func() {
		l.Info("closing Middleware InfluxDB connection")
		influxClient.Close()
	}

	return &influxDbMiddleware{
		InfluxDBClient: influxClient,
		Bucket:         influxCf.Bucket,
		log:            l,
	}, cleanup, nil
}

func (ifm *influxDbMiddleware) WriteMetric2InfluxDB() {

	rtmPoints, err := utility.GetGoRuntimeMetrics()
	if err != nil {
		log.Warnf("Fail on get runtime matix: %v", err)
	}

	options := influxdb3.WriteOptions{
		Database: ifm.Bucket,
	}

	if err := ifm.InfluxDBClient.WritePointsWithOptions(context.Background(), &options, rtmPoints...); err != nil {
		log.Warnf("error while writing point to InfluxD: %v", err)
	}

	if err != nil {
		ifm.log.Warnf("Fail on write runtime matix to influxdb: %v", err)
	}
}

func (ifm *influxDbMiddleware) runtimeMetricInfluxDbMiddleware(handler middleware.Handler) middleware.Handler {

	return func(ctx context.Context, req interface{}) (interface{}, error) {

		ifm.log.Debug("\nRuntime Metric InfluxDb middleware in", req)
		reply, err := handler(ctx, req)
		ifm.WriteMetric2InfluxDB()
		ifm.log.Debug("Runtime Metric InfluxDb middleware out", reply)

		return reply, err
	}
}
