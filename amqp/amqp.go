package amqp

import (
	"github.com/incpac/quiet/config"
	"pack.ag/amqp"
)

// Client provides and interface for connecting to either AMQP v0.9.1 or AMQP v1.0
type Client interface {
	Post(string) error
	Watch(func(string))
	Close()
}

// NewClient creates a new AMQP client
func NewClient(c config.Connection) (Client, error) {
	// From research there appears to be two main versions of the AMQP protocol in the wild,
	// v0.9.1 provided by RabbitMQ and others, and v1.0 provided by Apache ActiveMQ
	// We're using two different libraries depending on the version, so we first try to connect 
	// to using version 1.0. This is provided by vcabbage's library.
	_, err := amqp.Dial(c.ConnectionString(false), amqp.ConnSASLPlain(c.Username, c.Password))

	if err != nil {
		if err.Error() == "unexpected protocol version 0.9.1" {
			// v 1.0
			return newStreadway(c)
		} else {
			return nil, err
		}
	} 

	//v 0.9.1
	return newVCabbage(c)
}