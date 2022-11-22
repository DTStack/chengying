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
	"strings"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

type Result interface{}

type Handler func(context.Context) Result

type Docs struct {
	Desc                   string
	GET, POST, PUT, DELETE *ApiDoc
}

type Route struct {
	Path                   string
	Middlewares            []context.Handler
	GET, POST, PUT, DELETE Handler
	Docs                   Docs
	SubRoutes              []Route
}

func WrapHandler(handler Handler, doc *ApiDoc, restrictMode bool) context.Handler {
	return func(ctx context.Context) {
		if doc != nil {
			if !ApiCheckRequestParameters(ctx, doc, restrictMode) {
				return
			}
		}

		defer func() {
			if r := recover(); r != nil {
				Feedback(ctx, r)
			}
		}()

		Feedback(ctx, handler(ctx))
	}
}

func createRoute(parent router.Party, route *Route, restrictMode bool) error {

	routePath := route.Path
	hasSubRoutes := len(route.SubRoutes) > 0

	if !strings.HasPrefix(routePath, "/") {
		routePath = "/" + routePath
	}

	middlewares := route.Middlewares
	if middlewares != nil {
		middlewares = []context.Handler{}
	}

	methodInstallers := []struct {
		installer func(path string, handlers ...context.Handler) (*router.Route, error)
		handler   Handler
		doc       *ApiDoc
	}{
		{parent.Get, route.GET, route.Docs.GET},
		{parent.Post, route.POST, route.Docs.POST},
		{parent.Put, route.PUT, route.Docs.PUT},
		{parent.Delete, route.DELETE, route.Docs.DELETE},
	}

	for _, i := range methodInstallers {
		if i.handler != nil {
			handlers := append(middlewares, WrapHandler(i.handler, i.doc, restrictMode))
			_, err := i.installer(routePath, handlers...)
			if err != nil {
				return err
			}
		}
	}

	if hasSubRoutes {
		node := parent.Party(routePath, middlewares...)
		for _, sub := range route.SubRoutes {
			err := createRoute(node, &sub, restrictMode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func walkSchema(ref, doc *Route) {
	subRoutes := []string{}
	for _, s := range ref.SubRoutes {
		subRoutes = append(subRoutes, s.Path)
	}
	doc.GET = func(path string, docs Docs, subRoutes []string) Handler {
		return func(ctx context.Context) Result {
			return map[string]interface{}{
				"path":          path,
				"description":   docs.Desc,
				"method.get":    docs.GET,
				"method.post":   docs.POST,
				"method.put":    docs.PUT,
				"method.delete": docs.DELETE,
				"sub-routes":    subRoutes,
			}
		}
	}(ref.Path, ref.Docs, subRoutes)
	doc.SubRoutes = []Route{}
	for _, s := range ref.SubRoutes {
		path := s.Path
		if path[0] == '/' {
			path = path[1:]
		}
		if path[0] == '{' {
			commaIdx := strings.Index(path, ":")
			if commaIdx > 0 {
				path = path[1:commaIdx]
			} else {
				spaceIdx := strings.Index(path, " \t")
				if spaceIdx > 0 {
					path = path[1:spaceIdx]
				}
			}
		}
		sub := Route{
			Path: path,
		}
		walkSchema(&s, &sub)
		doc.SubRoutes = append(doc.SubRoutes, sub)
	}
}

func GenerateDocSchema(root *Route) {
	doc := Route{
		Path: "schema",
	}
	walkSchema(root, &doc)
	root.SubRoutes = append(root.SubRoutes, doc)
}

func InitSchema(app *iris.Application, root *Route, restrictMode bool) error {
	GenerateDocSchema(root)
	return createRoute(app, root, restrictMode)
}
