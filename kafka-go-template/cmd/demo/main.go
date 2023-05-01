package main

import (
	"flag"
	"fmt"
	"ashkal/kafka-template/pkg/customer"
	"ashkal/kafka-template/pkg/kafka"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	configfile *string
)

func init() {
	configfile = flag.String("config", "config", "target config file")
}

// Producer API App
func main() {
	flag.PrintDefaults()
	flag.Parse()
	config, err := kafka.InitConfig(*configfile)
	if err != nil {
		fmt.Println(err)
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/customers", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello")
	})

	r.POST("/customers", func(c *gin.Context) {
		var customer customer.Customer
		if err := c.BindJSON(&customer); err != nil {
			panic(err)
		}
		m, err := producer.Produce(config.TopicName, &customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, m)
	})

	go r.Run()

	_, stopchan := kafka.StartConsumer[customer.Customer](config, customer.NewHandler())

	// wait for consumer to stop, typically with a system int/stop signal
	<-stopchan
}
