package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type Producer struct {
	kafkaProducer *kafka.Producer
	srClient      schemaregistry.Client
	ser           *avro.SpecificSerializer
}

// Build a new producer with the provided config
func NewProducer(config Config) (*Producer, error) {
	kafkaConfigMap := kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"client.id":         config.ClientId,
	}

	for k, v := range config.Properties {
		kafkaConfigMap[k] = v
	}
	for k, v := range config.Producer {
		kafkaConfigMap[k] = v
	}

	kafkaProducer, err := kafka.NewProducer(&kafkaConfigMap)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		return nil, err
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
		return nil, err
	}

	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())

	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		return nil, err
	}

	return &Producer{
		kafkaProducer: kafkaProducer,
		srClient:      client,
		ser:           ser,
	}, nil
}

func (p *Producer) Produce(topic string, payload interface{}) (*kafka.Message, error) {

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)
	var err error

	var serialized []byte
	switch data := payload.(type) {
	case []byte:
		serialized = data
	case string:
		serialized = []byte(data)
	case interface{}:
		serialized, err = p.ser.Serialize(topic, data)
		if err != nil {
			fmt.Printf("Failed to serialize payload: %s\n", err)
			return nil, err
		}
	}

	err = p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          serialized,
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		return nil, err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return m, nil
}
