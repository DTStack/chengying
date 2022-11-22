
import React from 'react'

import { Form, InputNumber, message, Modal } from 'antd'
import { FormComponentProps } from 'antd/es/form'

import Api from '@/services/clusterInspectionService'

import './index.scss'

interface Prop extends FormComponentProps {
    visible: boolean;
    ModalClose: () => void
}
interface IState {
    isModalVisible: boolean;
    resParams: any;
}

class StatisticalSetting extends React.Component<Prop, IState>{

    state: IState = {
        isModalVisible: false,
        resParams: {
            fullGC_time: 1,
            fullGC_frequency: 0,
            dir_size: 0,
            node_cpu_usage: 0,
            node_mem_usage: 0,
            node_disk_usage: 0,
            node_inode_usage: 0,
        }
    }
    /**
     *  modal 设置
     */
    ModalOnOk = (e: any) => {
        // 表单验证,验证成功后退出
        this.props.form.validateFieldsAndScroll((err, values) => {
            if (!err) {
                this.updateSettingData()
            } else {
                // 校验失败
                e.preventDefault();
            }
        });
    }

    ModaOnCancel = () => {
        const { ModalClose } = this.props
        ModalClose()
    }

    changeItem = (value, key) => {
        const { resParams } = this.state;
        let obj = resParams
        // 如果是fullGC_time
        let initValue = value
        if (!value) {
            if (key === 'fullGC_time') {
                initValue = 1
            } else {
                initValue = 0
            }
        }
        obj[key] = initValue
        this.setState({resParams: obj})
    }
    componentDidMount() {
        this.getSetting()
    }

    updateSettingData = async () => {
        let res = await Api.getClusterInspectionStatisSetData(this.state.resParams)

        if (res.data.code == 0) {
            message.success('操作成功!')
            this.props.ModalClose()
        } else {
            message.error(res.msg)
        }
    }

    getSetting = async () => {
        const res = await Api.getSettingData()
        if (res.data.code == 0) {
            this.setState({resParams: res.data.data})
        } else {
            message.error(res.data.msg)
        }
    }


    render() {
        const { getFieldDecorator } = this.props.form
        const { visible } = this.props;
        const { resParams} = this.state;

        return (
            <>
                <Modal title="统计设置"
                    className='statisticalSetting'
                    visible={visible}
                    onOk={(e) => this.ModalOnOk(e)}
                    onCancel={this.ModaOnCancel}>
                        <Form>
                            <div className='settingNavTop'>
                                <div className='navTopLine'></div>
                                <div className='navTopTxt'>指标汇总</div>
                            </div>
                            <Form.Item label="服务Full GC统计">
                                <div style={{display: 'flex'}}>
                                <Form.Item style={{display: 'block'}}>
                                {getFieldDecorator('node_cpu_usage', {
                                        initialValue: resParams.node_cpu_usage,
                                        rules: [
                                            { pattern: /^[1-9]*$/, message: '请输入大于等于1的正整数' },
                                        ],
                                    })(
                                        <div>
                                            <span className='labelLeft'>最近</span>
                                            <InputNumber
                                                step={1}
                                                value={resParams.fullGC_time}
                                                onChange={(value) => this.changeItem(value, 'fullGC_time')}
                                            />分钟
                                        </div>
                                    )}
                                </Form.Item>
                                <Form.Item style={{display: 'block'}}>
                                {getFieldDecorator('fullGC_frequency', {
                                        initialValue: resParams.fullGC_frequency,
                                        rules: [
                                            { pattern: /^[0-9]*$/, message: '请输入大于等于1的正整数' },
                                        ],
                                    })(
                                        <div>
                                            <span className='labelLeft'>Full GC次数大于等于</span>
                                            <InputNumber
                                                step={1}
                                                value={resParams.fullGC_frequency}
                                                onChange={(value) => this.changeItem(value, 'fullGC_frequency')}
                                            />次
                                        </div>
                                    )}
                                </Form.Item>
                                </div>
                            </Form.Item>
                            <Form.Item label="/opt/dtstack 目录统计(最新)">
                                {getFieldDecorator('dir_size', {
                                    initialValue: resParams.dir_size,
                                    rules: [
                                        { pattern: /^[0-9]*$/, message: '请输入大于等于0的正整数' },
                                    ],
                                })(
                                    <div>
                                        <span className='labelLeft'>目录大小大于等于</span>
                                        <InputNumber
                                            step={1}
                                            value={resParams.dir_size}
                                            onChange={(value) => this.changeItem(value, 'dir_size')}
                                        />G
                                    </div>
                                )}
                            </Form.Item>
                            <div className='settingNavTop'>
                                <div className='navTopLine'></div>
                                <div className='navTopTxt'>节点运行</div>
                            </div>
                            <Form.Item label="CPU使用率统计">
                                {getFieldDecorator('node_cpu_usage', {
                                    initialValue: resParams.node_cpu_usage,
                                    rules: [
                                        { pattern: /^[0-9]*$/, message: '请输入大于等于0的正整数' },
                                    ],
                                })(
                                    <div>
                                        <span className='labelLeft'>CPU使用率大于等于</span>
                                        <InputNumber
                                            step={1}
                                            value={resParams.node_cpu_usage}
                                            onChange={(value) => this.changeItem(value, 'node_cpu_usage')}
                                        />%
                                    </div>
                                )}
                            </Form.Item>
                            <Form.Item label="内存使用趋势">
                                {getFieldDecorator('node_mem_usage', {
                                    initialValue: resParams.node_mem_usage,
                                    rules: [
                                        { pattern: /^[0-9]*$/, message: '请输入大于等于0的正整数' },
                                    ],
                                })(
                                    <div>
                                        <span className='labelLeft'>内存使用趋势大于等于</span>
                                        <InputNumber
                                            step={1}
                                            value={resParams.node_mem_usage}
                                            onChange={(value) => this.changeItem(value, 'node_mem_usage')}
                                        />%
                                    </div>
                                )}
                            </Form.Item>
                            <Form.Item label="磁盘使用趋势">
                                {getFieldDecorator('node_disk_usage', {
                                    initialValue: resParams.node_disk_usage,
                                    rules: [
                                        { pattern: /^[0-9]*$/, message: '请输入大于等于0的正整数' },
                                    ],
                                })(
                                    <div>
                                        <span className='labelLeft'>磁盘使用率大于等于</span>
                                        <InputNumber
                                            step={1}
                                            value={resParams.node_disk_usage}
                                            onChange={(value) => this.changeItem(value, 'node_disk_usage')}
                                        />%
                                    </div>
                                )}
                            </Form.Item>
                            <Form.Item label="inode使用率趋势">
                                {getFieldDecorator('node_inode_usage', {
                                    initialValue: resParams.node_inode_usage,
                                    rules: [
                                        { pattern: /^[0-9]*$/, message: '请输入大于等于0的正整数' },
                                    ],
                                })(
                                    <div>
                                        <span className='labelLeft'>inode使用率大于等于</span>
                                        <InputNumber
                                            step={1}
                                            value={resParams.node_inode_usage}
                                            onChange={(value) => this.changeItem(value, 'node_inode_usage')}
                                        />%
                                    </div>
                                )}
                            </Form.Item>
                        </Form>
                </Modal>
            </>
        )
    }
}

export default Form.create<Prop>({})(StatisticalSetting)