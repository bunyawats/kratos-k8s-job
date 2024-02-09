package data

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"kratos-k8s-job/internal/biz"
	"kratos-k8s-job/internal/conf"
	"time"
)

type (
	rAdapter struct {
		AmqpConn *amqp.Connection
		log      *log.Helper
	}
)

func NewRabbitMqAdapter(c *conf.Data, logger log.Logger) (biz.RabbitMqAdapter, func(), error) {

	l := log.NewHelper(logger)

	amqpCf := c.Amqp
	l.Debug("rabbitmq address: ", amqpCf.GetAddr())
	conn, err := amqp.Dial("amqp://" + amqpCf.GetAddr())
	if err != nil {
		l.Error("Fail on connect to RabbitMq")
		return nil, nil, err
	}

	cleanup := func() {
		l.Info("closing rabbitmq connection")
		conn.Close()
	}

	return &rAdapter{
		AmqpConn: conn,
		log:      log.NewHelper(logger),
	}, cleanup, nil
}

func (r *rAdapter) SendMessage2RabbitMQ(ctx context.Context, messages []biz.Message) error {

	ch, err := r.AmqpConn.Channel()
	if err != nil {
		log.Error(err, "Failed to open a channel")
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Error(err, "Failed to declare a queue")
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, message := range messages {

		body, err := json.Marshal(message)

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/json",
				Body:        body,
			})
		if err != nil {
			log.Error(err, "Failed to publish a message")
			return err
		}

		log.Infof(" [x] Sent %s\n", string(body))
	}

	return nil
}
