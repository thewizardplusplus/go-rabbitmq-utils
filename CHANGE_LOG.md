# Change Log

## [v1.1.3](https://github.com/thewizardplusplus/go-rabbitmq-utils/tree/v1.1.3) (2021-08-04)

Adding of the getting of a single message to the [RabbitMQ](https://www.rabbitmq.com/) client.

- client:
  - operations:
    - with messages:
      - getting of a single message:
        - check the specified queue name based on the declared queues;
        - block the execution flow until the message is received or an error occurs.

## [v1.1.2](https://github.com/thewizardplusplus/go-rabbitmq-utils/tree/v1.1.2) (2021-07-25)

Fixing of the bugs.

- fix the bugs:
  - in the starting of the message consuming;
  - in the example.

## [v1.1.1](https://github.com/thewizardplusplus/go-rabbitmq-utils/tree/v1.1.1) (2021-06-17)

Adding of the integration tests and example.

- adding of the integration tests:
  - for the message publishing via the client;
  - for the message consuming via the client:
    - including the message consumer using;
- adding of the example:
  - with the using:
    - of the message consumer;
    - of the wrappers for an outer message handler:
      - acknowledger;
      - JSON message handler.

## [v1.1](https://github.com/thewizardplusplus/go-rabbitmq-utils/tree/v1.1) (2021-06-10)

Checking of the queue name used in the [RabbitMQ](https://www.rabbitmq.com/) client operations based on the declared queues and extending of auxiliary message fields with a message ID and message timestamp.

- client:
  - options (will be applied on connecting):
    - message ID generator for them automatic generating;
  - operations:
    - with messages:
      - message publishing:
        - check the specified queue name based on the declared queues;
        - setting of auxiliary message fields:
          - setting of a message ID:
            - receiving of a custom message ID (optionally);
            - automatic generating of a message ID (optionally);
          - setting of a message timestamp;
      - starting of message consuming:
        - check the specified queue name based on the declared queues;
      - cancelling of message consuming:
        - check the specified queue name based on the declared queues.

## [v1.0](https://github.com/thewizardplusplus/go-rabbitmq-utils/tree/v1.0) (2021-05-28)

Major version. Implementing of the [RabbitMQ](https://www.rabbitmq.com/) client and utilities for message consuming.
