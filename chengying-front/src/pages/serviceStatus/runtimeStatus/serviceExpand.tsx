import * as React from 'react';
import { Form, InputNumber, Modal } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
const FormItem = Form.Item;

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 6 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16 },
  },
};

interface IProps extends FormComponentProps {
  visible: boolean;
  defaultValue?: number;
  onOk: Function;
  onCancel: (e: React.MouseEvent<HTMLElement, MouseEvent>) => void;
}

const PodModal: React.FC<IProps> = (props) => {
  const { visible, defaultValue, form } = props;
  const { getFieldDecorator, validateFields } = form;
  const [loading, setLoading] = React.useState<boolean>(false);

  // 校验数量
  function validatePodNum(rule: any, value: any, callback: any) {
    if (value || value === 0) {
      if (!/^-?[0-9]\d*$/.test(value)) {
        callback('请输入整数');
      }
      if (value < 1 || value > 1000) {
        callback('请输入1-1000内的数值');
      }
    }
    callback();
  }

  // 提交
  function handleSubmit() {
    validateFields((err: any, values: any) => {
      if (!err) {
        setLoading(true);
        props.onOk(values.replica);
      }
    });
  }

  return (
    <Modal
      title={'服务扩缩容'}
      visible={visible}
      confirmLoading={loading}
      onOk={handleSubmit}
      onCancel={props.onCancel}>
      <Form>
        <FormItem label="Replica" {...formItemLayout}>
          {getFieldDecorator('replica', {
            initialValue: defaultValue || undefined,
            rules: [
              { required: true, message: '请输入扩缩容节点数' },
              { validator: validatePodNum },
            ],
          })(
            <InputNumber
              placeholder="请输入扩缩容节点数"
              style={{ width: '100%' }}
            />
          )}
        </FormItem>
      </Form>
    </Modal>
  );
};
export default Form.create<IProps>()(PodModal);
