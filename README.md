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

## License

The MIT License (MIT)

Copyright &copy; 2021 thewizardplusplus
