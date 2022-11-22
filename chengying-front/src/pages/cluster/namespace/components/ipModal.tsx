import * as React from 'react';
import { Button, Modal, Form, Input, message, Alert } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { ClusterNamespaceService } from '@/services';
const FormItem = Form.Item;

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 6 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 15 },
  },
};

interface IProps extends FormComponentProps {
  namespace: string;
  visible: boolean;
  handleCancel: (e?: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  getTableList: Function;
}

const IpModal: React.FC<IProps> = (props) => {
  const { namespace, visible, handleCancel } = props;
  const { getFieldDecorator, validateFields } = props.form;
  const [isConnect, setIsConnect] = React.useState<boolean>(false);
  const [connectLoading, setConnectLoading] = React.useState<boolean>(false);
  const [saveLoading, setSaveLoding] = React.useState<boolean>(false);

  // 校验
  function validateCheck(callback) {
    validateFields((err: any, values: any) => {
      if (err) {
        return;
      }
      callback(values);
    });
  }

  // 测试连通性
  function handleConnect() {
    validateCheck((values: any) => {
      setConnectLoading(true);
      ClusterNamespaceService.pingConnect({
        ...values,
        namespace,
      }).then((response: any) => {
        const res = response.data;
        const { code, msg } = res;
        if (code === 0) {
          setIsConnect(true);
          message.success('测试连通性通过');
        } else {
          message.error(msg);
        }
        setConnectLoading(false);
      });
    });
  }

  // 确定绑定
  function handleOk() {
    validateCheck((values: any) => {
      setSaveLoding(true);
      ClusterNamespaceService.saveNamespace({
        ...values,
        namespace,
        type: 'agent',
      }).then((response: any) => {
        const res = response.data;
        const { code, msg } = res;
        if (code === 0) {
          message.success('绑定成功');
          props.getTableList();
          handleCancel();
        } else {
          message.error(msg);
        }
        setSaveLoding(false);
      });
    });
  }

  return (
    <Modal
      title="绑定IP"
      visible={visible}
      footer={
        <React.Fragment>
          <Button
            type="primary"
            ghost
            onClick={handleConnect}
            loading={connectLoading}>
            测试连通性
          </Button>
          <Button type="default" onClick={handleCancel}>
            取消
          </Button>
          <Button
            type="primary"
            disabled={!isConnect}
            onClick={handleOk}
            loading={saveLoading}>
            确定
          </Button>
        </React.Fragment>
      }
      bodyStyle={{ padding: '0 0 8px 0' }}
      onCancel={handleCancel}>
      <Form>
        <Alert
          className="mb-20"
          type="info"
          message="请联系集群管理员，暴露easyagent服务供集群外访问"
        />
        <FormItem label="IP" {...formItemLayout}>
          {getFieldDecorator('ip', {
            rules: [
              { required: true, message: '请输入IP地址' },
              {
                pattern:
                  /^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$/,
                message: '请输入正确的IP地址',
              },
            ],
          })(<Input placeholder="请输入IP地址" />)}
        </FormItem>
        <FormItem label="端口" {...formItemLayout}>
          {getFieldDecorator('port', {
            rules: [
              { required: true, message: '请输入端口' },
              {
                pattern:
                  /^([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
                message: '请输入正确的端口',
              },
            ],
          })(<Input placeholder="请输入端口" />)}
        </FormItem>
      </Form>
    </Modal>
  );
};
export default Form.create<IProps>()(IpModal);
