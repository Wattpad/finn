package finn

import (
	"fmt"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

// Default configuration for RabbitMQ
var defaultConfig = QueueConfig{
	"host":     "localhost",
	"port":     "5672",
	"user":     "guest",
	"password": "guest",
	"prefetch": "3",
}

// RabbitQueue represents a connection to RabbitMQ
type RabbitQueue struct {
	BaseQueue
	connection       *amqp.Connection
	topics           map[string]*RabbitTopic
	disconnectSignal chan string
	reconnectSignal  chan bool
	running          bool
}

// Initialize sets up the initial connection to RabbitMQ, and initalizes variables.
// Should only be called once.
func (self *RabbitQueue) Initialize(config QueueConfig) error {

	if self.connection != nil {
		return fmt.Errorf("RabbitQueue was already initialized!")
	}

	self.SetConfig(config, defaultConfig)
	if err := self.connect(self.config["user"], self.config["password"], self.config["host"], self.config["port"]); err != nil {
		return err
	}

	self.topics = make(map[string]*RabbitTopic)
	self.disconnectSignal = make(chan string)
	self.reconnectSignal = make(chan bool)
	self.running = true

	go self.disconnectionHandler()

	return nil
}

// connect is a helper function for Initialize that sets up the connection to the RabbitMQ server.
// Should only be called once.
func (self *RabbitQueue) connect(user string, password string, host string, port string) error {
	var err error
	if self.connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)); err != nil {
		return err
	}

	return nil
}

// disconnectionHandler handles disconnection errors and transparently attempts reconnection
func (self *RabbitQueue) disconnectionHandler() {
	// Number of topics that have been disconnected
	var numDisconnected int = 0

	// Wait for disconnects
	for _ = range self.disconnectSignal {
		numDisconnected++

		// If all topics have disconnected, then attempt to reconnect to RabbitMQ server
		if numDisconnected >= len(self.topics) {
			LogInfo("Attempting reconnection to RabbitMQ...")

			err := self.connect(self.config["user"], self.config["password"], self.config["host"], self.config["port"])
			for err != nil {
				// Wait a bit and then try again
				LogInfo("Reconnection failed, sleeping and trying again...")
				time.Sleep(30 * time.Second)
				err = self.connect(self.config["user"], self.config["password"], self.config["host"], self.config["port"])
			}

			// Successful reconnection, now reconnect topics
			for _, topic := range self.topics {
				topic.channel, _ = self.connection.Channel()
				topic.reconnect()
			}

			// Reset disconnections
			numDisconnected = 0

			// Wake up topics
			for _ = range self.topics {
				self.reconnectSignal <- true
			}
		}
	}
}

// NewTopic creates a new channel, and then listens to a topic on that channel
func (self *RabbitQueue) NewTopic(topic string) (GenericTopic, error) {

	channel, err := self.connection.Channel()
	if err != nil {
		return nil, err
	}

	// Set max num of messages to prefetch before ack-ing
	prefetch, err := strconv.Atoi(self.config["prefetch"])
	if err != nil {
		return nil, err
	}
	channel.Qos(prefetch, 0, false)

	// Declare queue to ensure that it exists
	queue, err := channel.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		channel.Close()
		return nil, err
	}

	self.topics[topic] = &RabbitTopic{
		channel:      channel,
		queue:        &queue,
		disconnected: self.disconnectSignal,
		reconnected:  self.reconnectSignal,
		running:      true,
	}

	return self.topics[topic], nil
}

// Close shuts down the connection to RabbitMQ
func (self *RabbitQueue) Close() error {
	self.running = false
	for _, topic := range self.topics {
		topic.Close()
	}
	self.connection.Close()
	return nil
}

// RabbitTopic is an abstraction on top of RabbitMQ queues
type RabbitTopic struct {
	channel      *amqp.Channel        // RabbitMQ channel
	queue        *amqp.Queue          // RabbitMQ queue
	deliveries   <-chan amqp.Delivery // Channel for messages in
	stream       chan []byte          // Channel for messages out
	disconnected chan string          // Global disconnect signal channel
	reconnected  chan bool            // Global reconnection signal channel
	running      bool
}

// Stream returns a blocking channel with deliveries from the topic.
// It also translates RabbitMQ messages into general-purpose ones.
func (self *RabbitTopic) Stream() (<-chan []byte, error) {
	var err error
	self.deliveries, err = self.channel.Consume(self.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	self.stream = make(chan []byte)

	go func() {
		for self.running {
			for delivery := range self.deliveries {
				delivery.Ack(false)
				self.stream <- delivery.Body
			}

			// If not an intentional shutdown, attempt reconnection
			if self.running {
				LogError(fmt.Errorf("Deliveries stopped on topic: %s", self.queue.Name))

				// Send disconnected signal, then wait for reconnection
				self.disconnected <- self.queue.Name
				<-self.reconnected
			}
		}
	}()

	return self.stream, nil
}

// reconnect is a helper function for RabbitQueue.disconnectionHandler
func (self *RabbitTopic) reconnect() error {
	var err error
	self.deliveries, err = self.channel.Consume(self.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

// Close shuts down the topic listener, as well as closes the output stream
func (self *RabbitTopic) Close() error {
	self.running = false
	self.channel.Close()
	close(self.stream)
	return nil
}
