package manager

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//getchOsSignal
func TestGetchOsSignal(t *testing.T) {
	//go InitConfig()
	Convey("Test_getchOsSignal", t, func() {
		Convey("test case 1", func() {
			signal := getchOsSignal()
			So(signal, ShouldNotBeNil)
		})
	})
}
