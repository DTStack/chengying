import * as React from 'react';
import FileViewModal from '@/components/fileViewModal';
import { message } from 'antd';

interface IProps {
  wsUrl: string;
}

interface IState {
  fileContent: string;
  fileViewShow: boolean;
}

export default class FileLogShow extends React.PureComponent<IProps, IState> {
  state: IState = {
    fileContent: '',
    fileViewShow: false,
  };

  private ws: any;

  // 查看完整日志
  handleShowLog = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
    e.preventDefault();
    this.setState(
      {
        fileViewShow: true,
      },
      () => {
        this.websocketLink();
      }
    );
  };

  // 开启websocket链接
  websocketLink = () => {
    let error = false;
    const { wsUrl } = this.props;
    this.ws = new WebSocket(wsUrl);
    // this.ws = new WebSocket('ws://localhost:9999');
    this.ws.onmessage = (res: any) => {
      this.setState({
        fileContent: res.data,
      });
    };
    this.ws.onerror = (res: any) => {
      error = true;
      message.error('实时日志获取失败！');
    };
    this.ws.onclose = () => {
      // 断开重连
      if (this.state.fileViewShow && !error) {
        this.websocketLink();
      }
    };
  };

  render() {
    const { fileViewShow, fileContent } = this.state;
    return (
      <React.Fragment>
        <a
          style={{
            fontSize: '12px',
            display: 'inline-block',
            marginTop: '10px',
          }}
          onClick={this.handleShowLog}>
          查看完整日志
        </a>
        {fileViewShow && (
          <FileViewModal
            title="日志"
            serviceData={null}
            maskClosable={true}
            visible={fileViewShow}
            content={fileContent}
            fileGoBottom={true}
            loading={true}
            onCancel={() => {
              this.setState(
                {
                  fileViewShow: false,
                  fileContent: '',
                },
                () => {
                  this.ws.close();
                }
              );
            }}
          />
        )}
      </React.Fragment>
    );
  }
}
