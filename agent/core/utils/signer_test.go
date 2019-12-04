package utils

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

//NewSigner
func TestNewSigner(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewSigner", t, func() {
		Convey("test case 1", func() {
			signer := NewSigner("")
			So(signer, ShouldNotBeNil)
		})
	})
}

//getReqTime
func TestGetReqTime(t *testing.T) {
	//go InitConfig()
	Convey("Test_getReqTime", t, func() {
		Convey("test case 1", func() {
			signer := &Signer{}
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			reqTime, err := signer.getReqTime(request)
			So(reqTime, ShouldNotBeNil)
			So(err, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			signer := &Signer{}
			request := &http.Request{
			//Header:map[string][]string{"x-sdk-date":{"123"}},
			}
			signer.getReqTime(request)
		})
	})
}

//CanonicalRequest
func TestCanonicalRequest(t *testing.T) {
	//go InitConfig()
	Convey("Test_CanonicalRequest", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(requestPayload, func(r *http.Request) ([]byte, error) {
				return nil, errors.New("123")
			})
			s, e := CanonicalRequest(nil)
			So(s, ShouldBeBlank)
			So(e.Error(), ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(requestPayload, func(r *http.Request) ([]byte, error) {
				return nil, nil
			})
			mockfn.Replace(hexEncodeSHA256Hash, func(body []byte) (string, error) {
				return "", errors.New("123")
			})
			s, e := CanonicalRequest(nil)
			So(s, ShouldBeBlank)
			So(e.Error(), ShouldEqual, "123")
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(requestPayload, func(r *http.Request) ([]byte, error) {
				return nil, nil
			})
			mockfn.Replace(hexEncodeSHA256Hash, func(body []byte) (string, error) {
				return "", nil
			})
			//canonicalURI
			mockfn.Replace(canonicalURI, func(r *http.Request) string {
				return ""
			})
			//canonicalQueryString
			mockfn.Replace(canonicalQueryString, func(r *http.Request) string {
				return ""
			})
			//canonicalHeaders
			mockfn.Replace(canonicalHeaders, func(r *http.Request) string {
				return ""
			})
			//signedHeaders
			mockfn.Replace(signedHeaders, func(r *http.Request) string {
				return ""
			})
			request := &http.Request{}
			s, e := CanonicalRequest(request)
			So(s, ShouldNotBeBlank)
			So(e, ShouldBeNil)
		})
	})
}

//requestPayload
func TestRequestPayload(t *testing.T) {
	//go InitConfig()
	Convey("Test_requestPayload", t, func() {
		Convey("test case 1", func() {
			request := &http.Request{}
			bytes, e := requestPayload(request)
			So(bytes, ShouldBeEmpty)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
				return []byte("123"), nil
			})
			request := &http.Request{Body: ioutil.NopCloser(strings.NewReader("12123"))}
			bytes, e := requestPayload(request)
			So(bytes, ShouldNotBeEmpty)
			So(e, ShouldBeNil)
		})
	})
}

//hexEncodeSHA256Hash
func TestHexEncodeSHA256Hash(t *testing.T) {
	//go InitConfig()
	Convey("Test_hexEncodeSHA256Hash", t, func() {
		Convey("test case 1", func() {
			s, e := hexEncodeSHA256Hash(nil)
			So(s, ShouldNotBeBlank)
			So(e, ShouldBeNil)
		})
	})
}

//canonicalURI
func TestCanonicalURI(t *testing.T) {
	//go InitConfig()
	Convey("Test_canonicalURI", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(strings.HasSuffix, func(s, suffix string) bool {
				return true
			})
			request := &http.Request{
				URL: &url.URL{
					Path: "/1/2/3/4//./../",
				},
			}
			uri := canonicalURI(request)
			So(uri, ShouldNotBeBlank)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(strings.HasSuffix, func(s, suffix string) bool {
				return false
			})
			request := &http.Request{
				URL: &url.URL{
					Path: "/1/2/3/4//./../",
				},
			}
			uri := canonicalURI(request)
			So(uri, ShouldNotBeBlank)
		})
	})
}

//canonicalQueryString
func TestCanonicalQueryString(t *testing.T) {
	//go InitConfig()
	Convey("Test_canonicalQueryString", t, func() {
		Convey("test case 1", func() {
			request := &http.Request{
				URL: &url.URL{
					RawQuery: "123",
				},
			}
			queryString := canonicalQueryString(request)
			So(queryString, ShouldEqual, "123=")
		})
	})
}

//canonicalHeaders
func TestCanonicalHeaders(t *testing.T) {
	//go InitConfig()
	Convey("Test_canonicalHeaders", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			headers := canonicalHeaders(request)
			So(headers, ShouldNotBeBlank)
		})
	})
}

//trimString
func TestTrimString(t *testing.T) {
	//go InitConfig()
	Convey("Test_trimString", t, func() {
		Convey("test case 1", func() {
			s := trimString("  1 2324 /")
			So(s, ShouldNotBeBlank)
		})
	})
}

//credentialScope
func TestCredentialScope(t *testing.T) {
	//go InitConfig()
	Convey("Test_credentialScope", t, func() {
		Convey("test case 1", func() {
			scope := credentialScope(time.Now(), "", "")
			So(scope, ShouldNotBeBlank)
		})
	})
}

//stringToSign
func TestStringToSign(t *testing.T) {
	//go InitConfig()
	Convey("Test_stringToSign", t, func() {
		Convey("test case 1", func() {
			sign := stringToSign("", "", time.Now())
			So(sign, ShouldNotBeBlank)
		})
	})
}

//generateSigningKey
func TestGenerateSigningKey(t *testing.T) {
	//go InitConfig()
	Convey("Test_generateSigningKey", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(hmacsha256, func(key []byte, data string) ([]byte, error) {
				return nil, errors.New("123")
			})
			bytes, e := generateSigningKey("", "", "", time.Now())
			So(bytes, ShouldBeEmpty)
			So(e.Error(), ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(hmacsha256, func(key []byte, data string) ([]byte, error) {
				return nil, nil
			})
			bytes, e := generateSigningKey("", "", "", time.Now())
			So(bytes, ShouldBeEmpty)
			So(e, ShouldBeNil)
		})
	})
}

//signStringToSign
func TestSignStringToSign(t *testing.T) {
	//go InitConfig()
	Convey("Test_signStringToSign", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(hmacsha256, func(key []byte, data string) ([]byte, error) {
				return nil, errors.New("123")
			})
			bytes, e := signStringToSign("", nil)
			So(bytes, ShouldBeEmpty)
			So(e.Error(), ShouldEqual, "123")
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(hmacsha256, func(key []byte, data string) ([]byte, error) {
				return nil, nil
			})
			bytes, e := signStringToSign("", nil)
			So(bytes, ShouldBeEmpty)
			So(e, ShouldBeNil)
		})
	})
}

//hmacsha256
func TestHmacsha256(t *testing.T) {
	//go InitConfig()
	Convey("Test_hmacsha256", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(hash.Hash.Write, func(h hash.Hash, p []byte) (n int, err error) {
				return 0, errors.New("123")
			})
			bytes, e := hmacsha256(nil, "")
			So(bytes, ShouldNotBeEmpty)
			So(e, ShouldBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(hash.Hash.Write, func(h hash.Hash, p []byte) (n int, err error) {
				return 0, nil
			})
			bytes, e := hmacsha256(nil, "")
			So(bytes, ShouldNotBeEmpty)
			So(e, ShouldBeNil)
		})
	})
}

//signedHeaders
func TestSignedHeaders(t *testing.T) {
	//go InitConfig()
	Convey("Test_signedHeaders", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			headers := signedHeaders(request)
			So(headers, ShouldNotBeBlank)
		})
	})
}

//AuthHeaderValue
func TestAuthHeaderValue(t *testing.T) {
	//go InitConfig()
	Convey("Test_AuthHeaderValue", t, func() {
		Convey("test case 1", func() {
			value := AuthHeaderValue("", "", "", "")
			So(value, ShouldNotBeBlank)
		})
	})
}

//GetAuthorization
func TestGetAuthorization(t *testing.T) {
	//go InitConfig()
	Convey("Test_GetAuthorization", t, func() {
		Convey("test case 1", func() {
			signer := &Signer{}
			s, e := signer.GetAuthorization(nil)
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 2", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*Signer).getReqTime, func(s *Signer, req *http.Request) (reqTime time.Time, err error) {
				return time.Now(), errors.New("123")
			})
			signer := &Signer{}
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			s, e := signer.GetAuthorization(request)
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 3", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*Signer).getReqTime, func(s *Signer, req *http.Request) (reqTime time.Time, err error) {
				return time.Now(), nil
			})
			//CanonicalRequest
			mockfn.Replace(CanonicalRequest, func(req *http.Request) (string, error) {
				return "", errors.New("123")
			})
			signer := &Signer{}
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			s, e := signer.GetAuthorization(request)
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 4", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*Signer).getReqTime, func(s *Signer, req *http.Request) (reqTime time.Time, err error) {
				return time.Now(), nil
			})
			//CanonicalRequest
			mockfn.Replace(CanonicalRequest, func(req *http.Request) (string, error) {
				return "", nil
			})
			//generateSigningKey
			mockfn.Replace(generateSigningKey, func(secretKey, regionName, serviceName string, t time.Time) ([]byte, error) {
				return nil, errors.New("123")
			})
			signer := &Signer{}
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			s, e := signer.GetAuthorization(request)
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 5", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*Signer).getReqTime, func(s *Signer, req *http.Request) (reqTime time.Time, err error) {
				return time.Now(), nil
			})
			//CanonicalRequest
			mockfn.Replace(CanonicalRequest, func(req *http.Request) (string, error) {
				return "", nil
			})
			//generateSigningKey
			mockfn.Replace(generateSigningKey, func(secretKey, regionName, serviceName string, t time.Time) ([]byte, error) {
				return nil, nil
			})
			//signStringToSign
			mockfn.Replace(signStringToSign, func(stringToSign string, signingKey []byte) (string, error) {
				return "", errors.New("123")
			})
			signer := &Signer{}
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			s, e := signer.GetAuthorization(request)
			So(s, ShouldBeBlank)
			So(e, ShouldNotBeNil)
		})
		Convey("test case 6", func() {
			defer mockfn.RevertAll()
			mockfn.Replace((*Signer).getReqTime, func(s *Signer, req *http.Request) (reqTime time.Time, err error) {
				return time.Now(), nil
			})
			//CanonicalRequest
			mockfn.Replace(CanonicalRequest, func(req *http.Request) (string, error) {
				return "", nil
			})
			//generateSigningKey
			mockfn.Replace(generateSigningKey, func(secretKey, regionName, serviceName string, t time.Time) ([]byte, error) {
				return nil, nil
			})
			//signStringToSign
			mockfn.Replace(signStringToSign, func(stringToSign string, signingKey []byte) (string, error) {
				return "", nil
			})
			signer := &Signer{}
			request := &http.Request{
				Header: map[string][]string{"x-sdk-date": {"123"}},
			}
			s, e := signer.GetAuthorization(request)
			So(s, ShouldNotBeBlank)
			So(e, ShouldBeNil)
		})
	})
}

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

func TestGetAuthorization1(t *testing.T) {

	/*createSampleConfFile()
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
	}*/
}
