import * as React from 'react';
import { DatePicker, Input, Select, message } from 'antd';
import { SecurityAuditService } from '@/services';
import { AuditReqParamsType } from './interface';
import AuditTable from './auditTable';
import './style.scss';

const { RangePicker } = DatePicker;
const { Option } = Select;

const searchTypeList = [
  {
    value: 'content',
    text: '详细内容',
    searchKey: 'content',
  },
  {
    value: 'operator',
    text: '操作人',
    searchKey: 'operator',
  },
  {
    value: 'ip',
    text: 'IP',
    searchKey: 'ip',
  },
];

interface IState {
  moduleList: any[];
  operationList: any[];
  reqParams: AuditReqParamsType;
  searchType: {
    value: string;
    text: string;
    searchKey: string;
  };
  searchVal: string;
}

export default class SecurityAudit extends React.PureComponent<any, IState> {
  state: IState = {
    moduleList: [],
    operationList: [],
    reqParams: {
      from: undefined,
      to: undefined,
      module: undefined,
      operation: undefined,
      content: undefined,
      operator: undefined,
      ip: undefined,
    },
    searchType: searchTypeList[0],
    searchVal: '',
  };

  componentDidMount() {
    this.getAuditModuleList();
  }

  // 获取操作模块列表
  getAuditModuleList = () => {
    SecurityAuditService.getAuditModule().then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({
          moduleList: res.data ? res.data.list : [],
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 获取动作列表
  getAuditOperation = (module: string) => {
    if (!module) {
      this.setState({ operationList: [] });
      return;
    }
    SecurityAuditService.getAuditOperation({
      module,
    }).then((res: any) => {
      res = res.data;
      if (res.code === 0) {
        this.setState({
          operationList: res.data ? res.data.list : [],
        });
      } else {
        message.error(res.msg);
      }
    });
  };

  // 时间选取
  handleTimeRangeChange = (
    dates: [moment.Moment, moment.Moment],
    dateStrings: [string, string]
  ) => {
    const timeRange = {
      from: dates[0] ? dates[0].startOf('day').unix() : undefined,
      to: dates[1] ? dates[1].endOf('day').unix() : undefined,
    };
    this.onFilterChange(timeRange);
  };

  // 操作模块变更
  handleOperateModalChange = (value: string) => {
    const params = {
      module: value,
      operation: undefined,
    };
    this.getAuditOperation(value);
    this.onFilterChange(params);
  };

  // 选择动作
  handleActionChange = (operation: string) => {
    this.onFilterChange({ operation });
  };

  // 搜索类型变更
  handleSearchTypeChange = (type: string) => {
    const searchType = searchTypeList.find((item) => item.value === type);
    this.setState({
      searchType: searchType,
      searchVal: undefined,
      reqParams: Object.assign({}, this.state.reqParams, {
        operator: '',
        ip: '',
        content: '',
      }),
    });
  };

  // 搜索关键词变更
  handleSearchValChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    this.setState({ searchVal: e.target.value });
  };

  // 按关键词搜索
  handleSearch = (searchVal: string) => {
    const { searchType } = this.state;
    this.onFilterChange({
      [searchType.searchKey]: searchVal,
    });
  };

  // 筛选条件变更
  onFilterChange = (params: any) => {
    const reqParams = {
      ...this.state.reqParams,
      ...params,
    };
    this.setState({ reqParams });
  };

  render() {
    const { reqParams, operationList, moduleList, searchType, searchVal } =
      this.state;
    const { module, operation } = reqParams;
    return (
      <div className="security-audit-page">
        <div className="audit-header">
          <RangePicker
            className="dt-form-shadow-bg mr-20"
            style={{ width: 300 }}
            allowClear
            onChange={this.handleTimeRangeChange}
          />
          <div className="mr-20">
            <span>操作模块：</span>
            <Select
              className="dt-form-shadow-bg"
              style={{ width: 200 }}
              placeholder="选择操作模块"
              allowClear
              value={module}
              onChange={this.handleOperateModalChange}>
              {Array.isArray(moduleList) &&
                moduleList.map((item: string) => (
                  <Option key={item} value={item}>
                    {item}
                  </Option>
                ))}
            </Select>
          </div>
          <div className="mr-20">
            <span>动作：</span>
            <Select
              className="dt-form-shadow-bg"
              style={{ width: 200 }}
              placeholder="选择动作"
              allowClear
              value={operation}
              onChange={this.handleActionChange}>
              {Array.isArray(operationList) &&
                operationList.map((item: string) => (
                  <Option key={item} value={item}>
                    {item}
                  </Option>
                ))}
            </Select>
          </div>
          <Input.Search
            className="dt-select-search-bar"
            style={{ width: 300 }}
            placeholder={`按${searchType.text}搜索`}
            value={searchVal}
            addonBefore={
              <Select
                value={searchType.value}
                onChange={this.handleSearchTypeChange}>
                {searchTypeList.map((item: any) => (
                  <Option key={item.value} value={item.value}>
                    {item.text}
                  </Option>
                ))}
              </Select>
            }
            onChange={this.handleSearchValChange}
            onSearch={this.handleSearch}
          />
        </div>
        <AuditTable reqParams={reqParams} />
      </div>
    );
  }
}
