import * as React from 'react';

const TimerInc = ({ initial }) => {
  const timerRef = React.useRef(null);
  const [timer, setTimer] = React.useState(initial);
  React.useEffect(() => {
    setTimer(initial);
    timerRef.current = setInterval(() => {
      setTimer((timer) => timer + 1);
    }, 1000);
    return () => clearInterval(timerRef.current);
  }, [initial]);
  return <div>{timer + 's'}</div>;
};

export default TimerInc;
