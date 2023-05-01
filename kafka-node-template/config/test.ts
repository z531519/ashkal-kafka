import { SchemaRegistryAPIClientArgs } from "@kafkajs/confluent-schema-registry/dist/api";

import Kafka from "node-rdkafka";

const config = {
  kafka: {
    ProducerGlobalConfig: {
      "metadata.broker.list": "localhost:19092"
    } as Kafka.ProducerGlobalConfig,
    ConsumerGlobalConfig: {
      "metadata.broker.list": "localhost:19092",
      "topic.metadata.refresh.interval.ms": 1000
    } as Kafka.ConsumerGlobalConfig
  }
};

module.exports = config;
