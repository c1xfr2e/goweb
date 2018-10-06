package logging

import "github.com/sirupsen/logrus"

// NewStatsField stats values
func NewStatsField(ns string, value map[string]interface{}) logrus.Fields {
	if len(ns) == 0 || len(value) == 0 {
		return logrus.Fields{}
	}
	return logrus.Fields{
		"stats": map[string]interface{}{
			ns: value,
		},
	}
}
