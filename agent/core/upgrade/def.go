package upgrade

// Info ...
type Info struct {
	Version string `json:"version"`
	File    string `json:"file"`
	Size    int    `json:"size"`
	Md5     string `json:"md5"`
}
