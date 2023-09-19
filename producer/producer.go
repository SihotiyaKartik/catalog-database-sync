//producer.go

package producer

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)


func GetProducer()(*kafka.Producer, error) {
	config := &kafka.ConfigMap{
		"metadata.broker.list": os.Getenv("BROKER_LIST"),
		"security.protocol":    "SASL_SSL",
		"sasl.mechanisms":      "SCRAM-SHA-256",
		"sasl.username":        os.Getenv("KAFKA_USER"),
		"sasl.password":        os.Getenv("KAFKA_PASSWORD"),
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
		//"debug":                           "generic,broker,security",
	}
	
	p, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return &kafka.Producer{}, err
	}
	fmt.Printf("Created Producer %v\n", p)

	return p, nil
	
}