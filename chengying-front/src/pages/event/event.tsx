import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators, Dispatch } from 'redux';
import { get } from 'lodash';
import {
  Layout,
  Input,
  Select,
  Form,
  Row,
  Col,
  Tabs,
  Radio,
  List,
  Avatar,
  Table,
  DatePicker,
  message,
} from 'antd';
import { AppStoreTypes } from '@/stores';
import { ServiceStore, HeaderStoreType } from '@/stores/modals';
import { FormComponentProps } from 'antd/lib/form/Form';
import { servicePageService, eventService } from '@/services';
import * as ServiceActions from '@/actions/serviceAction';
import Echarts from './echarts';
import moment from 'moment';
// import { component,eventList } from './mock'
import './style.scss';
const { TabPane } = Tabs;
const { Content } = Layout;
const { RangePicker } = DatePicker;
const Search = Input.Search;

const mapStateToProps = (state: AppStoreTypes) => ({
  HeaderStore: state.HeaderStore,
  ServiceStore: state.ServiceStore,
});
const mapDispatchToProps = (dispatch: Dispatch) => ({
  actions: bindActionCreators(Object.assign({}, ServiceActions), dispatch),
});

interface Prop {
  HeaderStore: HeaderStoreType;
  ServiceStore: ServiceStore;
  actions: any;
  match?: any;
  location?: any;
}
interface State {
  time: string;
  activeKey: string;
  eventType: any[];
  currentEventType: string;
  products: any[];
  serviceGroup: any[];
  hostGroup: any[];
  servicProduct: {};
  currentProduct: any[];
  currentService: any[];
  currentHost: any[];
  count: {
    count: number;
    product_count: number;
    service_count: number;
  };
  echartsTime: any[];
  echartsData: any[];
  componentData: any[];
  serviceData: any[];
  startDate: string;
  endDate: string;
  eventListData: any[];
  eventListColumns: any[];
  eventListCount: number;
  currentList: number;
  keyWord: string;
}
@(connect(mapStateToProps, mapDispatchToProps) as any)
class Event extends React.Component<Prop & FormComponentProps, State> {
  constructor(props: any) {
    super(props);
  }

  state: State = {
    time: 'day',
    activeKey: '1',
    eventType: [],
    currentEventType: '',
    products: [],
    serviceGroup: [],
    hostGroup: [],
    servicProduct: {},
    currentProduct: [],
    currentService: [],
    currentHost: [],
    count: {
      count: null,
      product_count: null,
      service_count: null,
    },
    echartsTime: [],
    echartsData: [],
    componentData: [],
    serviceData: [],
    startDate: '',
    endDate: '',
    eventListData: [],
    eventListColumns: [],
    eventListCount: 0,
    currentList: 1,
    keyWord: '',
  };

  componentDidMount = () => {
    const startDate = moment().startOf('day').subtract(30, 'd');
    const endDate = moment().endOf('day');
    this.setState({ startDate, endDate });
    const { cur_parent_product } = this.props.HeaderStore;
    this.handleGetEventType(cur_parent_product);
    this.getProductList(cur_parent_product);
  };

  // 获得事件类型
  handleGetEventType = (curParentProduct: any) => {
    const self = this;
    eventService
      .getEventType({ limit: '', parentProductName: curParentProduct })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          self.setState(
            {
              eventType: res.data.list,
              currentEventType: res.data.list[0],
            },
            () => {
              this.handleEventTypeChange(this.state.currentEventType);
            }
          );
        } else {
          message.error(res.msg);
          self.setState({ eventType: [] });
        }
      });
  };

  // 切换事件类型
  handleEventTypeChange = (e: any) => {
    const startDate = moment().startOf('day').subtract(30, 'd');
    const endDate = moment().endOf('day');
    this.props.form.setFieldsValue({
      eventTypes: e,
      component: [],
      service: [],
      host: [],
    });
    this.setState(
      {
        time: 'day',
        serviceGroup: [],
        hostGroup: [],
        currentProduct: [],
        currentHost: [],
        currentService: [],
        currentEventType: e,
        startDate,
        endDate,
        eventListColumns: [],
        eventListData: [],
        eventListCount: 0,
      },
      () => {
        this.handleChangeTime('day');
      }
    );
  };

  // 获得组件
  getProductList = (curParentProduct: any) => {
    const self = this;
    servicePageService
      .getProductName({ parentProductName: curParentProduct })
      .then((res: any) => {
        res = res.data;
        if (res.code === 0 && res.data?.list) {
          const list = res.data.list;
          self.setState({
            products: list,
          });
        } else {
          message.error(res.msg);
          this.setState({ products: [] });
        }
      });
  };

  // 切换组件
  handleSelectProduct = (item: any) => {
    const serviceGroup = [];
    const servicProduct = {};
    this.setState({
      serviceGroup: [],
      hostGroup: [],
      currentProduct: item,
      currentHost: [],
      currentService: [],
    });
    this.props.form.setFieldsValue({
      component: item,
      service: [],
      host: [],
    });
    item.map((value) => {
      servicePageService
        .getServiceGroup({ product_name: value })
        .then((res: any) => {
          res = res.data;
          if (res.code === 0) {
            if (res.data.count > 0) {
              const groups = res.data.groups;
              for (const i in groups) {
                serviceGroup.push({
                  name: i,
                  list: groups[i],
                });
                groups[i].map((a) => {
                  servicProduct[a.service_name] = value;
                });
              }
              this.setState({ serviceGroup, servicProduct });
            } else {
              message.error(res.msg);
            }
          }
        });
    });
  };

  onDeselectComponent = (value) => {
    const { currentProduct, time } = this.state;
    this.setState(
      {
        currentProduct: currentProduct.filter((item) => item != value),
      },
      () => {
        this.handleChangeTime(time);
      }
    );
  };

  // 切换服务
  handSelectService = (item: any) => {
    const { servicProduct } = this.state;
    this.setState({
      currentService: item,
      currentHost: [],
      hostGroup: [],
    });
    this.props.form.setFieldsValue({
      service: item,
      host: [],
    });
    const hostGroup = new Set();
    item.map((value: any) => {
      servicePageService
        .getServiceHostsList({
          product_name: servicProduct[value],
          service_name: value,
        })
        .then((data: any) => {
          data = data.data;
          if (data.code === 0) {
            if (data.data.count > 0) {
              data.data.list.forEach((item: any) => {
                hostGroup.add(item.ip);
              });
              this.setState({ hostGroup: Array.from(hostGroup) });
            }
          } else {
            message.error(data.msg);
          }
        });
    });
  };

  onDeselectService = (value) => {
    const { currentService, time } = this.state;
    this.setState(
      {
        currentService: currentService.filter((item) => item != value),
      },
      () => {
        this.handleChangeTime(time);
      }
    );
  };

  // 切换主机
  handSelectHost = (item: any) => {
    this.setState({
      currentHost: item,
    });
    this.props.form.setFieldsValue({
      host: item,
    });
  };

  onDeselectHost = (value) => {
    const { currentHost, time } = this.state;
    this.setState(
      {
        currentHost: currentHost.filter((item) => item != value),
      },
      () => {
        this.handleChangeTime(time);
      }
    );
  };

  // 统计
  handleGetEventCount = (params: any) => {
    console.log('统计');
    eventService.getEventCount(params).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({ count: res.data });
      } else {
        message.error(res.msg);
        this.setState({
          count: {
            count: null,
            product_count: null,
            service_count: null,
          },
        });
      }
    });
  };

  // 切换时间
  handleChangeTime = (e: any) => {
    const timeType = typeof e === 'string' ? e : e.target.value;
    this.setState({ time: timeType });
    const { cur_parent_product } = this.props.HeaderStore;
    const { currentProduct, currentService, currentHost, currentEventType } =
      this.state;
    if (currentEventType) {
      const end = moment()
        .set({
          hour: 23,
          minute: 59,
          second: 59,
        })
        .unix();
      let dateFrom: any;
      let start: any;
      switch (timeType) {
        case 'day':
          dateFrom = moment().subtract(7, 'd').format('YYYY-MM-DD');
          start = moment(dateFrom).unix();
          break;
        case 'mouth':
          dateFrom = moment().subtract(30, 'd').format('YYYY-MM-DD');
          start = moment(dateFrom).unix();
          break;
        case 'year':
          dateFrom = moment().startOf('year').format('YYYY-MM-DD');
          start = moment(dateFrom).unix();
          break;
      }
      const params = {
        eventType: currentEventType,
        parentProductName: cur_parent_product,
        productNames: currentProduct.join(),
        serviceNames: currentService.join(),
        hosts: currentHost.join(),
        from: start,
        to: end,
      };
      const params1 = {
        eventType: currentEventType,
        parentProductName: cur_parent_product,
        productNames: '',
        serviceNames: '',
        hosts: '',
        from: start,
        to: end,
      };
      this.handleGetEventEcharts(params);
      this.handleGetEventProductRank(params1);
      this.handleGetEventCount(params);
      this.setState({ currentList: 1 }, () => {
        this.handleGetEventList();
      });
    } else {
      message.error('请先选择事件类型！');
    }
  };

  // echarts图
  handleGetEventEcharts = (params: any) => {
    eventService.getEventEcharts(params).then((res: any) => {
      res = res.data;
      const echartsTime = [];
      const echartsData = [];
      if (res.code === 0) {
        const list = res.data.list;
        list.map((item: any) => {
          echartsTime.push(item.date.slice(5));
          echartsData.push(item.count);
        });
        this.setState({ echartsTime, echartsData });
      } else {
        message.error(res.msg);
        this.setState({ echartsTime, echartsData });
      }
    });
  };

  // 切换Tabs
  handleSelectTabs = (e) => {
    this.setState({ activeKey: e });
  };

  // 跳转到事件列表
  handleToEventList = () => {
    this.setState({ activeKey: '2' });
    const { time } = this.state;
    let startData: any;
    const endData = moment().endOf('day');
    switch (time) {
      case 'day':
        startData = moment().startOf('day').subtract(7, 'd');
        break;
      case 'mouth':
        startData = moment().startOf('day').subtract(30, 'd');
        break;
      case 'year':
        startData = moment().startOf('year');
        break;
    }
    this.props.form.setFieldsValue({
      time: [startData, endData],
    });
    this.handleGetEventList();
  };

  // 获取排行
  handleGetEventProductRank = (config: any) => {
    eventService
      .getEventProductRank({ product_or_service: 'product' }, config)
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.setState({ componentData: res.data });
        } else {
          message.error(res.msg);
          this.setState({ componentData: [] });
        }
      });
    eventService
      .getEventProductRank({ product_or_service: 'service' }, config)
      .then((res: any) => {
        res = res.data;
        if (res.code === 0) {
          this.setState({ serviceData: res.data });
        } else {
          message.error(res.msg);
          this.setState({ serviceData: [] });
        }
      });
  };

  // 获取事件列表
  handleGetEventList = () => {
    const {
      currentProduct,
      currentService,
      currentHost,
      currentEventType,
      keyWord,
      currentList,
    } = this.state;
    const { cur_parent_product } = this.props.HeaderStore;
    const time = this.props.form.getFieldValue('time');
    console.log('按照时间筛选，', time);
    const params = {
      eventType: currentEventType,
      parentProductName: cur_parent_product,
      productNames: currentProduct.join(),
      serviceNames: currentService.join(),
      hosts: currentHost.join(),
      keyWord: keyWord,
      limit: 10,
      start: 10 * (currentList - 1),
      'sort-by': 'create_time',
      'sort-dir': 'desc',
      from: time[0].unix(),
      to: time[1].unix(),
    };
    eventService.getEventList(params).then((res: any) => {
      res = res.data;
      const eventListColumns = [];
      const eventListData = [];
      if (res.code === 0) {
        const list = get(res, 'data.list', null);
        if (list) {
          list.map((item, index) => {
            item = JSON.parse(item);
            const p = {};
            for (const i in item) {
              let k = {};
              p[i] = item[i];
              if (index === 0) {
                if (item[i].desc === '操作') {
                  k = {
                    title: item[i].desc,
                    dataIndex: i,
                    key: i,
                    render: (value: any) => {
                      return (
                        <a
                          href={`/api/v2/instance/${value.instance}/logdown?logfile=${value.path}`}
                          download={value.path}>
                          {value.value}
                        </a>
                      );
                    },
                  };
                } else {
                  k = {
                    title: item[i].desc,
                    dataIndex: i,
                    key: i,
                    render: (value: any) => {
                      return <span>{value.value}</span>;
                    },
                  };
                }
                eventListColumns.push(k);
              }
            }
            eventListData.push(p);
          });
          this.setState({
            eventListColumns,
            eventListData,
            eventListCount: res.data.count,
          });
        } else {
          this.setState({ eventListData: [], eventListCount: res.data.count });
        }
      } else {
        message.error(res.msg);
        this.setState({ eventListColumns, eventListData, eventListCount: 0 });
      }
    });
  };

  // 按照文件名搜索
  handleSearchFile = (e) => {
    this.setState({ keyWord: e, currentList: 1 }, () => {
      this.handleGetEventList();
    });
  };

  // 按时间筛选
  handleSelectTime = (value) => {
    this.props.form.setFieldsValue({
      time: value,
    });
    this.setState({ currentList: 1 }, () => {
      this.handleGetEventList();
    });
  };

  // 切换分页
  handleChangepagination = (e) => {
    this.setState({ currentList: e.current }, () => {
      this.handleGetEventList();
    });
  };

  render() {
    const {
      time,
      activeKey,
      eventType,
      products,
      count,
      eventListColumns,
      eventListData,
      componentData,
      serviceData,
      currentEventType,
      echartsTime,
      echartsData,
      serviceGroup,
      hostGroup,
      endDate,
      startDate,
      currentList,
      eventListCount,
    } = this.state;
    const { getFieldDecorator } = this.props.form;
    const { Option, OptGroup } = Select;
    const pagination = {
      current: currentList,
      pageSize: 10,
      total: eventListCount,
    };
    let listTitle;
    switch (time) {
      case 'day':
        listTitle = '最近7天';
        break;
      case 'mouth':
        listTitle = '最近30天';
        break;
      case 'year':
        listTitle = '当年';
        break;
      default:
        listTitle = '最近7天';
    }
    const listData = [
      {
        title: `${listTitle}${currentEventType}次数`,
      },
      {
        title: `发生${currentEventType}组件个数`,
      },
      {
        title: `发生${currentEventType}服务个数`,
      },
    ];
    const formItemLayout = {
      labelCol: { span: 4 },
      wrapperCol: { span: 20 },
    };
    const rankColumns = [
      {
        title: '排行',
        dataIndex: 'rank',
        key: 'rank',
        render: (value) => {
          return value > 3 ? (
            <Avatar
              className="avatar_img"
              style={{ background: '#f5f5f5', color: '#000' }}>
              {value}
            </Avatar>
          ) : (
            <Avatar
              className="avatar_img"
              style={{ background: '#BEDEFD', color: '#43A0F8' }}>
              {value}
            </Avatar>
          );
        },
      },
      {
        title: '组件',
        dataIndex: 'name',
        key: 'name',
        render: (value) => {
          return <span style={{ color: '#43A0F8' }}>{value}</span>;
        },
      },
      {
        title: `${currentEventType}次数`,
        dataIndex: 'times',
        key: 'times',
      },
    ];
    const rankColumns1 = [
      {
        title: '排行',
        dataIndex: 'rank',
        key: 'rank',
        render: (value) => {
          return value > 3 ? (
            <Avatar
              className="avatar_img"
              style={{ background: '#f5f5f5', color: '#000' }}>
              {value}
            </Avatar>
          ) : (
            <Avatar
              className="avatar_img"
              style={{ background: '#BEDEFD', color: '#43A0F8' }}>
              {value}
            </Avatar>
          );
        },
      },
      {
        title: '服务',
        dataIndex: 'name',
        key: 'name',
        render: (value) => {
          return <span style={{ color: '#43A0F8' }}>{value}</span>;
        },
      },
      {
        title: `${currentEventType}次数`,
        dataIndex: 'times',
        key: 'times',
      },
    ];
    return (
      <Layout id="monitorDataFlowContainer">
        <Content>
          <div className="dash-page">
            <div className="top-navbar clearfix">
              <Form className="ant-advanced-search-form">
                <Row className="row_first" type="flex">
                  <Form.Item label="事件类型:">
                    {getFieldDecorator('eventTypes', {
                      initialValue: eventType[0],
                    })(
                      <Select
                        style={{ width: '243px' }}
                        className="dt-form-shadow-bg"
                        onChange={this.handleEventTypeChange}>
                        {eventType
                          ? eventType.map((item, index) => (
                              <Option value={item} key={index}>
                                {item}
                              </Option>
                            ))
                          : ''}
                      </Select>
                    )}
                  </Form.Item>
                  <Form.Item label="组件:">
                    {getFieldDecorator(
                      'component',
                      {}
                    )(
                      <Select
                        mode="multiple"
                        className="dt-form-shadow-bg"
                        style={{ width: '243px' }}
                        showArrow={true}
                        onChange={this.handleSelectProduct}
                        onBlur={() => this.handleChangeTime(time)}
                        onDeselect={(value) => this.onDeselectComponent(value)}>
                        {products
                          ? products.map((item, index) => (
                              <Option value={item.product_name} key={index}>
                                {item.product_name}
                              </Option>
                            ))
                          : ''}
                      </Select>
                    )}
                  </Form.Item>
                  <Form.Item label="服务:">
                    {getFieldDecorator(
                      'service',
                      {}
                    )(
                      <Select
                        mode="multiple"
                        className="dt-form-shadow-bg"
                        showArrow={true}
                        style={{ width: '243px' }}
                        onChange={this.handSelectService}
                        onBlur={() => this.handleChangeTime(time)}
                        onDeselect={(value) => this.onDeselectService(value)}>
                        {serviceGroup
                          ? serviceGroup.map((item, inde) => {
                              return (
                                <OptGroup label={item.name} key={inde}>
                                  {item.list.map((i, index) => (
                                    <Option key={index} value={i.service_name}>
                                      {i.service_name}
                                    </Option>
                                  ))}
                                </OptGroup>
                              );
                            })
                          : []}
                      </Select>
                    )}
                  </Form.Item>
                  <Form.Item label="主机:">
                    {getFieldDecorator(
                      'host',
                      {}
                    )(
                      <Select
                        mode="multiple"
                        className="dt-form-shadow-bg"
                        style={{ width: '243px' }}
                        showArrow={true}
                        onChange={this.handSelectHost}
                        onBlur={() => this.handleChangeTime(time)}
                        onDeselect={(value) => this.onDeselectHost(value)}>
                        {hostGroup
                          ? hostGroup.map((item, index) => (
                              <Option key={index} value={item}>
                                {item}
                              </Option>
                            ))
                          : []}
                      </Select>
                    )}
                  </Form.Item>
                </Row>
              </Form>
            </div>
            <div className="event_contant box-shadow-style">
              <Tabs
                className="c-tabs-padding"
                activeKey={activeKey}
                onChange={this.handleSelectTabs}>
                <TabPane tab="事件全景" key="1">
                  <Radio.Group
                    value={time}
                    onChange={this.handleChangeTime}
                    style={{ marginBottom: 12 }}>
                    <Radio.Button value="day">最近7天</Radio.Button>
                    <Radio.Button value="mouth">最近30天</Radio.Button>
                    <Radio.Button value="year">当年</Radio.Button>
                  </Radio.Group>
                  <Row className="event_contant_chart">
                    <Col span={6} style={{ paddingRight: '20px' }}>
                      <List itemLayout="horizontal" bordered>
                        <List.Item>
                          <List.Item.Meta
                            avatar={
                              <Avatar
                                src={require('./img/restart.png')}
                                style={{ width: 48, height: 48 }}
                              />
                            }
                            title={listData[0].title}
                            description={
                              <a
                                style={{ fontSize: '24px' }}
                                onClick={this.handleToEventList}>
                                {count.count}
                              </a>
                            }
                          />
                        </List.Item>
                        <List.Item>
                          <List.Item.Meta
                            avatar={
                              <Avatar
                                src={require('./img/component.png')}
                                style={{ width: 48, height: 48 }}
                              />
                            }
                            title={listData[1].title}
                            description={
                              <a
                                style={{ fontSize: '24px' }}
                                onClick={this.handleToEventList}>
                                {count.product_count}
                              </a>
                            }
                          />
                        </List.Item>
                        <List.Item>
                          <List.Item.Meta
                            avatar={
                              <Avatar
                                src={require('./img/service.png')}
                                style={{ width: 48, height: 48 }}
                              />
                            }
                            title={listData[2].title}
                            description={
                              <a
                                style={{ fontSize: '24px' }}
                                onClick={this.handleToEventList}>
                                {count.service_count}
                              </a>
                            }
                          />
                        </List.Item>
                      </List>
                    </Col>
                    <Col span={18} className="event_charts">
                      <Echarts
                        echartsTime={echartsTime}
                        echartsData={echartsData}
                        currentEventType={currentEventType}
                      />
                    </Col>
                  </Row>
                  <Row>
                    <Col span={12} style={{ paddingRight: '10px' }}>
                      <div className="event_arrangement_title">
                        组件发生{currentEventType}次数排行
                      </div>
                      <div className="event_arrangement_div">
                        <Table
                          className="event_arrangement_table"
                          size="small"
                          pagination={false}
                          dataSource={componentData}
                          columns={rankColumns}
                        />
                      </div>
                    </Col>
                    <Col span={12} style={{ paddingLeft: '10px' }}>
                      <div className="event_arrangement_title">
                        服务发生{currentEventType}次数排行
                      </div>
                      <div className="event_arrangement_div">
                        <Table
                          className="event_arrangement_table"
                          size="small"
                          pagination={false}
                          dataSource={serviceData}
                          columns={rankColumns1}
                        />
                      </div>
                    </Col>
                  </Row>
                </TabPane>
                <TabPane tab="事件列表" key="2" className="event_list_tab">
                  <div
                    style={
                      {
                        /* height: 667 */
                      }
                    }>
                    <Row>
                      <Form.Item style={{ display: 'inline-block' }}>
                        <Search
                          placeholder="搜索"
                          onSearch={this.handleSearchFile}
                          style={{ width: 298 }}
                        />
                      </Form.Item>
                      <Form.Item
                        label="日期"
                        {...formItemLayout}
                        style={{ display: 'inline-block', marginLeft: '20px' }}>
                        {getFieldDecorator('time', {
                          initialValue: [startDate, endDate],
                        })(
                          <RangePicker
                            allowClear={false}
                            style={{ width: 236 }}
                            onChange={this.handleSelectTime}
                          />
                        )}
                      </Form.Item>
                    </Row>
                    <Table
                      className="event_list_table"
                      pagination={pagination}
                      dataSource={eventListData}
                      columns={eventListColumns}
                      onChange={this.handleChangepagination}
                    />
                  </div>
                </TabPane>
              </Tabs>
            </div>
          </div>
        </Content>
      </Layout>
    );
  }
}
export default Form.create<any>()(Event);
