import path from "path";
import startConsumer from "../../src/customer/consumer";
import { sink } from "../../src/customer/sink";
import Kafka from "node-rdkafka";

import { DockerComposeEnvironment, StartedDockerComposeEnvironment } from "testcontainers";
import config from "config";
import { buildCustomer, Customer, Schema } from "../../src/customer/customer";

import waitForExpect from "wait-for-expect";
import { AvroConsumer, AvroProducer } from "../../src/kafka";

jest.mock("../../src/customer/sink", () => {
  return {
    sink: jest.fn().mockImplementation(data => {
      console.log(data);
    })
  };
});

describe("consumer handler integration", () => {
  let environment:StartedDockerComposeEnvironment;
  let producer: AvroProducer;
  let consumer: AvroConsumer;

  beforeAll(async () => {
    // startup our kafka container for testing
    const composeFilePath = path.resolve(__dirname, "../..");
    const composeFile = "docker-compose.yml";

    console.log("Testcontainers startup");
    environment = await new DockerComposeEnvironment(
      composeFilePath,
      composeFile
    ).up();
    console.log("Testcontainers startup complete");

    const client = Kafka.AdminClient.create(
      config.get("kafka.ProducerGlobalConfig")
    );

    client.createTopic({
      topic: "customers",
      num_partitions: 4,
      replication_factor: 1
    });
    
    client.disconnect();

    // get our producer and consumer running
    producer = await new AvroProducer(config.get("schemaRegistry")).create(
      config.get("kafka.ProducerGlobalConfig"),
      config.get("kafka.ProducerTopicConfig")
    );
    await producer.connect();
    
    
    let metadata
    await waitForExpect(() => {      
      producer.producer.getMetadata( {
        timeout: 1000,
        topic: "customers"
      },  function(err, m) {
        if (err) {
          console.error(err);
        } else {
          metadata = m;
        }
      });
      
      expect(metadata.topics.map( i => i.name).includes('customers')).toBeTruthy();
    }, 10000, 1000);

    consumer = await startConsumer();
  });

  afterAll(async () => {
    // make sure producer/consumer are stopped before shutting down containers
    producer ? await producer.producer.disconnect() : undefined;

    consumer ? consumer.consumer.commit() : undefined;
    consumer ? consumer.consumer.unsubscribe() : undefined;
    consumer ? await consumer.consumer.disconnect() : undefined;
    
    // we are ready
    if (environment) {
      console.log("Testcontainers teardown");
      await environment.down();      
      console.log("Testcontainers teardown complete");
    }
  });

  const testCustomer = birthdt =>
    buildCustomer({
      fname: "fname",
      id: "1",
      birthdt,
      fullname: "fullname",
      gender: "M",
      lname: "lname",
      mname: "mname",
      suffix: "suffix",
      title: "title",
      largePayload: "largePayload"
    });

  it("Should produce and consume adult consumer record()", async () => {
    const customer = testCustomer("1990-01-01");
    await producer.sendAvro<Customer>("customers", customer.id, customer);
    await waitForExpect(
      () => {
        expect(sink).toHaveBeenCalledWith(
          "Customer: fullname, birthdt: 1990-01-01"
        );
      },
      30000,
      1000
    );
  });

  it("Should produce and consume MINOR consumer record()", async () => {
    const customer = testCustomer("2010-01-01");
    await producer.sendAvro<Customer>("customers", customer.id, customer);
    await waitForExpect(
      () => {
        expect(sink).toHaveBeenCalledWith(
          "MINOR> Customer: fullname, birthdt: 2010-01-01"
        );
      },
      30000,
      1000
    );
  });
});
