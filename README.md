# Integration and Unit Testing in Golang

Test is in producer folder

## Setup

### Install Fiber

> Fiber is an Express inspired web framework written in Go

```bash
go get -u github.com/gofiber/fiber/v2
```

### Install Testify

> Testify is a toolkit with common assertions and mocks that plays nicely with the standard library.

```bash
go get -u github.com/stretchr/testify/mock
```

### Install Mockery

> Mockery is a mock code autogenerator for Golang.

```bash
go get github.com/vektra/mockery/v2/.../
```

### Go coverage highlighting (File should place beside of testing file)

> Go coverage highlighting is a plugin for highlighting lines in the Go source code that are not covered by tests.

place in json settings file

```json
"go.coverageOptions": "showBothCoveredAndUncoveredCode",
  "go.coverageDecorator": {
    "type": "highlight",
    "coveredHighlightColor": "rgba(59, 255, 125, 0.03)",
    "uncoveredHighlightColor": "rgba(255, 100, 81, 0.03)",
    "coveredBorderColor": "rgba(59, 255, 125, 0.03)",
    "uncoveredBorderColor": "rgba(255, 100, 81, 0.03)",
    "coveredGutterStyle": "blockgreen",
    "uncoveredGutterStyle": "blockred"
  },
  "editor.unicodeHighlight.nonBasicASCII": false,
```

## Mockery

> **--dir**: The directory in which to search for interfaces. (required)<br> **--name**: The name of the interface for which to generate a mock. (required)<br> **--output**: The directory to which the mock file will be written. (required)<br> **--outpkg**: The name of the package in which the mock should be generated. (required)<br> **--filename**: The name of the file to which the mock will be written. (optional)

> <ins>NOTE</ins> (backslash) at the end of each line indicates that the command continues on the next line.

<ins>mock with a specified, clear path</ins>

- **when to use**: `When you want to mock a function that is not in the same package.`

> #### services
>
> ```bash
> mockery \
> --dir=services/cacheService \
> --name=ICacheService \
> --output=services/mock \
> --outpkg=mockService \
> --filename=ICacheService.go
> ```

> #### repositories
>
> ```bash
> mockery \
> --dir=repositories/privilegeRepository \
> --name=IPrivilegeRepository \
> --output=repositories/mock \
> --outpkg=mockRepo \
> --filename=IPrivilegeRepository.go
> ```

<ins>mock all interfaces in a package</ins>

> #### services
>
> ```bash
> mockery \
> --dir=services \
> --output=services/mock \
> --outpkg=mockService \
> --all
> ```

> #### repositories
>
> ```bash
> mockery \
> --dir=repositories \
> --output=repositories/mock \
> --outpkg=mockRepo \
> --all
> ```

## How to run test

### Unit Test

> Unit tests are tests that test a single unit of code.

```bash
go test ./...
```

### Integration Test

#### 1. Not mark integration test on top of package in test file

```
//go:build integration
```

> Integration tests are tests that test multiple units of code together.

#### 2. Mark integration test on top of package

> Integration tests are skipped by default. To run integration tests, use the -tags=integration flag.

```
//go:build integration
```

```bash
go test ./... -tags=integration
```

## This is a simple example of using Kafka with Golang and MariaDB.

### Prerequisites:

- For Mac users: Install Kafka using Homebrew. if not you should install kafka for run 2. step
  ```bash
  brew install kafka
  ```
- Install Docker Desktop for Mac or Windows.
- Install Docker Compose.
- Docker pull kafka image
  ```bash
  docker pull bitnami/kafka
  ```
- Docker pull zookeeper image
  ```bash
  docker pull zookeeper
  ```
- Docker pull mariadb image
  ```bash
  docker pull mariadb
  ```

### Steps:

#### 1. Clean Up Data Directories:

> On the first clone, in the Server folder if there's existing data in the Kafka and ZooKeeper folders, remove it.

```
rm -rf ./kafka/*
rm -rf ./zookeeper/*
```

#### 2. Start Zookeeper & Kafka with Docker Compose:

> Run the following command to start ZooKeeper and Kafka servers using Docker Compose.

```
docker-compose up
```

#### 3. Run the Consumer (in consumer folder):

> Open a new terminal window and run the following command to start the consumer.

```
go run main.go
```

#### 4. Run the Producer (in producer folder):

> Open a new terminal window and run the following command to start the producer.

```
go run main.go
```

#### 5. Test the Application:

> Data will be sent to the Kafka server and saved to the MariaDB database.

> POST http://localhost:8000/openAccount

```json
{
  "AccountHolder": "kafkaman",
  "AccountType": 1,
  "OpeningBalance": 20000
}
```

> POST http://localhost:8000/depositFund

```json
{
  "ID": "cfbd34d7-fb3e-42db-b66a-ae9e55b16aec",
  "Amount": 3000
}
```

> POST http://localhost:8000/withdrawFund

```json
{
  "ID": "cfbd34d7-fb3e-42db-b66a-ae9e55b16aec",
  "Amount": 3000
}
```

> POST http://localhost:8000/closeAccount

```json
{
  "ID": "105fc312-af37-4057-bc11-325839c3ccee"
}
```

### Explanation:

> The producer will send a message to the Kafka server. The consumer will receive the message from the Kafka server and save it to the MariaDB database.

### Explain the Code:

#### 1. Producer:

> The producer will send a message to the Kafka server.

#### 2. Consumer:

> The consumer will receive the message from the Kafka server and save it to the MariaDB database.

#### 3. Docker Compose:

> Docker Compose is a tool for defining and running multi-container Docker applications. With Compose, you use a YAML file to configure your application's services. Then, with a single command, you create and start all the services from your configuration.

#### 4. Docker Compose File:

> The Docker Compose file is a YAML file defining services, networks, and volumes. The Docker Compose file is used to configure the application's services. The Docker Compose file is used to configure the application's networks. The Docker Compose file is used to configure the application's volumes.

### What is Kafka ZooKeeper, Kafka Cluster, and Kafka Broker?

#### 1. Kafka ZooKeeper:

> ZooKeeper is a centralized service for maintaining configuration information, naming, providing distributed synchronization, and providing group services. All of these kinds of services are used in some form or another by distributed applications.

#### 2. Kafka Cluster:

> A Kafka cluster is a group of brokers that work together to process and store records. Brokers are stateless, so they use ZooKeeper for maintaining their cluster state. Kafka uses ZooKeeper to store offsets of messages consumed for a specific topic and partition by a specific consumer group. Kafka uses ZooKeeper to perform leader election for partitions. ZooKeeper sends changes of the topology to Kafka, so each node in the cluster knows when a new broker has come online or gone offline.

#### 3. Kafka Broker:

> A Kafka broker is a server that runs in a Kafka cluster. Kafka brokers form a cluster. Each broker is identified by its ID (integer). Each broker contains certain topic partitions. After connecting to any broker, you will be connected to the entire cluster. A Kafka cluster can contain multiple brokers. Kafka brokers are stateless, so they use ZooKeeper for maintaining their cluster state.

### Note:

- "partitions" is the number of partitions in a topic. The default value is 1.

- "group" is the consumer group name. If the Kafka broker goes down, messages in the group will be sent to consumers once the broker is back up.

- "topic" is the topic name.

- "bootstrap-server" is the Kafka server address.

- "broker-list" is the Kafka server address.

- "create" is used to create a topic.

- "list" is used to list all topics.

- "rm -rf" is used to remove a folder.

- "rm" is used to remove a file.

- "up" is used to start Docker Compose.

- "down" is used to stop Docker Compose.

- "kafka-topics" is a command for creating and listing topics.

- "kafka-console-consumer" is a command for creating consumers.

- "kafka-console-producer" is a command for creating producers.
