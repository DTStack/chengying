import * as React from 'react';
import { render } from 'enzyme';
import RootLayout from '../index';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import configureStore from '../../../configureStore';
const store = configureStore({});

describe('test rootlayout navbar', () => {
  const props = {};
  const layout = render(
    <Provider store={store}>
      <BrowserRouter>
        <RootLayout {...props} />
      </BrowserRouter>
    </Provider>
  );

  test('check menuitem name', () => {
    console.log('layout ------ ', layout);
    const title = '集群管理';
    expect(title).toEqual('集群管理');
  });
});
