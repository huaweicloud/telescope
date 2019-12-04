package utils

import (
	"errors"
	"github.com/go-osstat/uptime"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//GetUptimeInSeconds
func Test_GetUptimeInSeconds(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetUptimeInSeconds", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(uptime.Get, func() (time.Duration, error) {
				return time.Second, nil
			})
			i, e := GetUptimeInSeconds()
			So(i, ShouldEqual, 1)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(uptime.Get, func() (time.Duration, error) {
				return time.Second, errors.New("123")
			})
			i, e := GetUptimeInSeconds()
			So(i, ShouldEqual, -1)
			So(e.Error(), ShouldEqual, "123")
		})
	})
}
