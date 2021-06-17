# go-rabbitmq-utils

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-rabbitmq-utils?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-rabbitmq-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-rabbitmq-utils)](https://goreportcard.com/report/github.com/thewizardplusplus/go-rabbitmq-utils)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-rabbitmq-utils.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-rabbitmq-utils)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-rabbitmq-utils/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-rabbitmq-utils)

The library that provides utility entities for working with [RabbitMQ](https://www.rabbitmq.com/).

## Features

- client:
  - options (will be applied on connecting):
    - maximal queue size (optionally);
    - queues for declaring:
      - will survive server restarts and remain without consumers;
    - message ID generator for them automatic generating;
  - operations:
    - with a connection:
      - opening;
      - closing;
    - with messages:
      - message publishing:
        - check the specified queue name based on the declared queues;
        - automatic marshalling of a message data to JSON;
        - setting of auxiliary message fields:
          - setting of a message ID:
            - receiving of a custom message ID (optionally);
            - automatic generating of a message ID (optionally);
          - setting of a message timestamp;
      - starting of message consuming:
        - check the specified queue name based on the declared queues;
        - automatic generating of a consumer name;
      - cancelling of message consuming:
        - check the specified queue name based on the declared queues;
        - automatic generating of a consumer name;
- message consumer:
  - arguments:
    - client;
    - queue name;
    - outer message handler;
  - operations:
    - message consuming:
      - starting;
      - cancelling;
    - message handling:
      - support of concurrent handling;
- wrappers for an outer message handler:
  - acknowledger:
    - processing on success:
      - acknowledging of the message;
      - logging of the success fact;
    - processing on failure:
      - rejecting of the message:
        - with once message handling;
        - with twice message handling (i.e. once requeue);
      - logging of the error;
  - JSON message handler:
    - automatical creating of a receiver for a message data by its specified type;
    - unmarshalling of a message data from JSON to the created receiver.

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/go-rabbitmq-utils.git
$ cd go-rabbitmq-utils
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

## Examples

```go
package main

import (
	"fmt"
	stdlog "log"
	"os"
	"reflect"
	"runtime"
	"sync"

	"github.com/go-log/log/print"
	rabbitmqutils "github.com/thewizardplusplus/go-rabbitmq-utils"
)

type exampleMessage struct {
	FieldOne int
	FieldTwo string
}

type messageHandler struct {
	waitGroupInstance *sync.WaitGroup
	locker            sync.Mutex
	messages          []exampleMessage
}

func (messageHandler *messageHandler) MessageType() reflect.Type {
	return reflect.TypeOf(exampleMessage{})
}

func (messageHandler *messageHandler) HandleMessage(message interface{}) error {
	defer messageHandler.waitGroupInstance.Done()

	messageHandler.locker.Lock()
	defer messageHandler.locker.Unlock()

	messageHandler.messages =
		append(messageHandler.messages, message.(exampleMessage))
	return nil
}

func main() {
	dsn, ok := os.LookupEnv("MESSAGE_BROKER_ADDRESS")
	if !ok {
		dsn = "amqp://rabbitmq:rabbitmq@localhost:5672"
	}

	// prepare the client
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	client, err :=
		rabbitmqutils.NewClient(dsn, rabbitmqutils.WithQueues([]string{"example"}))
	if err != nil {
		logger.Fatal(err)
	}
	defer client.Close()

	// start the message consuming
	var waitGroupInstance sync.WaitGroup
	messageHandler := messageHandler{waitGroupInstance: &waitGroupInstance}
	messageConsumer, err := rabbitmqutils.NewMessageConsumer(
		client,
		"example",
		rabbitmqutils.Acknowledger{
			MessageHandling: rabbitmqutils.OnceMessageHandling,
			MessageHandler: rabbitmqutils.JSONMessageHandler{
				MessageHandler: &messageHandler,
			},
			// wrap the standard logger via the github.com/go-log/log package
			Logger: print.New(logger),
		},
	)
	go messageConsumer.StartConcurrently(runtime.NumCPU())
	defer messageConsumer.Stop()
	if err != nil {
		logger.Fatal(err)
	}

	// publish the messages
	for i := 0; i < 10; i++ {
		waitGroupInstance.Add(1)

		err = client.PublishMessage("example", "", exampleMessage{
			FieldOne: 10 + i,
			FieldTwo: fmt.Sprintf("message data #%d", i),
		})
		if err != nil {
			logger.Fatal(err)
		}
	}
	waitGroupInstance.Wait()

	// print the results
	for _, message := range messageHandler.messages {
		fmt.Printf("%+v\n", message)
	}

	// Unordered output:
	// {FieldOne:10 FieldTwo:message data #0}
	// {FieldOne:11 FieldTwo:message data #1}
	// {FieldOne:12 FieldTwo:message data #2}
	// {FieldOne:13 FieldTwo:message data #3}
	// {FieldOne:14 FieldTwo:message data #4}
	// {FieldOne:15 FieldTwo:message data #5}
	// {FieldOne:16 FieldTwo:message data #6}
	// {FieldOne:17 FieldTwo:message data #7}
	// {FieldOne:18 FieldTwo:message data #8}
	// {FieldOne:19 FieldTwo:message data #9}
}
```

## License

The MIT License (MIT)

Copyright &copy; 2021 thewizardplusplus
