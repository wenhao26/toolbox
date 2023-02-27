package kafka

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// Connection configuration
type Option struct {
	Addr []string
}

// Serve structure
type Serve struct {
	Producer      sarama.SyncProducer
	AsyncProducer sarama.AsyncProducer
	Consumer      sarama.Consumer
}

type Result struct {
	Partition int32
	Offset    int64
}

// Create SyncProducer Serve instance
func NewSyncProducer(option *Option) (*Serve, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewSyncProducer(option.Addr, config)
	if err != nil {
		return nil, err
	}
	return &Serve{Producer: producer}, nil
}

// Create AsyncProducer Serve instance
func NewAsyncProducer(option *Option) (*Serve, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewAsyncProducer(option.Addr, config)
	if err != nil {
		return nil, err
	}
	return &Serve{AsyncProducer: producer}, nil
}

// Create consumer Serve instance
func NewConsumer(option *Option) (*Serve, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	//config.Version = sarama.V2_1_1_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second * 5
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewClient(option.Addr, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &Serve{Consumer: consumer}, nil
}

// Close SyncProducer
func (s *Serve) CloseSyncProducer() {
	s.Producer.Close()
}

// Close AsyncProducer
func (s *Serve) CloseAsyncProducer() {
	s.AsyncProducer.Close()
}

// Close Consumer
func (s *Serve) CloseConsumer() {
	s.Consumer.Close()
}

// Publish message to topic
func (s *Serve) Release(topic string, data []byte) (Result, error) {
	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		//Key: sarama.StringEncoder(""),
		Value: sarama.ByteEncoder(data),
	}
	partition, offset, err := s.Producer.SendMessage(producerMessage)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Partition: partition,
		Offset:    offset,
	}, nil
}

// Async Publish message to topic
func (s *Serve) AsyncRelease(topic string, data []byte) {
	producer := s.AsyncProducer
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case r := <-p.Successes():
				fmt.Println("offset: ", r.Offset, "partitions: ", r.Partition, "timestamp: ", r.Timestamp.String())
			case r := <-p.Errors():
				fmt.Println("error: ", r.Error())
			}
		}
	}(producer)

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		//Key:   sarama.StringEncoder(""),
		Value: sarama.ByteEncoder(data),
	}
	producer.Input() <- producerMessage
}

// Receive message from topic
func (s *Serve) Receive(topic string) {
	partitions, err := s.Consumer.Partitions(topic)
	if err != nil {
		log.Println(err)
		return
	}
	defer s.CloseConsumer()

	var wg sync.WaitGroup
	for partition := range partitions {
		pc, err := s.Consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to start the user of partition [%d]：%s\n", partition, err)
			return
		}
		defer pc.Close()

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			if err := recover(); err != nil {
				return
			}
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	wg.Wait()
}
