package services

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/aggregate"
	"github.com/huaweicloud/telescope/agent/core/ces/collectors"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/process"
)

// StartMetricCollectTask cron job for metric collect
func StartMetricCollectTask(data chan *model.InputMetric, agData chan model.InputMetricSlice) {

	var collectorList []collectors.CollectorInterface

	// simultaneously modify collectorNum in StartAggregateTask when modify the length of collectorList
	collectorList = append(collectorList, &collectors.CPUCollector{})
	collectorList = append(collectorList, &collectors.MemCollector{})
	collectorList = append(collectorList, &collectors.DiskCollector{})
	collectorList = append(collectorList, &collectors.NetCollector{})
	collectorList = append(collectorList, &collectors.LoadCollector{})
	collectorList = append(collectorList, &collectors.ProcStatusCollector{})

	metricSliceArr := make([]model.InputMetricSlice, len(collectorList))

	counter := 0

	nowSecond := time.Now().Second()

	if nowSecond != 0 {
		time.Sleep(time.Duration(59-(nowSecond%60)) * time.Second)
	}

	c := cron.New()

	c.AddFunc("*/10 * * * * *", func() {
		if config.GetConfig().Enable {
			collectTime := time.Now().Unix() * 1000

			allMetric := new(model.InputMetric)
			allMetric.CollectTime = collectTime

			allMetricData := []model.Metric{}

			for i, collector := range collectorList {

				tmp := collector.Collect(collectTime)

				if tmp != nil {
					for _, value := range tmp.Data {
						allMetricData = append(allMetricData, value)
					}
					metricSliceArr[i] = append(metricSliceArr[i], tmp)
				}
			}

			enableProcessList := config.GetConfig().EnableProcessList

			if len(enableProcessList) > 0 {
				processSliceArr := make([]model.InputMetricSlice, len(enableProcessList))

				for j, eachProcess := range enableProcessList {
					pid := eachProcess.Pid
					isExist, _ := process.PidExists(pid)
					if isExist {
						eachProcess, _ := process.NewProcess(pid)
						eachProcessCollector := new(collectors.ProcessCollector)
						eachProcessCollector.Process = eachProcess
						eachRes := eachProcessCollector.Collect(collectTime)
						if eachRes != nil {
							for _, value := range eachRes.Data {
								allMetricData = append(allMetricData, value)
							}
							processSliceArr[j] = append(processSliceArr[j], eachRes)

						}
					}

				}

				newMetricSliceArr := make([]model.InputMetricSlice, len(collectorList)+len(enableProcessList))
				copy(newMetricSliceArr, metricSliceArr)
				copy(newMetricSliceArr[len(collectorList):(len(collectorList)+len(enableProcessList))], processSliceArr)
				metricSliceArr = newMetricSliceArr
			}

			allMetric.Data = allMetricData
			data <- allMetric
			counter++

			if counter == 6 {
				for i, eachMetricSlice := range metricSliceArr {
					agData <- eachMetricSlice
					metricSliceArr[i] = metricSliceArr[i][:0]
				}
				metricSliceArr = metricSliceArr[:len(collectorList)]
				counter = 0
			}
		}
	})

	c.Start()

}

// StartAggregateTask task for aggregate metric in 1 minute
func StartAggregateTask(agRes chan *model.InputMetric, agData chan model.InputMetricSlice) {

	var aggregatorList []aggregate.AggregatorInterface

	aggregatorList = append(aggregatorList, &aggregate.AvgValue{})
	// don't open, first we only support average
	// aggregatorList = append(aggregatorList, &aggregate.MaxValue{})
	// aggregatorList = append(aggregatorList, &aggregate.MinValue{})

	allMetric := new(model.InputMetric)
	tmpData := []model.Metric{}
	count := 0
	collectorNum := 6

	for {
		tmp := <-agData
		for _, aggregator := range aggregatorList {

			eachRes := aggregator.Aggregate(tmp)

			if eachRes != nil {
				for _, value := range eachRes.Data {
					tmpData = append(tmpData, value)
				}

				if allMetric.CollectTime == 0 {
					allMetric.CollectTime = eachRes.CollectTime
				}
			}

		}
		// count length is the collectorNum-1, now the num of collector is 6, and the enabled processes should be considered
		enableProcessList := config.GetConfig().EnableProcessList
		if count < collectorNum+len(enableProcessList)-1 {
			count++
			continue
		} else {

			allMetric.Data = tmpData
			agRes <- allMetric
			tmpData = []model.Metric{}
			count = 0
			allMetric = new(model.InputMetric)
		}

	}

}

// BuildURL build URL string by URI
func BuildURL(destURI string) string {
	var url string
	url = config.GetConfig().Endpoint + "/" + utils.API_CES_VERSION + "/" + utils.GetConfig().ProjectId + destURI
	return url
}

// SendMetricTask task for post metric data
func SendMetricTask(data, agRes chan *model.InputMetric) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport, Timeout: utils.HTTP_CLIENT_TIME_OUT * time.Second}
	for {

		select {
		case metricDataOrigin := <-data:
			logs.GetCesLogger().Debugf("origin data is: %v", *metricDataOrigin)
			go report.SendMetricData(client, BuildURL(cesUtils.PostRawMetricDataURI), metricDataOrigin, false)
		case metricDataAggregate := <-agRes:
			logs.GetCesLogger().Debugf("aggregate data is %v", *metricDataAggregate)
			go report.SendMetricData(client, BuildURL(cesUtils.PostAggregatedMetricDataURI), metricDataAggregate, true)
		}

	}
}
