/*
Package finn is a framework for easily creating workers that listen to jobs on a message queue (e.g. RabbitMQ, Kafka)

The following is an example of a worker that we can write.

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

Here we register our worker with finn, then set the queue and queue configuration.
After that we start finn up.

	package main

	import "github.com/Wattpad/finn"

	func main() {
		finn.AddWorker(&finn.MyWorker{})
		finn.SetQueue(&finn.RabbitQueue{}, finn.QueueConfig{"host": "localhost"})
		finn.Listen()
	}
*/
package finn
