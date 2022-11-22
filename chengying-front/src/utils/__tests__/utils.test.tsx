import Utils from '../utils';
// import { fireEvent, render, RenderResult, cleanup } from '@testing-library/react';

describe('Utils unit test', () => {
  test('percent unit', () => {
    expect(Utils.percent(2, 100)).toBe('100%');
  });

  test('getCssText unit', () => {
    expect(
      Utils.getCssText({
        'font-size': '14px',
        color: '#3f87ff',
      })
    ).toBe('font-size:14px;color:#3f87ff;');
  });

  test('formatGBUnit unit', () => {
    expect(Utils.formatGBUnit('1TB')).toBe(1024);
  });

  test('jsonToQuery unit', () => {
    expect(
      Utils.jsonToQuery({
        name: 'tt',
        id: 1,
      })
    ).toBe(`name=tt&id=1`);
  });

  test('formateDateTime unit', () => {
    expect(Utils.formateDateTime(1639380097000)).toBe(`2021-12-13 15:21:37`);
  });

  test('formateDate unit', () => {
    expect(Utils.formateDate(1639380097000)).toBe(`2021-12-13`);
  });

  test('trim unit', () => {
    expect(Utils.trim(` 1639380097000   `)).toBe(`1639380097000`);
  });

  test('checkNullObj unit', () => {
    expect(Utils.checkNullObj({})).toBe(true);
  });
});
