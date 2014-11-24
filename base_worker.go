package finn

import "time"

// BaseWorker contains common methods used by the different worker implementations
type BaseWorker struct {
	RunAttempts int
	StartStamp  int64
}

func (self *BaseWorker) Attempts() int {
	return self.RunAttempts
}

func (self *BaseWorker) IncreaseAttempts() {
	self.RunAttempts++
}

func (self *BaseWorker) MaxAttempts() int {
	return 3
}

func (self *BaseWorker) CanRun() bool {
	if self.StartStamp > 0 {
		return time.Now().After(self.StartTime())
	}

	return true
}

func (self *BaseWorker) RunDelay() time.Duration {
	if self.StartStamp > 0 {
		return self.StartTime().Sub(time.Now())
	}

	return 0
}

func (self *BaseWorker) RetryDelaySeconds() int {
	return 60
}

func (self *BaseWorker) RetryDelay(delay int) time.Duration {
	return time.Duration(self.Attempts()*delay) * time.Second
}

func (self *BaseWorker) StartTime() time.Time {
	return time.Unix(self.StartStamp, 0)
}

func (self *BaseWorker) NextStartStamp(delay int) int64 {
	return time.Now().Add(self.RetryDelay(delay)).Unix()
}

func (self *BaseWorker) SetStartStamp(stamp int64) {
	self.StartStamp = stamp
}
