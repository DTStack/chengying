import * as React from 'react';

interface State {
  isFullScreen: boolean;
}
class FullScreen extends React.Component<any, State> {
  state: State = {
    isFullScreen: false,
  };

  componentDidMount() {
    this.watchFullScreen();
  }

  handleFullScreen = () => {
    console.log('fullscreen:', this.state.isFullScreen);
    if (!this.state.isFullScreen) {
      this.requestFullScreen();
    } else {
      this.setState({ isFullScreen: false });
      this.exitFullscreen();
    }
  };

  // 进入全屏
  requestFullScreen = () => {
    console.log('requestFullScreen');
    this.setState({ isFullScreen: true });
    const de: any = document.getElementById(this.props.idName);
    if (de.requestFullscreen) {
      de.requestFullscreen();
    } else if (de.mozRequestFullScreen) {
      de.mozRequestFullScreen();
    } else if (de.webkitRequestFullScreen) {
      de.webkitRequestFullScreen();
    }
  };

  // 退出全屏
  exitFullscreen = () => {
    var de: any = document;
    if (de.exitFullscreen) {
      de.exitFullscreen();
    } else if (de.mozCancelFullScreen) {
      de.mozCancelFullScreen();
    } else if (de.webkitCancelFullScreen) {
      de.webkitCancelFullScreen();
    }
  };

  // 监听事件
  watchFullScreen = () => {
    const _self = this;
    const de: any = document;
    document.addEventListener(
      'fullscreenchange',
      function () {
        _self.setState({
          isFullScreen: de.fullscreen,
        });
      },
      false
    );
    document.addEventListener(
      'webkitfullscreenchange',
      function () {
        _self.setState({
          isFullScreen: de.webkitIsFullScreen,
        });
      },
      false
    );
    document.addEventListener(
      'mozfullscreenchange',
      function () {
        _self.setState({
          isFullScreen: de.mozFullScreen,
        });
      },
      false
    );
  };

  render() {
    const { isFullScreen } = this.state;
    return !isFullScreen ? (
      <span onClick={this.handleFullScreen}>
        <img src={require('./img/fullScreen.png')} style={{ width: '40px' }} />
      </span>
    ) : (
      <span onClick={this.handleFullScreen}>
        <img src={require('./img/smallScreen.png')} style={{ width: '40px' }} />
      </span>
    );
  }
}

export default FullScreen;
