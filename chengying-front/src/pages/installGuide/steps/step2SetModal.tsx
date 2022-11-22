import * as React from 'react';
import { Input, Tooltip, Icon, Table} from 'antd';
import '../style.scss';

class Step1SetModal extends React.Component<any, any> {
  constructor(props: any) {
    super(props);
  }

  state = {
    testData: this.props.data.list
  }

  changeText = (value, record) => {
    const { testData } = this.state;
    const index = testData.findIndex(item => record.product_version === item.product_version)
    let arr = testData
    arr[index].field = value.target.value
    this.setState({testData: arr})
  }

  componentDidUpdate(prevProps: any) {
    if (prevProps.data?.list !== this.props.data?.list) {
      this.setState({ testData:  this.props.data.list});
    }
  }

  render() {
    const columns = [
      {
        title: '组件版本号',
        dataIndex: 'product_version',
        width: 160,
        ellipsis: true,
      },
      {
        title: (
            <Tooltip
            title={
              '应用启动用户名，将初始化进产品包的各服务中， 即产品包下各服务部署时，默认以此用户名启动。同时，各服务也支持使用不同的启动用户名， 故实际部署以最终细粒度的服务启动用户名为准。'
            }>
              应用启动用户名
              <Icon
                type="question-circle"
                style={{
                  fontSize: '16px',
                  marginLeft: '10px',
                }}
              />
          </Tooltip>
        ),
        dataIndex: 'field',
        width: 360,
        render: (text, record) => {
          return (
            <Input placeholder="请输入应用启动用户名" value={text} onChange={(val) => this.changeText(val, record)}/>
          )
        }
      },
    ]
    const { data } = this.props;
    return (
      <div className='step-set-modal'>
          <div className='step-set-title'>组件名称: {data.product_name}</div>
          <Table
            dataSource={this.state.testData}
            columns={columns}
            pagination={false}
          />
      </div>
    );
  }
}
export default Step1SetModal;
