import * as React from 'react';
import { message, Button, Input, Divider, Tag } from 'antd';
import { clusterManagerService } from '@/services';

interface IProps {
  clusterInfo: any;
}

interface IState {
  data: {
    secure: string;
    insecure: string;
    secure_v1beta1: string;
    insecure_v1beta1: string;
  };
}

export default class KubernetesCmd extends React.PureComponent<IProps, IState> {
  state: IState = {
    data: {
      secure: '',
      insecure: '',
      secure_v1beta1: '',
      insecure_v1beta1: '',
    },
  };

  componentDidMount() {
    this.getData();
  }

  // 获取命令行信息
  getData() {
    const { clusterInfo } = this.props;
    clusterManagerService
      .getKubernetesInstallCmd({
        cluster_id: clusterInfo.id,
        cluster_name: clusterInfo.name,
      })
      .then((res) => {
        res = res.data;
        if (res.code === 0) {
          this.setState({ data: res.data });
        } else {
          message.error(res.msg);
        }
      });
  }

  render() {
    const { data } = this.state;
    return (
      <div className="cluster-cmd-page">
        <CmdComponent
          version={'>=1.16'}
          secure={data.secure}
          insecure={data.insecure}
        />
        <Divider dashed style={{ margin: '20px 0' }} />
        <CmdComponent
          version={'1.12-1.15'}
          secure={data.secure_v1beta1}
          insecure={data.insecure_v1beta1}
        />
      </div>
    );
  }
}

interface CmdComponentProps {
  version: string;
  secure: string;
  insecure: string;
}

const CmdComponent = (props: CmdComponentProps) => {
  const { version, secure, insecure } = props;
  const TEXT_ID = 'J_CMDCode' + (version === '>=1.16' ? '' : '_v1beta1');
  const description = [
    `Kubernetes 版本${version}，运行下面的kubectl命令将其导入到ChengYing：`,
    '如果因为ChengYing使用不受信任 / 自签名的SSL证书而出现“由未知颁发机构签名的证书”错误，请运行下面的命令以绕过证书检查：',
  ];
  const tagStyle = {
    color: '#3f87ff',
    background: '#f0f8ff',
    borderColor: '#3f87ff',
    borderRadius: '11px',
  };
  return (
    <div>
      <p>
        Kubernetes 版本 <Tag style={tagStyle}>{version}</Tag>
        ，运行下面的kubectl命令将其导入到ChengYing：
      </p>
      <p className="code-content mb-20">{secure}</p>
      <p>{description[1]}</p>
      <p className="code-content">{insecure}</p>
      <React.Fragment>
        <Input.TextArea
          id={TEXT_ID}
          style={{ opacity: 0, height: 12, minHeight: 20 }}
          value={`${description[0]}\n${secure}\n${description[1]}\n${insecure}`}
        />
        <Button
          type="default"
          id="J_CopyBtn"
          data-clipboard-action="cut"
          data-clipboard-target={`#${TEXT_ID}`}>
          复制命令
        </Button>
      </React.Fragment>
    </div>
  );
};
