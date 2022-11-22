import React, { useState, useEffect } from 'react';
import { Input, Table, Form, Button, message, Modal, List } from 'antd';
import api from '@/services/backupConfigService';
import { FormComponentProps } from 'antd/es/form';
import './style.scss';

const BackupConfig = (props: FormComponentProps) => {
  const { setFieldsValue, getFieldDecorator, setFields } = props.form;
  const [dataSource, setDataSource] = useState([]);
  const [loading, setLoading] = useState(false)
  let [arr, setArr] = useState([]);
  let [hostInfo, setHostInfo] = useState({
    visible: false,
    info: {
      clusterName: '',
      cluster_id: '',
      connect_error: [],
      permissions_error: [],
    },
  });
  const columns = [
    {
      title: '集群名称',
      dataIndex: 'clusterName',
      key: 'clusterName',
    },
    {
      title: '应用备份升级路径',
      dataIndex: 'path',
      key: 'path',
      render: (text: string, record: any) => {
        return (
          <Form.Item>
            {getFieldDecorator(`${record.clusterId}`, {
              rules: [
                { required: true, message: '请输入路径!' },
                { validator: checkingAuthority },
              ],
            })(
              <Input
                style={{ width: '90%' }}
                placeholder="请输入备份升级路径！"
              />
            )}
          </Form.Item>
        );
      },
    },
  ];

  const getFieldValue = (dataSource: any[]) => {
    dataSource.map((item) => {
      setFieldsValue({
        [item.clusterId]: item.path,
      });
    });
  };

  const getBackupsList = async () => {
    let res = await api.queryBuildBackupPath({});
    res = res.data;
    if (res.code !== 0) {
      message.error(res.msg);
    } else {
      setDataSource(res.data.data);
      getFieldValue(res.data.data);
    }
  };

  const saveBackup = async (params: any[]) => {
    setArr([])
    setLoading(true)
    let res = await api.SetUpBackupPath(params);
    res = res.data;
    if (res.code !== 0) {
      setLoading(false)
      message.error(res.msg);
    } else if (res.data !== null) {
      await setArr(res.data);
      res.data.map((o) =>
        props.form.validateFields([`${o.cluster_id}`], { force: true })
      );
      setLoading(false)
    } else {
      setLoading(false)
      message.success('操作成功！');
    }
  };

  const handleSubmit = (e: any) => {
    let params = [];
    e.preventDefault();
    props.form.validateFields((err, values) => {
        for (let key in values) {
          params.push({ clusterId: +key, path: values[key] });
        }
        if (!err) {
          saveBackup(params);
        } else {
          // 判断错误类型, 如果不是非空错误，保存按钮可以调用
          for (let errTxt in err) {
            const ob = err[errTxt]?.errors[0]?.message
            setFields({[errTxt]: {value: values[errTxt], errors: null}})
            if (typeof ob != 'string') {
              saveBackup(params);
            }
          }
        }
    });
  };

  const checkingAuthority = (rule: any, value: any, callback: any) => {
    if (arr.length && arr.some((o) => o.cluster_id === +rule.field)) {
      const newArr = arr.filter((o) => o.cluster_id === +rule.field);
      const curretData = dataSource.filter((o) => o.clusterId === +rule.field);
      let errTxt = '';
      if (newArr[0]?.connect_error.length > 0) {
        if (newArr[0]?.permissions_error.length > 0) {
          errTxt = '存在主机agent连接失败，路径权限不足！';
        } else {
          errTxt = '存在主机agent连接失败！';
        }
      } else {
        if (newArr[0]?.permissions_error.length > 0) {
          errTxt = '路径权限不足！';
        }
      }
      const errorMsg = (
        <div>
          {' '}
          {errTxt}
          <span
            onClick={() => handleGetHost(newArr[0], curretData[0].clusterName)}
            style={{ color: '#3F87FF', marginLeft: 8, cursor: 'pointer' }}>
            查看主机
          </span>
        </div>
      );
      callback(errorMsg);
    } else {
      callback();
    }
  };

  useEffect(() => {
    getBackupsList();
  }, []);

  const handleGetHost = (hostError?: any, clusterName?: string) => {
    const infos = {
      visible: !hostInfo.visible,
      info: {
        clusterName,
        ...hostError,
      },
    };
    setHostInfo(infos);
  };

  return (
    <div className="backupConfig">
      <Form onSubmit={handleSubmit}>
        <Table
          rowKey="clusterId"
          dataSource={dataSource}
          columns={columns}
          pagination={false}
          style={{ marginBottom: 20 }}
        />
        <Button type="primary" htmlType="submit" loading={loading}>
          保存
        </Button>
      </Form>
      {hostInfo.visible && (
        <Modal
          className="error-list"
          title={`查看主机/${hostInfo.info?.clusterName}`}
          visible={hostInfo.visible}
          onCancel={handleGetHost}
          footer={null}>
          {hostInfo.info?.connect_error.length > 0 && (
            <List
              size="small"
              header={<div>以下主机agent连接失败</div>}
              dataSource={hostInfo.info?.connect_error}
              renderItem={(item) => <List.Item>{item}</List.Item>}
            />
          )}
          {hostInfo.info?.permissions_error.length > 0 && (
            <List
              size="small"
              header={<div>以下主机路径权限不足</div>}
              dataSource={hostInfo.info?.permissions_error}
              renderItem={(item) => <List.Item>{item}</List.Item>}
            />
          )}
        </Modal>
      )}
    </div>
  );
};

export default Form.create({ name: 'backup' })(BackupConfig);
