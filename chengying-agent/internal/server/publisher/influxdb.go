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
	"strings"

	"easyagent/internal/server/log"
	"github.com/elastic/go-ucfg"
	"github.com/influxdata/influxdb/client/v2"
)

var (
	OutputNameInfluxdb = "influxdb"
	InfluxdbDatabase   = "Metrics"
)

type InfluxdbClienter struct {
	client client.Client
}

func init() {
	if err := Publish.RegisterOutputer(OutputNameInfluxdb, InfluxdbClient); err != nil {
		panic(err)
	}
}

func InfluxdbClient(configContent map[string]*ucfg.Config) (Outputer, error) {
	cfg := influxdbConfig{}

	if _, ok := configContent[OutputNameInfluxdb]; !ok {
		return nil, nil
	}
	if err := configContent[OutputNameInfluxdb].Unpack(&cfg); err != nil {
		return nil, err
	}

	if len(cfg.UserName) > 0 {
		cli, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     strings.Join(cfg.Urls, ","),
			Username: cfg.UserName,
			Password: cfg.PassWord,
		})
		if err != nil {
			return nil, err
		}
		// create databases
		_, err = queryDB(cli, fmt.Sprintf("CREATE DATABASE %s", InfluxdbDatabase))
		if err != nil {
			log.Errorf("Error influxdb: ", err)
		}
		return &InfluxdbClienter{client: cli}, nil
	} else {
		cli, err := client.NewHTTPClient(client.HTTPConfig{
			Addr: strings.Join(cfg.Urls, ","),
		})
		if err != nil {
			return nil, err
		}
		// create databases
		_, err = queryDB(cli, fmt.Sprintf("CREATE DATABASE %s", InfluxdbDatabase))
		if err != nil {
			log.Errorf("Error influxdb: ", err)
		}
		return &InfluxdbClienter{client: cli}, nil
	}
}

//  方法名称前 加 (cli *handle);   func (cli *handle) Name () string{ }
//表示在此方法是 cli *handle的一个子项
func (cli *InfluxdbClienter) Name() string {
	return OutputNameInfluxdb
}

func queryDB(cli client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: InfluxdbDatabase,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (cli *InfluxdbClienter) OutputJson(ctx context.Context, id, index string, tpy string, jsonBody interface{}, key []byte) error {

	//set the Precision 时间是纳秒
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  InfluxdbDatabase,
		Precision: "ns",
	})

	if err != nil {
		log.Errorf("Error influxdb: ", err)
	}

	// Create a point and add to batch
	tags := map[string]string{"date": index}
	fields := map[string]interface{}{
		"type":   tpy,
		"metric": jsonBody,
	}

	pt, err := client.NewPoint("mertic", tags, fields)
	if err != nil {
		log.Errorf("Error influxdb: ", err)
	}
	bp.AddPoint(pt)

	// write
	if err := cli.client.Write(bp); err != nil {
		log.Errorf("Error influxdb: ", err)
	}
	return err

}

func (cli *InfluxdbClienter) Close() {
	cli.client.Close()
}
