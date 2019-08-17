package quiet 

import (
	"fmt"
	"strings"
	"github.com/incpac/quiet/amqp"
	"github.com/incpac/quiet/config"
)

// Client provides an interface for sending and receiving messages from MQ servers
type Client interface {
	Post(string) error
	Watch(func(string))
	Get() (string, error)
	Close()
}

type protocolError struct {
	protocol 	string 
}

func (e *protocolError) Error() string {
	return fmt.Sprintf("unsupported protocol - %s", e.protocol)
}

// NewClient creates an new MQ client by connecting to the MQ server specified in the provided configuration object
func NewClient(c config.Connection) (Client, error){
	switch (strings.ToLower(c.Protocol)) {
	case "amqp":
		return amqp.NewClient(c)
	case "amqps":
		return amqp.NewClient(c)
	default:
		return nil, &protocolError{c.Protocol}
	}
}

// NewClientFromString creates a new Client from a connection string such as "amqp://myawesomeserver:9000/epicqueue"
func NewClientFromString(c string) (Client, error) {
	conn := config.ParseString(c)

	return NewClient(conn)
}

