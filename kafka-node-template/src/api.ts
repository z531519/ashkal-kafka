import express from "express";
import config from "config";
import { buildCustomer, Customer } from "./customer/customer";
import { AvroProducer } from "./kafka";
import bodyParser from "body-parser";

const startApi = async () => {
  const topic: string = config.get("kafka.CustomerTopic");
  const producer = await new AvroProducer(config.get("schemaRegistry")).create(
    config.get("kafka.ProducerGlobalConfig"),
    config.get("kafka.ProducerTopicConfig")
  );
  
  await producer.connect();

  const app = express();
  app.use(bodyParser.json({ limit: "10mb" }));

  const port = process.env.PORT || 8080;

  app.get("/customers", async (req, res) => {
    return res.send("Hello");
  });

  app.post("/customers", async (req, res) => {
    const customer: Customer = buildCustomer(req.body);
    const result = await producer.sendAvro<Customer>(
      topic,
      customer.id,
      customer
    );
    delete result.value;
    return res.send(result);
  });

  app.listen(port, () => {
    console.log(`⚡️[server]: Server is running at https://localhost:${port}`);
  });
};

export default startApi;
