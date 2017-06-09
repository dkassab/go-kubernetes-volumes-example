package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
func main() {

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	listlen, err := c.Do("LLEN", "employeelist")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"employee", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	for i := 0; i < int(listlen.(int64)); i++ {
		record, err := redis.String(c.Do("LPOP", "employeelist"))
		if err != nil {
			log.Fatal(err)
			fmt.Println("key not found")
		}
		fmt.Println(record)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(record),
			})
		failOnError(err, "Failed to publish a message")
	}
}
