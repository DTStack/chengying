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

package config

import (
	"errors"

	"github.com/elastic/go-ucfg"
	"github.com/satori/go.uuid"
)

var _ ucfg.Unpacker = &WrapperUUID{}

type WrapperUUID struct {
	uuid.UUID
}

func (w *WrapperUUID) Unpack(v interface{}) error {
	input, ok := v.(string)
	if !ok {
		return errors.New("need string")
	}
	u, err := uuid.FromString(input)
	if err != nil {
		return err
	}
	copy(w.UUID[:], u[:])
	return nil
}
