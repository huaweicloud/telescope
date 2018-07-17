package channel

var chLtsConfigChan chan string

// Initialize the data channel
func init() {
	chLtsConfigChan = make(chan string, 20)
}

// Get the data channel
func GetLtsConfigChan() chan string {
	return chLtsConfigChan
}
