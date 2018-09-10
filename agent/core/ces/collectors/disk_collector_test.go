package collectors

import (
	"github.com/shirou/gopsutil/disk"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
)

//Collect
func TestDiskCollect(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {

			collector := DiskCollector{}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
		Convey("test case2", func() {
			//disk.IOCounters()
			mockfn.Replace(disk.IOCounters, func(names ...string) (map[string]disk.IOCountersStat, error) {
				stats := make(map[string]disk.IOCountersStat)
				stats[""] = disk.IOCountersStat{}
				return stats, nil
			})
			collector := DiskCollector{}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
	})
}
