package collectors

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//Collect
func TestLoadCollect(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {

			collector := LoadCollector{}
			collect := collector.Collect(1)
			So(collect, ShouldBeNil)
		})
	})
}
