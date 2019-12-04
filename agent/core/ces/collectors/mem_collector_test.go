package collectors

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//Collect
func TestMemCollect(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {

			collector := MemCollector{}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
	})
}
