package services

import (
	"sync"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/aggregate"
	"github.com/huaweicloud/telescope/agent/core/ces/collectors"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"github.com/shirou/gopsutil/process"
)

// simultaneously modify collectorNum in aggregateMetric when modify the length of collectorList
var (
	collectorList = []collectors.CollectorInterface{
		&collectors.CPUCollector{},
		&collectors.MemCollector{},
		&collectors.DiskCollector{},
		&collectors.NetCollector{},
		&collectors.LoadCollector{},
		&collectors.ProcStatusCollector{},
	}
)

// StartMetricCollectTask cron job for metric collect
func StartMetricCollectTask() {
	time.Sleep(time.Duration(5) * time.Second)
	counter := 0
	// set default
	cronTime := utils.DisableDetailDataCronJobTimeSecond
	// just send one during a cron interval
	sendCount := 1

	// check if detail metric is enable
	if utils.GetConfig().DetailMonitorEnable {
		logs.GetCesLogger().Infof("Detail data monitor is enabled.")
		sendCount = cesUtils.SendTotal
		cronTime = utils.DetailDataCronJobTimeSecond
	}
	ticker := time.NewTicker(time.Duration(cronTime) * time.Second)
	// metricClass store different type of metric
	collectorNum := len(collectorList)
	var enableProcesses = config.GetConfig().EnableProcessList
	var metricClass = make([]model.InputMetricSlice, collectorNum+len(enableProcesses))
	for range ticker.C {
		if !config.GetConfig().Enable {
			continue
		}

		// allMetricData store all collect metrics
		var allMetricData []model.Metric
		now := utils.GetCurrTSInMs()
		for i, collector := range collectorList {
			// collect basic 6 system metrics
			metric := collector.Collect(now)
			if metric != nil {
				allMetricData = append(allMetricData, metric.Data...)
			}
			metricClass[i] = append(metricClass[i], metric)
		}

		// collect customized processes metrics
		if len(enableProcesses) > 0 {
			wg := &sync.WaitGroup{}
			wg.Add(len(enableProcesses))
			for j, proc := range enableProcesses {
				go func(hb config.HbProcess, j int) {
					defer wg.Done()

					pid := hb.Pid
					exist, err := process.PidExists(pid)
					if nil != err || !exist {
						return
					}
					p, err := process.NewProcess(pid)
					if nil != err {
						return
					}
					metricChan := make(chan *model.InputMetric, 1)
					go func() {
						pc := &collectors.ProcessCollector{Process: p}
						metricChan <- pc.Collect(now)
					}()
					select {
					case metric := <-metricChan:
						if metric != nil {
							allMetricData = append(allMetricData, metric.Data...)
						}
						metricClass[collectorNum+j] = append(metricClass[collectorNum+j], metric)
					case <-time.After(3 * time.Second):
						logs.GetCesLogger().Errorf("collect processes metrics timeout(PID:%d)", pid)
					}
				}(proc, j)
			}
			wg.Wait()
		}
		allMetric := &model.InputMetric{
			CollectTime: now,
			Data:        allMetricData,
		}
		logs.GetCesLogger().Debugf("origin data is: %v", allMetric)

		if utils.GetConfig().DetailMonitorEnable {
			logs.GetCesLogger().Debugf("Begin to send detail metric data.")
			go report.SendMetricData(BuildURL(cesUtils.PostRawMetricDataURI), allMetric, false)
		}
		if counter++; counter < sendCount {
			continue
		}
		// reset timer counter per minute
		counter = 0
		// reset processes per minute
		enableProcesses = config.GetConfig().EnableProcessList
		// collect process count by input process command line argument
		processCmdlines := config.GetConfig().SpecifiedProcList
		logs.GetCesLogger().Debugf("process cmdline are %v", processCmdlines)
		if lineCount := len(processCmdlines); lineCount > 0 {
			spc := &collectors.SpeProcCountCollector{CmdLines: processCmdlines}
			metric := spc.Collect(now)
			if nil != metric {
				for _, m := range metric.Data {
					metricClass = append(metricClass, model.InputMetricSlice{
						&model.InputMetric{
							Data:        []model.Metric{m},
							Type:        "cmdline",
							CollectTime: now,
						},
					})
				}
			}
		}
		go aggregateMetric(metricClass)
		// reset metricClass after report data per minute
		metricClass = make([]model.InputMetricSlice, collectorNum+len(enableProcesses))
	}
}

// aggregateMetric aggregate metric per minute
func aggregateMetric(metricClass []model.InputMetricSlice) {
	var aggregatorList []aggregate.AggregatorInterface
	aggregatorList = append(aggregatorList, &aggregate.AvgValue{})
	// TODO only support average now, need configuration to enable max/min function
	// aggregatorList = append(aggregatorList, &aggregate.MaxValue{})
	// aggregatorList = append(aggregatorList, &aggregate.MinValue{})

	allMetric := new(model.InputMetric)
	for _, class := range metricClass {
		if nil == class || len(class) == 0 {
			continue
		}
		for _, aggregator := range aggregatorList {
			metric := aggregator.Aggregate(class)
			if metric == nil {
				continue
			}
			for _, value := range metric.Data {
				allMetric.Data = append(allMetric.Data, value)
			}

			if allMetric.CollectTime == 0 {
				allMetric.CollectTime = metric.CollectTime
			}
		}
	}
	logs.GetCesLogger().Debugf("aggregate data is %v", allMetric)
	report.SendMetricData(BuildURL(cesUtils.PostAggregatedMetricDataURI), allMetric, true)
}

// BuildURL build URL string by URI
func BuildURL(destURI string) string {
	var url string
	url = config.GetConfig().Endpoint + utils.SLASH + utils.APICESVersion + utils.SLASH + utils.GetConfig().ProjectId + destURI
	return url
}
