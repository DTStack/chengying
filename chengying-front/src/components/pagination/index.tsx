import * as React from 'react';
import { Icon } from 'antd';
import './style.scss';

interface Props {
  style?: any;
  handleClickTop?: () => void;
  handleClickUp?: () => void;
  handleClickDown?: () => void;
  handleClickNew?: () => void;
}
// interface State {

// }
class SpecialPagination extends React.Component<Props> {
  constructor(props: any) {
    super(props as Props);
  }

  render() {
    const {
      style,
      handleClickTop,
      handleClickUp,
      handleClickDown,
      handleClickNew,
    } = this.props;
    return (
      <div className="special_pagination" style={style}>
        <button onClick={handleClickTop}>TOP</button>
        <button onClick={handleClickUp}>
          <Icon type="up" />
        </button>
        <button onClick={handleClickDown}>
          <Icon type="down" />
        </button>
        <button onClick={handleClickNew}>NEW</button>
      </div>
    );
  }
}
export default SpecialPagination;
