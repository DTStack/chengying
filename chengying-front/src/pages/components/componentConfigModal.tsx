import * as React from 'react';
import { get } from 'lodash';
import { servicePageService, productService } from '@/services';
import FileViewModal from '@/components/fileViewModal';
import { ServiceTree, ServiceFile } from '@/model/apis';
import { QueryParams } from './container';

interface Props extends QueryParams {
  componentData: any;
  onClose?: Function;
  shouldGetAll?: boolean;
}

interface State {
  loading: boolean;
  configList: any[];
  showModal: boolean;
  serviceGroup: any;
  modalContent: string;
  selectedFile: string;
  selectedService: string;
  visibleConfigInfo: boolean;
}

class ComponentConfigModal extends React.Component<Props, State> {
  state: State = {
    loading: false,
    selectedService: null,
    configList: [],
    showModal: false,
    modalContent: null,
    visibleConfigInfo: false,
    selectedFile: undefined,
    serviceGroup: null,
  };

  componentDidMount() {
    const { componentData } = this.props;
    if (this.props.componentData) {
      this.onShowConfigInfo(componentData);
    }
  }

  componentDidUpdate(prevProps) {
    if (prevProps.componentData !== this.props.componentData) {
      this.onShowConfigInfo(this.props.componentData);
    }
  }

  loadServiceGroup = async (params: any) => {
    const { shouldGetAll } = this.props;
    const res = shouldGetAll
      ? await servicePageService.getAllServiceGroup(params)
      : await servicePageService.getServiceGroup(params);
    if (res.data.code === 0) {
      this.setState({
        serviceGroup: get(res, 'data.data.groups', {}),
      });
    }
  };

  onShowConfigInfo = (record: any) => {
    const { shouldGetAll } = this.props;
    if (record && record.product_name) {
      this.setState(
        {
          visibleConfigInfo: true,
        },
        () => {
          const params: any = {
            product_name: record.product_name,
          };
          if (shouldGetAll) {
            params.product_version = record.product_version;
          }
          this.loadServiceGroup(params);
        }
      );
    }
  };

  getConfigList = (params: ServiceTree) => {
    productService.getServiceTree(params).then((res: any) => {
      if (res.data.code === 0) {
        this.setState({
          configList: get(res, 'data.data.list', []),
        });
      }
    });
  };

  getConfigContent = (params: ServiceFile) => {
    productService.getServiceFile(params).then((res: any) => {
      if (res.data.code === 0) {
        this.setState({
          modalContent: get(res, 'data.data', ''),
        });
      }
    });
  };

  onSelectedFile = (file: string) => {
    const { componentData } = this.props;
    this.setState({ selectedFile: file });
    const { selectedService } = this.state;
    this.getConfigContent({
      file: file,
      serviceName: selectedService,
      productName: componentData.product_name,
      productVersion: componentData.product_version,
    });
  };

  onSelectedConfigService = (service: string) => {
    const { componentData } = this.props;
    this.setState({
      selectedService: service,
      selectedFile: undefined,
      configList: [],
    });
    this.getConfigList({
      serviceName: service,
      productName: componentData.product_name,
      productVersion: componentData.product_version,
    });
  };

  onCloseModal = () => {
    this.setState(
      {
        visibleConfigInfo: false,
        selectedService: null,
        selectedFile: undefined,
        modalContent: null,
      },
      () => {
        this.props.onClose();
      }
    );
  };

  render = () => {
    const { visibleConfigInfo } = this.state;

    return (
      <React.Fragment>
        <FileViewModal
          key="configModal"
          title="配置信息"
          fileList={this.state.configList}
          visible={visibleConfigInfo}
          content={this.state.modalContent}
          serviceData={this.state.serviceGroup}
          onCancel={this.onCloseModal}
          selectedFile={this.state.selectedFile}
          onSelectedFile={this.onSelectedFile}
          onSelectedService={this.onSelectedConfigService}
        />
      </React.Fragment>
    );
  };
}

export default ComponentConfigModal;
