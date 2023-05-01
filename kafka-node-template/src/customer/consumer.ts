import { Customer } from "./customer";

import config from "config";
import { AvroConsumer } from "../kafka/avro";
import { sink } from "./sink";

const YEAR_MS = 31556926000;

export const eachMessageHandler = function(err, data) {
  const customer: Customer = data.value;
  const today = new Date();
  const birthdt = new Date(customer.birthdt);
  const diff = (Number(today) - Number(birthdt)) / YEAR_MS;
  diff < 21
    ? sink(
        `MINOR> Customer: ${customer.fullname}, birthdt: ${customer.birthdt}`
      )
    : sink(`Customer: ${customer.fullname}, birthdt: ${customer.birthdt}`);
};

const startConsumer = async () => {
  const topic: string = config.get("kafka.CustomerTopic");
  const consumer = new AvroConsumer(config.get("schemaRegistry")).create(
    config.get("kafka.ConsumerGlobalConfig"),
    config.get("kafka.ConsumerTopicConfig")
  );

  await consumer.connect();
  consumer.subscribe([topic]);
  consumer.runAvro<Customer>(eachMessageHandler);
  return consumer;
};

export default startConsumer;
