package monitor

import (
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"time"
)

var (
	Duration              = 60
	STATUS_OK_HOSTS       = host.InitInitializeShOk
	STATUS_OK_KUBEZRNETES = host.InitInitializeShOk
)

type HostStatus struct {
	Status     int       `db:"status" json:"status"`
	Steps      int       `db:"steps" json:"steps"`
	UpdateDate base.Time `db:"updated" json:"updated"`
}

func StartClusterStatusM() {
	duration := time.Duration(Duration) * time.Second
	for {
		time.Sleep(duration)
		log.Infof("StartClusterStatusM ...")

		list, err := model.DeployClusterList.GetClusterList()

		if err != nil {
			log.Errorf("%v", err.Error())
			continue
		}
		hostPersist := "deploy_host"
		for _, cluster := range list {
			//获取主机集群状态
			query := "SELECT h.status,steps,updated " +
				"FROM deploy_cluster_host_rel " +
				"LEFT JOIN deploy_cluster_list ON deploy_cluster_list.id = deploy_cluster_host_rel.clusterId " +
				"LEFT JOIN %s as h ON deploy_cluster_host_rel.sid = h.sid " +
				"WHERE deploy_cluster_host_rel.is_deleted=0 and deploy_cluster_list.id = ? and deploy_cluster_list.is_deleted=0 and h.isDeleted=0 and h.status is not NULL"
			hostPersist = "deploy_host"
			query = fmt.Sprintf(query, hostPersist)
			hostList := make([]HostStatus, 0)
			if _ = model.USE_MYSQL_DB().Select(&hostList, query, cluster.Id); len(hostList) == 0 {
				model.DeployClusterList.UpdateClusterStatus(cluster.Id, model.DEPLOY_CLUSTER_STATUS_WAITING)
				continue
			}
			switch cluster.Type {
			case model.DEPLOY_CLUSTER_TYPE_HOSTS:
				processHostClusterStaus(cluster, hostList)
			case model.DEPLOY_CLUSTER_TYPE_KUBERNETES:
				processK8SClusterStaus(cluster, hostList)
			}
		}
	}
}

func processHostClusterStaus(cluster model.ClusterInfo, hostList []HostStatus) {
	var errCount, okCount int
	for _, h := range hostList {
		if h.Status < 0 {
			errCount++
		}
		if h.Status >= STATUS_OK_HOSTS {
			okCount++
		}
	}
	if errCount > 0 {
		model.DeployClusterList.UpdateClusterStatus(cluster.Id, model.DEPLOY_CLUSTER_STATUS_ERROR)
		return
	}
	if okCount > 0 {
		model.DeployClusterList.UpdateClusterStatus(cluster.Id, model.DEPLOY_CLUSTER_STATUS_RUNNING)
		return
	}
	model.DeployClusterList.UpdateClusterStatus(cluster.Id, model.DEPLOY_CLUSTER_STATUS_PENDING)
}

func processK8SClusterStaus(cluster model.ClusterInfo, hostList []HostStatus) {
	var errCount, okCount int
	for _, h := range hostList {
		if h.Status < 0 {
			errCount++
		}
		if h.Status >= STATUS_OK_KUBEZRNETES {
			okCount++
		}
	}
	if errCount > 0 {
		model.DeployClusterList.UpdateClusterStatus(cluster.Id, model.DEPLOY_CLUSTER_STATUS_ERROR)
		return
	}
	if okCount > 0 && okCount == len(hostList) {
		model.DeployClusterList.UpdateClusterStatus(cluster.Id, model.DEPLOY_CLUSTER_STATUS_RUNNING)
		return
	}
	model.DeployClusterList.UpdateClusterStatus(cluster.Id, model.DEPLOY_CLUSTER_STATUS_PENDING)
}
