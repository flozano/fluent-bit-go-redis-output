package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	tf = "2006-01-02 15:04:05.999999999 -0700 MST"
)

func TestNewMetricRecord(t *testing.T) {
	record := make(map[interface{}]interface{})
	record["a"] = "abcd1234"
	record["m"] = "requests"
	record["v"] = int64(125)
	ts, _ := time.Parse(tf, "2006-01-02 15:04:05.999999999 -0700 MST")
	metric, err := NewMetricRecord(ts, record)

	if err != nil {
		assert.Fail(t, "it is not expected that the call to NewMetricRecord fails:%v", err)
	}
	assert.NotNil(t, metric, "metric must not be nil")
	assert.Equal(t, record["a"], metric.app)
	assert.Equal(t, record["m"], metric.name)
	assert.Equal(t, record["v"], metric.value)
	assert.Equal(t, false, metric.discrete)
	assert.Equal(t, "2006-01-02_abcd1234", metric.ToKey())
}
