package main

import (
	"encoding/json"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func replacePlaceholders(message string, data map[string]interface{}) string {
	for key, value := range data {
		placeholder := "<" + key + ">"
		message = strings.Replace(message, placeholder, value.(string), -1)
	}
	return message
}

func sendMessage() {
	remoteHost := "amqp://Gleam:gleamadmin@93.188.162.243:5672/" // Replace with the actual remote host
	conn, err := amqp.Dial(remoteHost)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queueName := "EmailQueue" // Use the same queue name as the routing key
	routingKey := queueName

	bossInfo := map[string]interface{}{
		"email": "inasy@lkjsdfk.com",
		"name":  "Inayath",
		"otp":   "9029",
	}

	messageTemplate := "Hi <name>, Your OTP is <otp>. Welcome to our website <email>"
	message := replacePlaceholders(messageTemplate, bossInfo)

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	// Declare the queue
	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Publish the message to the default exchange with the routing key (queue name)
	err = ch.Publish("", routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        messageJSON,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sent: %s", messageJSON)
}

func main() {
	sendMessage()
}
