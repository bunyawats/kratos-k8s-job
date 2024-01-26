package data

import (
	"database/sql"
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
	MySqlDB  *sql.DB
	AmqpConn *amqp091.Connection
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	dbCf := c.Database
	db, err := sql.Open(dbCf.GetDriver(), dbCf.GetSource())
	if err != nil {
		log.NewHelper(logger).Error("Fail on connect to MySql")
		return nil, nil, err
	}

	amqpCf := c.Amqp
	conn, err := amqp.Dial("amqp://" + amqpCf.GetAddr())
	if err != nil {
		log.NewHelper(logger).Error("Fail on connect to RabbitMq")
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		db.Close()
		conn.Close()
	}

	return &Data{
		MySqlDB:  db,
		AmqpConn: conn,
	}, cleanup, nil
}
