# Delfare

Delfare is a project designed to demonstrate the concept of microservices and how they communicate with each other. Each microservice is implemented in Golang and follows the Domain-Driven Design (DDD) pattern. They communicate asynchronously using NATS.

## Microservices Overview

### Daemon Service:
- A Golang daemon that runs continuously to generate random JSON events.
- Each event contains the following fields:
    - `criticality` (integer)
    - `timestamp` (ISO 8601 format string)
    - `eventMessage` (string)
- Publishes these events on NATS using the subject `events`. Published events are encrypted.

### Client Service:
- A Golang client that queries NATS messages to retrieve the **last 10 events with a criticality level higher than `x`**, where `x` is set via an environment variable.
- The processed events are displayed in logs.

### Reader Service:
- A Golang client that services NATS requests to query **InfluxDB** to retrieve the **last `x` events with a criticality higher than `y`**, where `x` and `y` are specified in the request.
- Responds via NATS with the requested events.

### Writer Service:
- A Golang client that subscribes to events published on NATS and writes them to **InfluxDB**.
- Responds via NATS with the requested events.

## Installation
To install the project, follow these steps:

```bash
# Clone the repository
git clone https://github.com/yakob-abada/delfare.git

# Navigate to the project directory
cd delfare

# Run the main services daemon-service and client-service
mv .env_example .env
docker-compose up -d
```

Set up db by going to `http://localhost:8086/` use bucket `event-bucket` and org `event-org` then copy the obtain token into .env files inside 
`writer-service/` and `reader-service/` reader directories.

```bash

# Run writer-service
cd writer-service/ 
mv .env_example .env
docker-compose up -d

# Run reader-service
cd reader-service/
mv .env_example .env
docker-compose up -d
```

## Areas for Improvement
- Secure NATS using username & password mechanism is not working. Needs further investigation.
- Increase test coverage.
- Enhance error handling and logging.
- Improve documentation for better usability.
- Optimize performance for high-throughput scenarios.

## Contact
For any questions or contributions, feel free to reach out via yakob.abada@gmail.com.
