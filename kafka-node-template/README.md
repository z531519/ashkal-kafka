

# Kafka NodeJS Template

This template project showcases Basic Kafka Consumers and Producer applications with NodeJS.

The following npm modules are used:
- [node-rdkafka](https://github.com/Blizzard/node-rdkafka)
- [Schema Registry](https://github.com/kafkajs/confluent-schema-registry#readme)

This project uses Typescript.

## How about KafkaJS

[KafkaJS](https://github.com/tulios/kafkajs) is 100% javascript Kafka library and is gaining popularity.  However it appears that it is still lacking support in several features.  In our case, we need to be able to support [large payloads](#large-payloads) and as of current, KafkaJS does not support these configurations.  


**The Demo App** 
- expose an API `/customers` that can receive a Customer Record in avro-json format and send the record to the `customers` topic
- consumes the customer record and display details on the cnsole
  - add the prefix "MINOR", if the customer is less than 21yrs old
- support large payloads (1.5MB)
- use Avro message type and Schema Registry

## Prerequisites

- This project depends on the [Kafka Environment Project](../kafka-environment).  Please make sure you pull down this repo locally.


## Running the Demo App Locally

Once you have your kafka services running locally, you can run the Demo App with
```
npm run start
```


**Producing Messages**

The Demo App provides an API endpoint we will use to produce messages.  Here's a sample curl command


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

The consumer will pick up the messages and prints out this message to the log
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

Below is the configuration file that our Demo App uses [config/default.ts](./config/default.ts). 
```javascript
/ all available configurations are documented here https://github.com/confluentinc/librdkafka/blob/v2.0.2/CONFIGURATION.md

/**
 * Loads environment variables based on the group name prefix
 */
const loadEnvironmentVars = (group):any =>{
  const env = process.env;
  const map = {};
  let re = /_/gi;
  Object.keys(env).forEach(function(key) {
    let normalizedKey;
    if (key.indexOf(group) != -1) {
      if (group === key) { // exact match of key name         
        normalizedKey = key.substring(key.lastIndexOf("_") + 1).replace(re, ".").toLowerCase();
      } else {
        normalizedKey = key.replace(group + "_", "").replace(re, ".").toLowerCase();
      }
      
      map[normalizedKey] = env[key];
    }
  });
  return map;
}

const config = {
  kafka: {
    CustomerTopic: "customers",
    ProducerGlobalConfig: {      
      "client.id": "nodejs-demo",
      "metadata.broker.list": "localhost:19092,localhost:29092,localhost:39092",
      "message.max.bytes": 2097152, //2MB      
      ...loadEnvironmentVars("ASHKAL_GLOBAL"),
      ...loadEnvironmentVars("ASHKAL_PRODUCER_GLOBAL"),
    } as Kafka.ProducerGlobalConfig,
    ProducerTopicConfig: {
      "compression.type": "gzip",
      ...loadEnvironmentVars("ASHKAL_PRODUCER_TOPIC"),
    } as Kafka.ProducerTopicConfig,
    ConsumerGlobalConfig: {
      "metadata.broker.list": "localhost:19092,localhost:29092,localhost:39092",
      "group.id": "template-nodejs-consumer",
      ...loadEnvironmentVars("ASHKAL_GLOBAL"),
      ...loadEnvironmentVars("ASHKAL_CONSUMER_GLOBAL"),
    } as Kafka.ConsumerGlobalConfig,
    ConsumerTopicConfig: {
      "auto.offset.reset": "earliest",
      ...loadEnvironmentVars("ASHKAL_CONSUMER_TOPIC"),
    } as Kafka.ConsumerTopicConfig
  },
  schemaRegistry: {
    host: "http://localhost:8081",
    ...loadEnvironmentVars("ASHKAL_SCHEMA_REGISTRY_HOST"),
    auth: loadEnvironmentVars("ASHKAL_SCHEMA_REGISTRY_AUTH"),
  } as SchemaRegistryAPIClientArgs
};

```

As you can tell, the configuration can be overriden/supplemented with environment variables.  It simply follows a basic environment variable naming convention.  For example, setting the bootstrap server should look like
```
ASHKAL_GLOBAL_METADATA_BROKER_LIST=your-broker-server
```
Where the property is prefixed with `ASHKAL_GLOBAL_METADATA` followed by uppercased property name and dots replaced with underscores.


### Payload Too Large

If the customer record exceeds the producer `message.max.bytes`, the producer will emit this error
```text
PayloadTooLargeError: request entity too large
```

This size is based on the raw payload size prior to any compression.


## Producing Customer Records using Producer App Command Line Utility

This project comes with a command line too for producing Customer Records via command line.   Make sure you generate the `customers.json` file using the [mockdata](../kafka-environment/mockdata) project

```shell
# defaults to reading customers.json
npm run start:producer
```


### Running the Demo App with Confluent Cloud

The default configuration works with the local kafka services provisioned by docker compose.  In order to run with Confluent Cloud, we simply need to provide the proper connection configuration similar to the example below.  

```typescript

const brokerConfig = {
  "metadata.broker.list": "{{ BOOTSTRAP_SERVER }}",
  "security.protocol": "SASL_SSL",
  "sasl.mechanism": "PLAIN",
  "sasl.username": "{{ CLUSTER_API_KEY }}",
  "sasl.password": "{{ CLUSTER_API_SECRET }}"
};

const config = {
  kafka: {
    CustomerTopic: "customers",
    ProducerGlobalConfig: {
      ...brokerConfig,            
    } as Kafka.ProducerGlobalConfig,    
    ConsumerGlobalConfig: {
      ...brokerConfig,
    } as Kafka.ConsumerGlobalConfig,
    ConsumerTopicConfig: {
      "auto.offset.reset": "earliest"
    } as Kafka.ConsumerTopicConfig
  },
  schemaRegistry: {
    host: "{{ SR_URL }}",
    auth: {
      username: "{{ SR_USER }}",
      password: "{{ SR_SECRET }}"

    }
  } as SchemaRegistryAPIClientArgs
};
```

NOTE: we use environment variables for supplying all these configuration settings when running in a deployment environment.  The above example is shown only to demonstrate how you will run this in your local environment.  Please refer to the Helm chart configuration that showcases how we use Kubernetes secrets and environment variables to inject all these configuration settings.

## Running Tests

The project has both unit and integration test.  The integration test uses [testcontainers](https://github.com/testcontainers/testcontainers-node) so make sure you have docker running when running integration tests

```shell
npm run test  # runs unit tests
npm run test:integration # runs integration tests
```


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
  helm upgrade node-kafka-demo helm/node-kafka-demo --namespace ashkal --install -f ./helm/values-cloud.yaml
  ```
- Inspect the deployment and container logs to verify that the Demo App is able to connect to confluent cloud  



Refer to the [Dockerfile](./Dockerfile) and [Helm chart](./helm/node-kafka-demo) provided for more details.

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



