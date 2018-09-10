package collectors

import (
	"github.com/shirou/gopsutil/process"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//ProcCollect
func TestProcessProcCollect(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {
			collector := ProcessCollector{Process: &process.Process{}}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
	})
}
