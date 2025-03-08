package probe_kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	DEBUG             = false
	KAFKA_ADDR        = "localhost:9092"
	PROCESSABLE_COUNT = 50000
	TIMEOUT_SECONDS   = 10
)

func logDebug(msg string) {
	if !DEBUG {
		return
	}
	log.Println(msg)
}

func Probe(ctx context.Context) ([]string, error) {
	topic := "test-topic"

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{KAFKA_ADDR},
		Balancer:     &kafka.Hash{},
		RequiredAcks: -1,
		Async:        true,
		ReadTimeout:  10 * time.Second,
		BatchSize:    2048,
		BatchTimeout: 10 * time.Millisecond,
	})
	defer kafkaWriter.Close()

	logDebug("Creating consumer")
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               []string{KAFKA_ADDR},
		GroupID:               "",
		Topic:                 topic,
		MinBytes:              1,
		MaxBytes:              1000000,
		MaxWait:               100 * time.Millisecond,
		WatchPartitionChanges: true,
		CommitInterval:        time.Second,
	})
	defer kafkaReader.Close()

	logDebug("Producing data")
	var withErrs int = 0
	var kafkaMsg kafka.Message
	for i := range PROCESSABLE_COUNT {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			msgContent := fmt.Sprintf("msg_%d", i)
			if DEBUG {
				fmt.Printf("Producer: %s", fmt.Sprintf("Messages %d out of %d\r", i+1, PROCESSABLE_COUNT))
			}
			kafkaMsg = kafka.Message{Key: []byte(fmt.Sprintf("msg_%d", i)), Value: []byte(msgContent), Topic: topic}
			if err := kafkaWriter.WriteMessages(ctx, kafkaMsg); err != nil {
				logDebug(err.Error())
				withErrs++
			}
		}
	}
	if DEBUG {
		fmt.Println()
	}

	logDebug("Consuming data")
	msgs := make([]string, PROCESSABLE_COUNT-withErrs)
	for i := range PROCESSABLE_COUNT - withErrs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			msg, err := kafkaReader.FetchMessage(ctx)
			if err != nil {
				return nil, err
			}
			msgs[i] = string(msg.Value)
		}
	}

	return msgs, nil
}
