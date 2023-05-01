//go:build integration
// +build integration

package integrationtest

import (
	"fmt"
	"ashkal/kafka-template/pkg/customer"
	"ashkal/kafka-template/pkg/kafka"
	"testing"
	"time"

	"github.com/ecodia/golang-awaitility/awaitility"
	"github.com/stretchr/testify/assert"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestProduceConsumeBasicString(t *testing.T) {

	config, err := kafka.InitConfig("../../config")
	assert.Nil(t, err)

	config.TopicName = "topic-basic-string"

	producer, err := kafka.NewProducer(config)
	assert.Nil(t, err)

	payload := "Basic String"

	_, err = producer.Produce(config.TopicName, payload)
	assert.Nil(t, err)

	var receivedMessage interface{}
	consumer, _ := kafka.StartConsumer[string](config, func(value interface{}, msg *confluentKafka.Message) error {
		fmt.Println(value)
		receivedMessage = value
		return nil
	})

	err = awaitility.Await(500*time.Millisecond, 20*time.Second, func() bool {
		return receivedMessage != nil
	})
	assert.Nil(t, err)
	consumer.Stop()

	fmt.Println("Consumer Stopped")
	assert.NotNil(t, receivedMessage)
	assert.Equal(t, payload, receivedMessage)

}

func TestProduceConsumeCustomerAvro(t *testing.T) {

	config, err := kafka.InitConfig("../../config")
	assert.Nil(t, err)

	config.TopicName = "topic-avro"

	producer, err := kafka.NewProducer(config)
	assert.Nil(t, err)

	payload := createTestCustomer("2020-01-01")

	_, err = producer.Produce(config.TopicName, &payload)
	assert.Nil(t, err)

	var receivedMessage interface{}
	consumer, _ := kafka.StartConsumer[customer.Customer](config, func(value interface{}, msg *confluentKafka.Message) error {
		fmt.Println(value)
		receivedMessage = value
		return nil
	})

	err = awaitility.Await(500*time.Millisecond, 20*time.Second, func() bool {
		return receivedMessage != nil
	})
	assert.Nil(t, err)
	consumer.Stop()

	fmt.Println("Consumer Stopped")
	assert.NotNil(t, receivedMessage)
	assert.Equal(t, payload, receivedMessage)

}
