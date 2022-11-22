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

	"easyagent/internal/server/log"
	"encoding/json"
	sarama "github.com/Shopify/sarama"
	"github.com/elastic/go-ucfg"
)

var (
	OutputNameKafka = "kafka"
)

type KafkaClienter struct {
	client sarama.SyncProducer
}

func init() {
	if err := Publish.RegisterOutputer(OutputNameKafka, NewKafkaClient); err != nil {
		panic(err)
	}
}

func NewKafkaClient(configContent map[string]*ucfg.Config) (Outputer, error) {
	cfg := kafkaConfig{}
	if _, ok := configContent[OutputNameKafka]; !ok {
		return nil, nil
	}
	if err := configContent[OutputNameKafka].Unpack(&cfg); err != nil {
		return nil, err
	}
	k := sarama.NewConfig()
	if len(cfg.UserName) > 0 {
		k.Net.SASL.Enable = true
		k.Net.SASL.User = cfg.UserName
		k.Net.SASL.Password = cfg.PassWord
	}

	k.Net.DialTimeout = cfg.Timeout
	k.Net.ReadTimeout = cfg.Timeout
	k.Net.WriteTimeout = cfg.Timeout

	k.Producer.Return.Successes = true // enable return channel for signaling
	k.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(cfg.Urls, k)
	if err != nil {
		return nil, err
	}
	return &KafkaClienter{client: producer}, nil
}

func (cli *KafkaClienter) Name() string {
	return OutputNameKafka
}

func (cli *KafkaClienter) OutputJson(ctx context.Context, id, index string, tpy string, jsonBody interface{}, key []byte) error {
	content, err := json.Marshal(jsonBody)
	if err != nil {
		log.Errorf("KafkaClienter OutputJson Marshal to string err %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: tpy,
		Value: sarama.StringEncoder(content),
	}

	partition, offset, err := cli.client.SendMessage(msg)

	if err != nil {
		log.Errorf("KafkaClienter OutputJson %v", err)
		return err
	} else {
		log.Debugf("> message %v sent to partition %d at offset %d\n", jsonBody, partition, offset)
	}
	return nil
}

func (cli *KafkaClienter) Close() {
	if cli.client != nil {
		cli.client.Close()
	}
}
