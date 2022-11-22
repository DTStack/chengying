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
	"strconv"

	"github.com/kataras/iris/context"
)

type ApiParam struct {
	Type     string `json:"type"`
	Desc     string `json:"desc"`
	Default  string `json:"default"`
	Required bool   `json:"required"`
}

type ApiParams map[string]ApiParam

type ApiReturn struct {
	Type string `json:"type"`
	Desc string `json:"value"`
}

type ResultFields map[string]ApiReturn

type ApiReturnGroup struct {
	Desc   string       `json:"desc"`
	Fields ResultFields `json:"fields"`
}

type ApiDoc struct {
	Name    string           `json:"api-name" validate:"required"`
	Desc    string           `json:"desc"`
	Path    ApiParams        `json:"path-params"`
	Query   ApiParams        `json:"query-params"`
	Body    ApiParams        `json:"body-params"`
	Returns []ApiReturnGroup `json:"returns"`
}

func checkValue(errs *ApiParameterErrors, name string, value string, param *ApiParam, where string) {
	switch param.Type {
	case "int":
		if value == "" && param.Required {
			errs.AppendError(name, "No value for parameter in %s", where)
			return
		}
		if _, err := strconv.Atoi(value); err != nil {
			errs.AppendError(name, "Illegal value (%s) for int parameter in %s", value, where)
		}
	case "string", "":
		if value == "" && param.Required {
			errs.AppendError(name, "No value for string parameter '%s' in %s", where)
		}
	case "bool":
		if value != "" {
			if _, err := strconv.ParseBool(value); err != nil {
				errs.AppendError(name, "Invalid bool value for parameter in %s", name, where)
			}
		}
	default:
		errs.AppendError(name, "Unexpected type (%s) for the parameter in '%s'", param.Type, where)
	}
}

func ApiCheckRequestParameters(ctx context.Context, conf *ApiDoc, restrictMode bool) (ret bool) {
	errs := NewApiParameterErrors()
	defer func() {
		if r := recover(); r != nil {
			Feedback(ctx, r)
			ret = false
		} else {
			//ctx.Next()
			ret = true
		}
	}()

	for name, p := range conf.Path {
		v := ctx.Params().Get(name)
		if v == "" {
			if p.Required {
				errs.AppendError(name, "Missing required path parameter")
			}
			continue
		}
		checkValue(errs, name, v, &p, "path")
	}
	if restrictMode {
		ctx.Params().Visit(func(key string, value string) {
			if _, checked := conf.Path[key]; !checked {
				errs.AppendError(key, "Unexpected value '%s'", value)
			}
		})
	}

	queryParams := ctx.URLParams()
	for qname, q := range conf.Query {
		value, existed := queryParams[qname]
		if !existed {
			if q.Required {
				errs.AppendError(qname, "Missing required query parameter")
			}
			continue
		}
		checkValue(errs, qname, value, &q, "query")
	}

	if restrictMode {
		unexpectedQueryParams := []string{}
		for name, _ := range queryParams {
			_, checked := conf.Query[name]
			if !checked {
				unexpectedQueryParams = append(unexpectedQueryParams, name)
			}
		}
		if len(unexpectedQueryParams) > 0 {
			for _, p := range unexpectedQueryParams {
				errs.AppendError(p, "Unexpected query parameter as validated in restrict mode")
			}
		}
	}

	// checking body json parameters
	//if len(conf.Body) > 0 {
	//	json := map[string]interface{}{}
	//	for jpath, c := range conf.Body {
	//		result, err := JsonPathLookup(json, jpath)
	//		if err != nil {
	//			if IsJsonValueNotFound(err) {
	//				if c.Required {
	//					errs.AppendError(jpath, "required in body parameter")
	//				}
	//				continue
	//			} else {
	//				errs.AppendError(jpath, "Cannot find value by '%s' in json body: %s", jpath, err)
	//			}
	//		}
	//		switch reflect.TypeOf(result).Kind() {
	//		case reflect.Int, reflect.Uint, reflect.Float32, reflect.Float64:
	//			if c.Type != "int" && c.Type != "number" {
	//				errs.AppendError(jpath, "Unmatched number expect: %s", c.Type)
	//			}
	//		case reflect.String:
	//			if c.Type != "string" {
	//				errs.AppendError(jpath, "Unmatched string expect: %s", c.Type)
	//			}
	//		case reflect.Bool:
	//			if c.Type != "bool" {
	//				errs.AppendError(jpath, "Unmatched bool expect: %s", c.Type)
	//			}
	//		case reflect.Map:
	//			if c.Type != "object" && c.Type != "obj" && c.Type != "map" {
	//				errs.AppendError(jpath, "Unmatched object expect: %s", c.Type)
	//			}
	//		case reflect.Slice:
	//			if c.Type != "array" && c.Type != "list" {
	//				errs.AppendError(jpath, "Unmatched array expect: %s", c.Type)
	//			}
	//		}
	//	}
	//
	//	if restrictMode {
	//		// FIXME: iterate all json values to match the jsonpaths
	//	}
	//}

	errs.CheckAndThrowApiParameterErrors()
	return true
}
