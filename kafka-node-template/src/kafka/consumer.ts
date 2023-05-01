import Kafka, { SubscribeTopicList } from "node-rdkafka";

export interface IConsumer {
  connect(): Promise<Consumer>;
  subscribe(topics: SubscribeTopicList);
  run(cb: (err: Kafka.LibrdKafkaError, messages: any) => void): void;
}

export class Consumer implements IConsumer {
  connecting: boolean = false;
  connected: boolean = false;

  config: any;
  consumer: Kafka.KafkaConsumer;

  constructor() {}

  connect(): Promise<Consumer> {
    var _this = this;

    return new Promise<Consumer>(function(resolve, reject) {
      if (_this.connected) {
        resolve(_this);
      } else {
        _this.consumer.connect({}, function(err, res) {
          if (err) return reject(err);
        });
        _this.connecting = true;
        _this.consumer.on("ready", function() {
          _this.connected = true;
          _this.connecting = false;
          resolve(_this);
        });       
        _this.consumer.on("event.error", function(...ev) {
          console.log(ev);
        });
        _this.consumer.on("event.log", function(ev) {
          console.log(ev.message);
        });
        _this.consumer.on("rebalance", function(...ev) {
          console.log(ev);
        });
      }
    });
  }

  subscribe(topics: SubscribeTopicList) {
    this.consumer.subscribe(topics);
  }

  run(cb: (err: Kafka.LibrdKafkaError, messages: any) => void): void {
    this.consumer.consume(cb);
  }

  create(config, topicCfg) {
    this.consumer = new Kafka.KafkaConsumer(
      {
        ...config
      },
      topicCfg
    );
    this.config = config;
    return this;
  }
}
