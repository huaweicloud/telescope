package collectors

import (
	"github.com/shirou/gopsutil/net"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
)

//Collect
func TestNetCollect(t *testing.T) {
	Convey("Test_Collect", t, func() {
		Convey("test case1", func() {
			collector := NetCollector{}
			collect := collector.Collect(1)
			So(collect, ShouldBeNil)
		})
		Convey("test case2", func() {
			//LastStates
			collector := NetCollector{LastStates: &NetStates{}}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
		Convey("test case3", func() {
			//disk.IOCounters()
			mockfn.Replace(net.IOCounters, func(pernic bool) ([]net.IOCountersStat, error) {
				stats := []net.IOCountersStat{{}}
				return stats, nil
			})
			defer mockfn.RevertAll()
			collector := NetCollector{LastStates: &NetStates{byteRecv: 1, byteSent: 1, uptimeInSeconds: -1}}
			collect := collector.Collect(1)
			So(collect, ShouldNotBeNil)
		})
	})
}
