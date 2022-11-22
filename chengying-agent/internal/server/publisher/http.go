/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package publisher

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/elastic/go-ucfg"

	"fmt"

	"easyagent/internal/server/log"
	"easyagent/internal/server/util"
)

var (
	OutputNameHttp = "http"
)

const (
	HTTP_METHOD_GET  = "GET"
	HTTP_METHOD_POST = "POST"
)

type HttpClienter struct {
	httpClient *http.Client
	baseUrl    *url.URL
	restUrl    string
}

type httpServerResponse struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func init() {
	if err := Publish.RegisterOutputer(OutputNameHttp, NewHttpClient); err != nil {
		panic(err)
	}
}

// use to create a http request
func getTlsConfig() *tls.Config {
	var tlsConfig *tls.Config
	tlsConfig = &tls.Config{InsecureSkipVerify: true}
	return tlsConfig
}

func NewHttpClient(configContent map[string]*ucfg.Config) (Outputer, error) {
	cfg := HttpConfig{}
	if _, ok := configContent[OutputNameHttp]; !ok {
		return nil, nil
	}
	if err := configContent[OutputNameHttp].Unpack(&cfg); err != nil {
		return nil, err
	}
	client := &HttpClienter{httpClient: util.NewHTTPClient(getTlsConfig()), restUrl: cfg.ApiUri}
	var err error
	client.baseUrl, err = url.Parse("http://" + cfg.ApiHost)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (cli *HttpClienter) Name() string {
	return OutputNameHttp
}

func (cli *HttpClienter) OutputJson(ctx context.Context, id, index string, tpy string, jsonBody interface{}, key []byte) error {
	if id == "" {
		//filter os usage info
		return nil
	}
	if cli.httpClient == nil {
		return fmt.Errorf("[OutputJson]httpClient is null")
	}
	c := util.NewClient(cli.httpClient)
	c.BaseURL = cli.baseUrl
	respBody := new(httpServerResponse)

	r, err := c.NewRequest(HTTP_METHOD_POST, fmt.Sprintf(cli.restUrl, id, tpy), nil, jsonBody)

	if err != nil {
		log.Errorf("[HttpClienter] OutputJson Can not initialize REST request, id: %v, uri: %v, body: %v", id, fmt.Sprintf(cli.restUrl, tpy), jsonBody)
		return err
	}
	resp, err := c.Do(r, respBody)
	if err != nil {
		log.Errorf("[HttpClienter]id: %v,  OutputJson do request err: %v", id, err)
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if respBody.Code != 0 {
			log.Errorf("[HttpClienter] id:%v, OutputJson response code: %v", id, respBody)
			return fmt.Errorf("[HttpClienter]id:%v, OutputJson response code: %v", id, respBody)
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("[HttpClienter]id: %v,  OutputJson response err: %v", id, resp)
	}
	return nil
}

func (cli *HttpClienter) Close() {
}
