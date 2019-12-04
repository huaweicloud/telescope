package collectors

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//ProcCollect
func TestProcCollect(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {
			collector := ProcStatusCollector{}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
	})
}
