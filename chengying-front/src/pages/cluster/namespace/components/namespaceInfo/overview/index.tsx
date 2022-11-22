import * as React from 'react';
import { Select, Input, message, Tag, Empty } from 'antd';
import ServiceOverview from '@/components/serviceOverview';
import { ClusterNamespaceService } from '@/services';
import * as Cookies from 'js-cookie';
import { COOKIES } from '@/constants/const';
import utils from '@/utils/utils';
import './style.scss';

const Search = Input.Search;
const Option = Select.Option;

interface IProps {
  history?: any;
  namespace: string;
}
interface ReqParamsTypes {
  namespace: string;
  parentProductName: string;
  productName: string;
  serviceName: string;
}
const OverviewPage: React.FC<IProps> = (props) => {
  const { namespace } = props;
  const [reqParams, setReqParams] = React.useState<ReqParamsTypes>({
    namespace,
    parentProductName: 'all',
    productName: 'all',
    serviceName: 'all',
  });
  const [parentProductList, setParentProductList] = React.useState<string[]>(
    []
  );
  const [collapseKey, setCollapseKey] = React.useState<string[] | string>([]);
  const [serviceList, setServiceList] = React.useState<any[]>([]);
  const [serviceGroup, setServiceGroup] = React.useState<any[]>([]);

  React.useEffect(() => {
    getParentProductList();
  }, []);

  React.useEffect(() => {
    console.log(reqParams);
    if (parentProductList.length > 0) {
      getServiceList();
    }
  }, [reqParams]);

  // 获取已部署产品列表
  function getParentProductList() {
    ClusterNamespaceService.getServiceLists({
      namespace,
      parentProductName: 'all',
    }).then((response: any) => {
      const res = response.data;
      const { code, data, msg } = res;
      if (code === 0 && data) {
        const list = Object.keys(data.products);
        const parentProductName =
          list.find((name: string) => name === 'DTinsight') || list[0];
        setParentProductList(list);
        setReqParams({
          ...reqParams,
          parentProductName,
        });
      } else {
        message.error(msg);
      }
    });
  }

  // 获取服务列表
  function getServiceList() {
    ClusterNamespaceService.getServiceLists({
      namespace,
      ...reqParams,
    }).then((response: any) => {
      const res = response.data;
      const { code, data, msg } = res;
      if (code === 0 && data) {
        const parentProduct: any = Object.values(data.products)[0];
        const products = [];
        const collapseKey = [];
        const serviceGroup = [];
        // 重置数据结构
        for (const i in parentProduct) {
          // 健康服务个数
          let health_service_count = 0;
          // 重置数据service_list结构
          const cur_product = parentProduct[i];
          const service_names = Object.keys(cur_product);
          const service_list = Object.values(cur_product);
          service_list.forEach((service: any, index: number) => {
            Object.assign(service, {
              service_name: service_names[index],
              status: service.service_status === 'available' ? '正常' : '异常',
              status_count: service.service_status_count,
            });
            // 判断服务是否健康
            if (service.health_state === 'healthy') {
              ++health_service_count;
            }
            // 分组
            service.group && serviceGroup.push(service.group);
          });
          // 重置product结构
          const service_list_len = service_list.length;
          const product = {
            product_name: i,
            service_count: service_list_len,
            service_status:
              service_list_len === health_service_count
                ? 'healthy'
                : 'unhealthy',
            service_list,
          };
          /**
           * 排序，异常在前，正常在后
           * 展开，异常展开，正常不展开
           */
          if (product.service_status === 'healthy') {
            products.push(product);
          } else {
            products.unshift(product);
            collapseKey.push(product.product_name);
          }
        }
        setServiceGroup(serviceGroup);
        setServiceList(products);
        setCollapseKey(collapseKey);
      } else {
        message.error(msg);
      }
    });
  }

  // 选择产品
  function handleProductChange(parentProductName: string) {
    setReqParams({
      ...reqParams,
      parentProductName,
    });
  }

  // 服务名称搜索
  function handleServiceNameSearch(serviceName: string) {
    setReqParams({
      ...reqParams,
      serviceName: serviceName || 'all',
    });
  }

  // 点击表格行进行的操作
  function onTableRowClick(e, service) {
    Cookies.set(COOKIES.NAMESPACE, namespace);
    Cookies.set('em_product_name', service.product_name);
    utils.setNaviKey('menu_ops_center', 'sub_menu_service');
    props.history.push(
      `/opscenter/service?component=${service.product_name}&service_group=${e.group}&service=${e.service_name}`
    );
  }

  return (
    <div className="namespace-overview-page">
      <div className="header-container mb-12">
        <span className="mr-20">
          <span>产品名称：</span>
          <Select
            placeholder="请选择产品"
            style={{ width: 264 }}
            value={
              reqParams.parentProductName === 'all'
                ? undefined
                : reqParams.parentProductName
            }
            onChange={handleProductChange}>
            {parentProductList.map((name: string) => (
              <Option key={name} value={name}>
                {name}
              </Option>
            ))}
          </Select>
        </span>
        <Search
          placeholder="输入服务名称进行搜索"
          style={{ width: 264 }}
          onSearch={handleServiceNameSearch}
        />
      </div>
      {Array.isArray(serviceList) && serviceList.length ? (
        <ServiceOverview
          className="c-namespace__overview"
          activeKey={collapseKey}
          onChange={(key) => setCollapseKey(key)}
          serviceList={serviceList}
          serviceGroup={serviceGroup}
          serviceNameRender={(value: string, record: any) => (
            <React.Fragment>
              <span>{value}</span>
              <Tag className="c-overview__ant-tag ml-20">{record.version}</Tag>
            </React.Fragment>
          )}
          onTableRowClick={onTableRowClick}
        />
      ) : (
        <div className="overview-empty">
          <Empty className="c-ant_empty-center" />
        </div>
      )}
    </div>
  );
};
export default OverviewPage;
