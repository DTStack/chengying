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

package apibase

import (
	"net/http"
	"runtime/debug"

	"github.com/kataras/iris/context"
)

const (
	UNKNOWN_ERR = iota + 100
	API_PARAM_ERR
	DB_MODEL_ERR
	RPC_HANDLE_ERR
)

type ApiResult struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func Feedback(ctx context.Context, result interface{}) {
	if err, ok := result.(error); ok {
		if IsApiParameterErrors(err) {
			errs, _ := err.(*ApiParameterErrors)
			data := map[string]string{}
			for pname, err := range errs.errors {
				data[pname] = err.Error()
			}
			ctx.JSON(&ApiResult{
				Code: API_PARAM_ERR,
				Msg:  "请求参数错误: Invalid parameter(s)",
				Data: data,
			})
		} else if IsDBModelError(err) {
			e, _ := err.(*DBModelError)
			ctx.JSON(&ApiResult{
				Code: DB_MODEL_ERR,
				Msg:  "获取agent数据失败: DB Model error",
				Data: e.err.Error(),
			})
		} else if IsRpcHandleError(err) {
			e, _ := err.(*RpcHandleError)
			ctx.JSON(&ApiResult{
				Code: RPC_HANDLE_ERR,
				Msg:  "与agent通信失败: Rpc handle error",
				Data: e.err.Error(),
			})
		} else {
			debug.PrintStack()
			ctx.JSON(&ApiResult{
				Code: UNKNOWN_ERR,
				Msg:  err.Error(),
				Data: err.Error(),
			})
		}
	} else if str, ok := result.(string); ok {
		ctx.WriteString(str)
	} else if bin, ok := result.([]byte); ok {
		ctx.Write(bin)
	} else {
		ctx.JSON(&ApiResult{
			Code: 0,
			Msg:  "ok",
			Data: result,
		})
	}
	ctx.StatusCode(http.StatusOK)
	ctx.Done()
}
