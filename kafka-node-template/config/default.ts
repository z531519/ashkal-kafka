import { SchemaRegistryAPIClientArgs } from "@kafkajs/confluent-schema-registry/dist/api";

import Kafka from "node-rdkafka";

// all available configurations are documented here https://github.com/confluentinc/librdkafka/blob/v2.0.2/CONFIGURATION.md

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



module.exports = config;
