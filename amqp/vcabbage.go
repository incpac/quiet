package amqp

import (
	"context"
	"time"
	"github.com/incpac/quiet/config"
	"pack.ag/amqp"
)

type vcabbage struct {
	config		config.Connection
	client		*amqp.Client 
	session		*amqp.Session
	context		context.Context
	receiver	*amqp.Receiver
}

func newVCabbage(config config.Connection) (*vcabbage, error) {
	c := new(vcabbage)
	c.config = config
	
	client, err := amqp.Dial(c.config.ConnectionString(false), amqp.ConnSASLPlain(c.config.Username, c.config.Password))
	if err != nil {
		return nil, err 
	}
	c.client = client

	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	c.session = session

	ctx := context.Background()
	c.context = ctx

	return c, nil
}

func (c vcabbage) Post(m string) error {

	sender, err := c.session.NewSender(amqp.LinkTargetAddress("/"+ c.config.Queue))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.context, 5*time.Second)

	err = sender.Send(ctx, amqp.NewMessage([]byte(m)))
	if err != nil {
		cancel()
		return err
	}

	sender.Close(ctx)
	cancel()

	return nil
}

func (c vcabbage) watch(f func(string)) error {
	receiver, err := c.session.NewReceiver(amqp.LinkSourceAddress("/" + c.config.Queue), amqp.LinkCredit(10))
	if err != nil {
		return err
	}
	c.receiver = receiver

	for {
		msg, err := receiver.Receive(c.context)
		if err != nil {
			return err
		}

		msg.Accept()

		f(string(msg.GetData()))
	}
} 

func (c vcabbage) Watch(f func(string)) {	
	go c.watch(f)
}

func (c vcabbage) Get() (string, error) {
	receiver, err := c.session.NewReceiver(amqp.LinkSourceAddress("/" + c.config.Queue), amqp.LinkCredit(10))
	if err != nil {
		return "", err
	}

	msg, err := receiver.Receive(c.context)
	if err != nil {
		return "", err
	}

	msg.Accept()
	receiver.Close(c.context)

	return string(msg.GetData()), nil
}

func (c vcabbage) Close() {
	if c.receiver != nil {
		c.receiver.Close(c.context)
	}

	c.client.Close()
}