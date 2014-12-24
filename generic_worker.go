package finn

import "time"

// GenericWorker is the interface that all workers using Finn must match
type GenericWorker interface {
	// Name returns the name of the worker
	Name() string

	// Run performs the actual work.
	// Returns an error and whether to retry or not (bool)
	Run() (error, bool)

	// NewInstance returns a new instance of the worker
	NewInstance() GenericWorker

	// RunDelay returns the duration to sleep before running
	RunDelay() time.Duration

	// RetryDelay returns the duration to wait until the retry should start
	RetryDelay(int) time.Duration

	// RetryDelaySeconds returns the base number of seconds to wait before retrying
	RetryDelaySeconds() int

	// StartTime returns the time at which the worker should start running
	StartTime() time.Time

	// SetStartStamp sets the time at which the worker should start running
	SetStartStamp(int64)

	// NextStartStamp returns the next time the worker should start, based on RetryDelay()
	NextStartStamp(int) int64

	// Attempts returns the run attempt number
	Attempts() int

	// MaxAttempts returns the maximum number of attempts
	MaxAttempts() int

	// IncreaseAttempts increments the attempts counter by one
	IncreaseAttempts()

	// TopicName returns the queue/topic that the worker should listen on
	TopicName() string
}
