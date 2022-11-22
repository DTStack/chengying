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

package api

import (
	"fmt"

	apibase "easyagent/go-common/api-base"
	"easyagent/internal/server/base"
	"easyagent/internal/server/log"

	"github.com/kataras/iris/context"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var SystemApiRoutes = apibase.Route{
	Path: "sys",
	SubRoutes: []apibase.Route{{
		Path: "version",
		GET: func(ctx context.Context) apibase.Result {
			return map[string]interface{}{
				"version": base.VERSION,
			}
		},
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取产品版本号",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.version": apibase.ApiReturn{"string", "产品版本号"},
					},
				}},
			},
		},
	}, {
		Path: "logs",
		GET: func(ctx context.Context) apibase.Result {
			return map[string]interface{}{
				"log-dir":         log.LOGDIR,
				"logger-max-size": log.LOGGER_MAX_SIZE,
				"days-to-keep":    log.LOGGER_MAX_AGE,
				"logger-rotates":  log.LOGGER_MAX_BKS,
			}
		},
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取日志设置",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.log-dir":         apibase.ApiReturn{"string", "日志输出目录"},
						"$.logger-max-size": apibase.ApiReturn{"int", "单种日志的文件最大值"},
						"$.days-to-keep":    apibase.ApiReturn{"int", "日志文件保留天数"},
						"$.logger-rotates":  apibase.ApiReturn{"int", "备份日志文件个数"},
					},
				}},
			},
		},
	}, {
		Path: "host",
		GET: func(ctx context.Context) apibase.Result {
			result := map[string]interface{}{}
			cpuInfo, err := cpu.Info()
			if err != nil {
				return fmt.Errorf("Cannot get cpuinfo: %v", err)
			}
			if len(cpuInfo) >= 1 {
				core0 := cpuInfo[0]
				result["cpu"] = fmt.Sprintf("%s-%s", core0.Family, core0.ModelName)
			}
			result["cpu-num"] = len(cpuInfo)

			memInfo, err := mem.VirtualMemory()
			if err != nil {
				return fmt.Errorf("Cannot get meminfo: %v", err)
			}
			result["mem-total"] = memInfo.Total
			result["mem-usage"] = memInfo.UsedPercent

			hostInfo, err := host.Info()
			if err != nil {
				return fmt.Errorf("Cannot get hostinfo: %v", err)
			}
			result["os"] = fmt.Sprintf("%s-%s", hostInfo.OS, hostInfo.Platform)
			result["os-version"] = hostInfo.PlatformVersion

			return result
		},
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取所在服务器信息",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.cpu":        apibase.ApiReturn{"object", "cpu型号"},
						"$.cpu-num":    apibase.ApiReturn{"int", "cpu内核数"},
						"$.mem-total":  apibase.ApiReturn{"int", "内存总数"},
						"$.mem-usage":  apibase.ApiReturn{"float", "内存使用百分比"},
						"$.os":         apibase.ApiReturn{"string", "OS类型"},
						"$.os-version": apibase.ApiReturn{"string", "操作系统版本"},
					},
				}},
			},
		},
	}, {
		Path: "rpc-server",
		GET: func(ctx context.Context) apibase.Result {
			return map[string]interface{}{
				"port":    base.RPC_SERVER_PORT,
				"use-tls": base.RPC_USE_TLS,
			}
		},
		Docs: apibase.Docs{
			GET: &apibase.ApiDoc{
				Name: "获取RPC服务信息",
				Returns: []apibase.ApiReturnGroup{{
					Fields: apibase.ResultFields{
						"$.port":    apibase.ApiReturn{"int", "rpc端口"},
						"$.use-tls": apibase.ApiReturn{"bool", "是否使用TLS"},
					},
				}},
			},
		},
	}},
}
