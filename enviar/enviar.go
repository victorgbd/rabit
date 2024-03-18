// Sender.go

package main

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@64.225.48.231:5672/")
	failOnError(err, "error en la conexion a RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "fallo al abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mensajes", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "fallo al declarar la cola")

	body := "Hello World!"
	err = ch.PublishWithContext(context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "fallo al publicar el mensaje")
	fmt.Println(" [x] enviado", body)
}
