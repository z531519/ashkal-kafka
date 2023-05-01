
# Kafka Java Template

This template project showcases Basic Kafka Consumers and Producer applications using Spring Boot.

**The Demo App** 
- expose an API `/customers` that can receive a Customer Record in avro-json format and send the record to the `customers` topic
- consumes the customer record and display details on the cnsole
  - add the prefix "MINOR", if the customer is less than 21yrs old
- support large payloads (1.5MB)
- use Avro message type and Schema Registry

## Prerequisites

- This project depends on the [Kafka Environment Project](../../../kafka-environment).  Please make sure you pull down this repo locally.


## Spring Boot Configuration

The application configuration is primarily driven by the `application*.yaml`.  Provided here are the [default, local, cloud] configurations.  The default points to the local kafka services provisioned by docker compose, and the cloud to the confluent cloud instance.  

No code change are necessary to switch between the different target environments.


## Running the Demo App Locally

First, make sure that you have pulled down the supporting repo [Kafka Environment Project](../../../kafka-environment).
Follow the instructions to start up your local Kafka Environment.

Once we have kafka services running, let's start our Customer Demo App with
```shell
./gradlew bootRun
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

The consumer should pick up the messages and should print ou this message in the log
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

Shown here is the default SpringBoot configuration.

```yaml
spring:
  kafka:
    producer:
      key-serializer: org.apache.kafka.common.serialization.StringSerializer
      value-serializer: io.confluent.kafka.serializers.KafkaAvroSerializer
      # the following settings used to support large payloads and compression
      compression-type: gzip 
      max-request-size: 2MB
    consumer:
      client-id: template
      group-id: template-consumer
      auto-offset-reset: earliest
      key-deserializer: org.apache.kafka.common.serialization.StringDeserializer
      value-deserializer: io.confluent.kafka.serializers.KafkaAvroDeserializer
      properties:
        specific.avro.reader: true

topic.name: customers
```

Here that we are specifying a compression of `gzip` and also setting the max request size to `2MB`.  The broker and topic configurations mentioned above are configured during provisioning.  You can see your current configuration here:
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

You can use the [mockdata](../../../kafka-environment//tree/main/mockdata) project to generate sample Customer Records with large payloads.


### Payload Too Large

If the customer record exceeds the producer `message.max.bytes`, the producer will emit this error
```text
org.apache.kafka.common.errors.RecordTooLargeException: The message is 13145926 bytes when serialized which is larger than 1048576, which is the value of the max.request.size configuration.
```

This size is based on the raw payload size prior to any compression.

## Running the Demo App with Confluent Cloud

This is done simply by switching our spring boot profile to `cloud`.  This will load the [application-cloud.yml](./src/main/resources/application-cloud.yml) as the primary configuration.

The sample below shows a simple plain sasl connection and will likely look different from your organization setup (typically secured with SSL).

TBD: We will need to swap this out with how ashkal manages its connection to confluent cloud

You will need to swap out the variables shown in this sample configuration with your cloud's settings.  Make sure not to commit the secrets in git repo.

```yaml
spring:
  kafka:
    properties:
      sasl:
        jaas:
          config: org.apache.kafka.common.security.plain.PlainLoginModule required
            username='${CLUSTER_API_KEY}' password='${CLUSTER_API_SECRET}';
        mechanism: PLAIN
      session:
        timeout:
          ms: '45000'
      security:
        protocol: SASL_SSL
      basic.auth.credentials.source: USER_INFO
      basic.auth.user.info: ${SR_API_KEY}:${SR_API_SECRET}
      schema.registry.url: ${SR_URL}
      bootstrap:
        servers: ${CLUSTER_SERVER}
```

Running the producer and consumer applications should be the similar commands when running locally but this time we will need to set the spring profile to `cloud`. For example,

```shell
SPRING_PROFILES_ACTIVE=cloud ./gradlew bootRun
```

Depending on your cluster setup, you may have to create the `customers` topic first.  In most production settings, `auto.create.topics.enable` is disabled.

## Running Tests

A sample unit and integration test is provided and can be run with the following commands

```shell
./gradlew test
./gradlew testIntegration
```

The integration test uses [testcontainers](https://www.testcontainers.org/) and requires Docker.  Make sure that you spin down the docker compose kafka environment before running the integration test to make sure there are no conflicts. 


# Docker Image and Helm

Provided in this project is a sample Dockerfile and Helm Chart.  These sample configuration should be a good starting point for building up your application images and deploying them to your Kubernetes Cluster.

Run this script to build a local image `java-kafka-demo` and deploy it to your local Kubernetes.
```shell
./helm-deploy.sh
```
The output should look like this
```text
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace ashkal -l "app=java-kafka-demo" -o jsonpath="{.items[0].metadata.name}")
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
  helm upgrade java-kafka-demo helm/java-kafka-demo --namespace ashkal --install -f ./helm/values-cloud.yaml
  ```
- Inspect the deployment and container logs to verify that the Demo App is able to connect to confluent cloud  


Please refer to the [Dockerfile](./Dockerfile) and [Helm chart](./helm/java-kafka-demo) for more details.



