package logdumper

import "os"

var sampleConfStr string = "{\"UserId\": \"1234567\",\"DefaultDestAddr\": \"127.0.0.1\",\"GroupName\": \"myGroup\",\"Topics\": [{\"LogTopicName\": \"shit\",\"Path\":\"\"}]}"

func createSampleConfFile() {
	dir, _ := os.Getwd()
	path := dir + "/conf.json"
	_, err := os.Stat(path)
	var file *os.File

	if os.IsNotExist(err) {
		file, err = os.Create(path)
	}

	defer file.Close()
	file, err = os.OpenFile(path, os.O_RDWR, 0644)
	_, err = file.WriteString(sampleConfStr)
	err = file.Sync()
}

func deleteSampleConfFile() {
	dir, _ := os.Getwd()
	path := dir + "/conf.json"

	os.Remove(path)
}
