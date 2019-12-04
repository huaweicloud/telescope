package task

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//getCSTLocation
func TestGetCSTLocation(t *testing.T) {
	Convey("Test_getCSTLocation", t, func() {
		Convey("test case 1", func() {
			location := getCSTLocation()
			So(location, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.LoadLocation, func(name string) (*time.Location, error) {
				return nil, errors.New("123")
			})
			location := getCSTLocation()
			So(location, ShouldNotBeNil)
		})
	})
}

//isFinalState
func TestUtils_isFinalState(t *testing.T) {

	isFinalState(true, "")
	//STATE_CANCELED
	isFinalState(true, STATE_CANCELED)
	isFinalState(false, STATE_CANCELED)
	Convey("Test_isFinalState", t, func() {
		Convey("test case 1", func() {
			state := isFinalState(true, "")
			So(state, ShouldBeFalse)
		})
		Convey("test case 2", func() {
			state := isFinalState(true, STATE_CANCELED)
			So(state, ShouldBeTrue)
		})
		Convey("test case 3", func() {
			state := isFinalState(false, STATE_CANCELED)
			So(state, ShouldBeTrue)
		})
	})
}
