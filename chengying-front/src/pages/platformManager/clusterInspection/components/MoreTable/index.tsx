import React, { useState, useEffect } from 'react';
import { Table, message } from 'antd';

import Api from '@/services/clusterInspectionService';
import '../../index.scss';

const MoreTable = (props) => {
  const { listItem, cluster_id, showTableVisible } = props;
  const [list, setList] = useState([]);
  const [visible, setVisible] = useState(showTableVisible);
  const [columnsData, setColumnsData] = useState([]);

  const setTableData = async () => {
    const resparams = {
      cluster_id,
      form_title: listItem.form_title,
    };
    const res = await Api.getClusterInspectionTableData(resparams);
    if (res.data.code == 0) {
      const { form_head, form_value } = res.data.data;
      let columns = [];
      let dataList = [];
      form_head.map((item, index) => {
        columns.push({
          title: item,
          dataIndex: item,
          key: `${item}-${index}`,
        });
      });
      form_value.forEach((value, index) => {
        let obj1: any = {};
        for (let i = 0; i < form_head.length; i++) {
          obj1[form_head[i]] = value[i];
          obj1.key = `${value[i]}-${index}`;
        }
        dataList.push(obj1);
      });
      setColumnsData(columns);
      setList(dataList);
    } else {
      message.error(res.data.msg);
    }
  };

  useEffect(() => {
    setVisible(showTableVisible);
    if (cluster_id && showTableVisible) {
      setTableData();
    }
  }, [cluster_id, showTableVisible]);

  return (
    <>
      {visible ? (
        <div className="moreTable">
          <div className="moreTable-title">{listItem.desc}</div>
          <Table
            className="dt-table-border dt-table-last-row-noborder"
            columns={columnsData}
            dataSource={list}
            pagination={false}
          />
        </div>
      ) : null}
    </>
  );
};

export default MoreTable;
