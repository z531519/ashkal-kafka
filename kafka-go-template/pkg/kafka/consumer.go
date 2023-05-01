package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

// callback function for each message record received.
type ConsumerHandlerFunc func(value interface{}, msg *kafka.Message) error

type Consumer struct {
	stopchan chan struct{}
}

// Starts a kafka consumer in the background using the provided config.  The `handler` callback function is
// called for each message received.  This consumer can handle either a string or an avro message by simply passing
// either a string or a struct type in the T generic parameter.  A non-string type will automatically attempt an avro deserialization.
//
// Returns the Consumer instance, and a channel that can be used to wait for this consumer termination.
func StartConsumer[T any](config Config, handler ConsumerHandlerFunc) (Consumer, chan struct{}) {
	topic := config.TopicName

	kafkaConfigMap := kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"client.id":         config.ClientId,
	}

	for k, v := range config.Properties {
		kafkaConfigMap[k] = v
	}
	for k, v := range config.Consumer {
		kafkaConfigMap[k] = v
	}

	// Create Consumer instance
	c, err := kafka.NewConsumer(&kafkaConfigMap)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	var client schemaregistry.Client
	if config.SchemaRegistryUser != nil {
		client, err = schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
			config.SchemaRegistryUrl, *config.SchemaRegistryUser, *config.SchemaRegistryPassword))
	} else {
		client, err = schemaregistry.NewClient(schemaregistry.NewConfig(config.SchemaRegistryUrl))
	}

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	deser, err := avro.NewSpecificDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())

	if err != nil {
		fmt.Printf("Failed to create deserializer: %s\n", err)
		os.Exit(1)
	}

	// Subscribe to topic
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	stopchan := make(chan struct{})

	// Process messages
	go func() {
		run := true
		for run {
			select {
			case <-stopchan:
				fmt.Println("Stopping Consumer")
				run = false
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				close(stopchan)
				run = false
			default:
				msg, err := c.ReadMessage(100 * time.Millisecond)
				if err != nil {
					// Errors are informational and automatically handled by the consumer
					continue
				}

				var target T

				if isAvro(target) {
					err = deser.DeserializeInto(*msg.TopicPartition.Topic, msg.Value, &target)
					if err != nil {
						fmt.Printf("Failed to deserialize payload: %s\n", err)
						panic(err)
					}
					err = handler(target, msg)
					if err != nil {
						panic(err)
					}
				} else {
					err = handler(string(msg.Value), msg)
					if err != nil {
						panic(err)
					}
				}

			}
		}

		fmt.Printf("Closing consumer\n")
		c.Close()
	}()
	return Consumer{
		stopchan,
	}, stopchan
}

func (c Consumer) Stop() {
	close(c.stopchan)
}

func isAvro(msg interface{}) bool {
	switch msg.(type) {
	case string:
		return false
	default:
		return true
	}
}
