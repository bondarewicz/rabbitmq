package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func init() {

}

func Get(user string, password string, host string, vhost string, queue string) {

	conn, err := amqp.Dial("amqp://%s:%s@%s:5672/%s", user, password, host, vhost)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ch.Qos(32, 0, false)

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"k6",   // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		true,   // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			log.Printf("Received a message: %s", d.Body)

			d.Ack(true)
			return d.Body
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
