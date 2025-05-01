package entities

type KafkaMessage struct {
	TopicPartition string
	Value          []byte
}
