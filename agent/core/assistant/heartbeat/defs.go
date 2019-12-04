package heartbeat

import "time"

// Heartbeat ...
type Heartbeat struct {
	InstanceID string      `json:"instance_id"`
	Version    string      `json:"version,omitempty"`
	User       string      `json:"user,omitempty"`
	PID        int         `json:"pid,omitempty"`
	Status     string      `json:"status"`
	Metric     *Metric     `json:"metric,omitempty"`
	Extension  []Extension `json:"extension,omitempty"`
	EnvInfo    *Env        `json:"env_info,omitempty"`
}

// HBResp ...
type HBResp struct {
	Config *HBConfig `json:"config,omitempty"`
}

// HBConfig ...
type HBConfig struct {
	AssistSwitch bool `json:"assist_switch,omitempty"`
}

// TODO cpu -> cpu_usage, mem -> mem_usage, port -> port_used; make it detailed; 尽量命名准确
// Metric ...
type Metric struct {
	cpu        string
	mem        string
	port       string
	open_files string
	net_in     string
	net_out    string
	io         string
}

// TODO 下一期再说
// Extension ...
type Extension struct {
	Name    string
	Version string
	User    string
	PID     string
	Status  string
	Metric  string
	Meta    string
}

// Env ...
type Env struct {
	Hostname      string `json:"hostname"`
	IP            string `json:"ip"`
	Arch          string `json:"arch"`
	OS            string `json:"os"`
	Distro        string `json:"distro,omitempty"`
	DistroVersion string `json:"distro_version,omitempty"`
}

// HeartbeatState ...
type HeartbeatState struct {
	LastUpdateTime          time.Time
	LastUpdateSucceededTime time.Time
	IsEnvReported           bool
}

const (
	// AGENT_STATUS_RUNNING ...
	AGENT_STATUS_RUNNING = "RUNNING"
)

var (
	// HBS ...
	HBS HeartbeatState
)
