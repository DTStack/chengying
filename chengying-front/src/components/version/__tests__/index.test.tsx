import * as React from 'react';
import Version from '../index';
import { render } from '@testing-library/react';
const version = JSON.stringify(require('../../../../package.json').version);

describe('version render', () => {
  test('current version', () => {
    const wrapper = render(<Version />);
    const versionComponent = wrapper.getByTestId('version-component');
    expect(versionComponent.textContent).toMatch(version);
  });
});
