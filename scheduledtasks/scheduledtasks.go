package main

import "github.com/streadway/amqp"

func main() {
	url := "amqp://guest:guest@localhost:5672"
	connection, err := amqp.Dial(url)

	if err != nil {
		panic("connection error rabbitmq: " + err.Error())
	}

	channel, err := connection.Channel()

	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

}
