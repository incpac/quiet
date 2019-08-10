# Quiet/AMQP

The AMQP module provides functionality for connecting to message servers providing either AMQP v0.9.1 and v1.0.

It will first connect try connecting using AMQP v1.0. If this fails, the client will try connecting with v0.9.1.

## Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/incpac/quiet/amqp"
    "github.com/incpac/quiet/config"
)

func main() {
    conf := config.ParseString("amqp://guest:guest@localhost:5672/demo")

    client, err := amqp.NewClient(conf)
    if err != nil {
        log.Fatalf("Failed to create client:", err.Error)
    }

    client.Watch(func(s string) {
        log.Printf("Message received: %s", s)
    })

    log.Print("Watching...")

    for i := 0; i < 12; i++ {
        client.Post(fmt.Sprintf("Hello World #%d", i))
        time.Sleep(time.Second * 5)
    }

    client.Close()
}
```
