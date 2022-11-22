// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package kube

import (
	"context"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/view/request"
	monitorevents "dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	maxTryCount = 3
	baseurl     = "/api/v1/kube/resource/"
	Post        = "POST"
	Get         = "GET"
)

type RestClientCache struct {
	workspaceClient map[string]*RestClient
	mu              sync.RWMutex
}

type RestClient struct {
	workspace string
	c         *util.Client
}

func newClient(host, workspace string) (*RestClient, error) {
	parsedUrl, err := url.Parse(host)
	if err != nil {
		log.Errorf("parse url error : %s", host)
		return nil, err
	}
	client := util.NewClient(util.NewHTTPClient(util.GetTlsConfig()))
	client.BaseURL = parsedUrl

	return &RestClient{
		workspace: workspace,
		c:         client,
	}, nil
}

func (c *RestClientCache) Connect(server, workspace string) error {
	client, err := newClient(server, workspace)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.workspaceClient == nil {
		c.workspaceClient = make(map[string]*RestClient)
	}
	c.workspaceClient[workspace] = client
	return nil
}

func (c *RestClientCache) GetClient(workspace string) Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.workspaceClient == nil {
		return nil
	}
	return c.workspaceClient[workspace]
}

func (c *RestClientCache) DeleteClient(workspace string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.workspaceClient == nil {
		return
	}
	delete(c.workspaceClient, workspace)
}

func (c RestClientCache) Copy() ClientCache {
	return &c
}

func (c *RestClient) Apply(ctx context.Context, object runtime.Object) error {
	url := baseurl + "apply"
	method := Post
	body := toResouce(object)
	_, err := c.rest(method, url, nil, body)
	return err
}

func (c *RestClient) Status(ctx context.Context, object runtime.Object) error {
	url := baseurl + "status"
	method := Post
	body := toResouce(object)
	_, err := c.rest(method, url, nil, body)
	return err
}

func (c *RestClient) Create(ctx context.Context, object runtime.Object) error {
	url := baseurl + "create"
	method := Post
	body := toResouce(object)
	_, err := c.rest(method, url, nil, body)
	return err
}

func (c *RestClient) Update(ctx context.Context, object runtime.Object) error {
	url := baseurl + "update"
	method := Post
	body := toResouce(object)
	_, err := c.rest(method, url, nil, body)
	return err
}

func (c *RestClient) Delete(ctx context.Context, object runtime.Object) error {
	url := baseurl + "delete"
	method := Post
	body := toResouce(object)
	_, err := c.rest(method, url, nil, body)
	return err
}

func (c *RestClient) Get(ctx context.Context, object runtime.Object) (bool, error) {
	url := baseurl + "get"
	method := Post
	body := toResouce(object)
	rspdata, err := c.rest(method, url, nil, body)
	if err != nil {
		return false, err
	}
	if rspdata == nil {
		return false, nil
	}
	bts, err := json.Marshal(rspdata)
	if err != nil {
		log.Errorf("[rest_client]: mashal get rspdata %+v error %v", rspdata, err)
		return false, err
	}
	err = json.Unmarshal(bts, object)
	if err != nil {
		log.Errorf("[rest_client]: unmashl %s to obj %T error %v", string(bts), object, err)
		return false, err
	}
	return true, nil
}

func (c *RestClient) List(ctx context.Context, object runtime.Object, namespace string) error {
	url := baseurl + "list"
	method := Post
	gvks, _, _ := base.Schema.ObjectKinds(object)
	gvk := gvks[0]
	body := &request.ResourceList{
		Namespace: namespace,
		Group:     gvk.Group,
		Kind:      gvk.Kind,
		Version:   gvk.Version,
	}
	rspdata, err := c.rest(method, url, nil, body)
	if err != nil {
		return err
	}
	bts, err := json.Marshal(rspdata)
	if err != nil {
		log.Errorf("[rest_client]: mashal list rspdata %+v error %v", rspdata, err)
		return err
	}
	err = json.Unmarshal(bts, object)
	if err != nil {
		log.Errorf("[rest_client]: unmashl %s to obj %T error %v", string(bts), object, err)
		return err
	}
	return nil
}

func (c *RestClient) DryRun(action base.DryRunAction, object runtime.Object) error {
	url := baseurl + "dryrun"
	method := Post
	body := toResouce(object)
	body.Action = action
	_, err := c.rest(method, url, nil, body)
	return err
}

func (c *RestClient) Events(events *[]monitorevents.Event) error {
	url := baseurl + "events"
	method := Get
	resp, err := c.rest(method, url, nil, nil)
	if err != nil {
		return err
	}
	bts, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("[rest_client]: mashal array rspdata %+v error %v", resp, err)
		return err
	}
	err = json.Unmarshal(bts, events)
	if err != nil {
		log.Errorf("[rest_client]: unmashl %s to  monitorevents error %v", string(bts), err)
		return err
	}
	return nil
}

func (c *RestClient) rest(method, url string, params map[string]string, body interface{}) (interface{}, error) {
	var err error
	var req *http.Request
	resp := &apibase.ApiResult{}
	tryCount := 1
	for {
		if tryCount >= maxTryCount {
			break
		}
		req, err = c.c.NewRequest(method, url, params, body, "")
		if err != nil {
			log.Errorf("new request uri %s error: %v", url, err)
			return nil, err
		}
		_, err = c.c.Do(req, resp)
		if err != nil && tryCount < maxTryCount {
			tryCount++
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log.Errorf("[rest client]: request to %s error :%v", url, err)
			return nil, err
		}
		if resp.Code != 0 {
			return nil, fmt.Errorf("error from agent :%v", resp.Data)
		}
		return resp.Data, nil
	}
	return nil, fmt.Errorf("try three times, last error : %v", err)
}

func toResouce(object runtime.Object) *request.Resource {
	data, _ := json.Marshal(object)
	gvks, _, _ := base.Schema.ObjectKinds(object)
	gvk := gvks[0]
	return &request.Resource{
		Data:    data,
		Group:   gvk.Group,
		Kind:    gvk.Kind,
		Version: gvk.Version,
	}

}
