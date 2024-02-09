package data

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rabbitmq/amqp091-go"
	"kratos-k8s-job/internal/biz"
	"time"
)

func (r *jobRepo) SendMessage2RabbitMQ(ctx context.Context, messages []biz.Message) error {

	ch, err := r.data.AmqpConn.Channel()
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
			amqp091.Publishing{
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
