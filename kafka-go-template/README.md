

# Kafka Golang Template

This template project showcases Basic Kafka Consumers and Producer applications with Golang.

The following go modules are used:
- [confluent-kafka go module](https://github.com/confluentinc/confluent-kafka-go)
- [gogen-avro](https://github.com/actgardner/gogen-avro)

**The Demo App** 
- expose an API `/customers` that can receive a Customer Record in avro-json format and send the record to the `customers` topic
- consumes the customer record and display details on the cnsole
  - add the prefix "MINOR", if the customer is less than 21yrs old
- support large payloads (1.5MB)
- use Avro message type and Schema Registry

## Prerequisites

- This project depends on the [Kafka Environment Project](../../../kafka-environment).  Please make sure you pull down this repo locally.


## Running the Demo App Locally

Once you have your kafka services running locally, you can run the Demo App with
```
go run ./cmd/demo
```

**Producing Messages**

The Demo App provides an API endpoint that we can use to produce messages.  Here's a sample curl command


```shell
curl --location --request POST 'http://localhost:8080/customers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "birthdt": "1900-12-30",
    "fname": "Samuel",
    "fullname": "Samuel Francine Leffler",
    "gender": "N",
    "id": "asd",
    "lname": "Leffler",
    "mname": "Francine",
    "suffix": "V",
    "title": "Global Applications Designer",
    "largePayload": "Samuel Francine Leffler"
}'
```

The API should receive the message and send it to the `customers` Kafka Topic.

The consumer should pick up the messages and print out this message in the log
```text
Customer: Samuel Francine Leffler, birthdt: 1990-12-30
```

Play with changing the `birthdt` field to `2020-12-30`.  The output should show up like this
```text
MINOR> Customer: Samuel Francine Leffler, birthdt: 2020-12-30
```

You now have successfully produced and consumed messages in Kafka!

## Large Payloads

Kafka has an extensive configuration that can be customized to your specific needs.  Our Demo App is able to support payload sizes upwards of 1.5MB.  To do that, we have to adjust below properties from Kafka's default settings. 

Producer configuration
- [max.request.size](https://kafka.apache.org/documentation/#producerconfigs_max.request.size)
- [compression.type](https://kafka.apache.org/documentation/#brokerconfigs_compression.type)

Broker configuration
- [message.max.bytes](https://kafka.apache.org/documentation/#brokerconfigs_message.max.bytes)

Topic configuration
- max.message.bytes (https://kafka.apache.org/documentation/#topicconfigs_max.message.bytes)


The broker and topic configurations mentioned above are configured during provisioning.  You can see your current configuration here:
- broker `message.max.bytes`
  ```shell
  kafka-configs --bootstrap-server localhost:19092 \
    --describe \
    --broker 1 -all | grep message.max.bytes
  ```
- customers topic `max.message.bytes`
    ```shell
    kafka-configs --bootstrap-server localhost:19092 \
        --describe \
        --topic test \
        --all \
        | grep max.message.bytes
    ```

Below is the configuration file that our Demo App uses [config.yaml](./config.yaml). Note that you can override these configuration with environment variables.
```yaml
## Env ASHKAL_BOOTSTRAP_SERVERS
bootstrap.servers: localhost:19092
## Env ASHKAL_CLIENT_ID
client.id: ashkal
# add broker properties here, environment variable name follows a similar pattern prefixed with ASHKAL_KAFKA,
# for example: ASHKAL_KAFKA_SASL_USERNAME for property sasl.username
properties:  
  # security.protocol: SASL_SSL
  # sasl.mechanisms: PLAIN
  # sasl.username:
  # sasl.password:

## Env ASHKAL_SCHEMA_REGISTRY_URL
schema.registry.url: http://localhost:8081
## Env ASHKAL_SCHEMA_REGISTRY_USER
# schema.registry.user: 
## Env ASHKAL_SCHEMA_REGISTRY_PASSWORD
# schema.registry.password: 


# add consumer properties here, environment variable substition follows a similar pattern, with prefix of
#  ASHKAL_CONSUMER_property
consumer:
  group.id: template-go-consumer
  auto.offset.reset: earliest

# add producer properties here, environment variable substition follows a similar pattern, with prefix of
#  ASHKAL_PRODUCER_property
producer:
  compression.type: gzip
  message.max.bytes: 2097152

# topic used by the demo app
topic.name: customers
```

### Payload Too Large

If the customer record exceeds the producer `message.max.bytes`, the producer will emit this error
```text
Produce failed: Broker: Message size too large
```

This size is based on the raw payload size prior to any compression.

## Producing Customer Records using Producer App

This project comes with a command line too for producing Customer Records via command line.   Make sure you generate the `customers.json` file using the [mockdata](../../../kafka-environment/mockdata) project

```shell
# defaults to reading customers.json
go run ./cmd/producer
```

Below are the accepted arguments
```
  -config string
        target config file (default "config")
  -f string
        the customer file (default "customers.json")
```



### Generated Code

This project contains auto-generated code under `pkg/customer` folder.  The auto-generated code consists of the following:
- avro ser/deser from the customer.asvc schema
- mock implementations for sink.go

You should run `go generate ./...` whenever changes are made to the source code

### Running the Demo App with Confluent Cloud

The provided [config.yaml](./config.yaml) works with the local kafka services provisioned by docker compose.  In order to run with Confluent Cloud, we simply need to provide the proper connection configuration like the example below. 

```yaml
bootstrap.servers: {{ BOOTSTRAP_SERVER }}
client.id: ashkal
properties:  
  security.protocol: SASL_SSL
  sasl.mechanisms: PLAIN
  sasl.username: {{ CLUSTER_API_KEY }}
  sasl.password: {{ CLUSTER_API_SECRET }}


schema.registry.url: {{ SR_URL }}
schema.registry.user: {{ SR_KEY }}
schema.registry.password: {{ SR_SECRET}}

consumer:
  group.id: template-consumer
  auto.offset.reset: earliest
  session.timeout.ms: 45000 # recommended setting
producer:

topic.name: customers
```

As mentioned before, this configuration can be overriden/supplemented with environment variables.  In fact, our helm chart uses the default [config.yaml](./config.yaml) and uses environment variables to inject all the sensitive connection parameters.


## Running Tests

The project has both unit and integration test.  The integration test uses [testcontainers](https://golang.testcontainers.org/) so make sure you have docker running when running integration tests

```shell
go test ./...  # runs unit tests
go test ./... -tags=integration # runs integration tests
```



# Docker Image and Helm

Provided in this project is a sample Dockerfile and Helm Chart.  These sample configuration should be a good starting point for building up your application images and deploying them to your Kubernetes Cluster.

Run this script to build a local image `go-kafka-demo` and deploy it to your local Kubernetes.
```shell
./helm-deploy.sh
```

The output should look like this
```text
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace ashkal -l "app=go-kafka-demo" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace ashkal $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace ashkal port-forward $POD_NAME 8080:$CONTAINER_PORT
```

You can inspect your deployment with:
```shell
kubectl get deployments --namespace ashkal
```

The installed demo app connects to your local Kafka services.  If you look at [helm-deploy.sh](./helm-deploy.sh), it uses the [local values file](./helm/values-local.yaml) during the helm chart install.

## Helm and Confluent Cloud configuration

You can also install the demo app to connect to confluent cloud.  Our helm chart relies on Kubernetes Secrets as the source for our connection credentials so we will need create this secret prior to installing our Helm chart.

Follow these steps:
- create a kubernetes secret containing the credentials used for connecting to confluent cloud.  
  - Replace the values in the provided [sample-secrets.yaml](./helm/sample-secrets.yaml).
    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: cloud-kafka
      namespace: ashkal
    type: Opaque
    data:
      CLUSTER_API_KEY: your base64 encoded secret here
      CLUSTER_API_SECRET: your base64 encoded secret here
      SR_API_KEY: your base64 encoded secret here
      SR_API_SECRET: your base64 encoded secret here 
    ```
  - apply this resource with `kubectl apply -f helm/sample-secrets.yaml`
- Replace the values in the provided [values-cloud.yaml](./helm/values-cloud.yaml)

  ```yaml
  kafkaSecretName: cloud-kafka
  data:
    spring.profiles.active: cloud
    kafka.bootstrap_servers: # your cloud bootstrap server here
    kafka.schema_registry_servers: # your cloud schema registry here
  ```
- Install the helm chart 
  ```shell
  helm upgrade go-kafka-demo helm/go-kafka-demo --namespace ashkal --install -f ./helm/values-cloud.yaml
  ```
- Inspect the deployment and container logs to verify that the Demo App is able to connect to confluent cloud  



Refer to the [Dockerfile](./Dockerfile) and [Helm chart](./helm/go-kafka-demo) provided for more details.

## Using SSL

The confluent-kafka-go is a lightweight wrapper around librdkafka.  In order to use SSL, please follow the configuration as mentioned in the official librdkafka documentation.
https://github.com/confluentinc/librdkafka/wiki/Using-SSL-with-librdkafka#configure-librdkafka-client

Note that `kcat` uses rdkafka as well, and is configured in similar fashion.  You can use this utility tool to verify that you are able to connect to your broker.  For example,

```shell
# point to your secrets files
kcat -L -b yourbroker:9092 \
  -X security.protocol=ssl \
  -X ssl.ca.location=secrets/ca-cert \
  -X ssl.certificate.location=secrets/client_ashkal_client.pem \
  -X ssl.key.location=secrets/client_ashkal_client.key \
  -X ssl.key.password=*redacted*
  -X debug='security'
```