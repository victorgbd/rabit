// Sender.go

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", "postgres://jheyson:laparadpieza1@64.225.48.231/rabbitdb?sslmode=disable")
	failOnError(err, "Fallo al conectar con la base de datos")
	defer db.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "fallo al registrar el consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Mensaje recibido: %s", d.Body)
			_, err := db.Exec("INSERT INTO mensajes (mensaje) VALUES ($1)", d.Body)
			failOnError(err, "Fallo al insertar el mensaje en la base de datos")
			err = ch.PublishWithContext(context.Background(),
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(d.Body),
				})
			failOnError(err, "fallo al publicar el mensaje")
			fmt.Println(" [x] enviado", d.Body)
			if err := d.Ack(false); err != nil {
				log.Fatalf("Error al eliminar el mensaje: %v", err)
			}
			log.Printf("Mensaje eliminado: %s", d.Body)
		}
	}()

	fmt.Println(" [*] Esperando mensajes nuevos. To exit press CTRL+C")
	<-forever

}
