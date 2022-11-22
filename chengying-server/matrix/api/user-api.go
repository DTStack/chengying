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

package api

import (
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/api/impl"
	"github.com/kataras/iris/context"
)

var UserOperationEasyMatrixAPIRoutes = apibase.Route{
	Path: "user",
	SubRoutes: []apibase.Route{{
		Path: "identity",
		GET:  impl.Identity,
	}, {
		Path: "modifyInfo",
		POST: impl.ModifyInfoById,
	}, {
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		Path: "modifyInfoByAdmin",
		POST: impl.ModifyInfoByAdmin,
	}, {
		Path: "modifyPwd",
		POST: impl.ModifyPwdById,
	}, {
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		Path: "register",
		POST: impl.Register,
	}, {
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		Path: "list",
		GET:  impl.UserInfo,
	}, {
		Path: "login",
		POST: impl.Login,
	}, {
		Path: "getPublicKey",
		GET:  impl.GetPublicKey,
	}, {
		Path: "logout",
		POST: impl.LogOut,
	}, {
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		Path: "remove",
		POST: impl.RemoveUserById,
	}, {
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		Path: "resetPwdByAdmin",
		POST: impl.ResetPwdByAdmin,
	}, {
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		Path: "enable",
		POST: impl.Enable,
	}, {
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		Path: "disable",
		POST: impl.Disable,
	}, {
		Path: "personal",
		POST: impl.Personal,
	}, {
		Path: "getCaptcha",
		GET:  impl.GetCaptcha,
	}, {
		Path: "processCaptcha",
		POST: impl.ProcessCaptcha,
	}, {
		Path: "showCaptcha/{captcha}",
		GET:  impl.ShowCaptcha,
	}, {
		Path: "sys_config",
		Middlewares: []context.Handler{
			apibase.CheckPermission1,
		},
		SubRoutes: []apibase.Route{
			{
				Path: "platformSecurity",
				GET:  impl.GetSysconfigPlatformSecurity,
				POST: impl.ModifySysconfigPlatformSecurity,
			},
			{
				Path: "globalConfig",
				GET:  impl.GetGlobalConfig,
				POST: impl.ModifyGlobalConfig,
			},
		},
	}},
}
