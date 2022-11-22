import * as React from 'react';
import * as HeaderAction from '@/actions/headerAction';
import { Icon, Tooltip, Cascader } from 'antd';
import * as Cookie from 'js-cookie';
import utils from '@/utils/utils';
import './style.scss';

interface Props {
  location?: any;
  HeaderStore?: any;
  actions?: any;
}

interface States {
  tooltipVisible: boolean;
  popupVisible: boolean;
}

export default class Header extends React.Component<Props, States> {
  state: States = {
    tooltipVisible: false,
    popupVisible: false,
  };

  componentDidMount() {
    this.props.actions.getClusterProductList();
  }

  componentWillUnmount() {
    this.props.actions.setCurrentParentProduct('');
  }

  setOptions = (list: HeaderAction.ClusterItem[]) => {
    const options = [];
    Array.isArray(list) &&
      list.forEach((item) => {
        const { subdomain = {}, mode, clusterName } = item;
        const option: any = {
          value: clusterName,
          label: clusterName,
        };
        if (mode === 0) {
          const { products = [] } = subdomain;
          option.children = products.map((item) => ({
            value: item,
            label: item,
          }));
        } else {
          option.children = [];
          Object.keys(subdomain).forEach((key) => {
            option.children.push({
              value: key,
              label: key,
              children: subdomain[key].map((p) => ({
                value: p,
                label: p,
              })),
            });
          });
        }
        options.push(option);
      });
    return options;
  };

  handleSwitchProduct = (selectedKeys: any[]) => {
    const { parentProducts } = this.props.HeaderStore;
    const cluster = parentProducts.find(
      (item) => item.clusterName === selectedKeys[0]
    );
    // console.log(cluster)
    if (cluster?.mode === 0) {
      this.props.actions.setCurrentParentProduct(selectedKeys[1]);
      Cookie.set('em_current_parent_product', selectedKeys[1]);
      Cookie.set('em_current_k8s_namespace', '');
    } else {
      this.props.actions.setCurrentParentProduct(selectedKeys[2]);
      Cookie.set('em_current_parent_product', selectedKeys[2]);
      Cookie.set('em_current_k8s_namespace', selectedKeys[1]);
    }
    Cookie.set('em_current_cluster_id', cluster.clusterId);
    Cookie.set('em_current_cluster_type', cluster.clusterType);
    utils.setNaviKey(
      'menu_ops_center',
      cluster?.type === 'hosts' ? 'sub_menu_overview' : 'sub_menu_service'
    );
    window.location.href =
      cluster?.type === 'hosts' ? '/opscenter/overview' : '/opscenter/service';
  };

  render() {
    const { HeaderStore } = this.props;
    const { popupVisible } = this.state;
    const {
      parentProducts = [],
      cur_parent_product,
      cur_parent_cluster,
    } = HeaderStore;
    const options = this.setOptions(parentProducts);
    const cur_namespace = Cookie.get('em_current_k8s_namespace');
    const showStr = `${cur_parent_cluster.name}/${
      cur_namespace && cur_namespace !== 'undefined'  ? `${cur_namespace}/` : ''
    }${cur_parent_product}`;
    return (
      <Cascader
        popupClassName="header-cascader-pop"
        options={options}
        onChange={this.handleSwitchProduct}
        onPopupVisibleChange={(value: boolean) =>
          this.setState({ popupVisible: value })
        }>
        <Tooltip
          title={showStr}
          placement="bottom"
          visible={this.state.tooltipVisible}
          align={{ offset: [0, -7] }}
          onVisibleChange={(bool: boolean) =>
            this.setState({ tooltipVisible: popupVisible ? false : bool })
          }>
          <span
            className="c-text-ellipsis header-cascader-title"
            onClick={() => this.setState({ tooltipVisible: false })}>
            {showStr}
          </span>
          <Icon type="down" />
        </Tooltip>
      </Cascader>
    );
  }
}
