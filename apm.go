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

package apm

import (
	"github.com/go-chassis/go-chassis-apm/client"
	"github.com/go-chassis/go-chassis-apm/common"
	"github.com/go-mesh/openlogging"
)

var apmClientPlugins = make(map[string]func(common.Options) (client.ApmClient, error))
var apmClients = make(map[string]client.ApmClient)

//InstallClientPlugins register apmclient create func
func InstallClientPlugins(name string, f func(common.Options) (client.ApmClient, error)) {
	apmClientPlugins[name] = f
	openlogging.Info("Install apm client: " + name)
}

//CreateEntrySpan create entry span
func CreateEntrySpan(s *common.SpanContext, op common.Options) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateEntrySpan: %v", op)
	if client, ok := apmClients[op.APMName]; ok {
		openlogging.GetLogger().Debugf("CreateEntrySpan: %v", op)
		return client.CreateEntrySpan(s)
	}
	var spans interface{}
	return spans, nil
}

//CreateExitSpan create exit span
func CreateExitSpan(s *common.SpanContext, op common.Options) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateExitSpan: %v", op)
	if client, ok := apmClients[op.APMName]; ok {
		openlogging.GetLogger().Debugf("CreateExitSpan: %v", op)
		return client.CreateExitSpan(s)
	}
	var span interface{}
	return span, nil
}

//EndSpan end span
func EndSpan(span interface{}, status int, op common.Options) error {
	openlogging.GetLogger().Debugf("EndSpan: %v, status:%d", op, status)
	if client, ok := apmClients[op.APMName]; ok {
		openlogging.GetLogger().Debugf("EndSpan: %v, status:%d", op, status)
		return client.EndSpan(span, status)
	}
	return nil
}

//Init apm client
func Init(op common.Options) {
	openlogging.Info("Apm Init " + op.APMName + " " + op.ServerUri)
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
