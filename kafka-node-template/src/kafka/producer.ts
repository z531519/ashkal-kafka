import Kafka from "node-rdkafka";
import { v4 as uuidv4 } from "uuid";

export interface IProducer {
  connect(): Promise<IProducer>;
  onError(err: Error): void;
  send(key: any, value: Buffer, topic?: string): Promise<RecordMetadata>;
}

export interface RecordMetadata {
  topic: string;
  offset: number;
  partition: number;
  size: 10;
  key?: Buffer;
  value?: Buffer;
  timestamp: number;
}

export class Producer implements IProducer {
  connecting: boolean;
  connected: boolean;
  producer: Kafka.Producer;
  topic: string;
  config: any;
  producerPromisesCb: any = {};

  constructor() {}

  create(config, topicCfg?) {
    this.producer = new Kafka.Producer(
      {
        ...config,
        dr_cb: true,
        dr_msg_cb: true
      },
      topicCfg
    );
    this.config = config;
    return this;
  }

  connect(): Promise<Producer> {
    var _this = this;
    return new Promise(function(resolve, reject) {
      if (_this.connected) {
        resolve(_this);
      } else {
        _this.producer.connect({}, function(err, res) {
          if (err) return reject(err);
        });
        _this.connecting = true;
        _this.producer.on("ready", function() {
          _this.connected = true;
          _this.connecting = false;
          resolve(_this);
        });
        _this.producer.on("event.error", function(err) {
          if (!_this.connected) {
            return reject(err);
          }
          _this.onError(err);
        });
        _this.producer.on("delivery-report", function(err, report) {
          if (_this.producerPromisesCb[report.opaque]) {
            const cb = _this.producerPromisesCb[report.opaque];
            if (err !== null) {
              cb.reject(err);
            } else {
              cb.resolve(report);
            }
            delete _this.producerPromisesCb[report.opaque];
          }
        });
        _this.producer.setPollInterval(10);
      }
    });
  }

  send(topic, key, value): Promise<RecordMetadata> {
    var _this = this;
    const uuid = uuidv4();
    const p = new Promise<RecordMetadata>(function(resolve, reject) {
      _this.producerPromisesCb[uuid] = {
        resolve,
        reject
      };
      if (!_this.connected) {
        reject("Producer not connected");
      }

      const res = _this.producer.produce(topic, null, value, key, null, uuid);
    });

    return p;
  }

  onError(err) {
    console.error("[Producer] - ", err);
  }
}
