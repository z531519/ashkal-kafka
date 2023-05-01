import { SchemaRegistry, SchemaType } from "@kafkajs/confluent-schema-registry";
import { SchemaRegistryAPIClientArgs } from "@kafkajs/confluent-schema-registry/dist/api";
import Kafka from "node-rdkafka";
import { Consumer } from "./consumer";
import { Producer, RecordMetadata } from "./producer";

/**
 * Avro Model classes can implement this interface.
 * The class should then provide the avro schema in the Schema property
 */
export interface AvroModel {
  Schema: string;
}

export interface AvroRecordMetadata<T extends AvroModel> {
  topic: string;
  offset: number;
  partition: number;
  size: 10;
  key?: Buffer;
  timestamp: number;

  value?: T;
}

/**
 * Consumer that supports Avro messages using SchemaRegistry
 */
export class AvroConsumer extends Consumer {
  registry: SchemaRegistry;

  constructor(opts: SchemaRegistryAPIClientArgs) {
    super();
    this.registry = new SchemaRegistry(opts);
  }

  runAvro<T extends AvroModel>(
    cb: (err: Kafka.LibrdKafkaError, messages: AvroRecordMetadata<T>) => void
  ): void {
    const _this = this;
    super.run(async (err: Kafka.LibrdKafkaError, messages: any) => {
      // Discard invalid avro messages
      if (messages.value.readUInt8(0) !== 0) return;

      const value = await _this.registry.decode(messages.value);

      const decodedMessage = {
        ...messages,
        value
      };

      cb(err, decodedMessage);
    });
  }

  run(cb: (err: Kafka.LibrdKafkaError, messages: any) => void): void {
    const _this = this;
    super.run(async (err: Kafka.LibrdKafkaError, messages: any) => {
      // Discard invalid avro messages
      if (messages.value.readUInt8(0) !== 0) return;

      const value = await _this.registry.decode(messages.value);

      const decodedMessage = {
        ...messages,
        value
      };

      cb(err, decodedMessage);
    });
  }
}
/**
 * a Producer that can handle Avro messages
 */
export class AvroProducer extends Producer {
  registry: SchemaRegistry;

  constructor(opts: SchemaRegistryAPIClientArgs) {
    super();
    this.registry = new SchemaRegistry(opts);
  }

  async sendAvro<T extends AvroModel>(
    topic,
    key,
    value: T
  ): Promise<RecordMetadata> {
    const encoded = await this.encodeMessages(
      topic,
      value.Schema,
      value as any
    );
    return super.send(topic, key, encoded);
  }

  async sendAvroWithSchema(topic, key, value, schema): Promise<RecordMetadata> {
    const encoded = await this.encodeMessages(topic, schema, value);
    return super.send(topic, key, encoded);
  }

  async encodeMessages(
    topic: string,
    schema: string,
    value: Buffer
  ): Promise<Buffer> {
    const { id } = await this.registry.register(
      {
        type: SchemaType.AVRO,
        schema
      },
      { subject: `${topic}-value` } // use topic naming strategy
    );

    const encoded = await this.registry.encode(id, value);
    return encoded;
  }
}
