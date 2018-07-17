package utils

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

var sampleConfStr string = "{\"AccessKey\": \"2QRHG32RYGWZP08RGGVY\",\"SecretKey\": \"bOWkhWqj0BPc3eRkQsrhwse9JyYikTuuDa5EfHlB\",\"RegionId\": \"cn-north-1\",\"service\": \"CES\"}"

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

	err := os.Remove(path)
	if err != nil {
		fmt.Printf("delete sample file error, %v\n", err)
	}

}

func TestGetAuthorization(t *testing.T) {

	createSampleConfFile()
	defer deleteSampleConfFile()
	s := NewSigner("TEST")
	r, _ := http.NewRequest("GET", "https://ces.cn-north-1.myhwclouds.com/V1.0/5e6f18955f9a452d91205bf1b8911163/favorite-metrics", nil)
	r.Header.Add("X-Sdk-Date", "20170612T194640Z")
	r.Header.Add("Content-Type", "application/json")
	authkey, err := s.GetAuthorization(r)

	if err != nil {
		t.Fatal("generate fail, err:", err.Error())
	} else {
		t.Logf("generate success,%v.", authkey)
	}
}
