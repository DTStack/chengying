import * as React from 'react';
import { Upload, Icon, message } from 'antd';
import { isEqual } from 'lodash';
import './style.scss';
const { Dragger } = Upload;

export interface DraggerUploadProps {
  onChange: Function;
  icon?: React.ReactNode;
  defaultFileList?: any[];
}

// 深度比较
const useDeepEffect = (callback: Function, deps: any) => {
  const isFirst = React.useRef(true); // 是否是第一次进来
  const prevDepsRef = React.useRef(deps); // 前一次props

  React.useEffect(() => {
    const isFirstEffect = isFirst.current;
    const prevDeps = prevDepsRef.current;
    const isSame = isEqual(deps, prevDeps);

    isFirst.current = false;
    prevDepsRef.current = deps;

    if (isFirstEffect || !isSame) {
      return callback();
    }
  }, deps);
};

const DraggerUpload: React.FC<DraggerUploadProps> = (props) => {
  const { onChange, icon, defaultFileList } = props;
  const [fileList, setFileList] = React.useState<any[]>([]);
  const uploadProps = {
    name: 'file',
    fileList: fileList,
    beforeUpload: (file: any) => {
      const sizeLimit = 50 * 1024 * 1024; // 50MB
      if (file?.size > sizeLimit) {
        message.error('本地上传文件不可超过50MB!');
        onFileListChange([]);
        return false;
      }
      const fileList = [file];
      onFileListChange(fileList);
      console.log(file);
      return false;
    },
    onRemove: () => {
      onFileListChange([]);
    },
  };

  useDeepEffect(() => {
    setFileList(defaultFileList);
  }, [defaultFileList]);

  function onFileListChange(fileList: any[]) {
    setFileList(fileList);
    onChange(fileList);
  }

  return (
    <Dragger className="c-dragger_ant-upload" {...uploadProps}>
      <p className="ant-upload-drag-icon">
        {icon || <Icon type="cloud-upload" />}
      </p>
      <p className="ant-upload-text">点击或将文件拖拽到此处上传</p>
    </Dragger>
  );
};
export default DraggerUpload;
