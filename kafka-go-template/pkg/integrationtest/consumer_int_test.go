//go:build integration
// +build integration

package integrationtest

import (
	"fmt"
	"ashkal/kafka-template/pkg/customer"
	mock_customer "ashkal/kafka-template/pkg/customer/mocks"
	"ashkal/kafka-template/pkg/kafka"
	"testing"
	"time"

	"github.com/ecodia/golang-awaitility/awaitility"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestCustomer(birthDt string) customer.Customer {
	payload := customer.NewCustomer()
	payload.Fname = "fname"
	payload.Fullname = "fullname"
	payload.Lname = "lname"
	payload.Gender = "M"
	payload.Mname = "mname"
	payload.Id = "1"
	payload.LargePayload = "abc"
	payload.Birthdt = birthDt
	payload.Title = "title"
	payload.Suffix = "suffix"
	return payload
}

// Produce and Consume Customer records, both adult and minor
func TestCustomer(t *testing.T) {
	tests := []struct {
		name     string
		customer customer.Customer
		expected string
	}{
		{name: "adult", customer: createTestCustomer("1990-01-01"), expected: "Customer: fullname, birthdt: 1990-01-01"},
		{name: "minor", customer: createTestCustomer("2020-01-01"), expected: "MINOR> Customer: fullname, birthdt: 2020-01-01"},
	}

	config, err := kafka.InitConfig("../../config")
	assert.Nil(t, err)

	config.TopicName = config.TopicName + uuid.New().String()
	fmt.Println(config.TopicName)
	producer, err := kafka.NewProducer(config)
	assert.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sink := mock_customer.NewMockSink(ctrl)
	customerHandler := customer.NewHandlerWithSink(sink)
	consumer, _ := kafka.StartConsumer[customer.Customer](config, customerHandler)
	defer consumer.Stop()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			var receivedMessage interface{}
			sink.EXPECT().Print(tc.expected).Do(func(val string) {
				fmt.Println(val)
				receivedMessage = val
			}).Times(1)

			_, err = producer.Produce(config.TopicName, &tc.customer)
			assert.Nil(t, err)

			err = awaitility.Await(500*time.Millisecond, 30*time.Second, func() bool {
				return receivedMessage != nil
			})
			assert.Nil(t, err)
			assert.NotNil(t, receivedMessage)
			assert.Equal(t, tc.expected, receivedMessage)
		})
	}
}
