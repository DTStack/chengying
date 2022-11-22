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

package model

import (
	"database/sql"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"errors"
	. "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/wangqi811/gomonkey/v2"
	"reflect"
	"testing"
)

type CustomResult struct {
	AffectedRows int64
	InsertId     int64
}

func (res *CustomResult) LastInsertId() (int64, error) {
	return res.InsertId, nil
}

func (res *CustomResult) RowsAffected() (int64, error) {
	return res.AffectedRows, nil
}

func TestSwitchRecord_GetRecordById(t *testing.T) {
	Convey("Test Get Record by id", t, func() {
		Convey("Test Get record success", func() {
			ctrl := NewController(t)
			defer ctrl.Finish()

			mock := NewMockSwitchRecordI(ctrl)
			mock.EXPECT().GetRecordById(Eq(1)).Return(&SwitchRecordInfo{Id: 1}, nil)
			v, err := mock.GetRecordById(1)
			So(err, ShouldBeNil)
			So(v.Id, ShouldEqual, 1)
		})

		Convey("Test Get record error norow", func() {
			ctrl := NewController(t)
			defer ctrl.Finish()

			mock := NewMockSwitchRecordI(ctrl)
			mock.EXPECT().GetRecordById(Eq(1)).Return(nil, sql.ErrNoRows)
			v, err := mock.GetRecordById(1)
			So(err, ShouldEqual, sql.ErrNoRows)
			So(v, ShouldBeNil)
		})
	})
}

func TestNewSwitchRecord(t *testing.T) {
	Convey("Test new switch record", t, func() {
		Convey("Test new switch record fail", func() {
			record := &dbhelper.DbTable{}
			patches := ApplyMethod(reflect.TypeOf(record), "GetWhere",
				func(_ *dbhelper.DbTable, _ []string, _ dbhelper.WhereCause, _ interface{}) error {
					return sql.ErrNoRows
				})
			defer patches.Reset()

			patches.ApplyMethod(reflect.TypeOf(record), "InsertWhere",
				func(_ *dbhelper.DbTable, _ dbhelper.UpdateFields) (sql.Result, error) {
					return nil, errors.New("error")
				})

			res, err := SwitchRecord.NewSwitchRecord("", "", "", "", "", "", 1, 1)
			So(res, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Test new switch record success", func() {
			record := &dbhelper.DbTable{}
			patches := ApplyMethod(reflect.TypeOf(record), "GetWhere",
				func(_ *dbhelper.DbTable, _ []string, _ dbhelper.WhereCause, _ interface{}) error {
					return sql.ErrNoRows
				})
			defer patches.Reset()

			patches.ApplyMethod(reflect.TypeOf(record), "InsertWhere",
				func(_ *dbhelper.DbTable, _ dbhelper.UpdateFields) (sql.Result, error) {
					res := &CustomResult{InsertId: 1, AffectedRows: 0}
					return res, nil
				})

			res, err := SwitchRecord.NewSwitchRecord("", "", "", "", "", "", 1, 1)
			So(res, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})

		Convey("Test update switch record fail", func() {
			record := &dbhelper.DbTable{}
			patches := ApplyMethod(reflect.TypeOf(record), "GetWhere",
				func(_ *dbhelper.DbTable, _ []string, _ dbhelper.WhereCause, _ interface{}) error {
					return nil
				})
			defer patches.Reset()

			patches.ApplyMethod(reflect.TypeOf(record), "UpdateWhere",
				func(_ *dbhelper.DbTable, _ dbhelper.WhereCause, _ dbhelper.UpdateFields, _ bool) error {
					return errors.New("error")
				})

			res, err := SwitchRecord.NewSwitchRecord("", "", "", "", "", "", 1, 1)
			So(res, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Test update switch record success", func() {
			record := &dbhelper.DbTable{}
			patches := ApplyMethod(reflect.TypeOf(record), "GetWhere",
				func(_ *dbhelper.DbTable, _ []string, _ dbhelper.WhereCause, _ interface{}) error {
					return nil
				})
			defer patches.Reset()

			patches.ApplyMethod(reflect.TypeOf(record), "UpdateWhere",
				func(_ *dbhelper.DbTable, _ dbhelper.WhereCause, _ dbhelper.UpdateFields, _ bool) error {
					return nil
				})

			res, err := SwitchRecord.NewSwitchRecord("", "", "", "", "", "", 1, 1)
			So(res, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestSwitchRecord_GetCurrentSwitchRecord(t *testing.T) {
	Convey("Test Get Current Switch Record", t, func() {
		Convey("Test Get Current Switch Record success", func() {
			record := &dbhelper.DbTable{}
			patches := ApplyMethod(reflect.TypeOf(record), "GetWhere",
				func(_ *dbhelper.DbTable, _ []string, _ dbhelper.WhereCause, _ interface{}) error {
					return nil
				})
			defer patches.Reset()
			res, err := SwitchRecord.GetCurrentSwitchRecord(1, "", "", "")
			So(res, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Test Get Current Switch Record fail", func() {
			record := &dbhelper.DbTable{}
			patches := ApplyMethod(reflect.TypeOf(record), "GetWhere",
				func(_ *dbhelper.DbTable, _ []string, _ dbhelper.WhereCause, _ interface{}) error {
					return sql.ErrNoRows
				})
			defer patches.Reset()
			res, err := SwitchRecord.GetCurrentSwitchRecord(1, "", "", "")
			So(res, ShouldBeNil)
			So(err, ShouldEqual, sql.ErrNoRows)
		})
	})
}
