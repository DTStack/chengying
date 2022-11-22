/*
 * @Description: 脚本回显列表
 * @Author: wulin
 * @Date: 2021-06-01 14:00:41
 * @LastEditors: wulin
 * @LastEditTime: 2021-06-21 20:15:53
 */
import * as React from 'react';
import { Input, DatePicker, Select, AutoComplete, Icon, message } from 'antd';
import { operationMaps } from './index';
import { echoService } from '@/services';
import * as Cookie from 'js-cookie';
import moment from 'moment';
import './style.scss';
const { Option } = Select;
const { RangePicker } = DatePicker;

type ISearchProps = {
  handleEvent: (params) => void;
};

const AutoSearch: React.FC<ISearchProps> = ({ handleEvent }) => {
  const [dataSource, getDataSource] = React.useState<string[]>([]);
  const [list, getDataList] = React.useState<string[]>([]);

  React.useEffect(() => {
    getEchoList();
  }, []);

  /**
   * 获取列表数据
   */
  const getEchoList = () => {
    echoService
      .getEchoSearchList({
        clusterId: Cookie.get('em_current_cluster_id'),
      })
      .then((res: any) => {
        const ret = res.data;
        if (ret.code === 0) {
          getDataSource(ret.data || []);
          getDataList(ret.data || []);
        } else {
          message.error(ret.msg);
        }
      });
  };

  const handleFilter = (value) => {
    let list = dataSource;
    if (!value) {
      handleEvent((prevValues) => ({
        ...prevValues,
        objectValue: value,
        page: 1,
      }));
    } else {
      list = dataSource.filter(
        (o) => value && o.toLowerCase().indexOf(value.toLowerCase()) > -1
      );
    }
    getDataList(list);
  };

  const renderOption = (item) => <Option key={item}>{item}</Option>;

  return (
    <AutoComplete
      dataSource={list.map(renderOption)}
      onSelect={(value: string): void => {
        handleEvent((prevValues) => ({
          ...prevValues,
          objectValue: value,
          page: 1,
        }));
      }}
      onChange={(value: string): void => handleFilter(value)}
      allowClear>
      <Input placeholder="请输入搜索关键词" suffix={<Icon type="search" />} />
    </AutoComplete>
  );
};

const EchoSearch: React.FC<ISearchProps> = ({ handleEvent }) => (
  <div className="echo-form clearfix mb-12">
    <span className="mr-20">
      对象：
      <AutoSearch handleEvent={handleEvent} />
    </span>
    <span className="mr-20">
      操作：
      <Select
        placeholder="请选择命令"
        showSearch
        onChange={(value: string): void => {
          handleEvent((prevValues) => ({
            ...prevValues,
            operationType: value,
            page: 1,
          }));
        }}
        allowClear>
        {Object.keys(operationMaps).map((item) => (
          <Option key={item} value={item}>
            {operationMaps[item]}
          </Option>
        ))}
      </Select>
    </span>
    <span className="mr-20">
      状态：
      <Select
        placeholder="请选择状态"
        showSearch
        onChange={(value: string): void => {
          handleEvent((prevValues) => ({
            ...prevValues,
            status: value,
            page: 1,
          }));
        }}
        allowClear>
        <Option value="2">正常</Option>
        <Option value="3">失败</Option>
        <Option value="1">进行中</Option>
      </Select>
    </span>
    <span className="mr-20">
      开始时间：
      <RangePicker
        format="YYYY-MM-DD HH:mm"
        disabledDate={disabledDate}
        showTime
        onChange={(date, dateString: string[]): React.ReactNode | void => {
          if (
            dateString.length &&
            dateString[0] === '' &&
            dateString[1] === ''
          ) {
            handleEvent((prevValues) => ({
              ...prevValues,
              startTime: dateString[0],
              endTime: dateString[1],
              page: 1,
            }));
          }
        }}
        onOk={(dates): React.ReactNode | void => {
          const startTime =
            dates && dates.length
              ? moment(dates[0]).format('YYYY-MM-DD HH:mm')
              : '';
          const endTime =
            dates && dates.length
              ? moment(dates[1]).format('YYYY-MM-DD HH:mm')
              : '';
          if (startTime && endTime && startTime === endTime) {
            return message.warning('开始时间必须小于结束时间！');
          }
          handleEvent((prevValues) => ({
            ...prevValues,
            startTime,
            endTime,
            page: 1,
          }));
        }}
      />
    </span>
  </div>
);

export default EchoSearch;

function disabledDate(current) {
  // Can not select days before today and today
  return current && current.valueOf() > Date.now();
}
