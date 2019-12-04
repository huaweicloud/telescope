package model

import (
	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"regexp"
	"testing"
)

//BuildMetric
func TestBuildMetric(t *testing.T) {
	//go InitConfig()
	Convey("Test_BuildMetric", t, func() {
		Convey("test case 1", func() {
			metric := BuildMetric(0, nil)
			So(metric, ShouldNotBeNil)
		})
	})
}

//BuildCesMetricData
func TestBuildCesMetricData(t *testing.T) {
	//go InitConfig()
	Convey("Test_BuildCesMetricData", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{ExternalService: ces_utils.ExternalServiceBMS}
			})
			//AliasMetricName
			mockfn.Replace(AliasMetricName, func(metricName string) string {
				return "123"
			})
			metric := &InputMetric{
				Data: []Metric{
					{
						MetricName:   "name",
						MetricPrefix: "volumeSlAsH",
					},
				},
				CollectTime: 123,
			}
			data := BuildCesMetricData(metric, true)
			So(data, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{ExternalService: ces_utils.ExternalServiceBMS}
			})
			metric := &InputMetric{
				Data: []Metric{
					{MetricName: "name",
						MetricPrefix: "123"},
				},
				CollectTime: 123,
			}
			data := BuildCesMetricData(metric, false)
			So(data, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(utils.GetConfig, func() *utils.GeneralConfig {
				return &utils.GeneralConfig{}
			})
			mockfn.Replace(AliasMetricName, func(metricName string) string {
				return "123"
			})
			metric := &InputMetric{
				Data: []Metric{
					{MetricName: "name"},
				},
				CollectTime: 123,
			}
			data := BuildCesMetricData(metric, false)
			So(data, ShouldNotBeNil)
		})
	})
}

//getOldMetricName
func TestGetOldMetricName(t *testing.T) {
	//go InitConfig()
	Convey("Test_getOldMetricName", t, func() {
		Convey("test case 1", func() {
			name := getOldMetricName("", "")
			So(name, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			name := getOldMetricName("disk_123", "m")
			So(name, ShouldEqual, "m_disk_123")
		})
		Convey("test case 3", func() {
			name := getOldMetricName("proc_123", "m")
			So(name, ShouldEqual, "proc_m_123")
		})
		Convey("test case 4", func() {
			name := getOldMetricName("gpu_123", "m")
			So(name, ShouldEqual, "slotm_gpu_123")
		})
		Convey("test case 5", func() {
			name := getOldMetricName("1proc_123_device", "mdmm")
			So(name, ShouldEqual, "mdmm_1proc_123_device")
		})
		Convey("test case 6", func() {
			name := getOldMetricName("1proc_123_device1", "1mdmm")
			So(name, ShouldBeBlank)
		})
	})
}

//setOldMetricData
func TestSetOldMetricData(t *testing.T) {
	//go InitConfig()
	Convey("Test_setOldMetricData", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(getOldMetricName, func(metricName, MetricPrefix string) string {
				return "123"
			})
			mockfn.Replace(AliasMetricName, func(metricName string) string {
				return "123"
			})
			arrs := []CesMetricData{}
			data := setOldMetricData(arrs, CesMetricData{}, Metric{})
			So(data, ShouldNotBeNil)
		})
	})
}

//AliasMetricName
func TestAliasMetricName(t *testing.T) {
	//go InitConfig()
	Convey("Test_AliasMetricName", t, func() {
		Convey("test case 1", func() {
			name := AliasMetricName("123")
			So(name, ShouldBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(regexp.MatchString, func(pattern string, s string) (matched bool, err error) {
				return true, nil
			})
			name := AliasMetricName("123")
			So(name, ShouldNotBeBlank)
		})
		Convey("test case 3", func() {
			name := AliasMetricName("abcdefghjkabcdefghjkabcdefghjkabcdefghjkabcdefghjkabcdefghjkabcdefghjk")
			So(name, ShouldNotBeBlank)
		})
	})
}

//generateHashID
func TestGenerateHashID(t *testing.T) {
	//go InitConfig()
	Convey("Test_generateHashID", t, func() {
		Convey("test case 1", func() {
			id := generateHashID("")
			So(id, ShouldNotBeBlank)
		})
	})
}

//getUnitByMetric
func Test_getUnitByMetric(t *testing.T) {
	//go InitConfig()
	Convey("Test_getUnitByMetric", t, func() {
		Convey("test case 1", func() {
			metric := getUnitByMetric("")
			So(metric, ShouldNotBeNil)
		})
	})
}

//GetMountPrefix
func TestGetMountPrefix(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetMountPrefix", t, func() {
		Convey("test case 1", func() {
			prefix := GetMountPrefix("")
			So(prefix, ShouldNotBeBlank)
		})
	})
}
