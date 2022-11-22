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

import "dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"

type KubeClientCache struct {
	c *base.ClientCache
}

func (c *KubeClientCache) GetClient(workspace string) Client{
	if c.c == nil{
		c.c = &base.ClientCache{}
	}
	return c.c.GetClient(workspace)
}

func (c *KubeClientCache) Connect(kubeconfig,workspace string) error{
	if c.c == nil{
		c.c = &base.ClientCache{}
	}
	return c.c.Connect(kubeconfig,workspace)
}

func (c *KubeClientCache) DeleteClient(workspace string){
	if c.c == nil{
		c.c = &base.ClientCache{}
	}
	c.c.DeleteClient(workspace)
}

func (c KubeClientCache) Copy() ClientCache{
	return &c
}
