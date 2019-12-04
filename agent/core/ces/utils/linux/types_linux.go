package linux

// FSMountStat ...
type FSMountStat struct {
	Partition  string `json:"partition"`
	MountPoint string `json:"mountPoint"`
	State      int64  `json:"state"`
}
