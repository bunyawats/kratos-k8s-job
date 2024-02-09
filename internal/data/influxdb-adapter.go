package data

import (
	"context"
	"fmt"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-k8s-job/internal/biz"
)

type (
	I struct {
		data *Data
		log  *log.Helper
	}
)

func NewInfluxDbAdapter(data *Data, logger log.Logger) biz.InfluxDbAdapter {
	return &I{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *I) ReadInfluxDB(ctx context.Context) error {

	// Execute query
	query := `SELECT *
          FROM 'census'
          WHERE time >= now() - interval '12 hour'
            AND ('bees' IS NOT NULL OR 'ants' IS NOT NULL)`

	queryOptions := influxdb3.QueryOptions{
		Database: r.data.Bucket,
	}
	iterator, err := r.data.InfluxDBClient.QueryWithOptions(context.Background(), &queryOptions, query)

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
