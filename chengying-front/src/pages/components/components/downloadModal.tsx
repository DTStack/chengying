import * as React from 'react';
import { Button, Form, Modal, Select } from 'antd';
import { formItemBaseLayout } from '@/constants/formLayout';
import { FormComponentProps } from 'antd/lib/form';
const FormItem = Form.Item;
const { Option } = Select;

interface IProps extends FormComponentProps {
  visible: boolean;
  onOk?: (e?: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  onCancel?: (e?: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  defaultValue?: {
    clusterId: number;
    parentProductName?: string;
  };
  data: any[];
}
const DownLoadModal: React.FC<IProps> = (props) => {
  const { visible, defaultValue, data, onCancel } = props;
  const { getFieldDecorator, setFieldsValue, validateFields } = props.form;
  const [clusterId, setClusterId] = React.useState<number>(
    defaultValue?.clusterId
  );
  const [productList, setProductList] = React.useState<any[]>([]);

  React.useEffect(() => {
    const cluster = data.find((item) => item.clusterId === clusterId) || {};
    const cProducts = cluster.subdomain?.products || [];
    const defaultProductName = defaultValue?.parentProductName || 'DTinsight';
    const product =
      cProducts.find((item: string) => item === defaultProductName) ||
      cProducts[0];
    setProductList(cProducts);
    setFieldsValue({
      parentProductName: product?.productName,
    });
  }, [clusterId]);

  // 集群变更
  function handleClusterChange(value: number) {
    setClusterId(value);
  }

  // 确定下载
  function handleOk() {
    validateFields((err: any, values: any) => {
      if (!err) {
        const { clusterId, parentProductName } = values;
        // 创建隐藏的可下载链接
        const link = document.createElement('a');
        link.style.display = 'none';
        link.href = `/api/v2/cluster/productsInfo?clusterId=${clusterId}&parentProductName=${parentProductName}`;
        link.download = parentProductName + '_' + clusterId;
        // 点击触发
        document.body.appendChild(link);
        link.click();
        // 移除
        document.body.removeChild(link);
        onCancel();
      }
    });
  }
  return (
    <Modal
      title="下载产品部署内容"
      visible={visible}
      footer={[
        <Button key="cancel" onClick={onCancel}>
          取消
        </Button>,
        <Button key="download" type="primary" onClick={handleOk}>
          下载
        </Button>,
      ]}
      onCancel={onCancel}>
      <Form>
        <FormItem {...formItemBaseLayout} label="集群">
          {getFieldDecorator('clusterId', {
            initialValue: defaultValue?.clusterId,
            rules: [{ required: true, message: '请选择集群' }],
          })(
            <Select placeholder="请选择集群" onChange={handleClusterChange}>
              {Array.isArray(data) &&
                data.map((item: any) => (
                  <Option key={item.clusterId} value={item.clusterId}>
                    {item.clusterName}
                  </Option>
                ))}
            </Select>
          )}
        </FormItem>
        <FormItem {...formItemBaseLayout} label="产品">
          {getFieldDecorator('parentProductName', {
            rules: [{ required: true, message: '请选择产品' }],
          })(
            <Select placeholder="请选择产品">
              {Array.isArray(productList) &&
                productList.map((item: string, index: number) => (
                  <Option key={item + index} value={item}>
                    {item}
                  </Option>
                ))}
            </Select>
          )}
        </FormItem>
      </Form>
    </Modal>
  );
};
export default Form.create<IProps>()(DownLoadModal);
