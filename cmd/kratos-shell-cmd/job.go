package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	mysqlQuery "kratos-k8s-job/internal/data/mysql"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "kratos-k8s-job/api/helloworld/v1"
)

func runCommand(context.Context) error {
	msg := "message from K8S"
	fmt.Printf("Hello Kratos Application: %v\n\n", msg)

	//err := queryMySqlDB()
	//if err != nil {
	//	log.Println(err)
	//}
	//sendMessage2RabbitMQ(msg)

	callGRPC(msg)

	done <- true
	return nil
}

func callGRPC(msg string) {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: msg})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)

	// returns error
	//	_, err = client.SayHello(context.Background(), &v1.HelloRequest{Name: "error"})
	//	if err != nil {
	//		log.Printf("[grpc] SayHello error: %v\n", err)
	//	}
	//	if errors.IsBadRequest(err) {
	//		log.Printf("[grpc] SayHello error is invalid argument: %v\n", err)
	//	}
}

func queryMySqlDB() error {
	ctx := context.Background()

	fmt.Println("Call queryMySqlDB")

	dbCf := bc.Data.Database
	db, err := sql.Open(dbCf.GetDriver(), dbCf.GetSource())
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

	amqpCf := bc.Data.Amqp
	conn, err := amqp.Dial("amqp://" + amqpCf.GetAddr())
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
