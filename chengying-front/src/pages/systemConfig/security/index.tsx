import React, { useState, useEffect } from 'react';
import {
  InputNumber,
  Form,
  Button,
  Select,
  Switch,
  message,
  Alert,
} from 'antd';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
import { securityService } from '@/services';
import { FormComponentProps } from 'antd/es/form';
import './style.scss';

const Option = Select.Option;
const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 4 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 12 },
  },
};

interface ISafeConfProps {
  login_encrypt: string;
  force_reset_password: number;
  account_login_lock_switch: number;
  account_login_limit_error: number;
  account_login_lock_time: number;
  account_logout_sleep_time: number;
}

interface ISecurityProps extends FormComponentProps {
  authorityList: Array<{ [key: string]: boolean }>;
}

const mapStateToProps = (state: AppStoreTypes) => ({
  authorityList: state.UserCenterStore.authorityList,
});

const SecurityConfig: React.FC<ISecurityProps> = (props) => {
  const { authorityList, form } = props;
  const { getFieldDecorator, validateFields } = form;
  const [safeConf, setSafeConf] = useState<ISafeConfProps>({
    login_encrypt: 'rsa',
    force_reset_password: 0,
    account_login_lock_switch: 0,
    account_login_limit_error: 3,
    account_login_lock_time: 5,
    account_logout_sleep_time: 1440,
  });
  const [hasLock, setLockStatus] = useState(true);
  const hasEdit: boolean =
    authorityList?.sub_menu_configuration_platformsecurity_config_edit || false;

  useEffect(() => {
    getSecurityConf();
  }, []);

  const getSecurityConf = async () => {
    const { data } = await securityService.getSecurity();
    if (data.code !== 0) {
      return message.error(data.msg);
    }
    setSafeConf((prevValue) => ({ ...prevValue, ...data.data }));
  };

  const saveSecurityConf = async (params) => {
    const { data } = await securityService.setSecurity(params);
    if (data.code !== 0) {
      return message.error(data.msg);
    }
    getSecurityConf();
    message.success('保存成功');
  };

  const handleSubmit = (e: any) => {
    e.preventDefault();
    validateFields((err, values) => {
      if (err) {
        return console.log(err);
      }
      values.force_reset_password = values.force_reset_password ? 1 : 0;
      values.account_login_lock_switch = values.account_login_lock_switch
        ? 1
        : 0;
      saveSecurityConf(values);
      setLockStatus(true);
    });
  };

  const validSleepRule = (rules, value, callback) => {
    const regx = /^[1-9]\d{0,3}$/;
    if ((value || value === 0) && !regx.test(value)) {
      callback('请输入正整数');
    }
    callback();
  };

  const validLockRule = (rules, value, callback) => {
    const regx = /^[1-9]\d{0,2}$/;
    if ((value || value === 0) && !regx.test(value)) {
      callback('请输入正整数');
    }
    callback();
  };

  return (
    <Form
      className="security-wrapper"
      {...formItemLayout}
      onSubmit={handleSubmit}>
      <Alert
        type="info"
        showIcon
        message="为保证密码传输安全性，平台提供以下传输加密方式供用户按照需要进行选择。"
      />
      <Form.Item label="密码传输加密">
        {getFieldDecorator('login_encrypt', {
          initialValue: safeConf.login_encrypt,
        })(
          <Select onChange={(e) => setLockStatus(false)}>
            <Option value="rsa">rsa加密</Option>
            <Option value="sm2">sm加密</Option>
          </Select>
        )}
      </Form.Item>
      <Alert
        type="info"
        showIcon
        message="初入平台由系统管理员为其开通账号，若勾选下方选项则要求用户在使用初始密码登录后进行密码强制修改，否则不予使用平台功能。"
      />
      <Form.Item label="强制用户修改初始密码">
        {getFieldDecorator('force_reset_password', {
          initialValue: safeConf.force_reset_password === 1,
          valuePropName: 'checked',
        })(
          <Switch
            onChange={(e) => {
              setSafeConf((preValue) => ({
                ...preValue,
                force_reset_password:
                  safeConf.force_reset_password === 0 ? 1 : 0,
              }));
              setLockStatus(false);
            }}
          />
        )}
      </Form.Item>
      <Alert
        type="info"
        showIcon
        message="为保证平台安全性，系统管理员可为所有用户设置登录时允许密码出错的次数，达到次数后该账户将被锁定，且支持设置锁定时长，达到时长后用户可继续尝试登录。"
      />
      <Form.Item label="密码出错锁定">
        {getFieldDecorator('account_login_lock_switch', {
          initialValue: safeConf.account_login_lock_switch === 1,
          valuePropName: 'checked',
        })(
          <Switch
            onChange={(e) => {
              setSafeConf((preValue) => ({
                ...preValue,
                account_login_lock_switch:
                  safeConf.account_login_lock_switch === 0 ? 1 : 0,
              }));
              setLockStatus(false);
            }}
          />
        )}
      </Form.Item>
      <Form.Item label="密码出错锁定">
        {getFieldDecorator('account_login_limit_error', {
          initialValue: safeConf.account_login_limit_error,
          rules: [
            {
              required: true && safeConf.account_login_lock_switch === 1,
              message: '请输入次数',
            },
            {
              validator: validLockRule,
            },
          ],
        })(
          <InputNumber
            maxLength={3}
            placeholder="请输入次数"
            disabled={safeConf.account_login_lock_switch === 0}
            onChange={(e) => setLockStatus(false)}
          />
        )}
        <span className="ant-form-text">次</span>
      </Form.Item>
      <Form.Item label="锁定时长">
        {getFieldDecorator('account_login_lock_time', {
          initialValue: safeConf.account_login_lock_time,
          rules: [
            {
              required: true && safeConf.account_login_lock_switch === 1,
              message: '请输入时长',
            },
            {
              validator: validLockRule,
            },
          ],
        })(
          <InputNumber
            maxLength={3}
            placeholder="请输入时长"
            disabled={safeConf.account_login_lock_switch === 0}
            onChange={(e) => setLockStatus(false)}
          />
        )}
        <span className="ant-form-text">分钟</span>
      </Form.Item>
      <Alert
        type="info"
        showIcon
        message="为保证平台安全性，系统管理员可针对 “用户登录后在页面停留无操作” 行为进行时长限制，达到时长后将对该账户采取自动登出措施，用户可再次登录访问。"
      />
      <Form.Item label="自动登出时长限制">
        {getFieldDecorator('account_logout_sleep_time', {
          initialValue: safeConf.account_logout_sleep_time,
          rules: [
            {
              required: true,
              message: '自动登出时长限制不可为空',
            },
            {
              validator: validSleepRule,
            },
          ],
        })(
          <InputNumber
            maxLength={4}
            placeholder="请输入时长"
            onChange={(e) => setLockStatus(false)}
          />
        )}
        <span className="ant-form-text">分钟</span>
      </Form.Item>
      <Form.Item className="footer-submit" label={<span></span>} colon={false}>
        <Button htmlType="submit" type="primary" disabled={hasLock && hasEdit}>
          保存
        </Button>
      </Form.Item>
    </Form>
  );
};

export default connect(
  mapStateToProps,
  undefined
)(Form.create()(SecurityConfig));
