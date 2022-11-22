import * as React from 'react';
import { TransferListBodyProps } from 'antd/lib/transfer/renderListBody';
import Table, { ColumnProps } from 'antd/lib/table';
import difference from 'lodash/difference';
import './style.scss';

export interface TransferItem {
  id?: string | number;
  key: string;
  [key: string]: string | number;
}

export interface TransferTableBodyProps extends TransferListBodyProps {
  footer: (currentPageData: unknown) => React.ReactNode;
  renderItem: (param: TransferItem) => React.ReactNode;
}

const TransferTable: React.FC<TransferTableBodyProps> = (props) => {
  const {
    filteredItems,
    onItemSelectAll,
    onItemSelect,
    selectedKeys,
    renderItem,
    direction,
  } = props;
  const columns: ColumnProps<TransferItem>[] = [
    {
      title: '',
      render: (record: TransferItem) => renderItem(record),
    },
  ];
  const rowSelection = {
    onSelectAll(selected, selectedRows) {
      const treeSelectedKeys = selectedRows
        .filter((item) => item)
        .map(({ key }) => key);
      const diffKeys = selected
        ? difference(treeSelectedKeys, selectedKeys)
        : difference(selectedKeys, treeSelectedKeys);
      onItemSelectAll(diffKeys, selected);
    },
    onSelect({ key }, selected) {
      onItemSelect(key, selected);
    },
    getCheckboxProps: item => ({ disabled: item.disabled }),
    selectedRowKeys: selectedKeys
  };

  const pagination = {
    size: 'small',
    total: filteredItems.length,
    showSizeChanger: true,
    showQuickJumper: true,
  };
  return (
    <Table
      className="dt-table-last-row-noborder dt-table-fixed-contain-footer"
      rowSelection={rowSelection}
      showHeader={false}
      rowKey="id"
      data-testid={direction}
      columns={columns}
      scroll={{ y: true }}
      dataSource={filteredItems}
      size="small"
      onRow={({ key }) => ({
        onClick: () => {
          onItemSelect(key, !selectedKeys.includes(key));
        },
      })}
      style={{ height: 313 }}
      pagination={pagination}
      // pagination={false}
      // footer={footer}
    />
  );
};
export default TransferTable;
