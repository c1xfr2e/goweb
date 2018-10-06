package fluentd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"

	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
}

// Config fluent config
type Config struct {
	Host          string `json:"host" yaml:"host" env:"FLUENT_HOST"`
	Port          int    `json:"port" yaml:"port" env:"FLUENT_PORT"`
	Tag           string `json:"tag" yaml:"tag" env:"FLUENT_TAG"` // like debug.applicative
	ContainerName string `json:"containername" yaml:"containername" env:"CONTAINER_NAME"`
}

// FluentHook logrus.Hook
type FluentHook struct {
	writer      *fluent.Fluent
	containerID string
	tag         string
	levels      []logrus.Level
}

var _ logrus.Hook = &FluentHook{}

// NewFluentHook new fluent hook
func NewFluentHook(c Config) (*FluentHook, error) {
	if len(c.Host) == 0 {
		return nil, errors.New("invalid fluent host")
	}
	if c.Port == 0 {
		return nil, errors.New("invalid fluent port")
	}
	if len(c.Tag) == 0 {
		return nil, errors.New("invalid fluent tag")
	}

	fd, err := fluent.New(fluent.Config{
		FluentHost:   c.Host,
		FluentPort:   c.Port,
		RetryWait:    1000,
		MaxRetry:     math.MaxInt32,
		AsyncConnect: true,
	})
	if err != nil {
		return nil, err
	}

	cid := os.Getenv("CONTAINER_NAME")
	if len(cid) == 0 {
		cid = c.ContainerName
	}

	return &FluentHook{
		writer:      fd,
		tag:         c.Tag,
		containerID: cid,
		levels:      defaultLevels,
	}, nil
}

// Levels get hook levels
func (h *FluentHook) Levels() []logrus.Level {
	return h.levels
}

// Fire fire
func (h *FluentHook) Fire(entry *logrus.Entry) error {
	b, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}

	data := map[string]string{
		"source":       "logrus_fluentd",
		"log":          string(b),
		"container_id": h.containerID,
	}
	return h.writer.PostWithTime(h.tag, entry.Time, data)
}

func getContainerID() (string, error) {
	b, err := ioutil.ReadFile("/proc/self/cgroup")
	if err != nil {
		return "", err
	}

	re, err := regexp.Compile("/docker/([0-9a-f]{64})")
	if err != nil {
		return "", err
	}
	matchs := re.FindStringSubmatch(string(b))
	if len(matchs) > 0 && len(matchs[1]) > 0 {
		return matchs[1], nil
	}
	return "", fmt.Errorf("can not get container id from %s", string(b))
}
