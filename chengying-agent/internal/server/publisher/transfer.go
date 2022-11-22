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
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/server/log"
	"github.com/elastic/go-ucfg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	OutputNameTransfer = "transfer"
)

type TransferClienter struct {
	client proto.TransferServiceClient
	conn   *grpc.ClientConn
	ch     chan *proto.Msg
}

func init() {
	if err := Publish.RegisterOutputer(OutputNameTransfer, NewClientTransfer); err != nil {
		panic(err)
	}
}

func NewClientTransfer(configContent map[string]*ucfg.Config) (Outputer, error) {
	cfg := TransferConfig{Concurrency: 1, Timeout: 3 * time.Second}
	if _, ok := configContent[OutputNameTransfer]; !ok {
		return nil, nil
	}
	if err := configContent[OutputNameTransfer].Unpack(&cfg); err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{grpc.WithBlock()}
	if cfg.Tls {
		var cp *x509.CertPool
		if cfg.CertFile != "" {
			b, err := ioutil.ReadFile(cfg.CertFile)
			if err != nil {
				return nil, err
			}
			cp = x509.NewCertPool()
			if !cp.AppendCertsFromPEM(b) {
				return nil, errors.New("credentials: failed to append certificates")
			}
		}
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: cfg.TlsSkipVerify, RootCAs: cp})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	c, err := grpc.DialContext(ctx, net.JoinHostPort(cfg.Server, strconv.Itoa(cfg.Port)), opts...)
	if err != nil {
		return nil, err
	}

	cli := &TransferClienter{
		client: proto.NewTransferServiceClient(c),
		conn:   c,
		ch:     make(chan *proto.Msg),
	}
	for i := uint8(0); i < cfg.Concurrency; i++ {
		go cli.sendLoop()
	}
	return cli, nil
}

func (cli *TransferClienter) Name() string {
	return OutputNameTransfer
}

func (cli *TransferClienter) OutputJson(ctx context.Context, id, index string, tpy string, jsonBody interface{}, key []byte) error {
	jsonBytes, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case cli.ch <- &proto.Msg{Id: id, Index: index, Typ: tpy, JsonBytes: jsonBytes, Key: key}:
	}

	return err
}

func (cli *TransferClienter) sendLoop() {
	for msg := range cli.ch {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		if _, err := cli.client.Send(ctx, msg); err != nil {
			log.Errorf("send grpc error: %v", err)
		}
		cancel()
	}
}

func (cli *TransferClienter) Close() {
	if cli.conn != nil {
		cli.conn.Close()
	}
}
