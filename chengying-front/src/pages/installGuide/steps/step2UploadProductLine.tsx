import * as React from 'react';
import { Modal, Upload, message, Table, Button, Popover } from 'antd';
import { productLine } from '@/services';
interface IProps{
    callList: (id?: any) => void;
    visible: boolean;
    dataList: any;
    onCancel: () => void;
}

let selectRecord = null

const UploadProductLine: React.FC<IProps> = (props) => {
    const {visible, callList, onCancel, dataList} = props;
    const [popoverVisible, setPopoverVisible] = React.useState(false)

    React.useEffect(() => {
      // console.log(props.dataList)
    }, [])

    const doDelete = (record) => {
        productLine.deleteProductLine({id: record.id}).then((res) => {
            if (res.data.code == 0) {
                message.success('操作成功！')
                setPopoverVisible(false)
                callList(record.id)
            } else {
                message.error(res.data.msg)
            }
        })
    }

    const customRequest = async (options) => {
      const formData = new FormData();
      formData.append('file', options.file);
      const res = await productLine.uploadProductLine(formData)
      if (res.data.code == 0) {
        message.success('上传成功！');
        callList()
        // 刷新列表
      } else {
        message.error(res.data.msg)
      }

    }

    const changeVisible = (visible, record) => {
      setPopoverVisible(visible)
      selectRecord = record
    }

    const content = (
      <div className='propBox'>
        <Button size='small' style={{marginRight: 12}} onClick={() => changeVisible(false, null)}>取消</Button>
        <Button  size='small'  type="primary" onClick={() => doDelete(selectRecord)}>确定</Button>
      </div>
    )

    const uploadProductLine = {
        name: 'file',
        customRequest: (options) => customRequest(options),
        accept: '.json',
        fileList: []
    }

    const columns = [
        {
            title: '产品线名称',
            dataIndex: 'product_line_name',
            key: 'product_line_name',
            render: (text, record) => (<span>{text} （{record.product_line_version}）</span>),
        },
        {
            width: 80,
            title: '操作',
            key: 'action',
            className: 'actionName',
            render: (text, record) => (
              <Popover
                visible={popoverVisible && selectRecord?.product_line_version === record.product_line_version && selectRecord?.product_line_name === record.product_line_name} 
                placement="top" 
                title='确定删除产品线?' 
                content={content} 
                onVisibleChange={(visible) => changeVisible(visible, record)}
                trigger="click">
                <Button type="link">删除</Button>
            </Popover>
            ), 
        }
    ]


    return (
        <Modal
            visible={visible}
            onCancel={onCancel}
            className='uploadProductLineModal'
            title='上传产品线'
            footer={null}
        >
             <Upload {...uploadProductLine}>
                <Button icon="upload" size='small'>上传产品线</Button>
                <span className='uploadTxt'>仅支持.json格式文件</span>
              </Upload>
              <div className='uploadProductTable'>
                <Table 
                  scroll={{ y: 260 }}
                  columns={columns} 
                  dataSource={dataList} 
                  pagination={false} 
                  size='small'/>
              </div>
        </Modal>
    )
}

export default UploadProductLine;