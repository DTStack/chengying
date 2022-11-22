import * as React from 'react';

class TextMark extends React.Component<any, any> {
  public renderMark(text = '', markText = '') {
    const reg = new RegExp(markText, 'g');
    return text.replace(reg, `<span style="color: #2391F7">${markText}</span>`);
  }

  public render() {
    const { text, markText, ...others } = this.props;
    return (
      <span
        dangerouslySetInnerHTML={{ __html: this.renderMark(text, markText) }}
        {...others}></span>
    );
  }
}

export default TextMark;
