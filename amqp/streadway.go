package amqp 

import (
	"github.com/incpac/quiet/config"
	"github.com/streadway/amqp"
)

type streadway struct {
	config		config.Connection
	connection	*amqp.Connection
	channel		*amqp.Channel
	queue		amqp.Queue
}

func newStreadway(config config.Connection) (*streadway, error) {
	c := new(streadway)
	c.config = config

	conn, err := amqp.Dial(c.config.ConnectionString(true))
	if err != nil {
		return nil, err
	}
	c.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	c.channel = ch

	q, err := c.channel.QueueDeclare(
		c.config.Queue,
		false, 
		false,
		false,
		false, 
		nil,
	)
	if err != nil {
		return nil, err
	}
	c.queue = q

	return c, nil
}

func (c streadway) Post(m string) error {
	err := c.channel.Publish(
		"",
		c.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(m),
		},
	)

	return err
}

func (c streadway) watch(f func(string)) error {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			f(string(d.Body))
		}
	}()

	<- forever

	return nil
}

func (c streadway) Watch(f func(string)) {
	go c.watch(f)
}

func (c streadway) Get() (string, error) {	
	msg, _, err := c.channel.Get(
		c.queue.Name,
		true,
	)
	if err != nil {
		return "", err
	}

	return string(msg.Body), nil
}

func (c streadway) Close() {
	c.channel.Close()
	c.connection.Close()
}