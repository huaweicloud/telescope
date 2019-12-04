package channel

var chCesConfigChan chan string

// Initialize the data channel
func init() {
	chCesConfigChan = make(chan string, 20)
}

// GetCesConfigChan ...
func GetCesConfigChan() chan string {
	return chCesConfigChan
}
