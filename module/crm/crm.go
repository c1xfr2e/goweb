//  Created by paincompiler on 28/01/2018

package crm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"time"

	"errors"

	"github.com/allegro/bigcache"
	"github.com/sirupsen/logrus"
)

const (
	// ErrorAccessDeniedTooFrequently 请求频繁
	ErrorAccessDeniedTooFrequently = 110006
	// ErrorInvalidAccessToken 无效的访问token
	ErrorInvalidAccessToken = 20000002

	tokenCacheKey = "xiaoshouyi_token"
)

// Xiaoshouyi Xiaoshouyi storage
type Xiaoshouyi struct {
	ClientID     string `yaml:"clientid"`
	ClientSecret string `yaml:"clientsecret"`
	RedirectURI  string `yaml:"redirecturi"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`

	GetTokenURL string `yaml:"gettokenurl"`
	Expire      int64  `yaml:"expire"`

	CreateLeadURL string `yaml:"createleadurl"`

	LeadSourceID int64 `yaml:"leadsourceid"`
	HighSeaID    int64 `yaml:"highseaid"`
}

// GConfig for global config
var GConfig Xiaoshouyi
var cache *bigcache.BigCache

// Init for initialization of crm config
func Init(xiaoshouyi Xiaoshouyi) {
	GConfig = xiaoshouyi
}

// LeadObject 销售线索对象
type LeadObject struct {
	Public bool `json:"public"`
	Record struct {
		Name         string `json:"name"`
		CompanyName  string `json:"companyName"`
		Mobile       string `json:"mobile"`
		Email        string `json:"email"`
		WeChatC      string `json:"dbcVarchar1"`
		DataSourceC  string `json:"dbcVarchar2"`
		LeadSourceID int64  `json:"leadSourceId"`
		HighSeaID    int64  `json:"highSeaId"`
		Comment      string `json:"comment"`
	} `json:"record"`
}

// NewLeadObject new lead object
func NewLeadObject() *LeadObject {
	object := &LeadObject{
		Public: true,
	}
	object.Record.LeadSourceID = GConfig.LeadSourceID
	return object
}

func getTokenFromAPI() (token string, issuedAt int64, err error) {
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("client_id", GConfig.ClientID)
	form.Add("client_secret", GConfig.ClientSecret)
	form.Add("redirect_uri", GConfig.RedirectURI)
	form.Add("username", GConfig.Username)
	form.Add("password", GConfig.Password)
	req, err := http.NewRequest(http.MethodPost, GConfig.GetTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	ret := struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		Token            string `json:"access_token"`
		IssuedAt         int64  `json:"issued_at"`
	}{}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return
	}

	if len(ret.Error) > 0 {
		err = errors.New(fmt.Sprintf("%s:%s", ret.Error, ret.ErrorDescription))
		return
	}
	return ret.Token, ret.IssuedAt, nil
}

// GetToken get token
func GetToken(renew bool) (string, error) {
	if !renew && cache != nil {
		token, err := cache.Get(tokenCacheKey)
		if err == nil {
			return string(token), nil
		} else {
			logrus.Error("get xiaoshouyi token from cache failed for: ", err)
		}
	}
	var issuedAt int64
	token, issuedAt, err := getTokenFromAPI()
	if err != nil {
		return "", err
	}
	expire := GConfig.Expire - (time.Now().Unix() - issuedAt/1000)
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Duration(expire) * time.Second))
	if err != nil {
		return "", err
	}
	cache.Set(tokenCacheKey, []byte(token))
	if err != nil {
		return "", err
	}
	return token, nil
}

// Create create lead object
func (l *LeadObject) Create() error {
	token, err := GetToken(false)
	if err != nil {
		return err
	}

	bodyReq, err := json.Marshal(l)
	if err != nil {
		return err
	}

	ret := struct {
		ID      int64  `json:"id"`
		Error   int64  `json:"error_code"`
		Message string `json:"message"`
	}{}
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest(http.MethodPost, GConfig.CreateLeadURL, strings.NewReader(string(bodyReq)))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		bodyRet, err := ioutil.ReadAll(res.Body)
		if err != nil {
			res.Body.Close()
			return err
		}
		res.Body.Close()

		err = json.Unmarshal(bodyRet, &ret)
		if err != nil {
			return err
		}
		switch ret.Error {
		case ErrorAccessDeniedTooFrequently:
			// fixme add queue handler
			time.Sleep(time.Second)
			continue
		case ErrorInvalidAccessToken:
			token, err = GetToken(true)
			if err != nil {
				return err
			}
			continue
		}
		if ret.Error != 0 {
			return errors.New(fmt.Sprintf("%d:%s", ret.Error, ret.Message))
		}
		return nil
	}
	return errors.New(fmt.Sprintf("%d:%s", ret.Error, ret.Message))
}
