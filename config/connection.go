package config 

import (
	"strings"
)

// Connection breaks down a connection string into it's individual components.
type Connection struct {
	Server		string
	Username 	string 
	Password	string
	Port		string 
	Protocol 	string 
	Queue 		string 
}

// Internal helper function to generate string to append before server 
// name when providing inline username and password combination
func (c Connection) getUserDetails() string {
	return c.Username + ":" + c.Password + "@"
}

// ToString generates a single line string expressing the connection to a specific queue
// Example: amqp://myawesomeserver:9000/epicqueue
// Embedding the username and password is an optional flag
func (c Connection) ToString(appendUserDetails bool) string {
	return c.ConnectionString(appendUserDetails) + "/" + c.Queue 
}

// ConnectionString generates a single line string expressing the connection to a server
// Example: amqp://myawesomeserver:9000
// Embedding the username and password is an optional flag
func (c Connection) ConnectionString(appendUserDetails bool) string {
	userdetails := ""

	if appendUserDetails {
		userdetails = c.getUserDetails()
	}

	return c.Protocol + "://" + userdetails + c.Server + ":" + c.Port
}

// ParseString will take in a connection string (eg amqp://username:password@myawesomeserver:9000/epicqueue)
// and convert it into a configuration object
func ParseString(s string) Connection {
	c := Connection{}

	split1		:= strings.Split(s, ":") 			// [amqps //user password@ipaddress portnumber/folder]
	split2		:= strings.Split(split1[len(split1)-1], "/")	// [portnumber folder1 folder2 ...]

	c.Protocol	= split1[0]
	c.Port		= split2[0]
	c.Queue		= strings.Join(split2[1:len(split2)], "/")

	if strings.Contains(s, "@") {
		split3		:= strings.Split(split1[2], "@")	// [portnumber ipaddress]

		c.Username	= split1[1][2:len(split1[1])]
		c.Password	= split3[0]
		c.Server	= split3[1]
	} else {
		c.Server	= split1[1][2:len(split1[1])]
	}

	return c
}
