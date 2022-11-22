export const RESULT_STATUS = {
  UNNORMAL: 3,
  NORMAL: 2,
  RUN: 1,
  UNRUN: 0,
};
export const RESULT_FILTER = [
  { text: '正常', value: RESULT_STATUS.NORMAL },
  { text: '异常', value: RESULT_STATUS.UNNORMAL },
  { text: '运行中', value: RESULT_STATUS.RUN },
  { text: '未运行', value: RESULT_STATUS.UNRUN },
];

export const RESULT_FILTER_HISTORY = [
  { text: '正常', value: RESULT_STATUS.NORMAL },
  { text: '异常', value: RESULT_STATUS.UNNORMAL },
  { text: '运行中', value: RESULT_STATUS.RUN },
];
