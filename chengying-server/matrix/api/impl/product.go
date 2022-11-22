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

package impl

import (
	"archive/tar"
	"bytes"
	sysContext "context"
	"database/sql"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt/aes"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/matrix/model/upgrade"
	"dtstack.com/dtstack/easymatrix/matrix/workload"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/discover"
	"dtstack.com/dtstack/easymatrix/matrix/event"
	"dtstack.com/dtstack/easymatrix/matrix/group"
	"dtstack.com/dtstack/easymatrix/matrix/harole"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/instance"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"dtstack.com/dtstack/easymatrix/schema"
	"github.com/better0332/zip"
	"github.com/jmoiron/sqlx"
	errors2 "github.com/juju/errors"
	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
)

const (
	IP_LIST_SEP           = ","
	EVENT_LINE_CHAR       = "-"
	EVENT_CONTENT_CHAR    = "|"
	SLASH                 = "\\"
	BASE_SERVICE_OPTIONAL = "optional"
	BASE_SERVICE_BRIDGE   = "bridge"
	LOG_LINE_SEPRATOR     = "*************************** %v.%s ***************************"
	EVENT_LINE_SEPRATOR   = "+------------+---------------- %v --------------------------+"
	EVENT_CONTENT_URL     = "api/v2/instance/%v/event?eventId=%v"
	TEST_SET_SERVICE_NAME = "testSet"
)

const (
	NORMAL    = "正常"
	ABNORMAL  = "异常"
	RUNNING   = "running"
	HEALTHY   = 1
	UNHEALTHY = -1
)

var (
	BASE_SERVICE_DEFAUL_IPS   = []string{"127.0.0.1"}
	BASE_SERVICE_DEFAUL_HOSTS = []string{"127-0-0-1"}
	PRODUCT_WHITE_EXTENTION   = []string{".sh", ".sql", ".yml", ".xml", ".properties", ".conf", ".yaml", ".tar", ".txt"}
	MYSQL_SPECIAL_CHAR        = []string{"_", "\\"}
)

type sqlxer interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

var contextCancelMapMutex sync.Mutex
var contextCancelMap = map[uuid.UUID]sysContext.CancelFunc{}
var rateLimit = 20

func init() {
	envRateLimit := os.Getenv("Rolling_Update_Rate_Limit")
	if envRateLimit != "" {
		envRateLimitNum, err := strconv.Atoi(envRateLimit)
		if err == nil {
			rateLimit = envRateLimitNum
		}
	}
}

func sqlTransfer(input string) string {
	for _, c := range MYSQL_SPECIAL_CHAR {
		if strings.Contains(input, c) {
			input = strings.Replace(input, c, SLASH+c, -1)
			return input
		}
	}
	return input
}

func checkSmoothUpgradeServiceAddr(clusterId int, productName, productVersion string, isFinalUpgrade bool) error {
	//校验mysql地址更改
	var mysqlServiceName = "mysql"
	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		log.Errorf("get from deployproductlist err: %v", err)
		return err
	}
	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("Unmarshal err: %v", err)
		return err
	}
	if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
		log.Errorf("inheritBaseService err: %+v", err)
		return err
	}
	if err := setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
		log.Errorf("setSchemaFieldServiceAddr err: %v", err)
		return err
	}
	mysqlIPList, err := model.DeployMysqlIpList.GetMysqlIpList(clusterId, productName)
	if err != nil {
		log.Errorf("get mysql ip list error:%v", err)
		return err
	}
	var oldMysqlIpList []string
	if sc.Service[mysqlServiceName].ServiceAddr != nil {
		oldMysqlIpList = sc.Service[mysqlServiceName].ServiceAddr.IP
	}
	//平滑升级中提示修改数据库地址
	if !isFinalUpgrade && CompareIpList(mysqlIPList, oldMysqlIpList) {
		return fmt.Errorf("服务 `%v` 未完善资源分配", mysqlServiceName)
	}
	//最终一次平滑升级提示修改数据库地址
	if isFinalUpgrade && !CompareIpList(mysqlIPList, oldMysqlIpList) {
		return fmt.Errorf("服务 `%v` 未完善资源分配", mysqlServiceName)
	}

	//校验平滑升级的服务编排
	_, err = model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("%v", err)
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		smoothUpgradeList, err := upgrade.SmoothUpgrade.GetByProductName(productName)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		for _, info := range smoothUpgradeList {
			SelectIpList, err := model.DeployServiceIpList.GetServiceIpList(clusterId, productName, info.ServiceName)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
			var UnselectIpList []string
			instanceList, _ := model.DeployInstanceList.GetInstanceBelongService(productName, info.ServiceName, clusterId)
			for _, instance := range instanceList {
				UnselectIpList = append(UnselectIpList, instance.Ip)
			}
			if len(SelectIpList) == len(UnselectIpList) {
				return fmt.Errorf("服务 `%v` 未完善资源分配", info.ServiceName)
			}
		}
	}
	return nil
}

func inheritBaseService(clusterId int, sc *schema.SchemaConfig, s sqlxer) error {
	var err error
	for _, name := range sc.GetBaseService() {
		baseProduct := sc.Service[name].BaseProduct
		baseService := sc.Service[name].BaseService
		baseAtrri := sc.Service[name].BaseAtrribute
		// 获取该产品包下服务组件依赖的其它产品包下的服务组件的基本信息(服务组件的配置信息、ip、host、版本信息)
		baseConfigMap, ips, hosts, version, _, err_ := getBaseServicInfo(s, baseProduct, baseService, baseAtrri, "", clusterId)
		if err_ != nil {
			err = errors2.Wrap(err, fmt.Errorf("base service %v(BaseProduct:%v,  BaseService:%v) error:%v", name, baseProduct, baseService, err_))
			continue
		}

		//log.Debugf("service %v(BaseProduct:%v, BaseService:%v, baseConfigMap:%v, ips:%v, hosts:%v)",
		//	name, baseProduct, baseService, baseConfigMap, ips, hosts)
		//将上面获取依赖服务组件的相关信息覆盖到该组件对应的配置属性处
		sc.SetBaseService(name, baseConfigMap, ips, hosts, version)

	}
	return err
}

func replaceServiceConfig(sc *schema.SchemaConfig, serviceName string, values map[string]interface{}) {
	newSvc := sc.Service[serviceName]

	newConfig := make(schema.ConfigMap, len(values)+len(newSvc.Config))

	for key, value := range values {
		newConfig[key] = value
	}
	for okey, ovalue := range newSvc.Config {
		newConfig[okey] = ovalue
	}
	newSvc.Config = newConfig

	sc.Service[serviceName] = newSvc
}

// get service config that is final parsed. that means the service you want to get, it's config is all parsed.
// if the service you want to get is not deployed, the service is not parsed yet.
// And if it is optional, service's config will be set to nil,ips and hosts will be set to default
func getBaseServicInfo(s sqlxer, baseProduct, baseService, baseAttri string, relynamespace string, clusterId int) (configMap schema.ConfigMap, ips, hosts []string, version string, service *schema.ServiceConfig, err error) {
	var productParsed []byte
	// 检索获取deploy_cluster_product_rel表的product_parsed字段
	query := fmt.Sprintf("SELECT %s.product_parsed FROM %s LEFT JOIN %s ON %s.id = %s.pid WHERE"+
		" product_name=? AND clusterId=? AND %s.status=? AND %s.namespace=?",
		model.DeployClusterProductRel.TableName,
		model.DeployProductList.TableName,
		model.DeployClusterProductRel.TableName,
		model.DeployProductList.TableName,
		model.DeployClusterProductRel.TableName,
		model.DeployClusterProductRel.TableName,
		model.DeployClusterProductRel.TableName)
	log.Infof("get basic info sql: %s", query)

	if err = s.Get(&productParsed, query, baseProduct, clusterId, model.PRODUCT_STATUS_DEPLOYED, relynamespace); err == sql.ErrNoRows {
		if baseAttri != BASE_SERVICE_OPTIONAL {
			//若依赖关系为强依赖，则返回查询不到已部署的组件包
			err = fmt.Errorf("not found such deployed product %s", baseProduct)
			return
		}
	} else if err != nil {
		return
	}
	//若依赖关系为弱依赖（optional）且未查询到已部署的依赖组件，去除部署条件再次查询，查询是否存在部署失败的组件。若不存在返回默认
	if err == sql.ErrNoRows && baseAttri == BASE_SERVICE_OPTIONAL {
		query = fmt.Sprintf("SELECT %s.product_parsed FROM %s LEFT JOIN %s ON %s.id = %s.pid WHERE"+
			" product_name=? AND clusterId=? AND %s.namespace=?",
			model.DeployClusterProductRel.TableName,
			model.DeployProductList.TableName,
			model.DeployClusterProductRel.TableName,
			model.DeployProductList.TableName,
			model.DeployClusterProductRel.TableName,
			model.DeployClusterProductRel.TableName)
		log.Infof("get optional basic info sql: %s", query)
		if err = s.Get(&productParsed, query, baseProduct, clusterId, relynamespace); err == sql.ErrNoRows {
			configMap = nil
			ips = BASE_SERVICE_DEFAUL_IPS
			hosts = BASE_SERVICE_DEFAUL_HOSTS
			err = nil
			return
		} else if err != nil {
			return
		}
	}
	sc, err := schema.Unmarshal(productParsed)
	if err != nil {
		return
	}
	// maybe the product is exist,but the service is not exist. it is not forced.
	baseSvc, exist := sc.Service[baseService]
	if !exist && baseAttri == BASE_SERVICE_OPTIONAL {
		configMap = nil
		ips = BASE_SERVICE_DEFAUL_IPS
		hosts = BASE_SERVICE_DEFAUL_HOSTS
		err = nil
		return
	}
	configMap = baseSvc.Config
	version = baseSvc.Version
	if baseSvc.ServiceAddr.IP != nil {
		ips = baseSvc.ServiceAddr.IP
		hosts = baseSvc.ServiceAddr.Host
	}
	// 获取插件产品包部署后服务组件的信息
	if baseSvc.Workload == workload.PLUGIN {
		service = &baseSvc
	}
	return
}

func setSchemaFieldServiceAddr(clusterId int, sc *schema.SchemaConfig, s sqlxer, namespace string) error {
	var infoList []model.SchemaFieldModifyInfo
	query := "SELECT service_name, field_path, field FROM " + model.DeploySchemaFieldModify.TableName + " WHERE product_name=? AND cluster_id=? AND namespace=?"
	if err := model.USE_MYSQL_DB().Select(&infoList, query, sc.ProductName, clusterId, namespace); err != nil {
		return fmt.Errorf("query deploySchemaFieldModify error: %s", err)
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
				if hosts, err = getHostsFromIP(ips); err != nil {
					log.Errorf("get host from ip error: %v", err)
					hosts = ips
				}
			}
			sc.SetServiceAddr(name, ips, hosts)
		}
	}

	//listByClusterId, err := model.DeployHostList.GetHostListByClusterId(clusterId)
	//if err != nil {
	//	return err
	//}
	//
	//IpRoleMap := make(map[string]schema.IpRole)
	//for _, hInfo := range listByClusterId {
	//	if hInfo.RoleList.Valid && strings.TrimSpace(hInfo.RoleList.String) != "" {
	//		roleNameList, err := model.HostRole.GetRoleNameListStrByIdList(hInfo.RoleList.String)
	//		if err != nil {
	//			return err
	//		}
	//		IpRoleMap[hInfo.Ip] = schema.IpRole{
	//			IP:       hInfo.Ip,
	//			RoleList: roleNameList,
	//		}
	//	} else {
	//		IpRoleMap[hInfo.Ip] = schema.IpRole{
	//			IP:       hInfo.Ip,
	//			RoleList: nil,
	//		}
	//	}
	//}
	//for name, svc := range sc.Service {
	//	//每次都深拷贝 因为有 delete map操作
	//	deepCopyIpRoleMap := make(map[string]schema.IpRole)
	//	for k, v := range IpRoleMap {
	//		deepCopyIpRoleMap[k] = v
	//	}
	//
	//	var ipList string
	//	query := "SELECT ip_list FROM " + model.DeployServiceIpList.TableName + " WHERE product_name=? AND service_name=? AND cluster_id=? AND namespace=?"
	//	if err := s.Get(&ipList, query, sc.ProductName, name, clusterId, namespace); err != nil && err != sql.ErrNoRows {
	//		return fmt.Errorf("query deployServiceIpList error: %s", err)
	//	}
	//
	//	if ipList != "" {
	//		ips := strings.Split(ipList, IP_LIST_SEP)
	//		var hosts []string
	//		var err error
	//		if svc.Instance != nil && !svc.Instance.UseCloud && !svc.BaseParsed {
	//			if hosts, err = getHostsFromIP(ips); err != nil {
	//				log.Errorf("get host from ip error: %v", err)
	//				hosts = ips
	//			}
	//		}
	//		sc.SetServiceAddr(name, ips, hosts)
	//
	//	}
	//	//无论有没有 ip，都要设置 role info  因为 select 与 unselect 自动部署需要回显
	//	if sc.Service[name].ServiceAddr != nil {
	//		err = SetAddrWithRoleInfo(name, sc, deepCopyIpRoleMap, ipList)
	//		if err != nil {
	//			return err
	//		}
	//	} else {
	//		svc.ServiceAddr = &schema.ServiceAddrStruct{
	//			Host:        nil,
	//			IP:          nil,
	//			NodeId:      0,
	//			SingleIndex: 0,
	//			Select:      nil,
	//			UnSelect:    nil,
	//		}
	//		sc.Service[name] = svc
	//		err = SetAddrWithRoleInfo(name, sc, deepCopyIpRoleMap, ipList)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//}
	return nil
}

//reflush base cluster service ips and hosts
//surpport kubernetes mode with base cluster service of cloud mode
func setBaseServiceAddr(clusterId int, sc *schema.SchemaConfig, s sqlxer) error {
	var infoList []model.SchemaFieldModifyInfo
	if clusterId <= 0 {
		//do nothing
		return nil
	}
	query := "SELECT service_name, field_path, field FROM " + model.DeploySchemaFieldModify.TableName + " WHERE cluster_id=?"
	if err := model.USE_MYSQL_DB().Select(&infoList, query, clusterId); err != nil {
		return err
	}
	for _, modify := range infoList {
		if _, ok := sc.Service[modify.ServiceName]; ok {
			sc.SetField(modify.ServiceName+"."+modify.FieldPath, modify.Field)
		}
	}
	for name, svc := range sc.Service {
		if svc.BaseAtrribute != BASE_SERVICE_BRIDGE {
			continue
		}
		var ipList string
		query = "SELECT ip_list FROM " + model.DeployServiceIpList.TableName + " WHERE service_name=? AND cluster_id=? order by id desc"
		if err := s.Get(&ipList, query, name, clusterId); err != nil && err != sql.ErrNoRows {
			return err
		}
		if ipList != "" {
			ips := strings.Split(ipList, IP_LIST_SEP)
			var hosts []string
			var err error
			if svc.Instance != nil && !svc.Instance.UseCloud && !svc.BaseParsed {
				if hosts, err = getHostsFromIP(ips); err != nil {
					log.Errorf("get host from ip error: %v", err)
					hosts = ips
				}
			}
			sc.SetServiceAddr(name, ips, hosts)
		}
	}
	return nil
}

type zipFileWriter struct {
	file   *os.File
	writer *zip.Writer
}

var uploadLock sync.Mutex

// Upload
// @Description  	upload package
// @Summary      	上传产品包
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           package body string true "-F 'package=@Trino_0.359-tdh_centos7_x86_64.tar'"
// @Success         200 {object} string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/upload [post]
func Upload(ctx context.Context) (rlt apibase.Result) {
	uploadLock.Lock()
	defer uploadLock.Unlock()

	log.Debugf("[Upload] uploading product package")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		userId = 0
		log.Infof("upload no userId: %v", err.Error())
	}
	log.Infof("upload userId:  %v", userId)
	f, _, err := ctx.FormFile("package")
	if err != nil {
		return err
	}
	defer f.Close()
	return UnzipAndParse(f, userId)
}

func UnzipAndParse(f io.Reader, userId int) (rlt apibase.Result) {
	var sc *schema.SchemaConfig
	var patch schema.Patch
	var isPatch bool
	var pkgPath string
	var zipFileWriterMap = map[string]*zipFileWriter{}
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		for _, zfw := range zipFileWriterMap {
			if zfw.file != nil {
				zfw.file.Close()
			}
		}

		if _, ok := rlt.(error); ok {
			tx.Rollback()
			log.Errorf("deploy error(rollback): %v", rlt)
		}
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("deploy panic(rollback): %v", r)
			rlt = r
		}
	}()
	var pid int64
	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Debugf("file header: %+v", hdr)
		if sc == nil {
			if hdr.Name != schema.SCHEMA_FILE {
				return fmt.Errorf("tar package first file is %v, not %v", hdr.Name, schema.SCHEMA_FILE)
			}
			buf := bytes.NewBuffer(make([]byte, 0, hdr.Size))
			if _, err = io.Copy(buf, tr); err != nil {
				return err
			}
			if sc, err = schema.ParseSchemaConfigBytes(buf.Bytes()); err != nil {
				return err
			}
			log.Debugf("%v content: \n%s", schema.SCHEMA_FILE, buf.Bytes())

			if sc.ParentProductName == "" {
				sc.ParentProductName = sc.ProductName
			}

			//如果不同的 ParentProduct 下存在同名的 productName 则返回error
			var productName string
			productNameUniqueSql := "select product_name from " + model.DeployProductList.TableName + " where parent_product_name != ? and product_name = ? and product_type = ? limit 1"
			if err = model.DeployInstanceList.GetDB().Get(&productName, productNameUniqueSql, sc.ParentProductName, sc.ProductName, sc.ProductType); err == nil {
				log.Errorf("product_name must be unique， product_name=%v is already exist", productName)
				return fmt.Errorf("product_name must be unique， product_name=%v is already exist", productName)
			} else if err != sql.ErrNoRows {
				return err
			}

			product, err := json.Marshal(sc)
			if err != nil {
				return err
			}

			if sc.ProductNameDisplay == "" {
				sc.ProductNameDisplay = sc.ProductName
			}
			var productType int
			if sc.ProductType == "kubernetes" {
				productType = 1
			}

			info := model.DeployProductListInfo{
				DeployUUID:         "",
				ProductParsed:      []byte(""),
				ParentProductName:  sc.ParentProductName,
				ProductName:        sc.ProductName,
				ProductNameDisplay: sc.ProductNameDisplay,
				ProductVersion:     sc.ProductVersion,
				Product:            product,
				Status:             model.PRODUCT_STATUS_UNDEPLOYED,
				Schema:             product,
				UserId:             userId,
				ProductType:        productType,
			}

			query := "INSERT INTO " + model.DeployProductList.TableName +
				" (deploy_uuid,product_parsed,parent_product_name, product_name, product_name_display, product_version, product, is_current_version, `status`, `schema`, `user_id`, `product_type`) VALUES" +
				" (:deploy_uuid,:product_parsed,:parent_product_name, :product_name, :product_name_display, :product_version, :product, :is_current_version, :status, :schema, :user_id, :product_type)"
			var ret sql.Result
			if ret, err = tx.NamedExec(query, &info); err != nil {
				return err
			}
			pid, _ = ret.LastInsertId()
			//if err = inheritBaseService(sc, tx); err != nil {
			//	log.Errorf("Upload inheritBaseService warn: %+v", err)
			//}
			if err = sc.ParseVariable(); err != nil {
				return err
			}
			pkgPath = filepath.Join(base.WebRoot, sc.ProductName, sc.ProductVersion)
			if err = os.MkdirAll(pkgPath, 0755); err != nil {
				return err
			}
			continue
		}

		if !isPatch && hdr.Name == schema.PATCH_FILE {
			if err = gob.NewDecoder(tr).Decode(&patch); err != nil {
				return err
			}
			log.Debugf("patch: %+v", patch)

			if patch.ProductName != sc.ProductName || patch.NewProductVersion != sc.ProductVersion {
				return fmt.Errorf("the patch and schema info not consistent")
			}

			var id int
			query := "SELECT id FROM " + model.DeployProductList.TableName + " WHERE product_name=? AND product_version=?"
			if err = tx.Get(&id, query, patch.ProductName, patch.OldProductVersion); err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("product %v not found old version %v", patch.ProductName, patch.OldProductVersion)
				}
				return err
			}
			zipPattern := filepath.Join(base.WebRoot, patch.ProductName, patch.OldProductVersion, "*.zip")
			oldProductZips, _ := filepath.Glob(zipPattern)
			if len(oldProductZips) == 0 {
				return fmt.Errorf("not found any zip at %v", zipPattern)
			}
			for _, zipPath := range oldProductZips {
				svcName := strings.TrimSuffix(filepath.Base(zipPath), ".zip")
				newZipPath := filepath.Join(base.WebRoot, patch.ProductName, patch.NewProductVersion, svcName+".zip")
				if patch.IsDiffService(svcName) {
					log.Debugf("diff service %v: CopyFile from %v to %v", svcName, zipPath, newZipPath)
					if _, err = util.CopyFile(zipPath, newZipPath); err != nil {
						return err
					}
				} else if !patch.IsDeletedService(svcName) {
					log.Debugf("no change service %v: Link from %v to %v", svcName, zipPath, newZipPath)
					os.Remove(newZipPath)
					if err = os.Link(zipPath, newZipPath); err != nil {
						return err
					}
				}
			}
			isPatch = true
			continue
		}

		if isPatch {
			if i := strings.IndexByte(hdr.Name, '/'); i > 0 {
				svcName := hdr.Name[:i]
				targetFile := hdr.Name[i+1:]
				newZipPath := filepath.Join(base.WebRoot, patch.ProductName, patch.NewProductVersion, svcName+".zip")
				log.Debugf("append %v(type %c) to zip %v", hdr.Name, hdr.Typeflag, newZipPath)

				zfw, exist := zipFileWriterMap[svcName]
				if !exist {
					zfw = &zipFileWriter{}
					zipFileWriterMap[svcName] = zfw
					newZipFile, err := os.OpenFile(newZipPath, os.O_RDWR, 0)
					if err != nil {
						return fmt.Errorf("open err: %v", err)
					}
					zfw.file = newZipFile
					zipSize, _ := newZipFile.Seek(0, io.SeekEnd)
					zr, err := zip.NewReader(newZipFile, zipSize)
					if err != nil {
						return fmt.Errorf("zip NewReader err: %v", err)
					}
					zfw.writer = zr.Append(newZipFile)
				}
				zfh, _ := zip.FileInfoHeader(hdr.FileInfo())
				zfh.Name = targetFile
				if hdr.Typeflag == tar.TypeDir {
					zfh.Name += "/"
				}
				zfh.Method = zip.Deflate
				w, err := zfw.writer.CreateHeader(zfh)
				if err != nil {
					return fmt.Errorf("zip writer create err: %v", err)
				}
				if hdr.Typeflag == tar.TypeReg || hdr.Typeflag == tar.TypeSymlink {
					if _, err = io.Copy(w, tr); err != nil {
						return fmt.Errorf("zip writer copy err: %v", err)
					}
				}
				continue
			} else if hdr.Typeflag == tar.TypeDir {
				// i == 0
				continue
			}
		}

		log.Infof("decompress %v", hdr.Name)
		fh, err := os.Create(filepath.Join(pkgPath, hdr.Name))
		if err != nil {
			return err
		}
		if _, err = io.Copy(fh, tr); err != nil {
			fh.Close()
			return err
		}
		fh.Close()
	}

	if sc == nil {
		return fmt.Errorf("can't get %v in tar package", schema.SCHEMA_FILE)
	}
	if isPatch {
		for _, diff := range patch.DiffServices {
			for _, deleteFile := range diff.DeletedFiles {
				newZipPath := filepath.Join(base.WebRoot, patch.ProductName, patch.NewProductVersion, diff.ServiceName+".zip")
				log.Debugf("delete %v from zip %v", deleteFile, newZipPath)

				zfw, exist := zipFileWriterMap[diff.ServiceName]
				if !exist {
					zfw = &zipFileWriter{}
					zipFileWriterMap[diff.ServiceName] = zfw
					newZipFile, err := os.OpenFile(newZipPath, os.O_RDWR, 0)
					if err != nil {
						return fmt.Errorf("open err: %v", err)
					}
					zfw.file = newZipFile
					zipSize, _ := newZipFile.Seek(0, io.SeekEnd)
					zr, err := zip.NewReader(newZipFile, zipSize)
					if err != nil {
						return fmt.Errorf("zip NewReader err: %v", err)
					}
					zfw.writer = zr.Append(newZipFile)
				}
				zfw.writer.Delete(deleteFile)
			}
		}
		for svcName, zfw := range zipFileWriterMap {
			if zfw.writer != nil {
				if err := zfw.writer.Close(); err != nil {
					return fmt.Errorf("service %v zip close err: %v", svcName, err)
				}
			}
			if zfw.file != nil {
				zfw.file.Close()
			}
			delete(zipFileWriterMap, svcName)
		}
	}

	tmpDir, err := ioutil.TempDir("", sc.ProductName)
	defer os.RemoveAll(tmpDir)
	for name, svc := range sc.Service {
		// set fake ip just for check ConfigFiles
		sc.SetServiceAddr(name, []string{"127.0.0.1"}, []string{"localhost"})

		if svc.Instance == nil {
			continue
		}
		err = func() error {
			zipFile := filepath.Join(pkgPath, name+".zip")
			log.Infof("check zip package %v", zipFile)
			r, err := zip.OpenReader(zipFile)
			if err != nil {
				return err
			}
			defer r.Close()

			var configPathMap = make(map[string]struct{}, len(svc.Instance.ConfigPaths))
			for _, configPath := range svc.Instance.ConfigPaths {
				configPathMap[configPath] = struct{}{}
			}

			for _, zf := range r.File {
				for configPath := range configPathMap {
					if zf.Name == configPath {
						configFile := filepath.Join(pkgPath, name, zf.Name)
						if err = os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
							return err
						}
						rc, err := zf.Open()
						if err != nil {
							return err
						}
						fh, err := os.Create(configFile)
						if err != nil {
							rc.Close()
							return err
						}
						if _, err = io.Copy(fh, rc); err != nil {
							fh.Close()
							rc.Close()
							return err
						}
						fh.Close()
						rc.Close()

						delete(configPathMap, configPath)
						break
					}
				}
			}
			// 这里只解压常见文本文件，防止解压大文件;同时也是危险检测的目标
			util.UnZipCore(r, filepath.Join(pkgPath, name), PRODUCT_WHITE_EXTENTION)

			return nil
		}()
		if err != nil {
			return err
		}
	}
	// parse again because SetServiceAddr
	if err = sc.ParseVariable(); err != nil {
		return err
	}
	// check all ConfigFiles once
	if _, err = sc.ParseConfigFiles(pkgPath); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if risks := sc.CheckRisks(pkgPath); len(risks) > 0 {
		msg := ""
		for name, list := range risks {
			msg += fmt.Sprintf("服务%s, 存在危险脚本%s; ", name, strings.Join(list, ", "))
		}
		return msg
	} else {
		return map[string]interface{}{
			"pid":            pid,
			"productName":    sc.ProductName,
			"productVersion": sc.ProductVersion,
		}
	}
}

//4.0.9_Rel, 4.0.11_Rel
// a > b return 1
// a <= b return -1
// a is compare target
// 4.0.20_beta vs 4.1_beta
func CompareVersion(a, b string) int {
	//TODO judge version regex
	version_as := strings.Split(a, "_")
	version_bs := strings.Split(b, "_")

	version_as_s := strings.Split(version_as[0], ".")
	version_bs_s := strings.Split(version_bs[0], ".")

	for index, _ := range version_as_s {
		version_a, err := strconv.Atoi(version_as_s[index])
		if err != nil {
			log.Errorf("%v", err)
			return -1
		}
		if len(version_bs_s) > index {
			version_b, err := strconv.Atoi(version_bs_s[index])
			if err != nil {
				log.Errorf("%v", err)
				return -1
			}
			if version_a > version_b {
				return 1
			}
			if version_a == version_b {
				continue
			}
			if version_a < version_b {
				return -1
			}
		} else {
			//if len(target) > len(a)
			return 1
		}

	}
	return -1
}

//定制升级，包含Sql模块的支持 "升级" 模式
//Sql组件服务名模式{XXXXXSql}
func checkCanUpgrade(sc *schema.SchemaConfig) bool {
	haveSql := false
	//var serviceConfig schema.ServiceConfig
	//var serviceName string
	for name, _ := range sc.Service {
		if strings.HasSuffix(name, "Sql") {
			haveSql = true
			//serviceConfig = config
			//serviceName = name
			break
		}
	}
	//不存在sql模块
	if !haveSql {
		return false
	}
	// kubernetes not surpoort upgrade
	if sc.ProductType == "kubernetes" || sc.ProductType == "1" {
		return false
	}
	products, _ := model.DeployProductList.GetProductListByNameAndType(sc.ProductName, sc.ProductType, nil)
	if len(products) < 1 {
		return false
	}
	//只要有一个产品包可升级, 可升级判断条件：sql模块版本号大小
	for _, product := range products {
		if product.ProductType != 0 {
			continue
		}
		osc, err := schema.Unmarshal(product.Schema)
		if err != nil {
			log.Errorf("%v", err.Error())
			return false
		}
		log.Infof("current product %v version %v, target product %v version %v", sc.ProductName, sc.ProductVersion, osc.ProductName, osc.ProductVersion)
		if CompareVersion(osc.ProductVersion, sc.ProductVersion) > 0 {
			return true
		}
	}
	return false
}

func PatchUpload(ctx context.Context) apibase.Result {
	uploadLock.Lock()
	defer uploadLock.Unlock()

	tx := model.USE_MYSQL_DB().MustBegin()

	log.Debugf("[Patch Update] update product patches")

	// 获取入参以及判断入参是否为空
	paramErrs := apibase.NewApiParameterErrors()
	parentProductName := ctx.FormValue("parentProductName")
	productName := ctx.FormValue("product_name")
	productVersion := ctx.FormValue("version")
	path := ctx.FormValue("path")
	productType := ctx.FormValue("product_type")
	packageName := ctx.FormValue("package_name")

	if parentProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("parentProductName is empty"))
	}
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("version is empty"))
	}
	if path == "" {
		paramErrs.AppendError("$", fmt.Errorf("path is empty"))
	}
	if productType == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_type is empty"))
	}

	productTypeint, err := strconv.Atoi(productType)
	if err != nil {
		return fmt.Errorf("product_type string translation int failed：%s\n", err)
	}

	if productTypeint >= 1 {
		return fmt.Errorf("unsupported package type:%s\n", "only supported host package")
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	// 获取集群id
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	// 获取用户id
	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		userId = 0
		log.Infof("update no userId: %v", err.Error())
	}
	log.Infof("update userId:  %v", userId)

	// 拼接路径
	serviceName := strings.Split(path, "/")[0] //"dtuic/tools/arthas/arthas-spy.jar"-->"dtuic"
	update_dir := base.INSTALL_CURRRENT_PATH + path
	backup_dir := base.INSTALL_CURRRENT_PATH + productName + "/" + serviceName + "/" + "patch_backup"
	upload_dir := base.WebRoot + "/" + productName + "/" + productVersion + "/" + serviceName + "/" + "patches_package"

	// 获取补丁包
	f, _, formerr := ctx.FormFile("package")
	if f == nil && packageName != "" {
		// 如果用户没上传则判断本地补丁包是否存在
		_, err = os.Stat(upload_dir + "/" + packageName)
		if err != nil {
			if os.IsNotExist(err) {
				log.Errorf("%v", err)
				return fmt.Errorf("package is not exist: %s", packageName)
			}
			return err
		}
	} else {
		if formerr != nil {
			return formerr
		}
		defer f.Close()
		// 创建补丁包上传目录以及创建补丁包文件
		_, err = os.Stat(upload_dir)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(upload_dir, 0755)
				if err != nil {
					log.Errorf("%v", err)
					return err
				}
			}
		}
		out, err := os.OpenFile(upload_dir+"/"+packageName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		defer out.Close()
		io.Copy(out, f)
	}

	// 获取产品包对应的id
	var productInfo model.DeployProductListInfo
	getProductid := fmt.Sprintf("select * from %s where product_name=\"%s\" and product_version=\"%s\"", model.DeployProductList.TableName, productName, productVersion)
	if err := model.DeployProductList.GetDB().Get(&productInfo, getProductid); err != nil {
		return err
	}

	// 状态为更新中
	updateUUID := uuid.NewV4()
	productUpdateHistoryInfo := model.DeployProductUpdateHistoryInfo{
		Namespace:         "",
		ClusterId:         clusterId,
		UpdateUUID:        updateUUID,
		ProductName:       productName,
		ProductVersion:    productVersion,
		ParentProductName: parentProductName,
		UserId:            userId,
		PackageName:       packageName,
		UpdateDir:         update_dir,
		BackupDir:         backup_dir,
		Status:            "update",
		ProductId:         productInfo.ID,
	}
	//所有的 list 参数的 uuid 都要入库设置类型
	err = model.DeployUUID.InsertOne(updateUUID.String(), "", model.ManualDeployUUIDType, productInfo.ID)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	query := "INSERT INTO " + model.DeployProductUpdateHistory.TableName + " (namespace,cluster_id, product_name, update_uuid, product_version, `status`, parent_product_name, update_start_time, user_id, package_name,update_dir,backup_dir,product_id)" +
		"VALUES (:namespace,:cluster_id, :product_name, :update_uuid, :product_version, :status , :parent_product_name, NOW(), :user_id,:package_name,:update_dir,:backup_dir,:product_id)"
	if _, err := tx.NamedExec(query, &productUpdateHistoryInfo); err != nil {
		log.Errorf("%v", err)
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Errorf("%v", err)
		return err
	}
	defer func() {
		if err := addSafetyAuditRecord(ctx, "补丁包更新", "上传补丁包", "补丁包名称："+packageName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	return updateUUID
}

var wg sync.WaitGroup

func PatchUpdate(ctx context.Context) apibase.Result {
	var err error
	log.Debugf("[Patches Update] update product patches")

	// 获取入参以及判断入参是否为空
	paramErrs := apibase.NewApiParameterErrors()
	param := patchupdateParam{}
	if err = ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON ERROR: %v", err)
	}

	if param.ParentProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("parentProductName is empty"))
	}
	if param.ProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if param.Version == "" {
		paramErrs.AppendError("$", fmt.Errorf("version is empty"))
	}
	if param.Path == "" {
		paramErrs.AppendError("$", fmt.Errorf("path is empty"))
	}
	if param.PackageName == "" {
		paramErrs.AppendError("$", fmt.Errorf("package_name is empty"))
	}
	if param.UpdateUUID == "" {
		paramErrs.AppendError("$", fmt.Errorf("uuid is empty"))
	}

	updateUUID, err := uuid.FromString(param.UpdateUUID)
	if err != nil {
		return fmt.Errorf("UUID ERROR: %v", err)
	}

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	// 拼接路径
	serviceName := strings.Split(param.Path, "/")[0] //"dtuic/tools/arthas/arthas-spy.jar"-->"dtuic"
	upload_dir := base.WebRoot + "/" + param.ProductName + "/" + param.Version + "/" + serviceName + "/" + "patches_package"

	//防止补丁包被误删除，重新更新时出问题
	downloadfile := upload_dir + "/" + param.PackageName
	_, err = os.Stat(downloadfile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Errorf("%v", err)
			return fmt.Errorf("package is not exist: %s", param.PackageName)
		}
		return err
	}

	//获取服务组件所在节点信息,并开始更新
	serviceInfo, count := model.DeployInstanceList.GetInstanceBelongService(param.ProductName, serviceName, clusterId)
	counterr := 0
	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()

	err = model.OperationList.Insert(model.OperationInfo{
		ClusterId:       clusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.PatchUpdate.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Product.Code,
		ObjectValue:     param.ProductName,
	})
	if err != nil {
		return fmt.Errorf("OperationList Insert error: %v", err)
	}
	for _, svcinfo := range serviceInfo {
		wg.Add(1)
		go func() {

			errs := downloadPatches(updateUUID, svcinfo, param.Path, downloadfile, operationId)
			if errs != nil {
				counterr = counterr + 1
			}
			if counterr > 0 && counterr <= count {
				err = errs
			}
		}()
		wg.Wait()
	}

	defer func() {
		var status = "success"
		if err != nil {
			status = "fail"
		}
		query := "UPDATE " + model.DeployProductUpdateHistory.TableName + " SET `status`=?, update_end_time=NOW() WHERE update_uuid=? AND cluster_id=?"
		if _, err := model.DeployProductHistory.GetDB().Exec(query, status, updateUUID, clusterId); err != nil {
			log.Errorf("%v", err)
		}
	}()

	return nil
}

func downloadPatches(updateUUID uuid.UUID, serviceInfo model.InstanceAndProductInfo, path string, downloadfile string, operationId string) (err error) {

	updateQuery := "UPDATE " + model.DeployInstanceUpdateRecord.TableName + " SET `status`=?, status_message=?, progress=?, update_time=NOW() WHERE id=?"

	defer func() {
		// 实例更新失败改变deploy_instance_update_record表中记录
		if err != nil {
			if e, exist := err.(instanceErr); exist {
				if _, err := model.DeployInstanceUpdateRecord.GetDB().Exec(updateQuery, e.status, e.Error(), e.progress, e.id); err != nil {
					log.Errorf("%v", err)
				}
			}
		}
		wg.Done()
	}()

	// 初始化服务组件状态为更新中
	instanceupdateRecordInfo := model.DeployInstanceUpdateRecordInfo{
		UpdateUUID:         updateUUID,
		InstanceId:         serviceInfo.ID,
		Sid:                serviceInfo.Sid,
		Ip:                 serviceInfo.Ip,
		ProductName:        serviceInfo.ProductName,
		ProductNameDisplay: serviceInfo.ProductNameDisplay,
		ProductVersion:     serviceInfo.ProductVersion,
		Group:              serviceInfo.Group,
		ServiceName:        serviceInfo.ServiceName,
		ServiceNameDisplay: serviceInfo.ServiceNameDisplay,
		ServiceVersion:     serviceInfo.ServiceVersion,
		Status:             "update",
		Progress:           50,
	}
	insertQuery := "INSERT INTO " + model.DeployInstanceUpdateRecord.TableName +
		" (update_uuid, instance_id, sid, ip, product_name, product_name_display, product_version, `group`, service_name, service_name_display, service_version, `status`, progress) VALUES" +
		" (:update_uuid, :instance_id, :sid, :ip, :product_name, :product_name_display, :product_version, :group, :service_name, :service_name_display, :service_version, :status, :progress)"
	var rlt sql.Result
	if rlt, err = model.DeployInstanceUpdateRecord.GetDB().NamedExec(insertQuery, &instanceupdateRecordInfo); err != nil {
		log.Errorf("%v", err)
		return err
	}
	id, _ := rlt.LastInsertId()

	//创建实例对象进行补丁包更新操作
	var instancer instance.Instancer
	instancer = instance.NewCommonInstancer(serviceInfo.ClusterId, serviceInfo.Pid, serviceInfo.Ip, serviceInfo.Sid, serviceInfo.ServiceName, serviceInfo.AgentId, operationId)

	if err := instancer.PatchUpdate(serviceInfo.ProductName, serviceInfo.ServiceName, path, downloadfile); err != nil {
		log.Errorf("%v", err)
		return instanceErr{id, "fail", 50, err}
	}

	if _, err := model.DeployInstanceUpdateRecord.GetDB().Exec(updateQuery, "success", "", 100, id); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

// ProductInfo
// @Description  	GET Product Info
// @Summary      	获取产品信息
// @Tags         	product
// @Accept          application/json
// @Accept          application/json
// @Param           deploy_status query  string  false  "部署状态"
// @Param           parentProductName query  string  false  "父级产品包名称"
// @Param           productName query  string  false  "产品名称"
// @Param           clusterId query  string  false  "集群ID"
// @Param           product_type query  string  false  "产品类型"
// @Param           mode query  string  false  "模式"
// @Success         200  {object} string "{"msg":"ok","code":0,"data":{"list":"","count":""}}"
// @Router          /api/v2/product/{product_name}/version/{product_version} [get]
func ProductInfo(ctx context.Context) apibase.Result {
	var deployStatus []string
	var deployProductNames []string
	if status := ctx.URLParam("deploy_status"); status != "" {
		deployStatus = strings.Split(status, ",")
	}
	parentProductName := ctx.URLParam("parentProductName")
	if productNames := ctx.URLParam("productName"); productNames != "" {
		deployProductNames = strings.Split(productNames, ",")
	}
	//productVersionLike := sqlTransfer(ctx.URLParam("productVersion"))
	productVersionLike := ctx.URLParam("productVersion")
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	clusterId := ctx.URLParam("clusterId")
	productType := ctx.URLParam("product_type")
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	// 在获取已部署应用列表的时候置空cookie中的namespace
	mode := ctx.URLParam("mode")
	if mode != "" {
		namespace = ""
	}

	// 获取集群id
	var cid int
	var err error
	if clusterId == "" {
		cid, err = GetCurrentClusterId(ctx)
		if err != nil {
			log.Errorf("[ProductInfo] get cluster id from cookie error: %s", err)
			//return fmt.Errorf("[ProductInfo] get cluster id from cookie error: %s", err)
		}
	} else {
		cid, _ = strconv.Atoi(clusterId)
	}

	ProductListInfo := make([]model.DeployProductListInfoWithNamespace, 0)
	// 集群id为0时表示，获取上传的产品包信息，非0表示获取指定集群下已部署过的产品包的信息
	if cid == 0 {
		// 安装包管理下的所有上传的产品包信息，同时根据前端检索条件进行检索
		ProductListInfo, _ = model.DeployProductList.GetProductListInfo(parentProductName, productName, productVersionLike, productVersion, productType, cid, deployStatus, deployProductNames, namespace)
	} else {
		// 获取集群下所有部署的产品包信息，同时根据前端检索条件进行检索
		ProductListInfo, err = model.DeployClusterProductRel.GetDeployClusterProductList(parentProductName, productName, productVersionLike, productVersion, productType, cid, deployStatus, deployProductNames, namespace)
		if err != nil {
			return fmt.Errorf("ProductInfo query error %s", err)
		}
	}

	type SmoothUpgradeInfoRes struct {
		Id                 int    `json:"id"`
		ProductName        string `json:"product_name"`
		ProductNameDisplay string `json:"product_name_display"`
		ProductVersion     string `json:"product_version"`
		IsCurrentVersion   int    `json:"is_current_version"`
		Status             string `json:"status"`
		NameSpace          string `json:"name_space"`
		DeployUUID         string `json:"deploy_uuid"`
		ProductType        int    `json:"product_type"`
		CanRollback        bool   `json:"can_rollback"`
		DeployTime         string `json:"deploy_time"`
		CreateTime         string `json:"create_time"`
	}
	type ProductInfoRes struct {
		Id                   int                   `json:"id"`
		ProductName          string                `json:"product_name"`
		ProductNameDisplay   string                `json:"product_name_display"`
		ProductVersion       string                `json:"product_version"`
		UserName             string                `json:"user_name"`
		IsCurrentVersion     int                   `json:"is_current_version"`
		Status               string                `json:"status"`
		NameSpace            string                `json:"name_space"`
		DeployUUID           string                `json:"deploy_uuid"`
		ProductType          int                   `json:"product_type"`
		CanUpgrade           bool                  `json:"can_upgrade"`
		CanRollback          bool                  `json:"can_rollback"`
		CanSmoothUpgrade     bool                  `json:"can_smooth_upgrade"`
		SmoothUpgradeProduct *SmoothUpgradeInfoRes `json:"smooth_upgrade_product,omitempty"`
		UpgradeService       []string              `json:"upgrade_service,omitempty"`
		DeployTime           string                `json:"deploy_time"`
		CreateTime           string                `json:"create_time"`
	}

	resultList := make([]ProductInfoRes, 0)
	smoothUpgradeMap := make(map[string][]string)
	for index, s := range ProductListInfo {
		//获取平滑升级信息
		smoothUpgradeList, err := upgrade.SmoothUpgrade.GetByProductName(s.ProductName)
		if err != nil {
			log.Errorf("SmoothUpgrade-query db error: %v", err)
			return err
		}
		for _, su := range smoothUpgradeList {
			smoothUpgradeMap[su.ProductName] = append(smoothUpgradeMap[su.ProductName], su.ServiceName)
		}

		var vFound, sFound bool
		sc, err := schema.Unmarshal(s.Product)
		if err != nil {
			log.Errorf("[ProductInfo] Unmarshal err: %v", err)
		}
		m := ProductInfoRes{}
		m.Id = s.ID
		m.ProductName = s.ProductName
		m.ProductNameDisplay = s.ProductNameDisplay
		m.ProductVersion = s.ProductVersion
		if strings.Contains(s.ProductVersion, productVersionLike) {
			vFound = true
		}
		//m["product"] = sc
		if ProductListInfo[index].UserId > 0 {
			if err, userInfo := model.UserList.GetInfoByUserId(ProductListInfo[index].UserId); err != nil {
				m.UserName = ""
			} else {
				m.UserName = userInfo.UserName
			}
		} else {
			m.UserName = ""
		}
		// 根据已安装的产品包id，获取该产品包对应的详细信息
		installed, err := model.DeployClusterProductRel.GetProductByPid(s.ID)
		if err == nil && len(installed) > 0 {
			m.IsCurrentVersion = 1
		} else {
			suInstalled, err := model.DeployClusterSmoothUpgradeProductRel.GetProductByPid(s.ID)
			if err == nil && len(suInstalled) > 0 {
				m.IsCurrentVersion = 1
			}
		}
		m.Status = s.Status
		if len(deployStatus) == 0 {
			sFound = true
		}
		for _, status := range deployStatus {
			if status == s.Status {
				sFound = true
				break
			}
		}
		m.NameSpace = s.Namespace
		m.DeployUUID = s.DeployUUID
		m.ProductType = s.ProductType
		//kubernetes type not surpport upgrade
		if cid > 0 {
			m.CanUpgrade = checkCanUpgrade(sc)
			m.CanRollback = CanRollback(cid, s.ProductName, s.ProductVersion)
			if _, ok := smoothUpgradeMap[s.ProductName]; ok {
				m.CanSmoothUpgrade = true
				info, err := model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterIdNamespace(s.ProductName, cid, namespace)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("smooth Upgrade ProductInfo query error %s", err)
				}
				if err == nil {
					var smoothUpgradeInfo SmoothUpgradeInfoRes
					smoothUpgradeInfo.IsCurrentVersion = 1
					smoothUpgradeInfo.Id = info.ID
					smoothUpgradeInfo.ProductName = info.ProductName
					smoothUpgradeInfo.ProductNameDisplay = info.ProductNameDisplay
					smoothUpgradeInfo.ProductVersion = info.ProductVersion
					if strings.Contains(info.ProductVersion, productVersionLike) {
						vFound = true
					}
					smoothUpgradeInfo.Status = info.Status
					for _, status := range deployStatus {
						if status == info.Status {
							sFound = true
							break
						}
					}
					smoothUpgradeInfo.NameSpace = info.Namespace
					smoothUpgradeInfo.DeployUUID = info.DeployUUID
					smoothUpgradeInfo.ProductType = info.ProductType
					if info.Status == model.PRODUCT_STATUS_DEPLOYED || info.Status == model.PRODUCT_STATUS_DEPLOY_FAIL {
						smoothUpgradeInfo.CanRollback = true
					}
					if info.Status == model.PRODUCT_STATUS_DEPLOYING || info.Status == model.PRODUCT_STATUS_DEPLOY_FAIL {
						m.CanUpgrade = false
					}
					if m.Status == model.PRODUCT_STATUS_DEPLOYING {
						smoothUpgradeInfo.CanRollback = false
					}
					if info.DeployTime.Valid == true {
						smoothUpgradeInfo.DeployTime = info.DeployTime.Time.Format(base.TsLayout)
					} else {
						smoothUpgradeInfo.DeployTime = ""
					}
					if info.CreateTime.Valid == true {
						smoothUpgradeInfo.CreateTime = info.CreateTime.Time.Format(base.TsLayout)
					} else {
						smoothUpgradeInfo.CreateTime = ""
					}
					m.SmoothUpgradeProduct = &smoothUpgradeInfo
				}
				m.UpgradeService = smoothUpgradeMap[s.ProductName]
			} else {
				m.CanSmoothUpgrade = false
			}
		}
		if s.DeployTime.Valid == true {
			m.DeployTime = s.DeployTime.Time.Format(base.TsLayout)
		} else {
			m.DeployTime = ""
		}

		if s.CreateTime.Valid == true {
			m.CreateTime = s.CreateTime.Time.Format(base.TsLayout)
		} else {
			m.CreateTime = ""
		}
		// 当产品包的部署状态变为undeploying或deploying状态时重写can_upgrade属性
		if s.Status == model.PRODUCT_STATUS_UNDEPLOYING || s.Status == model.PRODUCT_STATUS_DEPLOYING {
			m.CanUpgrade = false
			m.CanRollback = false
		}
		if vFound && sFound {
			resultList = append(resultList, m)
		}
	}

	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, nil)
	switch pagination.SortBy {
	case "product_version":
		sort.SliceStable(resultList, func(i, j int) bool {
			if pagination.SortDesc {
				return resultList[i].ProductVersion > resultList[j].ProductVersion
			} else {
				return resultList[i].ProductVersion < resultList[j].ProductVersion
			}
		})
	case "deploy_time":
		sort.SliceStable(resultList, func(i, j int) bool {
			if pagination.SortDesc {
				return resultList[i].DeployTime > resultList[j].DeployTime
			} else {
				return resultList[i].DeployTime < resultList[j].DeployTime
			}
		})
	case "create_time":
		sort.SliceStable(resultList, func(i, j int) bool {
			if pagination.SortDesc {
				return resultList[i].CreateTime > resultList[j].CreateTime
			} else {
				return resultList[i].CreateTime < resultList[j].CreateTime
			}
		})
	default:
		// 默认以服务名排序
		sort.SliceStable(resultList, func(i, j int) bool {
			return strings.Compare(resultList[i].ProductName, resultList[j].ProductName) == -1
		})
	}
	// 重写分页
	total := len(resultList) // result总数量
	if pagination.Start > 0 {
		if pagination.Start+pagination.Limit < total {
			resultList = resultList[pagination.Start : pagination.Start+pagination.Limit]
		} else if pagination.Start > total {
			resultList = nil
		} else {
			resultList = resultList[pagination.Start:total]
		}
	} else {
		if pagination.Limit == 0 {
			resultList = resultList[:total]
		} else if pagination.Limit < total {
			resultList = resultList[:pagination.Limit]
		}
	}

	return map[string]interface{}{
		"list":  resultList,
		"count": total,
	}
}

func ParentProductInfo(ctx context.Context) apibase.Result {
	parentProductNames := model.DeployProductList.GetDeployParentProductList()
	return parentProductNames
}

// ProductList
// @Description  	GET Product List
// @Summary      	获取产品包列表
// @Tags         	product
// @Accept          application/json
// @Accept          application/json
// @Param           product_name query  string  false  "产品名称"
// @Param           product_type query  string  false  "产品类型"
// @Param           deploy_status query  string  false  "部署状态"
// @Success         200  {object} string "{"msg":"ok","code":0,"data":{"list":"","count":""}}"
// @Router          /api/v2/product/productList [get]
func ProductList(ctx context.Context) apibase.Result {
	log.Debugf("[ProductList] ProductList from EasyMatrix API ")

	var deployStatus []string
	productName := ctx.URLParam("product_name")
	productType := ctx.URLParam("product_type")
	if status := ctx.URLParam("deploy_status"); status != "" {
		deployStatus = strings.Split(status, ",")
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.DeployProductListInfo{})
	list, _ := model.DeployProductList.GetProductListByNameAndType(productName, productType, pagination)
	result := make([]model.DeployProductListInfo, 0)
	for _, pro := range list {
		installed, err := model.DeployClusterProductRel.GetByPidAndClusterId(pro.ID, clusterId)
		if err == nil {
			pro.Status = installed.Status
		} else {
			pro.Status = model.PRODUCT_STATUS_UNDEPLOYED
		}
		if len(deployStatus) == 0 {
			result = append(result, pro)
			continue
		}
		for _, status := range deployStatus {
			if pro.Status == status {
				result = append(result, pro)
			}
		}
	}
	return map[string]interface{}{
		"list":  result,
		"count": len(result),
	}
}

func ProductName(ctx context.Context) apibase.Result {
	log.Debugf("[ProductName] ProductName from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	productType, err := ctx.URLParamInt64("product_type")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("type is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	productNames := model.DeployProductList.GetDeployProductName(int(productType))
	return productNames
}

// ProductStart
// @Description  	Start Product
// @Summary      	启动组件
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Success         200 {string} string   "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/{pid}/start [get]
func ProductStart(ctx context.Context) apibase.Result {
	log.Debugf("[ProductStart] ProductStart from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	pInfo, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		return err
	}

	err = model.OperationList.Insert(model.OperationInfo{
		ClusterId:       clusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.ProductStart.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Product.Code,
		ObjectValue:     pInfo.ProductName,
	})
	if err != nil {
		log.Errorf("OperationList Insert err:%v", err)
	}
	grouper, err := group.NewGrouper(pid, clusterId, "", operationId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	grouper.Start()
	//err = grouper.Start()
	//if err != nil {
	//	//启动之后，开启告警
	//	err = grafana.StartAlert(pid)
	//	if err != nil {
	//		log.Errorf("%v", err)
	//		return err
	//	}
	//}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		productInfo, err := model.DeployProductList.GetProductInfoById(pid)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "集群运维", "组件启动", "集群名称："+clusterInfo.Name+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
		if err := model.NotifyEvent.DeleteNotifyEvent(clusterId, 0, productInfo.ProductName, "", "", true); err != nil {
			log.Errorf("delete notify event error: %v", err)
		}
	}()
	return grouper.GetResult()
}

// ProductStop
// @Description  	Stop Product
// @Summary      	停止组件
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           type  query     string   false  "停止类型"
// @Success         200 {string} string   "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/{pid}/stop [get]
func ProductStop(ctx context.Context) apibase.Result {
	log.Debugf("[ProductStop] ProductStop from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	stopAgentOptionsType, err := ctx.URLParamInt("type")
	if err != nil {
		stopAgentOptionsType = agent.AGENT_STOP_RECOVER
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	grouper, err := group.NewGrouper(pid, clusterId, "", "")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	//停止之前，停止告警
	//err = grafana.StopAlert(pid)
	//if err != nil {
	//	log.Errorf("%v", err)
	//	return err
	//}

	grouper.Stop(stopAgentOptionsType)
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		productInfo, err := model.DeployProductList.GetProductInfoById(pid)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "集群运维", "组件停止", "集群名称："+clusterInfo.Name+", 组件名称："+productInfo.ProductName+productInfo.ProductVersion); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
		if err := model.NotifyEvent.UpdateProductStopped(clusterId, 0, productInfo.ProductName); err != nil {
			log.Errorf("failed to update product stopped, error: %v", err)
		}
	}()
	return grouper.GetResult()
}

// ProductUncheckedServices
// @Description  	Get Unchecked Services
// @Summary      	查看未部署组件
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           namespace  query     string   false  "命名空间"
// @Success         200 {string} string   "{"msg":"ok","code":0,"data":[""]}"
// @Router          /api/v2/product/{pid}/unchecked_services [get]
func ProductUncheckedServices(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ProductUncheckedServices] get unchecked services from EasyMatrix API ")

	pid, err := ctx.Params().GetInt("pid")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	product, err := model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	namespace := ctx.URLParam("namespace")

	//根据产品名、产品类型获取uncheckServices
	//同一个集群，同一个产品的uncheckServices 一致
	products, _ := model.DeployProductList.GetProductListByNameAndType(product.ProductName, strconv.Itoa(product.ProductType), nil)

	if len(products) == 0 {
		log.Errorf("not found product %v", product.ProductName)
		return err
	}
	var uncheckServices []string
	var uncheckedServiceList = make([]*model.DeployUncheckedServiceInfo, 0)
	for _, p := range products {
		info, err := model.DeployUncheckedService.GetUncheckedServicesByPidClusterId(p.ID, clusterId, namespace)
		if err != nil {
			return fmt.Errorf("UncheckedServices query deployUncheckedService error: %s", err)
		}
		uncheckedServiceList = append(uncheckedServiceList, info)
	}

	sort.SliceStable(uncheckedServiceList, func(i, j int) bool {
		return uncheckedServiceList[i].UpdateDate.Time.After(uncheckedServiceList[j].UpdateDate.Time)
	})
	uncheckServices = strings.Split(uncheckedServiceList[0].UncheckedServices, ",")
	originSchema, err := schema.Unmarshal(product.Schema)
	if err != nil {
		return err
	}
	//获取当前部署产品包内的所有服务，如果unchecksevice 不在当前部署包的服务里面，意味着之前的某个服务被改名或者移除了
	for idx, service := range uncheckServices {
		if _, ok := originSchema.Service[service]; !ok {
			uncheckServices = append(uncheckServices[:idx], uncheckServices[idx+1:]...)
		}
	}
	return uncheckServices
}

func ProductDelete(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ProductDelete] delete product from EasyMatrix API ")

	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")

	paramErrs := apibase.NewApiParameterErrors()
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		return err
	}
	products, _ := model.DeployClusterProductRel.GetProductByPid(info.ID)
	if len(products) > 0 {
		return fmt.Errorf("product's status must be " + model.PRODUCT_STATUS_UNDEPLOYED)
	}
	query := "DELETE FROM " + model.DeployProductList.TableName + " WHERE product_name=? AND product_version=? AND `status`=?"
	if _, err := model.DeployProductList.GetDB().Exec(query, productName, productVersion, info.Status); err != nil {
		log.Errorf("%v", err)
		return err
	}

	return os.RemoveAll(filepath.Join(base.WebRoot, productName, productVersion))
}

func UpdateHistory(ctx context.Context) apibase.Result {
	log.Debugf("[Product->History] return update history product info from EasyMatrix API ")
	var updateStatus []string
	var productType []int
	paramErrs := apibase.NewApiParameterErrors()

	productName := ctx.Params().Get("product_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	parentProductName := ctx.URLParam("parentProductName")
	clusterId := ctx.URLParam("clusterId")
	productNames := ctx.URLParam("productName")
	if status := ctx.URLParam("update_status"); status != "" {
		updateStatus = strings.Split(status, ",")
	}
	if pType := ctx.URLParam("product_type"); pType != "" {
		for _, pt := range strings.Split(pType, ",") {
			ipt, _ := strconv.Atoi(pt)
			productType = append(productType, ipt)
		}
	}
	productVersionLike := sqlTransfer(ctx.URLParam("productVersion"))

	paramErrs.CheckAndThrowApiParameterErrors()

	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.DeployProductUpdateHistoryInfo{})

	info, count := model.DeployProductUpdateHistory.GetDeployProductUpdateHistory(pagination, parentProductName, productNames, productType, updateStatus, productVersionLike, clusterId)

	list := []map[string]interface{}{}

	for _, s := range info {
		h := map[string]interface{}{}
		h["id"] = s.ID
		h["product_name"] = s.ProductName
		h["product_name_display"] = s.ProductNameDisplay
		h["product_version"] = s.ProductVersion
		h["status"] = s.Status
		h["update_uuid"] = s.UpdateUUID
		h["namespace"] = s.Namespace
		h["product_type"] = s.ProductType
		h["package_name"] = s.PackageName
		h["update_dir"] = s.UpdateDir
		h["backup_dir"] = s.BackupDir
		h["product_id"] = s.ProductId

		if s.UserId > 0 {
			if err, userInfo := model.UserList.GetInfoByUserId(s.UserId); err != nil {
				h["username"] = ""
			} else {
				h["username"] = userInfo.UserName
			}
		} else {
			h["username"] = ""
		}

		if s.CreateTime.Valid == true {
			h["create_time"] = s.CreateTime.Time.Format(base.TsLayout)
		} else {
			h["create_time"] = ""
		}

		list = append(list, h)
	}

	return map[string]interface{}{
		"list":  list,
		"count": count,
	}
}

func History(ctx context.Context) apibase.Result {
	log.Debugf("[Product->History] return history product info from EasyMatrix API ")
	var deployStatus []string
	var productType []int
	paramErrs := apibase.NewApiParameterErrors()

	productName := ctx.Params().Get("product_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	parentProductName := ctx.URLParam("parentProductName")
	clusterId := ctx.URLParam("clusterId")
	productNames := ctx.URLParam("productName")
	if status := ctx.URLParam("deploy_status"); status != "" {
		deployStatus = strings.Split(status, ",")
	}
	if pType := ctx.URLParam("product_type"); pType != "" {
		for _, pt := range strings.Split(pType, ",") {
			ipt, _ := strconv.Atoi(pt)
			productType = append(productType, ipt)
		}
	}
	productVersionLike := sqlTransfer(ctx.URLParam("productVersion"))

	paramErrs.CheckAndThrowApiParameterErrors()

	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.DeployProductHistoryInfo{})

	info, count := model.DeployProductHistory.GetDeployProductHistory(pagination, parentProductName, productNames, productType, deployStatus, productVersionLike, clusterId)

	list := []map[string]interface{}{}

	for _, s := range info {
		h := map[string]interface{}{}
		h["id"] = s.ID
		h["product_name"] = s.ProductName
		h["product_name_display"] = s.ProductNameDisplay
		h["product_version"] = s.ProductVersion
		h["status"] = s.Status
		h["deploy_uuid"] = s.DeployUUID
		h["namespace"] = s.Namespace
		h["product_type"] = s.ProductType

		if s.UserId > 0 {
			if err, userInfo := model.UserList.GetInfoByUserId(s.UserId); err != nil {
				h["username"] = ""
			} else {
				h["username"] = userInfo.UserName
			}
		} else {
			h["username"] = ""
		}

		if s.CreateTime.Valid == true {
			h["create_time"] = s.CreateTime.Time.Format(base.TsLayout)
		} else {
			h["create_time"] = ""
		}

		list = append(list, h)
	}

	return map[string]interface{}{
		"list":  list,
		"count": count,
	}
}

func HaRole(ctx context.Context) apibase.Result {
	log.Debugf("[Product->HaRole] return ha role info from EasyMatrix API ")
	paramErrs := apibase.NewApiParameterErrors()

	productName := ctx.Params().Get("product_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}

	serviceName := ctx.Params().Get("service_name")
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	info, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	ret := map[string]interface{}{}
	if err == sql.ErrNoRows {
		return ret
	}

	roleData := harole.RoleData(info.ID, serviceName)
	if roleData != nil {
		return roleData
	}

	return ret
}

func Service(ctx context.Context) apibase.Result {
	log.Debugf("Service: %v", ctx.Request().RequestURI)

	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	relyNamespace := ctx.URLParam("relynamespace")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("Get cluster id error:%v", err)
		return fmt.Errorf("Get cluster id error:%v", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	type serviceRes struct {
		ServiceName        string `json:"serviceName"`
		ServiceNameDisplay string `json:"serviceNameDisplay"`
		ServiceVersion     string `json:"serviceVersion"`
		BaseProduct        string `json:"baseProduct"`
		baseService        string `json:"baseService"`
	}

	clusterinfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	if err != nil {
		log.Errorf("[Product->Service] Get cluster info error:%v", err)
		return fmt.Errorf("Get cluster info error:%v", err)
	}

	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		return err
	}

	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("[Product->Service] Unmarshal err: %v", err)
		return err
	}

	// 获取该产品包下服务组件依赖对应服务的相关配置信息,主要用于获取依赖组件的版本号
	if clusterinfo.Type == model.DEPLOY_CLUSTER_TYPE_KUBERNETES {
		if err = inheritK8sBaseService(clusterId, relyNamespace, sc, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[Product->Service] inheritK8sBaseService error:%v", err)
		}
	} else {
		if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[Product->Service] inheritBaseService err: %+v", err)
		}
	}

	services := make([]serviceRes, 0)
	for name, svc := range sc.Service {
		serviceDisplay := svc.ServiceDisplay
		if serviceDisplay == "" {
			serviceDisplay = name
		}
		services = append(services, serviceRes{
			ServiceName:        name,
			ServiceNameDisplay: serviceDisplay,
			ServiceVersion:     svc.Version,
			BaseProduct:        svc.BaseProduct,
			baseService:        svc.BaseService,
		})
	}
	// 默认以服务名排序
	sort.SliceStable(services, func(i, j int) bool {
		return strings.Compare(services[i].ServiceName, services[j].ServiceName) == -1
	})

	return services
}

// get the services that baseservice is parsed.
// the baseservice will be set to default group.
func ServiceGroup(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	uncheckedServices := strings.Split(ctx.URLParam("unchecked_services"), ",")
	// k8s模式下产品包依赖的namespace
	relyNamespace := ctx.URLParam("relynamespace")
	namespace := ctx.URLParam("namespace")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	// 从url参数中获取集群的id
	clusterId, err := ctx.URLParamInt("clusterId")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	upgradeMode := ctx.URLParam("upgrade_mode")
	paramErrs.CheckAndThrowApiParameterErrors()
	// 获取正在部署的产品包信息
	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)

	if err != nil {
		return err
	}

	// 获取产品包的schema信息
	sc, err := schema.Unmarshal(info.Product)

	if err != nil {
		log.Errorf("[Product->ServiceGroup] Unmarshal err: %v", err)
		return err
	}

	err, userInfo := model.UserList.GetInfoByUserId(1)
	if err != nil {
		log.Errorf("GetInfoByUserId %v", err)
		return err
	}
	reg := regexp.MustCompile(`(?i).*password.*`)

	// just return the services that are grouped.
	// it is used to display on the front page to show the structure.
	if uncheckedServices[0] == "undefined" {
		res := sc.Group(nil)
		for _, group := range res {
			for _, svc := range group {
				for key, configItem := range svc.Config {

					if reg.Match([]byte(key)) {

						defaultValue, err := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(configItem.(schema.VisualConfig).Default.(*string))), userInfo.PassWord)

						if err != nil {
							return err
						}
						value, err := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(configItem.(schema.VisualConfig).Value.(*string))), userInfo.PassWord)
						if err != nil {
							return err
						}
						svc.Config[key] = schema.VisualConfig{
							Default: defaultValue,
							Desc:    configItem.(schema.VisualConfig).Desc,
							Type:    configItem.(schema.VisualConfig).Type,
							Value:   value,
						}
					}
				}
			}
		}
		return res
	}
	cluster, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	if err != nil {
		log.Errorf("[Product->ServiceGroup] get clusterinfo err : %v", err)
		return fmt.Errorf("[ServiceGroup] get current cluster err:%v", err)
	}
	// do not judge if is kubenetes type by product info.
	// now the product name and version identify a product, but as the kubenets type import,
	//it is not a good way to identify a product only by product name and version
	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_KUBERNETES {
		if err = inheritK8sBaseService(clusterId, relyNamespace, sc, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[Product->ServiceGroup] inheritK8sBaseService err: %v", err)
			return fmt.Errorf("[ServiceGroup] inheritK8sBaseService err: %v", err)
		}
		if err = setSchemaFieldDNS(clusterId, sc, model.USE_MYSQL_DB(), namespace); err != nil {
			log.Errorf("[Product->ServiceGroup] setSchemaFieldDNS err %v", err)
			return fmt.Errorf("[ServiceGroup] setSchemaFieldDNS err %v", err)
		}
		//beacause of the dtbase and all the product's depends is bridge.
		//if the baseservice is deployed in the same cluster.
		//and if the service's baseservice ip is modified.
		//then, all the product in the cluster will deploy can see the change of the baseservice

		//if the baseserivce is deployed in the other cluster.
		//then, all the product will show the ip of the baseservice pod's hostip
		//all of the service need to deploy will need to modify
		//if err = setBaseServiceAddr(baseClusterId, sc, model.USE_MYSQL_DB()); err != nil {
		//	log.Errorf("[Product->ServiceGroup] setBaseServiceAddr err: %v", err)
		//	return err
		//}
		if err = setSchemafieldModifyInfo(clusterId, sc, namespace); err != nil {
			log.Errorf("[product->ServiceGroup] service config modify fail,err: %v", err)
			return fmt.Errorf("[ServiceGroup] service config modify fail,err: %v", err)
		}
		if err = BaseServiceAddrModify(clusterId, sc, namespace); err != nil {
			log.Errorf("[product->ServiceGroup] set base service addr with modified fail,err: %v", err)
			return fmt.Errorf("[erviceGroup] set base service addr with modified fail,err: %v", err)
		}
	} else {
		if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[Product->ServiceGroup] inheritBaseService warn: %+v", err)
		}
		if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
			log.Errorf("[Product->ServiceGroup] setSchemaFieldServiceAddr err: %v", err)
			return err
		}
	}
	if err = handleUncheckedServicesCore(sc, uncheckedServices); err != nil {
		log.Errorf("[Product->ServiceGroup] handleUncheckedServicesCore warn: %+v", err)
	}
	if err = sc.ParseVariable(); err != nil {
		log.Errorf("[Product->ServiceGroup] ParseVariable err: %v", err)
		return err
	}
	//添加select unselect  信息
	if upgradeMode == upgrade.SMOOTH_UPGRADE_MODE {
		if err = SmoothUpgradeWithIpRoleInfo(clusterId, info, sc); err != nil {
			log.Debugf("[Product->ServiceGroup] SmoothUpgradeWithIpRoleInfo err: %v", err)
			return err
		}
	} else {
		if err = WithIpRoleInfo(clusterId, sc); err != nil {
			log.Debugf("[Product->ServiceGroup] WithIpRoleInfo err: %v", err)
			return err
		}
	}
	res := sc.Group(uncheckedServices)
	for _, group := range res {
		for svcName, svc := range group {
			for key, configItem := range svc.Config {
				var multiFieldList []model.SchemaMultiFieldInfo
				query := "select field, hosts from " + model.SchemaMultiField.TableName + " where cluster_id=? and product_name=? and service_name=? and field_path=? and is_deleted=0 order by id asc"
				if err := model.USE_MYSQL_DB().Select(&multiFieldList, query, clusterId, productName, svcName, "Config."+key+".Value"); err != nil {
					log.Errorf("%v", err)
					continue
				}
				if len(multiFieldList) > 0 {
					fieldHostMap := map[string][]string{}
					keyList := []string{}
					for _, multiField := range multiFieldList {
						var hostSlice []string
						if !contains(keyList, multiField.Field) {
							keyList = append(keyList, multiField.Field)
						}
						hostSlice, _ = fieldHostMap[multiField.Field]
						hostSlice = append(hostSlice, multiField.Hosts)
						fieldHostMap[multiField.Field] = hostSlice
					}
					var value []map[string]string
					for _, k := range keyList {
						v := fieldHostMap[k]
						value = append(value, map[string]string{
							"hosts": strings.Join(v, ","),
							"field": k,
						})
					}
					svc.Config[key] = schema.VisualConfig{
						Default: configItem.(schema.VisualConfig).Default,
						Desc:    configItem.(schema.VisualConfig).Desc,
						Type:    configItem.(schema.VisualConfig).Type,
						Value:   value,
					}
				} else {
					if reg.Match([]byte(key)) {
						log.Infof("Match uncheckedServices password key: %s", key)

						defaultValue, err := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(configItem.(schema.VisualConfig).Default.(*string))), userInfo.PassWord)
						if err != nil {
							return err
						}
						value, err := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(configItem.(schema.VisualConfig).Value.(*string))), userInfo.PassWord)
						if err != nil {
							return err
						}
						svc.Config[key] = schema.VisualConfig{
							Default: defaultValue,
							Desc:    configItem.(schema.VisualConfig).Desc,
							Type:    configItem.(schema.VisualConfig).Type,
							Value:   value,
						}
					}
				}

			}
		}
	}
	return res
}

func CompareIpList(a, b []string) bool {
	sort.Sort(sort.StringSlice(a))
	sort.Sort(sort.StringSlice(b))
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func CheckMysqlAddr(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		log.Errorf("CheckMysqlAddr-product_name is empty")
		return fmt.Errorf("product_name为空")
	}
	param := struct {
		ClusterId    int    `json:"cluster_id"`
		FinalUpgrade bool   `json:"final_upgrade"`
		Ip           string `json:"ip"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("CheckMysqlAddr-parse param error: %v", err)
		return err
	}

	mysqlIPList, err := model.DeployMysqlIpList.GetMysqlIpList(param.ClusterId, productName)
	if err != nil {
		log.Errorf("get mysql ip list error:%v", err)
		return err
	}
	//平滑升级中提示修改数据库地址
	if !param.FinalUpgrade && CompareIpList(mysqlIPList, strings.Split(param.Ip, ",")) {
		return fmt.Sprintf("请更换 IP 地址")
	}
	//最终一次平滑升级提示修改数据库地址
	if param.FinalUpgrade && !CompareIpList(mysqlIPList, strings.Split(param.Ip, ",")) {
		return fmt.Sprintf("请更换 IP 地址为: %v", strings.Join(mysqlIPList, ","))
	}

	return ""
}

func loadKeyWithFile(filePath string) []string {
	placeholderKey := []string{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("[ServiceConfig] Config File err: %v", err)
		return nil
	}
	reg1 := regexp.MustCompile(`\{{\.\w+}}`)
	reg2 := regexp.MustCompile(`\w+`)
	for _, v := range reg1.FindAllString(string(file), -1) {
		k := reg2.FindString(v)
		if k != "" {
			placeholderKey = append(placeholderKey, k)
		}
	}
	return placeholderKey
}

func ServiceGroupFile(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	// k8s模式下产品包依赖的namespace
	relyNamespace := ctx.URLParam("relynamespace")
	namespace := ctx.URLParam("namespace")
	serviceName := ctx.URLParam("servicename")
	fileName := ctx.URLParam("file")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("servicename is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	// 从url参数中获取集群的id
	clusterId, err := ctx.URLParamInt("clusterId")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("clusterId err %v", err))
	}
	// 获取正在部署的产品包信息
	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		return err
	}
	// 获取产品包的schema信息
	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("[Product->ServiceGroup] Unmarshal err: %v", err)
		return err
	}
	if _, ok := sc.Service[serviceName]; !ok {
		paramErrs.AppendError("$", fmt.Errorf("serviceName %s not exist", serviceName))
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	cluster, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	if err != nil {
		log.Errorf("[Product->ServiceGroup] get clusterinfo err : %v", err)
		return fmt.Errorf("[ServiceGroup] get current cluster err:%v", err)
	}

	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_KUBERNETES {
		// k8s的未测试过
		if err = inheritK8sBaseService(clusterId, relyNamespace, sc, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[Product->ServiceGroup] inheritK8sBaseService err: %v", err)
			return fmt.Errorf("[ServiceGroup] inheritK8sBaseService err: %v", err)
		}
		if err = setSchemaFieldDNS(clusterId, sc, model.USE_MYSQL_DB(), namespace); err != nil {
			log.Errorf("[Product->ServiceGroup] setSchemaFieldDNS err %v", err)
			return fmt.Errorf("[ServiceGroup] setSchemaFieldDNS err %v", err)
		}
		if err = setSchemafieldModifyInfo(clusterId, sc, namespace); err != nil {
			log.Errorf("[product->ServiceGroup] service config modify fail,err: %v", err)
			return fmt.Errorf("[ServiceGroup] service config modify fail,err: %v", err)
		}
		if err = BaseServiceAddrModify(clusterId, sc, namespace); err != nil {
			log.Errorf("[product->ServiceGroup] set base service addr with modified fail,err: %v", err)
			return fmt.Errorf("[erviceGroup] set base service addr with modified fail,err: %v", err)
		}
	} else {
		if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[Product->ServiceGroup] inheritBaseService warn: %+v", err)
		}
		if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
			log.Errorf("[Product->ServiceGroup] setSchemaFieldServiceAddr err: %v", err)
			return err
		}
	}
	var (
		svc            = sc.Service[serviceName]
		_, userInfo    = model.UserList.GetInfoByUserId(1)
		placeholderKey []string
	)

	if err = sc.ParseVariable(); err != nil {
		log.Errorf("[Product->ServiceGroup] ParseVariable err: %v", err)
		return err
	}

	//添加select unselect  信息
	if err = WithIpRoleInfo(clusterId, sc); err != nil {
		log.Debugf("[Product->ServiceGroup] WithIpRoleInfo err: %v", err)
		return err
	}

	if fileName == "" {
		svc.Config.AesEncryptByPassword(userInfo.PassWord)
		return svc.Config
	}

	targetFile := filepath.Join(base.WebRoot, productName, productVersion, serviceName, fileName)
	placeholderKey = append(placeholderKey, loadKeyWithFile(targetFile)...)

	res := map[string]schema.VisualConfig{}
	for _, key := range placeholderKey {
		var (
			multiFieldList = []model.SchemaMultiFieldInfo{}
			fieldHostMap   = map[string][]string{}
			keyList        = []string{}
		)
		vsConfig := svc.Config[key].(schema.VisualConfig)
		query := "select field, hosts from " + model.SchemaMultiField.TableName + " where cluster_id=? and product_name=? and service_name=? and field_path=? and is_deleted=0 order by id asc"
		if err := model.USE_MYSQL_DB().Select(&multiFieldList, query, clusterId, productName, serviceName, "Config."+key+".Value"); err != nil {
			log.Errorf("%v", err)
			continue
		}

		if len(multiFieldList) == 0 && regexp.MustCompile(`(?i).*password.*`).Match([]byte(key)) {
			vsConfig.AesEncryptByPassword(userInfo.PassWord)
		}

		if len(multiFieldList) > 0 {
			for _, multiField := range multiFieldList {
				var hostSlice []string
				if !contains(keyList, multiField.Field) {
					keyList = append(keyList, multiField.Field)
				}
				hostSlice, _ = fieldHostMap[multiField.Field]
				hostSlice = append(hostSlice, multiField.Hosts)
				fieldHostMap[multiField.Field] = hostSlice
			}
			value := []map[string]string{}
			for _, k := range keyList {
				v := fieldHostMap[k]
				value = append(value, map[string]string{
					"hosts": strings.Join(v, ","),
					"field": k,
				})
			}
			vsConfig.Value = value
		}

		res[key] = vsConfig
	}
	return res
}

func ServiceTree(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	pkgPath := filepath.Join(base.WebRoot, productName, productVersion, serviceName)

	tree, err := util.ListFiles(pkgPath)
	if err != nil {
		log.Errorf("[Product->ServiceTree] read service file err:%v", err)
		return err
	}
	var subTree []string
	for _, file := range tree {
		for _, ext := range PRODUCT_WHITE_EXTENTION {
			if strings.Contains(file, ext) {
				subTree = append(subTree, strings.Replace(file, pkgPath+"/", "", -1))
			}
		}
	}
	count := len(subTree)

	return map[string]interface{}{
		"count": count,
		"list":  subTree,
	}
}

func ServiceFile(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	file := ctx.URLParam("file")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	if file == "" {
		paramErrs.AppendError("$", fmt.Errorf("file is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	targetFile := filepath.Join(base.WebRoot, productName, productVersion, serviceName, file)

	fi, err := os.Open(targetFile)
	defer fi.Close()
	if err != nil {
		log.Errorf("[Product->ServiceFile] get service file err: %v", err)
		return err
	}
	content, err := ioutil.ReadAll(fi)

	if err != nil {
		log.Errorf("[Product->ServiceFile] read service file err: %v", err)
		return err
	}
	return string(content[:])
}

type serviceUpdateParam struct {
	ProductVersion string `json:"product_version"`
	FieldPath      string `json:"field_path"`
	Field          string `json:"field"`
}

func ServiceUpdate(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	params := make([]serviceUpdateParam, 0)
	if err := ctx.ReadJSON(&params); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}

	for _, param := range params {
		if param.ProductVersion == "" {
			return fmt.Errorf("product_version is empty")
		}
		info, err := model.DeployProductList.GetByProductNameAndVersion(productName, param.ProductVersion)
		if err != nil {
			log.Errorf("[Product->ServiceUpdate] read service file err:%v", err)
			return err
		}
		sc, err := schema.Unmarshal(info.Product)
		if err != nil {
			log.Errorf("[Product->ServiceGroup] Unmarshal err: %v", err)
			return err
		}
		for name, _ := range sc.Service {

			_, err := sc.SetField(name+"."+param.FieldPath, param.Field)
			if err != nil {
				log.Errorf("[Product->ServiceUpdate] read service file err:%v", err)
				continue
			}
			err = serviceUpdateDeployModifySchema(productName, name, param.FieldPath, param.Field, clusterId)
			if err != nil {
				log.Errorf("[Product->ServiceUpdate] read service file err:%v", err)
				return err
			}

			clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
			if err != nil {
				log.Errorf("%v\n", err)
				continue
			}
			if err := addSafetyAuditRecord(ctx, "集群运维", "服务参数修改", "集群名称："+clusterInfo.Name+", 组件名称："+productName+param.ProductVersion+
				", 服务组："+sc.Service[name].Group+", 服务名称："+name+sc.Service[name].Version+", 运行参数："+param.FieldPath+param.Field); err != nil {
				log.Errorf("failed to add safety audit record\n")
			}
		}

		schema, err := json.Marshal(sc)
		if err != nil {
			log.Errorf("[Product->ServiceUpdate] read service file err:%v", err)
			return err
		}
		query := "UPDATE " + model.DeployProductList.TableName + " SET `product`=? WHERE product_name=? AND product_version=? AND product!=?"
		if _, err := model.DeployProductList.GetDB().Exec(query, schema, sc.ProductName, sc.ProductVersion, schema); err != nil {
			log.Errorf("[Product->ServiceUpdate] read service file err:%v", err)
			return err
		}
	}
	return nil
}

func getRelativeHosts(productName string, serviceNames []string) ([]string, error) {
	hosts := []string{}
	//find the largest hosts of relative service
	for _, serviceName := range serviceNames {
		var ipList string
		query := fmt.Sprintf("SELECT ip_list FROM %s WHERE product_name=? AND service_name=?",
			model.DeployServiceIpList.TableName)
		log.Debugf("%v", query)
		if err := model.USE_MYSQL_DB().Get(&ipList, query, productName, serviceName); err != nil && err != sql.ErrNoRows {
			log.Errorf("%v", err)
			return hosts, err
		}
		ips := strings.Split(ipList, IP_LIST_SEP)
		if len(ipList) > len(hosts) {
			hosts = ips
		}
	}
	return hosts, nil
}

func selectCandidateHost(replica int, relatives, available []string) ([]string, error) {
	ipList := []string{}
	if len(relatives) < replica {
		ipList = append(ipList, relatives...)
	} else if len(relatives) == replica {
		ipList = relatives
	} else {
		ipList = relatives[:replica]
	}
	if len(ipList) < replica {
		ipList = append(ipList, available[:(replica-len(ipList))]...)
	}
	return ipList, nil
}

func doServiceGraphy(config schema.ServiceConfig, productName, serviceName string, clusterId int) error {
	maxReplica, err := strconv.Atoi(config.Instance.MaxReplica)
	if maxReplica == 0 {
		return nil
	}
	if err != nil {
		log.Errorf("[Product->ServiceGraphy] parse maxReplica error: %v", err)
		return err
	}
	relativeHosts, err := getRelativeHosts(productName, config.Relatives)
	if err != nil {
		log.Errorf("[Product->ServiceGraphy] get relative hosts error: %v", err)
		return err
	}

	hostInfo := []model.HostInfo{}
	ipFilter := "'',"
	for _, ip := range relativeHosts {
		ipFilter = ipFilter + "'" + ip + "',"
	}
	if ipFilter != "" {
		ipFilter = ipFilter[:len(ipFilter)-1]
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE steps=3 AND updated>SUBDATE(NOW(),INTERVAL 3 MINUTE) AND %s.ip NOT IN ("+ipFilter+")",
		model.DeployHostList.TableName,
		model.DeployHostList.TableName,
	)
	log.Debugf("%v", query)
	if err := model.USE_MYSQL_DB().Select(&hostInfo, query); err != nil {
		log.Errorf("%v", err)
		return err
	}
	availableHosts := []string{}
	for _, host := range hostInfo {
		availableHosts = append(availableHosts, host.Ip)
	}
	ipList, err := selectCandidateHost(maxReplica, relativeHosts, availableHosts)
	if err != nil {
		log.Errorf("[Product->ServiceGraphy] selectCandidateHost error: %v", err)
		return err
	}
	if err = model.DeployServiceIpList.SetServiceIp(productName, serviceName, strings.Join(ipList, IP_LIST_SEP), clusterId, ""); err != nil {
		log.Errorf("[SetIP] SetServiceIp err: %v", err)
		return err
	}
	return nil
}

func ServiceGraphy(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(productName)
	if err != nil {
		return err
	}
	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		return err
	}
	sc, err := schema.Unmarshal(info.Product)
	if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
		log.Errorf("[Product->ServiceGraphy] setSchemaFieldServiceAddr err: %v", err)
		return err
	}
	if err != nil {
		log.Errorf("[Product->ServiceGraphy] Unmarshal err: %v", err)
		return err
	}
	return doServiceGraphy(sc.Service[serviceName], productName, serviceName, clusterId)
}

func ServicesGraphy(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	uncheckedServices := strings.Split(ctx.URLParam("unchecked_services"), ",")
	paramErrs.CheckAndThrowApiParameterErrors()

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(productName)
	if err != nil {
		return err
	}
	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		return err
	}
	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("[Product->ServicesGraphy] Unmarshal err: %v", err)
		return err
	}
	if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
		log.Errorf("[Product->Current] setSchemaFieldServiceAddr err: %v", err)
		return err
	}
	for name, svc := range sc.Service {
		if uncheckedServices != nil && util.StringContain(uncheckedServices, name) {
			continue
		}
		if svc.Instance == nil || svc.Instance.UseCloud {
			continue
		}
		err = doServiceGraphy(svc, productName, name, clusterId)
		if err != nil {
			log.Errorf("[Product->ServicesGraphy] doServiceGraphy err: %v", err)
			return err
		}
	}
	return nil
}

func ServiceConfigFiles(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		return err
	}
	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("[Product->Service] Unmarshal err: %v", err)
		return err
	}
	configList := []string{}
	if _, ok := sc.Service[serviceName]; ok && sc.Service[serviceName].Instance != nil {
		configList = sc.Service[serviceName].Instance.ConfigPaths
	}
	return map[string]interface{}{
		"list":  configList,
		"count": len(configList),
	}
}

func ServiceConfigFile(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	file := ctx.URLParam("file")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	if file == "" {
		paramErrs.AppendError("$", fmt.Errorf("file is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	requestedFile := filepath.Clean(filepath.Join("/", productName, productVersion, serviceName, file))
	rel, err := filepath.Rel("/", requestedFile)
	if err != nil {
		log.Errorf("failed to get the relative path, err: %v", err)
		return err
	}
	absWebRootDir, err := filepath.Abs(base.WebRoot)
	if err != nil {
		log.Errorf("failed to get web root absolute path, err: %v", err)
		return err
	}
	targetFile := filepath.Join(absWebRootDir, rel)
	fi, err := os.Open(targetFile)
	defer fi.Close()
	if err != nil {
		log.Errorf("[Product->ServiceFile] get service file err: %v", err)
		return err
	}
	content, err := ioutil.ReadAll(fi)

	if err != nil {
		log.Errorf("[Product->ServiceFile] read service file err: %v", err)
		return err
	}
	return string(content[:])
}

func ConfigUpdate(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("[Product->ConfigUpdate] ConfigUpdate from EasyMatrix API ")
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	param := ConfigUpdateParam{Values: make(map[string]interface{})}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if len(param.Values) == 0 && len(param.Deleted) == 0 {
		log.Infof("[Product->ConfigUpdate] values and deleted is nil: %v", param)
		return nil
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	var query string
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(productName)
	if err != nil {
		return fmt.Errorf("[ConfigUpdate] check product ready for deploy error:%s", err)
	}
	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		return fmt.Errorf("[ConfigUpdate] get by product name and version error:%s", err)
	}
	// 反序列化解析为schema的结构体
	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("[Product->ConfigUpdate] Unmarshal err: %v", err)
		return err
	}
	//原始schema
	//2.8.2以前版本升级需要做数据库兼容 update deploy_product_list set `schema`=`product`
	panguSc, err := schema.Unmarshal(info.Schema)
	if err != nil {
		log.Errorf("[Product->ConfigUpdate] Unmarshal Pangu schema err: %v", err)
		return err
	}
	/* 参数删除逻辑：
	 * 1、判断要删除的参数是否为schema中默认配置的参数，若不是默认参数执行删除逻辑，
	 * 2、删除product中的参数和deploy_schema_field_modify里的该参数记录
	 */
	tx := model.USE_MYSQL_DB().MustBegin()
	deletedList := strings.Split(param.Deleted, ",")
	for _, key := range deletedList {
		if _, ok := panguSc.Service[serviceName].Config[key]; ok {
			log.Errorf("[Product->ConfigUpdate] init config %v cant be deleted", key)
			return fmt.Errorf("默认运行参数不能删除: %v", key)
		}
		delete(sc.Service[serviceName].Config, key)
		// 在deploy_schema_field_modify里 删除参数,namespace用于区分k8s模式和主机模式
		query = "DELETE FROM " + model.DeploySchemaFieldModify.TableName + " WHERE product_name=? AND service_name=? AND field_path=? AND cluster_id=? AND namespace=?"
		if _, err := tx.Exec(query, productName, serviceName, "Config."+key+".Value", clusterId, namespace); err != nil {
			log.Errorf("[Product->ConfigUpdate] delete record err:%v", err)
			return fmt.Errorf("[ConfigUpdate] delete schema field error:%s", err)
		}

		//query = "DELETE FROM " + model.DeploySchemaFieldModify.TableName + " WHERE product_name=? AND service_name=? AND field_path=? AND cluster_id=?"
		//if _, err := tx.Exec(query, productName, serviceName, "Config."+key+".Value", clusterId); err != nil {
		//	log.Errorf("[Product->ConfigUpdate] delete record err:%v", err)
		//	return err
		//}

		// 在deploy_schema_multi_fields 删除参数
		query = "update " + model.SchemaMultiField.TableName + " set is_deleted=1 where product_name=? and service_name=? and field_path=? and cluster_id=? and is_deleted=0"
		if _, err := tx.Exec(query, productName, serviceName, "Config."+key+".Value", clusterId); err != nil {
			log.Errorf("[Product->ConfigUpdate] delete multi record err:%v", err)
			return fmt.Errorf("[ConfigUpdate] delete schema mulit field error:%s", err)
		}
	}

	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()

	/* 参数新增逻辑：
	 * 1、为反序列化的product添加指定服务组件的参数配置
	 * 2、判断新增的字段合法性，如：是否为内置字段
	 * 3、序列化product，更新deploy_product_list表中的product字段
	 * 4、获取服务组件的配置文件，写入新增的配置参数
	 * 5、将新增的配置参数记录添加到deploy_schema_field_modify表中
	 */
	replaceServiceConfig(sc, serviceName, param.Values)

	err = sc.CheckServiceConfig(sc.Service[serviceName].Config)
	if err != nil {
		log.Errorf("[Product->ConfigUpdate] wrong values:%v", err)
		return err
	}
	schema, err := json.Marshal(sc)
	if err != nil {
		log.Errorf("[Product->ConfigUpdate] read service file err:%v", err)
		return err
	}
	// 更新deploy_product_list表中的product字段
	query = "UPDATE " + model.DeployProductList.TableName + " SET `product`=? WHERE product_name=? AND product_version=? AND product!=?"
	if _, err := tx.Exec(query, schema, sc.ProductName, sc.ProductVersion, schema); err != nil {
		log.Errorf("[Product->ConfigUpdate] read service file err:%v", err)
		return fmt.Errorf("[ConfigUpdate] update product schema error:%s", err)
	}
	if param.File != "" {
		// 获取要更新配置参数的文件路径
		targetFile := filepath.Join(base.WebRoot, productName, productVersion, serviceName, param.File)

		// 将新增的配置参数写入指定配置文件前先做备份
		fi, err := util.Create(targetFile)
		defer fi.Close()
		if err != nil {
			log.Errorf("[Product->ConfigUpdate] err:%v", err)
			return err
		}
		_, err = fi.Write([]byte(param.Content))
		if err != nil {
			log.Errorf("[Product->ConfigUpdate] err:%v", err)
			return err
		}
	}

	//在deploy_schema_field_modify里添加新参数
	for k, v := range param.Values {
		record := model.SchemaFieldModifyInfo{
			ClutserId:   clusterId,
			ProductName: productName,
			ServiceName: serviceName,
			FieldPath:   "Config." + k + ".Value",
			Field:       v.(string),
			Namespace:   namespace,
		}
		rlt = modifyField(&record, tx)
		if _, ok := rlt.(error); ok {
			return rlt
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[Product->ConfigUpdate] err:%v", err)
		return err
	}
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		paramsContent := ""
		for k, v := range param.Values {
			paramsContent += k + ":" + v.(string)
		}
		if err := addSafetyAuditRecord(ctx, "集群运维", "服务参数新增", "集群名称："+clusterInfo.Name+", 组件名称："+productName+productVersion+
			", 服务名称："+serviceName+", 部署参数："+paramsContent); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()

	return nil
}

func serviceUpdateDeployModifySchema(product_name, service_name, field_path, field string, clusterId int) error {
	info := model.SchemaFieldModifyInfo{
		ClutserId:   clusterId,
		ProductName: product_name,
		ServiceName: service_name,
		FieldPath:   field_path,
		Field:       field,
	}
	rlt := CommonModifySchemaField(&info)

	if err, ok := rlt.(error); ok {
		return err
	}
	return nil
}

func Current(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	// diff host model and k8s model
	var info *model.DeployProductListInfo

	info, err = model.DeployClusterProductRel.GetCurrentProductByProductNameClusterIdNamespace(productName, clusterId, namespace)
	if err != nil {
		return fmt.Errorf("[Current] Database query error %v", err)
	}

	ret := map[string]interface{}{}

	if err == sql.ErrNoRows {
		return ret
	}
	if err != nil {
		log.Errorf("[Product->Current] GetCurrentProductInfoByName err: %v", err)
		return err
	}

	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("[Product->Current] Unmarshal err: %v", err)
		return err
	}
	if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
		log.Errorf("[Product->Current] inheritBaseService warn: %+v", err)
	}
	if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), namespace); err != nil {
		log.Errorf("[Product->Current] setSchemaFieldServiceAddr err: %v", err)
		return err
	}
	if err = handleUncheckedServices(sc, info.ID, clusterId, namespace); err != nil {
		log.Errorf("[Product->Current] handleUncheckedServices warn: %+v", err)
	}
	if err = sc.ParseVariable(); err != nil {
		log.Errorf("[Product->Current] product info err: %v", err)
		return err
	}

	ret["id"] = info.ID
	ret["product_name"] = info.ProductName
	ret["product_version"] = info.ProductVersion
	ret["product"] = sc
	ret["is_current_version"] = 1
	ret["status"] = info.Status

	if info.DeployTime.Valid == true {
		ret["deploy_time"] = info.DeployTime.Time.Format(base.TsLayout)
	} else {
		ret["deploy_time"] = ""
	}

	if info.CreateTime.Valid == true {
		ret["create_time"] = info.CreateTime.Time.Format(base.TsLayout)
	} else {
		ret["create_time"] = ""
	}

	return ret
}

func ServicesStatus(productName string, clusterId int) apibase.Result {
	//log.Debugf("ServicesStatus: %v", ctx.Request().RequestURI)

	var query, status string
	var count, productPid, healthState, serviceCount int
	var serviceNames []string
	var err error

	type serviceBaseInfo struct {
		Group            string `db:"group" json:"group"`
		StatusCount      int    `json:"status_count"`
		HealthStateCount int    `json:"health_state_count"`
		ServiceName      string `json:"service_name"`
		Status           string `json:"status"`
		HealthState      int    `json:"health_state"`
	}

	type serviceInfo struct {
		serviceBaseInfo
		HealthStateOkCount int `db:"health_state_ok_count"`
		StatusOkCount      int `db:"status_ok_count"`
		ServiceCount       int `db:"service_count"`
	}

	var svi serviceInfo
	var serviceList []serviceBaseInfo

	//productName := ctx.Params().Get("product_name")

	product, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"count": 0,
			"list":  nil,
		}
	} else if err != nil {
		log.Errorf("err: %v", err.Error())
		apibase.ThrowDBModelError(err)
	}
	productPid = product.ID
	query = `SELECT COUNT(DISTINCT deploy_instance_list.service_name) FROM deploy_instance_list WHERE deploy_instance_list.pid=? AND cluster_id=?`
	if err = model.USE_MYSQL_DB().Get(&count, query, productPid, clusterId); err != nil {
		log.Errorf("service kind count query: %v, values %v, err: %v", query, productPid, err)
		apibase.ThrowDBModelError(err)
	}

	if count > 0 {
		query = `SELECT DISTINCT deploy_instance_list.service_name FROM deploy_instance_list WHERE deploy_instance_list.pid=? AND cluster_id=?`
		if err = model.USE_MYSQL_DB().Select(&serviceNames, query, productPid, clusterId); err != nil {
			log.Errorf("service name query: %v, values %v, err: %v", query, productPid, err)
			apibase.ThrowDBModelError(err)
		}
		for _, sn := range serviceNames {
			query = `SELECT DISTINCT deploy_instance_list.group, COUNT(deploy_instance_list.service_name) AS service_count FROM deploy_instance_list WHERE deploy_instance_list.service_name=? AND cluster_id=? AND pid=?`
			if err = model.USE_MYSQL_DB().Get(&svi, query, sn, clusterId, productPid); err != nil {
				log.Errorf("instance query: %v, values %v, err: %v", query, sn, err)
				apibase.ThrowDBModelError(err)
			}
			//running中的服务
			query = `SELECT COUNT(*) AS status_ok_count FROM deploy_instance_list WHERE deploy_instance_list.service_name=? AND deploy_instance_list.status=? AND cluster_id=? AND pid=?`
			if err = model.USE_MYSQL_DB().Get(&svi.StatusOkCount, query, sn, RUNNING, clusterId, productPid); err != nil {
				log.Errorf("statusOkCount query: %v, values %v, err: %v", query, sn, err)
				apibase.ThrowDBModelError(err)
			}
			query = `SELECT COUNT(*) AS health_state_ok_count FROM deploy_instance_list WHERE deploy_instance_list.service_name=? AND deploy_instance_list.health_state=? AND cluster_id=? AND pid=?`
			if err = model.USE_MYSQL_DB().Get(&svi.HealthStateOkCount, query, sn, HEALTHY, clusterId, productPid); err != nil {
				log.Errorf("healthStateOkCount query: %v, values %v, err: %v", query, sn, err)
				apibase.ThrowDBModelError(err)
			}

			svcList, err := model.DeployInstanceList.GetInstanceListByPidServiceName(productPid, clusterId, sn)
			if err != nil {
				log.Errorf("err: %v", err.Error())
				apibase.ThrowDBModelError(err)
			}

			var alertCount int
			svcMap := make(map[string]struct{}, 0)
			for _, svc := range svcList {
				if _, ok := svcMap[svc.ServiceName]; !ok {
					var alertState string
					var execFail bool
					alertState, err = getInstanceAlertState(productName, svc.ServiceName, "")
					if err != nil {
						log.Errorf("%v", err)
					}
					infoList, err := model.HealthCheck.GetInfoByClusterIdAndProductNameAndServiceName(clusterId, productName, svc.ServiceName, svc.Ip)
					if err != nil {
						log.Errorf("%v", err)
					}
					for _, info := range infoList {
						if info.ExecStatus == enums.ExecStatusType.Failed.Code {
							execFail = true
							break
						}
					}
					if alertState == "alert" || execFail == true {
						alertCount++
					}
				}
				svcMap[svc.ServiceName] = struct{}{}
			}
			healthState = HEALTHY
			status = NORMAL
			svi.HealthStateCount = svi.ServiceCount
			svi.StatusCount = svi.ServiceCount

			sc, err := schema.Unmarshal(product.Schema)
			if err != nil {
				return err
			}

			if srv, exist := sc.Service[sn]; exist && srv.Instance != nil && srv.Instance.HealthCheck == nil || srv.Instance.HealthCheck.Shell == "" {
				goto SUBSEQUENCE
			}
			if svi.HealthStateOkCount != svi.ServiceCount {
				svi.HealthStateCount = svi.ServiceCount - svi.HealthStateOkCount
				healthState = UNHEALTHY
			}
		SUBSEQUENCE:
			if svi.StatusOkCount != svi.ServiceCount {
				svi.StatusCount = svi.ServiceCount - svi.StatusOkCount
				status = ABNORMAL
			}

			if alertCount != 0 && status == NORMAL {
				svi.StatusCount = alertCount
				status = ABNORMAL
			}
			if status == NORMAL && healthState == 1 {
				continue
			}
			serviceList = append(serviceList, serviceBaseInfo{
				ServiceName:      sn,
				Group:            svi.Group,
				Status:           status,
				StatusCount:      svi.StatusCount,
				HealthState:      healthState,
				HealthStateCount: svi.HealthStateCount,
			})
			serviceCount += 1
		}
	}

	return map[string]interface{}{
		"count": serviceCount,
		"list":  serviceList,
	}
}

type rollingUpdate struct {
	schema            *schema.SchemaConfig
	svcMap            map[string]*sync.Once
	deployUUID        uuid.UUID
	operationId       string
	pid               int
	clusterId         int
	installMode       int
	uncheckedServices []string
	finalUpgrade      bool

	ctx       sysContext.Context
	rateLimit chan struct{}

	errMu sync.Mutex
	err   error

	sync.WaitGroup
}

type rollingUninstall struct {
	schema      *schema.SchemaConfig
	svcMap      map[string]*sync.Once
	deployUUID  uuid.UUID
	operationId string
	pid         int
	clusterId   int

	errMu sync.Mutex
	err   error

	sync.WaitGroup
}

func newRollingUpdate(sc *schema.SchemaConfig, deployUUID uuid.UUID, pid int, ctx sysContext.Context, uncheckedServices []string, clusterId int, installMode int, operationId string, finalUpgrade bool) *rollingUpdate {
	return &rollingUpdate{
		schema:            sc,
		svcMap:            make(map[string]*sync.Once, len(sc.Service)),
		deployUUID:        deployUUID,
		operationId:       operationId,
		pid:               pid,
		clusterId:         clusterId,
		installMode:       installMode,
		ctx:               ctx,
		rateLimit:         make(chan struct{}, rateLimit),
		uncheckedServices: uncheckedServices,
		finalUpgrade:      finalUpgrade,
	}
}

func newRollingUninstall(clusterId int, sc *schema.SchemaConfig, deployUUID uuid.UUID, pid int) *rollingUninstall {
	return &rollingUninstall{
		clusterId:  clusterId,
		schema:     sc,
		svcMap:     make(map[string]*sync.Once, len(sc.Service)),
		deployUUID: deployUUID,
		pid:        pid,
	}
}

type instanceErr struct {
	id       int64
	status   string
	progress uint
	err      error
}

func (e instanceErr) Error() string {
	return e.err.Error()
}

func (r *rollingUpdate) acquireForLimit() bool {
	select {
	case <-r.ctx.Done():
		return false
	default:
	}

	select {
	case r.rateLimit <- struct{}{}:
		return true
	case <-r.ctx.Done():
		return false
	}
}

func (r *rollingUpdate) releaseForLimit() {
	<-r.rateLimit
}

func (r *rollingUpdate) createInstance(
	instancer instance.Instancer,
	id int64,
	newSchema *schema.SchemaConfig,
	svcName string,
	wg *sync.WaitGroup,
	updateRlt *int64,
	waitOneCh chan struct{}) (err error) {

	updateQuery := "UPDATE " + model.DeployInstanceRecord.TableName + " SET `status`=?, status_message=?, progress=?, update_time=NOW() WHERE id=?"

	defer func() {
		instancer.Clear()

		if err != nil {
			atomic.StoreInt64(updateRlt, 1)
			if e, exist := err.(instanceErr); exist {
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, e.status, e.Error(), e.progress, e.id); err != nil {
					log.Errorf("%v", err)
				}
			}
		}
		waitOneCh <- struct{}{}
		wg.Done()
	}()

	if err := instancer.Install(false); err != nil {
		log.Errorf("%v", err)
		r.releaseForLimit()
		return instanceErr{id, model.INSTANCE_STATUS_INSTALL_FAIL, 30, err}
	}

	r.releaseForLimit()

	if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_INSTALLED, "", 30, id); err != nil {
		log.Errorf("%v", err)
	}
	if newSchema.Service[svcName].Instance.StartAfterInstall {
		return nil
	}
	if err := instancer.Start(); err != nil {
		log.Errorf("%v", err)
		return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 50, err}
	}
	discover.FlushServiceDiscover()
	if newSchema.Service[svcName].Instance.HealthCheck == nil {
		if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 100, id); err != nil {
			log.Errorf("%v", err)
		}
		return nil
	}
	if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 50, id); err != nil {
		log.Errorf("%v", err)
	}
	log.Debugf("waiting instance(%d) GetStatusChan...", instancer.ID())
	ev := <-instancer.GetStatusChan()
	log.Debugf("end instance(%d) GetStatusChan", instancer.ID())
	switch ev.Type {
	case event.REPORT_EVENT_HEALTH_CHECK:
		if ev.Data.(*agent.HealthCheck).Failed {
			err := fmt.Errorf("health check failed")
			log.Errorf("%v", err)
			return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_FAIL, 80, err}
		}
		if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_HEALTH_CHECKED, "", 100, id); err != nil {
			log.Errorf("%v", err)
		}
	case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
		err := fmt.Errorf("health check cancelled")
		log.Errorf("%v", err)
		return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_CANCELLED, 50, err}
	case event.REPORT_EVENT_INSTANCE_ERROR:
		err := fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
		log.Errorf("%v", err)
		return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 60, err}
	}

	return nil
}

func rollingInstances(
	rollingInstance []model.InstanceAndProductInfo,
	serviceIp []string,
	r *rollingUpdate,
	svcName string,
	idToIndex map[uint]uint,
	updateRlt *int64) (err error) {

	updateQuery := "UPDATE " + model.DeployInstanceRecord.TableName + " SET `status`=?, status_message=?, progress=?, update_time=NOW() WHERE id=?"

	defer func() {
		if err != nil {
			atomic.StoreInt64(updateRlt, 1)
			if e, exist := err.(instanceErr); exist {
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, e.status, e.Error(), e.progress, e.id); err != nil {
					log.Errorf("%v", err)
				}
			}
		}
	}()

	for _, info := range rollingInstance {
		newSchema, err := schema.Clone(r.schema)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		node, err := model.GetServiceIpNode(r.clusterId, r.schema.ProductName, svcName, info.Ip)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		if err := newSchema.SetServiceNodeIP(svcName, util.FoundIpIdx(serviceIp, info.Ip), node.NodeId, idToIndex); err != nil {
			log.Errorf("%v", err)
			return err
		}
		// 若配置项设置了多个值，在此处替换掉
		multiFields, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndIp(r.clusterId, r.schema.ProductName, svcName, info.Ip)
		if err != nil {
			log.Errorf("%v", err)
			break
		}
		for _, multiField := range multiFields {
			newSchema.SetField(multiField.ServiceName+"."+multiField.FieldPath, multiField.Field)
		}

		if err := newSchema.ParseServiceVariable(svcName); err != nil {
			log.Errorf("%v", err)
			return err
		}

		// the ServiceAddrStruct can set nil after ParseServiceVariable
		svc := newSchema.Service[svcName]
		svc.ServiceAddr = nil
		newSchema.Service[svcName] = svc

		newSchemaJson, err := json.Marshal(newSchema.Service[svcName])
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		instanceRecordInfo := model.DeployInstanceRecordInfo{
			DeployUUID:         r.deployUUID,
			InstanceId:         info.ID,
			Sid:                info.Sid,
			Ip:                 info.Ip,
			ProductName:        newSchema.ProductName,
			ProductNameDisplay: newSchema.ProductNameDisplay,
			Group:              info.Group,
			ServiceName:        info.ServiceName,
			ServiceNameDisplay: info.ServiceNameDisplay,
		}
		insertQuery := "INSERT INTO " + model.DeployInstanceRecord.TableName +
			" (deploy_uuid, instance_id, sid, ip, product_name, product_name_display, product_version, `group`, service_name, service_name_display, service_version, `status`, progress) VALUES" +
			" (:deploy_uuid, :instance_id, :sid, :ip, :product_name, :product_name_display, :product_version, :group, :service_name, :service_name_display, :service_version, :status, :progress)"

		if bytes.Equal(newSchemaJson, info.Schema) {
			instancer, err := instance.NewInstancer(info.Pid, info.Ip, svcName, newSchema, r.operationId)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
			mergeStatus := info.Status
			if info.Status == model.INSTANCE_STATUS_RUNNING && info.HealthState == model.INSTANCE_HEALTH_OK {
				mergeStatus = model.INSTANCE_STATUS_HEALTH_CHECKED
			}
			if info.Status == model.INSTANCE_STATUS_RUNNING && info.HealthState == model.INSTANCE_HEALTH_BAD {
				mergeStatus = model.INSTANCE_STATUS_HEALTH_CHECK_FAIL
			}
			if info.Status == model.INSTANCE_STATUS_RUNNING && info.HealthState == model.INSTANCE_HEALTH_WAITING {
				mergeStatus = model.INSTANCE_STATUS_HEALTH_CHECK_WAITING
			}
			switch mergeStatus {
			case model.INSTANCE_STATUS_RUN_FAIL, model.INSTANCE_STATUS_HEALTH_CHECK_FAIL, model.INSTANCE_STATUS_HEALTH_CHECK_WAITING,
				model.INSTANCE_STATUS_STOPPED, model.INSTANCE_STATUS_STOP_FAIL:
				instanceRecordInfo.ProductVersion = info.ProductVersion
				instanceRecordInfo.ServiceVersion = info.ServiceVersion
				instanceRecordInfo.Status = model.INSTANCE_STATUS_STOPPING
				instanceRecordInfo.Progress = 0
				rlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
				if err != nil {
					log.Errorf("%v", err)
					instancer.Clear()
					return err
				}
				id, _ := rlt.LastInsertId()

				if err := instancer.Stop(); err != nil {
					log.Errorf("%v", err)
					instancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_STOP_FAIL, 30, err}
				}

				updateRecord := "UPDATE " + model.DeployInstanceRecord.TableName + " SET `status`=?, product_version=?, update_time=NOW() WHERE id=?"
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateRecord, model.INSTANCE_STATUS_STOPPED, newSchema.ProductVersion, id); err != nil {
					log.Errorf("%v", err)
				}

				if err := instancer.Start(); err != nil {
					log.Errorf("%v", err)
					instancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 50, err}
				}
				if newSchema.Service[svcName].Instance.HealthCheck == nil {
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 100, id); err != nil {
						log.Errorf("%v", err)
					}
					instancer.Clear()
					continue
				}
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 50, id); err != nil {
					log.Errorf("%v", err)
				}
				log.Debugf("waiting instance(%d) GetStatusChan...", instancer.ID())
				ev := <-instancer.GetStatusChan()
				log.Debugf("end instance(%d) GetStatusChan", instancer.ID())
				switch ev.Type {
				case event.REPORT_EVENT_HEALTH_CHECK:
					if ev.Data.(*agent.HealthCheck).Failed {
						err := fmt.Errorf("health check failed")
						log.Errorf("%v", err)
						instancer.Clear()
						return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_FAIL, 80, err}
					}
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_HEALTH_CHECKED, "", 100, id); err != nil {
						log.Errorf("%v", err)
					}
				case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
					err := fmt.Errorf("health check cancelled")
					log.Errorf("%v", err)
					instancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_CANCELLED, 50, err}
				case event.REPORT_EVENT_INSTANCE_ERROR:
					err := fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
					log.Errorf("%v", err)
					instancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 60, err}
				}
				instancer.Clear()
			case model.INSTANCE_STATUS_RUNNING, model.INSTANCE_STATUS_HEALTH_CHECKED:
				if info.Pid != r.pid {
					if err := instancer.SetPid(r.pid); err != nil {
						log.Errorf("%v", err)
						instancer.Clear()
						return err
					}
					discover.FlushServiceDiscover()
				}
				instancer.Clear()
			default:
				instanceRecordInfo.ProductVersion = info.ProductVersion
				instanceRecordInfo.ServiceVersion = info.ServiceVersion
				instanceRecordInfo.Status = model.INSTANCE_STATUS_UNINSTALLING
				instanceRecordInfo.Progress = 0
				rlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
				if err != nil {
					log.Errorf("%v", err)
					instancer.Clear()
					return err
				}
				id, _ := rlt.LastInsertId()

				if err := instancer.UnInstall(false); err != nil {
					log.Errorf("%v", err)
					instancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_UNINSTALL_FAIL, 80, err}
				}
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLED, "", 100, id); err != nil {
					log.Errorf("%v", err)
				}
				instancer.Clear()

				if !r.acquireForLimit() {
					err := errors.New(model.INSTANCE_STATUS_INSTALLING_CANCELLED)
					log.Errorf("%v", err)
					return err
				}

				// must NewInstancer after UnInstall
				newInstancer, err := instance.NewInstancer(r.pid, info.Ip, svcName, newSchema, r.operationId)
				if err != nil {
					log.Errorf("%v", err)
					r.releaseForLimit()
					return err
				}
				instanceRecordInfo.InstanceId = newInstancer.ID()
				instanceRecordInfo.ProductVersion = newSchema.ProductVersion
				instanceRecordInfo.ServiceVersion = newSchema.Service[svcName].Version
				instanceRecordInfo.Status = model.INSTANCE_STATUS_INSTALLING
				instanceRecordInfo.Progress = 0
				newRlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
				if err != nil {
					log.Errorf("%v", err)
					newInstancer.Clear()
					r.releaseForLimit()
					return err
				}
				newId, _ := newRlt.LastInsertId()

				newInstancer.SetMode(r.installMode)

				if err := newInstancer.Install(false); err != nil {
					log.Errorf("%v", err)
					newInstancer.Clear()
					r.releaseForLimit()
					return instanceErr{newId, model.INSTANCE_STATUS_INSTALL_FAIL, 30, err}
				}

				r.releaseForLimit()

				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_INSTALLED, "", 30, newId); err != nil {
					log.Errorf("%v", err)
				}
				if newSchema.Service[svcName].Instance.StartAfterInstall {
					newInstancer.Clear()
					continue
				}
				if err := newInstancer.Start(); err != nil {
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{newId, model.INSTANCE_STATUS_RUN_FAIL, 50, err}
				}
				discover.FlushServiceDiscover()
				if newSchema.Service[svcName].Instance.HealthCheck == nil {
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 100, newId); err != nil {
						log.Errorf("%v", err)
					}
					newInstancer.Clear()
					continue
				}
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 50, newId); err != nil {
					log.Errorf("%v", err)
				}
				log.Debugf("waiting instance(%d) GetStatusChan...", instancer.ID())
				ev := <-newInstancer.GetStatusChan()
				log.Debugf("end instance(%d) GetStatusChan", instancer.ID())
				switch ev.Type {
				case event.REPORT_EVENT_HEALTH_CHECK:
					if ev.Data.(*agent.HealthCheck).Failed {
						err := fmt.Errorf("health check failed")
						log.Errorf("%v", err)
						newInstancer.Clear()
						return instanceErr{newId, model.INSTANCE_STATUS_HEALTH_CHECK_FAIL, 80, err}
					}
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_HEALTH_CHECKED, "", 100, newId); err != nil {
						log.Errorf("%v", err)
					}
				case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
					err := fmt.Errorf("health check cancelled")
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_CANCELLED, 50, err}
				case event.REPORT_EVENT_INSTANCE_ERROR:
					err := fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{newId, model.INSTANCE_STATUS_RUN_FAIL, 60, err}
				}
				newInstancer.Clear()
			}
		} else {
			log.Debugf("<-------------------------->")
			log.Debugf("old schema: %s", info.Schema)
			log.Debugf("<==========================>")
			log.Debugf("new schema: %s", newSchemaJson)
			log.Debugf("<-------------------------->")

			newsc := schema.ServiceConfig{}
			if err := json.Unmarshal(newSchemaJson, &newsc); err != nil {
				log.Errorf("%v", err)
				return err
			}
			oldsc := schema.ServiceConfig{}
			if err := json.Unmarshal(info.Schema, &oldsc); err != nil {
				log.Errorf("%v", err)
				return err
			}
			oldInstancer, err := instance.NewInstancer(info.Pid, info.Ip, svcName, &schema.SchemaConfig{
				ProductName:    newSchema.ProductName,
				ProductVersion: info.ProductVersion,
				Service:        map[string]schema.ServiceConfig{svcName: oldsc},
			}, r.operationId)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}

			if !needReInstallAgent(&oldsc, &newsc) {
				if reflect.DeepEqual(oldsc.Config, newsc.Config) &&
					info.Status == model.INSTANCE_STATUS_RUNNING &&
					(info.HealthState == model.INSTANCE_HEALTH_NOTSET || info.HealthState == model.INSTANCE_HEALTH_OK) {
					// NewInstancer will update instance schema, just only update for this
					newInstancer, err := instance.NewInstancer(info.Pid, info.Ip, svcName, newSchema, r.operationId)
					if err != nil {
						log.Errorf("%v", err)
						return err
					}
					if info.Pid != r.pid {
						if err = newInstancer.SetPid(r.pid); err != nil {
							log.Errorf("%v", err)
							newInstancer.Clear()
							oldInstancer.Clear()
							return err
						}
						discover.FlushServiceDiscover()
					}
					newInstancer.Clear()
					oldInstancer.Clear()
					continue
				}

				var id int64
				var newInstancer instance.Instancer

				switch info.Status {
				case model.INSTANCE_STATUS_INSTALLING, model.INSTANCE_STATUS_INSTALLED, model.INSTANCE_STATUS_INSTALL_FAIL,
					model.INSTANCE_STATUS_UNINSTALLING, model.INSTANCE_STATUS_UNINSTALLED, model.INSTANCE_STATUS_UNINSTALL_FAIL:
					instanceRecordInfo.ProductVersion = info.ProductVersion
					instanceRecordInfo.ServiceVersion = info.ServiceVersion
					instanceRecordInfo.Status = model.INSTANCE_STATUS_UNINSTALLING
					instanceRecordInfo.Progress = 0
					rlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
					if err != nil {
						log.Errorf("%v", err)
						oldInstancer.Clear()
						return err
					}
					id, _ = rlt.LastInsertId()

					if err := oldInstancer.UnInstall(false); err != nil {
						log.Errorf("%v", err)
						oldInstancer.Clear()
						return instanceErr{id, model.INSTANCE_STATUS_UNINSTALL_FAIL, 80, err}
					}
					oldInstancer.Clear()

					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLED, "", 100, id); err != nil {
						log.Errorf("%v", err)
					}

					if !r.acquireForLimit() {
						err := errors.New(model.INSTANCE_STATUS_INSTALLING_CANCELLED)
						log.Errorf("%v", err)
						return err
					}

					newInstancer, err = instance.NewInstancer(info.Pid, info.Ip, svcName, newSchema, r.operationId)
					if err != nil {
						log.Errorf("%v", err)
						r.releaseForLimit()
						return err
					}

					instanceRecordInfo.InstanceId = newInstancer.ID()
					instanceRecordInfo.ProductVersion = newSchema.ProductVersion
					instanceRecordInfo.ServiceVersion = newSchema.Service[svcName].Version
					instanceRecordInfo.Status = model.INSTANCE_STATUS_INSTALLING
					instanceRecordInfo.Progress = 0
					rlt, err = model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
					if err != nil {
						log.Errorf("%v", err)
						newInstancer.Clear()
						r.releaseForLimit()
						return err
					}
					id, _ = rlt.LastInsertId()
					newInstancer.SetMode(r.installMode)
					if err := newInstancer.Install(false); err != nil {
						log.Errorf("%v", err)
						newInstancer.Clear()
						r.releaseForLimit()
						return instanceErr{id, model.INSTANCE_STATUS_INSTALL_FAIL, 30, err}
					}

					r.releaseForLimit()

					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_INSTALLED, "", 30, id); err != nil {
						log.Errorf("%v", err)
					}
					if newSchema.Service[svcName].Instance.StartAfterInstall {
						newInstancer.Clear()
						continue
					}
				default:
					log.Infof("--------------------->")
					log.Infof("enter default")
					if r.finalUpgrade {
						log.Infof("is final upgrade")
						if strings.HasSuffix(svcName, "Sql") || strings.HasSuffix(svcName, "sql") {
							instanceRecordInfo.ProductVersion = info.ProductVersion
							instanceRecordInfo.ServiceVersion = info.ServiceVersion
							instanceRecordInfo.Status = model.INSTANCE_STATUS_UNINSTALLING
							instanceRecordInfo.Progress = 0
							rlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
							if err != nil {
								log.Errorf("%v", err)
								oldInstancer.Clear()
								return err
							}
							id, _ = rlt.LastInsertId()

							if err := oldInstancer.UnInstall(false); err != nil {
								log.Errorf("%v", err)
								oldInstancer.Clear()
								return instanceErr{id, model.INSTANCE_STATUS_UNINSTALL_FAIL, 80, err}
							}
							oldInstancer.Clear()

							if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLED, "", 100, id); err != nil {
								log.Errorf("%v", err)
							}

							if !r.acquireForLimit() {
								err := errors.New(model.INSTANCE_STATUS_INSTALLING_CANCELLED)
								log.Errorf("%v", err)
								return err
							}

							newInstancer, err = instance.NewInstancer(info.Pid, info.Ip, svcName, newSchema, r.operationId)
							if err != nil {
								log.Errorf("%v", err)
								r.releaseForLimit()
								return err
							}

							instanceRecordInfo.InstanceId = newInstancer.ID()
							instanceRecordInfo.ProductVersion = newSchema.ProductVersion
							instanceRecordInfo.ServiceVersion = newSchema.Service[svcName].Version
							instanceRecordInfo.Status = model.INSTANCE_STATUS_INSTALLING
							instanceRecordInfo.Progress = 0
							rlt, err = model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
							if err != nil {
								log.Errorf("%v", err)
								newInstancer.Clear()
								r.releaseForLimit()
								return err
							}
							id, _ = rlt.LastInsertId()
							newInstancer.SetMode(r.installMode)
							if err := newInstancer.Install(false); err != nil {
								log.Errorf("%v", err)
								newInstancer.Clear()
								r.releaseForLimit()
								return instanceErr{id, model.INSTANCE_STATUS_INSTALL_FAIL, 30, err}
							}

							r.releaseForLimit()

							if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_INSTALLED, "", 30, id); err != nil {
								log.Errorf("%v", err)
							}
							if newSchema.Service[svcName].Instance.StartAfterInstall {
								newInstancer.Clear()
								continue
							}
						} else {
							continue
						}
					} else {
						log.Infof("-------------------------->")
						log.Infof("restart service")
						if r.installMode == 3 {
							log.Infof("skip restart service because of smooth upgrade.")
							continue
						}
						instanceRecordInfo.ProductVersion = info.ProductVersion
						instanceRecordInfo.ServiceVersion = info.ServiceVersion
						instanceRecordInfo.Status = model.INSTANCE_STATUS_STOPPING
						instanceRecordInfo.Progress = 0
						rlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
						if err != nil {
							log.Errorf("%v", err)
							oldInstancer.Clear()
							return err
						}
						id, _ = rlt.LastInsertId()

						if err := oldInstancer.Stop(); err != nil {
							log.Errorf("%v", err)
							oldInstancer.Clear()
							return instanceErr{id, model.INSTANCE_STATUS_STOP_FAIL, 30, err}
						}
						oldInstancer.Clear()

						newInstancer, err = instance.NewInstancer(info.Pid, info.Ip, svcName, newSchema, r.operationId)
						if err != nil {
							log.Errorf("%v", err)
							return err
						}
						if info.Pid != r.pid {
							if err = newInstancer.SetPid(r.pid); err != nil {
								log.Errorf("%v", err)
							}
						}
						updateRecord := "UPDATE " + model.DeployInstanceRecord.TableName + " SET `status`=?, product_version=?, update_time=NOW() WHERE id=?"
						if _, err := model.DeployInstanceRecord.GetDB().Exec(updateRecord, model.INSTANCE_STATUS_STOPPED, newSchema.ProductVersion, id); err != nil {
							log.Errorf("%v", err)
						}
						if err := newInstancer.UpdateConfig(); err != nil {
							log.Errorf("%v", err)
							newInstancer.Clear()
							return instanceErr{id, model.INSTANCE_STATUS_UPDATE_CONFIG_FAIL, 50, err}
						}
					}
				}

				if err := newInstancer.Start(); err != nil {
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 50, err}
				}
				discover.FlushServiceDiscover()
				if newSchema.Service[svcName].Instance.HealthCheck == nil {
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 100, id); err != nil {
						log.Errorf("%v", err)
					}
					newInstancer.Clear()
					continue
				}
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 50, id); err != nil {
					log.Errorf("%v", err)
				}
				log.Debugf("waiting instance(%d) GetStatusChan...", newInstancer.ID())
				ev := <-newInstancer.GetStatusChan()
				log.Debugf("end instance(%d) GetStatusChan", newInstancer.ID())
				switch ev.Type {
				case event.REPORT_EVENT_HEALTH_CHECK:
					if ev.Data.(*agent.HealthCheck).Failed {
						err := fmt.Errorf("health check failed")
						log.Errorf("%v", err)
						newInstancer.Clear()
						return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_FAIL, 80, err}
					}
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_HEALTH_CHECKED, "", 100, id); err != nil {
						log.Errorf("%v", err)
					}
				case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
					err := fmt.Errorf("health check cancelled")
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_CANCELLED, 50, err}
				case event.REPORT_EVENT_INSTANCE_ERROR:
					err := fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 60, err}
				}
				newInstancer.Clear()
			} else {
				instanceRecordInfo.ProductVersion = info.ProductVersion
				instanceRecordInfo.ServiceVersion = info.ServiceVersion
				instanceRecordInfo.Status = model.INSTANCE_STATUS_UNINSTALLING
				instanceRecordInfo.Progress = 0
				oldRlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
				if err != nil {
					oldInstancer.Clear()
					log.Errorf("%v", err)
					return err
				}
				oldId, _ := oldRlt.LastInsertId()

				switch info.Status {
				case model.INSTANCE_STATUS_RUN_FAIL, model.INSTANCE_STATUS_RUNNING, model.INSTANCE_STATUS_STOP_FAIL:
					if err := oldInstancer.Stop(); err != nil {
						log.Errorf("%v", err)
						oldInstancer.Clear()
						return instanceErr{oldId, model.INSTANCE_STATUS_STOP_FAIL, 30, err}
					}
				}
				if err := oldInstancer.UnInstall(oldsc.Version == newsc.Version); err != nil {
					log.Errorf("%v", err)
					oldInstancer.Clear()
					return instanceErr{oldId, model.INSTANCE_STATUS_UNINSTALL_FAIL, 80, err}
				}
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLED, "", 100, oldId); err != nil {
					log.Errorf("%v", err)
				}
				oldInstancer.Clear()

				if !r.acquireForLimit() {
					err := errors.New(model.INSTANCE_STATUS_INSTALLING_CANCELLED)
					log.Errorf("%v", err)
					return err
				}

				newInstancer, err := instance.NewInstancer(r.pid, info.Ip, svcName, newSchema, r.operationId)
				if err != nil {
					log.Errorf("%v", err)
					r.releaseForLimit()
					return err
				}
				instanceRecordInfo.InstanceId = newInstancer.ID()
				instanceRecordInfo.ProductVersion = newSchema.ProductVersion
				instanceRecordInfo.ServiceVersion = newSchema.Service[svcName].Version
				instanceRecordInfo.Status = model.INSTANCE_STATUS_INSTALLING
				instanceRecordInfo.Progress = 0
				rlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
				if err != nil {
					log.Errorf("%v", err)
					newInstancer.Clear()
					r.releaseForLimit()
					return err
				}
				id, _ := rlt.LastInsertId()
				newInstancer.SetMode(r.installMode)
				if err := newInstancer.Install(oldsc.Version == newsc.Version); err != nil {
					log.Errorf("%v", err)
					newInstancer.Clear()
					r.releaseForLimit()
					return instanceErr{id, model.INSTANCE_STATUS_INSTALL_FAIL, 30, err}
				}

				r.releaseForLimit()

				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_INSTALLED, "", 30, id); err != nil {
					log.Errorf("%v", err)
				}
				if newSchema.Service[svcName].Instance.StartAfterInstall {
					newInstancer.Clear()
					continue
				}
				if err := newInstancer.Start(); err != nil {
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 50, err}
				}
				discover.FlushServiceDiscover()
				if newSchema.Service[svcName].Instance.HealthCheck == nil {
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 100, id); err != nil {
						log.Errorf("%v", err)
					}
					newInstancer.Clear()
					continue
				}
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 50, id); err != nil {
					log.Errorf("%v", err)
				}
				log.Debugf("waiting instance(%d) GetStatusChan...", newInstancer.ID())
				ev := <-newInstancer.GetStatusChan()
				log.Debugf("end instance(%d) GetStatusChan", newInstancer.ID())
				switch ev.Type {
				case event.REPORT_EVENT_HEALTH_CHECK:
					if ev.Data.(*agent.HealthCheck).Failed {
						err := fmt.Errorf("health check failed")
						log.Errorf("%v", err)
						newInstancer.Clear()
						return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_FAIL, 80, err}
					}
					if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_HEALTH_CHECKED, "", 100, id); err != nil {
						log.Errorf("%v", err)
					}
				case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
					err := fmt.Errorf("health check cancelled")
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_HEALTH_CHECK_CANCELLED, 50, err}
				case event.REPORT_EVENT_INSTANCE_ERROR:
					err := fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
					log.Errorf("%v", err)
					newInstancer.Clear()
					return instanceErr{id, model.INSTANCE_STATUS_RUN_FAIL, 60, err}
				}
				newInstancer.Clear()
			}
		}
	}

	return nil
}

func needReInstallAgent(oldsc, newsc *schema.ServiceConfig) bool {
	if oldsc.Version != newsc.Version {
		return true
	}
	if oldsc.Instance == nil || newsc.Instance == nil {
		return true
	}
	if oldsc.Instance.Cmd != newsc.Instance.Cmd {
		return true
	}
	if !reflect.DeepEqual(oldsc.Instance.Environment, newsc.Instance.Environment) {
		return true
	}
	if !reflect.DeepEqual(oldsc.Instance.HealthCheck, newsc.Instance.HealthCheck) {
		return true
	}

	return false
}

func startInstalledInstances(clusterId int, productName, productVersion, serviceName string, deployUUID uuid.UUID, operationId string) error {
	var instanceInfo []model.DeployInstanceInfo
	query := fmt.Sprintf("SELECT %s.* FROM %s LEFT JOIN %s ON pid=%s.id "+
		"WHERE product_name=? AND product_version=? AND service_name=? AND %s.status=? AND cluster_id=?",
		model.DeployInstanceList.TableName,
		model.DeployInstanceList.TableName,
		model.DeployProductList.TableName,
		model.DeployProductList.TableName,
		model.DeployInstanceList.TableName,
	)
	log.Debugf("%v", query)
	if err := model.USE_MYSQL_DB().Select(&instanceInfo, query, productName, productVersion, serviceName, model.INSTANCE_STATUS_INSTALLED, clusterId); err != nil {
		log.Errorf("%v", err)
		return err
	}

	var startRlt int64 // 0 success, 1 fail
	var updateQuery = "UPDATE " + model.DeployInstanceRecord.TableName +
		" SET `status`=?, status_message=?, progress=?, update_time=NOW() WHERE deploy_uuid=? AND instance_id=?"

	wg := sync.WaitGroup{}
	for _, info := range instanceInfo {
		wg.Add(1)
		go func(info model.DeployInstanceInfo) (err error) {
			defer func() {
				if err != nil {
					atomic.StoreInt64(&startRlt, 1)
					if e, exist := err.(instanceErr); exist {
						if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, e.status, e.Error(), e.progress, deployUUID, e.id); err != nil {
							log.Errorf("%v", err)
						}
					}
				}
				wg.Done()
			}()

			sc := schema.ServiceConfig{}
			if err := json.Unmarshal(info.Schema, &sc); err != nil {
				log.Errorf("%v", err)
				return instanceErr{int64(info.ID), model.INSTANCE_STATUS_RUN_FAIL, 50, err}
			}
			instancer, err := instance.NewInstancer(info.Pid, info.Ip, serviceName, &schema.SchemaConfig{
				ProductName:    productName,
				ProductVersion: productVersion,
				Service:        map[string]schema.ServiceConfig{info.ServiceName: sc},
			}, operationId)
			if err != nil {
				log.Errorf("%v", err)
				return instanceErr{int64(info.ID), model.INSTANCE_STATUS_RUN_FAIL, 50, err}
			}
			if err := instancer.Start(); err != nil {
				log.Errorf("%v", err)
				return instanceErr{int64(info.ID), model.INSTANCE_STATUS_RUN_FAIL, 50, err}
			}
			discover.FlushServiceDiscover()
			if sc.Instance.HealthCheck == nil {
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 100, deployUUID, info.ID); err != nil {
					log.Errorf("%v", err)
				}
				return nil
			}
			if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_RUNNING, "", 50, deployUUID, info.ID); err != nil {
				log.Errorf("%v", err)
			}
			log.Debugf("waiting instance(%d) GetStatusChan...", instancer.ID())
			ev := <-instancer.GetStatusChan()
			log.Debugf("end instance(%d) GetStatusChan", instancer.ID())
			switch ev.Type {
			case event.REPORT_EVENT_HEALTH_CHECK:
				if ev.Data.(*agent.HealthCheck).Failed {
					err := fmt.Errorf("health check failed")
					log.Errorf("%v", err)
					return instanceErr{int64(info.ID), model.INSTANCE_STATUS_HEALTH_CHECK_FAIL, 80, err}
				}
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_HEALTH_CHECKED, "", 100, deployUUID, info.ID); err != nil {
					log.Errorf("%v", err)
				}
			case event.REPORT_EVENT_HEALTH_CHECK_CANCEL:
				err := fmt.Errorf("health check cancelled")
				log.Errorf("%v", err)
				return instanceErr{int64(info.ID), model.INSTANCE_STATUS_HEALTH_CHECK_CANCELLED, 50, err}
			case event.REPORT_EVENT_INSTANCE_ERROR:
				err := fmt.Errorf("%v", ev.Data.(*agent.AgentError).ErrStr)
				log.Errorf("%v", err)
				return instanceErr{int64(info.ID), model.INSTANCE_STATUS_RUN_FAIL, 60, err}
			}

			return nil
		}(info)
	}
	wg.Wait()

	if startRlt == 1 {
		//当不需要部署产品包的时候  没有 exec shell
		count, err := model.ExecShellList.GetCountByOperationId(operationId)
		if err != nil {
			log.Errorf("ExecShellList.GetCountByOperationId error: %v", err)
		}
		if *count == 0 {
			err := model.OperationList.UpdateStatusByOperationId(operationId, enums.ExecStatusType.Failed.Code, dbhelper.NullTime{Time: time.Now(), Valid: true}, sql.NullFloat64{Float64: 0, Valid: true})
			if err != nil {
				log.Errorf("OperationList UpdateStatusByOperationId error %v", err)
			}

		}
		return fmt.Errorf("some installed instance start fail")
	} else {
		//当不需要部署产品包的时候  没有 exec shell
		count, err := model.ExecShellList.GetCountByOperationId(operationId)
		if err != nil {
			log.Errorf("ExecShellList.GetCountByOperationId error: %v", err)
		}
		if *count == 0 {
			err := model.OperationList.UpdateStatusByOperationId(operationId, enums.ExecStatusType.Success.Code, dbhelper.NullTime{Time: time.Now(), Valid: true}, sql.NullFloat64{Float64: 0, Valid: true})
			if err != nil {
				log.Errorf("OperationList UpdateStatusByOperationId error %v", err)
			}

		}
	}

	return nil
}

func (r *rollingUpdate) rollingUpdateCore(svcName string) error {
	var ipList string
	query := "SELECT ip_list FROM " + model.DeployServiceIpList.TableName + " WHERE product_name=? AND service_name=? AND cluster_id=?"
	if err := model.DeployServiceIpList.GetDB().Get(&ipList, query, r.schema.ProductName, svcName, r.clusterId); err != nil {
		log.Errorf("%v", err)
		return err
	}
	serviceIp := strings.Split(ipList, IP_LIST_SEP)
	log.Debugf("cluster %v found %v old instance ip: %+v", r.clusterId, svcName, serviceIp)

	var instanceInfo []model.InstanceAndProductInfo
	query = fmt.Sprintf("SELECT %s.*, product_name, product_name_display, product_version FROM %s LEFT JOIN %s ON pid=%s.id WHERE product_name=? AND service_name=? AND cluster_id=?",
		model.DeployInstanceList.TableName,
		model.DeployInstanceList.TableName,
		model.DeployProductList.TableName,
		model.DeployProductList.TableName,
	)
	if err := model.USE_MYSQL_DB().Select(&instanceInfo, query, r.schema.ProductName, svcName, r.clusterId); err != nil {
		log.Errorf("%v", err)
		return err
	}

	newInstance, rollingInstance, oldInstance := filterInstance(serviceIp, instanceInfo)
	log.Debugf("found %v new instance ip: %+v", svcName, newInstance)

	idToIndex, err := assignNodeID(r.clusterId, r.schema.ProductName, svcName, serviceIp, newInstance, oldInstance, rollingInstance)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	//非平滑升级，卸载原服务
	if r.installMode != 3 && r.schema.Service[svcName].Instance.UpdateRecreate {
		if err := uninstallInstances(r.clusterId, oldInstance, r.deployUUID, r.operationId); err != nil {
			return err
		}
	}

	var updateRlt int64 // 0 success, 1 fail
	wg := sync.WaitGroup{}

	// start install new instance
	waitOneCh := make(chan struct{}, len(newInstance))
	if len(newInstance) == 0 {
		close(waitOneCh)
	}
	for _, ip := range newInstance {
		var newSchema *schema.SchemaConfig
		if newSchema, err = schema.Clone(r.schema); err != nil {
			log.Errorf("%v", err)
			break
		}
		var node *model.ServiceIpNode
		node, err = model.GetServiceIpNode(r.clusterId, r.schema.ProductName, svcName, ip)
		if err != nil {
			log.Errorf("%v", err)
			break
		}
		if err = newSchema.SetServiceNodeIP(svcName, util.FoundIpIdx(serviceIp, ip), node.NodeId, idToIndex); err != nil {
			log.Errorf("%v", err)
			break
		}
		// 若配置项设置了多个值，在此处替换掉
		multiFields, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndIp(r.clusterId, r.schema.ProductName, svcName, ip)
		if err != nil {
			log.Errorf("%v", err)
			break
		}
		for _, multiField := range multiFields {
			newSchema.SetField(multiField.ServiceName+"."+multiField.FieldPath, multiField.Field)
		}
		if err = newSchema.ParseServiceVariable(svcName); err != nil {
			log.Errorf("%v", err)
			break
		}

		// the ServiceAddr can set nil after ParseServiceVariable
		// the ServiceAddrStruct can set nil after ParseServiceVariable
		svc := newSchema.Service[svcName]
		svc.ServiceAddr = nil
		newSchema.Service[svcName] = svc

		if !r.acquireForLimit() {
			err = errors.New(model.INSTANCE_STATUS_INSTALLING_CANCELLED)
			break
		}

		var instancer instance.Instancer
		if instancer, err = instance.NewInstancer(r.pid, ip, svcName, newSchema, r.operationId); err != nil {
			log.Errorf("%v", err)
			r.releaseForLimit()
			break
		}
		var info *model.DeployInstanceInfo
		if err, info = instancer.GetInstanceInfo(); err != nil {
			log.Errorf("%v", err)
			instancer.Clear()
			r.releaseForLimit()
			break
		}
		instanceRecordInfo := model.DeployInstanceRecordInfo{
			DeployUUID:         r.deployUUID,
			InstanceId:         info.ID,
			Sid:                info.Sid,
			Ip:                 info.Ip,
			ProductName:        newSchema.ProductName,
			ProductNameDisplay: newSchema.ProductNameDisplay,
			ProductVersion:     newSchema.ProductVersion,
			Group:              info.Group,
			ServiceName:        info.ServiceName,
			ServiceNameDisplay: info.ServiceNameDisplay,
			ServiceVersion:     info.ServiceVersion,
			Status:             model.INSTANCE_STATUS_INSTALLING,
			Progress:           0,
		}
		insertQuery := "INSERT INTO " + model.DeployInstanceRecord.TableName +
			" (deploy_uuid, instance_id, sid, ip, product_name, product_name_display, product_version, `group`, service_name, service_name_display, service_version, `status`, progress) VALUES" +
			" (:deploy_uuid, :instance_id, :sid, :ip, :product_name, :product_name_display, :product_version, :group, :service_name, :service_name_display, :service_version, :status, :progress)"
		var rlt sql.Result
		if rlt, err = model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo); err != nil {
			log.Errorf("%v", err)
			instancer.Clear()
			r.releaseForLimit()
			break
		}
		id, _ := rlt.LastInsertId()

		wg.Add(1)
		go r.createInstance(instancer, id, newSchema, svcName, &wg, &updateRlt, waitOneCh)
	}

	if err != nil {
		wg.Wait()
		return err
	}

	<-waitOneCh
	if updateRlt == 1 {
		wg.Wait()
		return fmt.Errorf("some instance of %v update fail", svcName)
	}

	// start rolling instance
	rollingInstances(rollingInstance, serviceIp, r, svcName, idToIndex, &updateRlt)

	wg.Wait()
	if updateRlt == 1 {
		return fmt.Errorf("cluster %v some instance of %v update fail", r.clusterId, svcName)
	}

	return startInstalledInstances(r.clusterId, r.schema.ProductName, r.schema.ProductVersion, svcName, r.deployUUID, r.operationId)
}

// rollingUpdate should be goroutine and recursive
func (r *rollingUpdate) rollingUpdate(svcName string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	r.svcMap[svcName].Do(func() {
		wgThis := sync.WaitGroup{}
		for _, dname := range r.schema.Service[svcName].DependsOn {
			wgThis.Add(1)
			go r.rollingUpdate(dname, &wgThis)
		}
		wgThis.Wait()

		if err = r.getError(); err != nil {
			return
		}

		if r.schema.Service[svcName].Instance != nil && !r.schema.Service[svcName].Instance.UseCloud && !util.StringContain(r.uncheckedServices, svcName) {
			log.Infof("cluster %v rollingUpdateCore %v ...", r.clusterId, svcName)

			if err = r.rollingUpdateCore(svcName); err != nil {
				r.setError(err)
			}

			log.Infof("rollingUpdateCore %v finish(%v)", svcName, err)
		}
	})

	return
}

func (r *rollingUpdate) setError(err error) {
	r.errMu.Lock()
	defer r.errMu.Unlock()

	r.err = err
}

func (r *rollingUpdate) getError() error {
	r.errMu.Lock()
	defer r.errMu.Unlock()

	return r.err
}

func (r *rollingUpdate) run() error {
	for name := range r.schema.Service {
		r.svcMap[name] = &sync.Once{}
	}
	for name := range r.schema.Service {
		r.Add(1)
		go r.rollingUpdate(name, &r.WaitGroup)
	}
	r.Wait()

	return r.getError()
}

func deploy(sc *schema.SchemaConfig, deployUUID uuid.UUID, pid int, ctx sysContext.Context, uncheckedServices []string, clusterId int, installMode int, operationId, sourceVersion string, finalUpgrade bool) {
	var err error
	var query string

	defer func() {
		var status = model.PRODUCT_STATUS_DEPLOYED
		var productParsed []byte

		if err != nil {
			status = model.PRODUCT_STATUS_DEPLOY_FAIL
		}

		if err = sc.ParseVariable(); err == nil {
			if productParsed, err = json.Marshal(sc); err != nil {
				log.Errorf("%v", err)
			}
		} else {
			log.Errorf("%v", err)
			status = model.PRODUCT_STATUS_DEPLOY_FAIL
		}
		sourceProduct, err := model.DeployProductList.GetByProductNameAndVersion(sc.ProductName, sourceVersion)
		if err != nil {
			log.Errorf("%v", err)
		}
		switch installMode {
		case 3:
			query = "UPDATE " + model.DeployClusterSmoothUpgradeProductRel.TableName + " SET status=?, product_parsed=?, update_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
			if _, err := model.DeployClusterSmoothUpgradeProductRel.GetDB().Exec(query, status, productParsed, pid, clusterId); err != nil {
				log.Errorf("%v", err)
			}
			instances, err := model.DeployInstanceList.GetInstanceListByClusterId(clusterId, sourceProduct.ID)
			if err != nil {
				log.Errorf("%v", err)
			}
			if len(instances) == 0 {
				suRel, err := model.DeployClusterSmoothUpgradeProductRel.GetSmoothUpgradeProductRelByClusterIdAndPid(clusterId, pid)
				if err != nil {
					log.Errorf("%v", err)
				}
				query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET pid=?, deploy_uuid=?, `product_parsed`=?, `status`=?, `user_id`=?, deploy_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
				if _, err := model.DeployClusterProductRel.Exec(query, suRel.Pid, suRel.DeployUUID, suRel.ProductParsed, suRel.Status, suRel.UserId, sourceProduct.ID, clusterId); err != nil {
					log.Errorf("%v", err)
				}
				query = "UPDATE " + model.DeployClusterSmoothUpgradeProductRel.TableName + " SET is_deleted=1, update_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
				if _, err := model.DeployClusterProductRel.GetDB().Exec(query, pid, clusterId); err != nil {
					log.Errorf("%v", err)
				}
				query = "UPDATE " + upgrade.UpgradeHistory.TableName + " SET upgrade_mode=? WHERE cluster_id=? AND product_name=? AND source_version=? AND upgrade_mode=? AND is_deleted=0"
				if _, err := model.DeployClusterProductRel.GetDB().Exec(query, "", clusterId, sourceProduct.ProductName, sourceProduct.ProductVersion, upgrade.SMOOTH_UPGRADE_MODE); err != nil {
					log.Errorf("%v", err)
				}
				if err := model.DeployMysqlIpList.Delete("", sourceProduct.ProductName, clusterId); err != nil {
					log.Errorf("%v", err)
				}
			}
		case 2:
			query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET status=?, product_parsed=?, update_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, status, productParsed, pid, clusterId); err != nil {
				log.Errorf("%v", err)
			}
			query = "UPDATE " + model.DeployClusterSmoothUpgradeProductRel.TableName + " SET is_deleted=1, update_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, sourceProduct.ID, clusterId); err != nil {
				log.Errorf("%v", err)
			}
			query = "UPDATE " + upgrade.UpgradeHistory.TableName + " SET upgrade_mode=? WHERE cluster_id=? AND product_name=? AND target_version=? AND upgrade_mode=? AND is_deleted=0"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, "", clusterId, sourceProduct.ProductName, sourceProduct.ProductVersion, upgrade.SMOOTH_UPGRADE_MODE); err != nil {
				log.Errorf("%v", err)
			}
			if err := model.DeployMysqlIpList.Delete("", sourceProduct.ProductName, clusterId); err != nil {
				log.Errorf("%v", err)
			}
		default:
			query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET status=?, product_parsed=?, update_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, status, productParsed, pid, clusterId); err != nil {
				log.Errorf("%v", err)
			}
			productInfo, err := model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterId(sc.ProductName, clusterId)
			if err != nil {
				log.Errorf("%v", err)
			}
			query = "UPDATE " + model.DeployClusterSmoothUpgradeProductRel.TableName + " SET is_deleted=1, update_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, productInfo.ID, clusterId); err != nil {
				log.Errorf("%v", err)
			}
			query = "UPDATE " + upgrade.UpgradeHistory.TableName + " SET upgrade_mode=? WHERE cluster_id=? AND product_name=? AND upgrade_mode=? AND is_deleted=0"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, "", clusterId, sc.ProductName, upgrade.SMOOTH_UPGRADE_MODE); err != nil {
				log.Errorf("%v", err)
			}
		}
		query = "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? AND cluster_id=?"
		if _, err := model.DeployProductHistory.GetDB().Exec(query, status, deployUUID, clusterId); err != nil {
			log.Errorf("%v", err)
		}

		contextCancelMapMutex.Lock()
		delete(contextCancelMap, deployUUID)
		contextCancelMapMutex.Unlock()
	}()

	log.Infof("cluster %v installing new instance and rolling update ...", clusterId)

	if err = newRollingUpdate(sc, deployUUID, pid, ctx, uncheckedServices, clusterId, installMode, operationId, finalUpgrade).run(); err != nil {
		log.Errorf("%v update error: %v", deployUUID, err)
		return
	}

	log.Infof("cluster %v uninstalling old instance ...", clusterId)

	var serviceIpInfo []model.DeployServiceIpInfo
	query = "SELECT * FROM " + model.DeployServiceIpList.TableName + " WHERE product_name=? AND cluster_id=?"
	if err = model.USE_MYSQL_DB().Select(&serviceIpInfo, query, sc.ProductName, clusterId); err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Debugf("cluster %v found %v service ip_list: %+v", clusterId, sc.ProductName, serviceIpInfo)

	var instanceInfo []model.InstanceAndProductInfo
	query = fmt.Sprintf("SELECT %s.*, product_name, product_name_display, product_version FROM %s LEFT JOIN %s ON pid=%s.id WHERE product_name=? AND cluster_id=?",
		model.DeployInstanceList.TableName,
		model.DeployInstanceList.TableName,
		model.DeployProductList.TableName,
		model.DeployProductList.TableName,
	)
	if err = model.USE_MYSQL_DB().Select(&instanceInfo, query, sc.ProductName, clusterId); err != nil {
		log.Errorf("%v", err)
		return
	}

	for _, info := range instanceInfo {
		sc := schema.ServiceConfig{}
		if err = json.Unmarshal(info.Schema, &sc); err != nil {
			log.Errorf("%v", err)
			return
		}
		healthCheckinfos, err := model.HealthCheck.GetInfoByClusterIdAndProductNameAndServiceName(info.ClusterId, info.ProductName, info.ServiceName, info.Ip)
		if err != nil {
			log.Errorf("err: %v", err.Error())
			return
		}
		if len(healthCheckinfos) != 0 {
			query := fmt.Sprintf("DELETE FROM %s WHERE cluster_id=? AND product_name=? AND service_name=? AND ip=?", model.HealthCheck.TableName)
			if _, err := model.HealthCheck.GetDB().Exec(query, info.ClusterId, info.ProductName, info.ServiceName, info.Ip); err != nil {
				log.Errorf("%v", err)
				return
			}
		}
		for _, v := range sc.Instance.ExtendedHealthCheck {
			query = "INSERT INTO " + model.HealthCheck.TableName + " (cluster_id, product_name, pid, service_name, agent_id, sid, ip, " +
				"script_name, script_name_display, auto_exec, period, retries, exec_status, create_time) " +
				"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?, NOW())"
			if _, err := model.HealthCheck.GetDB().Exec(query, info.ClusterId, info.ProductName, info.Pid, info.ServiceName, info.AgentId, info.Sid, info.Ip,
				v.Shell, v.Name, v.AutoExec, v.Period, v.Retries, 0); err != nil {
				log.Errorf("%v", err)
				return
			}
		}
	}
	//非平滑升级，卸载原服务
	if installMode != 3 {
		var oldInstance []model.InstanceAndProductInfo
		for _, info := range instanceInfo {
			var bReserve bool
			if srvConfig, exist := sc.Service[info.ServiceName]; info.ProductVersion == sc.ProductVersion && exist {
				if srvConfig.Instance != nil && !srvConfig.Instance.UseCloud && !util.StringContain(uncheckedServices, info.ServiceName) {
					for _, ipInfo := range serviceIpInfo {
						// strictly match ip
						isIncludeIp := false
						for _, v := range strings.Split(ipInfo.IpList, ",") {
							if v == info.Ip {
								isIncludeIp = true
								break
							}
						}
						if ipInfo.ServiceName == info.ServiceName && isIncludeIp {
							log.Debugf("instance reserve (ProductName:%v, ProductVersion:%v, ServiceName: %v, Ip: %v)",
								sc.ProductName, info.ProductVersion, info.ServiceName, info.Ip)
							bReserve = true
							break
						}
					}
				}
			}
			if !bReserve {
				oldInstance = append(oldInstance, info)
			}
		}
		if err = uninstallInstances(clusterId, oldInstance, deployUUID, operationId); err != nil {
			return
		}
	}

	log.Infof("deploy %v(%v) success", sc.ProductName, sc.ProductVersion)
}

func (r *rollingUninstall) rollingUninstallCore(clusterId int, svcName string) error {
	var oldInstance []model.InstanceAndProductInfo
	query := fmt.Sprintf("SELECT %s.*, product_name, product_name_display, product_version FROM %s LEFT JOIN %s ON pid=%s.id WHERE pid=? AND service_name=? AND cluster_id=?",
		model.DeployInstanceList.TableName,
		model.DeployInstanceList.TableName,
		model.DeployProductList.TableName,
		model.DeployProductList.TableName,
	)
	if err := model.USE_MYSQL_DB().Select(&oldInstance, query, r.pid, svcName, clusterId); err != nil {
		return err
	}
	err := uninstallInstances(r.clusterId, oldInstance, r.deployUUID, r.operationId)
	if err == nil {
		model.DeleteNodeByClusterIdProductService(clusterId, r.schema.ProductName, svcName)
	}
	return err
}

func (r *rollingUninstall) findBeDepends(name string) []string {
	list := make([]string, 0)
	for svcName, svc := range r.schema.Service {
		for _, dname := range svc.DependsOn {
			if dname == name {
				list = append(list, svcName)
			}
		}
	}
	return list
}

// rollingUninstall should be goroutine and recursive
func (r *rollingUninstall) rollingUninstall(clusterId int, svcName string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	r.svcMap[svcName].Do(func() {
		wgThis := sync.WaitGroup{}
		for _, beDname := range r.findBeDepends(svcName) {
			wgThis.Add(1)
			go r.rollingUninstall(clusterId, beDname, &wgThis)
		}
		wgThis.Wait()

		if err = r.getError(); err != nil {
			return
		}

		if r.schema.Service[svcName].Instance != nil && !r.schema.Service[svcName].Instance.UseCloud {
			log.Infof("rollingUninstallCore %v ...", svcName)

			if err = r.rollingUninstallCore(clusterId, svcName); err != nil {
				r.setError(err)
			}

			log.Infof("rollingUninstallCore %v finish(%v)", svcName, err)
		}
	})

	return
}

func (r *rollingUninstall) setError(err error) {
	r.errMu.Lock()
	defer r.errMu.Unlock()

	r.err = err
}

func (r *rollingUninstall) getError() error {
	r.errMu.Lock()
	defer r.errMu.Unlock()

	return r.err
}

func (r *rollingUninstall) run() error {
	for name := range r.schema.Service {
		r.svcMap[name] = &sync.Once{}
	}
	for name := range r.schema.Service {
		r.Add(1)
		go r.rollingUninstall(r.clusterId, name, &r.WaitGroup)
	}
	r.Wait()

	return r.getError()
}

func undeploy(clusterId int, sc *schema.SchemaConfig, deployUUID uuid.UUID, pid int) {
	var err error
	var query string

	defer func() {
		var status = model.PRODUCT_STATUS_UNDEPLOYED
		if err != nil {
			status = model.PRODUCT_STATUS_UNDEPLOY_FAIL
		}

		if status == model.PRODUCT_STATUS_UNDEPLOYED {
			query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET status=?, is_deleted=?, update_time=NOW() WHERE pid=? AND clusterId=?"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, status, 1, pid, clusterId); err != nil {
				log.Errorf("%v", err)
			}
		} else {
			query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET status=?, update_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, status, pid, clusterId); err != nil {
				log.Errorf("%v", err)
			}
		}
		query = "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? AND cluster_id=?"
		if _, err := model.DeployProductHistory.GetDB().Exec(query, status, deployUUID, clusterId); err != nil {
			log.Errorf("%v", err)
		}
	}()

	log.Infof("cluster %v undeploy instance ...", clusterId)

	if err = newRollingUninstall(clusterId, sc, deployUUID, pid).run(); err != nil {
		log.Errorf("%v uninstall error: %v", deployUUID, err)
		return
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE cluster_id=? AND pid=?", model.HealthCheck.TableName)
	if _, err := model.HealthCheck.GetDB().Exec(query, clusterId, pid); err != nil {
		log.Errorf("%v", err)
		return
	}

	log.Infof("undeploy %v(%v) success", sc.ProductName, sc.ProductVersion)
}

func uninstallInstances(clusterId int, instanceInfo []model.InstanceAndProductInfo, deployUUID uuid.UUID, operationId string) error {
	var removeRlt int64 // 0 success, 1 fail
	wg := sync.WaitGroup{}

	for _, info := range instanceInfo {
		log.Debugf("cluster %v instance remove (ProductName:%v, ProductVersion:%v, ServiceName: %v, IP: %v)",
			clusterId, info.ProductName, info.ProductVersion, info.ServiceName, info.Ip)
		var err error
		oldsc := schema.ServiceConfig{}
		if err = json.Unmarshal(info.Schema, &oldsc); err != nil {
			log.Errorf("%v", err)
			return err
		}
		err, hostInfo := model.DeployHostList.GetHostInfoByIp(info.Ip)
		if err != nil {
			return err
		}
		var instancer instance.Instancer
		if instancer, err = instance.NewInstancer(info.Pid, info.Ip, info.ServiceName, &schema.SchemaConfig{
			ProductName:    info.ProductName,
			ProductVersion: info.ProductVersion,
			Service:        map[string]schema.ServiceConfig{info.ServiceName: oldsc},
		}, operationId); err != nil {
			log.Errorf("%v", err)
			return err
		}
		defer instancer.Clear()

		wg.Add(1)
		go func(info model.InstanceAndProductInfo) (err error) {
			updateQuery := "UPDATE " + model.DeployInstanceRecord.TableName + " SET `status`=?, status_message=?, progress=?, update_time=NOW() WHERE id=?"

			defer func() {
				if err != nil {
					atomic.StoreInt64(&removeRlt, 1)
					if e, exist := err.(instanceErr); exist {
						if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, e.status, e.Error(), e.progress, e.id); err != nil {
							log.Errorf("%v", err)
						}
					}
				}
				wg.Done()
			}()

			instanceRecordInfo := model.DeployInstanceRecordInfo{
				DeployUUID:         deployUUID,
				InstanceId:         info.ID,
				Sid:                info.Sid,
				Ip:                 info.Ip,
				ProductName:        info.ProductName,
				ProductNameDisplay: info.ProductNameDisplay,
				ProductVersion:     info.ProductVersion,
				Group:              info.Group,
				ServiceName:        info.ServiceName,
				ServiceNameDisplay: info.ServiceNameDisplay,
				ServiceVersion:     info.ServiceVersion,
				Status:             model.INSTANCE_STATUS_UNINSTALLING,
				Progress:           0,
			}
			insertQuery := "INSERT INTO " + model.DeployInstanceRecord.TableName +
				" (deploy_uuid, instance_id, sid, ip, product_name, product_name_display, product_version, `group`, service_name, service_name_display, service_version, `status`, progress) VALUES" +
				" (:deploy_uuid, :instance_id, :sid, :ip, :product_name, :product_name_display, :product_version, :group, :service_name, :service_name_display, :service_version, :status, :progress)"
			rlt, err := model.DeployInstanceRecord.GetDB().NamedExec(insertQuery, &instanceRecordInfo)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
			id, _ := rlt.LastInsertId()
			if hostInfo.Status == host.SidecarOffline {
				if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLED, "", 10, id); err != nil {
					log.Errorf("%v", err)
					return instanceErr{id, model.INSTANCE_STATUS_UNINSTALL_FAIL, 80, err}
				}
				return nil
			}

			switch info.Status {
			case model.INSTANCE_STATUS_RUN_FAIL, model.INSTANCE_STATUS_RUNNING, model.INSTANCE_STATUS_STOP_FAIL:
				if err := instancer.Stop(); err != nil {
					log.Errorf("%v", err)
					return instanceErr{id, model.INSTANCE_STATUS_STOP_FAIL, 30, err}
				}
			}
			if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLING, "", 10, id); err != nil {
				log.Errorf("%v", err)
			}
			if err := instancer.UnInstall(false); err != nil {
				log.Errorf("%v", err)
				return instanceErr{id, model.INSTANCE_STATUS_UNINSTALL_FAIL, 80, err}
			}
			discover.FlushServiceDiscover()
			if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLED, "", 100, id); err != nil {
				log.Errorf("%v", err)
			}
			query := fmt.Sprintf("DELETE FROM %s WHERE cluster_id=? AND product_name=? AND service_name=? AND pid=? AND ip=?", model.HealthCheck.TableName)
			if _, err := model.HealthCheck.GetDB().Exec(query, info.ClusterId, info.ProductName, info.ServiceName, info.Pid, info.Ip); err != nil {
				log.Errorf("%v", err)
				return err
			}

			return nil
		}(info)
	}
	wg.Wait()

	if removeRlt == 1 {
		return fmt.Errorf("some instance remove fail")
	}
	return nil
}

func filterInstance(serviceIp []string, instanceInfo []model.InstanceAndProductInfo) (newInstance []string, rollingInstance, oldInstance []model.InstanceAndProductInfo) {
	for _, ip := range serviceIp {
		bFound := false
		for _, info := range instanceInfo {
			if ip == info.Ip {
				// found rolling ip(instance)
				rollingInstance = append(rollingInstance, info)
				bFound = true
				break
			}
		}
		if !bFound {
			// found new ip(instance)
			newInstance = append(newInstance, ip)
		}
	}
	for _, ii := range instanceInfo {
		bFound := false
		for _, ri := range rollingInstance {
			if ii.Ip == ri.Ip {
				bFound = true
				break
			}
		}
		if !bFound {
			// found old ip(instance)
			oldInstance = append(oldInstance, ii)
		}
	}
	return
}

func assignNodeID(clusterId int, productName, serviceName string, serviceIp, newInstance []string, oldInstance, rollingInstance []model.InstanceAndProductInfo) (map[uint]uint, error) {

	for _, info := range oldInstance {
		n, err := model.GetServiceIpNode(clusterId, productName, serviceName, info.Ip)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err := n.Delete(); err != nil {
			return nil, err
		}
	}
	for _, ip := range newInstance {
		n := model.ServiceIpNode{ClusterId: clusterId, ProductName: productName, ServiceName: serviceName, Ip: ip}
		if err := n.Delete(); err != nil {
			return nil, err
		}
	}
	nodeIdMap, err := getNodeIdMap(clusterId, productName, serviceName)
	if err != nil {
		return nil, err
	}
	//add new instance ip node
	var floorNodeId uint
	for i := 0; i < len(newInstance); i++ {
		n := model.ServiceIpNode{ClusterId: clusterId, ProductName: productName, ServiceName: serviceName, Ip: newInstance[i]}
		for n.NodeId = floorNodeId + 1; nodeIdMap[n.NodeId]; n.NodeId++ {
		}
		if err := n.Create(); err != nil {
			return nil, err
		}
		floorNodeId = n.NodeId
	}
	//补全UnInstall 失败导致service_ip_node丢失情况
	//部署产品包->服务主机编排->去掉部分ip->部署失败（导致service_ip_node 丢失）
	//部署产品包->服务主机编排->勾选上次部署去掉的ip->此处逻辑会修复丢失的service_ip_node
	nodeIdMapR, err := getNodeIdMap(clusterId, productName, serviceName)
	if err != nil {
		return nil, err
	}
	var rollNodeId uint
	for _, roll := range rollingInstance {
		_, err := model.GetServiceIpNode(clusterId, productName, serviceName, roll.Ip)
		if err == sql.ErrNoRows {
			n := model.ServiceIpNode{ClusterId: clusterId, ProductName: productName, ServiceName: serviceName, Ip: roll.Ip}
			for n.NodeId = rollNodeId + 1; nodeIdMapR[n.NodeId]; n.NodeId++ {
			}
			if err := n.Create(); err != nil {
				return nil, err
			}
			rollNodeId = n.NodeId
		}
	}
	//如果服务nodeid不是从1开始，订正最小的服务nodeid为1，确保部署不阻塞
	nodeIdMapN, err := getNodeIdMap(clusterId, productName, serviceName)
	if err != nil {
		return nil, err
	}
	//如果服务nodeid不是从1开始，订正最小的服务nodeid为1，确保依赖nodeid为1的服务部署不阻塞
	if _, ok := nodeIdMapN[1]; !ok {
		var nodeId uint
		for nodeId = 1; nodeId > 0; nodeId++ {
			if _, ok := nodeIdMapN[nodeId]; ok {
				break
			}
		}
		//log.Infof("fix nodeId to %v, cluster %v, productName %v, serviceName %v, nodeId %v")
		err := model.UpdateNodeIdWithNodeId(clusterId, productName, serviceName, nodeId, 1)
		if err != nil {
			log.Errorf("%v", err.Error())
		}
	}
	return getIdToIndex(clusterId, productName, serviceName, serviceIp)
}

func getNodeIdMap(clusterId int, productName, serviceName string) (map[uint]bool, error) {
	nodes, err := model.GetServiceNodes(clusterId, productName, serviceName)
	if err != nil {
		return nil, err
	}
	nodeIdMap := make(map[uint]bool, len(nodes))
	for _, n := range nodes {
		nodeIdMap[n.NodeId] = true
	}
	return nodeIdMap, nil
}

func getIdToIndex(clusterId int, productName, serviceName string, serviceIp []string) (map[uint]uint, error) {
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

func getHostsFromIP(ips []string) (hosts []string, err error) {
	if len(ips) == 0 {
		return
	}

	ipHostMap := make(map[string]string, len(ips))
	ipHostInfo := make([]model.HostInfo, 0)
	query := "SELECT ip, hostname FROM " + model.DeployHostList.TableName
	if err = model.DeployHostList.GetDB().Select(&ipHostInfo, query); err != nil {
		return
	}
	for _, info := range ipHostInfo {
		ipHostMap[info.Ip] = info.HostName
	}
	ipNodeInfo := make([]model.NodeInfo, 0)
	query = "SELECT ip, hostname FROM " + model.DeployNodeList.TableName
	if err = model.DeployNodeList.GetDB().Select(&ipNodeInfo, query); err != nil {
		return
	}
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

func handleUncheckedServices(sc *schema.SchemaConfig, pid, clusterId int, namespace string) error {
	var info *model.DeployUncheckedServiceInfo
	var err error

	info, err = model.DeployUncheckedService.GetUncheckedServicesByPidClusterId(pid, clusterId, namespace)
	if err != nil {
		return err
	}

	//info, err = model.DeployUncheckedService.GetUncheckedServicesByPidClusterId(pid, clusterId, "")
	//if err != nil {
	//	return err
	//}

	if info.UncheckedServices != "" {
		err = handleUncheckedServicesCore(sc, strings.Split(info.UncheckedServices, ","))
	}
	return err
}

func handleUncheckedServicesCore(sc *schema.SchemaConfig, uncheckedServices []string) error {
	var err error
	for _, name := range uncheckedServices {
		if _, ok := sc.Service[name]; !ok {
			err = errors2.Wrap(err, fmt.Errorf("unchecked service `%v` not exist", name))
			continue
		}
		if err_ := sc.SetServiceEmpty(name); err_ != nil {
			err = errors2.Wrap(err, err_)
		}
	}
	return err
}

type patchupdateParam struct {
	ParentProductName string `json:"parentProductname"`
	ProductName       string `json:"product_name"`
	Version           string `json:"version"`
	Path              string `json:"path"`
	ProductType       int    `json:"product_type"`
	PackageName       string `json:"package_name"`
	UpdateUUID        string `json:"uuid"`
}

type deployParam struct {
	UncheckedServices []string `json:"unchecked_services,omitempty"`
	ClusterId         int      `json:"clusterId"`
	Pid               int      `json:"pid"`
	Namespace         string   `json:"namespace,omitempty"`
	RelyNamespace     string   `json:"relyNamespace,omitempty"`
	DeployMode        int      `json:"deployMode,omitempty"`
	SourceVersion     string   `json:"source_version,omitempty"`
	FinalUpgrade      bool     `json:"final_upgrade"`
}

type unDeployParam struct {
	Namespace string `json:"namespace,omitempty"`
	ClusterId int    `json:"clusterId"`
}

func Deploy(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	if productName == "" || productVersion == "" {
		return fmt.Errorf("product_name or product_version is empty")
	}
	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		return fmt.Errorf("get userId err: %v", err)
	}
	param := deployParam{}
	if err = ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if param.ClusterId == 0 {
		param.ClusterId, err = GetCurrentClusterFromParam(ctx)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
	}
	if param.DeployMode == 3 {
		if err := checkSmoothUpgradeServiceAddr(param.ClusterId, productName, productVersion, param.FinalUpgrade); err != nil {
			return err
		}
	}
	log.Infof("deploy product_name:%v, product_version: %v, userId: %v, clusterId: %v", productName, productVersion, userId, param.ClusterId)
	cluster, err := model.DeployClusterList.GetClusterInfoById(param.ClusterId)
	if err != nil {
		return err
	}
	defer func() {
		if err := addSafetyAuditRecord(ctx, "部署向导", "产品部署", "集群名称："+cluster.Name+", 组件名称："+productName+productVersion); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
		if err := model.NotifyEvent.DeleteNotifyEvent(cluster.Id, 0, productName, "", "", false); err != nil {
			log.Errorf("delete notify event error: %v", err)
		}
	}()
	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_KUBERNETES {
		return DealK8SDeploy(param.Namespace, param.UncheckedServices, userId, param.ClusterId, param.RelyNamespace, param.Pid)
	} else {
		return DealDeploy(productName, productVersion, param.SourceVersion, param.UncheckedServices, userId, param.ClusterId, param.DeployMode, param.FinalUpgrade)
	}
}

func DeployForDevOps(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	if productName == "" || productVersion == "" {
		return fmt.Errorf("product_name or product_version is empty")
	}
	userId := 1
	param := deployParam{}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if param.ClusterId == 0 {
		return fmt.Errorf("clusterId empty")
	}
	log.Infof("deploy product_name:%v, product_version: %v, userId: %v, clusterId: %v", productName, productVersion, userId, param.ClusterId)
	cluster, err := model.DeployClusterList.GetClusterInfoById(param.ClusterId)
	if err != nil {
		return err
	}
	defer func() {
		if err := addSafetyAuditRecord(ctx, "部署向导", "产品部署", "集群名称："+cluster.Name+", 组件名称："+productName+productVersion); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
		if err := model.NotifyEvent.DeleteNotifyEvent(cluster.Id, 0, productName, "", "", false); err != nil {
			log.Errorf("delete notify event error: %v", err)
		}
	}()
	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_KUBERNETES {
		return DealK8SDeploy(param.Namespace, param.UncheckedServices, userId, param.ClusterId, param.RelyNamespace, param.Pid)
	} else {
		return DealDeploy(productName, productVersion, "", param.UncheckedServices, userId, param.ClusterId, param.DeployMode, false)
	}
}

func DealDeploy(productName, productVersion, sourceVersion string, uncheckedServices []string, userId, clusterId int, installMode int, finalUpgrade bool) (rlt interface{}) {
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()
	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	err := model.OperationList.InsertWithTx(tx, model.OperationInfo{
		ClusterId:       clusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.ProductDeploy.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Product.Code,
		ObjectValue:     productName,
	})
	if err != nil {
		log.Errorf("OperationList Insert err:%v", err)
	}

	var productListInfo model.DeployProductListInfo
	query := "SELECT id, product, parent_product_name, product_type, product_name FROM " + model.DeployProductList.TableName + " WHERE product_name=? AND product_version=?"
	if err := tx.Get(&productListInfo, query, productName, productVersion); err != nil {
		return err
	}

	sc, err := schema.Unmarshal(productListInfo.Product) // now product
	if err != nil {
		return err
	}
	if err = inheritBaseService(clusterId, sc, tx); err != nil {
		return err
	}
	if err = setSchemaFieldServiceAddr(clusterId, sc, tx, ""); err != nil {
		return err
	}
	if err = handleUncheckedServicesCore(sc, uncheckedServices); err != nil {
		return err
	}
	if err = sc.CheckServiceAddr(); err != nil {
		log.Errorf("%v", err)
		return err
	}
	if installMode == 3 {
		if err = model.DeployClusterSmoothUpgradeProductRel.CheckProductReadyForDeploy(productName); err != nil {
			return err
		}
	} else {
		if err = model.DeployClusterProductRel.CheckProductReadyForDeploy(productName); err != nil {
			return err
		}
	}

	deployUUID := uuid.NewV4()
	rel := model.ClusterProductRel{
		Namespace:     "",
		ProductParsed: []byte(""),
		Pid:           productListInfo.ID,
		ClusterId:     clusterId,
		Status:        model.PRODUCT_STATUS_DEPLOYING,
		DeployUUID:    deployUUID.String(),
		UserId:        userId,
	}

	var productRelTable string
	var oldProductListInfo *model.DeployProductListInfo
	var _err error
	if installMode == 3 {
		productRelTable = model.DeployClusterSmoothUpgradeProductRel.TableName
		oldProductListInfo, _err = model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	} else {
		productRelTable = model.DeployClusterProductRel.TableName
		oldProductListInfo, _err = model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	}
	//升级或者重新部署
	if _err == nil {
		query = "UPDATE " + productRelTable + " SET pid=?, user_id=?, `status`=?, `deploy_uuid`=?, deploy_time=NOW() WHERE pid=? AND clusterId=? AND is_deleted=0"
		if _, err := tx.Exec(query, productListInfo.ID, userId, model.PRODUCT_STATUS_DEPLOYING, deployUUID, oldProductListInfo.ID, clusterId); err != nil {
			log.Errorf("%v", err)
			return err
		}
		//如果产品包id变更了
		if oldProductListInfo.ID != productListInfo.ID {
			currSelect, err := model.ProductSelectHistory.GetPidListStrByClusterId(clusterId)
			if err != nil {
				log.Errorf("model.ProductSelectHistory.GetPidListStrByClusterId err %v", err)
			} else {
				//移除老的 添加新的
				if strings.TrimSpace(currSelect) != "" {
					oldPids := strings.Split(currSelect, ",")
					for idx, pidStr := range oldPids {
						//如果升级的包存在老的 select 中  那么需要更新该 id
						if strconv.Itoa(oldProductListInfo.ID) == pidStr {
							oldPids = append(oldPids[:idx], oldPids[idx+1:]...)
							break
						}
					}
					newPids := append(oldPids, strconv.Itoa(productListInfo.ID))
					err := model.ProductSelectHistory.SetPidListStrByClusterId(strings.Join(newPids, ","), clusterId)
					if err != nil {
						log.Errorf("model.ProductSelectHistory.SetPidListStrByClusterId err %v", err)
					}
				}
			}
		}

		//	安装
	} else if _err == sql.ErrNoRows {
		query = "INSERT INTO " + productRelTable + " (namespace,product_parsed,pid, clusterId, deploy_uuid, user_id, deploy_time, status) VALUES" +
			" (:namespace,:product_parsed,:pid, :clusterId, :deploy_uuid, :user_id, NOW(), :status)"
		if _, err = tx.NamedExec(query, &rel); err != nil {
			log.Errorf("%v", err)
			return err
		}
		err = model.ProductSelectHistory.AddPidByClusterId(strconv.Itoa(productListInfo.ID), clusterId)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
	} else {
		log.Errorf("%v", err)
		return err
	}

	// 新一次部署之前删除之前产品包的unchecked_service
	products, _ := model.DeployProductList.GetProductListByNameAndType(productListInfo.ProductName, strconv.Itoa(productListInfo.ProductType), nil)
	for _, product := range products {
		query = "DELETE FROM " + model.DeployUncheckedService.TableName + " WHERE pid=? AND cluster_id=?"
		if _, err = tx.Exec(query, product.ID, clusterId); err != nil && err != sql.ErrNoRows {
			log.Errorf("%v", err)
			return err
		}
	}

	if len(uncheckedServices) > 0 {
		uncheckedServiceInfo := model.DeployUncheckedServiceInfo{ClusterId: clusterId, Pid: productListInfo.ID, UncheckedServices: strings.Join(uncheckedServices, ","), Namespace: ""}
		query = "INSERT INTO " + model.DeployUncheckedService.TableName + " (pid, cluster_id, unchecked_services,namespace) VALUES" +
			" (:pid, :cluster_id, :unchecked_services,:namespace) ON DUPLICATE KEY UPDATE unchecked_services=:unchecked_services, update_time=NOW()"
		if _, err = tx.NamedExec(query, &uncheckedServiceInfo); err != nil {
			log.Errorf("%v", err)
			return err
		}
	} else {
		query = "DELETE FROM " + model.DeployUncheckedService.TableName + " WHERE pid=? AND cluster_id=?"
		if _, err = tx.Exec(query, productListInfo.ID, clusterId); err != nil && err != sql.ErrNoRows {
			log.Errorf("%v", err)
			return err
		}
	}

	productHistoryInfo := model.DeployProductHistoryInfo{
		Namespace:          "",
		ClusterId:          clusterId,
		DeployUUID:         deployUUID,
		ProductName:        productName,
		ProductNameDisplay: productListInfo.ProductNameDisplay,
		ProductVersion:     productVersion,
		Status:             model.PRODUCT_STATUS_DEPLOYING,
		ParentProductName:  productListInfo.ParentProductName,
		UserId:             userId,
	}
	sc.ParentProductName = productListInfo.ParentProductName

	query = "INSERT INTO " + model.DeployProductHistory.TableName + " (namespace,cluster_id, product_name, product_name_display, deploy_uuid, product_version, `status`, parent_product_name, deploy_start_time, user_id) " +
		"VALUES (:namespace,:cluster_id, :product_name, :product_name_display, :deploy_uuid, :product_version, :status , :parent_product_name, NOW(), :user_id)"
	if _, err := tx.NamedExec(query, &productHistoryInfo); err != nil {
		log.Errorf("%v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("%v", err)
		return err
	}

	ctx, cancel := sysContext.WithCancel(sysContext.Background())
	contextCancelMapMutex.Lock()
	contextCancelMap[deployUUID] = cancel
	contextCancelMapMutex.Unlock()
	//所有的list 接口用到接收的 uuid 参数 都要在该表中有记录 用以判断该 uuid 类型
	err = model.DeployUUID.InsertOne(deployUUID.String(), "", model.ManualDeployUUIDType, productListInfo.ID)
	if err != nil {
		log.Errorf("%v", err)
		return nil
	}
	go deploy(sc, deployUUID, productListInfo.ID, ctx, uncheckedServices, clusterId, installMode, operationId, sourceVersion, finalUpgrade)

	return map[string]interface{}{"deploy_uuid": deployUUID}
}

func Undeploy(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	if productName == "" || productVersion == "" {
		return fmt.Errorf("product_name or product_version is empty")
	}
	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		return fmt.Errorf("get userId err: %v", err)
	}
	param := unDeployParam{}
	if err = ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	var clusterId int
	if param.ClusterId > 0 {
		clusterId = param.ClusterId
	} else {
		clusterId, err = GetCurrentClusterFromParam(ctx)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
	}

	log.Infof("cluster %v undeploy product_name:%v, product_version: %v, userId: %v", clusterId, productName, productVersion, userId)

	defer func() {
		info, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "部署向导", "产品卸载", "集群名称："+info.Name+", 组件名称："+productName+productVersion); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	return DealUndeploy(productName, productVersion, param.Namespace, userId, clusterId)
}

func DealUndeploy(productName, productVersion, namespace string, userId, clusterId int) (rlt interface{}) {

	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()
	// get info about product that will be undeploy(from deploy_product_list table)
	var productListInfo model.DeployProductListInfo
	query := "SELECT id, product, parent_product_name FROM " + model.DeployProductList.TableName + " WHERE product_name=? AND product_version=?"
	if err := tx.Get(&productListInfo, query, productName, productVersion); err != nil {
		return err
	}

	// return if product is deploying
	for {
		kubeLock, _ := model.DeployKubeProductLock.GetByPidAndClusterIdAndNamespace(productListInfo.ID, clusterId, namespace)
		if kubeLock.IsDeploy == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	sc, err := schema.Unmarshal(productListInfo.Product)
	if err != nil {
		return err
	}
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(productName)
	if err != nil {
		return err
	}
	// 获取指定集群下要卸载的产品包详细信息，如:(schema、产品包名称、版本等)
	var installedProductListInfo *model.DeployProductListInfo
	installedProductListInfo, err = model.DeployClusterProductRel.GetCurrentProductByProductNameClusterIdNamespace(productName, clusterId, namespace)
	if err != nil {
		return err
	}

	if installedProductListInfo.ID != productListInfo.ID {
		return fmt.Errorf("no target pid: %v", productListInfo.ID)
	}
	// 获取指定集群下产品包部署的关联记录
	currentProductRel, err := model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(productListInfo.ID, clusterId, namespace)
	if err != nil {
		return err
	}

	if currentProductRel.Status == model.PRODUCT_STATUS_UNDEPLOYING {
		return fmt.Errorf("product is undeploying, can't undeploy again")
	}

	// update if product is deploying
	if currentProductRel.Status == model.PRODUCT_STATUS_DEPLOYING {
		query = "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? AND cluster_id=?"
		if _, err := model.DeployProductHistory.GetDB().Exec(query, model.PRODUCT_STATUS_DEPLOY_FAIL, currentProductRel.DeployUUID, clusterId); err != nil {
			return err
		}
	}

	deployUUID := uuid.NewV4()
	query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET pid=?, user_id=?, `status`=?, `deploy_uuid`=?, deploy_time=NOW() WHERE clusterId=? AND pid=? AND is_deleted=0 AND namespace=?"
	if _, err := tx.Exec(query, productListInfo.ID, userId, model.PRODUCT_STATUS_UNDEPLOYING, deployUUID, clusterId, installedProductListInfo.ID, namespace); err != nil {
		log.Errorf("%v", err)
		return err
	}
	//所有的list 接口用到接收的 uuid 参数 都要在该表中有记录 用以判断该 uuid 类型
	err = model.DeployUUID.InsertOne(deployUUID.String(), "", model.ManualDeployUUIDType, productListInfo.ID)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	// insert into undeploy history
	productHistoryInfo := model.DeployProductHistoryInfo{
		ClusterId:          clusterId,
		Namespace:          namespace,
		DeployUUID:         deployUUID,
		ProductName:        productName,
		ProductNameDisplay: productListInfo.ProductNameDisplay,
		ProductVersion:     productVersion,
		ProductType:        installedProductListInfo.ProductType,
		Status:             model.PRODUCT_STATUS_UNDEPLOYING,
		ParentProductName:  productListInfo.ParentProductName,
		UserId:             userId,
	}
	sc.ParentProductName = productListInfo.ParentProductName

	query = "INSERT INTO " + model.DeployProductHistory.TableName + " (cluster_id, namespace, product_name, product_name_display, deploy_uuid, product_version, `status`, parent_product_name, deploy_start_time, user_id, product_type) " +
		"VALUES (:cluster_id, :namespace, :product_name, :product_name_display, :deploy_uuid, :product_version, :status , :parent_product_name, NOW(), :user_id, :product_type)"
	if _, err := tx.NamedExec(query, &productHistoryInfo); err != nil {
		log.Errorf("%v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("%v", err)
		return err
	}

	cluster, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	if err != nil {
		return err
	}
	if cluster.Type == model.DEPLOY_CLUSTER_TYPE_KUBERNETES {
		go K8sUndeploy(clusterId, sc, deployUUID, productListInfo.ID, namespace)
	} else {
		go undeploy(clusterId, sc, deployUUID, productListInfo.ID)
	}

	//维护deploy_product_select_history表  卸载的时候维护  例如自动部署了 a b c  手动卸载了 b  回显应该回显 a b
	err = model.ProductSelectHistory.RemovePidByClusterId(strconv.Itoa(productListInfo.ID), clusterId)
	if err != nil {
		log.Errorf("ProductSelectHistory RemovePidByClusterId err %v", err)
	}
	return map[string]interface{}{"deploy_uuid": deployUUID}
}

func buildDeployEvents(events []model.DeployInstanceEventInfo) string {
	var serviceEvent []string
	for index, e := range events {
		var message string
		event := &agent.InstanceEvent{}
		err := json.Unmarshal([]byte(e.Content), event)
		if err != nil {
			message = err.Error()
		}
		message = event.GetMessage()
		line := fmt.Sprintf(EVENT_LINE_SEPRATOR, index+1)
		url := fmt.Sprintf(EVENT_CONTENT_URL, e.InstanceId, e.ID)
		content := fmt.Sprintf(`%s
| 组件信息     | %v
| 事件类型     | %v
| 事件时间     | %v
| 事件结果     | %v
| 事件详情     | %v
`,
			line,
			e.InstanceId,
			e.EventType,
			e.CreateDate,
			message,
			url,
		)
		serviceEvent = append(serviceEvent, string(content[:]))
	}
	return strings.Join(serviceEvent, "\r\n")
}

func buildDeplogLogs(info []model.DeployInstanceRecordByDeployIdInfo) string {
	ret := ""
	for index, s := range info {
		r := map[string]interface{}{}
		if s.Schema.Valid {
			sc := &schema.ServiceConfig{}
			json.Unmarshal([]byte(s.Schema.String), sc)
			r["schema"], _ = json.MarshalIndent(sc, "", "\t")
		} else {
			r["schema"] = "[]"
		}
		r["product_name"] = s.ProductName
		r["product_version"] = s.ProductVersion
		r["group"] = s.Group
		r["service_name"] = s.ServiceName
		r["service_version"] = s.ServiceVersion
		r["status"] = s.Status
		r["status_message"] = s.StatusMessage
		r["progress"] = s.Progress

		if s.UpdateDate.Valid == true {
			r["update_time"] = s.UpdateDate.Time.Format(base.TsLayout)
		} else {
			r["update_time"] = ""
		}
		if s.CreateDate.Valid == true {
			r["create_time"] = s.CreateDate.Time.Format(base.TsLayout)
		} else {
			r["create_time"] = ""
		}
		events, err := model.DeployInstanceEvent.GetDeployInstanceEventById(nil, int64(s.InstanceId), "")
		if err != nil {
			r["service_event"] = ""
			log.Errorf("%v", err)
		} else {
			r["service_event"] = buildDeployEvents(events)
		}
		line := fmt.Sprintf(LOG_LINE_SEPRATOR, index+1, s.ServiceName)
		content := fmt.Sprintf(`%s
部署服务名称: %s
部署服务版本: %s
部署开始时间: %s
部署结束时间: %s
部署结果: %s
部署摘要: %s
部署配置: 
%s
部署事件: 
%s
`,
			line,
			r["service_name"],
			r["service_version"],
			r["create_time"],
			r["update_time"],
			r["status"],
			r["status_message"],
			r["schema"],
			r["service_event"],
		)
		ret = ret + content
	}
	return ret
}

func DeployLogs(ctx context.Context) apibase.Result {
	log.Debugf("[Product->DeployLogs] return  product deploy logs from EasyMatrix API ")
	paramErrs := apibase.NewApiParameterErrors()

	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	if productName == "" || productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name or product_version is empty"))
	}
	deployUUID := ctx.URLParam("deploy_uuid")
	if deployUUID == "" {
		paramErrs.AppendError("$", fmt.Errorf("deploy_uuid is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	serviceName := ctx.URLParam("service")

	var status []string
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.DeployInstanceRecordInfo{})
	info, _, _ := model.DeployInstanceRecord.GetDeployInstanceRecordByDeployId(pagination, deployUUID, status, serviceName)
	logs := buildDeplogLogs(info)

	return map[string]interface{}{"result": logs}
}

func Cancel(ctx context.Context) apibase.Result {
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	if productName == "" || productVersion == "" {
		return fmt.Errorf("product_name or product_version or clusterId is empty")
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	productType := ctx.FormValue("product_type")
	if productType == "" {
		return fmt.Errorf("product_type is empty")
	}
	deployModeStr := ctx.FormValue("deploy_mode")
	deployMode, err := strconv.Atoi(deployModeStr)
	if err != nil {
		return fmt.Errorf("deploy_mode is not number")
	}
	//pagination := apibase.GetPaginationFromQueryParameters(nil, ctx)
	var _err error
	var product = &model.DeployProductListInfo{}
	if deployMode == 3 {
		//平滑升级
		product, _err = model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	} else {
		//非平滑升级
		product, _err = model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	}
	if _err != nil || product.Status != model.PRODUCT_STATUS_DEPLOYING {
		return errors.New("product is not deploying")
	}
	deployUUID, _ := uuid.FromString(product.DeployUUID)
	contextCancelMapMutex.Lock()
	if cancel, ok := contextCancelMap[deployUUID]; ok {
		cancel()
	}
	contextCancelMapMutex.Unlock()

	instances, err := model.DeployInstanceList.GetInstanceListByPidServiceName(product.ID, clusterId, "")
	if err != nil {
		return err
	}

	//params := agent.CancelParams{
	//	Agents: make(map[string][]string),
	//}
	for _, ins := range instances {
		//params.Agents[ins.Sid] = append(params.Agents[ins.Sid], ins.AgentId)

		// cancel health check
		ev := &event.Event{
			AgentId: ins.AgentId,
			Type:    event.REPORT_EVENT_HEALTH_CHECK_CANCEL,
			Data:    nil,
		}
		event.GetEventManager().EventReciever(ev)
	}

	// 往sidecar发送取消安装功能先不启用。因为sidecar的controller在install和unisntall过程中并没有将agentId放入agents map，
	// 导致即使发送rpc取消指令，因找不到agent实例而无效。其次，旧版sidecar在收到不支持指令时并不会发report或event消息，
	// 会引起响应超时，影响功能体验或增加性能负担。
	// err, respBody := agent.AgentClient.AgentCancel(&params)
	// if err != nil {
	//	log.Errorf("[Cancel] cancel %v error: %v", params)
	//	return err
	// }
	// if respBody.Data == nil {
	// 	log.Errorf("[Cancel] cancel %v error, data nil: %v", params)
	// 	return fmt.Errorf("empty response")
	// }
	// return respBody.Data

	return nil
}

// GetIP
// @Description  	Get Service Host
// @Summary      	查看服务的主机
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Success         200         {string}  string "{"msg":"ok","code":0,"data":{"ip":[],"count":""}}"
// @Router          /api/v2/product/{product_name}/service/{service_name}/get_ip [get]
func GetIP(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()

	productName := ctx.Params().Get("product_name")
	serviceName := ctx.Params().Get("service_name")

	var ipList []string

	if productName == "" || serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name or service_name is empty"))
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	err, info := model.DeployServiceIpList.GetServiceIpListByName(productName, serviceName, clusterId, "")
	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"ip":    ipList,
			"count": len(ipList),
		}
	}
	if err != nil {
		log.Errorf("[GetIP] GetServiceIpListByName err: %v", err)
		return err
	}
	ipList = strings.Split(info.IpList, IP_LIST_SEP)
	if len(ipList) == 0 {
		log.Errorf("[GetIP] service ip is null")
		return fmt.Errorf("[GetIP] service %v ip is null", serviceName)
	}
	return map[string]interface{}{
		"ip":    ipList,
		"count": len(ipList),
	}
}

func checkSortIpList(ips string) (ipList []string, err error) {
	ipMap := map[string]struct{}{}
	for _, ip := range strings.Split(ips, IP_LIST_SEP) {
		//if net.ParseIP(ip) == nil {
		//	return nil, fmt.Errorf("ip format invalid")
		//}
		if _, exist := ipMap[ip]; !exist {
			ipMap[ip] = struct{}{}
			ipList = append(ipList, ip)
		} else {
			return nil, fmt.Errorf("ip duplicate")
		}
	}
	sort.Strings(ipList)

	return
}

// SetIP
// @Description  	Set Service Host
// @Summary      	查看服务的主机
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           namespace      query     string     false  "命名空间"
// @Success         200         {string}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/{product_name}/service/{service_name}/set_ip [post]
func SetIP(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()

	productName := ctx.Params().Get("product_name")
	serviceName := ctx.Params().Get("service_name")
	//baseClusterId, _ := strconv.Atoi(ctx.FormValue("baseClusterId"))
	namespace := ctx.FormValue("namespace")

	if productName == "" || serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name or service_name is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	//不从cookie中获取集群id，从参数获取
	clusterId, err := GetCurrentClusterFromParam(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	//if baseClusterId > 0 {
	//	clusterId = baseClusterId
	//}
	//info, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(productName, clusterId)
	//if err == nil {
	//	sc, err := schema.Unmarshal(info.Schema)
	//	if err != nil {
	//		log.Errorf("%v", err)
	//		return err
	//	}
	//	baseClusterId, _ := strconv.Atoi(ctx.FormValue("baseClusterId"))
	//	if baseClusterId > 0 {
	//		clusterId = baseClusterId
	//		for _, config := range sc.Service {
	//			if config.BaseService == serviceName {
	//				productName = config.BaseProduct
	//				break
	//			}
	//		}
	//	}
	//}

	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(productName)
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}
	var query string
	//选择 ip 的时候 维护deploy_product_select_history 表
	if ctx.FormValue("ip") == "" {
		// delete ip
		query = "DELETE FROM " + model.DeployServiceIpList.TableName + " WHERE product_name=? AND service_name=? AND cluster_id=? AND namespace=?"
		if _, err := model.DeployServiceIpList.GetDB().Exec(query, productName, serviceName, clusterId, namespace); err != nil {
			return err
		}
		// 检查deploy_schema_multi_field表是否存在记录，若存在则删除所有对应服务的参数修改配置
		if err := model.SchemaMultiField.DeleteByProductNameAndServiceNameAndClusterId(productName, serviceName, clusterId, tx); err != nil {
			return err
		}
	} else {
		// 获取当前服务中编排的所有ip
		oldIpList, err := model.DeployServiceIpList.GetServiceIpList(clusterId, productName, serviceName)
		if err != nil {
			return err
		}
		// 检测ip是否重复，同时排列下ip
		ipList, err := checkSortIpList(ctx.FormValue("ip"))
		if err != nil {
			return err
		}
		for _, ip := range oldIpList {
			if !contains(ipList, ip) {
				// 该ip被移除
				if err := model.SchemaMultiField.DeleteByIp(clusterId, ip); err != nil {
					return err
				}
			}
		}
		for _, ip := range ipList {
			if !contains(oldIpList, ip) {
				// 新编排主机
				multiFields, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndClusterId(productName, serviceName, clusterId)
				if err != nil {
					log.Errorf("Get multi field when set ip error, product:%s, service:%s, err:%v", productName, serviceName, err)
					return err
				}
				var configMap = make(map[string]struct{}, 0)
				for _, multiField := range multiFields {
					if _, ok := configMap[multiField.FieldPath]; !ok {
						err := addToHead(clusterId, productName, serviceName, multiField.FieldPath, ip)
						if err != nil {
							log.Errorf("add to head error: %v", err)
						}
						configMap[multiField.FieldPath] = struct{}{}
					}
				}
			}
		}
		ValidateSchemaFields(clusterId, productName, serviceName)
		// 更新或增加服务组件和host的关联关系
		if err = model.DeployServiceIpList.SetServiceIp(productName, serviceName, strings.Join(ipList, IP_LIST_SEP), clusterId, namespace); err != nil {
			log.Errorf("[SetIP] SetServiceIp err: %v", err)
			return err
		}

		// 不为空的的时候去添加
		//err = model.ProductSelectHistory.AddPidByClusterId(strconv.Itoa(info.ID), clusterId)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

// 对于存在多配置的配置项，将新编排的主机关联到到第一个配置
func addToHead(clusterId int, productName, serviceName, fieldPath, ip string) error {
	tx := model.USE_MYSQL_DB().MustBegin()
	multiList, err := model.SchemaMultiField.GetListByFieldPath(clusterId, productName, serviceName, fieldPath)
	if err != nil {
		log.Errorf("get multi list failed error: %v", err)
		return err
	}
	// 构建添加的第一条记录
	if len(multiList) > 0 {
		oldFirstElement := multiList[0]
		newFirstElement := model.SchemaMultiFieldInfo{
			ClusterId:   oldFirstElement.ClusterId,
			ProductName: oldFirstElement.ProductName,
			ServiceName: oldFirstElement.ServiceName,
			FieldPath:   oldFirstElement.FieldPath,
			Field:       oldFirstElement.Field,
			Hosts:       ip,
		}
		multiList = append([]model.SchemaMultiFieldInfo{newFirstElement}, multiList...)
	}
	// 删除当前存在的所有该配置项的配置
	if err := model.SchemaMultiField.DeleteByFieldPath(clusterId, productName, serviceName, fieldPath); err != nil {
		return err
	}
	for _, multiField := range multiList {
		_ = modifyMultiField(&multiField, tx)
	}
	//if err := tx.Commit(); err != nil {
	//	log.Errorf("%v", err)
	//	return err
	//}
	return tx.Commit()
}

func ValidateSchemaFields(clusterId int, productName, serviceName string) error {
	fieldPathCountList := model.SchemaMultiField.GetDistinctPathCount(clusterId, productName, serviceName)
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, fieldPathCount := range fieldPathCountList {
		if fieldPathCount.Count == 1 {
			err := model.SchemaMultiField.DeleteByFieldPath(clusterId, productName, serviceName, fieldPathCount.FieldPath)
			if err != nil {
				log.Errorf("%v", err)
			}
			schemaModifyField := model.SchemaFieldModifyInfo{
				ClutserId:   clusterId,
				ProductName: productName,
				ServiceName: serviceName,
				FieldPath:   fieldPathCount.FieldPath,
				Field:       fieldPathCount.Field,
			}
			_ = modifyField(&schemaModifyField, tx)
		} else if fieldPathCount.Count > 1 {
			fieldCountList := model.SchemaMultiField.GetDistinctValueCount(clusterId, productName, serviceName)
			if len(fieldCountList) == 1 {
				fieldCount := fieldCountList[0]
				err := model.SchemaMultiField.DeleteByFieldPath(clusterId, productName, serviceName, fieldPathCount.FieldPath)
				if err != nil {
					log.Errorf("%v", err)
				}
				schemaModifyInfo := model.SchemaFieldModifyInfo{
					ClutserId:   clusterId,
					ProductName: productName,
					ServiceName: serviceName,
					FieldPath:   fieldCount.FieldPath,
					Field:       fieldCount.Field,
				}
				_ = modifyField(&schemaModifyInfo, tx)
			}
		}
	}
	return tx.Commit()
}

func ModifyAll(ctx context.Context) (rlt apibase.Result) {
	info := model.SchemaFieldModifyInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if info.ProductName == "" || info.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}
	paramErrs := apibase.NewApiParameterErrors()
	params := make(map[string]string)
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	var err error
	// 从cookie中获取集群id
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	info.ClutserId = clusterId

	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	if namespace != "" {
		info.Namespace = namespace
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(info.ProductName)
	if err != nil {
		return err
	}
	reg := regexp.MustCompile(`(?i).*password.*`)

	err, userInfo := model.UserList.GetInfoByUserId(1)
	if err != nil {
		log.Errorf("GetInfoByUserId %v", err)
		return err
	}
	// 接收修改的参数对deploy_schema_field_modify表进行修改
	for k, v := range params {
		if reg.Match([]byte(k)) {
			password, err := encrypt.PlatformEncrypt.SchemeDecrypt(v, userInfo.PassWord)
			if err != nil {
				return err
			}
			v = password
		}

		fieldpath := strings.Replace(k, "current", "Value", 1)
		info.FieldPath = fieldpath
		info.Field = v

		rlt = modifyField(&info, tx)
		if _, ok := rlt.(error); ok {
			return rlt
		}
		CheckMultiFieldConfig(clusterId, info.ProductName, info.ServiceName, fieldpath)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	// 添加集群审计记录
	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		paramsContent := ""
		for k, v := range params {
			paramsContent += k + ":" + v
		}
		if err := addSafetyAuditRecord(ctx, "集群运维", "服务参数修改", "集群名称："+clusterInfo.Name+", 组件名称："+info.ProductName+
			", 服务名称："+info.ServiceName+", 部署参数："+paramsContent); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	return
}

func ModifyMultiField(ctx context.Context) (rlt apibase.Result) {
	multiFieldInfo := model.SchemaMultiFieldInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if multiFieldInfo.ProductName == "" || multiFieldInfo.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}

	clusterId, err := GetCurrentClusterFromParam(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	multiFieldInfo.ClusterId = clusterId
	fieldPath := ctx.FormValue("field_path")
	fieldValue := ctx.FormValue("field")
	var fieldConfigSlice []model.FieldConfig
	if err = json.Unmarshal([]byte(fieldValue), &fieldConfigSlice); err != nil {
		log.Errorf("%v", err)
		return err
	}
	ipList, err := model.DeployServiceIpList.GetServiceIpList(clusterId, multiFieldInfo.ProductName, multiFieldInfo.ServiceName)
	if err != nil {
		log.Errorf("get service ip list error:%v", err)
	}
	// validate config value
	_, currentList := removeDuplicate(fieldConfigSlice)
	//if len(values) != len(fieldConfigSlice) {
	//	fieldArr := strings.Split(fieldPath, ".")
	//	return errors.New(fmt.Sprintf("(%s)输入值重复", fieldArr[1]))
	//}
	if len(ipList) != len(currentList) {
		return errors.New("存在未关联主机")
	}
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	CheckMultiFieldConfig(clusterId, multiFieldInfo.ProductName, multiFieldInfo.ServiceName, fieldPath)

	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(multiFieldInfo.ProductName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	multiFieldInfo.FieldPath = fieldPath

	for _, item := range fieldConfigSlice {
		multiFieldInfo.Field = item.Field
		hostList := strings.Split(item.Hosts, ",")
		for _, hostIp := range hostList {
			multiFieldInfo.Hosts = hostIp
			currentInfo := *&multiFieldInfo
			rlt = modifyMultiField(&currentInfo, tx)
			if _, ok := rlt.(error); ok {
				return rlt
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		if err = addSafetyAuditRecord(ctx, "集群运维", "服务参数修改", "集群名称: "+clusterInfo.Name+", 组件名称: "+
			multiFieldInfo.ProductName+", 服务名称: "+multiFieldInfo.ServiceName+", 部署参数: "+fieldPath+":"+strings.Replace(fieldValue, " ", "", -1)); err != nil {
			log.Errorf("add safety audit record error: %v", err)
			return
		}
		ValidateSchemaFields(clusterId, multiFieldInfo.ProductName, multiFieldInfo.ServiceName)

	}()

	return
}

func ModifyAllSchemaMultiField(ctx context.Context) (rlt apibase.Result) {
	// 参数校验
	multiFieldInfo := model.SchemaMultiFieldInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if multiFieldInfo.ProductName == "" || multiFieldInfo.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}

	paramErrs := apibase.NewApiParameterErrors()
	params := make(map[string][]model.FieldConfig)
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
	}
	multiFieldInfo.ClusterId = clusterId
	ipList, err := model.DeployServiceIpList.GetServiceIpList(clusterId, multiFieldInfo.ProductName, multiFieldInfo.ServiceName)
	if err != nil {
		log.Errorf("get service ip list error: %v", err)
	}
	var invalidParams []string
	// validate config value
	for k, v := range params {
		result, currentList := removeDuplicate(v)
		if len(result) != len(v) {
			fieldArr := strings.Split(k, ".")
			invalidParams = append(invalidParams, fieldArr[1])
		}
		if len(ipList) != len(currentList) {
			return errors.New("存在未关联主机")
		}
	}
	if len(invalidParams) > 0 {
		return errors.New(fmt.Sprintf("(%s)输入值重复", strings.Join(invalidParams, ",")))
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()

	// 检测产品包是否正在部署或是卸载
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(multiFieldInfo.ProductName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for k, v := range params {
		fieldPath := strings.Replace(k, "current", "Value", 1)
		multiFieldInfo.FieldPath = fieldPath
		// 判断数据的记录是否需要修改
		if model.SchemaMultiField.WhetherChangeField(clusterId, multiFieldInfo.ProductName, multiFieldInfo.ServiceName, fieldPath, v) {
			log.Infof("service multi field no change")
			continue
		}
		CheckMultiFieldConfig(clusterId, multiFieldInfo.ProductName, multiFieldInfo.ServiceName, fieldPath)
		for _, item := range v {
			multiFieldInfo.Field = item.Field
			hostList := strings.Split(item.Hosts, ",")
			for _, hostIp := range hostList {
				multiFieldInfo.Hosts = hostIp
				currentInfo := *&multiFieldInfo
				rlt := modifyMultiField(&currentInfo, tx)
				if _, ok := rlt.(error); ok {
					return rlt
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	defer func() {
		clusterInfo, err := model.DeployClusterList.GetClusterInfoById(clusterId)
		if err != nil {
			log.Errorf("%v", err)
			return
		}

		paramsContent := ""
		for k, v := range params {
			jsonStr, err := json.Marshal(v)
			if err != nil {
				log.Errorf("unmarshal json error: %v", err)
				return
			}
			paramsContent += k + ":" + string(jsonStr)
		}

		if err := addSafetyAuditRecord(ctx, "集群运维", "服务参数修改", "集群名称: "+clusterInfo.Name+", 组件名称: "+
			multiFieldInfo.ProductName+", 服务名称: "+multiFieldInfo.ServiceName+", 服务参数: "+paramsContent); err != nil {
			log.Errorf("failed to add safety record, error: %v", err)
		}
	}()

	return
}

// ModifySchemaField
// @Description  	Modify Schema Field
// @Summary      	修改schema字段
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           field_path formData string true "field key"
// @Param           field      formData string true "field value"
// @Success         200        {string}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/{product_name}/service/{service_name}/modify_schema_field [post]
func ModifySchemaField(ctx context.Context) (rlt apibase.Result) {
	info := model.SchemaFieldModifyInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if info.ProductName == "" || info.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}
	clusterId, err := GetCurrentClusterFromParam(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	info.ClutserId = clusterId
	info.FieldPath = ctx.FormValue("field_path")
	fp := strings.Split(info.FieldPath, ".")
	if fp[0] != "Instance" && fp[0] != "Config" {
		return fmt.Errorf("field_path format error")
	}
	info.Field = ctx.FormValue("field")

	// 检测修改的配置是否之前配置过多个配置值，若存在，删除之前的记录
	CheckMultiFieldConfig(clusterId, info.ProductName, info.ServiceName, info.FieldPath)

	// 获取k8s模式部署下当前应用所处的命名空间
	namespace := ctx.FormValue("namespace")
	if namespace != "" {
		info.Namespace = namespace
	}

	return CommonModifySchemaField(&info)
}

// ModifySchemaFieldBatch
// @Description  	Modify Schema Field
// @Summary      	修改scheme多个字段
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           message body string true "[{"clusterId":"","field_path":"","field":"","namespace":""}]"
// @Success         200 {string}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/{product_name}/service/{service_name}/modify_schema_field_batch [post]
func ModifySchemaFieldForDevOps(ctx context.Context) (rlt apibase.Result) {
	info := model.SchemaFieldModifyInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if info.ProductName == "" || info.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}

	var reqParams struct {
		ClusterId int    `json:"clusterId"`
		FieldPath string `json:"field_path"`
		Field     string `json:"field"`
	}
	err := ctx.ReadJSON(&reqParams)

	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	clusterId := reqParams.ClusterId

	info.ClutserId = clusterId
	info.FieldPath = reqParams.FieldPath
	fp := strings.Split(info.FieldPath, ".")
	if fp[0] != "Instance" && fp[0] != "Config" {
		return fmt.Errorf("field_path format error")
	}
	info.Field = reqParams.Field

	// 检测修改的配置是否之前配置过多个配置值，若存在，删除之前的记录
	CheckMultiFieldConfig(clusterId, info.ProductName, info.ServiceName, info.FieldPath)

	return CommonModifySchemaField(&info)
}

func ModifySchemaFieldBatch(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("[Product->ModifySchemaFieldBatch] ConfigAlterGroups from EasyMatrix API ")
	info := model.SchemaFieldModifyInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if info.ProductName == "" || info.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}
	var reqParams []struct {
		ClusterId int    `json:"clusterId"`
		FieldPath string `json:"field_path"`
		Field     string `json:"field"`
		Namespace string `json:"namespace"`
	}
	err := ctx.ReadJSON(&reqParams)

	if err != nil {
		return err
	}

	for _, reqParam := range reqParams {

		if reqParam.Namespace != "" {
			info.Namespace = reqParam.Namespace
		}
		info.FieldPath = reqParam.FieldPath
		info.ClutserId = reqParam.ClusterId
		log.Infof("clusterid: %v", info.ClutserId)
		// 检测修改的配置是否之前配置过多个配置值，若存在，删除之前的记录
		CheckMultiFieldConfig(info.ClutserId, info.ProductName, info.ServiceName, info.FieldPath)
		fp := strings.Split(info.FieldPath, ".")
		if fp[0] != "Instance" && fp[0] != "Config" {
			return fmt.Errorf("field_path format error")
		}
		info.Field = reqParam.Field
		if CommonModifySchemaField(&info) != nil {
			return fmt.Errorf("CommonModifySchemaField err")
		}
	}
	return nil
}

func CheckMultiFieldConfig(clusterId int, productName, serviceName, fieldPath string) (rlt interface{}) {
	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()

	multiFieldIdList, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndPath(clusterId, productName, serviceName, fieldPath)
	if err != nil {
		log.Errorf("get multi field list error: %v", err)
	}

	if len(multiFieldIdList) > 0 {
		query := fmt.Sprintf("update %s set is_deleted = 1 where id in (%s)", model.TBL_SCHEMA_MULTI_FIELD, strings.Join(multiFieldIdList, ","))
		if _, err := tx.Exec(query); err != nil {
			log.Errorf("%v", err)
			return err
		}

		if err = tx.Commit(); err != nil {
			log.Errorf("%v", err)
			return err
		}
	}
	return nil

}

func CommonModifySchemaField(info *model.SchemaFieldModifyInfo) (rlt interface{}) {

	var err error

	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(info.ProductName)
	if err != nil {
		return err
	}
	err, userInfo := model.UserList.GetInfoByUserId(1)
	if err != nil {
		log.Errorf("GetInfoByUserId %v", err)
		return err
	}
	reg := regexp.MustCompile(`(?i).*password.*`)
	if reg.Match([]byte(info.FieldPath)) {
		password, err := encrypt.PlatformEncrypt.SchemeDecrypt(info.Field, userInfo.PassWord)
		if err != nil {
			return err
		}
		info.Field = password
		// password, err := aes.AesDecryptByPassword(info.Field, userInfo.PassWord)
		//if err != nil {
		//	return err
		//} else {
		//	info.Field = password
		//}
	}

	rlt = modifyField(info, tx)
	if _, ok := rlt.(error); ok {
		return rlt
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return rlt
}

func modifyMultiField(info *model.SchemaMultiFieldInfo, tx *sqlx.Tx) (rlt interface{}) {
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE cluster_id=? AND product_name=? AND service_name=? AND field_path=? "+
		"AND hosts=?", model.TBL_SCHEMA_MULTI_FIELD)
	if _, err := tx.Exec(deleteQuery, info.ClusterId, info.ProductName, info.ServiceName, info.FieldPath, info.Hosts); err != nil {
		log.Errorf("%v", err)
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (cluster_id, product_name, service_name, field_path, field, hosts)"+
		" VALUES (:cluster_id, :product_name, :service_name, :field_path, :field, :hosts) "+
		"ON DUPLICATE KEY UPDATE field=:field, update_time=NOW(), is_deleted=0", model.TBL_SCHEMA_MULTI_FIELD)
	if _, err := tx.NamedExec(query, info); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func modifyField(info *model.SchemaFieldModifyInfo, tx *sqlx.Tx) (rlt interface{}) {
	var err error

	var risks []string
	switch info.FieldPath {
	case "Instance.PostDeploy", "Instance.PostUndeploy":
		if currentInfo, err := model.DeployClusterProductRel.GetCurrentProductByProductNameClusterId(info.ProductName, info.ClutserId); err == nil {
			dir := filepath.Join(base.WebRoot, info.ProductName, currentInfo.ProductVersion)
			risks = schema.GetRisks(dir, info.Field)
		}
	}
	//if the filed's json tag is omitempty and the field is showed in the front
	//i think it can not be a zero value
	if fpaths := strings.Split(info.FieldPath, "."); fpaths[0] == "Instance" {
		ins, _ := reflect.TypeOf(schema.ServiceConfig{}).FieldByName("Instance")
		fpath, _ := ins.Type.Elem().FieldByName(fpaths[1])
		jsontag := fpath.Tag.Get("json")
		if strings.Contains(jsontag, "omitempty") {
			zeroValue := zeroValue(reflect.Zero(fpath.Type))
			if zeroValue == info.Field {
				return fmt.Errorf("%s can not be zero value `%s`", fpaths[1], info.Field)
			}
		}
	}

	log.Debugf("SchemaFieldModifyInfo: %+v", info)
	value, err := model.DeploySchemaFieldModify.GetFieldValue(info.ClutserId, info.ProductName, info.ServiceName, info.FieldPath)
	if err == nil && value == info.Field {
		log.Debugf("db %s field no change", info.FieldPath)
		return nil
	}
	query := "INSERT INTO " + model.DeploySchemaFieldModify.TableName +
		" (cluster_id, product_name, service_name, field_path, field, namespace) VALUES" +
		" (:cluster_id, :product_name, :service_name, :field_path, :field, :namespace) ON DUPLICATE KEY UPDATE field=:field,namespace=:namespace, update_time=NOW()"
	if _, err = tx.NamedExec(query, &info); err != nil {
		log.Errorf("%v", err)
		return fmt.Errorf("[modifyField] insert schema field error:%s", err)
	}

	if len(risks) > 0 {
		return "参数中包含危险命令： " + strings.Join(risks, ", ")
	} else {
		return nil
	}
}

// return the zero value of a type and format it as a string type.
// eg: type string's zero value is ""
//    type int's zero value is 0
func zeroValue(v reflect.Value) string {

	switch v.Kind() {
	case reflect.String:
		return v.String()
	//case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
	//	return strconv.FormatInt(v.Int(), 10)
	//case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
	//	return strconv.FormatUint(v.Uint(), 10)
	//case reflect.Float32:
	//	return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	//case reflect.Float64:
	//	return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	//case reflect.Bool:
	//	return strconv.FormatBool(v.Bool())
	default:
		return ""
	}
}

// ResetSchemaField
// @Description  	Reset Schema Field
// @Summary      	重置schema字段
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           field_path formData string true "field key"
// @Param           product_version      formData string true "版本号"
// @Success         200         {string}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/{product_name}/service/{service_name}/reset_schema_field [post]
func ResetSchemaField(ctx context.Context) (rlt apibase.Result) {
	var err error

	info := model.SchemaFieldModifyInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if info.ProductName == "" || info.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}
	info.FieldPath = strings.Replace(ctx.FormValue("field_path"), "current", "Value", 1)
	fp := strings.Split(info.FieldPath, ".")
	if fp[0] != "Instance" && fp[0] != "Config" {
		return fmt.Errorf("field_path format error")
	}
	productVersion := ctx.FormValue("product_version")
	if productVersion == "" {
		return fmt.Errorf("product_version is empty")
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	info.ClutserId = clusterId

	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	if namespace != "" {
		info.Namespace = namespace
	}

	log.Infof("ResetieldModifyInfo: %+v", info)

	tx := model.USE_MYSQL_DB().MustBegin()
	defer func() {
		if _, ok := rlt.(error); ok {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
			rlt = r
		}
	}()
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(info.ProductName)
	if err != nil {
		return fmt.Errorf("[ResetSchemaField]  check product read for deploy err: %s", err)
	}
	isNewConfig := false
	// only the service's config can be added in the front.
	if fp[0] == "Config" {
		productinfo, err := model.DeployProductList.GetByProductNameAndVersion(info.ProductName, productVersion)
		if err != nil {
			log.Errorf("[ResetSchemaField] get from deployproductlist err %v,"+
				"productname is %s, productversion is %s", err, info.ProductName, productVersion)
			return fmt.Errorf("[ResetSchemaField] get by product name and version err: %s", err)
		}
		oriproduct, err := schema.Unmarshal(productinfo.Schema)
		if err != nil {
			log.Errorf("[ResetSchemaField] Unmarshal productlist schema err: %v", err)
			return err
		}
		//the config would like to be reset is added after upload
		if _, exist := oriproduct.Service[info.ServiceName].Config[fp[1]]; !exist {
			isNewConfig = true
			product, err := schema.Unmarshal(productinfo.Product)
			if err != nil {
				log.Errorf("[ResetSchemaField] Unmarshal productlist product err: %v", err)
				return err
			}
			defaultValue := product.Service[info.ServiceName].Config[fp[1]].(schema.VisualConfig).Default.(*string)
			info.Field = *defaultValue

			// 区分k8s模式和主机模式产品包的参数配置重置
			query := "UPDATE " + model.DeploySchemaFieldModify.TableName + " SET field=:field " +
				"WHERE cluster_id=:cluster_id AND product_name=:product_name " +
				"AND service_name=:service_name AND field_path=:field_path AND namespace=:namespace"
			if _, err = tx.NamedExec(query, &info); err != nil {
				log.Errorf("[ResetSchemaField] update deployschemamodifyfiled err: %v,"+
					" filed is %s,clusterid is %d, productname is %s, servicename is %s, filedpath is %s, namespace is %s",
					err, info.Field, info.ClutserId, info.ProductName, info.ServiceName, info.FieldPath, info.Namespace)
				return fmt.Errorf("[ResetSchemaField] update deploy_schema_field_mpdify err: %s", err)
			}
		}
	}
	if !isNewConfig {
		// 区分k8s模式和主机模式下产品包的参数配置删除
		query := "DELETE FROM " + model.DeploySchemaFieldModify.TableName +
			" WHERE cluster_id=:cluster_id AND product_name=:product_name " +
			"AND service_name=:service_name AND field_path=:field_path AND namespace=:namespace"
		if _, err = tx.NamedExec(query, &info); err != nil {
			log.Errorf("[ResetSchemaField] delete deployschemamodifyfiled %v,"+
				" clusterid is %d, productname is %s, servicename is %s, fieldpath is %s, namespace is %s",
				err, info.ClutserId, info.ProductName, info.ServiceName, info.FieldPath, info.Namespace)
			return fmt.Errorf("[ResetSchemaField] delete deploy_schema_field_mpdify err: %s", err)
		}
	}

	query := fmt.Sprintf("update %s set is_deleted = 1 where cluster_id=:cluster_id and product_name=:product_name and service_name=:service_name and field_path=:field_path and is_deleted=0", model.TBL_SCHEMA_MULTI_FIELD)
	if _, err := tx.NamedExec(query, &info); err != nil {
		log.Errorf("%v", err)
		return fmt.Errorf("[ResetSchemaField] update deploy_schema_multi_field err: %s", err)
	}

	return tx.Commit()
}

// ResetSchemaMultiField
// @Description  	Modify Schema Field
// @Summary      	修改schema字段
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           field_path formData string true "field key"
// @Param           product_version      formData string true "版本号"
// @Success         200         {string}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/{product_name}/service/{service_name}/reset_multi_schema_field [post]

func ResetSchemaMultiField(ctx context.Context) apibase.Result {
	var err error
	info := model.SchemaMultiFieldInfo{
		ProductName: ctx.Params().Get("product_name"),
		ServiceName: ctx.Params().Get("service_name"),
	}
	if info.ProductName == "" || info.ServiceName == "" {
		return fmt.Errorf("product_name or service_name is empty")
	}
	info.FieldPath = strings.Replace(ctx.FormValue("field_path"), "current", "Value", 1)
	fp := strings.Split(info.FieldPath, ".")
	if fp[0] != "Config" {
		return fmt.Errorf("field_path format error")
	}
	productVersion := ctx.FormValue("product_version")
	hosts := ctx.FormValue("hosts")
	hostList := strings.Split(hosts, ",")
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	info.ClusterId = clusterId
	err = model.DeployClusterProductRel.CheckProductReadyForDeploy(info.ProductName)
	if err != nil {
		return fmt.Errorf("reset schema muliti field error: %v", err)
	}

	if fp[0] == "Config" {
		productInfo, err := model.DeployProductList.GetByProductNameAndVersion(info.ProductName, productVersion)
		if err != nil {
			log.Errorf("get product info error:%v, product_name:%s", err, info.ProductName)
			return err
		}
		product, err := schema.Unmarshal(productInfo.Product)
		if err != nil {
			log.Errorf("unmarshal product error: %v", err)
		}
		defaultValue := product.Service[info.ServiceName].Config[fp[1]].(schema.VisualConfig).Default.(*string)
		info.Field = *defaultValue
		for _, ip := range hostList {
			if ip != "" {
				info.Hosts = ip
				currentInfo := *&info
				err := model.SchemaMultiField.UpdateField(currentInfo)
				if err != nil {
					log.Errorf("update field error: %v", err)
				}
			}
		}
	}
	defer func() {
		tx := model.USE_MYSQL_DB().MustBegin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		ValidateSchemaFields(clusterId, info.ProductName, info.ServiceName)
		_ = tx.Commit()
	}()
	return nil
}

func AvailableHosts(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	serviceName := ctx.Params().Get("service_name")

	if productName == "" || serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name or service_name is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	clusterId, err := ctx.URLParamInt("baseClusterId")
	if err != nil {
		clusterId, err = ctx.URLParamInt("clusterId")
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
	}
	var ipList string
	query := fmt.Sprintf("SELECT ip_list FROM %s WHERE product_name=? AND service_name=? AND cluster_id=?",
		model.DeployServiceIpList.TableName)
	log.Debugf("%v", query)
	if err := model.USE_MYSQL_DB().Get(&ipList, query, productName, serviceName, clusterId); err != nil && err != sql.ErrNoRows {
		log.Errorf("%v", err)
		return err
	}

	hostInfo := []model.HostInfo{}
	ipFilter := ""
	for _, ip := range strings.Split(ipList, IP_LIST_SEP) {
		ipFilter = ipFilter + "'" + ip + "',"
	}
	if ipFilter != "" {
		ipFilter = ipFilter[:len(ipFilter)-1]
	}
	query = fmt.Sprintf("SELECT %s.* FROM %s "+
		"LEFT JOIN %s on %s.sid=%s.sid "+
		"WHERE steps>=3 AND updated>SUBDATE(NOW(),INTERVAL 3 MINUTE) AND %s.ip NOT IN ("+ipFilter+") AND clusterId=? AND %s.is_deleted=0",
		model.DeployHostList.TableName,
		model.DeployHostList.TableName,
		model.DeployClusterHostRel.TableName,
		model.DeployHostList.TableName,
		model.DeployClusterHostRel.TableName,
		model.DeployHostList.TableName,
		model.DeployClusterHostRel.TableName,
	)

	log.Debugf("%v", query)
	if err := model.USE_MYSQL_DB().Select(&hostInfo, query, clusterId); err != nil {
		log.Errorf("%v", err)
		return err
	}

	return map[string]interface{}{"hosts": hostInfo}
}

func SelectedHosts(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	serviceName := ctx.Params().Get("service_name")

	if productName == "" || serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name or service_name is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	clusterId, err := ctx.URLParamInt("clusterId")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	var ipList string
	query := fmt.Sprintf("SELECT ip_list FROM %s WHERE cluster_id=? and product_name=? and service_name=?", model.DeployServiceIpList.TableName)
	if err := model.USE_MYSQL_DB().Get(&ipList, query, clusterId, productName, serviceName); err != nil {
		log.Errorf("%v", err)
		if err == sql.ErrNoRows {
			var result = make([]string, 0)
			return map[string]interface{}{
				"hosts": result,
			}
		}
		return err
	}

	hostList := strings.Split(ipList, IP_LIST_SEP)
	return map[string]interface{}{
		"hosts": hostList,
	}

}

func ProductStatus(ctx context.Context) apibase.Result {
	parentProductName := ctx.URLParam("parentProductName")
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return map[string]interface{}{
			"count": 0,
			"list":  []map[string]interface{}{},
		}
	}
	productNames, err := model.DeployProductList.GetDeploySonProductName(parentProductName, clusterId)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	count := len(productNames)
	list := []map[string]interface{}{}
	for _, productName := range productNames {
		m := map[string]interface{}{}
		m["product_name"] = productName
		status := ServicesStatus(productName, clusterId)
		switch status.(type) {
		case error:
			return status.(error)
		case map[string]interface{}:
			m["service_count"] = status.(map[string]interface{})["count"]
			m["service_list"] = status.(map[string]interface{})["list"]
		}
		list = append(list, m)
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i]["service_count"].(int) > list[j]["service_count"].(int)
	})
	return map[string]interface{}{
		"count": count,
		"list":  list,
	}
}

func GetProductAnomalyService(ctx context.Context) apibase.Result {
	parentProductName := ctx.URLParam("parentProductName")
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return map[string]interface{}{
			"count": 0,
			"list":  []map[string]string{},
		}
	}
	productNames, err := model.DeployProductList.GetDeploySonProductName(parentProductName, clusterId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	var list = make([]map[string]string, 0)
	for _, productName := range productNames {
		status := ServicesStatus(productName, clusterId)
		switch status.(type) {
		case error:
			return status.(error)
		case map[string]interface{}:
			anomalyList := make([]map[string]interface{}, 0)
			j, _ := json.Marshal(status.(map[string]interface{})["list"])
			json.Unmarshal(j, &anomalyList)
			for _, anomaly := range anomalyList {
				anomalyService := map[string]string{
					"product_name": productName,
					"service_name": anomaly["service_name"].(string),
					"group":        anomaly["group"].(string),
				}
				list = append(list, anomalyService)
			}
		}
	}
	return map[string]interface{}{
		"count": len(list),
		"list":  list,
	}
}

// 目录树节点
type Node struct {
	Path     string  `json:"path"`
	Name     string  `json:"name"`
	Children []*Node `json:"children"`
	indexMap map[string]int
}

// 目录树根节点
type NodeRoot struct {
	root *Node
}

func (r *NodeRoot) add(path string) {
	pathNames := strings.Split(path, "/")
	// 生成前缀树父节点
	nd := r.root
	for index, pathname := range pathNames {
		// 判断目标对象是否存在前缀树的子节点中，若不存在则开始添加子节点信息
		i, exist := nd.indexMap[pathname]
		if !exist {
			nd.indexMap[pathname] = len(nd.Children)
			i = nd.indexMap[pathname]
			// 往前缀树中添加子节点
			nd.Children = append(nd.Children, &Node{
				Path:     strings.Join(pathNames[0:index+1], "/"),
				Name:     pathname,
				Children: []*Node{},
				indexMap: map[string]int{},
			})
		}
		nd = nd.Children[i]
	}
}

func PatchPath(ctx context.Context) apibase.Result {
	// 获取入参以及判断入参是否为空
	paramErrs := apibase.NewApiParameterErrors()
	productId := ctx.URLParam("product_id")

	if productId == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_id is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	// 根据产品包id获取产品包相关信息
	var productInfo model.DeployProductListInfo
	getProductinfo := fmt.Sprintf("select * from %s where id=%s", model.DeployProductList.TableName, productId)
	if err := model.DeployProductList.GetDB().Get(&productInfo, getProductinfo); err != nil {
		return err
	}

	rootpath := "/matrix/easyagent" + "/" + productInfo.ProductName + "/" + productInfo.ProductVersion

	// 判断目录是否存在
	_, err := os.Lstat(rootpath)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	dirpathlist := []string{}
	zipfilelist := []string{}

	// 获取指定产品包版本目录下的文件列表
	f, err := os.Open(rootpath)
	if err != nil {
		log.Errorf("PATH ERROR: %v", err)
		return err
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		log.Errorf("GET FILE NAME ERROR: %v", err)
		return err
	}
	f.Close()

	//过滤出zip包文件
	for _, i := range names {
		if strings.Contains(i, ".zip") {
			zipfilelist = append(zipfilelist, i)
		}
	}

	// 获取指定版本的产品包目录下zip包文件内容
	zipFile := zip.ReadCloser{}
	for _, z := range zipfilelist {
		zipFile, err := zip.OpenReader(rootpath + "/" + z)
		if err != nil {
			log.Errorf("READ ZIP FILE ERROR: %v", err)
			return err
		}
		serviceName := strings.Split(z, ".zip")[0]
		for _, f := range zipFile.File {
			fio := f.FileInfo()
			if fio.IsDir() {
				continue
			}

			dirpathlist = append(dirpathlist, serviceName+"/"+f.Name)
		}
	}
	defer zipFile.Close()

	// 生成目录树结构
	nodes := NodeRoot{root: &Node{
		Path:     productInfo.ProductName + "/" + productInfo.ProductVersion,
		Name:     productInfo.ProductName + "/" + productInfo.ProductVersion,
		Children: []*Node{},
		indexMap: map[string]int{},
	}}

	for _, s := range dirpathlist {
		nodes.add(s)
	}

	list := []interface{}{nodes.root}
	return map[string]interface{}{
		"path": list,
	}
}

func ServiceConfig(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ServiceConfig] ServiceConfig from EasyMatrix API ")

	// 获取入参以及判断入参是否为空
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	file := ctx.URLParam("file")

	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if productVersion == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_version is empty"))
	}
	if serviceName == "" {
		paramErrs.AppendError("$", fmt.Errorf("service_name is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	paramErrs.CheckAndThrowApiParameterErrors()

	// 根据产品包名称和版本从表deployment_product_list获取该产品包的信息
	info, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err: %v", err)
		return err
	}
	// 根据产品包id和集群id获取从deploy_cluster_product_list_rel表中获取产品包的信息
	var productRel model.ClusterProductRel

	productRel, err = model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(info.ID, clusterId, namespace)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err: %v", err)
		return fmt.Errorf("[ServiceConfig] Database err: %v", err)
	}

	//productRel, err = model.DeployClusterProductRel.GetByPidAndClusterId(info.ID, clusterId)
	//if err != nil {
	//	log.Errorf("[ServiceConfig] Database err: %v", err)
	//	return fmt.Errorf("[ServiceConfig] Database err: %v", err)
	//}

	// 反序列化，获取product_parsed，为schema提供解析过的Default值
	originParsed, err := schema.Unmarshal(productRel.ProductParsed)
	if err != nil {
		log.Errorf("[ServiceConfig] Unmarshal err: %v", err)
		return err
	}
	// 反序列化获取schema，代表产品包默认信息，用来与product比较，判断新增和修改的参数
	originSchema, err := schema.Unmarshal(info.Schema)
	if err != nil {
		log.Errorf("[ServiceConfig] Unmarshal err: %v", err)
		return err
	}
	// 获取product,记录当前存在的服务组件的配置参数，em-2.13之前版本的参数只新增在product上
	product, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("[ServiceConfig] Unmarshal err: %v", err)
		return err
	}
	// 从表deploy_schema_field_modify获取记录各个集群下的组件修改和新增配置记录
	schemaFileds, err := model.DeploySchemaFieldModify.GetByProductNameAndServiceAndClusterId(productName, serviceName, clusterId, namespace)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err:%v", err.Error())
		return fmt.Errorf("[ServiceConfig] Database err:%v", err)
	}

	multiSchemaFields, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndClusterId(productName, serviceName, clusterId)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err:%v", err)
		return err
	}

	// 存放全局配置参数列表
	list := []map[string]interface{}{}
	// string记录config_name, int记录config_name在list中的位置
	schemaSet := map[string]int{}
	count := 0
	// 存放解析后的配置参数map，{参数名：参数值}
	parsedSet := make(map[string]string)

	/*@获取产品包中原始schema的服务配置内容和解析后的内容：
	 * 将产品包中schema文件中默认存在的配置参数解析存入schemaSet map.
	 * 将schema文件中配置参数经过平台解析后的结果存入parsedSet map.
	 */
	for k := range originSchema.Service[serviceName].Config { // 标记默认存在的config，默认中没有的参数即为新增
		switch originParsed.Service[serviceName].Config[k].(schema.VisualConfig).Default.(type) {
		case *string:
			schemaSet[k] = -1
			//替换schema参数为解析后的结果
			parsedSet[k] = *originParsed.Service[serviceName].Config[k].(schema.VisualConfig).Default.(*string)
		}
	} //此时若schemaSet[k] == 0 则说明是新的config

	/*
	 * @产品包下指定服务组件全部参数获取：
	 *  1、获取服务组件下所有配置文件，过滤出配置项的占位符
	 *  2、获取product map中service部分的配置参数，和parsedSet map进行匹配，
	 *     在存入全局配置参数列表的同时判断配置参数是否为新增以及是否被引用.
	 */
	// 获取服务用到的所有配置文件"[conf/a.conf,conf/b.conf]"
	configFileList := []string{}
	if _, ok := product.Service[serviceName]; ok && product.Service[serviceName].Instance != nil {
		configFileList = product.Service[serviceName].Instance.ConfigPaths
	}

	// 获取配置文件中的占位符
	placeholderKey := make(map[string]string)
	for _, cfile := range configFileList {
		targetFile := filepath.Join(base.WebRoot, productName, productVersion, serviceName, cfile)
		file, err := ioutil.ReadFile(targetFile)
		if err != nil {
			log.Errorf("[ServiceConfig] Config File err: %v", err)
			return fmt.Errorf("[ServiceConfig] Get Config File err: %v", err)
		}
		reg1 := regexp.MustCompile(`\{{\.\w+}}`)
		reg2 := regexp.MustCompile(`\w+`)
		for _, v := range reg1.FindAllString(string(file), -1) {
			placeholderKey[reg2.FindString(v)] = "1"
		}
	}

	//获取product map中service部分的配置参数存入全局配置参数列表，同时判断配置参数是否为新增以及是否被引用
	isNewConfig := make(map[int]bool) //补丁：判断新增
	var keys []string
	for k := range product.Service[serviceName].Config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	configMap := product.Service[serviceName].Config
	for _, k := range keys {
		// 忽略带有"${@"格式的参数在列表中出现
		if strings.HasPrefix(*configMap[k].(schema.VisualConfig).Default.(*string), "${@") {
			continue
		}
		m := map[string]interface{}{}
		// 判断是否新增参数
		if _, exists := parsedSet[k]; exists {
			m["config"] = k
			m["default"] = parsedSet[k]
			m["current"] = parsedSet[k]
			m["isnew"] = false
		} else {
			m["config"] = k
			m["default"] = *configMap[k].(schema.VisualConfig).Default.(*string)
			m["current"] = *configMap[k].(schema.VisualConfig).Value.(*string)
			m["isnew"] = true
		}

		m["updated"] = ""
		count++
		// 将新增加的参数添加到schemaSet字典中，同时设置key的位置
		schemaSet[k] = count

		// 判断参数是否被引用
		if _, ok := placeholderKey[k]; ok {
			m["nouse"] = false
		} else {
			m["nouse"] = true
		}

		// 生成服务的全局配置参数列表
		list = append(list, m)
		if m["isnew"].(bool) {
			isNewConfig[count] = true
		}
	}
	/*
	 *@配置参数更新逻辑：
	 * 将schemaSet内容与deploy_schema_filed_modify表中的记录比较
	 * 如果参数有修改则更新全局配置参数列表list中的参数
	 */
	for _, v := range schemaFileds {
		// 从deploy_schema_filed_modify表中获取新增字段或修改的字段名称
		config := strings.Split(v.FieldPath, ".")[1]
		if schemaSet[config] != 0 { // 如果list存在这个config_name
			// 获取参数列表中对应的新增字段或修改字段的原始值，这里将其原始默认值替换为修改后的值
			list[schemaSet[config]-1]["current"] = v.Field
			if list[schemaSet[config]-1]["default"] != list[schemaSet[config]-1]["current"] || list[schemaSet[config]-1]["isnew"] == true {
				list[schemaSet[config]-1]["updated"] = v.UpdateDate.String()
				isNewConfig[schemaSet[config]] = false
			}
		}
	}

	configFieldMap := map[string]map[string][]string{}
	multiUpdateTimeMap := map[string]interface{}{}
	// 保存有序的map顺序
	keyMap := map[string][]string{}
	for _, v := range multiSchemaFields {
		config := strings.Split(v.FieldPath, ".")[1]
		_, ok := configFieldMap[config]
		if !ok {
			configFieldMap[config] = make(map[string][]string, 0)
			keyMap[config] = make([]string, 0)
		}
		fieldConfigMap := configFieldMap[config]
		_, ok = fieldConfigMap[v.Field]
		if !ok {
			fieldConfigMap[v.Field] = make([]string, 0)
			keyList := keyMap[config]
			if !contains(keyList, v.Field) {
				keyList = append(keyList, v.Field)
			}
			keyMap[config] = keyList
		}
		hostSlice := fieldConfigMap[v.Field]
		hostSlice = append(hostSlice, v.Hosts)
		fieldConfigMap[v.Field] = hostSlice
		configFieldMap[config] = fieldConfigMap
		multiUpdateTimeMap[config] = v.UpdateTime.String()
	}

	for k, v := range configFieldMap {
		if schemaSet[k] != 0 {
			fieldConfigSlice := []model.FieldConfig{}
			keyList := keyMap[k]
			for _, field := range keyList {
				hosts := v[field]
				fieldConfigSlice = append(fieldConfigSlice, model.FieldConfig{Field: field, Hosts: strings.Join(hosts, ",")})
			}
			list[schemaSet[k]-1]["current"] = fieldConfigSlice
			if updateTime, ok := multiUpdateTimeMap[k]; ok {
				list[schemaSet[k]-1]["updated"] = updateTime
			}
		}
	}

	//补丁，填充老版本新增修改时间
	for i := 1; i <= count; i++ {
		if isNewConfig[i] {
			list[i-1]["updated"] = time.Time{}
		}
	}
	//配置文件中所有 key 包含 password 不区分大小写的配置 value 要用对称加密后返回 秘钥为 admin 的密码
	err, userInfo := model.UserList.GetInfoByUserId(1)
	if err != nil {
		log.Errorf("GetInfoByUserId %v", err)
		return fmt.Errorf("[ServiceConfig] Database Query GetInfoByUserId: %v", err)
	}

	// 判断file入参数是否指定具体文件
	if file == "" {
		//util.FilterAndEncryptedPassword(fileConfigList,userInfo.PassWord)
		reg := regexp.MustCompile(`(?i).*password.*`)
		for _, fileConfig := range list {
			if reg.Match([]byte(fileConfig["config"].(string))) {
				currentValue, err := aes.AesEncryptByPassword(fileConfig["current"].(string), userInfo.PassWord)
				if err != nil {
					return err
				}
				defaultValue, err := aes.AesEncryptByPassword(fileConfig["current"].(string), userInfo.PassWord)
				if err != nil {
					return err
				}
				fileConfig["current"] = currentValue
				fileConfig["default"] = defaultValue
			}
		}
		return map[string]interface{}{
			"count": count,
			"list":  list,
		}
	} else {
		/*
		*@对于指定组件的配置文件当中的参数显示如下：
		*  1、获取指定配置文件内容过滤出配置项的占位符
		*  2、正则表达式匹配配置文件中的占位符
		*  3、将匹配到的占位符当中的key存入map
		*  4、从全局配置参数列表过滤出来指定配置文件拥有的参数列表
		 */
		// 获取指定的配置文件路径
		targetFile := filepath.Join(base.WebRoot, productName, productVersion, serviceName, file)

		file, err := ioutil.ReadFile(targetFile)
		if err != nil {
			log.Errorf("[]%v", err)
			return fmt.Errorf("[ServiceConfig] Get Config File err: %v", err)
		}
		reg1 := regexp.MustCompile(`\{{\.\w+}}`)
		reg2 := regexp.MustCompile(`\w+`)

		placeholderKey := make(map[string]string)
		for _, v := range reg1.FindAllString(string(file), -1) {
			placeholderKey[reg2.FindString(v)] = "1"
		}

		// 全局参数列表中过滤出指定配置文件的配置参数列表
		fileConfigkeyCount := 0
		fileConfigList := []map[string]interface{}{}
		for _, v := range list {
			key := v["config"].(string)
			if _, ok := placeholderKey[key]; ok {
				fileConfigList = append(fileConfigList, v)
				fileConfigkeyCount++
			}
		}

		//util.FilterAndEncryptedPassword(fileConfigList,userInfo.PassWord)
		reg := regexp.MustCompile(`(?i).*password.*`)
		for _, fileConfig := range fileConfigList {
			if reg.Match([]byte(fileConfig["config"].(string))) {
				currentValue, err := aes.AesEncryptByPassword(fileConfig["current"].(string), userInfo.PassWord)
				if err != nil {
					return err
				}
				defaultValue, err := aes.AesEncryptByPassword(fileConfig["current"].(string), userInfo.PassWord)
				if err != nil {
					return err
				}
				fileConfig["current"] = currentValue
				fileConfig["default"] = defaultValue
			}
		}
		return map[string]interface{}{
			"count": fileConfigkeyCount,
			"list":  fileConfigList,
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ConfigAlterGroups(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ConfigAlterGroups] ConfigAlterGroups from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	parentProductName := ctx.URLParam("parentProductName")

	if parentProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("parentProductName is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	list := []map[string]interface{}{}
	count := 0

	productNames, err := model.DeployProductList.GetDeploySonProductName(parentProductName, clusterId)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err:%v", err.Error())
		return err
	}
	for _, productName := range productNames { // 记录有变化或者新增的config
		m := map[string]interface{}{}
		schemaFileds, err := model.DeploySchemaFieldModify.GetByProductNameClusterId(productName, clusterId) // 获取schema_field_modify，记录各个集群的组件的修改和新增配置
		if err != nil {
			log.Errorf("[ServiceConfig] Database err:%v", err.Error())
			return err
		}
		info, err := model.DeployProductList.GetByProductNameClusterId(productName, clusterId)
		if err != nil {
			log.Errorf("[ServiceConfig] Database err:%v", err.Error())
			return err
		}

		serviceNameSet, err := model.DeployInstanceList.GetServiceNameSetByClusterIdAndPid(clusterId, info.ID)
		if err != nil {
			log.Errorf("[ServiceConfig] get serviceNameSet err:%v", err.Error())
			return err
		}

		originSchema, err := schema.Unmarshal(info.Schema) // 获取schema，用来与product比较，判断新增和修改的参数
		if err != nil {
			log.Errorf("[ServiceConfig] Unmarshal err: %v", err)
			return err
		}
		schemaSet := make(map[string]map[string]string)
		configCheck := make(map[string]map[string]bool)
		for serviceName, serviceConfig := range originSchema.Service { // 标记默认存在的config，默认中没有的参数即为新增

			if !serviceNameSet[serviceName] { // 筛选属于该集群的服务
				continue
			}

			schemaSet[serviceName] = make(map[string]string)
			configCheck[serviceName] = make(map[string]bool)

			for config, v := range serviceConfig.Config {
				switch v.(schema.VisualConfig).Default.(type) {
				case *string:
					schemaSet[serviceName][config] = *v.(schema.VisualConfig).Default.(*string)
					configCheck[serviceName][config] = true
				}
			}
		}
		//兼容之前版本
		originProduct, err := schema.Unmarshal(info.Product) // 获取product,记录当前存在的config，em-2.13之前版本的参数只新增在product上
		if err != nil {
			log.Errorf("[ServiceConfig] Unmarshal err:%v", err.Error())
			return err
		}
		newConfig := make(map[string]bool)
		productSet := make(map[string]map[string]bool)
		for serviceName, serviceConfig := range originProduct.Service {

			if !serviceNameSet[serviceName] { // 筛选属于该集群的服务
				continue
			}

			productSet[serviceName] = make(map[string]bool)
			for config := range serviceConfig.Config {
				if !configCheck[serviceName][config] { // 补丁：判断是否为新的config
					newConfig[serviceName] = true
				}
				productSet[serviceName][config] = true
			}
		}

		mSet := make(map[string]bool) // 去重
		m["groups"] = make([]string, 0)

		// 补丁：如果是new config所在的group，直接添加
		for serviceName := range newConfig {
			serviceGroup := originProduct.Service[serviceName].Group
			if !mSet[serviceGroup] {
				m["groups"] = append(m["groups"].([]string), serviceGroup)
				mSet[serviceGroup] = true
			}
		}

		for _, field := range schemaFileds { // 与deploy_schema_filed_modify表中的记录比较
			config := strings.Split(field.FieldPath, ".")
			if config[0] != "Config" {
				continue
			}

			//兼容之前版本
			if !productSet[field.ServiceName][config[1]] { //不存在于当前版本的config
				continue
			}

			if !configCheck[field.ServiceName][config[1]] || schemaSet[field.ServiceName][config[1]] != field.Field { // 如果参数为新增或者有修改，则加入group
				serviceGroup := originSchema.Service[field.ServiceName].Group
				if !mSet[serviceGroup] { // 去掉重复的group
					m["groups"] = append(m["groups"].([]string), serviceGroup)
					mSet[serviceGroup] = true
				}
			}
		}
		if len(m["groups"].([]string)) > 0 {
			m["product_name"] = productName
			list = append(list, m)
			count++
		}
	}
	return map[string]interface{}{
		"count": count,
		"list":  list,
	}
}

func ConfigAlteration(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ConfigAlteration] ConfigAlteration from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	parentProductName := ctx.URLParam("parentProductName")
	productName := ctx.URLParam("ProductName")
	serviceGroup := ctx.URLParam("group")

	if parentProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("parentProductName is empty"))
	}
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("ProductName is empty"))
	}
	if serviceGroup == "" {
		paramErrs.AppendError("$", fmt.Errorf("group is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	list := []map[string]interface{}{}
	count := 0

	info, err := model.DeployProductList.GetByProductNameAndParentProductNameClusterId(productName, parentProductName, clusterId)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err:%v", err.Error())
		return err
	}

	serviceNameSet, err := model.DeployInstanceList.GetServiceNameSetByClusterIdAndPid(clusterId, info.ID)
	if err != nil {
		log.Errorf("[ServiceConfig] get serviceNameSet err:%v", err.Error())
		return err
	}

	productRel, err := model.DeployClusterProductRel.GetByPidAndClusterId(info.ID, clusterId)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err: %v", err)
		return err
	}
	originParsed, err := schema.Unmarshal(productRel.ProductParsed) // 获取product_parsed，为schema提供解析过的Default值
	if err != nil {
		log.Errorf("[ServiceConfig] Unmarshal err:%v", err.Error())
		return err
	}
	originSchema, err := schema.Unmarshal(info.Schema) // 获取schema，代表产品包默认信息，用来与product比较，判断新增和修改的参数
	if err != nil {
		log.Errorf("[ServiceConfig] Unmarshal err:%v", err.Error())
		return err
	}
	schemaSet := make(map[string]map[string]string)
	configCheck := make(map[string]map[string]bool)                // 判断参数是否默认存在
	for serviceName, serviceConfig := range originSchema.Service { // 遍历schema下所有服务

		if !serviceNameSet[serviceName] { // 筛选属于该集群的服务
			continue
		}

		schemaSet[serviceName] = make(map[string]string)
		configCheck[serviceName] = make(map[string]bool)
		for config := range serviceConfig.Config { // 标记默认存在的config，默认中没有的参数即为新增
			switch originParsed.Service[serviceName].Config[config].(schema.VisualConfig).Default.(type) {
			case *string:
				//替换schema参数为解析后的结果
				schemaSet[serviceName][config] = *originParsed.Service[serviceName].Config[config].(schema.VisualConfig).Default.(*string)
				configCheck[serviceName][config] = true
			}
		}
	}

	//兼容之前版本
	originProduct, err := schema.Unmarshal(info.Product) // 获取product,记录当前存在的config，em-2.13之前版本的参数只新增在product上
	if err != nil {
		log.Errorf("[ServiceConfig] Unmarshal err:%v", err.Error())
		return err
	}

	productSet := make(map[string]map[string]bool)
	productDefault := make(map[string]map[string]string)

	// 补丁
	newConfig := make(map[string]map[string]string)
	newConfigSet := make(map[string]map[string]bool)

	for serviceName, serviceConfig := range originProduct.Service { // 遍历product包含的所有服务

		if !serviceNameSet[serviceName] { // 筛选属于该集群的服务
			continue
		}

		productSet[serviceName] = make(map[string]bool)
		productDefault[serviceName] = make(map[string]string)
		newConfig[serviceName] = make(map[string]string)
		newConfigSet[serviceName] = make(map[string]bool)
		for config, v := range serviceConfig.Config {
			switch v.(schema.VisualConfig).Default.(type) {
			case *string:
				if !configCheck[serviceName][config] { // 补丁：获取new的config
					newConfig[serviceName][config] = *v.(schema.VisualConfig).Value.(*string)
					newConfigSet[serviceName][config] = true
				}
				productSet[serviceName][config] = true
				productDefault[serviceName][config] = *v.(schema.VisualConfig).Default.(*string)
			}
		}
	}

	schemaFileds, err := model.DeploySchemaFieldModify.GetByProductNameClusterId(productName, clusterId) // 获取schema_field_modify
	if err != nil {
		log.Errorf("[ServiceConfig] Database err:%v", err.Error())
		return err
	}
	mSet := map[string]int{} //记录service_name在list中的位置

	for _, field := range schemaFileds { // 与deploy_schema_filed_modify表中的记录比较，如果参数有修改则更新参数
		alter := map[string]interface{}{}
		config := strings.Split(field.FieldPath, ".")
		if config[0] != "Config" || serviceGroup != originSchema.Service[field.ServiceName].Group {
			continue
		}
		if schemaSet[field.ServiceName][config[1]] == field.Field && configCheck[field.ServiceName][config[1]] {
			continue
		}

		//兼容之前版本
		if !productSet[field.ServiceName][config[1]] { //不存在于当前版本的config
			continue
		}

		if !configCheck[field.ServiceName][config[1]] {
			alter["isnew"] = true
			newConfigSet[field.ServiceName][config[1]] = false // 补丁 ：筛掉modify表里存在的new的config
			alter["default"] = productDefault[field.ServiceName][config[1]]
		} else {
			alter["isnew"] = false
			alter["default"] = schemaSet[field.ServiceName][config[1]]
		}
		alter["config"] = config[1]
		alter["current"] = field.Field
		alter["updated"] = field.UpdateDate
		if mSet[field.ServiceName] == 0 {
			list = append(list, map[string]interface{}{
				"service_name": field.ServiceName,
				"alteration":   []map[string]interface{}{},
			})
			count++
			mSet[field.ServiceName] = count

		}
		list[mSet[field.ServiceName]-1]["alteration"] = append(list[mSet[field.ServiceName]-1]["alteration"].([]map[string]interface{}), alter)
	}

	// 补丁 兼容em 2.13之前版本，获取只在product中存在，modify中没有的新增参数
	for myServiceName, serviceConfig := range newConfig {
		for configName, value := range serviceConfig {
			if !newConfigSet[myServiceName][configName] { // 补丁：如果modify表中存在则跳过
				continue
			}
			if strings.HasPrefix(value, "${@") {
				continue
			}
			alter := map[string]interface{}{}
			if serviceGroup != originSchema.Service[myServiceName].Group {
				continue
			}
			alter["isnew"] = true
			alter["config"] = configName
			alter["default"] = productDefault[myServiceName][configName]
			alter["current"] = value
			alter["updated"] = time.Time{}
			if mSet[myServiceName] == 0 {
				list = append(list, map[string]interface{}{
					"service_name": myServiceName,
					"alteration":   []map[string]interface{}{},
				})
				count++
				mSet[myServiceName] = count

			}
			list[mSet[myServiceName]-1]["alteration"] = append(list[mSet[myServiceName]-1]["alteration"].([]map[string]interface{}), alter)
		}
	}
	for i := range list {
		sort.SliceStable(list[i]["alteration"].([]map[string]interface{}), func(j, k int) bool {
			return list[i]["alteration"].([]map[string]interface{})[j]["updated"].(time.Time).After(list[i]["alteration"].([]map[string]interface{})[k]["updated"].(time.Time))
		})
	}
	return map[string]interface{}{
		"count": count,
		"list":  list,
	}
}

func ConfigAlterAll(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ConfigAlterGroups] ConfigAlterGroups from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	parentProductName := ctx.URLParam("parentProductName")

	if parentProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("parentProductName is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	productNames, err := model.DeployProductList.GetDeploySonProductName(parentProductName, clusterId)
	if err != nil {
		log.Errorf("[ServiceConfig] Database err:%v", err.Error())
		return err
	}
	list := []map[string]interface{}{}
	count := 0
	for _, productName := range productNames {
		p := map[string]interface{}{}
		g := map[string]map[string]interface{}{}
		info, err := model.DeployProductList.GetByProductNameAndParentProductNameClusterId(productName, parentProductName, clusterId)
		if err != nil {
			log.Errorf("[ServiceConfig] Database err:%v", err.Error())
			return err
		}
		originSchema, err := schema.Unmarshal(info.Schema) // 获取原始schema,记录config默认值
		if err != nil {
			log.Errorf("[ServiceConfig] Unmarshal err:%v", err.Error())
			return err
		}
		schemaSet := make(map[string]map[string]string)
		for serviceName, serviceConfig := range originSchema.Service {
			schemaSet[serviceName] = make(map[string]string)
			for config, v := range serviceConfig.Config {
				schemaSet[serviceName][config] = *v.(schema.VisualConfig).Default.(*string)
			}
		}
		schemaFileds, err := model.DeploySchemaFieldModify.GetByProductNameClusterId(productName, clusterId) // 获取schema_field_modify
		if err != nil {
			log.Errorf("[ServiceConfig] Database err:%v", err.Error())
			return err
		}
		gSet := map[string]bool{}
		for _, field := range schemaFileds {
			alter := map[string]interface{}{}
			config := strings.Split(field.FieldPath, ".")
			serviceGroup := originSchema.Service[field.ServiceName].Group
			if config[0] != "Config" {
				continue
			}
			if schemaSet[field.ServiceName][config[1]] == field.Field && schemaSet[field.ServiceName][config[1]] != "" {
				continue
			}
			if schemaSet[field.ServiceName][config[1]] == "" {
				alter["isnew"] = true
				alter["default"] = field.Field
			} else {
				alter["isnew"] = false
				alter["default"] = schemaSet[field.ServiceName][config[1]]
			}
			alter["config"] = config[1]
			alter["current"] = field.Field
			alter["updated"] = field.UpdateDate
			if !gSet[serviceGroup] {
				g[serviceGroup] = make(map[string]interface{})
				g[serviceGroup]["service_name"] = field.ServiceName
				g[serviceGroup]["alteration"] = []map[string]interface{}{}
				gSet[serviceGroup] = true
			}
			g[serviceGroup]["alteration"] = append(g[serviceGroup]["alteration"].([]map[string]interface{}), alter)
		}
		if len(g) > 0 {
			p["product_name"] = productName
			p["groups"] = g
			list = append(list, p)
			count++
		}
	}
	return map[string]interface{}{
		"count": count,
		"list":  list,
	}
}

func GetProductGroupList(ctx context.Context) apibase.Result {
	log.Debugf("GetProductGroupList: %v", ctx.Request().RequestURI)
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		return fmt.Errorf("product name is null")
	}
	productVersion := ctx.Params().Get("product_version")
	if productVersion == "" {
		return fmt.Errorf("product version is null")
	}

	product, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		log.Errorf(err.Error())
		return fmt.Errorf("database err: %v", err)
	}

	type resultList struct {
		ServiceName        string `json:"service_name"`
		ServiceNameDisplay string `json:"service_name_display"`
		Alert              bool   `json:"alert"`
	}
	groupAndServices := map[string][]resultList{}

	originSchema, err := schema.Unmarshal(product.Schema) // 获取schem
	if err != nil {
		log.Errorf("schema Unmarshal err: %v", err)
		return err
	}

	for name, config := range originSchema.Service {
		if config.Instance == nil {
			continue
		}
		if config.ServiceDisplay == "" {
			config.ServiceDisplay = name
		}
		groupAndServices[config.Group] = append(groupAndServices[config.Group], resultList{
			name,
			config.ServiceDisplay,
			false,
		})
	}
	return map[string]interface{}{
		"groups": groupAndServices,
		"count":  len(groupAndServices),
	}
}

func GetUpgradeCandidateList(ctx context.Context) apibase.Result {
	log.Debugf("GetProductGroupList: %v", ctx.Request().RequestURI)
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		return fmt.Errorf("product name is null")
	}
	productVersion := ctx.Params().Get("product_version")
	if productVersion == "" {
		return fmt.Errorf("product version is null")
	}
	productType := ctx.URLParam("product_type")
	upgradeMode := ctx.URLParam("upgrade_mode")
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	type resultList struct {
		Id             int    `json:"id"`
		ProductVersion string `json:"product_version"`
	}
	var availables []resultList

	if upgradeMode == upgrade.SMOOTH_UPGRADE_MODE {
		infoList, err := upgrade.UpgradeHistory.GetTargetVersionInfo(clusterId, productName, productVersion, productType, upgradeMode)
		if err != nil {
			log.Errorf("GetUpgradeCandidateList-query db error: %v", err)
			return err
		}
		if len(infoList) == 0 {
			goto SUBSEQUENCE
		}
		//存在平滑升级版本
		for _, info := range infoList {
			availables = append(availables, resultList{
				info.ID,
				info.ProductVersion,
			})
		}
		return map[string]interface{}{
			"list":  availables,
			"count": len(availables),
		}
	}
SUBSEQUENCE:
	products, _ := model.DeployProductList.GetProductListByNameAndType(productName, productType, nil)
	for _, p := range products {
		if p.ProductType != 0 {
			continue
		}
		log.Infof("%v, %v, %v", p.ProductName, p.ProductVersion, p.ProductType)
		if CompareVersion(p.ProductVersion, productVersion) > 0 {
			availables = append(availables, resultList{
				p.ID,
				p.ProductVersion,
			})
		}
	}
	return map[string]interface{}{
		"list":  availables,
		"count": len(availables),
	}
}

func removeDuplicate(source []model.FieldConfig) ([]model.FieldConfig, []string) {
	result := []model.FieldConfig{}
	ipList := []string{}
	tempMap := map[string]byte{}
	for _, e := range source {
		l := len(tempMap)
		tempMap[e.Field] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
		ipList = append(ipList, strings.Split(e.Hosts, ",")...)
	}
	return result, ipList
}

type backUpInfo []struct {
	Product string `json:"product"`
	Version string `json:"version"`
	Service []struct {
		Name     string `json:"name"`
		HostInfo []struct {
			Sid string `json:"sid"`
			Ip  string `json:"ip"`
		} `json:"host_info"`
	} `json:"service"`
}

// GetBackupPackage
// @Description  	Get Cluster Service Backup
// @Summary      	获取集群服务备份目录
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           clusterId query string  true "集群ID"
// @Success         200 {object} backUpInfo
// @Router          /api/v2/product/backup [get]
func GetBackupPackage(ctx context.Context) apibase.Result {
	cid, err := strconv.Atoi(ctx.URLParam("clusterId"))
	if err != nil {
		return fmt.Errorf("clusterId is null")
	}
	var (
		path           = model.ClusterBackupConfig.GetPathByClusterId(cid)
		cmd            = fmt.Sprintf("#!/bin/sh\nfind %s -regex  \"^.*-.*-.*~$\"", path)
		backUpInfoResp = backUpInfo{}
		hostList, _    = model.DeployClusterHostRel.GetClusterHostRelList(cid)
		syncMap        = sync.Map{}
		wg             = sync.WaitGroup{}
	)

	for _, host := range hostList {
		wg.Add(1)
		go func(host model.ClusterHostRel) {
			defer wg.Done()
			backUpDirListStr, err := agent.AgentClient.ToExecCmd(host.Sid, "", cmd, "")
			if err != nil {
				return
			}
			syncMap.Store(host.Sid, backUpDirListStr)
		}(host)
	}
	wg.Wait()
	for _, host := range hostList {
		err, hostInfo := model.DeployHostList.GetHostInfoBySid(host.Sid)
		if err != nil {
			return err
		}
		ip := hostInfo.Ip
		// /opt/dtstack/easymanager/mysql-1610970097_-4.0.2~
		// /opt/dtstack/easymanager/mysql-1610970097-4.0.3~
		// /opt/dtstack/easymanager/clean_history-1610970097-4.0.3~
		// /data/bbb/ddd/DTFront/tengine-1641540471-4.2.57_rel_CentOS7_x86_64~
		backUpDirListStr, err := agent.AgentClient.ToExecCmd(host.Sid, "", cmd, "")
		if load, ok := syncMap.Load(host.Sid); ok {
			if strValue, typeOk := load.(string); typeOk {
				backUpDirListStr = strValue
			} else {
				return fmt.Errorf("类型转换失败")
			}
		}
		//去除空格与换行
		backUpDirListStr = strings.TrimSpace(backUpDirListStr)
		if backUpDirListStr == "" {
			log.Infof("sid: %s not found backup dir", host.Sid)
			continue
		}
		// /opt/dtstack/DTBase/mysql-1610970098-4.0.2~
		// /data/bbb/ddd/DTFront/tengine-1641540471-4.2.57_rel_CentOS7_x86_64~
		for _, path := range strings.Split(backUpDirListStr, "\n") {
			var (
				arrPath  = strings.Split(path, "/")
				product  = arrPath[len(arrPath)-2]
				fileName = arrPath[len(arrPath)-1]
				service  = strings.Split(fileName, "-")[0]
				version  = strings.TrimRight(strings.Split(fileName, "-")[2], "~")
			)
			log.Infof("product: %s service %s version %s", product, service, version)

			//第一次 backUpInfoResp 为空
			if len(backUpInfoResp) == 0 {
				backUpInfoResp = append(backUpInfoResp, struct {
					Product string `json:"product"`
					Version string `json:"version"`
					Service []struct {
						Name     string `json:"name"`
						HostInfo []struct {
							Sid string `json:"sid"`
							Ip  string `json:"ip"`
						} `json:"host_info"`
					} `json:"service"`
				}{
					Product: product,
					Version: version,
					Service: []struct {
						Name     string `json:"name"`
						HostInfo []struct {
							Sid string `json:"sid"`
							Ip  string `json:"ip"`
						} `json:"host_info"`
					}{{
						Name: service,
						HostInfo: []struct {
							Sid string `json:"sid"`
							Ip  string `json:"ip"`
						}{
							{Sid: host.Sid, Ip: ip}},
					}},
				})

				continue
			}

			for productIdx, _ := range backUpInfoResp {
				//info:=backUpInfoResp[productIdx]
				if backUpInfoResp[productIdx].Product == product && backUpInfoResp[productIdx].Version == version {
					for serviceIdx, _ := range backUpInfoResp[productIdx].Service {
						//svc:=backUpInfoResp[productIdx].Service[serviceIdx]
						//如果找到了服务
						if backUpInfoResp[productIdx].Service[serviceIdx].Name == service {

							for hostIdx, _ := range backUpInfoResp[productIdx].Service[serviceIdx].HostInfo {
								info := backUpInfoResp[productIdx].Service[serviceIdx].HostInfo[hostIdx]
								//如果 sid 在列表中了 直接 break
								if info.Sid == host.Sid {
									break
								}
								//如果 sid 不在列表中了
								if hostIdx == len(backUpInfoResp[productIdx].Service[serviceIdx].HostInfo)-1 {
									backUpInfoResp[productIdx].Service[serviceIdx].HostInfo = append(backUpInfoResp[productIdx].Service[serviceIdx].HostInfo, struct {
										Sid string `json:"sid"`
										Ip  string `json:"ip"`
									}{
										Sid: host.Sid,
										Ip:  ip,
									})
									break
								}
							}
							break
						}
						//如果没找到该服务 直接添加
						if serviceIdx == len(backUpInfoResp[productIdx].Service)-1 {
							backUpInfoResp[productIdx].Service = append(backUpInfoResp[productIdx].Service, struct {
								Name     string `json:"name"`
								HostInfo []struct {
									Sid string `json:"sid"`
									Ip  string `json:"ip"`
								} `json:"host_info"`
							}{
								Name: service,
								HostInfo: []struct {
									Sid string `json:"sid"`
									Ip  string `json:"ip"`
								}{
									{Sid: host.Sid, Ip: ip}},
							})
							break
						}
					}
					break
				}
				//如果没找到该版本的产品包 直接添加
				if productIdx == len(backUpInfoResp)-1 {
					backUpInfoResp = append(backUpInfoResp, struct {
						Product string `json:"product"`
						Version string `json:"version"`
						Service []struct {
							Name     string `json:"name"`
							HostInfo []struct {
								Sid string `json:"sid"`
								Ip  string `json:"ip"`
							} `json:"host_info"`
						} `json:"service"`
					}{
						Product: product,
						Service: []struct {
							Name     string `json:"name"`
							HostInfo []struct {
								Sid string `json:"sid"`
								Ip  string `json:"ip"`
							} `json:"host_info"`
						}{{
							Name: service,
							HostInfo: []struct {
								Sid string `json:"sid"`
								Ip  string `json:"ip"`
							}{
								{Sid: host.Sid, Ip: ip}},
						}},
						Version: version,
					})
					break
				}
			}
		}
	}
	return backUpInfoResp
}

// CleanBackupPackage
// @Description  	Clean Cluster Backup
// @Summary      	清理备份文件夹
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           message body backUpInfo  true "备份信息"
// @Success         200   {string}  string  "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/clean [post]
func CleanBackupPackage(ctx context.Context) apibase.Result {
	var backup backUpInfo
	err := ctx.ReadJSON(&backup)
	if err != nil {
		return err
	}
	log.Debugf("CleanBackupPackage params %v", backup)
	cid, _ := GetCurrentClusterId(ctx)
	//sid -> pathList
	cleanMap := make(map[string][]string)
	for _, info := range backup {
		for _, svc := range info.Service {
			// 防止恶意删除
			if strings.Contains(info.Product, "*") || strings.Contains(svc.Name, "*") || strings.Contains(info.Version, "*") {
				return fmt.Errorf("包含非法通配符 * ")
			}

			for _, host := range svc.HostInfo {
				sid := host.Sid
				configPath := model.ClusterBackupConfig.GetPathByClusterId(cid)
				path := fmt.Sprintf("%s/%s-*-%s~", path.Join(configPath, info.Product), svc.Name, info.Version)
				if pathList, ok := cleanMap[sid]; ok {
					///opt/dtstack/DTBase/mysql_1610970097_4.0.2~
					cleanMap[sid] = append(pathList, path)
				} else {
					cleanMap[sid] = []string{path}
				}
			}
		}
	}

	for sid, pathList := range cleanMap {
		var cmd strings.Builder
		cmd.WriteString("#!/bin/sh\n")
		for _, path := range pathList {
			cmd.WriteString(fmt.Sprintf("rm -rf %s;\n", path))
		}
		log.Infof("sid: %s  backup package clean shell : %s", sid, strings.TrimSpace(cmd.String()))
		content, err := agent.AgentClient.ToExecCmd(sid, "", strings.TrimSpace(cmd.String()), "")
		if err != nil {
			log.Errorf("execCmd failed, content: %s ;err: %v", content, err)
			return err
		}
	}
	return nil
}

//AutoTestHistory 获取当前产品最近的一次自动化测试记录
func AutoTestHistory(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.Params().Get("product_name")
	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("get cluster id from cookie error"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	//获取deployment_product_list表中产品包信息
	product, err := model.DeployProductList.GetByProductNameClusterId(productName, clusterId)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("product %s is not installed", productName)
	} else if err != nil {
		log.Errorf("err: %v", err.Error())
		return err
	}
	//反序列化，产品包默认schema信息
	originSchema, err := schema.Unmarshal(product.Schema)
	if err != nil {
		log.Errorf("err: %v", err.Error())
		return err
	}
	serviceConfig := originSchema.Service
	//自动化测试开关打开，则返回记录
	if _, ok := serviceConfig[TEST_SET_SERVICE_NAME]; ok && serviceConfig[TEST_SET_SERVICE_NAME].Instance != nil &&
		serviceConfig[TEST_SET_SERVICE_NAME].Instance.TestOn {
		info, err := model.AutoTest.GetByClusterIdAndProductName(clusterId, product.ProductName)
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]interface{}{
				"auto_test": true,
			}
		} else if err != nil {
			log.Errorf("err: %v", err.Error())
			return err
		}
		r := map[string]interface{}{}
		r["cluster_id"] = info.ClusterId
		r["product_name"] = info.ProductName
		r["auto_test"] = true
		r["exec_status"] = info.ExecStatus
		r["report_url"] = info.ReportUrl
		if info.CreateTime.Valid == true {
			r["create_time"] = info.CreateTime.Time.Format(base.TsLayout)
		} else {
			r["create_time"] = ""
		}
		if info.EndTime.Valid == true {
			r["end_time"] = info.EndTime.Time.Format(base.TsLayout)
		} else {
			r["end_time"] = ""
		}
		return r
	}
	return map[string]interface{}{
		"auto_test": false,
	}
}

//StartAutoTest 启动当前产品自动化测试
func StartAutoTest(ctx context.Context) apibase.Result {
	paramErrs := apibase.NewApiParameterErrors()
	var param = struct {
		ClusterId   int    `json:"cluster_id"`
		ProductName string `json:"product_name"`
	}{}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if param.ProductName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if param.ClusterId == 0 {
		cid, err := GetCurrentClusterId(ctx)
		if err != nil {
			paramErrs.AppendError("$", fmt.Errorf("get cluster id from cookie error"))
		}
		param.ClusterId = cid
	}
	paramErrs.CheckAndThrowApiParameterErrors()
	log.Debugf("autoTest cluster_id: %v,product_id: %v", param.ClusterId, param.ProductName)

	//获取deployment_product_list表中产品包信息
	info, err := model.DeployProductList.GetByProductNameClusterId(param.ProductName, param.ClusterId)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("product %s is not installed", param.ProductName)
	} else if err != nil {
		log.Errorf("err: %v", err.Error())
		return err
	}
	//反序列化，产品包默认schema信息
	originSchema, err := schema.Unmarshal(info.Schema)
	if err != nil {
		log.Errorf("err: %v", err.Error())
		return err
	}
	//获取deploy_schema_field_modify表中修改和新增的配置记录
	schemaFileds := make([]model.SchemaFieldModifyInfo, 0)
	query := "SELECT service_name, field_path, field FROM " + model.DeploySchemaFieldModify.TableName + " WHERE cluster_id=? AND product_name=? AND service_name=?"
	if err := model.USE_MYSQL_DB().Select(&schemaFileds, query, param.ClusterId, info.ProductName, TEST_SET_SERVICE_NAME); err != nil {
		log.Errorf("err: %v", err.Error())
		return err
	}

	if _, ok := originSchema.Service[TEST_SET_SERVICE_NAME]; ok && originSchema.Service[TEST_SET_SERVICE_NAME].Instance != nil {
		//deploy_schema_field_modify表中有修改和新增的配置记录
		newSchema, err := schema.Clone(originSchema)
		if err != nil {
			log.Errorf("[Switcher] Clone schema error: %v", err.Error())
			return err
		}
		for _, modify := range schemaFileds {
			newSchema.SetField(modify.ServiceName+"."+modify.FieldPath, modify.Field)
		}
		if err = inheritBaseService(param.ClusterId, newSchema, model.USE_MYSQL_DB()); err != nil {
			log.Errorf("[ProductInfo] inheritBaseService warn: %+v", err.Error())
		}
		if err = setSchemaFieldServiceAddr(param.ClusterId, newSchema, model.USE_MYSQL_DB(), ""); err != nil {
			log.Debugf("[ProductInfo] setSchemaFieldServiceAddr err: %v", err)
			return err
		}
		if err = newSchema.ParseVariable(); err != nil {
			log.Errorf("[ProductInfo] product info err: %v", err.Error())
			return err
		}
		if newSchema.Service[TEST_SET_SERVICE_NAME].Instance.TestOn && newSchema.Service[TEST_SET_SERVICE_NAME].Instance.TestScript != "" {
			dependsInfo := make([]model.DeployInstanceInfo, 0)
			testDepends := newSchema.Service[TEST_SET_SERVICE_NAME].Instance.TestDepends
			testScript := newSchema.Service[TEST_SET_SERVICE_NAME].Instance.TestScript
			//查询自动化测试产品包信息
			query := fmt.Sprintf("SELECT %s.* FROM %s LEFT JOIN %s ON %s.id = %s.pid WHERE cluster_id =? AND product_name =? AND %s.status =?",
				model.DeployInstanceList.TableName,
				model.DeployProductList.TableName,
				model.DeployInstanceList.TableName,
				model.DeployProductList.TableName,
				model.DeployInstanceList.TableName,
				model.DeployInstanceList.TableName,
			)
			if err := model.USE_MYSQL_DB().Select(&dependsInfo, query, strconv.Itoa(param.ClusterId), testDepends, model.INSTANCE_STATUS_RUNNING); err != nil {
				log.Errorf("err: %v", err.Error())
				return err
			}
			if len(dependsInfo) == 0 {
				return fmt.Errorf("depends %s is not installed", testDepends)
			}
			svc := dependsInfo[0]
			//写入操作数据
			operationId := uuid.NewV4().String()
			if err := model.AutoTest.Insert(model.AutoTestInfo{
				ClusterId:   param.ClusterId,
				ProductName: info.ProductName,
				OperationId: operationId,
				TestScript:  testScript,
				ExecStatus:  enums.ExecStatusType.Running.Code,
			}); err != nil {
				log.Errorf("err: %v", err.Error())
				return fmt.Errorf("operationInfo insert error: %v", err)
			}
			//执行自动化测试
			var err error
			var reportUrl string
			cmd := fmt.Sprintf("#!/bin/sh\n %s", testScript)
			timeOut := fmt.Sprintf("%dm", cache.SysConfig.GlobalConfig.AutoTestTimeoutLimit)
			content, err := agent.AgentClient.ToExecCmdWithTimeout(svc.Sid, svc.AgentId, strings.TrimSpace(cmd), timeOut, "", "")
			if err == nil {
				reg := regexp.MustCompile(`reportUrl:https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
				matches := reg.FindStringSubmatch(content)
				if matches != nil {
					reportUrl = strings.SplitN(matches[0], ":", 2)[1]
				} else {
					err = errors.New("unknown error")
				}
			}
			if err != nil {
				if err := model.AutoTest.UpdateStatusByOperationId(operationId, "", err.Error(), enums.ExecStatusType.Failed.Code, dbhelper.NullTime{Time: time.Now(), Valid: true}); err != nil {
					log.Errorf("err: %v", err.Error())
					return err
				}
				log.Errorf("execCmd failed, content: %s ;err: %v", content, err.Error())
				return err
			}
			log.Debugf("preview response: %v", content)
			//更新操作数据
			if err := model.AutoTest.UpdateStatusByOperationId(operationId, reportUrl, "", enums.ExecStatusType.Success.Code, dbhelper.NullTime{Time: time.Now(), Valid: true}); err != nil {
				log.Errorf("err: %v", err.Error())
				return err
			}

			return map[string]interface{}{
				"cluster_id":   param.ClusterId,
				"product_name": info.ProductName,
				"exec_status":  enums.ExecStatusType.Success.Code,
				"report_url":   reportUrl,
			}
		}
	}

	return map[string]interface{}{
		"auto_test": false,
	}
}

type checkError struct {
	ClusterId      int      `json:"cluster_id"`
	ConnectErr     []string `json:"connect_error"`
	PermissionsErr []string `json:"permissions_error"`
}

func checkPathPermissions(hosts []model.HostInfo, path string) checkError {
	script := fmt.Sprintf(`#!/bin/env bash
backPath=%s
who=$(whoami)

if sudo [ -f ${backPath} ]
then
  echo -n "file exist"
  exit 0
fi

if sudo [ -d ${backPath} ];then
  testDir=${backPath}/32656d2b-7821-4d10-976b-5e68190ce3cc
  mkdir -p ${testDir} 2>/dev/null
  if [[ $? -eq 0 ]] && [[ -d ${testDir} ]];then
      rm -fr ${testDir}
  else
      echo -n "permission denid ${backPath}"
      exit 0
  fi
else
  mkdir -p ${backPath} 2>/dev/null
  if [ $? -ne 0 ];then 
      echo -n "permission denid ${backPath}"
      exit 0
  fi
fi
`, path)
	var (
		wg       = sync.WaitGroup{}
		duration = "20s"
	)
	wg.Add(len(hosts))
	errs := checkError{
		ConnectErr:     []string{},
		PermissionsErr: []string{},
	}
	for _, host := range hosts {
		go func(host model.HostInfo) {
			defer wg.Done()
			content, err := agent.AgentClient.ToExecCmdWithTimeout(host.SidecarId, "", script, duration, "", "")
			if err != nil {
				log.Errorf("%s checkPathPermissions ToExecCmdWithTimeout %v", host.Ip, err)
				errs.ConnectErr = append(errs.ConnectErr, host.Ip)
				return
			}
			if strings.Contains(content, "file exist") {
				log.Errorf("%s checkPathPermissions %s is file", host.Ip, path)
				errs.PermissionsErr = append(errs.PermissionsErr, host.Ip)
				return
			}
			if strings.Contains(content, "permission denid") {
				log.Errorf("%s checkPathPermissions %s permission error", host.Ip, path)
				errs.PermissionsErr = append(errs.PermissionsErr, host.Ip)
				return
			}
		}(host)
	}
	wg.Wait()
	return errs
}

// SetClusterBackupPATH
// @Description  	set cluster backup path
// @Summary      	设置备份路径
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           message body string true "[{"clusterId":"1","path":"/opt/dtstack"}]"
// @Success         200 {object} string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/backup/setconfig [post]
func SetClusterBackupPATH(ctx context.Context) (rlt apibase.Result) {
	var requestParam []struct {
		ClusterId int    `json:"clusterId"`
		Path      string `json:"path"`
	}

	if err := ctx.ReadJSON(&requestParam); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}

	checkParamData := func(id int, path string) error {
		info, err := model.DeployClusterList.GetClusterInfoById(id)
		if err != nil {
			return fmt.Errorf("cluster %d not exist", id)
		}
		if info.Type != "hosts" {
			return fmt.Errorf("cluster %d type not hosts", id)
		}
		if !filepath.IsAbs(path) || strings.ContainsAny(path, " #$%^&*()+}{\":?><⌘;!?*=[]°§~@\\") {
			return fmt.Errorf("directory %s is illegal", path)
		}
		return nil
	}

	for _, param := range requestParam {
		if err := checkParamData(param.ClusterId, param.Path); err != nil {
			return err
		}
	}

	checkErrs := []checkError{}
	clusterWg := sync.WaitGroup{}
	clusterWg.Add(len(requestParam))
	for _, param := range requestParam {
		clusterId := param.ClusterId
		path := param.Path
		go func() {
			defer clusterWg.Done()
			hosts := model.DeployHostList.GetRunHostListByClusterId(clusterId)
			checkErr := checkPathPermissions(hosts, path)
			if len(checkErr.ConnectErr) != 0 || len(checkErr.PermissionsErr) != 0 {
				checkErr.ClusterId = clusterId
				checkErrs = append(checkErrs, checkErr)
			}
		}()
	}
	clusterWg.Wait()
	if len(checkErrs) != 0 {
		return checkErrs
	}

	for _, m := range requestParam {
		if err := model.ClusterBackupConfig.SaveClusterBackupConfig(m.ClusterId, m.Path); err != nil {
			return err
		}
	}

	return nil
}

// GetClusterBackupPATH
// @Description  	get cluster backup path
// @Summary      	查看备份路径
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           message body string true "[{"clusterId":"1","path":"/opt/dtstack"}]"
// @Success         200 {object} string "{"msg":"ok","code":0,"data":{"count":4,"data":{"clusterId":"","clusterName":"","path":""}}}"
// @Router          /api/v2/product/backup/getconfig [post]
func GetClusterBackupPATH(ctx context.Context) (rlt apibase.Result) {
	var (
		data []map[string]interface{}
	)
	hs, _ := model.DeployClusterList.SelectHostClusterList()
	for _, h := range hs {
		m := map[string]interface{}{
			"clusterId":   h.Id,
			"clusterName": h.Name,
			"path":        model.ClusterBackupConfig.GetPathByClusterId(h.Id),
		}
		data = append(data, m)
	}

	response := map[string]interface{}{
		"count": len(data),
		"data":  data,
	}
	return response
}

func ServiceConfigDiff(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ServiceConfigDiff] ServiceConfigDiff from EasyMatrix API ")

	var before, after interface{}
	productName := ctx.Params().Get("product_name")
	productVersion := ctx.Params().Get("product_version")
	serviceName := ctx.Params().Get("service_name")
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	ip := ctx.URLParam("ip")
	file := ctx.URLParam("file")
	if productName == "" {
		return fmt.Errorf("product_name is empty")
	}
	if productVersion == "" {
		return fmt.Errorf("product_version is empty")
	}
	if serviceName == "" {
		return fmt.Errorf("service_name is empty")
	}
	if ip == "" {
		return fmt.Errorf("ip is empty")
	}
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("Get cluster id error: %v", err)
		return fmt.Errorf("Get cluster id error: %v", err)
	}
	if file == "" {
		return map[string]interface{}{
			"before": before,
			"after":  after,
		}
	}
	//获取节点上当前已下发的配置文件信息
	instanceInfo := model.InstanceAndProductInfo{}
	query := "SELECT IL.*,PL.product_name, PL.product_name_display, PL.product_version FROM " +
		model.DeployInstanceList.TableName + " AS IL LEFT JOIN " + model.DeployProductList.TableName +
		" AS PL ON IL.pid = PL.id " + "WHERE PL.product_name =? AND IL.service_name =? AND IL.cluster_id =? AND ip = ?"
	if err = model.USE_MYSQL_DB().Get(&instanceInfo, query, productName, serviceName, clusterId, ip); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("[Product->ServiceConfigDiff] get instance error: %v", err)
		return err
	} else if err == nil {
		content, err := model.DeployInstanceList.GetInstanceServiceConfig(instanceInfo.ID, file)
		if err != nil && !strings.Contains(err.Error(), "No such file or directory") {
			log.Errorf("[Product->ServiceConfigDiff] read service file err: %v", err)
			return err
		}
		if err == nil {
			before = content
		}
	}
	//获取预下发的配置文件信息
	productInfo, err := model.DeployProductList.GetByProductNameAndVersion(productName, productVersion)
	if err != nil {
		log.Errorf("[Product->ServiceConfigDiff] get by product name and version error: %v", err)
		return err
	}
	sc, err := schema.Unmarshal(productInfo.Product)
	if err != nil {
		log.Errorf("[Product->ServiceConfigDiff] unmarshal err: %v", err)
		return err
	}
	if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
		log.Debugf("[Product->ServiceConfigDiff] inheritBaseService warn: %+v", err)
	}
	if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
		log.Debugf("[Product->ServiceConfigDiff] setSchemaFieldServiceAddr err: %v", err)
		return err
	}
	var node *model.ServiceIpNode
	node, err = model.GetServiceIpNode(clusterId, sc.ProductName, serviceName, ip)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("%v", err)
		return err
	}
	if err == nil {
		var ipList string
		query := "SELECT ip_list FROM " + model.DeployServiceIpList.TableName + " WHERE product_name=? AND service_name=? AND cluster_id=?"
		if err := model.DeployServiceIpList.GetDB().Get(&ipList, query, productName, serviceName, clusterId); err != nil {
			log.Errorf("%v", err)
			return err
		}
		serviceIp := strings.Split(ipList, IP_LIST_SEP)
		idToIndex, err := getIdToIndex(clusterId, productName, serviceName, serviceIp)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		if err = sc.SetServiceNodeIP(serviceName, util.FoundIpIdx(serviceIp, ip), node.NodeId, idToIndex); err != nil {
			log.Errorf("%v", err)
			return err
		}
	}
	if err = sc.ParseVariable(); err != nil {
		log.Errorf("[Product->ServiceConfigDiff] ParseVariable err: %v", err)
		return err
	}
	if err = handleUncheckedServices(sc, productInfo.ID, clusterId, namespace); err != nil {
		log.Errorf("[Product->ServiceConfigDiff] handleUncheckedServices warn: %+v", err)
	}
	multiFields, err := model.SchemaMultiField.GetByProductNameAndServiceNameAndIp(clusterId, sc.ProductName, serviceName, ip)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("[Product->ServiceConfigDiff] get schema multi field err: %v", err)
	}
	for _, multiField := range multiFields {
		sc.SetField(multiField.ServiceName+"."+multiField.FieldPath, multiField.Field)
	}
	if err = WithIpRoleInfo(clusterId, sc); err != nil {
		log.Debugf("[Product->ServiceConfigDiff] WithIpRoleInfo err: %v", err)
		return err
	}
	if sc.Service[serviceName].Instance == nil || sc.Service[serviceName].Instance.UseCloud {
		return errors.New("schema instance invalid")
	}
	if !sc.Service[serviceName].Instance.EmptyCar && len(sc.Service[serviceName].Instance.ConfigPaths) > 0 {
		baseDir := filepath.Join(base.WebRoot, sc.ProductName, sc.ProductVersion)
		cfgContents, err := sc.ParseServiceConfigFiles(baseDir, serviceName)
		if err != nil {
			msg := fmt.Sprintf("parse service config err: %v, service: %v", err, serviceName)
			return msg
		}
		var index int
		for i, path := range sc.Service[serviceName].Instance.ConfigPaths {
			if path == file {
				index = i
				break
			}
		}
		after = string(cfgContents[index])
	}

	return map[string]interface{}{
		"before": before,
		"after":  after,
	}
}

type ProductNameStruct struct {
	ID          int                       `json:"id"`
	Name        string                    `json:"product_name"`
	Version     string                    `json:"product_version"`
	Product     *ProductNameProductStruct `json:"product"`
	ProductType int                       `json:"product_type"`
}

type ProductNameProductStruct struct {
	ParentProductName string
}

func GetProductNameList(ctx context.Context) apibase.Result {
	type res struct {
		List []ProductNameStruct `json:"list"`
	}
	var result res
	var deployProductNames []string
	var deployStatus []string
	parentProductName := ctx.URLParam("parentProductName")
	if productNames := ctx.URLParam("productName"); productNames != "" {
		deployProductNames = strings.Split(productNames, ",")
	}
	clusterId := ctx.URLParam("clusterId")
	//namespace从cookie中获取，若mode不为空代表主机模式，将namespace置空
	namespace := ctx.GetCookie(COOKIE_CURRENT_K8S_NAMESPACE)
	// 在获取已部署应用列表的时候置空cookie中的namespace
	mode := ctx.URLParam("mode")
	if mode != "" {
		namespace = ""
	}
	//若前端未传clusterId,从cookie中获取实际clusterId
	var cid int
	var err error
	if clusterId == "" {
		cid, err = GetCurrentClusterId(ctx)
		if err != nil {
			log.Errorf("[GetProductNameList] get cluster id from cookie error: %s", err)
			//return fmt.Errorf("[ProductInfo] get cluster id from cookie error: %s", err)
		}
	} else {
		cid, _ = strconv.Atoi(clusterId)
	}
	ProductListInfo := make([]model.DeployProductListInfoWithNamespace, 0)
	// 集群id为0时表示，获取上传的产品包信息，非0表示获取指定集群下已部署过的产品包的信息
	if cid == 0 {
		// 安装包管理下的所有上传的产品包信息，同时根据前端检索条件进行检索
		ProductListInfo, _ = model.DeployProductList.GetProductListInfo(parentProductName, "", "", "", "", cid, deployStatus, deployProductNames, namespace)
	} else {
		// 获取集群下所有部署的产品包信息，同时根据前端检索条件进行检索
		ProductListInfo, err = model.DeployClusterProductRel.GetDeployClusterProductList(parentProductName, "", "", "", "", cid, deployStatus, deployProductNames, namespace)
		if err != nil {
			return fmt.Errorf("ProductInfo query error %s", err)
		}
	}

	//listMap := make(map[string]bool, 0)

	// 默认以服务名排序
	sort.SliceStable(ProductListInfo, func(i, j int) bool {
		return strings.Compare(ProductListInfo[i].ProductName, ProductListInfo[j].ProductName) == -1
	})

	for _, v := range ProductListInfo {
		temp := ProductNameStruct{}
		temp.ID = v.ID
		temp.Name = v.ProductName
		temp.Version = v.ProductVersion
		sc, err := schema.Unmarshal(v.Product)
		if err != nil {
			log.Errorf("[GetProductNameList-ProductInfo] Unmarshal err: %v", err)
		}
		temp.Product = &ProductNameProductStruct{
			ParentProductName: sc.ParentProductName,
		}
		temp.ProductType = v.ProductType
		result.List = append(result.List, temp)

	}

	return result
}

func GetServiceConfig(ctx context.Context) apibase.Result {
	log.Debugf("[Product->ServiceConfig] ServiceConfig from EasyMatrix API ")

	// 获取入参以及判断入参是否为空
	paramErrs := apibase.NewApiParameterErrors()
	productName := ctx.URLParam("product_name")
	pid, err := ctx.URLParamInt("pid")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("pid is not int"))
	}
	clusterId, err := ctx.URLParamInt("clusterId")
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("clusterId is not int"))
	}
	configPath := ctx.URLParam("configPath")

	if productName == "" {
		paramErrs.AppendError("$", fmt.Errorf("product_name is empty"))
	}
	if configPath == "" {
		paramErrs.AppendError("$", fmt.Errorf("configKey is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	log.Infof("GetServiceConfig, product name %v, clusterId %v, configPath %v", productName, clusterId, configPath)

	var info *model.DeployProductListInfo

	info, err = model.DeployProductList.GetProductInfoById(pid)
	if err != nil {
		return fmt.Errorf("Database query error %v", err)
	}
	sc, err := schema.Unmarshal(info.Product)
	if err != nil {
		log.Errorf("Unmarshal err: %v", err)
		return err
	}
	if err = inheritBaseService(clusterId, sc, model.USE_MYSQL_DB()); err != nil {
		log.Errorf("inheritBaseService warn: %+v", err)
	}
	if err = setSchemaFieldServiceAddr(clusterId, sc, model.USE_MYSQL_DB(), ""); err != nil {
		log.Errorf("setSchemaFieldServiceAddr err: %v", err)
		return err
	}
	if err = handleUncheckedServices(sc, info.ID, clusterId, ""); err != nil {
		log.Errorf("handleUncheckedServices warn: %+v", err)
	}
	if err = sc.ParseVariable(); err != nil {
		log.Errorf("product info err: %v", err)
		return err
	}
	//PublicService@Config.mysql_db
	configGroup := strings.Split(configPath, "@")
	if len(configGroup) < 2 {
		log.Errorf("configKey format err: %v", configPath)
		return fmt.Errorf("configKey format err: %v", configPath)
	}
	serviceName := configGroup[0]
	serviceConfig := configGroup[1]
	config := strings.Split(serviceConfig, ".")

	if len(config) < 2 || config[0] != "Config" {
		log.Errorf("configKey format err: %v", configPath)
		return fmt.Errorf("configKey format err: %v", configPath)
	}
	configKey := strings.Join(config[1:], ".")
	var configValue interface{}
	var service schema.ServiceConfig
	if _, ok := sc.Service[serviceName]; !ok {
		log.Errorf("service  not exist: %v", configKey)
		return fmt.Errorf("service  not exist: %v", configKey)
	}
	service = sc.Service[serviceName]
	if _, ok := service.Config[configKey]; !ok {
		log.Errorf("configKey not exist: %v", configKey)
		return fmt.Errorf("configKey not exist: %v", configKey)
	}
	switch service.Config[configKey].(schema.VisualConfig).Default.(type) {
	case *string:
		configValue = service.Config[configKey].(schema.VisualConfig).Default.(*string)
	case *schema.ServiceAddrStruct:
		configValue = service.Config[configKey].(schema.VisualConfig).Default.(*schema.ServiceAddrStruct).IP[0]
	}
	var infoList []model.SchemaFieldModifyInfo
	query := "SELECT service_name, field_path, field FROM " + model.DeploySchemaFieldModify.TableName + " WHERE product_name=? AND cluster_id=? AND namespace=?"
	if err := model.USE_MYSQL_DB().Select(&infoList, query, sc.ProductName, clusterId, ""); err != nil {
		return fmt.Errorf("query deploySchemaFieldModify error: %s", err)
	}
	var modifyValue interface{}
	for _, modify := range infoList {
		if modify.FieldPath == serviceConfig+".Value" {
			modifyValue = modify.Field
		}
	}
	log.Infof("configValue :%v, modifyValue %v", configValue, modifyValue)
	return map[string]interface{}{
		"pid":         pid,
		"clusterId":   clusterId,
		"productName": productName,
		"serviceName": serviceName,
		"configKey":   configKey,
		"configValue": configValue,
		"modifyvalue": modifyValue,
	}
}

// CheckDeployCondition
// @Description  	Check Deploy Condition
// @Summary      	检查当前组件的部署条件
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           cluster_id query  int  false  "集群id"
// @Param           auto_deploy query  bool  false  "自动部署"
// @Param           product_name query  string  false  "产品名称"
// @Param           product_line_name query  string  false  "产品线名称"
// @Param          product_line_version query  string  false  "产品线版本"
// @Param           product_type query  int  false  "产品类型"
// @Success         200  {object} string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/deployCondition [post]
func CheckDeployCondition(ctx context.Context) apibase.Result {
	log.Debugf("[Product->CheckDeployCondition] CheckDeployCondition from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	var param struct {
		ClusterId          int    `json:"cluster_id"`
		AutoDeploy         bool   `json:"auto_deploy"`
		ProductName        string `json:"product_name,omitempty"`
		ProductLineName    string `json:"product_line_name,omitempty"`
		ProductLineVersion string `json:"product_line_version,omitempty"`
		ProductType        int    `json:"product_type"`
	}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if param.ClusterId == 0 {
		paramErrs.AppendError("$", fmt.Errorf("cluster_id is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	if param.AutoDeploy {
		//自动部署，检查当前组件（包括其依赖组件）是否存在平滑升级版本、产品线是否为空、是否缺失组件包
		if param.ProductLineName != "" && param.ProductLineVersion != "" {
			info, err := model.DeployProductLineList.GetProductLineListByNameAndVersion(param.ProductLineName, param.ProductLineVersion)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Errorf("[Product->CheckDeployCondition] get product line err: %v", err)
				return err
			}
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("产品线 `%v(%v)` 不存在", param.ProductLineName, param.ProductLineVersion)
			}
			serials := make([]model.ProductSerial, 0)
			if err := json.Unmarshal(info.ProductSerial, &serials); err != nil {
				log.Errorf("[Product->CheckDeployCondition] json unmarshal error: %v", err)
				return err
			}
			productList := make([]string, 0)
			productRelList := make([]string, 0)
			for _, serial := range serials {
				smoothUpgradeProductRel, err := model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterId(serial.ProductName, param.ClusterId)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("query smoothUpgradeProductRel error: %v", err)
				} else if err == nil {
					productRelList = append(productRelList, smoothUpgradeProductRel.ProductName)
				}
			}
			if len(productRelList) > 0 {
				return fmt.Errorf("组件 `%v` 存在平滑升级版本，请勿操作部署。", strings.Join(productRelList, ","))
			}
			for _, serial := range serials {
				productInfoList, _ := model.DeployProductList.GetProductList(serial.ProductName, strconv.Itoa(param.ProductType), nil, nil)
				if len(productInfoList) == 0 {
					productList = append(productList, serial.ProductName)
				}
			}
			if len(productList) > 0 {
				return fmt.Errorf("缺失组件包`%v`，请先上传", strings.Join(productList, ","))
			}
		} else if param.ProductLineName == "" && param.ProductLineVersion == "" {
			return fmt.Errorf("请选择产品线")
		} else {
			return fmt.Errorf("product_line_name or product_line_version is empty")
		}
	} else {
		//手动部署，检查当前组件是否有平滑升级版本
		if param.ProductName == "" {
			return fmt.Errorf("product_name is empty")
		}
		products, _ := model.DeployProductList.GetProductListByNameAndType(param.ProductName, strconv.Itoa(param.ProductType), nil)
		if len(products) == 0 {
			log.Errorf("not found product %v", param.ProductName)
			return fmt.Errorf("组件 `%v` 不存在", param.ProductName)
		}
		smoothUpgradeProductRel, err := model.DeployClusterSmoothUpgradeProductRel.GetCurrentProductByProductNameClusterId(param.ProductName, param.ClusterId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("query smoothUpgradeProductRel error: %s", err)
		} else if err == nil {
			return fmt.Errorf("组件 `%v` 存在平滑升级版本，请勿操作部署。", smoothUpgradeProductRel.ProductName)
		}
	}

	return nil
}
