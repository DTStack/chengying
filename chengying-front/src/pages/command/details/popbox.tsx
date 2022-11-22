/*
 * @Description: 弹窗
 * @Author: wulin
 * @Date: 2021-06-01 14:00:41
 * @LastEditors: wulin
 * @LastEditTime: 2021-06-23 17:34:29
 */
import * as React from 'react';
import { Modal, Empty, Button, message, Spin } from 'antd';
import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/mode-golang';
import 'ace-builds/src-noconflict/mode-powershell';
import 'ace-builds/src-noconflict/theme-kuroir';
import { IPopboxState } from './detail';
import '../style.scss';
import echoService from '@/services/echoService';
import moment from 'moment';
import * as Cookies from 'js-cookie';
import * as _ from 'lodash';

const download = (url: string, filename: string, config: Object = {}) => {
  const completeUrl = Object.keys(config)
    .reduce((temp, key) => {
      return `${temp}${key}=${config[key]}&`;
    }, `${url}?`)
    .replace(/&$/, '');

  fetch(completeUrl)
    .then((res) => res.blob())
    .then((res) => {
      const objectUrl = window.URL.createObjectURL(res);
      const ele = document.createElement('a');
      ele.href = objectUrl;
      ele.download = filename;
      ele.click();
      window.URL.revokeObjectURL(ele.href);
    })
    .catch((err) => {
      message.error(err.message);
    });
};

type ICommandProps = IPopboxState & {
  onColse: () => void;
  showFooter?: boolean;
};

const CommandPopbox = (props: ICommandProps) => {
  const { title, visible, type, onColse, execId, showFooter = true } = props;

  const refLog = React.useRef(null);

  const handleDownload = React.useCallback(({ type, payload }) => {
    switch (type) {
      case 'log':
        download(
          '/api/v2/cluster/downLoadShellLog',
          `${execId}-shell.log`,
          payload
        );
        break;
      case 'echo':
        download(
          '/api/v2/cluster/downLoadShellContent',
          `${execId}-content.sh`,
          payload
        );
        break;
      default:
    }
  }, []);

  // if (!execId) return null;
  const footer: JSX.Element =
    type === 'echo' ? (
      <Button
        type="primary"
        disabled={!execId}
        onClick={_.debounce(
          () =>
            handleDownload({
              type: 'echo',
              payload: {
                execId,
              },
            }),
          200
        )}>
        下载完整脚本
      </Button>
    ) : (
      <Button
        type="primary"
        disabled={!execId}
        onClick={_.debounce(() => {
          handleDownload({
            type: 'log',
            payload: {
              clusterId: Cookies.get('em_current_cluster_id'),
              execId,
            },
          });
        }, 200)}>
        下载完整日志
      </Button>
    );

  const handleClose = (e) => {
    if (refLog.current) {
      refLog.current.wsClose();
    }
    onColse();
  };

  return (
    <Modal
      className="code-view-modal"
      width={type === 'echo' ? 520 : 1050}
      title={title}
      visible={visible}
      footer={showFooter ? footer : null}
      onCancel={handleClose}>
      {type === 'echo' ? (
        <AceContent execId={execId} />
      ) : (
        <LogContent execId={execId} cref={refLog} />
      )}
    </Modal>
  );
};

/**
 * 脚本查看
 * @param type
 * @returns
 */
const AceContent: React.FC<any> = ({ execId }) => {
  const [content, setContent] = React.useState('');

  const showShellContent = (execId) => {
    echoService
      .showShellContent({ execId })
      .then(({ data: res }) => {
        const { data, code, msg } = res;
        if (code === 0) {
          return setContent(data);
        }
        message.error(msg);
      })
      .catch((err) => {
        message.error(err.message);
      });
  };

  React.useEffect(() => {
    execId && showShellContent(execId);
  }, [execId]);

  return (
    <AceEditor
      className="ace-code-portal"
      mode="golang"
      theme="kuroir"
      value={content}
      readOnly={true}
      width="400px"
      height="100%"
      showGutter={false}
    />
  );
};

/**
 * 日志查看
 * @param type
 * @returns
 */
const LogContent: React.FC<any> = ({ execId, cref }) => {
  const [shellRecords, setShellRecords] = React.useState<
    {
      time: string;
      message: string;
    }[]
  >([]);
  const [loading, changeLoadingStatus] = React.useState<boolean>(false);
  const ws = React.useRef(null);

  React.useImperativeHandle(cref, () => {
    return {
      wsClose: () => {
        ws.current && ws.current.close();
        ws.current = null;
      },
    };
  });

  const bindWsEvent = () => {
    ws.current.onmessage = (e) => {
      setShellRecords((records) => [
        ...records,
        { time: `[${moment().format('hh:mm:ss')}]`, message: e.data },
      ]);
      changeLoadingStatus(false);
    };
    ws.current.onopen = (e) => {
      console.log('websocket open');
      changeLoadingStatus(true);
    };
    ws.current.onerror = (err) => {
      console.error(err);
      changeLoadingStatus(false);
    };
    ws.current.onclose = () => {
      console.log('onclose');
      changeLoadingStatus(false);
    };
  };

  React.useEffect(() => {
    if (!execId) return;
    try {
      ws.current = new WebSocket(
        `ws://${window.location.hostname}:${window.location.port}/api/v2/cluster/showShellLog?execId=${execId}`
      );
      bindWsEvent();
    } catch (error) {
      console.log(error);
      changeLoadingStatus(false);
    }
  }, [execId]);

  // 日志渲染
  const logInfoRender = (records): string => {
    let content: string = '';
    for (let i in records) {
      const item: string = `
        <li>
          <pre>${records[i].message}</pre>
        </li>
      `;
      content += item;
    }
    return content;
  };

  return (
    <Spin spinning={loading}>
      {shellRecords.length ? (
        <div className="log-info">
          <ul
            dangerouslySetInnerHTML={{ __html: logInfoRender(shellRecords) }}
          />
        </div>
      ) : (
        <Empty />
      )}
    </Spin>
  );
};

export default React.memo(CommandPopbox, (prevProps, nextProps) => {
  return prevProps.execId === nextProps.execId;
});
