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

package dbhelper

import (
	"database/sql"
	apibase "easyagent/go-common/api-base"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
)

type DbTable struct {
	GetDB     func() *sqlx.DB
	TableName string
}

func (h *DbTable) formatQuery(query string, args ...interface{}) (string, []interface{}) {
	argList := []interface{}{}
	valList := []interface{}{}
	argIdx := 0
	tblNameBegin := 0
	tblNameIndexs := []int{}

	for i, r := range query {
		switch r {
		case '{':
			tblNameBegin = 1
		case 'T':
			if tblNameBegin == 1 {
				tblNameBegin = 2
			} else {
				tblNameBegin = 0
			}
		case '}':
			if tblNameBegin == 2 {
				tblNameIndexs = append(tblNameIndexs, i-2)
			}
			tblNameBegin = 0
		case '?':
			valList = append(valList, args[argIdx])
			argIdx++
			tblNameBegin = 0
		case '%':
			argList = append(argList, args[argIdx])
			argIdx++
			tblNameBegin = 0
		default:
			tblNameBegin = 0
		}
	}

	if len(tblNameIndexs) > 0 {
		for i := len(tblNameIndexs) - 1; i >= 0; i-- {
			index := tblNameIndexs[i]
			query = query[0:index] + h.TableName + query[index+3:]
		}
	}

	return fmt.Sprintf(query, argList...), valList
}

func (h *DbTable) Exec(query string, args ...interface{}) (sql.Result, error) {
	query, args = h.formatQuery(query, args...)
	return h.GetDB().Exec(query, args...)
}

func (h *DbTable) MustExec(query string, args ...interface{}) sql.Result {
	query, args = h.formatQuery(query, args...)
	return h.GetDB().MustExec(query, args...)
}

func (h *DbTable) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	query, args = h.formatQuery(query, args...)
	return h.GetDB().Queryx(query, args...)
}

func (h *DbTable) QueryRow(query string, args ...interface{}) *sqlx.Row {
	query, args = h.formatQuery(query, args...)
	return h.GetDB().QueryRowx(query, args...)
}

func (h *DbTable) Get(dest interface{}, query string, args ...interface{}) error {
	query, args = h.formatQuery(query, args...)
	return h.GetDB().Get(dest, query, args...)
}

func (h *DbTable) QueryWithPagination(query string, pagination *apibase.Pagination, args ...interface{}) (*sqlx.Rows, error) {
	query, args = h.formatQuery(query+" "+pagination.AsQuery(), args...)
	return h.GetDB().Queryx(query, args...)
}

func makeSelectors(keys []string) string {
	str := ""
	for i, k := range keys {
		if i > 0 {
			str += ","
		}
		str += "`" + k + "`"
	}
	return str
}

type causeExpression struct {
	expr        string
	value       interface{}
	multiValues bool
}
type WhereCause []causeExpression

func wrapWhere(w string) string {
	if strings.Contains(w, ".") {
		l := []string{}
		for _, w0 := range strings.Split(w, ".") {
			l = append(l, "`"+w0+"`")
		}
		w = strings.Join(l, ".")
	} else {
		w = "`" + w + "`"
	}
	return w
}

func MakeWhereCause() WhereCause {
	return WhereCause{}
}

func (w WhereCause) Equal(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " = ? ", value, false})
}

func (w WhereCause) NotEqual(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " != ? ", value, false})
}

func (w WhereCause) GreaterThan(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " > ? ", value, false})
}

func (w WhereCause) GreaterAndEqualThan(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " >= ? ", value, false})
}

func (w WhereCause) LittleThan(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " < ? ", value, false})
}

func (w WhereCause) LittleAndEqualThan(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " <= ? ", value, false})
}

func (w WhereCause) Like(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " LIKE ? ", value, false})
}

func (w WhereCause) NotLike(where string, value interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " NOT LIKE ? ", value, false})
}

func (w WhereCause) Included(where string, values ...interface{}) WhereCause {
	where = wrapWhere(where)
	whereCause := where + " IN ("
	for i := 0; i < len(values); i++ {
		if i > 0 {
			whereCause += ","
		}
		whereCause += "?"
	}
	return append(w, causeExpression{whereCause + ")", values, true})
}

func (w WhereCause) Between(where string, from, to interface{}) WhereCause {
	where = wrapWhere(where)
	return append(w, causeExpression{where + " BETWEEN ? AND ? ", []interface{}{from, to}, true})
}

func (w WhereCause) And() WhereCause {
	return append(w, causeExpression{" AND ", nil, false})
}

func (w WhereCause) Or() WhereCause {
	return append(w, causeExpression{" OR ", nil, false})
}

func (w WhereCause) SQL() (string, []interface{}) {
	str := ""
	if len(w) == 0 {
		return str, nil
	}
	str += "WHERE "
	values := []interface{}{}
	for _, c := range w {
		str += c.expr
		if c.value != nil {
			if c.multiValues {
				if l, ok := c.value.([]interface{}); ok {
					for _, v := range l {
						values = append(values, v)
					}
				}
			} else {
				values = append(values, c.value)
			}
		}
	}
	return str, values
}

func (h *DbTable) selectCauseAndValues(keys []string, whereCause WhereCause, pagination *apibase.Pagination) (string, []interface{}) {
	query := "SELECT "
	if keys != nil && len(keys) > 0 {
		query += makeSelectors(keys) + " "
	} else {
		query += "* "
	}
	query += "FROM `" + h.TableName + "` "
	values := []interface{}{}
	if whereCause != nil {
		var wc string
		wc, values = whereCause.SQL()
		if wc != "" {
			query += wc
		}
	}
	if pagination != nil {
		query += pagination.AsQuery()
	}
	return query, values
}

func (h *DbTable) countQuery(whereCause WhereCause) (string, []interface{}) {
	query := "SELECT COUNT(*) FROM `" + h.TableName + "` "
	values := []interface{}{}
	if whereCause != nil {
		var wc string
		wc, values = whereCause.SQL()
		if wc != "" {
			query += wc
		}
	}
	return query, values
}

func (h *DbTable) SelectWhere(keys []string, whereCause WhereCause, pagination *apibase.Pagination) (*sqlx.Rows, int, error) {
	query, values := h.selectCauseAndValues(keys, whereCause, pagination)
	count := 0
	if pagination != nil {
		countQuery, countValues := h.countQuery(whereCause)
		if err := h.GetDB().Get(&count, countQuery, countValues...); err != nil {
			return nil, 0, err
		}
	}
	rows, err := h.GetDB().Queryx(query, values...)
	return rows, count, err
}

func (h *DbTable) SelectOneWhere(keys []string, whereCause WhereCause) *sqlx.Row {
	query, values := h.selectCauseAndValues(keys, whereCause, nil)
	return h.GetDB().QueryRowx(query, values...)
}

func (h *DbTable) GetWhere(keys []string, cause WhereCause, dest interface{}) error {
	query, values := h.selectCauseAndValues(keys, cause, nil)
	return h.GetDB().Get(dest, query, values...)
}

type UpdateFields map[string]interface{}

func (up UpdateFields) GetKeys() []string {
	keys := []string{}
	for k, _ := range up {
		keys = append(keys, k)
	}
	return keys
}

func (h *DbTable) InsertWhere(fields UpdateFields) (sql.Result, error) {
	query := "INSERT INTO `" + h.TableName + "`("
	first := true
	quotes := ""
	values := []interface{}{}
	for k, v := range fields {
		if !first {
			query += ","
			quotes += ","
		} else {
			first = false
		}
		query += wrapWhere(k)
		quotes += "?"
		values = append(values, v)
	}
	query += ") VALUES (" + quotes + ") "
	ret, err := h.GetDB().Exec(query, values...)
	return ret, err
}

func getAsInterfaceValueAndType(data interface{}) (reflect.Type, reflect.Value, bool) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Struct {
		data = v.Interface()
	} else if t.Kind() == reflect.Ptr {
		data = v.Elem().Interface()
	} else if t.Kind() != reflect.Interface {
		return reflect.TypeOf(nil), reflect.ValueOf(nil), false
	}
	v = reflect.ValueOf(data)
	t = reflect.TypeOf(data)
	return t, v, true
}

func GetDBColumnNamesFrom(data interface{}, tag string) []string {
	list := []string{}
	t, _, ok := getAsInterfaceValueAndType(data)
	if !ok {
		return list
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(tag)
		if tag != "" {
			list = append(list, tag)
		}
	}
	return list
}

func (h *DbTable) InsertOne(data interface{}) error {
	t, v, ok := getAsInterfaceValueAndType(data)
	if !ok {
		return nil
	}

	fields := UpdateFields{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		value := v.Field(i).Interface()
		tag := f.Tag.Get("db")
		if tag != "" {
			fields[tag] = value
		}
	}
	_, err := h.InsertWhere(fields)
	return err
}

func (h *DbTable) UpdateWhere(whereCause WhereCause, updateFields UpdateFields, insertIfNotExists bool) error {
	if whereCause == nil {
		return fmt.Errorf("Missing where causes")
	}
	var one *sqlx.Row
	if insertIfNotExists {
		one = h.SelectOneWhere(nil, whereCause)
	} else {
		one = h.SelectOneWhere(updateFields.GetKeys(), whereCause)
	}
	if one != nil {
		dest := UpdateFields{}
		if err := one.MapScan(dest); err != nil {
			return err
		}

		update_fields := []string{}
		for k, v := range updateFields {
			existed, ok := dest[k]
			if !ok || !reflect.DeepEqual(v, existed) {
				update_fields = append(update_fields, k)
			}
		}

		if len(update_fields) > 0 {
			query := "UPDATE " + h.TableName + " SET "
			values := []interface{}{}
			for i, f := range update_fields {
				if i > 0 {
					query += ","
				}
				query += wrapWhere(f) + " = ? "
				values = append(values, updateFields[f])
			}
			w, whereValues := whereCause.SQL()
			query += w
			values = append(values, whereValues...)
			_, err := h.GetDB().Exec(query, values...)
			return err
		}

	} else if insertIfNotExists {
		_, err := h.InsertWhere(updateFields)
		return err
	} else {
		//return fmt.Errorf("Not exist!")
	}
	return nil
}
