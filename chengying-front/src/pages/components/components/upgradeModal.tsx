import React, { useMemo, useState, useEffect } from 'react';
import { Modal, Form, Select, Alert, Button, Icon, message, Radio } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { formItemBaseLayout } from '@/constants/formLayout';
import { deployService } from '@/services';
import CommandPopbox from '@/pages/command/details/popbox';



const FormItem = Form.Item;
const { Option } = Select;

interface IProps extends FormComponentProps {
  visible: boolean;
  type: string;
  onOk?: (url: string, e?: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  onCancel: (e?: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  changeUpgradeType: (type: string) => void;
  options: any[];
  clusterId: string | number;
  record: any;
  isFirstSmooth: boolean;
}

const UpgradeModal: React.FC<IProps> = (props) => {
  const { visible, type, onOk, onCancel, options, clusterId, record, form, changeUpgradeType } =
    props;
  const { getFieldDecorator, validateFields } = form;
  const [backUpInfo, setBackUpInfo] = useState({
    status: '',
    exec_id: '',
    backup_name: '',
  });
  const [visibleInfo, changePopboxInfo] = useState({
    visible: false,
    title: '查看日志',
    type: 'log',
    execId: '',
  });
  const [loading, setLoading] = useState<boolean>(false)
  const [versions, setVersionList] = useState<string[]>([]);
  const [backupList, setBackUpList] = useState<string[]>([]);
  const [backUpable, setBackUpableStatus] = useState(false);
  const [defaultRadio, setDefaultRadio] = useState(record?.smooth_upgrade_product ? 'smooth' : 'normal');
  const [title, tips] = useMemo(() => {
    const title = type === 'upgrade' ? '升级' : '回滚';
    const tips: string =
      type === 'upgrade'
        ? '请先备份库，再执行升级部署，部署完成后将自动执行版本间的增量SQL；开始备份后，请勿退出页面，否则将取消本次升级。'
        : '仅支持回滚至 “升级” 操作前的版本，请先备份库，再执行回滚；开始备份后，请勿退出页面，否则将取消本次回滚。';
    return [title, tips];
  }, []);

  useEffect(() => {
    if (type !== 'upgrade') {
      getRollBackList();
    }
  }, []);

  /**
   * 备份
   */
  async function handleBackUp() {
    const param = {
      cluster_id: clusterId,
      source_version: record?.product_version,
      target_version: form.getFieldValue('target_version'),
    };
    if (backUpInfo.status !== 'running') {
      setBackUpInfo((preValue) => ({ ...preValue, status: 'running' }));
      setBackUpableStatus(false);
    }
    const res = await deployService.handleBackUp(
      {
        productName: record?.product_name,
      },
      param
    );
    if (res.data.code === 0) {
      const { data } = res.data;
      // 备份成功，保存备份信息到sessionStorage中,提供升级使用
      if (data.status !== 'running') {
        setBackUpInfo(data);
        setBackUpableStatus(true);
        if (type === 'upgrade' && data.status === 'success') {
          sessionStorage.setItem(
            'product_backup_info',
            JSON.stringify({
              ...param,
              backup_name: data.backup_name,
              backup_sqls: data.backup_sqls,
            })
          );
        }
      }
    } else {
      setBackUpInfo((preValue) => ({ ...preValue, status: 'fail' }));
      setBackUpableStatus(true);
      message.error(res.data.msg);
    }
  }

  function handleStatus (e) {
    if (e) {
      if (record?.smooth_upgrade_product || record?.isChildren) {
        setBackUpInfo((preValue) => ({ ...preValue, status: 'success' }))
      }
    }
  }

  async function getRollBackList() {
    let params = {
      cluster_id: clusterId,
      product_version: record?.product_version,
    }
    if (record?.smooth_upgrade_product || record?.isChildren) {
      Object.assign(params, {upgrade_mode: 'smooth'})
    }
    const { data } = await deployService.getRollBackList( {productName: record?.product_name}, params);
    if (data.code === 0) {
      setVersionList(data.data || []);
    }
  }

  async function getBackupTimes(target_version) {
    let params = {
      cluster_id: clusterId,
      target_version,
    }
    if (record?.smooth_upgrade_product || record?.isChildren) {
      Object.assign(params, {upgrade_mode: 'smooth'})
    }
    const { data } = await deployService.getBackupTimes(
      {
        productName: record?.product_name,
      },
      params,
    );
    if (data.code === 0) {
      setBackUpList(data.data || []);
    }
  }

  const handleEvent = () => {
    changePopboxInfo((prevValue) => ({
      visible: !prevValue.visible,
      title: '日志查看',
      type: 'log',
      execId: backUpInfo.exec_id,
    }));
  };

  // 改变升级状态
  const changeRadioStatus = (e) => {
    changeUpgradeType(e.target.value)
    setDefaultRadio(e.target.value)
  }

  // 进入部署向导
  function handleOk() {
    validateFields((err: any, values: any) => {
      if (!err) {
        let param = values.target_version;
        // 回滚时参数
        if (type === 'rollback') {
          param = {
            cluster_id: clusterId,
            source_version: record.product_version,
            target_version: values.target_version,
            backup_name: values.backup_name,
          };
          setLoading(true)
        }
        onOk(param);
      }
    });
  }

  // 取消
  function handleCancel() {
    setLoading(false)
    onCancel();
  }

  function handleStatusAble(target_version) {
    setBackUpableStatus(true);
    getBackupTimes(target_version);
  }

  return (
    <>
      <Modal
        title={`组件${title}`}
        visible={visible}
        onCancel={handleCancel}
        className="product-modal"
        footer={[
          <Button key="cancel" onClick={handleCancel}>
            取消
          </Button>,
          <Button
            loading={loading}
            key="ok"
            type="primary"
            onClick={handleOk}
            disabled={backUpInfo?.status !== 'success'}>
            {title}
          </Button>,
        ]}>
        <Form {...formItemBaseLayout}>
          <Alert className="mb-20" type="info" showIcon message={tips} />
          {type === 'upgrade' && record?.can_smooth_upgrade &&
          <FormItem label="升级模式">
          <Radio.Group 
            disabled={record?.smooth_upgrade_product ? true : false}
            onChange={changeRadioStatus}
            name="radiogroup" 
            defaultValue={defaultRadio}>
            <Radio value={'normal'}>普通升级</Radio>
            <Radio value={'smooth'}>平滑升级</Radio>
          </Radio.Group>
          <div className='modeTips'>
            {defaultRadio === 'normal' ? '一键升级至目标版本，过程中会有短暂停服。' : '通过多次升级平滑过渡至目标版本，升级过程不停服。'}
          </div>
          </FormItem>}
          <FormItem label="目标组件版本">
            {getFieldDecorator('target_version', {
              rules: [{ required: true, message: '请选择目标组件版本' }],
            })(
              <Select
                onChange={(e) => handleStatusAble(e)}
                placeholder="请选择目标组件版本">
                {type === 'upgrade'
                  ? options.map((item: any) => (
                      <Option key={item.id} value={item.product_version}>
                        {item.product_version}
                      </Option>
                    ))
                  : versions.map((item: any) => (
                      <Option key={item} value={item}>
                        {item}
                      </Option>
                    ))}
              </Select>
            )}
          </FormItem>
          {type !== 'upgrade' && (
            <FormItem label="备份库还原">
              {getFieldDecorator('backup_name', {
                rules: [{ required: true, message: '备份库缺失' }],
              })(
                <Select placeholder="请选择备份时间" disabled={!backUpable} onChange={(e) => handleStatus(e)}>
                  {backupList.map((item: any) => (
                    <Option key={item} value={item}>
                      {item}
                    </Option>
                  ))}
                </Select>
              )}
            </FormItem>
          )}
          {(type === 'upgrade' || (!record?.smooth_upgrade_product && !record?.isChildren) && type !== 'upgrade') &&
          <FormItem label="备份当前库" style={{ marginBottom: 0 }}>
            <Button
              type="primary"
              onClick={handleBackUp}
              disabled={!backUpable}>
              开始备份
            </Button>
            <div className="tips">
              {backUpInfo?.status === 'running' && (
                <div className="backuping">
                  <Icon type="reload" spin /> 备份中，请勿退出
                </div>
              )}
              {backUpInfo?.status === 'success' && (
                <div className="backupsuccess">
                  <span>
                    <Icon type="check-circle" theme="filled" /> 备份成功
                  </span>
                  <span>{backUpInfo.backup_name}</span>
                </div>
              )}
              {backUpInfo?.status === 'fail' && (
                <div className="backupfail">
                  <span>
                    <Icon type="close-circle" theme="filled" /> 备份失败
                  </span>
                  <a onClick={(e) => handleEvent()}>查看日志</a>
                </div>
              )}
            </div>
          </FormItem>}
        </Form>
      </Modal>
      {visibleInfo.visible && (
        <CommandPopbox
          {...visibleInfo}
          showFooter={false}
          onColse={handleEvent}
        />
      )}
    </>
  );
};

export default Form.create<IProps>()(UpgradeModal);
