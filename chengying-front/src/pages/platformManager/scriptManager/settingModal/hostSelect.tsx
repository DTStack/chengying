import React, { useState, useEffect } from 'react';
import { AutoComplete, Transfer } from 'antd';
import { installGuideService, Service } from '@/services';
import './style.scss';
interface IProps {
  defaultHosts: any[];
  selectHosts: (hosts) => void;
}

const HostSelect: React.FC<IProps> = (props) => {
  const { defaultHosts, selectHosts } = props;
  const [clusterList, setClusterList] = useState([]);
  const [clusterId, setClusterId] = useState('');
  const [hostList, setHostList] = useState([]);
  const [targetKeys, setTargetKeys] = useState([]);

  // 获取集群列表
  const getClusterList = async () => {
    const list = [];
    const response = await installGuideService.getInstallClusterList({
      'sort-by': 'id',
      'sort-dir': 'desc',
      type: 'hosts',
      limit: 0,
      start: 0,
    });
    const res = response.data;
    res.data.counts > 0 &&
      res.data.clusters.forEach((item) => {
        list.push({ clusterId: item.id, value: item.name });
      });
    setClusterList(list);
  };

  // 获取host列表
  const getHosts = async () => {
    setHostList([]);
    const res = await Service.getClusterHostList(
      {
        limit: 0,
        start: 0,
        cluster_id: clusterId,
      },
      'hosts'
    );
    const currentHost = res?.data?.data?.hosts;
    currentHost?.map((item) => {
      item.key = item.id;
      return item;
    });
    if (defaultHosts.length > 0) {
      defaultHosts.forEach((item: any) => {
        currentHost.forEach((pitem: any) => {
          if (item.ip === pitem.ip) {
            if (!targetKeys.includes(pitem.id)) {
              targetKeys.push(pitem.id);
            }
          }
        });
      });
      selectHosts(targetKeys);
    }
    setHostList(currentHost);
  };

  useEffect(() => {
    getClusterList();
  }, []);

  useEffect(() => {
    if (clusterId) {
      getHosts();
    }
  }, [clusterId]);

  useEffect(() => {
    if (defaultHosts.length > 0) {
      setClusterId(defaultHosts[0]?.cluster_id);
    } else {
      setClusterId(clusterList[0]?.clusterId);
    }
  }, [defaultHosts, clusterList]);

  // 选中集群
  const selectCluster = (id) => {
    setClusterId(id);
  };

  // 选中host
  const handleChange = (nextTargetKeys) => {
    selectHosts(nextTargetKeys);
    setTargetKeys(nextTargetKeys);
  };

  return (
    <div className="hostSelect">
      <div className="settingNavTop">
        <div className="navTopLine"></div>
        <div className="navTopTxt">下发主机</div>
      </div>
      <div className="hostSelectContent">
        <div className="hostSelectContent-cluster">
          <div className="hostSelectContent-cluster-title">请选择集群：</div>
          {clusterList.length > 0 &&
            clusterList.map((item) => {
              return (
                <div
                  onClick={() => {
                    selectCluster(item.clusterId);
                  }}
                  className={`hostSelectContent-cluster-content ${
                    item.clusterId === clusterId ?? clusterList[0].clusterId
                      ? 'active'
                      : null
                  }`}
                  key={item.clusterId}>
                  {item.value}
                </div>
              );
            })}
        </div>
        <div className="hostSelectContent-host">
          <Transfer
            dataSource={hostList}
            showSearch
            listStyle={{
              width: 250,
              height: 300
            }}
            titles={['未选', '已选']}
            targetKeys={targetKeys}
            onChange={handleChange}
            render={(item) => item.ip}
          />
        </div>
      </div>
    </div>
  );
};
export default HostSelect;
