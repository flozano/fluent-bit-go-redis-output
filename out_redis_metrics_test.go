package main

import (
	"testing"
	"time"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/stretchr/testify/assert"
)

const (
	timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

func TestParseMap(t *testing.T) {

	nestedMap := map[interface{}]interface{}{
		"pod_name":       []byte("test_pod"),
		"container_name": "test_container",
		"annotations": map[interface{}]interface{}{
			"namespace_name":  []byte("test_namespace"),
			"checksum/config": "2e239b0ee49b0803c617dea3",
		},
	}
	pm := parseMap(nestedMap)

	assert.Equal(t, "test_pod", pm["pod_name"])
	assert.Equal(t, "test_container", pm["container_name"])
	assert.Equal(t, "test_namespace",
		pm["annotations"].(map[string]interface{})["namespace_name"])
	assert.Equal(t, "2e239b0ee49b0803c617dea3",
		pm["annotations"].(map[string]interface{})["checksum/config"])
}

type testrecord struct {
	rc   int
	ts   interface{}
	data map[interface{}]interface{}
}

type testFluentPlugin struct {
	hosts    string
	db       string
	records  []testrecord
	position int
	metrics  []*MetricRecord
}

func (p *testFluentPlugin) Environment(ctx unsafe.Pointer, key string) string {
	switch key {
	case "Hosts":
		return p.hosts
	case "Password":
		return "mypasswd"
	case "Key":
		return "testkey"
	case "DB":
		return p.db
	case "UseTLS":
		return "false"
	case "TLSSkipVerify":
		return "false"
	}
	return "unknown-" + key
}

func (p *testFluentPlugin) Unregister(ctx unsafe.Pointer)                                 {}
func (p *testFluentPlugin) NewDecoder(data unsafe.Pointer, length int) *output.FLBDecoder { return nil }
func (p *testFluentPlugin) Exit(code int)                                                 {}
func (p *testFluentPlugin) Send(values []*MetricRecord) error {
	p.metrics = append(p.metrics, values...)
	return nil
}
func (p *testFluentPlugin) GetRecord(dec *output.FLBDecoder) (int, interface{}, map[interface{}]interface{}) {
	if p.position < len(p.records) {
		r := p.records[p.position]
		p.position++
		return r.rc, r.ts, r.data
	}
	return -1, nil, nil
}
func (p *testFluentPlugin) addrecord(rc int, ts interface{}, data map[interface{}]interface{}) {
	p.records = append(p.records, testrecord{rc: rc, ts: ts, data: data})
}

func TestPluginInitialization(t *testing.T) {
	plugin = &testFluentPlugin{hosts: "hosta hostb", db: "0"}
	res := FLBPluginInit(nil)
	assert.Equal(t, output.FLB_OK, res)
	assert.Len(t, rc.pools.pools, 2)
}

func TestPluginInitializationFailure(t *testing.T) {
	plugin = &testFluentPlugin{hosts: "hosta hostb", db: "a"}
	res := FLBPluginInit(nil)
	assert.Equal(t, output.FLB_ERROR, res)
}

func TestPluginFlusher(t *testing.T) {
	testplugin := &testFluentPlugin{hosts: "hosta hostb", db: "0"}
	ts := time.Date(2018, time.February, 10, 10, 11, 12, 0, time.UTC)
	testrecords := map[interface{}]interface{}{
		"a": "bbbccc123",
		"m": "files",
		"v": 10,
	}
	testplugin.addrecord(0, output.FLBTime{Time: ts}, testrecords)
	testplugin.addrecord(0, uint64(ts.Unix()), testrecords)
	testplugin.addrecord(0, 0, testrecords)
	plugin = testplugin
	res := FLBPluginFlush(nil, 0, nil)
	assert.Equal(t, output.FLB_OK, res)
	assert.Len(t, testplugin.metrics, len(testplugin.records))
	assert.Equal(t, testrecords["a"], testplugin.metrics[0].app)
	assert.Equal(t, testrecords["m"], testplugin.metrics[0].name)
	assert.Equal(t, ts.Format("2001-01-01"), testplugin.metrics[0].ts.Format("2001-01-01"))

	assert.Equal(t, testrecords["a"], testplugin.metrics[1].app)
	assert.Equal(t, testrecords["m"], testplugin.metrics[1].name)
	assert.Equal(t, ts.Format("2001-01-01"), testplugin.metrics[1].ts.Format("2001-01-01"))

	assert.Equal(t, testrecords["a"], testplugin.metrics[2].app)
	assert.Equal(t, testrecords["m"], testplugin.metrics[2].name)
	assert.NotEqual(t, ts.Format("2001-01-01"), testplugin.metrics[2].ts.Format("2001-01-01"))
}
