package apm

import (
	"github.com/go-chassis/go-chassis-apm/tracing"
	"github.com/go-mesh/openlogging"
)

//TracingClient for apm interface
type TracingClient interface {
	CreateEntrySpan(sc *tracing.SpanContext) (interface{}, error)
	CreateExitSpan(sc *tracing.SpanContext) (interface{}, error)
	EndSpan(sp interface{}, statusCode int) error
}

var apmClientPlugins = make(map[string]func(tracing.TracingOptions) (TracingClient, error))
var apmClients = make(map[string]TracingClient)

//InstallClientPlugins register TracingClient create func
func InstallClientPlugins(name string, f func(tracing.TracingOptions) (TracingClient, error)) {
	apmClientPlugins[name] = f
	openlogging.Info("Install apm client: " + name)
}

//CreateEntrySpan create entry span
func CreateEntrySpan(s *tracing.SpanContext, op tracing.TracingOptions) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateEntrySpan: %v", op)
	if client, ok := apmClients[op.APMName]; ok {
		openlogging.GetLogger().Debugf("CreateEntrySpan: %v", op)
		return client.CreateEntrySpan(s)
	}
	var spans interface{}
	return spans, nil
}

//CreateExitSpan create exit span
func CreateExitSpan(s *tracing.SpanContext, op tracing.TracingOptions) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateExitSpan: %v", op)
	if client, ok := apmClients[op.APMName]; ok {
		openlogging.GetLogger().Debugf("CreateExitSpan: %v", op)
		return client.CreateExitSpan(s)
	}
	var span interface{}
	return span, nil
}

//EndSpan end span
func EndSpan(span interface{}, status int, op tracing.TracingOptions) error {
	openlogging.GetLogger().Debugf("EndSpan: %v, status:%d", op, status)
	if client, ok := apmClients[op.APMName]; ok {
		openlogging.GetLogger().Debugf("EndSpan: %v, status:%d", op, status)
		return client.EndSpan(span, status)
	}
	return nil
}

//Init apm client
func Init(op tracing.TracingOptions) {
	openlogging.Info("Apm Init " + op.APMName + " " + op.ServerURI)
	f, ok := apmClientPlugins[op.APMName]
	if ok {
		client, err := f(op)
		if err == nil {
			apmClients[op.APMName] = client
		} else {
			openlogging.Error("apmClients init failed. " + err.Error())
		}
	}
}
