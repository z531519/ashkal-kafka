import startConsumer from "./customer/consumer";
import startApi from "./api";

import Kafka from "node-rdkafka";
import config from "config";

const app = async (): Promise<void> => {
  console.log("Ashkal Kafka Template");
  
  console.log(`rdkafka version: ${Kafka.librdkafkaVersion}`);
  console.log(`Features: ${Kafka.features}`);

  console.log("Configuration:");
  console.log(config);

  startConsumer();
  startApi();

  return Promise.resolve();
};

export default app;

app();
