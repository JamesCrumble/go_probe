package probe_kafka

import (
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func CreateTopic(addr string, topic string) error {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{{Topic: topic, NumPartitions: 1, ReplicationFactor: 1}}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTopic(addr string, topic string) error {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	err = controllerConn.DeleteTopics(topic)
	if err != nil {
		return err
	}

	return nil
}
