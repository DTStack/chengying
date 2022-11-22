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
	"fmt"
	"time"

	"easyagent/internal/server/log"
	"github.com/elastic/go-ucfg"
)

type jsonBody struct {
	id        string
	index     string
	tpy       string
	jsonBytes interface{}
	key       []byte
}

type Publisher interface {
	RegisterOutputer(name string, c OutputCreater) error
	ConfigOutput(configContent map[string]*ucfg.Config) error
	OutputJson(ctx context.Context, id, index, tpy string, jsonBody interface{}, key []byte) error
	Close()
}

var (
	Publish Publisher = &publish{
		pubChan:        make(chan *jsonBody),
		outputRegister: make(map[string]OutputCreater),
		outputer:       make(map[string]Outputer),
	}
)

type publish struct {
	pubChan        chan *jsonBody
	outputRegister map[string]OutputCreater
	outputer       map[string]Outputer
}

func (p *publish) RegisterOutputer(name string, c OutputCreater) error {
	if _, ok := p.outputRegister[name]; ok {
		fmt.Println("Execution driver named " + name + " is already registered")
		return nil
	}
	p.outputRegister[name] = c
	return nil
}

func (p *publish) ConfigOutput(configContent map[string]*ucfg.Config) error {
	for outputName, _ := range p.outputRegister {
		log.Debugf("config output:%s\n", outputName)
		output, err := p.outputRegister[outputName](configContent)
		if err != nil {
			return err
		}
		if output != nil {
			p.outputer[outputName] = output
		}
	}
	go func() {
		for {
			select {
			case ctl, ok := <-p.pubChan:
				if !ok {
					log.Errorf("all outputs are closed")
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
				for _, outputer := range p.outputer {
					err := outputer.OutputJson(ctx, ctl.id, ctl.index, ctl.tpy, ctl.jsonBytes, ctl.key)
					if err != nil {
						log.Errorf("publish %v json data error: %v", outputer.Name(), err)
					}
				}
				cancel()
			}
		}
	}()
	return nil
}

func (p *publish) OutputJson(ctx context.Context, id, index, tpy string, js interface{}, key []byte) error {
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case p.pubChan <- &jsonBody{
		id:        id,
		index:     index,
		tpy:       tpy,
		jsonBytes: js,
		key:       key,
	}:
	}
	return err
}

func (p *publish) Close() {
	for _, outputer := range p.outputer {
		outputer.Close()
	}
}
