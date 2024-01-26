package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rabbitmq/amqp091-go"
	"kratos-k8s-job/internal/biz"
	"kratos-k8s-job/internal/data/mysql"
	"reflect"
	"time"
)

func callJob(r *greeterRepo, g *biz.Greeter) error {
	err := queryMySqlDB(r.data.MySqlDB, r.log)
	if err != nil {
		return err
	}
	err = sendMessage2RabbitMQ(g.Hello, r.data.AmqpConn, r.log)
	if err != nil {
		return err
	}
	return nil
}

func queryMySqlDB(db *sql.DB, log *log.Helper) error {
	ctx := context.Background()

	fmt.Println("Call queryMySqlDB")

	queries := mysql.New(db)

	// get current template
	currentTemplate, err := queries.GetCurrentTemplate(ctx, 1)
	if err != nil {
		return err
	}
	log.Info(currentTemplate)

	// create new current template
	result, err := queries.CreateCurrentTemplate(ctx, mysql.CreateCurrentTemplateParams{
		TemplateName: "Data Privacy",
		Version:      "1.0",
	})
	if err != nil {
		return err
	}

	insertedCurrentTemplateID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Info("insertedCurrentTemplateID", insertedCurrentTemplateID)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetCurrentTemplate(ctx, insertedCurrentTemplateID)
	if err != nil {
		return err
	}

	// prints true
	log.Info(reflect.DeepEqual(insertedCurrentTemplateID, fetchedAuthor.ID))
	return nil
}

func sendMessage2RabbitMQ(msg string, conn *amqp091.Connection, log *log.Helper) error {

	ch, err := conn.Channel()
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

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		log.Error(err, "Failed to publish a message")
		return err
	}

	log.Infof(" [x] Sent %s\n", msg)
	return nil
}
