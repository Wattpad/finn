package finn

import (
	"math"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInitialAttempts(t *testing.T) {
	worker := BaseWorker{}
	Convey("Initial attempts should be zero", t, func() {
		So(worker.Attempts(), ShouldEqual, 0)
	})
}

func TestIncreaseAttempts(t *testing.T) {
	worker := BaseWorker{}
	Convey("Increase should add 1 to attempts", t, func() {
		worker.IncreaseAttempts()
		So(worker.Attempts(), ShouldEqual, 1)
	})
}

func TestMaxAttempts(t *testing.T) {
	worker := BaseWorker{}
	Convey("Increase until we reach max attempts", t, func() {
		for i := 0; i < worker.MaxAttempts(); i++ {
			worker.IncreaseAttempts()
		}
		So(worker.Attempts(), ShouldEqual, worker.MaxAttempts())
	})
}

func TestCanRun(t *testing.T) {
	worker := BaseWorker{}
	Convey("Can run by default", t, func() {
		So(worker.CanRun(), ShouldBeTrue)
	})
	Convey("Time in past can run", t, func() {
		worker.SetStartStamp(time.Now().Add(-5000).Unix())
		So(worker.CanRun(), ShouldBeTrue)
	})
	Convey("Time in future can not", t, func() {
		worker.SetStartStamp(time.Now().Add(5 * time.Minute).Unix())
		So(worker.CanRun(), ShouldBeFalse)
	})
}

func TestRunDelay(t *testing.T) {
	worker := BaseWorker{}
	Convey("No delay by default", t, func() {
		So(worker.RunDelay(), ShouldEqual, 0)
	})
	Convey("Should return duration", t, func() {
		worker.SetStartStamp(time.Now().Add(5 * time.Minute).Unix())
		So(math.Ceil(worker.RunDelay().Minutes()), ShouldEqual, 5)
	})
}

func TestRetryDelay(t *testing.T) {
	worker := BaseWorker{}
	Convey("No delay by default", t, func() {
		So(worker.RetryDelay(worker.RetryDelaySeconds()), ShouldEqual, 0)
	})
	Convey("Based on our attempts and RunDelaySeconds", t, func() {
		for i := 0; i < worker.MaxAttempts(); i++ {
			worker.IncreaseAttempts()
			So(worker.RetryDelay(worker.RetryDelaySeconds()).Seconds(), ShouldEqual, worker.Attempts()*worker.RetryDelaySeconds())
		}
	})
}

func TestNextStartStamp(t *testing.T) {
	worker := BaseWorker{}
	Convey("Should be 0 by default", t, func() {
		worker.SetStartStamp(worker.NextStartStamp(worker.RetryDelaySeconds()))
		So(math.Ceil(worker.StartTime().Sub(time.Now()).Seconds()), ShouldEqual, 0)
	})

	Convey("Based on our attempts and RunDelaySeconds", t, func() {
		for i := 0; i < worker.MaxAttempts(); i++ {
			worker.IncreaseAttempts()
			worker.SetStartStamp(worker.NextStartStamp(worker.RetryDelaySeconds()))
			So(math.Ceil(worker.StartTime().Sub(time.Now()).Seconds()), ShouldEqual, worker.Attempts()*worker.RetryDelaySeconds())
		}
	})
}
