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

package apibase

import (
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"errors"
	"fmt"
	"net/http"
	_ "runtime/debug"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
)

const (
	JWT_SECRET = "SGVsbG8gV29ybGQK"
)

type user struct {
	UserId   int
	UserName string
}

var (
	userCache = map[int]user{}
)

func SetUserCache(userId int, username string) {
	userCache[userId] = user{UserId: userId, UserName: username}
}

func DeleteUserCache(userId int) {
	delete(userCache, userId)
}

func getUserCache(userId int) *user {
	userInfo, ok := userCache[userId]
	if ok {
		return &userInfo
	}
	return nil
}

// 定义api接口认证相关检查，非认证的接口不进行cookie的验证
func ApiValidateDtstackCookies(ctx context.Context) bool {

	// 定义非授权的路由
	token := ctx.GetCookie("em_token")

	if strings.HasSuffix(ctx.Path(), "register") ||
		strings.HasSuffix(ctx.Path(), "login") ||
		strings.HasSuffix(ctx.Path(), "license") ||
		strings.HasSuffix(ctx.Path(), "dt_agent_health_check") ||
		strings.HasSuffix(ctx.Path(), "dt_agent_error") ||
		strings.HasSuffix(ctx.Path(), "dt_agent_host_resource") ||
		strings.HasSuffix(ctx.Path(), "dt_agent_performance") ||
		strings.HasSuffix(ctx.Path(), "callback") ||
		strings.HasSuffix(ctx.Path(), "identity") ||
		strings.Contains(ctx.Path(), "Captcha") ||
		strings.Contains(ctx.Path(), "listwatch") ||
		strings.Contains(ctx.Path(), "seqReport") ||
		strings.Contains(ctx.Path(), "isShowLog") ||
		strings.Contains(ctx.Path(), "shellStatusReport") ||
		strings.Contains(ctx.Path(), "getPublicKey") {
		return true
	}
	if len(token) == 0 && strings.HasSuffix(ctx.Path(), "upload") {
		return true
	}
	if len(token) == 0 && strings.HasSuffix(ctx.Path(), "workloadinit") {
		return true
	}

	// 定义认证失败信息回显方法
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
			Feedback(ctx, &AccessDeniedError{errors.New(r.(string))})
		}
	}()

	// 解析token为jwt的结构体，然后进行验证
	tk, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		if tk.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("Unexpected auth method: %s", tk.Method.Alg())
		}
		return []byte(JWT_SECRET), nil
	})
	reject(err, "Validate em_token failure: %s", err)
	reject(!tk.Valid, "Invalid em_token: %s", token)

	// 验证cookie当中携带的用户id，用户名，以及token过期时间
	var expiration time.Time
	var info *user
	if claims, ok := tk.Claims.(jwt.MapClaims); ok {
		if !cache.ValidationSessionId(ctx.Path(), claims["session_hash_id"].(string)) {
			ctx.RemoveCookie("em_token")
			reject(true, "No operation for a long time", nil)
		}
		if id, ok := claims["user_id"].(float64); ok {
			info = getUserCache(int(id))
			reject(info == nil, "Invalid user_id (%s) in em_token", id)
			reject(int(id) != info.UserId, "Unmatched user_id in cookies")
		} else {
			reject(true, "Missing 'user_id' in em_token")
		}
		if claim_username, ok := claims["user_name"].(string); ok {
			reject(claim_username != info.UserName, "Unmatched user_name in cookies")
		} else {
			reject(true, "Missing 'user_name' in em_token")
		}
		exp, ok := claims["exp"].(float64)
		reject(!ok, "Missing 'exp' in em_token")
		expiration = time.Unix(int64(exp), 0)
		reject(time.Now().After(expiration), "em_token is expired")
	} else {
		reject(true, "Unable to map em_token")
	}

	// shift to next
	ctx.Values().Set("userId", info.UserId)
	ctx.Values().Set("username", info.UserName)
	ctx.Values().Set("expiration", expiration)

	return true
}

func CreateToken(username string, userId, permission int) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["user_name"] = username
	claims["user_id"] = userId
	claims["user_permission"] = permission
	sessionId := cache.CreateSessionId()
	claims["session_hash_id"] = sessionId
	// 设置token过期时间为三天
	claims["exp"] = time.Now().Unix() + 259200

	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(JWT_SECRET))

	SetUserCache(userId, username)
	cache.SetSessionCache(sessionId)

	return tokenString
}

var (
	SetCookieExpiration = time.Duration(259200) * time.Second
)

func SetCookie(ctx context.Context, name, value string) {
	c := &http.Cookie{}
	c.Name = name
	c.Value = value
	c.HttpOnly = true
	c.Expires = time.Now().Add(SetCookieExpiration)
	c.MaxAge = int(SetCookieExpiration.Seconds())
	ctx.SetCookie(c)

	ctx.Header(name, value)
}

func GetTokenUserPermission(ctx context.Context) (int, error) {
	token := ctx.GetCookie("em_token")
	tk, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		if tk.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("Unexpected auth method: %s", tk.Method.Alg())
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return 0, fmt.Errorf("validate em_token failure: %s", err)
	}
	if !tk.Valid {
		return 0, fmt.Errorf("invalid em_token: %s", err)
	}
	if claims, ok := tk.Claims.(jwt.MapClaims); ok {
		return int(claims["user_permission"].(float64)), nil
	}
	return 0, fmt.Errorf("unable to map em_token")
}
