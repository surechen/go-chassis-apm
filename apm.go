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
	"strconv"
)

//ApmClient is an interface for application performance manage client

var apmClientPlugins = make(map[string]func(common.Options) (client.ApmClient, error))
var apmClients = make(map[string]client.ApmClient)

//InstallClientPlugins register apmclient create func
func InstallClientPlugins(name string, f func(common.Options) (client.ApmClient, error)) {
	apmClientPlugins[name] = f
	openlogging.Info("Install apm client: " + name)
}

//CreateSpans use invocation to make spans for apm
func CreateSpans(s *common.SpanContext, op common.Options) ([]interface{}, error) {
	openlogging.Info("CreateSpans")
	if client, ok := apmClients[op.APMName]; ok {
		openlogging.Info("client.CreateSpans")
		return client.CreateSpans(s)
	}
	var spans []interface{}
	return spans, nil
}

//EndSpans use invocation to make spans of apm end
func EndSpans(spans []interface{}, status int, op common.Options) error {
	openlogging.Info("EndSpans" + strconv.Itoa(status))
	if client, ok := apmClients[op.APMName]; ok {
		return client.EndSpans(spans, status)
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
