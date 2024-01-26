package main

import (
	"context"
	"database/sql"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	mysqlQuery "kratos-k8s-job/internal/data/mysql"
)

func runCommand(context.Context) error {
	msg := "message from K8S"
	fmt.Printf("Hello Kratos Application: %v\n\n", msg)

	err := queryMySqlDB()
	if err != nil {
		log.Println(err)
	}
	sendMessage2RabbitMQ(msg)

	done <- true
	return nil
}

func queryMySqlDB() error {
	ctx := context.Background()

	fmt.Println("Call queryMySqlDB")

	db, err := sql.Open("mysql", "test:test@/test?parseTime=true")
	if err != nil {
		fmt.Println("Connect to database error", err)
		return err
	}

	queries := mysqlQuery.New(db)

	// get current template
	currentTemplate, err := queries.GetCurrentTemplate(ctx, 1)
	if err != nil {
		return err
	}
	log.Println(currentTemplate)

	// create new current template
	result, err := queries.CreateCurrentTemplate(ctx, mysqlQuery.CreateCurrentTemplateParams{
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
	log.Println("insertedCurrentTemplateID", insertedCurrentTemplateID)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetCurrentTemplate(ctx, insertedCurrentTemplateID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedCurrentTemplateID, fetchedAuthor.ID))
	return nil
}

func sendMessage2RabbitMQ(msg string) {
	conn, err := amqp.Dial("amqp://user:smd95nzXiN30SAXt@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msg)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
