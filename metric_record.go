package main

import (
	"errors"
	"github.com/spf13/cast"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

type MetricRecord struct {
	ts       time.Time
	app      string
	name     string
	value    int64
	discrete bool
}

func NewMetricRecord(ts time.Time, rec map[interface{}]interface{}) (*MetricRecord, error) {
	var result MetricRecord
	var err error
	var exists bool
	var val interface{}

	result.ts, err = cast.ToTimeE(ts)
	if err != nil {
		return nil, err
	}

	val, exists = rec["a"]
	if !exists {
		return nil, errors.New("`a` missing from record")
	}
	result.app, err = cast.ToStringE(val)
	if err != nil {
		return nil, err
	}

	val, exists = rec["m"]
	if !exists {
		return nil, errors.New("`m` missing from record")
	}
	result.name, err = cast.ToStringE(val)
	if err != nil {
		return nil, err
	}

	val, exists = rec["d"]
	if exists {
		result.discrete = true
		result.value, err = cast.ToInt64E(val)
		if err != nil {
			return nil, err
		}
	} else {
		val, exists = rec["v"]
		if !exists {
			return nil, errors.New("both `d` and `v` are missing")
		}
		result.discrete = false
		result.value, err = cast.ToInt64E(val)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (mr *MetricRecord) ToKey() string {
	return mr.ts.Format(layoutISO) + "_" + mr.app
}
