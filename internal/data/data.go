package data

import (
	"database/sql"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/rabbitmq/amqp091-go"
	"kratos-k8s-job/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	MySqlDB        *sql.DB
	AmqpConn       *amqp091.Connection
	InfluxDBClient *influxdb3.Client
	Bucket         string
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	l := log.NewHelper(logger)

	dbCf := c.Database
	l.Debug("mysql source: ", dbCf.GetSource())
	db, err := sql.Open(dbCf.GetDriver(), dbCf.GetSource())
	if err != nil {
		l.Error("Fail on connect to MySql")
		return nil, nil, err
	}

	amqpCf := c.Amqp
	l.Debug("rabbitmq address: ", amqpCf.GetAddr())
	conn, err := amqp.Dial("amqp://" + amqpCf.GetAddr())
	if err != nil {
		l.Error("Fail on connect to RabbitMq")
		return nil, nil, err
	}

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
		l.Info("closing the data resources")
		db.Close()
		conn.Close()
		influxClient.Close()
	}

	return &Data{
		MySqlDB:        db,
		AmqpConn:       conn,
		InfluxDBClient: influxClient,
		Bucket:         influxCf.Bucket,
	}, cleanup, nil
}
