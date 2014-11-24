package finn

// QueueConfig is used to do queue-specific configuration (e.g. will be different for RabbitMQ vs Kafka)
type QueueConfig map[string]string

// GenericQueue is the interface all queues using Finn must match (ie maps 1:1 with RabbitMQ, Kafka)
type GenericQueue interface {
	// Initialize sets up the initial connection to the queue service
	Initialize(QueueConfig) error

	// NewTopic creates a new topic to listen on
	NewTopic(string) (GenericTopic, error)

	// Close shuts down the queue
	Close() error
}

// GenericTopic represents a stream of messages from a topic in a queue
type GenericTopic interface {
	// Stream returns a channel with a stream of messages from the topic
	Stream() (<-chan []byte, error)

	// Close shuts down the connection to the topic
	Close() error
}
