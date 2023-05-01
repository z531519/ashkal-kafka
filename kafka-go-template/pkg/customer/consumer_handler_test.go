package customer

import (
	mock_customer "ashkal/kafka-template/pkg/customer/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

////
// Unit test for the Consumer Handler

func CreateTestCustomer(birthDt string) Customer {
	payload := NewCustomer()
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

func TestHandleMinorCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sink := mock_customer.NewMockSink(ctrl)

	sink.EXPECT().Print("MINOR> Customer: fullname, birthdt: 2020-01-01").Times(1)

	customerHandler := NewHandlerWithSink(sink)

	payload := CreateTestCustomer("2020-01-01")
	err := customerHandler(payload, nil)
	assert.Nil(t, err)
}

func TestHandleAdultCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sink := mock_customer.NewMockSink(ctrl)

	sink.EXPECT().Print("Customer: fullname, birthdt: 1960-01-01").Times(1)

	customerHandler := NewHandlerWithSink(sink)

	payload := CreateTestCustomer("1960-01-01")
	err := customerHandler(payload, nil)
	assert.Nil(t, err)
}
