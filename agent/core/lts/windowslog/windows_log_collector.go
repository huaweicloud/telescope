// +build windows

package windowslog

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/lts/model"
	lts_utils "github.com/huaweicloud/telescope/agent/core/lts/utils"
	sys "github.com/elastic/beats/winlogbeat/sys"
	wineventlog "github.com/elastic/beats/winlogbeat/sys/wineventlog"
	"golang.org/x/sys/windows"
)

const CollectorRenderBufferSize = 1024 * 1024 * 2

type WindowsLogCollector struct {
	Query         string
	Channel       string
	Subscription  wineventlog.EvtHandle
	ReadMaxCount  int
	LastRecordId  uint64
	Render        func(event wineventlog.EvtHandle, out io.Writer) error
	RenderBuf     []byte
	OutputBuf     *sys.ByteBuffer
	WindowsOsLogs []model.WindowsSystemEvent
}

func NewWindowsLogCollector(Channel string) (*WindowsLogCollector, error) {
	query, err := wineventlog.Query{
		Log:         Channel,
		IgnoreOlder: time.Millisecond * lts_utils.LOG_File_VALID_DURATION, //收集七天内的
		Level:       "",
		EventID:     "",
		Provider:    []string{},
	}.Build()
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to build query, error: %s", err.Error())
		return nil, err
	}

	c := &WindowsLogCollector{
		Query:        query,
		Channel:      Channel,
		ReadMaxCount: lts_utils.WINDOWS_OS_LOG_PER_COLLECT_MAX_NUMBER,
		RenderBuf:    make([]byte, CollectorRenderBufferSize),
		OutputBuf:    sys.NewByteBuffer(CollectorRenderBufferSize),
	}

	c.Render = func(event wineventlog.EvtHandle, out io.Writer) error {
		return wineventlog.RenderEvent(event, 0, c.RenderBuf, nil, out)
	}
	return c, nil

}

func (c *WindowsLogCollector) Open(recordId uint64) error {
	bm, err := wineventlog.CreateBookmark(c.Channel, recordId)
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to create book mark, error: %s", err.Error())
		return err
	}
	defer wineventlog.Close(bm)
	signalEvent, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		logs.GetLtsLogger().Error("Failed to create event, error: %s", err.Error())
		return nil
	}
	subscriptionHandle, err := wineventlog.Subscribe(
		0, // Session - nil for localhost
		signalEvent,
		"",      // Channel - empty b/c channel is in the query
		c.Query, // Query - nil means all events
		bm,      // Bookmark - for resuming from a specific event
		wineventlog.EvtSubscribeStartAfterBookmark)
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to subscribe, error: %s", err.Error())
		return err
	}

	c.Subscription = subscriptionHandle
	return nil
}

func (c *WindowsLogCollector) Read() []model.WindowsOsLogEventXml {
	handles, _, err := c.eventHandles(c.ReadMaxCount)
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to read windows os log error: %s", err.Error())
		return nil
	}
	if len(handles) == 0 {
		logs.GetLtsLogger().Debug("No newest data return.")
		return nil
	}

	defer func() {
		for _, handle := range handles {
			wineventlog.Close(handle)
		}
	}()

	var events []model.WindowsOsLogEventXml
	for _, h := range handles {
		c.OutputBuf.Reset()
		err := c.Render(h, c.OutputBuf)
		if bufErr, ok := err.(sys.InsufficientBufferError); ok {
			c.RenderBuf = make([]byte, bufErr.RequiredSize)
			c.OutputBuf.Reset()
			err = c.Render(h, c.OutputBuf)
		}
		if err != nil && c.OutputBuf.Len() == 0 {
			continue
		}

		r, err := c.ParseWindowsEventFromXML(c.OutputBuf.Bytes(), err)
		if err != nil {
			continue
		}
		events = append(events, r)
		c.LastRecordId = r.RecordID
	}
	return events
}

func (c *WindowsLogCollector) eventHandles(maxRead int) ([]wineventlog.EvtHandle, int, error) {
	handles, err := wineventlog.EventHandles(c.Subscription, maxRead)
	if err != nil {
		logs.GetLtsLogger().Debugf("Failed to event handle, error: %s", err.Error())
	}

	switch err {
	case nil:
		return handles, maxRead, nil
	case wineventlog.ERROR_NO_MORE_ITEMS:
		return nil, maxRead, nil
	case wineventlog.RPC_S_INVALID_BOUND:
		if err := c.Close(); err != nil {
			logs.GetLtsLogger().Errorf("Failed to close windows log collector, error: %s", err.Error())
			return nil, 0, err
		}
		if err := c.Open(c.LastRecordId); err != nil {
			logs.GetLtsLogger().Errorf("Failed to open windows log collector, error: %s", err.Error())
			return nil, 0, err
		}
		return c.eventHandles(maxRead / 2)
	default:
		return nil, 0, err
	}
}

func (c *WindowsLogCollector) Close() error {
	return wineventlog.Close(c.Subscription)
}

func (c *WindowsLogCollector) ParseWindowsEventFromXML(bytesArr []byte, err error) (model.WindowsOsLogEventXml, error) {
	e, err := sys.UnmarshalEventXML(bytesArr)
	if err != nil {
		return model.WindowsOsLogEventXml{}, fmt.Errorf("Failed to unmarshal XML='%s'. %v", bytesArr, err)
	}

	sys.PopulateAccount(&e.User)

	r := model.WindowsOsLogEventXml{
		Event: e,
	}
	return r, nil
}

func (c *WindowsLogCollector) ReadWindowsOsLogFromRecordId(recordId uint64) error {
	err := c.Open(recordId)
	if err != nil {
		logs.GetLtsLogger().Errorf("Failed to open windows log collector, error: %s", err.Error())
		return err
	}

	eventXmls := c.Read()
	defer c.Close()
	events := make([]model.WindowsSystemEvent, 0)
	for _, eventXml := range eventXmls {
		event := model.WindowsSystemEvent{}
		event.ProviderName = eventXml.Provider.Name
		event.EventSourceName = eventXml.Provider.EventSourceName
		event.EventId = strconv.FormatUint(uint64(eventXml.EventIdentifier.ID), 10)
		event.Version = strconv.FormatUint(uint64(eventXml.Version), 10)
		level, _ := strconv.Atoi(eventXml.Level)
		event.Level = wineventlog.EventLevelToString[wineventlog.EventLevel(level)]
		event.Task = eventXml.Task
		event.Opcode = eventXml.Opcode
		event.TimeCreated = eventXml.TimeCreated.SystemTime.UnixNano() / 1000000
		event.RecordId = strconv.FormatUint(eventXml.RecordID, 10)
		event.ActivityId = eventXml.Correlation.ActivityID
		event.RelatedActivityID = eventXml.Correlation.RelatedActivityID
		event.ProcessId = strconv.FormatUint(uint64(eventXml.Execution.ProcessID), 10)
		event.ThreadId = strconv.FormatUint(uint64(eventXml.Execution.ThreadID), 10)
		event.Channel = eventXml.Channel
		event.HostName = eventXml.Computer
		event.UserId = eventXml.User.Identifier
		events = append(events, event)
	}
	c.WindowsOsLogs = events
	return nil
}
