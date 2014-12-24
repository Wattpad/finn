package finn

// MockQueue is a fake queue that is meant for testing workers
type MockQueue struct {
	BaseQueue
	Topics      map[string]*MockTopic
	initialized bool
}

func (self *MockQueue) Initialize(config QueueConfig) error {
	if !self.initialized {
		self.SetConfig(config, QueueConfig{})
		self.Topics = make(map[string]*MockTopic)
		self.initialized = true
	}
	return nil
}

func (self *MockQueue) NewTopic(topic string) (GenericTopic, error) {
	if self.Topics[topic] == nil {
		self.Topics[topic] = &MockTopic{}
	}

	return self.Topics[topic], nil
}

// GetTopic returns a topic for testing
func (self *MockQueue) GetTopic(topicName string) *MockTopic {
	topic, _ := self.NewTopic(topicName)
	return topic.(*MockTopic)
}

func (self *MockQueue) Close() error {
	return nil
}

type MockTopic struct {
	stream chan []byte
}

func (self *MockTopic) Stream() (<-chan []byte, error) {
	if self.stream == nil {
		self.stream = make(chan []byte)
	}
	return self.stream, nil
}

func (self *MockTopic) Close() error {
	return nil
}

// Put is a non-blocking testing function that adds a message to the stream
func (self *MockTopic) Put(message []byte) error {
	if self.stream == nil {
		self.stream = make(chan []byte)
	}

	go func() {
		self.stream <- message
	}()

	return nil
}

// PutWorker is a non-blocking testing function that adds a worker to the stream
func (self *MockTopic) PutWorker(worker GenericWorker) {
	self.Put(Pack(worker))
}
