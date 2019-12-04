package collectors

import (
	"github.com/shirou/gopsutil/cpu"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//getTotalCPUTime
func TestGetTotalCPUTime(t *testing.T) {
	Convey("Test_NewWindowsLogCollector", t, func() {
		Convey("test case1", func() {
			cpuTimes, _ := cpu.Times(false)
			cpuTime := getTotalCPUTime(cpuTimes[0])
			ShouldNotBeNil(cpuTime)
		})
	})
}

//Collect
func TestCpuCollect(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {

			collector := CPUCollector{}
			collect := collector.Collect(1)
			So(collect, ShouldBeNil)
		})
		Convey("test case2", func() {

			collector := CPUCollector{LastStates: &CPUStates{}}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
	})
}
