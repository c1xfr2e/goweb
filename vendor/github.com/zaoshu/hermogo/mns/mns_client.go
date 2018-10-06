package mns

import (
	"bytes"
	"crypto/hmac"
	cmd5 "crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/zaoshu/hermogo/proto"
)

// MNSClient mns client
type MNSClient struct {
	Host      string
	Version   string
	AccessID  string
	AccessKey string
}

var _ Client = MNSClient{}

// NewClient new mns client
func NewClient(host, accessID, accessKey string) (Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	if len(u.Scheme) == 0 {
		u.Scheme = "http"
	}
	if strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path[0 : len(u.Path)-1]
	}

	return MNSClient{u.String(), APIVersion, accessID, accessKey}, nil
}

func (c MNSClient) getHeader(method, path string, data []byte, header map[string]string) map[string]string {
	v := map[string]string{
		"Content-MD5":   "",
		"Content-Type":  RequestContentType,
		"Date":          time.Now().UTC().Format(http.TimeFormat),
		"x-mns-version": c.Version,
	}

	for key, value := range header {
		v[key] = value
	}
	if data != nil {
		v["Content-MD5"] = c.getMD5String(data)
	}
	v["Authorization"] = fmt.Sprintf(
		"MNS %s:%s",
		c.AccessID,
		c.getSignature(method, path, v),
	)
	return v
}

func (c MNSClient) getMD5String(data []byte) string {
	h := cmd5.New()
	_, err := h.Write(data)
	if err != nil {
		return ""
	}
	digest := h.Sum(nil)
	hexDigest := hex.EncodeToString(digest)
	return base64.StdEncoding.EncodeToString([]byte(hexDigest))
}

// see https://help.aliyun.com/document_detail/27487.html
func (c MNSClient) getSignature(method, path string, header map[string]string) string {
	values := make([]string, 0, 8)
	values = append(values, method)
	values = append(values, header["Content-MD5"])
	values = append(values, header["Content-Type"])
	values = append(values, header["Date"])
	for _, k := range c.getSortedMNSHeader(header) {
		values = append(values, k+":"+header[k])
	}
	values = append(values, path)
	mac := hmac.New(sha1.New, []byte(c.AccessKey))
	_, err := mac.Write([]byte(strings.Join(values, "\n")))
	if err != nil {
		return ""
	}
	checksum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(checksum)
}

func (c MNSClient) getSortedMNSHeader(header map[string]string) []string {
	sorted := []string{}
	for k := range header {
		if strings.HasPrefix(k, "x-mns-") {
			sorted = append(sorted, k)
		}
	}
	sort.Strings(sorted)
	return sorted
}

func (c MNSClient) do(method, path string, data []byte, header map[string]string) (*http.Response, error) {
	var body io.Reader
	if data != nil {
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.Host+path, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.getHeader(method, path, data, header) {
		req.Header.Set(k, v)
	}

	return http.DefaultClient.Do(req)
}

// DoRequest send request
func (c MNSClient) DoRequest(method string, path string, header map[string]string, req interface{}, resp interface{}) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid path %s, must start with /", path)
	}

	var (
		body []byte
		err  error
	)

	if req != nil {
		body, err = xml.Marshal(req)
		if err != nil {
			return err
		}
	}

	httpResp, err := c.do(method, path, body, header)
	if err != nil {
		return err
	}
	defer func() {
		err2 := httpResp.Body.Close()
		if err2 != nil {
			return
		}
	}()

	b, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	if httpResp.StatusCode >= 200 && httpResp.StatusCode < 300 {
		if len(b) == 0 || resp == nil {
			return nil
		}
		return proto.UnmarshalFromXML(b, resp)
	} else if httpResp.StatusCode >= 400 && httpResp.StatusCode < 600 {
		err = proto.UnmarshalFromXML(b, resp)
		if err != nil {
			return proto.NewErrorFromXML(b)
		}
		return nil
	} else {
		return fmt.Errorf("can not handle http status %d, [%s]", httpResp.StatusCode, string(b))
	}
}
