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

package schema

import (
	"bytes"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt/aes"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"dtstack.com/dtstack/easymatrix/matrix/util"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
)

const (
	SCHEMA_FILE = "schema.yml"
)

var (
	reVar      = regexp.MustCompile(`\${.*?}`)
	reVarLimit = regexp.MustCompile(`^[a-zA-Z_]\w*$`)
)

type ConfigMap map[string]interface{}

type IpRole struct {
	IP       string
	RoleList []string
}

type ServiceAddrStruct struct {
	idToIndex map[uint]uint

	// must sort by field!!!
	Host        []string
	IP          []string
	NodeId      uint // default get first element
	SingleIndex uint
	// 选择的主机的信息
	Select []IpRole
	// 未选择的主机的信息
	UnSelect []IpRole
}

func (sa *ServiceAddrStruct) getIP() string {
	//现在ServiceAddrStruct 必定不为空 所以用 ip 判断
	if sa.IP == nil {
		return ""
	}
	return sa.IP[sa.SingleIndex]
}

func (sa *ServiceAddrStruct) getHost() string {
	if sa.Host == nil {
		return ""
	}
	return sa.Host[sa.SingleIndex]
}

func (sa *ServiceAddrStruct) getIpByNodeID(id uint) (string, error) {
	if sa.idToIndex == nil {
		return "", fmt.Errorf("serviceAddr is nil")
	}
	if i, ok := sa.idToIndex[id]; ok {
		return sa.IP[i], nil
	}
	return "", fmt.Errorf("nodeID %d not found %v", id, sa.idToIndex)
}

func (sa *ServiceAddrStruct) getHostByNodeID(id uint) (string, error) {
	if sa.idToIndex == nil {
		return "", fmt.Errorf("serviceAddr is nil")
	}
	if i, ok := sa.idToIndex[id]; ok {
		return sa.Host[i], nil
	}
	return "", fmt.Errorf("nodeID %d not found", id)
}

type AffinityStruct struct {
	//亲和性
	Affinity []string `config:"affinity" json:",omitempty"`
	// 反亲和性
	AntiAffinity []string `config:"anti_affinity" json:",omitempty"`
}

type ServiceConfig struct {
	ServiceDisplay  string `config:"service_display"`
	Version         string `config:"version"`
	IsJob           bool   `config:"is_job" json:",omitempty"`
	IsDeployIngress bool   `config:"is_deploy_ingress" json:",omitempty"`
	Workload        string `config:"workload" json:"workload,omitempty"`
	Instance        *struct {
		NeedStorage  *bool              `config:"need_storage" json:",omitempty"`
		StorageClass *string            `config:"storage_class" json:",omitempty"`
		Replica      int                `config:"replica" json:",omitempty"`
		Hostname     string             `config:"hostname" json:"hostname,omitempty"`
		EmptyCar     bool               `config:"empty_car" json:",omitempty"`
		UseCloud     bool               `config:"use_cloud" json:",omitempty"`
		ConfigPaths  []string           `config:"config_paths" json:",omitempty"`
		Logs         []string           `config:"logs" json:",omitempty"`
		DataDir      []string           `config:"data_dir" json:",omitempty"`
		Environment  map[string]*string `config:"environment" json:",omitempty"`
		Backup       string             `config:"backup" json:",omitempty"`
		Rollback     string             `config:"rollback" json:",omitempty"`
		HealthCheck  *struct {
			Shell       string `config:"shell" validate:"required"`
			Period      string `config:"period"`
			StartPeriod string `config:"start_period"`
			Timeout     string `config:"timeout"`
			Retries     int    `config:"retries"`
		} `config:"healthcheck" json:",omitempty"`
		ExtendedHealthCheck []*struct {
			Name     string `config:"name" validate:"required"`
			AutoExec bool   `config:"auto_exec" validate:"required"`
			Period   string `config:"period"`
			Retries  int    `config:"retries"`
			Shell    string `config:"shell" validate:"required"`
		} `config:"extended_health_check" json:",omitempty"`
		RunUser           string             `config:"run_user"`
		InstallPath       string             `config:"install_path" json:",omitempty"`
		Cmd               string             `config:"cmd" validate:"required"`
		Image             string             `config:"image" json:",omitempty"`
		HARoleCmd         string             `config:"ha_role_cmd" json:",omitempty"`
		HomePage          string             `config:"home_page" json:",omitempty"`
		Pseudo            bool               `config:"pseudo" json:",omitempty"`
		TestOn            bool               `config:"test_on" json:",omitempty"`
		TestDepends       string             `config:"test_depends" json:",omitempty"`
		TestScript        string             `config:"test_script" json:",omitempty"`
		PostDeploy        string             `config:"post_deploy"`
		PostUpGrade       string             `config:"post_upgrade"`
		PostUndeploy      string             `config:"post_undeploy"`
		UnInstall         string             `config:"uninstall"`
		Ports             []int              `config:"ports" json:",omitempty"`
		PrometheusPort    string             `config:"prometheus_port"`
		MaxReplica        string             `config:"max_replica" json:",omitempty"`
		StartAfterInstall bool               `config:"start_after_install" json:",omitempty"`
		UpdateRecreate    bool               `config:"update_recreate" json:",omitempty"`
		ResourceLimit     map[string]*string `config:"resource_limit" json:",omitempty"`
		ResourceRequest   map[string]*string `config:"resource_request" json:",omitempty"`
		NeedHostAlia      *bool              `config:"need_hostalia" json:",omitempty"`
		HostAlias         *string            `config:"hostalias" json:",omitempty"`
		PluginPath        string             `config:"plugin_path" json:",omitempty"`
		PluginInit        []string           `config:"plugin_init" json:",omitempty"`
		initService       []ServiceConfig
		Switch            map[string]*struct {
			Config       string `config:"config" validate:"required"`
			Desc         string `config:"desc"`
			IsOn         *bool  `config:"is_on" validate:"required"`
			OnScript     string `config:"on_script" validate:"required"`
			OffScript    string `config:"off_script"`
			PostOnScript *struct {
				Type  string `config:"type"`
				Value string `config:"value"`
				Desc  string `config:"desc"`
			} `config:"post_on_script"`
			PostOffScript *struct {
				Type  string `config:"type"`
				Value string `config:"value"`
				Desc  string `config:"desc"`
			} `config:"post_off_script"`
			Extention *struct {
				Type  string `config:"type"`
				Value string `config:"value"`
				Desc  string `config:"desc"`
			} `config:"extention"`
		} `config:"switch" json:",omitempty"`
	} `config:"instance" json:",omitempty"`
	Group              string             `config:"group"`
	DependsOn          []string           `config:"depends_on" json:",omitempty"`
	Relatives          []string           `config:"relatives" json:",omitempty"`
	Config             ConfigMap          `config:"config" json:",omitempty"`
	BaseProduct        string             `config:",ignore"`
	BaseProductVersion string             `config:",ignore"`
	BaseService        string             `config:",ignore"`
	BaseParsed         bool               `config:",ignore"`
	BaseAtrribute      string             `config:",ignore"`
	RelyOn             []string           `config:"rely_on" json:",omitempty"`
	ServiceAddr        *ServiceAddrStruct `config:",ignore" json:",omitempty"`

	serviceTree *nodeService

	//编排
	Orchestration *AffinityStruct `config:"orchestration" json:",omitempty"`
}

func (sc *ServiceConfig) ParseUnInstallScript(baseDir string) (string, error) {
	fileName := strings.Split(sc.Instance.UnInstall, " ")[0]
	tpl, err := template.ParseFiles(filepath.Join(baseDir, fileName))
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = tpl.Option("missingkey=error").Execute(buf, sc.Config); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (sc *ServiceConfig) ParseUnInstallParamter() ([]string, error) {
	splits := strings.Split(sc.Instance.UnInstall, " ")
	argLen := len(splits)
	if argLen <= 1 {
		return nil, nil
	}
	args := splits[1:argLen]
	paramter := []string{}
	for _, arg := range args {
		reg1 := regexp.MustCompile(`\${\w+}`)
		reg2 := regexp.MustCompile(`\w+`)
		vaule := reg1.FindString(arg)
		if vaule == "" {
			paramter = append(paramter, arg) // 没有{}包裹
			continue
		}
		k := reg2.FindString(vaule)
		mp, ok := sc.Config[k]
		if !ok {
			return nil, errors.New("arg key null in scheme")
		}
		paramter = append(paramter, fmt.Sprintf("%s", mp))
	}
	return paramter, nil
}

type SchemaConfig struct {
	ParentProductName  string                   `config:"parent_product_name" validate:"required"`
	ProductName        string                   `config:"product_name" validate:"required"`
	ProductNameDisplay string                   `config:"product_name_display"`
	ProductVersion     string                   `config:"product_version" validate:"required"`
	ProductType        string                   `config:"product_type" json:",omitempty"`
	DeployType         string                   `config:"deploy_type" json:",omitempty"`
	Service            map[string]ServiceConfig `config:"service" validate:"required"`
}

func (sc *SchemaConfig) validateServiceInstance(name string) error {
	svc, exist := sc.Service[name]
	if !exist {
		return fmt.Errorf("can't found service `%v`", name)
	}
	if svc.Instance != nil {
		if svc.Instance.MaxReplica != "" && !reVar.MatchString(svc.Instance.MaxReplica) {
			if _, err := strconv.ParseUint(svc.Instance.MaxReplica, 10, 64); err != nil {
				return fmt.Errorf("service %s's max_replica is not uint", name)
			}
		}
		if svc.Instance.HealthCheck != nil {
			if svc.Instance.HealthCheck.Period == "" {
				svc.Instance.HealthCheck.Period = "1m"
			}
			if svc.Instance.HealthCheck.Retries <= 0 {
				svc.Instance.HealthCheck.Retries = 1
			}

			var period time.Duration
			if !reVar.MatchString(svc.Instance.HealthCheck.Period) {
				var err error
				if period, err = time.ParseDuration(svc.Instance.HealthCheck.Period); err != nil {
					return fmt.Errorf("service %s's health.period is invalid", name)
				} else if period < time.Second {
					return fmt.Errorf("service %s's health.period less than 1s", name)
				}
			}
			if svc.Instance.HealthCheck.StartPeriod != "" && !reVar.MatchString(svc.Instance.HealthCheck.StartPeriod) {
				if startPeriod, err := time.ParseDuration(svc.Instance.HealthCheck.StartPeriod); err != nil {
					return fmt.Errorf("service %s's health.start_period is invalid", name)
				} else if startPeriod < 0 {
					return fmt.Errorf("service %s's health.start_period less than zero", name)
				}
			}
			if svc.Instance.HealthCheck.Timeout != "" && !reVar.MatchString(svc.Instance.HealthCheck.Timeout) {
				if timeout, err := time.ParseDuration(svc.Instance.HealthCheck.Timeout); err != nil {
					return fmt.Errorf("service %s's health.timeout is invalid", name)
				} else if timeout <= 0 || timeout > period {
					svc.Instance.HealthCheck.Timeout = svc.Instance.HealthCheck.Period
				}
			}
		}
		if svc.Instance.PrometheusPort != "" && !reVar.MatchString(svc.Instance.PrometheusPort) {
			if _, err := strconv.ParseUint(svc.Instance.PrometheusPort, 10, 16); err != nil {
				return fmt.Errorf("service %s's prometheus_port is not uint16", name)
			}
		}
		for i, path := range svc.Instance.ConfigPaths {
			svc.Instance.ConfigPaths[i] = strings.Replace(filepath.Clean(path), `\`, `/`, -1)
		}
	}
	return nil
}

func (sc *SchemaConfig) Validate() error {
	if ok, _ := regexp.MatchString(`[/\\]`, sc.ParentProductName); ok {
		return fmt.Errorf("parent product name is invalid")
	}
	if ok, _ := regexp.MatchString(`[/\\]`, sc.ProductName); ok {
		return fmt.Errorf("product name is invalid")
	}
	if ok, _ := regexp.MatchString(`^[\w\-.]+$`, sc.ProductVersion); !ok {
		return fmt.Errorf("product version is invalid")
	}

	baseServices := make([]string, 0)
	for name, svc := range sc.Service {
		if ok, _ := regexp.MatchString(`^[a-zA-Z_]\w*(@[^/\\]+\.*[a-zA-Z_]\w*)?$`, name); !ok {
			return fmt.Errorf("service name %s is invalid", name)
		}

		if err := sc.CheckServiceConfig(svc.Config); err != nil {
			return err
		}

		if svc.Group == "" {
			svc.Group = "default"
			sc.Service[name] = svc
		}

		if names := strings.SplitN(name, "@", 3); len(names) >= 2 && svc.Instance == nil {
			baseServices = append(baseServices, name)
		} else if svc.Version == "" {
			return fmt.Errorf("service %s version not set", name)
		}
		if svc.Version != "" {
			if ok, _ := regexp.MatchString(`^[\w\-.]+$`, svc.Version); !ok {
				return fmt.Errorf("service version %s is invalid", svc.Version)
			}
		}

		if err := sc.validateServiceInstance(name); err != nil {
			return err
		}
	}
	for _, s := range baseServices {
		value := sc.Service[s]
		names := strings.SplitN(s, "@", 3)
		dot := strings.LastIndexByte(names[1], '.')
		value.BaseProduct = names[1][:dot]
		if strings.Contains(value.BaseProduct, "#") {
			products := strings.Split(value.BaseProduct, "#")
			value.BaseProduct = products[0]
			value.BaseProductVersion = products[1]
		}
		value.BaseService = names[1][dot+1:]
		if len(names) > 2 {
			value.BaseAtrribute = names[2]
		}
		delete(sc.Service, s)
		sc.Service[names[0]] = value
	}
	if err := sc.checkDependencies(); err != nil {
		return err
	}
	return nil
}

type nodeService struct {
	name  string
	upper []*nodeService

	parsed bool
}

func (ns *nodeService) findService(name string) *nodeService {
	if ns.name == name {
		return ns
	}
	for _, uns := range ns.upper {
		if ns := uns.findService(name); ns != nil {
			return ns
		}
	}
	return nil
}

type pos struct {
	svcName string
	varName string
	record  map[string]struct{}
}

func (sc *SchemaConfig) parseVar(curPos pos, src *string) (dst interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	rs := reVar.ReplaceAllStringFunc(*src, func(s string) string {
		// s = ${xxxxxxx}
		var svcName, varName string
		var bVarAlone = *src == s

		if names := strings.Split(s[2:len(s)-1], "."); len(names) == 1 {
			if len(names[0]) > 0 && names[0][0] == '@' {
				// check and parse @atName
				atName := names[0][1:]
				if ok := reVarLimit.MatchString(atName); !ok {
					panic(fmt.Errorf("`%v` format error", s))
				}
				if _, exist := sc.Service[atName]; !exist {
					panic(fmt.Errorf("service `%v` not exist", atName))
				}
				if sc.Service[atName].ServiceAddr != nil {
					if bVarAlone {
						dst = sc.Service[atName].ServiceAddr
						return ""
					} else {
						return sc.Service[atName].ServiceAddr.getIP()
					}
				}
				// not set ServiceAddrStruct, just return s = ${@atName}, can't return error
				return s
			}
			// ${varName}
			svcName = curPos.svcName
			varName = names[0]
		} else if len(names) == 2 {
			svcName, varName = names[0], names[1]
		} else {
			panic(fmt.Errorf("`%v` format error", s))
		}
		if !reVarLimit.MatchString(svcName) || !reVarLimit.MatchString(varName) {
			panic(fmt.Errorf("`%v` format error", s))
		}

		if _, exist := sc.Service[svcName]; !exist {
			panic(fmt.Errorf("service `%v` not exist", svcName))
		}
		value, exist := sc.Service[svcName].Config[varName]
		if !exist {
			if sc.Service[svcName].BaseParsed || sc.Service[svcName].BaseService == "" {
				log.Errorf("service `%v` config `%v` not exist", svcName, varName)
				return ""
			} else {
				// BaseService not parsed, just return s = ${@varName}
				return s
			}
		}
		vc := value.(VisualConfig)
		if src, ok := vc.Value.(*string); ok {
			savePos := curPos
			if _, ok := curPos.record[svcName+":"+varName]; !ok {
				curPos.svcName, curPos.varName = svcName, varName
				curPos.record[svcName+":"+varName] = struct{}{}
			} else {
				panic(fmt.Errorf("service `%v` config `%v` recursive variable", svcName, varName))
			}
			dst, err = sc.parseVar(curPos, src)
			delete(curPos.record, svcName+":"+varName)
			curPos = savePos
			if err != nil {
				panic(err)
			}
			vc.Value = dst
			sc.Service[svcName].Config[varName] = vc
			if _, ok := dst.(*string); !ok {
				if bVarAlone == false {
					return vc.String()
				}
				return ""
			}
			return *dst.(*string)
		} else {
			return vc.String()
		}
	})

	if _, ok := dst.(*string); ok || dst == nil {
		dst = &rs
	}
	return
}

func (sc *SchemaConfig) parseVarInstance(svcName string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			return sc.parseVarInstance(svcName, v.Elem())
		}
	case reflect.Slice:
		if !v.IsNil() {
			for i := 0; i < v.Len(); i++ {
				if err := sc.parseVarInstance(svcName, v.Index(i)); err != nil {
					return err
				}
			}
		}
	case reflect.Map:
		if !v.IsNil() {
			for _, k := range v.MapKeys() {
				if err := sc.parseVarInstance(svcName, v.MapIndex(k)); err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := sc.parseVarInstance(svcName, v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.String:
		curPos := pos{svcName: svcName, record: map[string]struct{}{}}
		vv := v.String()
		dst, err := sc.parseVar(curPos, &vv)
		if err != nil {
			return err
		}
		switch d := dst.(type) {
		case *string:
			v.SetString(*d)
		case *ServiceAddrStruct:
			v.SetString(d.getIP())
		}
	}
	return nil
}

func (sc *SchemaConfig) parseVarConfig(svcName string, cfg ConfigMap) error {
	for key, value := range cfg {
		vc := value.(VisualConfig)
		if src, ok := vc.Value.(*string); ok {
			curPos := pos{
				svcName: svcName,
				varName: key,
				record:  map[string]struct{}{svcName + ":" + key: struct{}{}},
			}
			dst, err := sc.parseVar(curPos, src)
			if err != nil {
				return err
			}
			vc.Value = dst
		}

		if src, ok := vc.Default.(*string); ok {
			curPos := pos{
				svcName: svcName,
				varName: key,
				record:  map[string]struct{}{svcName + ":" + key: {}},
			}
			dst, err := sc.parseVar(curPos, src)
			if err != nil {
				return err
			}
			vc.Default = dst
		}

		cfg[key] = vc
	}
	return nil
}

func allElem(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) {
		v = v.Elem()
	}
	return v
}

// pathName like "dtlog.Instance.Environment.MYSQL_ADDRESS", "dtlog.Instance.Logs.0"
func (sc *SchemaConfig) SetField(pathName, field string) (old string, err error) {
	v := reflect.ValueOf(sc.Service)
	for _, name := range strings.Split(pathName, ".") {
		if name == "" {
			err = fmt.Errorf("path `%v` is invalid", pathName)
			return
		}

		switch v.Kind() {
		case reflect.Map:
			if v = allElem(v.MapIndex(reflect.ValueOf(name))); !v.IsValid() {
				err = fmt.Errorf("`%v` not found", name)
				return
			}
		case reflect.Struct:
			if v = allElem(v.FieldByName(name)); !v.IsValid() {
				err = fmt.Errorf("`%v` not found", name)
				return
			}
		case reflect.Slice, reflect.Array:
			var index int
			index, err = strconv.Atoi(name)
			if err != nil {
				err = fmt.Errorf("`%v` not slice or array index", name)
				return
			}
			if index > v.Len() || index < 0 {
				err = fmt.Errorf("`%v` out of range", name)
				return
			}
			v = allElem(v.Index(index))
		default:
			err = fmt.Errorf("`%v` type is %v, not support", name, v.Kind())
			return
		}
	}

	if !v.CanSet() {
		err = fmt.Errorf("path `%v` can't set", pathName)
		return
	}
	if v.Kind() == reflect.String {
		old = v.String()
		v.SetString(field)
	} else if v.Kind() == reflect.Bool {
		old = strconv.FormatBool(v.Bool())
		var b bool
		if b, err = strconv.ParseBool(field); err == nil {
			v.SetBool(b)
		}
	} else if v.Kind() == reflect.Int {
		old = strconv.FormatInt(v.Int(), 10)
		var i int64
		if i, err = strconv.ParseInt(field, 10, 64); err == nil {
			v.SetInt(i)
		}
	} else {
		err = fmt.Errorf("path `%v` is not string/bool", pathName)
	}

	return
}

// SetServiceNodeIP will affect ParseServiceVariable/ParseServiceConfigFiles result
func (sc *SchemaConfig) SetServiceNodeIP(name string, index, nodeId uint, idToIndex map[uint]uint) error {
	svc, ok := sc.Service[name]
	if !ok {
		return fmt.Errorf("service `%v` not exist", name)
	}
	if svc.ServiceAddr.IP == nil {
		return fmt.Errorf("service `%v` not set any ip", name)
	}

	svc.ServiceAddr.idToIndex = idToIndex
	svc.ServiceAddr.SingleIndex = index
	svc.ServiceAddr.NodeId = nodeId
	sc.Service[name] = svc
	return nil
}

func (sc *SchemaConfig) ParseServiceVariable(name string) error {
	svc, ok := sc.Service[name]
	if !ok {
		return fmt.Errorf("service `%v` not exist", name)
	}

	if err := sc.parseVarConfig(name, svc.Config); err != nil {
		return err
	}
	if err := sc.parseVarInstance(name, reflect.ValueOf(svc.Instance)); err != nil {
		return err
	}
	sc.Service[name] = svc
	return sc.validateServiceInstance(name)
}

func (sc *SchemaConfig) ParseVariable() error {
	for name := range sc.Service {
		if err := sc.ParseServiceVariable(name); err != nil {
			return err
		}
	}
	return nil
}

// 获取产品包下有baseProduct和baseService字段的service name
func (sc *SchemaConfig) GetBaseService() []string {
	nameList := make([]string, 0)
	for name, svc := range sc.Service {
		if svc.BaseProduct != "" || svc.BaseService != "" {
			nameList = append(nameList, name)
		}
	}
	return nameList
}

func (sc *SchemaConfig) SetBaseService(name string, config ConfigMap, ips, hosts []string, version string) error {
	svc, ok := sc.Service[name]
	if !ok {
		return fmt.Errorf("service `%v` not exist", name)
	}
	if config != nil {
		if svc.Config == nil && len(config) > 0 {
			svc.Config = make(ConfigMap, len(config))
		}
		for bKey, bValue := range config {
			svc.Config[bKey] = bValue
		}
	}
	svc.BaseParsed = true
	if svc.Version == "" {
		svc.Version = version
	}
	sc.Service[name] = svc

	if svc.Instance == nil {
		return sc.SetServiceAddr(name, ips, hosts)
	}
	return nil
}

// config file will check at baseDir/name/Instance.ConfigPaths
func (sc *SchemaConfig) ParseServiceConfigFiles(baseDir, name string) ([][]byte, error) {
	svc, ok := sc.Service[name]
	if !ok {
		return nil, fmt.Errorf("service `%v` not exist", name)
	}

	cfgs := make([][]byte, 0, len(svc.Instance.ConfigPaths))
	if svc.Instance != nil {
		for _, configPath := range svc.Instance.ConfigPaths {
			tpl, err := template.ParseFiles(filepath.Join(baseDir, name, configPath))
			if err != nil {
				return nil, err
			}
			buf := &bytes.Buffer{}
			if err = tpl.Option("missingkey=error").Execute(buf, svc.Config); err != nil { // Execute svc.Config[serviceName].VisualConfig.String()
				return nil, err
			}
			cfgs = append(cfgs, buf.Bytes())
		}
	}

	return cfgs, nil
}

// config file will check at baseDir/name/Instance.ConfigPaths
func (sc *SchemaConfig) ParseConfigFiles(baseDir string) (map[string][][]byte, error) {
	cfgM := make(map[string][][]byte)

	for svcName, svc := range sc.Service {
		if svc.Instance != nil && len(svc.Instance.ConfigPaths) > 0 {
			cfgs, err := sc.ParseServiceConfigFiles(baseDir, svcName)
			if err != nil {
				return nil, err
			}
			cfgM[svcName] = cfgs
		}
	}
	return cfgM, nil
}

func (sc *SchemaConfig) CheckServiceAddr() error {
	for name, svc := range sc.Service {
		if svc.ServiceAddr == nil || svc.ServiceAddr.IP == nil {
			return fmt.Errorf("service `%v` have not set ip", name)
		}
	}
	return nil
}

func (sc *SchemaConfig) CheckStorage() error {
	for name, svc := range sc.Service {
		if svc.Instance != nil && svc.Instance.StorageClass != nil && *svc.Instance.StorageClass == "" {
			return fmt.Errorf("service %s StorageClass is not set", name)
		}
	}
	return nil
}

func (sc *SchemaConfig) SetEmptyServiceAddr() error {
	for name, svc := range sc.Service {
		if svc.ServiceAddr != nil && svc.ServiceAddr.IP == nil {
			err := sc.SetServiceEmpty(name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sc *SchemaConfig) SetServiceEmpty(name string) error {
	svc, ok := sc.Service[name]
	if !ok {
		return fmt.Errorf("service `%v` not exist", name)
	}

	for key, value := range svc.Config {
		svc.Config[key] = value.(VisualConfig).SetEmpty()
	}
	sc.Service[name] = svc
	return sc.SetServiceAddr(name, []string{"127.0.0.1"}, []string{"localhost"})
}

func (sc *SchemaConfig) SetServiceAddr(name string, ip, host []string) error {
	svc, ok := sc.Service[name]
	if !ok {
		return fmt.Errorf("service `%v` not exist", name)
	}

	if len(ip) == 0 {
		return fmt.Errorf("service `%v` ip len is zero", name)
	}
	if host != nil && len(ip) != len(host) {
		return fmt.Errorf("service `%v` ip/host len not equal", name)
	}
	svc.ServiceAddr = &ServiceAddrStruct{
		IP:        make([]string, len(ip)),
		Host:      make([]string, len(host)),
		NodeId:    1,
		idToIndex: map[uint]uint{1: 0},
	}
	copy(svc.ServiceAddr.IP, ip)
	copy(svc.ServiceAddr.Host, host)
	sc.Service[name] = svc

	return nil
}

func (sc *SchemaConfig) parseDependencies(upper *nodeService, name string) error {
	svc, ok := sc.Service[name]
	if !ok {
		return fmt.Errorf("can't found service `%v`", name)
	}

	mapDependsOn := make(map[string]bool, len(svc.DependsOn))
	for _, lName := range svc.DependsOn {
		if mapDependsOn[lName] {
			return fmt.Errorf("service `%v` dependency duplicate(`%v`)", name, lName)
		}
		mapDependsOn[lName] = true
	}

	self := svc.serviceTree
	if upper != nil {
		self.upper = append(self.upper, upper)
	}
	if self.parsed == true {
		// lower dependencies already parsed, just return
		return nil
	}

	self.parsed = true
	for _, lName := range svc.DependsOn {
		if lName == name {
			return fmt.Errorf("service `%v` dependency self", name)
		}
		if err := sc.parseDependencies(self, lName); err != nil {
			return err
		}
	}
	return nil
}

func (sc *SchemaConfig) checkDependencies() error {
	// init every service's serviceTree
	for name, svc := range sc.Service {
		svc.serviceTree = &nodeService{name: name}
		sc.Service[name] = svc
	}

	// check every service
	for name := range sc.Service {
		if err := sc.parseDependencies(nil, name); err != nil {
			return err
		}
	}

	// check every service recursive dependency
	for name, svc := range sc.Service {
		for _, uns := range svc.serviceTree.upper {
			if uns.findService(name) != nil {
				return fmt.Errorf("service `%v` found recursive dependency", name)
			}
		}
	}

	return nil
}

func (sc *SchemaConfig) CheckRisks(pkgDir string) map[string][]string {
	risks := make(map[string][]string)
	for name, svc := range sc.Service {
		if svc.Instance == nil {
			continue
		}
		if svc.Instance.PostDeploy != "" {
			instanceDir := filepath.Join(pkgDir, name)
			risk := GetRisks(instanceDir, svc.Instance.PostDeploy)
			if len(risk) > 0 {
				risks[name] = append(risks[name], risk...)
			}
		}
		if svc.Instance.PostUndeploy != "" {
			instanceDir := filepath.Join(pkgDir, name)
			risk := GetRisks(instanceDir, svc.Instance.PostUndeploy)
			if len(risk) > 0 {
				risks[name] = append(risks[name], risk...)
			}
		}
	}
	for name := range risks {
		risks[name] = util.ArrayUniqueStr(risks[name])
	}

	return risks
}

// 返回产品包下服务组件的分组列表信息如{UIC:{uic:{ServiceConfig}},default:{mysql:{ServiceConfig}}}
func (sc *SchemaConfig) Group(unchecked []string) map[string]map[string]ServiceConfig {
	serviceGroup := map[string]map[string]ServiceConfig{}
	set := make(map[string]struct{}, len(unchecked))
	for _, us := range unchecked {
		set[us] = struct{}{}
	}
	for name, svc := range sc.Service {
		if _, ok := set[name]; ok {
			continue
		}
		services, ok := serviceGroup[svc.Group]
		if !ok {
			services = map[string]ServiceConfig{}
			serviceGroup[svc.Group] = services
		}
		if svc.ServiceDisplay == "" {
			svc.ServiceDisplay = name
		}
		services[name] = svc
	}
	return serviceGroup
}

func ParseSchemaConfigFile(schemaFile string) (*SchemaConfig, error) {
	configContent, err := yaml.NewConfigWithFile(schemaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema file not found")
		}
		return nil, err
	}

	cfg := SchemaConfig{}
	if err = configContent.Unpack(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ParseSchemaConfigBytes(yamlInput []byte) (*SchemaConfig, error) {
	configContent, err := yaml.NewConfig(yamlInput)
	if err != nil {
		return nil, err
	}

	cfg := SchemaConfig{}
	if err = configContent.Unpack(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Unmarshal(jsonInput []byte) (*SchemaConfig, error) {
	var sc *SchemaConfig
	if err := json.Unmarshal(jsonInput, &sc); err != nil {
		return nil, err
	}

	for name, svc := range sc.Service {
		if svc.Instance != nil && svc.Instance.NeedStorage != nil && *svc.Instance.NeedStorage == true {
			s := ""
			svc.Instance.StorageClass = &s
		}
		if svc.Instance != nil && svc.Instance.NeedHostAlia != nil && *svc.Instance.NeedHostAlia {
			s := ""
			svc.Instance.HostAlias = &s
		}
		for key, value := range svc.Config {
			b, _ := json.Marshal(value)
			vc := VisualConfig{}
			json.Unmarshal(b, &vc)

			switch v := vc.Default.(type) {
			case string:
				vc.Default = &v
			default:
				b, _ := json.Marshal(v)
				addr := ServiceAddrStruct{}
				json.Unmarshal(b, &addr)
				if len(addr.Host) == 0 && len(addr.IP) > 0 {
					addr.Host = nil
				}
				vc.Default = &addr
			}
			switch v := vc.Value.(type) {
			case string:
				vc.Value = &v
			default:
				b, _ := json.Marshal(v)
				addr := ServiceAddrStruct{}
				json.Unmarshal(b, &addr)
				if len(addr.Host) == 0 && len(addr.IP) > 0 {
					addr.Host = nil
				}
				vc.Value = &addr
			}

			svc.Config[key] = vc
		}
		if svc.ServiceAddr != nil && len(svc.ServiceAddr.Host) == 0 && len(svc.ServiceAddr.IP) > 0 {
			svc.ServiceAddr.Host = nil
			sc.Service[name] = svc
		}
	}

	return sc, nil
}

func Clone(org *SchemaConfig) (dst *SchemaConfig, err error) {
	var b []byte
	b, err = json.Marshal(org)
	if err != nil {
		return
	}
	return Unmarshal(b)
}

type VisualConfig struct {
	// must sort by field!!!
	Default interface{} `config:"default"`
	Desc    string      `config:"desc" validate:"required"`
	Type    string      `config:"type" validate:"required"`
	Value   interface{} `config:",ignore"`
}

func (vc VisualConfig) AesEncryptByPassword(adminPass string) {
	defaultValue, err := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(vc.Default.(*string))), adminPass)
	if err != nil {
		log.Errorf("VisualConfig.AesEncryptByPassword %v", err)
	}
	value, _ := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(vc.Value.(*string))), adminPass)
	vc.Default = defaultValue
	vc.Value = value
}

func (vc VisualConfig) String() string {
	switch v := vc.Value.(type) {
	case *ServiceAddrStruct:
		return v.getIP()
	default:
		return *vc.Value.(*string)
	}
}

func (vc VisualConfig) SetEmpty() VisualConfig {
	var empty string
	switch vc.Value.(type) {
	case *ServiceAddrStruct:
		vc.Value = &ServiceAddrStruct{IP: []string{"127.0.0.1"}, Host: []string{"localhost"}}
	default:
		vc.Value = &empty
	}
	return vc
}

func (sc *SchemaConfig) CheckServiceConfig(cfg ConfigMap) error {
	for key, value := range cfg {
		if ok := reVarLimit.MatchString(key); !ok {
			return fmt.Errorf("service config `%v` is invalid", key)
		}
		switch key {
		case "Join", "JoinHost", "Joinx", "JoinxHost", "NodeCount", "NodeIndex", "NodeID", "AddOne", "LastSegIP",
			"Hostname", "IPList", "HostList", "GetIpByNodeID", "GetHostByNodeID":
			return fmt.Errorf("service config `%v` is internal method name", key)
		}

		vc := VisualConfig{}
		def, val := "", ""
		switch value.(type) {
		case map[string]interface{}:
			v, err := ucfg.NewFrom(value)
			if err != nil {
				return err
			}
			if err = v.Unpack(&vc); err != nil {
				return err
			}
			switch vc.Type {
			case "filepath", "port", "ip", "number", "string":
			default:
				return fmt.Errorf("service config `%v` type %s not support", key, vc.Type)
			}
			switch vc.Default.(type) {
			case string:
				def = vc.Default.(string)
			case int64, uint64, float64:
				def = fmt.Sprint(vc.Default)
			case nil:
			default:
				return fmt.Errorf("service config `%v` default type %T not support", key, vc.Default)
			}
		case nil:
			vc.Type = "internal"
			vc.Desc = "internal"
		case string:
			vc.Type = "internal"
			vc.Desc = "internal"
			def = value.(string)
		case int64, uint64, float64:
			vc.Type = "internal"
			vc.Desc = "internal"
			def = fmt.Sprint(value)
		case VisualConfig:
			vc.Type = "internal"
			vc.Desc = "internal"
			def = value.(VisualConfig).String()
		default:
			return fmt.Errorf("service config `%v` struct %T not support", key, value)
		}

		val = def
		vc.Default, vc.Value = &def, &val
		cfg[key] = vc
	}
	return nil
}
