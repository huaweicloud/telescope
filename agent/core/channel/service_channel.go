package channel

type HBServiceData struct {
	Service string `json:"service"`
	Detail  string `json:"detail"`
}

var servicesChData chan HBServiceData

// Initialize the data channel
func init() {
	servicesChData = make(chan HBServiceData, 20)
}

// Get the data channel
func GetServicesChData() chan HBServiceData {
	return servicesChData
}
