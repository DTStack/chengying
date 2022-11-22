import * as React from 'react';
import { Progress } from 'antd';
import './style.scss';

const multiMapping = {
  B: 1,
  KB: 1024,
  MB: 1024 ** 2,
  GB: 1024 ** 3,
  TB: 1024 ** 4,
};

function formatNum(str: string | number): number {
  if (typeof str === 'string') {
    const regexExec = /(B|KB|MB|GB|TB)/.exec(str);
    const unit = regexExec ? regexExec[0] : null;
    const multi = unit ? multiMapping[unit] : 1;
    const num = str ? numNaN(parseFloat(str)) : 0;
    return multi * num;
  } else {
    return str ? numNaN(str) : 0;
  }
}

function numNaN(num) {
  return isNaN(num) ? 0 : num;
}

interface ProgressBarProps {
  unit?: string;
  percent?: number;
  now: number | string;
  total: number | string;
  pStyle?: any;
}

export default function ProgressBar(props: ProgressBarProps) {
  let { now, total, unit, pStyle } = props;
  const nowNum = formatNum(now);
  const totalNum = formatNum(total);

  const percent =
    props.percent !== undefined
      ? props.percent
      : totalNum
      ? (nowNum / totalNum) * 100
      : 0;

  if (unit) {
    now = `${now}${unit}`;
    total = `${total}${unit}`;
  }

  return (
    <div className="progress-bar">
      <p style={pStyle}>{`${now} / ${total}`}</p>
      <Progress
        strokeColor="#3F87FF"
        strokeWidth={4}
        percent={percent}
        showInfo={false}
      />
    </div>
  );
}
