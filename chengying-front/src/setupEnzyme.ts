import { configure } from 'enzyme';
import EnzymeAdapter from 'enzyme-adapter-react-16';
// tslint:disable
const _ = require('lodash');
const originalConsoleError = console.error;

configure({ adapter: new EnzymeAdapter() });
// 解决jsdom parse css报错问题
console.error = function (msg: any) {
  // if (_.startsWith(msg, '[vuex] unknown')) return
  if (_.startsWith(msg, 'Error: Could not parse CSS stylesheet')) {
    return;
  }
  originalConsoleError(msg);
};
