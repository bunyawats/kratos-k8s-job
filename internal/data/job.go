package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rabbitmq/amqp091-go"
	"kratos-k8s-job/internal/biz"
	"kratos-k8s-job/internal/data/mysql"
	"time"
)

type (
	Message struct {
		TemPlateName string `json:"templateName"`
		Version      string `json:"version"`
	}
)

func callJob(r *greeterRepo, g *biz.Greeter) error {
	messageList, err := queryMySqlDB(r.data.MySqlDB, r.log)
	if err != nil {
		return err
	}

	err = sendMessage2RabbitMQ(messageList, r.data.AmqpConn, r.log)
	if err != nil {
		return err
	}

	return nil
}

func queryMySqlDB(db *sql.DB, log *log.Helper) ([]Message, error) {
	ctx := context.Background()

	fmt.Println("Call queryMySqlDB")

	queries := mysql.New(db)

	var lastTempalteId int64
	// get current template
	currentTemplate, err := queries.GetCurrentTemplate(ctx)
	if err != nil {
		log.Warnf(err.Error())
		lastTempalteId = -1
	} else {
		lastTempalteId = currentTemplate.ConsentTemplateID
	}
	log.Infof("Found Last Tempalte ID: %v", lastTempalteId)

	lastUpdatedTemplateList, err := queries.ListAllLastUpdatedTemplate(ctx, lastTempalteId)
	if err != nil {
		log.Warnf(err.Error())
	}

	messageList := make([]Message, 0)

	for _, consentTemplate := range lastUpdatedTemplateList {
		result, err := queries.CreateCurrentTemplate(ctx, mysql.CreateCurrentTemplateParams{
			ConsentTemplateID: consentTemplate.ID,
			TemplateName:      consentTemplate.TemplateName,
			Version:           consentTemplate.Version,
		})
		if err != nil {
			return nil, err
		}
		insertedCurrentTemplateID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		log.Info("insertedCurrentTemplateID", insertedCurrentTemplateID)

		messageList = append(
			messageList,
			Message{
				TemPlateName: consentTemplate.TemplateName,
				Version:      consentTemplate.Version,
			},
		)
	}
	return messageList, nil
}

func sendMessage2RabbitMQ(messageList []Message, amqpConn *amqp091.Connection, log *log.Helper) error {

	ch, err := amqpConn.Channel()
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

	for _, message := range messageList {

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
