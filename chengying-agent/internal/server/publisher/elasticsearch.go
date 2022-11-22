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
	"time"

	"github.com/elastic/go-ucfg"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	OutputNameEs = "elasticsearch"

	indexMapping = `
{
  "mappings": {
    "metric": {
      "dynamic_templates": [
        {
          "em_string_template": {
            "match_mapping_type": "string",
            "mapping": {
              "ignore_above": "1024",
              "type": "keyword"
            }
          }
        }
      ]
    },
	"discover": {
      "dynamic_templates": [
        {
          "em_string_template": {
            "match_mapping_type": "string",
            "mapping": {
              "ignore_above": "1024",
              "type": "keyword"
            }
          }
        }
      ]
    }
  }
}
`
)

type ElasticClienter struct {
	client *elastic.Client
}

func init() {
	if err := Publish.RegisterOutputer(OutputNameEs, NewClient); err != nil {
		panic(err)
	}
}

func NewClient(configContent map[string]*ucfg.Config) (Outputer, error) {
	cfg := esConfig{}
	if _, ok := configContent[OutputNameEs]; !ok {
		return nil, nil
	}
	if err := configContent[OutputNameEs].Unpack(&cfg); err != nil {
		return nil, err
	}
	if len(cfg.UserName) > 0 {
		cli, err := elastic.NewClient(
			elastic.SetURL(cfg.Urls...),
			elastic.SetBasicAuth(cfg.UserName, cfg.PassWord),
		)
		if err != nil {
			return nil, err
		}
		return &ElasticClienter{client: cli}, nil
	} else {
		cli, err := elastic.NewClient(
			elastic.SetURL(cfg.Urls...),
		)
		if err != nil {
			return nil, err
		}
		return &ElasticClienter{client: cli}, nil
	}
}

func (cli *ElasticClienter) Name() string {
	return OutputNameEs
}

func (cli *ElasticClienter) OutputJson(ctx context.Context, id, index string, tpy string, jsonBody interface{}, key []byte) error {
	is := cli.client.Index().
		Index(index).
		Type(tpy).
		BodyJson(jsonBody).
		Refresh("false")
DO:
	_, err := is.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			if exists, _ := cli.client.IndexExists(index).Do(ctx); !exists {
				if _, err = cli.client.CreateIndex(index).Body(indexMapping).Do(ctx); err != nil {
					time.Sleep(2 * time.Second) // take a rest
				}
				goto DO
			}
		}
	}
	return err
}

func (cli *ElasticClienter) Close() {
	if cli.client != nil {
		cli.client.Stop()
	}
}
