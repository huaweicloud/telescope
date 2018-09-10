package collectors

import (
	"github.com/shirou/gopsutil/process"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
)

//SpeProcCountCollector
func TestSpeProcCountCollector(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {
			//process.Processes()
			mockfn.Replace(process.Processes, func() ([]*process.Process, error) {
				stats := []*process.Process{
					{Pid: 1},
				}
				return stats, nil
			})
			defer mockfn.RevertAll()
			//Cmdline
			mockfn.Replace((*process.Process).Cmdline, func(*process.Process) (string, error) {
				return "test", nil
			})
			collector := SpeProcCountCollector{Pname: "test"}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
	})
}
