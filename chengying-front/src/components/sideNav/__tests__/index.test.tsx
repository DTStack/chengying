import * as React from 'react';
import { createMemoryHistory } from 'history';
import { Router } from 'react-router-dom';
import SideNav from '../index';
import { renderWithRedux } from '@/utils/test';
import { authorityList } from '@/mocks';
import reducer from '@/stores';
import { navData } from '@/constants/navData';
import {
  fireEvent,
  RenderResult,
  cleanup,
  screen,
} from '@testing-library/react';

const defaultProps = {
  selectKeys: ['sub_menu_diagnose_log'],
  menuList: navData[0].children,
  defaultOpenKey: ['sub_menu_diagnose'],
  authorityList,
  history: {},
};

describe('sidenav render', () => {
  let wrapper: RenderResult;
  beforeEach(() => {
    const history = createMemoryHistory();
    const layout = (
      <Router history={history}>
        <SideNav {...defaultProps} />
      </Router>
    );
    wrapper = renderWithRedux(layout, reducer, {
      UserCenterStore: { authorityList },
    });
  });

  afterEach(cleanup);

  // menuList
  test('menuList render', () => {
    const menuItems = wrapper.container.getElementsByClassName('ant-menu-item');
    const menuItemNames = ['概览', '服务', '主机'];
    for (const i in menuItems) {
      if (menuItems[i] instanceof HTMLElement) {
        expect(menuItems[i].textContent).toBe(menuItemNames[i]);
      }
    }
  });

  // submenu
  // test('submenu render', () => {
  //     screen.debug();
  //     const subMenu = wrapper.container.getElementsByClassName('ant-menu-sub');
  //     expect(subMenu.container.getElementsByClassName('ant-menu-item')[0].textContent).toBe('日志查看');
  //     // defaultOpenKey
  //     const subMenuLi = subMenu.parentElement;
  //     expect(subMenuLi.className).toMatch('ant-menu-submenu-open');
  //     // 收起
  //     fireEvent.click(subMenu);
  //     expect(subMenuLi.className).not.toMatch('ant-menu-submenu-open');
  // })

  test('sider snapshot', () => {
    expect(wrapper.asFragment()).toMatchSnapshot();
  });
});
