finn
========

Finn is a Go module for creating workers from any message queue (e.g. RabbitMQ)

### Sample Worker

	package workers

	import "github.com/Wattpad/finn"

	type MyWorker struct {
		finn.BaseWorker
	}

	func (self *MyWorker) NewInstance() finn.GenericWorker {
		return &MyWorker{}
	}

	func (self *MyWorker) Name() string {
		return "MyWorker"
	}

	func (self *MyWorker) TopicName() string {
		return "my-topic"
	}

	func (self *MyWorker) Run() {
		// do something
	}

### Sample main()

	package main

	import "github.com/Wattpad/finn"

	func main() {
		finn.AddWorker(&finn.MyWorker{})
		finn.SetQueue(&finn.RabbitQueue{}, finn.QueueConfig{"host": "localhost"})
		finn.Listen()
	}
