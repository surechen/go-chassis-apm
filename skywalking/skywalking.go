/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package skywalking

import (
	"container/list"
	"github.com/go-chassis/go-chassis-apm"
	"github.com/go-chassis/go-chassis-apm/client"
	"github.com/go-chassis/go-chassis-apm/common"
	"github.com/go-mesh/openlogging"
	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/reporter"
	skycom "github.com/tetratelabs/go2sky/reporter/grpc/common"
	"strconv"
)

const (
	HTTPPrefix                = "http://"
	CrossProcessProtocolV2    = "Sw6"
	Name                      = "skywalking"
	DefaultTraceContext = ""
)

const (
	HTTPClientComponentID = 2
	HTTPServerComponentID = 49
)

//SkyWalkingClient for connecting and reporting to skywalking server
type SkyWalkingClient struct {
	reporter go2sky.Reporter
	tracer   *go2sky.Tracer
}

//CreateEntrySpan create entry span
func (s *SkyWalkingClient) CreateEntrySpan(sc *common.SpanContext) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateEntrySpan begin. spanctx:%#v", sc)
	span, ctx, err := s.tracer.CreateEntrySpan(sc.Ctx, sc.OperationName, func() (string, error) {
		if sc.TraceContext != nil {
			return sc.TraceContext[CrossProcessProtocolV2], nil
		}
		return DefaultTraceContext, nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateExitSpan error:%s", err.Error())
		return &span, err
	}
	span.Tag(go2sky.TagHTTPMethod, sc.Method)
	span.Tag(go2sky.TagURL, sc.URL)
	span.SetSpanLayer(skycom.SpanLayer_Http)
	span.SetComponent(HTTPServerComponentID)
	sc.Ctx = ctx
	return &span, err
}

//CreateExitSpan create end span
func (s *SkyWalkingClient) CreateExitSpan(sc *common.SpanContext) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateExitSpan begin. spanctx:%v", sc)
	span, err := s.tracer.CreateExitSpan(sc.Ctx, sc.OperationName, sc.Peer, func(header string) error {
		sc.TraceContext[CrossProcessProtocolV2] = header
		return nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateExitSpan error:%s", err.Error())
		return &span, err
	}
	span.Tag(go2sky.TagHTTPMethod, sc.Method)
	span.Tag(go2sky.TagURL, sc.URL)
	span.SetSpanLayer(skycom.SpanLayer_Http)
	span.SetComponent(HTTPClientComponentID)
	return &span, err
}

//EndSpan make span end and report to skywalking
func (s *SkyWalkingClient) EndSpan(sp interface{}, statusCode int) error {
	span, ok := (sp).(*go2sky.Span)
	if !ok || span == nil {
		return nil
	}
	(*span).Tag(go2sky.TagStatusCode, strconv.Itoa(statusCode))
	(*span).End()
	return nil
}

//CreateSpans create entry and exit spans for report
func (s *SkyWalkingClient) CreateSpans(sc *common.SpanContext) ([]interface{}, error) {
	openlogging.GetLogger().Debugf("CreateSpans begin. spanctx:%#v", sc)
	var spans []interface{}
	span, ctx, err := s.tracer.CreateEntrySpan(sc.Ctx, sc.OperationName, func() (string, error) {
		if sc.TraceContext != nil {
			return sc.TraceContext[CrossProcessProtocolV2], nil
		}
		return DefaultTraceContext, nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateSpans error:%s", err.Error())
		return spans, err
	}
	l := list.New()
	l.PushBack(1)
	span.Tag(go2sky.TagHTTPMethod, sc.Method)
	span.Tag(go2sky.TagURL, sc.URL)
	span.SetSpanLayer(skycom.SpanLayer_Http)
	span.SetComponent(HTTPServerComponentID)
	spans = append(spans, &span)
	spanExit, err := s.tracer.CreateExitSpan(ctx, sc.OperationName, sc.Peer, func(header string) error {
		sc.TraceContext[CrossProcessProtocolV2] = header
		return nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateSpans error:%s", err.Error())
		return spans, err
	}
	spanExit.Tag(go2sky.TagHTTPMethod, sc.Method)
	spanExit.Tag(go2sky.TagURL, sc.URL)
	spanExit.SetSpanLayer(skycom.SpanLayer_Http)
	spanExit.SetComponent(HTTPClientComponentID)
	spans = append(spans, &spanExit)
	return spans, nil

}

//EndSpans make spans end and report to skywalking
func (s *SkyWalkingClient) EndSpans(spans []interface{}, status int) error {
	openlogging.GetLogger().Debugf("EndSpans spans:%u status:%#v", len(spans), status)
	for i := len(spans) - 1; i >= 0; i-- {
		span, ok := (spans[i]).(*go2sky.Span)
		if !ok || spans[i] == nil {
			continue
		}
		(*span).Tag(go2sky.TagStatusCode, strconv.Itoa(status))
		(*span).End()
	}
	return nil
}

//NewApmClient init report and tracer for connecting and sending messages to skywalking server
func NewApmClient(op common.Options) (client.ApmClient, error) {
	var (
		err    error
		client SkyWalkingClient
	)
	client.reporter, err = reporter.NewGRPCReporter(op.ServerUri)
	if err != nil {
		openlogging.GetLogger().Errorf("NewGRPCReporter error:%s", err.Error())
		return &client, err
	}
	client.tracer, err = go2sky.NewTracer(op.MicServiceName, go2sky.WithReporter(client.reporter))
	//t.WaitUntilRegister()
	if err != nil {
		openlogging.GetLogger().Errorf("NewTracer error:%s", err.Error())
		return &client, err

	}
	openlogging.GetLogger().Debugf("NewApmClient succ. name:%s uri:%s", op.APMName, op.ServerUri)
	return &client, err
}

func init() {
	apm.InstallClientPlugins(Name, NewApmClient)
}
