package data

import (
	"context"
	"fmt"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-k8s-job/internal/biz"
	"kratos-k8s-job/internal/conf"
)

type (
	iAdapter struct {
		InfluxDBClient *influxdb3.Client
		Bucket         string
		log            *log.Helper
	}
)

func NewInfluxDbAdapter(c *conf.Data, logger log.Logger) (biz.InfluxDbAdapter, func(), error) {

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
		l.Info("closing InfluxDB connection")
		influxClient.Close()
	}

	return &iAdapter{
		InfluxDBClient: influxClient,
		Bucket:         influxCf.Bucket,
	}, cleanup, nil
}

func (i *iAdapter) ReadInfluxDB(ctx context.Context) error {

	// Execute query
	query := `SELECT *
          FROM 'census'
          WHERE time >= now() - interval '12 hour'
            AND ('bees' IS NOT NULL OR 'ants' IS NOT NULL)`

	queryOptions := influxdb3.QueryOptions{
		Database: i.Bucket,
	}
	iterator, err := i.InfluxDBClient.QueryWithOptions(context.Background(), &queryOptions, query)

	if err != nil {
		panic(err)
	}

	for iterator.Next() {
		value := iterator.Value()

		location := value["location"]
		ants := value["ants"]
		bees := value["bees"]
		fmt.Printf("in %s are %d ants and %d bees\n", location, ants, bees)
	}

	return nil
}

func (i *iAdapter) WriteMatrix2InfluxDB(runtimeMetrics map[string]interface{}) error {

	options := influxdb3.WriteOptions{
		Database: i.Bucket,
	}

	point := influxdb3.NewPointWithMeasurement("metrics")
	for metricName, metricValue := range runtimeMetrics {
		switch metricValue.(type) {
		case string:
			point.SetTag(metricName, fmt.Sprintf("%v", metricValue))
		default:
			point.SetField(metricName, metricValue)
		}

	}

	if err := i.InfluxDBClient.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		log.Warnf("error while writing point to InfluxD: %v", err)
		return err
	}
	return nil
}
