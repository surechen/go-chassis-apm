package apm_test

import (
	"context"
	"github.com/go-chassis/go-chassis-apm"
	"github.com/go-chassis/go-chassis-apm/tracing"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	op tracing.TracingOptions
	sc *tracing.SpanContext
)

//InitOption
func InitOption() {
	op = tracing.TracingOptions{
		APMName:        "testclient",
		ServerURI:      "192.168.88.64:8080",
		MicServiceName: "mesher",
		MicServiceType: 1}
}

//InitSpanContext
func InitSpanContext() {
	sc = &tracing.SpanContext{
		Ctx:           context.Background(),
		OperationName: "test",
		ParTraceCtx:   map[string]string{},
		TraceCtx:      map[string]string{},
		Peer:          "test",
		Method:        "get",
		URL:           "/etc/url",
		ComponentID:   "1",
		SpanLayerID:   "11",
		ServiceName:   "mesher"}
}

//TestClient
type TestClient struct {
}

//CreateEntrySpan
func (t *TestClient) CreateEntrySpan(sc *tracing.SpanContext) (interface{}, error) {
	return 1, nil
}

//CreateExitSpan
func (t *TestClient) CreateExitSpan(sc *tracing.SpanContext) (interface{}, error) {
	return 1, nil
}

//EndSpan
func (t *TestClient) EndSpan(sp interface{}, statusCode int) error {
	return nil
}

//NewApmClient
func NewApmClient(op tracing.TracingOptions) (apm.TracingClient, error) {
	var (
		err    error
		client TestClient
	)
	return &client, err
}

//InitApmClient
func InitApmClient() {
	apm.InstallClientPlugins("testclient", NewApmClient)
}

//TestInit
func TestInit(t *testing.T) {
	InitOption()
	InitApmClient()
	InitSpanContext()
	apm.Init(op)
}

//TestInstallClientPlugins
func TestInstallClientPlugins(t *testing.T) {
	InitOption()
	InitSpanContext()
	apm.InstallClientPlugins("testclient", NewApmClient)
	assert.Equal(t, nil, nil)
}

//TestCreateEntrySpan
func TestCreateEntrySpan(t *testing.T) {
	InitOption()
	InitSpanContext()
	apm.InstallClientPlugins("testclient", NewApmClient)
	span, err := apm.CreateEntrySpan(sc, op)
	assert.NotEqual(t, span, nil)
	assert.Equal(t, err, nil)

}

//TestCreateExitSpan
func TestCreateExitSpan(t *testing.T) {
	InitOption()
	InitSpanContext()
	apm.InstallClientPlugins("testclient", NewApmClient)
	span, err := apm.CreateExitSpan(sc, op)
	assert.NotEqual(t, span, nil)
	assert.Equal(t, err, nil)
}

//TestEndSpan
func TestEndSpan(t *testing.T) {
	InitOption()
	InitSpanContext()
	apm.InstallClientPlugins("testclient", NewApmClient)
	span, err := apm.CreateExitSpan(sc, op)
	assert.NotEqual(t, span, nil)
	assert.Equal(t, err, nil)
	err = apm.EndSpan(span, 1, op)
	assert.Equal(t, err, nil)
}
