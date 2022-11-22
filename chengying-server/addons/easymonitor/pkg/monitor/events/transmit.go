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

package events

import (
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"dtstack.com/dtstack/easymatrix/go-common/utils"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	MAX_EVENT_BUFFER = 1024
	REST_TRY_COUNT   = 3

	TRANSMIT_EVENT_URI = "/api/v2/cluster/kubernetes/listwatch/events"
)

type TransmitorInterface interface {
	Push(event Eventer)
	Process()
}

var Transmitor TransmitorInterface

//var (
//	Transmitor = &transmit{
//		c: make(chan Eventer, MAX_EVENT_BUFFER),
//	}
//)
//var Transmitor TransmitorInterface
func InitTransmiter(t TransmitorInterface) {
	Transmitor = t
}

//func InitTransmiter(matrix string, stopCh chan struct{}) error {
//	var err error
//	Transmitor.httpClient = utils.DefaultClient
//	Transmitor.baseUrl, err = url.Parse("http://" + matrix + ":8864")
//	Transmitor.stopCh = stopCh
//	if err != nil {
//		log.Errorf("[InitTransmiter] init err: %v", err)
//		return err
//	}
//	//start transmit event process
//	go Transmitor.Process()
//	return nil
//}
func SelfBuildModeTransmitor(matrix string, stopCh chan struct{}) error {
	var err error
	baseurl, err := url.Parse("http://" + matrix + ":8864")
	if err != nil {
		log.Errorf("[transmit]: parse url %s, error %v", baseurl, err)
	}
	t := &transmit{
		c:          make(chan Eventer, MAX_EVENT_BUFFER),
		httpClient: utils.DefaultClient,
		baseUrl:    baseurl,
		stopCh:     stopCh,
	}
	Transmitor = t
	//start transmit event process
	go Transmitor.Process()
	return nil
}

type transmit struct {
	c          chan Eventer
	httpClient *http.Client
	baseUrl    *url.URL
	stopCh     chan struct{}
}

func (t *transmit) Push(event Eventer) {
	t.c <- event
}

func (t *transmit) Process() {
	for {
		select {
		case ev := <-t.c:
			log.Infof("process event to matrix: %v", ev.Info())
			err := t.send(ev)
			if err != nil {
				log.Errorf("process event to matrix failed: %v", err.Error())
			}
		case <-t.stopCh:
			return
		}

	}
}

func (t *transmit) send(event Eventer) error {
	if t.httpClient == nil {
		return fmt.Errorf("http client is nil")
	}
	return t.restCall("POST", TRANSMIT_EVENT_URI, nil, event)
}

func (this *transmit) restCall(method, uri string, params map[string]string, body interface{}) error {
	c := util.NewClient(this.httpClient)
	c.BaseURL = this.baseUrl

	var err error
	var r *http.Request
	respBody := new(EventResponse)
	tryCount := 1
	for {
		if tryCount >= REST_TRY_COUNT {
			break
		}
		if r, err = c.NewRequest(method, uri, params, body, ""); err != nil {
			return err
		}
		_, err = c.Do(r, respBody)

		if err != nil && tryCount < REST_TRY_COUNT {
			tryCount++
			time.Sleep(1 * time.Second)
			continue
		}
		if respBody.Code != 0 {
			tryCount++
			err = fmt.Errorf("%v", respBody.Data)
			log.Errorf("restCall do request err: %v, try %v", err, tryCount)
			time.Sleep(1 * time.Second)
			continue
		} else {
			break
		}
	}
	return err
}
