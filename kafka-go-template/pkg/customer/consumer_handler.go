package customer

import (
	"fmt"
	"time"

	"ashkal/kafka-template/pkg/kafka"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewHandler() kafka.ConsumerHandlerFunc {
	customerHandler := CustomerHandler{
		Sink: CustomerSink{},
	}
	return customerHandler.handler
}

func NewHandlerWithSink(sink Sink) kafka.ConsumerHandlerFunc {
	customerHandler := CustomerHandler{
		Sink: sink,
	}
	return customerHandler.handler
}

type CustomerHandler struct {
	Sink Sink
}

type CustomerSink struct{}

func (CustomerSink) Print(data string) {
	fmt.Println(data)
}

// a basic message handler for Customer records
func (h CustomerHandler) handler(value interface{}, msg *confluentKafka.Message) error {
	customer := value.(Customer)
	now := time.Now()
	birthdt, err := time.Parse("2006-01-02", customer.Birthdt)
	cutoff := now.AddDate(-21, 0, 0)

	if birthdt.After(cutoff) {
		h.Sink.Print(fmt.Sprintf("MINOR> Customer: %s, birthdt: %s", customer.Fullname, customer.Birthdt))
	} else {
		h.Sink.Print(fmt.Sprintf("Customer: %s, birthdt: %s", customer.Fullname, customer.Birthdt))
	}

	return err
}
