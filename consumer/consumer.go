package consumer

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)


func GetConsumer()(*kafka.Consumer, error){

	config := &kafka.ConfigMap{
		"metadata.broker.list": os.Getenv("BROKER_LIST"),
		"security.protocol":    "SASL_SSL",
		"sasl.mechanisms":      "SCRAM-SHA-256",
		"sasl.username":        os.Getenv("KAFKA_USER"),
		"sasl.password":        os.Getenv("KAFKA_PASSWORD"),
		"group.id":                        "cctupsea-cloudkarafka-example",
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
		//"debug":                           "generic,broker,security",
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return &kafka.Consumer{}, err
	}
	fmt.Printf("Created Consumer %v\n", c)

	return c, nil

}