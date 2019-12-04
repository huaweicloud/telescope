package circuitbreaker

import (
	"container/list"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

//getCurrentUser
func TestDelOutdatedState(t *testing.T) {
	Convey("TestDelOutdatedState", t, func() {
		Convey("test case 1", func() {
			cb = &CircuitBreaker{
				StateList: list.New(),
			}
			now := time.Now()
			for i := 1; i <= 5; i++ {
				cb.StateList.PushBack(State{
					CPUPct: 99.01,
					Memory: 123,
					Time:   now.Add(-time.Duration(i) * time.Minute),
				})
			}
			delOutdatedState(now.Add(-63 * time.Second))
			So(cb.StateList.Len(), ShouldEqual, 1)
		})
		Convey("test case 2", func() {
			cb = &CircuitBreaker{
				StateList: list.New(),
			}
			now := time.Now()
			for i := 5; i >= 1; i-- {
				cb.StateList.PushBack(State{
					CPUPct: 99.01,
					Memory: 123,
					Time:   now.Add(-time.Duration(i) * time.Minute),
				})
			}
			delOutdatedState(now.Add(-63 * time.Second))
			So(cb.StateList.Len(), ShouldEqual, 1)
		})
		Convey("test case 3", func() {
			cb = &CircuitBreaker{
				StateList: list.New(),
			}
			now := time.Now()
			for i := 5; i >= 1; i-- {
				cb.StateList.PushBack(State{
					CPUPct: 99.01,
					Memory: 123,
					Time:   now.Add(-time.Duration(i) * time.Minute),
				})
			}
			delOutdatedState(now)
			So(cb.StateList.Len(), ShouldEqual, 0)
		})
		Convey("test case 4", func() {
			cb = &CircuitBreaker{
				StateList: list.New(),
			}
			now := time.Now()
			for i := 5; i >= 1; i-- {
				cb.StateList.PushBack(State{
					CPUPct: 99.01,
					Memory: 123,
					Time:   now.Add(-time.Duration(i) * time.Minute),
				})
			}
			cb.StateList.PushBack(State{CPUPct: 0, Memory: 0, Time: now})
			delOutdatedState(now)
			So(cb.StateList.Len(), ShouldEqual, 1)
		})
	})
}
