import React, { useEffect, useState } from 'react';
import { Modal, Icon, message, Drawer } from 'antd';
import AceEditor from 'react-ace';
import { scriptManager } from '@/services';
import 'ace-builds/src-noconflict/mode-golang';
import 'ace-builds/src-noconflict/mode-powershell';
import 'ace-builds/src-noconflict/theme-kuroir';
import * as _ from 'lodash';
import './style.scss'

interface IProps {
    id: number;
    title: string;
    visible: boolean;
    close: () => void;
    content: string;
    showUploadModal: (detailInfo: any, id: any) => void
}

const PreviewScript: React.FC<IProps> = (props) => {
    const {title, visible, close, id, showUploadModal} = props;
    const [content, setContent] = useState('')
    const [detailInfo, setDetailInfo] = useState({
        describe: '',
        execTimeout: '',
        logRetention: ''
    })
    useEffect(() => {
        getDetailContent()
    }, [props.id])

    const getDetailContent = async () => {
        const res = await scriptManager.getTaskContent({ id: id });
        const { data = {} } = res;
        if (data.code === 0) {
            setContent(data.data.script_content)
            setDetailInfo({
                describe: data.data.describe,
                execTimeout: data.data.exec_timeout,
                logRetention: data.data.log_retention
            })
          } else {
            message.error(data.msg);
          }
    }

    return (
        <Drawer
            title={title}
            visible={visible}
            onClose={close}
            width='520px'
            placement="right"
        >
            <div className='settingNavTop' style={{marginBottom: 20}}>
                <div className='navTopLine'></div>
                <div className='navTopTxt'>
                    基本信息
                    <Icon type="edit" style={{color: '#3F87FF', marginLeft: 8}} onClick={() => showUploadModal(detailInfo, id)}/>
                </div>
            </div>
            <div className='desBox'>
                <div className='desBox-title flex-item'>脚本描述：</div>
                <div className='desBox-content flex-item'>{detailInfo?.describe ? detailInfo?.describe : '--'}</div>
            </div>
            <div className='desBox'>
                <div className='desBox-title flex-item'>超时设置：</div>
                <div className='desBox-content flex-item'>{detailInfo?.execTimeout}s</div>
            </div>
            <div className='desBox'>
                <div className='desBox-title flex-item' style={{width: 108}}>执行历史保存周期：</div>
                <div className='desBox-content flex-item'>{detailInfo?.logRetention}天</div>
            </div>
            <div className='settingNavTop' style={{marginBottom: 20}}>
                <div className='navTopLine'></div>
                <div className='navTopTxt'>查看脚本</div>
            </div>
            <AceEditor
                className="ace-code-portal"
                mode="golang"
                theme="kuroir"
                value={content}
                readOnly={true}
                width="520"
                showGutter={false}
            />
        </Drawer>
    )
}

export default PreviewScript;

