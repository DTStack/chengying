import * as React from 'react';
import './style.scss';

interface IProps {
  className?: string;
  style?: React.CSSProperties;
  title: string;
  content?: string;
  imgSrc: any;
  handleTypeClick?: () => void;
  onFocus?: () => void;
  onBlur?: () => void;
}

export default class ChoiceCard extends React.PureComponent<IProps, any> {
  render() {
    const { title, content, imgSrc, className, style } = this.props;
    return (
      <div
        className={`c-card__container ${className || ''}`}
        style={style}
        onClick={this.props.handleTypeClick || null}
        tabIndex={0}
        onFocus={this.props.onFocus || null}
        onBlur={this.props.onBlur || null}>
        <img src={imgSrc} width={64} height={64} />
        <div className="card-info">
          <p className="text-title-bold">{title}</p>
          <p className="text-content">{content}</p>
        </div>
        <i className="emicon emicon-enter" />
      </div>
    );
  }
}
