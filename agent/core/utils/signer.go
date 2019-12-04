package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	// COMMA ...
	COMMA          string = ","
	SEMICOLON      string = ";"
	COLON          string = ":"
	EQUALSIGN      string = "="
	BLANK          string = " "
	SLASH          string = "/"
	LINE_SEPARATOR string = "\n"

	BasicDateFormat      = "20060102T150405Z"
	BasicDateFormatShort = "20060102"
	TerminationString    = "sdk_request"
	Algorithm            = "SDK-HMAC-SHA256"
	PreSKString          = "SDK"
	HeaderXDate          = "x-sdk-date"
	HeaderDate           = "date"
	HeaderHost           = "host"
	HeaderAuthorization  = "Authorization"
	HeaderProjectId      = "X-Project-Id"
	HeaderAuthToken      = "X-Auth-Token"

	HeaderConnection  = "connection"
	HeaderUserAgent   = "user-agent"
	HeaderContentType = "content-type"
)

// Signer ...
type Signer struct {
	AccessKey string
	SecretKey string
	AKSKToken string
	Region    string
	Service   string
	reqTime   time.Time
}

// NewSigner ...
func NewSigner(serviceName string) *Signer {
	return &Signer{
		AccessKey: GetConfig().AccessKey,
		SecretKey: GetConfig().SecretKey,
		AKSKToken: GetConfig().AKSKToken,
		Region:    GetConfig().RegionId,
		Service:   serviceName,
	}
}

func (sig *Signer) getReqTime(req *http.Request) (reqTime time.Time, err error) {
	if timeStr := req.Header.Get(HeaderXDate); len(timeStr) > 0 {
		return time.Parse(BasicDateFormat, timeStr)
	} else {
		return time.Time{}, errors.New("No x-sdk-date in request.")
	}
}

// CanonicalRequest ...
func CanonicalRequest(req *http.Request) (string, error) {
	data, err := requestPayload(req)
	if err != nil {
		return "", err
	}
	hexencode, err := hexEncodeSHA256Hash(data)
	if err != nil {
		return "", err
	}
	result := bytes.Buffer{}
	result.WriteString(req.Method)
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(canonicalURI(req))
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(canonicalQueryString(req))
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(canonicalHeaders(req))
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(signedHeaders(req))
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(hexencode)
	return result.String(), nil
}

// RequestPayload
func requestPayload(r *http.Request) ([]byte, error) {

	if r.Body == nil {
		return []byte(""), nil
	}
	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}
	return b, err
}

// hexEncodeSHA256Hash returns hexcode of sha256
func hexEncodeSHA256Hash(body []byte) (string, error) {
	hash := sha256.New()
	if body == nil {
		body = []byte("")
	}
	_, err := hash.Write(body)
	return hex.EncodeToString(hash.Sum(nil)), err
}

// CanonicalURI returns request uri
func canonicalURI(r *http.Request) string {
	pattens := strings.Split(r.URL.Path, SLASH)
	var uri []string
	for _, v := range pattens {
		switch v {
		case "":
			continue
		case ".":
			continue
		case "..":
			if len(uri) > 0 {
				uri = uri[:len(uri)-1]
			}
		default:
			uri = append(uri, url.QueryEscape(v))
		}
	}
	urlpath := SLASH + strings.Join(uri, SLASH)
	urlpath = strings.Replace(urlpath, "+", "%20", -1)
	if ok := strings.HasSuffix(urlpath, SLASH); ok {
		return urlpath
	} else {
		return urlpath + SLASH
	}
}

// CanonicalQueryString
func canonicalQueryString(r *http.Request) string {
	var a []string
	for key, value := range r.URL.Query() {
		k := url.QueryEscape(key)
		for _, v := range value {
			var kv string
			if v == "" {
				kv = k + "="
			} else {
				kv = k + "=" + url.QueryEscape(v)
			}
			a = append(a, strings.Replace(kv, "+", "%20", -1))
		}
	}
	sort.Strings(a)
	return strings.Join(a, "&")
}

// CanonicalHeaders
func canonicalHeaders(r *http.Request) string {
	var a []string
	for key, value := range r.Header {
		sort.Strings(value)
		var q []string
		for _, v := range value {
			q = append(q, trimString(v))
		}
		a = append(a, strings.ToLower(key)+":"+strings.Join(q, ","))
	}
	a = append(a, HeaderHost+":"+r.Host)
	sort.Strings(a)
	return strings.Join(a, "\n") + LINE_SEPARATOR
}

func trimString(s string) string {
	var trimedString []byte
	inQuote := false
	var lastChar byte
	s = strings.TrimSpace(s)
	for _, v := range []byte(s) {
		if byte(v) == byte('"') {
			inQuote = !inQuote
		}
		if lastChar == byte(' ') && byte(v) == byte(' ') && !inQuote {
			continue
		}
		trimedString = append(trimedString, v)
		lastChar = v
	}
	return string(trimedString)
}

// Return the Credential Scope. See http://docs.aws.amazon.com/general/latest/gr/sigv4-create-string-to-sign.html
func credentialScope(t time.Time, regionName, serviceName string) string {
	result := bytes.Buffer{}
	result.WriteString(t.UTC().Format(BasicDateFormatShort))
	result.WriteString(SLASH)
	result.WriteString(regionName)
	result.WriteString(SLASH)
	result.WriteString(serviceName)
	result.WriteString(SLASH)
	result.WriteString(TerminationString)
	return result.String()
}

// Create a "String to Sign". See http://docs.aws.amazon.com/general/latest/gr/sigv4-create-string-to-sign.html
func stringToSign(canonicalRequest, credentialScope string, t time.Time) string {
	hash := sha256.New()
	hash.Write([]byte(canonicalRequest))

	result := bytes.Buffer{}
	result.WriteString(Algorithm)
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(t.UTC().Format(BasicDateFormat))
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(credentialScope)
	result.WriteString(LINE_SEPARATOR)
	result.WriteString(hex.EncodeToString(hash.Sum(nil)))
	return result.String()
}

// Generate a "signing key" to sign the "String To Sign". See http://docs.aws.amazon.com/general/latest/gr/sigv4-calculate-signature.html
func generateSigningKey(secretKey, regionName, serviceName string, t time.Time) ([]byte, error) {

	key := []byte(PreSKString + secretKey)
	var err error
	dateStamp := t.UTC().Format(BasicDateFormatShort)
	data := []string{dateStamp, regionName, serviceName, TerminationString}
	for _, d := range data {
		key, err = hmacsha256(key, d)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

// Create the HWS Signature. See http://docs.aws.amazon.com/general/latest/gr/sigv4-calculate-signature.html
func signStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hm, err := hmacsha256(signingKey, stringToSign)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hm), nil
}

func hmacsha256(key []byte, data string) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(key))
	if _, err := h.Write([]byte(data)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// SignedHeaders
func signedHeaders(r *http.Request) string {
	var a []string
	for key := range r.Header {
		a = append(a, strings.ToLower(key))
	}
	a = append(a, HeaderHost)
	sort.Strings(a)
	return strings.Join(a, SEMICOLON)
}

// AuthHeaderValue Get the finalized value for the "Authorization" header. The signature parameter is the output from SignStringToSign
func AuthHeaderValue(signature, accessKey, credentialScope, signedHeaders string) string {
	result := bytes.Buffer{}
	result.WriteString(Algorithm)
	result.WriteString(" Credential=")
	result.WriteString(accessKey)
	result.WriteString(SLASH)
	result.WriteString(credentialScope)
	result.WriteString(", SignedHeaders=")
	result.WriteString(signedHeaders)
	result.WriteString(", Signature=")
	result.WriteString(signature)
	return result.String()
}

// GetAuthorization ...
func (sig *Signer) GetAuthorization(req *http.Request) (string, error) {
	var authorization string
	if req == nil {
		return authorization, errors.New("Verify failed, req is nil")
	}

	reqTime, err := sig.getReqTime(req)
	if err != nil {
		return authorization, err
	}
	sig.reqTime = reqTime

	canonicalRequest, err := CanonicalRequest(req)
	if err != nil {
		return authorization, err
	}
	credentialScope := credentialScope(sig.reqTime, sig.Region, sig.Service)
	stringToSign := stringToSign(canonicalRequest, credentialScope, sig.reqTime)

	key, err := generateSigningKey(sig.SecretKey, sig.Region, sig.Service, sig.reqTime)
	if err != nil {
		return authorization, err
	}
	signature, err := signStringToSign(stringToSign, key)
	if err != nil {
		return authorization, err
	}

	signedHeaders := signedHeaders(req)
	authValue := AuthHeaderValue(signature, sig.AccessKey, credentialScope, signedHeaders)

	return authValue, nil

}
