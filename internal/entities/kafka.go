package entities

type KafkaMessage struct {
	TopicPartition string
	Value          []byte
}

type WordMessage struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}
