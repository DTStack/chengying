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
	"fmt"
	"net/http"
	_ "runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
)

const (
	JWT_SECRET = "SGVsbG8gV29ybGQK"
)

func ApiValidateDtstackCookies(ctx context.Context) {
	reject := func(result interface{}, format string, args ...interface{}) {
		do_panic := false
		if b, ok := result.(bool); ok {
			if b {
				do_panic = true
			}
		} else if err, ok := result.(error); ok {
			if err != nil {
				do_panic = true
			}
		}
		if do_panic {
			panic(fmt.Sprintf(format, args...))
		}
	}

	defer func() {
		if r := recover(); r != nil {
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.Header("Location", "https://account.dtstack.com/login.html")
			ctx.WriteString(fmt.Sprintf("%s", r))
			ctx.Done()
		}
	}()

	token := ctx.GetCookie("dt_token")

	tenant_id, err := strconv.Atoi(ctx.GetCookie("dt_tenant_id"))
	reject(err, "Illegal dt_tenant_id value")

	user_id, err := strconv.Atoi(ctx.GetCookie("dt_user_id"))
	reject(err, "Illegal dt_user_id value")

	username := ctx.GetCookie("dt_username")
	username = strings.TrimRight(strings.TrimLeft(username, "\"'"), "\"'")

	session_id := ctx.GetCookie("sessionid")

	role_type, err := strconv.Atoi(ctx.GetCookie("role_type"))
	reject(err, "Illegal role_type value")

	// validate
	tk, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		if tk.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("Unexpected auth method: %s", tk.Method.Alg())
		}
		return []byte(JWT_SECRET), nil
	})
	reject(err, "Validate dt_token failure: %s", err)
	reject(!tk.Valid, "Invalid dt_token: %s", token)

	var expiration time.Time
	if claims, ok := tk.Claims.(jwt.MapClaims); ok {
		if claim_username, ok := claims["user_name"].(string); ok {
			reject(claim_username != username, "Unmatched dt_username in cookies")
		} else {
			reject(true, "Missing 'user_name' in dt_token")
		}
		if _id, ok := claims["tenant_id"].(string); ok {
			id, err := strconv.Atoi(_id)
			reject(err, "Invalid tenant_id (%s) in dt_token", _id)
			reject(id != tenant_id, "Unmatched dt_tenant_id in cookies: %s %v", id, tenant_id)
		} else {
			reject(true, "Missing 'tenant_id' in dt_token")
		}
		if _id, ok := claims["user_id"].(string); ok {
			id, err := strconv.Atoi(_id)
			reject(err, "Invalid user_id (%s) in dt_token", _id)
			reject(id != user_id, "Unmatched dt_user_id in cookies")
		} else {
			reject(true, "Missing 'user_id' in dt_token")
		}
		exp, ok := claims["exp"].(float64)
		reject(!ok, "Missing 'exp' in dt_token")
		expiration = time.Unix(int64(exp), 0)
		reject(time.Now().After(expiration), "dt_token is expired")
	} else {
		reject(true, "Unable to map dt_token")
	}

	// FIXME: validate session id
	session_id = session_id

	// shift to next
	ctx.Values().Set("tenant_id", tenant_id)
	ctx.Values().Set("user_id", user_id)
	ctx.Values().Set("username", username)
	ctx.Values().Set("role_type", role_type)
	ctx.Values().Set("expiration", expiration)
	ctx.Next()

}
