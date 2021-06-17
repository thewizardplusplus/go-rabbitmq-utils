// +build integration

package rabbitmqutils_test

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

func ExampleClient() {
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
