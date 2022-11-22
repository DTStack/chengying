import * as React from 'react';
import { BrowserRouter, Route, Switch, Redirect } from 'react-router-dom';
import { RouterConf, RouterConfItemType } from './routerConf';
import { alertModal } from '@/utils/modal';
import { connect } from 'react-redux';
import { AppStoreTypes } from '@/stores';
const mapStateToProps = (state: AppStoreTypes) => ({
  runtimeState: state.InstallGuideStore.runtimeState,
  deployState: state.InstallGuideStore.deployState,
});
function getRouterConf(routerConf: RouterConfItemType[]) {
  const myRoutes = [];
  const routerConfFormat = (
    fPath: string,
    routerConf: any[],
    routeLists: any[]
  ) => {
    routerConf.forEach((r) => {
      const { path, layout } = r;
      const childrenRouters = getChildrenRoutes(r);
      if (layout) {
        const routeSwitch = <Switch>{childrenRouters}</Switch>;
        const routeLayout = (
          <Route
            key={path}
            path={path}
            render={(props: any) =>
              React.createElement(layout, props, routeSwitch)
            }></Route>
        );
        addRouteElement(path, routeLists, routeLayout);
        return;
      }
      if (childrenRouters.length) {
        routeLists.push(...childrenRouters);
        return;
      }
      const newPath = getNewPath(fPath, path);
      const routeElement = createRouteElement(r, newPath);
      addRouteElement(path, routeLists, routeElement);
    });
  };
  routerConfFormat('/', routerConf, myRoutes);
  return myRoutes;

  // 顺序处理
  function addRouteElement(
    path: string,
    routeLists: any[],
    routeElement: React.ReactNode
  ) {
    path === '/' || path === '*'
      ? routeLists.push(routeElement)
      : routeLists.unshift(routeElement);
  }

  // 获取子组件列表
  function getChildrenRoutes({ path, children }) {
    const routes = [];
    if (Array.isArray(children)) {
      routerConfFormat(path, children, routes);
    }
    return routes;
  }

  // 获取路径
  function getNewPath(fPath: string, path: string) {
    return fPath === '/' ? path : path === '/' ? fPath : fPath + path;
  }

  // 创建路由
  function createRouteElement(
    { redirect, path, layout, component, children }: RouterConfItemType,
    newPath: string
  ) {
    let routeElement;
    if (redirect) {
      routeElement = (
        <Redirect exact key={redirect} from={path} to={redirect}></Redirect>
      );
    } else {
      routeElement = (
        <Route exact key={newPath} path={newPath} component={component}></Route>
      );
    }
    return routeElement;
  }
}

@(connect(mapStateToProps) as any)
class Routers extends React.Component<any> {
  getConfirmation = (message: any, callback: any) => {
    const { runtimeState, deployState } = this.props;
    alertModal(runtimeState, deployState, callback);
  };

  render() {
    return (
      <BrowserRouter getUserConfirmation={this.getConfirmation}>
        <Switch>{getRouterConf(RouterConf)}</Switch>
      </BrowserRouter>
    );
  }
}

export default Routers;
