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

//黑夜给了我黑色的眼睛，专治各种buf(g)；
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"database/sql"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/event"
	"dtstack.com/dtstack/easymatrix/matrix/instance"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"dtstack.com/dtstack/easymatrix/schema"
	errors2 "github.com/juju/errors"
)

const (
	STATUS_CHAN_BUFFER = 16
	STATUS_SUCCESS     = "success"
	STATUS_FAILED      = "failed"
)

const (
	IP_LIST_SEP           = ","
	SLASH                 = "\\"
	BASE_SERVICE_OPTIONAL = "optional"
	BASE_SERVICE_BRIDGE   = "bridge"
)

var (
	BASE_SERVICE_DEFAUL_IPS   = []string{"127.0.0.1"}
	BASE_SERVICE_DEFAUL_HOSTS = []string{"127-0-0-1"}
	PRODUCT_WHITE_EXTENTION   = []string{".sh", ".sql", ".yml", ".xml", ".properties", ".conf", ".yaml"}
	MYSQL_SPECIAL_CHAR        = []string{"_", "\\"}
)

type sqlxer interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

type Servicer interface {
	Start() error
	Stop(stopAgentOptionsTypeArr ...int) error
	RollingRestart() error
	RollingConfigUpdate(clusterId int) error

	GetStartStatusChan() <-chan string
}

type service struct {
	pid         int
	name        string
	ip          []string
	agentId     []string
	status      []string
	schema      []schema.ServiceConfig
	pschema     *schema.SchemaConfig
	info        *model.DeployProductListInfo
	operationId string

	//success;failed
	statusCh chan string
}

func canStart(status string) bool {
	for _, item := range model.OUT_OF_START_STATUS_LIST {
		if status == item {
			return false
		}
	}
	return true
}

func canStop(status string) bool {
	for _, item := range model.OUT_OF_STOP_STATUS_LIST {
		if status == item {
			return false
		}
	}
	return true
}

func NewServicer(pid, clusterId int, name string, operationId string) (Servicer, error) {
	list, err := model.DeployInstanceList.GetInstanceListByPidServiceName(pid, clusterId, name)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("not found `%v` service", name)
	}

	info, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		return nil, err
	}

	newServicer := &service{
		pid:         pid,
		name:        name,
		ip:          []string{},
		agentId:     []string{},
		status:      []string{},
		schema:      []schema.ServiceConfig{},
		info:        info,
		operationId: operationId,
		statusCh:    make(chan string, STATUS_CHAN_BUFFER),
	}
	//兼容免部署配置滚动更新和服务启停
	psc := &schema.SchemaConfig{}
	err = json.Unmarshal(info.Product, psc)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	newServicer.pschema = psc
	for _, instance := range list {
		newServicer.ip = append(newServicer.ip, instance.Ip)
		newServicer.agentId = append(newServicer.agentId, instance.AgentId)
		newServicer.status = append(newServicer.status, instance.Status)
		sc := &schema.ServiceConfig{}
		err := json.Unmarshal(instance.Schema, sc)
		if err != nil {
			log.Errorf("%v", err)
			return nil, err
		}
		newServicer.schema = append(newServicer.schema, *sc)
	}
	return newServicer, nil
}

func (this *service) GetStartStatusChan() <-chan string {
	return this.statusCh
}

func (this *service) checkHealthStatus(wg *sync.WaitGroup, hlRlt *int64) {
	wg.Wait()

	ret := STATUS_SUCCESS
	if *hlRlt > 0 {
		ret = STATUS_FAILED
	}
	select {
	case this.statusCh <- ret:
	default:
	}
}

func (this *service) Start() error {
	wg := sync.WaitGroup{}
	healthWg := sync.WaitGroup{}
	wg.Add(len(this.agentId))
	var startRlt int64  // 0 success, 1 fail
	var healthRlt int64 // 0 success, 1 fail

	for index := range this.agentId {
		go func(index int) {
			var err error
			defer wg.Done()

			if !canStart(this.status[index]) {
				return
			}

			var instancer instance.Instancer
			prodSchema := &schema.SchemaConfig{
				ProductName:    this.info.ProductName,
				ProductVersion: this.info.ProductVersion,
				Service:        map[string]schema.ServiceConfig{this.name: this.schema[index]},
			}

			if instancer, err = instance.NewInstancer(this.pid, this.ip[index], this.name, prodSchema, this.operationId); err != nil {
				log.Errorf("%v", err)
				atomic.StoreInt64(&startRlt, 1)
				return
			}
			err = instancer.Start()
			if err != nil {
				log.Errorf("%v", err)
				atomic.StoreInt64(&startRlt, 1)
				return
			}
			if this.schema[index].Instance.HealthCheck == nil {
				instancer.Clear()
				return
			}

			healthWg.Add(1)
			go func() {
				defer healthWg.Done()
				defer instancer.Clear()

				log.Debugf("waiting service instance(%d) GetStatusChan...", instancer.ID())
				ev := <-instancer.GetStatusChan()
				log.Debugf("end service instance(%d) GetStatusChan", instancer.ID())
				switch ev.Type {
				case event.REPORT_EVENT_HEALTH_CHECK:
					if ev.Data.(*agent.HealthCheck).Failed {
						err := fmt.Errorf("service health check failed")
						log.Errorf("%v", err)
						atomic.StoreInt64(&healthRlt, 1)
					}
				case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
					err := fmt.Errorf("service health check canceled")
					log.Errorf("%v", err)
					atomic.StoreInt64(&healthRlt, 1)
				case event.REPORT_EVENT_INSTANCE_ERROR:
					err := fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
					log.Errorf("%v", err)
					atomic.StoreInt64(&healthRlt, 1)
				}
			}()
		}(index)
	}
	wg.Wait()
	if startRlt > 0 {
		return errors.New("some service instance start err")
	}
	go this.checkHealthStatus(&healthWg, &healthRlt)

	return nil
}

func (this *service) Stop(stopAgentOptionsTypeArr ...int) error {
	wg := sync.WaitGroup{}
	wg.Add(len(this.agentId))

	stopAgentOptionsType := agent.AGENT_STOP_UNRECOVER
	if stopAgentOptionsTypeArr != nil {
		stopAgentOptionsType = stopAgentOptionsTypeArr[0]
	}

	var stopRlt int64 // 0 success, 1 fail
	for index := range this.agentId {
		go func(index int) {
			var err error
			defer wg.Done()

			if !canStop(this.status[index]) {
				return
			}

			var instancer instance.Instancer
			prodSchema := &schema.SchemaConfig{
				ProductName:    this.info.ProductName,
				ProductVersion: this.info.ProductVersion,
				Service:        map[string]schema.ServiceConfig{this.name: this.schema[index]},
			}

			if instancer, err = instance.NewInstancer(this.pid, this.ip[index], this.name, prodSchema, this.operationId); err != nil {
				log.Errorf("%v", err)
				atomic.StoreInt64(&stopRlt, 1)
				return
			}
			defer instancer.Clear()

			err = instancer.Stop(stopAgentOptionsType)
			if err != nil {
				log.Errorf("%v", err)
				atomic.StoreInt64(&stopRlt, 1)
				return
			}
		}(index)
	}
	wg.Wait()
	if stopRlt > 0 {
		return errors.New("some service instance stop err")
	}

	return nil
}

func (this *service) RollingRestart() error {
	for index := range this.agentId {
		instancer, err := instance.NewInstancer(this.pid, this.ip[index], this.name, &schema.SchemaConfig{
			ProductName:    this.info.ProductName,
			ProductVersion: this.info.ProductVersion,
			Service:        map[string]schema.ServiceConfig{this.name: this.schema[index]},
		}, this.operationId)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}

		if canStop(this.status[index]) {
			if err = instancer.Stop(); err != nil {
				this.status[index] = model.INSTANCE_STATUS_STOP_FAIL
				log.Errorf("%v", err)
				instancer.Clear()
				return err
			}
			this.status[index] = model.INSTANCE_STATUS_STOPPED
		}
		if canStart(this.status[index]) {
			if err = instancer.Start(); err != nil {
				this.status[index] = model.INSTANCE_STATUS_RUN_FAIL
				log.Errorf("%v", err)
				instancer.Clear()
				return err
			}
			this.status[index] = model.INSTANCE_STATUS_RUNNING

			if this.schema[index].Instance.HealthCheck == nil {
				instancer.Clear()
				continue
			}

			log.Debugf("waiting service instance(%d) GetStatusChan...", instancer.ID())
			ev := <-instancer.GetStatusChan()
			log.Debugf("end service instance(%d) GetStatusChan", instancer.ID())
			switch ev.Type {
			case event.REPORT_EVENT_HEALTH_CHECK:
				if ev.Data.(*agent.HealthCheck).Failed {
					err = fmt.Errorf("health check failed")
					log.Errorf("%v", err)
					instancer.Clear()
					return err
				}
			case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
				err = fmt.Errorf("health check cancelled")
				log.Errorf("%v", err)
				instancer.Clear()
				return err
			case event.REPORT_EVENT_INSTANCE_ERROR:
				this.status[index] = model.INSTANCE_STATUS_RUN_FAIL
				err = fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
				log.Errorf("%v", err)
				instancer.Clear()
				return err
			}
		}
		instancer.Clear()
	}
	return nil
}

func (this *service) getHostsFromIP(ips []string) (hosts []string, err error) {
	if len(ips) == 0 {
		return
	}

	ipHostMap := make(map[string]string, len(ips))
	// 宿主机模式下节点信息
	ipHostInfo := make([]model.HostInfo, 0)
	query := "SELECT ip, hostname FROM " + model.DeployHostList.TableName
	if err = model.DeployHostList.GetDB().Select(&ipHostInfo, query); err != nil {
		return
	}
	for _, info := range ipHostInfo {
		ipHostMap[info.Ip] = info.HostName
	}
	// k8s模式下节点信息
	ipNodeInfo := make([]model.NodeInfo, 0)
	query = "SELECT ip, hostname FROM " + model.DeployNodeList.TableName
	if err = model.DeployNodeList.GetDB().Select(&ipNodeInfo, query); err != nil {
		return
	}
	// 生成所有主机模式和k8s模式下的节点信息(ip:hostname)
	for _, info := range ipNodeInfo {
		ipHostMap[info.Ip] = info.HostName
	}
	for _, ip := range ips {
		if host, exist := ipHostMap[ip]; exist {
			if host == "" {
				return nil, fmt.Errorf("%v hostname is empty", ip)
			}
			hosts = append(hosts, host)
		} else {
			return nil, fmt.Errorf("%v not found in our system", ip)
		}
	}

	return
}

func (this *service) setSchemaFieldServiceAddr(clusterId int, sc *schema.SchemaConfig, s sqlxer) error {
	var infoList []model.SchemaFieldModifyInfo
	query := "SELECT service_name, field_path, field FROM " + model.DeploySchemaFieldModify.TableName + " WHERE product_name=? AND cluster_id=? "
	if err := model.USE_MYSQL_DB().Select(&infoList, query, sc.ProductName, clusterId); err != nil {
		return err
	}
	for _, modify := range infoList {
		sc.SetField(modify.ServiceName+"."+modify.FieldPath, modify.Field)
	}
	for name, svc := range sc.Service {
		var ipList string
		query = "SELECT ip_list FROM " + model.DeployServiceIpList.TableName + " WHERE product_name=? AND service_name=? AND cluster_id=?"
		if err := s.Get(&ipList, query, sc.ProductName, name, clusterId); err != nil && err != sql.ErrNoRows {
			return err
		}
		if ipList != "" {
			ips := strings.Split(ipList, IP_LIST_SEP)
			var hosts []string
			var err error
			if svc.Instance != nil && !svc.Instance.UseCloud && !svc.BaseParsed {
				// 获取产品包所在的节点地址
				if hosts, err = this.getHostsFromIP(ips); err != nil {
					log.Errorf("get host from ip error: %v")
					hosts = ips
				}
			}
			sc.SetServiceAddr(name, ips, hosts)
		}
	}
	return nil
}

func (this *service) getIdToIndex(clusterId int, productName, serviceName string, serviceIp []string) (map[uint]uint, error) {
	nodes, err := model.GetServiceNodes(clusterId, productName, serviceName)
	if err != nil {
		return nil, err
	}
	idToIndex := make(map[uint]uint, len(nodes))
	for _, n := range nodes {
		idToIndex[n.NodeId] = util.FoundIpIdx(serviceIp, n.Ip)
	}
	return idToIndex, nil
}

func (this *service) inheritBaseService(clusterId int, sc *schema.SchemaConfig, s sqlxer) error {
	var err error
	for _, name := range sc.GetBaseService() {
		baseProduct := sc.Service[name].BaseProduct
		baseService := sc.Service[name].BaseService
		baseAtrri := sc.Service[name].BaseAtrribute
		baseConfigMap, ips, hosts, version, err_ := this.getBaseServicInfo(s, baseProduct, baseService, baseAtrri, clusterId)
		if err_ != nil {
			err = errors2.Wrap(err, fmt.Errorf("base service %v(BaseProduct:%v,  BaseService:%v) error:%v", name, baseProduct, baseService, err_))
			continue
		}

		//log.Debugf("service %v(BaseProduct:%v, BaseService:%v, baseConfigMap:%v, ips:%v, hosts:%v)",
		//	name, baseProduct, baseService, baseConfigMap, ips, hosts)
		sc.SetBaseService(name, baseConfigMap, ips, hosts, version)
	}
	return err
}

func (this *service) getBaseServicInfo(s sqlxer, baseProduct, baseService, baseAttri string, clusterId int) (configMap schema.ConfigMap, ips, hosts []string, version string, err error) {
	var productParsed []byte
	query := fmt.Sprintf("SELECT %v.product_parsed FROM %v LEFT JOIN %v ON %v.id = %v.pid WHERE"+
		" product_name=? AND clusterId=? AND %v.status=?",
		model.DeployClusterProductRel.TableName,
		model.DeployProductList.TableName,
		model.DeployClusterProductRel.TableName,
		model.DeployProductList.TableName,
		model.DeployClusterProductRel.TableName,
		model.DeployClusterProductRel.TableName)

	if err = s.Get(&productParsed, query, baseProduct, clusterId, model.PRODUCT_STATUS_DEPLOYED); err == sql.ErrNoRows {
		err = fmt.Errorf("not found such deployed product")
		if baseAttri == BASE_SERVICE_OPTIONAL {
			configMap = nil
			ips = BASE_SERVICE_DEFAUL_IPS
			hosts = BASE_SERVICE_DEFAUL_HOSTS
			err = nil
		}
		return
	} else if err != nil {
		return
	}

	sc, err := schema.Unmarshal(productParsed)
	if err != nil {
		return
	}

	baseSvc := sc.Service[baseService]
	configMap = baseSvc.Config
	version = baseSvc.Version
	if baseSvc.ServiceAddr.IP != nil {
		ips = baseSvc.ServiceAddr.IP
		hosts = baseSvc.ServiceAddr.Host
	}
	return
}

func (this *service) RollingConfigUpdate(clusterId int) error {
	for index := range this.agentId {
		newSchema, err := schema.Clone(this.pschema)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		//兼容继承/分包
		// 填写该产品包下指定服务组件依赖的的配置信息(即Base*属性信息)
		if err = this.inheritBaseService(clusterId, newSchema, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[RollingConfigUpdate] inheritBaseService warn: %+v", err)
		}
		// 解析deploy_schema_field_modify中服务配置修改记录生成新的schema
		if err = this.setSchemaFieldServiceAddr(clusterId, newSchema, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[RollingConfigUpdate] setSchemaFieldServiceAddr err: %v", err)
			return err
		}
		var node *model.ServiceIpNode
		node, err = model.GetServiceIpNode(clusterId, newSchema.ProductName, this.name, this.ip[index])
		if err != nil {
			log.Errorf("%v", err)
			break
		}
		idToIndex, err := this.getIdToIndex(clusterId, newSchema.ProductName, this.name, this.ip)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		if err = newSchema.SetServiceNodeIP(this.name,
			util.FoundIpIdx(this.ip, this.ip[index]),
			node.NodeId,
			idToIndex); err != nil {
			log.Errorf("%v", err)
			return err
		}
		//兼容存在uncheck部署模式
		if err = newSchema.SetEmptyServiceAddr(); err != nil {
			log.Errorf("%v", err)
			return err
		}

		// 若配置项设置了多个值，在此处替换掉
		multiFields, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndIp(clusterId, newSchema.ProductName, this.name, this.ip[index])
		if err != nil {
			log.Errorf("%v", err)
			break
		}
		for _, multiField := range multiFields {
			newSchema.SetField(multiField.ServiceName+"."+multiField.FieldPath, multiField.Field)
		}

		if err := newSchema.ParseServiceVariable(this.name); err != nil {
			log.Errorf("%v", err)
			return err
		}
		instancer, err := instance.NewInstancer(this.pid, this.ip[index], this.name, newSchema, this.operationId)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		if err := instancer.UpdateConfig(); err != nil {
			log.Errorf("%v", err)
			instancer.Clear()
			return err
		}
		instancer.Clear()
	}
	return nil
}
