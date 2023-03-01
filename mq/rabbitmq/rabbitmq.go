package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// Connection configuration
type Option struct {
	// "amqp://test:test@127.0.0.1:5672/"
	Url string

	Queue    string
	Exchange string
	Key      string
}

// Serve structure
type Serve struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	queue    string
	exchange string
	key      string
}

func NewServe(option *Option) (*Serve, error) {
	conn, err := amqp.Dial(option.Url)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Serve{
		conn:    conn,
		channel: channel,

		queue:    option.Queue,
		exchange: option.Exchange,
		key:      option.Key,
	}, nil
}

func (s *Serve) Close() {
	s.conn.Close()
	s.channel.Close()
}

func (s *Serve) Publish(message string) (string, error) {
	_, err := s.channel.QueueDeclare(
		s.queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", err
	}

	err = s.channel.Publish(
		s.exchange,
		s.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (s *Serve) Receive() {
	q, err := s.channel.QueueDeclare(
		s.queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	messages, err := s.channel.Consume(
		q.Name,
		"TEST",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			// todo do something
			log.Printf("Received a message: %s", message.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
