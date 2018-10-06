package hermogo

import (
	"errors"
	"net/url"
	"strings"

	"github.com/zaoshu/hermogo/mns"
)

// Config mns client configure
type Config struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	AccessID  string `json:"accessid" yaml:"accessid"`
	AccessKey string `json:"accesskey" yaml:"accesskey"`

	accountID string
	region    string
}

var defaultClient mns.Client
var defaultConfig Config

// WaitSeconds wait seconds
var WaitSeconds = mns.DefaultPollingWaitSeconds

// Init init MNS client
func Init(c Config) error {
	// eg: https://xxxxx.mns.cn-beijing.aliyuncs.com/
	// eg: https://xxxxx.mns.cn-beijing-internal.aliyuncs.com/
	// eg: https://xxxxx.mns.cn-beijing-internal-vpc.aliyuncs.com/

	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return err
	}
	hostFields := strings.Split(u.Host, ".")
	if len(hostFields) != 5 {
		return errors.New("invalid mns host: " + c.Endpoint)
	}

	c.accountID = hostFields[0]
	if len(c.accountID) == 0 {
		return errors.New("invalid mns account id")
	}

	regions := strings.Split(hostFields[2], "-")
	if len(regions) < 2 {
		return errors.New("invalid mns region: " + hostFields[2])
	}
	c.region = regions[0] + "-" + regions[1]

	client, err := mns.NewClient(c.Endpoint, c.AccessID, c.AccessKey)
	if err != nil {
		return err
	}
	defaultClient = client
	defaultConfig = c
	return nil
}
