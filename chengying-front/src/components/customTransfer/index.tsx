import * as React from 'react';
import Transfer, { TransferProps, TransferListProps } from 'antd/lib/transfer';
import { TransferListBodyProps } from 'antd/lib/transfer/renderListBody';
import TransferTable, { TransferTableBodyProps } from './transferTable';
import './style.scss';

interface CustomTransferProps extends TransferProps {
  footerLeft?: (props: TransferListProps) => React.ReactNode;
  footerRight?: (props: TransferListProps) => React.ReactNode;
}

const CustomTransfer: React.FC<CustomTransferProps> = (props) => {
  const { footerLeft, footerRight, render } = props;
  return (
    <Transfer className="transfer-table" {...props}>
      {(listProps: TransferListBodyProps) => {
        const { direction } = listProps;
        const footer = direction === 'left' ? footerLeft : footerRight;
        const transferProps: TransferTableBodyProps = {
          ...listProps,
          footer,
          renderItem: render,
        };
        return <TransferTable {...transferProps} />;
      }}
    </Transfer>
  );
};

export default CustomTransfer;
