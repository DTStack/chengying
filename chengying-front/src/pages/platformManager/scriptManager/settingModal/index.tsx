import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, Button, message } from 'antd';
import HostSelect from './hostSelect';
import { scriptManager } from '@/services';
import './style.scss'

const FormItem = Form.Item
const { TextArea } = Input

const TimeSetting: React.FC = (props: any) => {
    const { title, isTimeSetting, visible, close, spec, hosts, id, sucCall } = props;
    const [isError, setIsError] = useState(false)
    const [hostList, setHostList] = useState([])
    const [isDisabled, setIsDisabled] = useState(true)
    const [specValue, setSpecValue] = useState(spec)
    const [errorMsg, setErrorMsg] = useState('')
    const [timeList, setTimeList] = useState([])
    const [isErrorHost, setIsErrorHost] = useState(false)
    const [isLoading, setIsLoading] = useState(false)

    const checkCorn = async() => {
        if (specValue === '') {
            setIsError(true)
            setErrorMsg('cron表达式不能为空')
        }
        const param = {
            id,
            spec: specValue,
            next: 10
        }
        const res = await scriptManager.checkCronParse(param)
        if (res.data.code == 0) {
            setTimeList(res.data?.data.next_time)
            setIsError(false)
            setIsDisabled(false)
        } else {
            setIsDisabled(true)
            setIsError(true)
            setTimeList([])
            setErrorMsg(res.data.msg)
        }
    }

    useEffect(() => {
        if (spec.length > 0 && isTimeSetting) {
            checkCorn()
        }
    }, [spec, isTimeSetting])

    // 定时设置弹框
    const doTimeSetting = async () => {
        const hostIds = hostList.join(',')
        const param = {
         id,
         'host_id': hostIds,
          spec: specValue,
        }
        const res = await scriptManager.getSettingStatus(param)
        if (res.data.code == 0) {
            setIsLoading(false)
            close()
            sucCall()
            message.success(res.data.msg)
        } else {
            setIsLoading(false)
            message.error(res.data.msg)
        }
    }

    // 手动执行
    const handSetting = async () => {
        const hostIds = hostList.join(',')
        const param = {
         id,
         'host_id': hostIds
        }
        const res = await scriptManager.taskRun(param)
        if (res.data.code == 0) {
            setIsLoading(false)
            close()
            sucCall()
            message.success(res.data.msg)
        } else {
            setIsLoading(false)
            message.error(res.data.msg)
        }
    }
    // 执行确定操作
   const doSetting = () => {
    if (hostList.length === 0) {
        setIsErrorHost(true)
        return
    }
    setIsLoading(true)
      if (isTimeSetting) {
        doTimeSetting()
      } else {
        handSetting()
      }
   }

   // 执行取消操作
   const cancelConfirm = () => {
    close()
   }

   // 选中的主机
   const selectHosts = (list) => {
       if (list.length > 0) {
        setIsErrorHost(false)
       }
       setHostList(list)
   }

   // cron语句有改变
   const changeCron = (e) => {
       const { value } = e.target;
       const { spec } = props;
       if (spec.trim() !== value.trim()) {
           setIsDisabled(true)
       } else {
        setIsDisabled(false) 
       }
       setSpecValue(value)
   }

   // 校验cron表达式
   const doCheck = () => {
    checkCorn()
   }

    return (
        <Modal
            title={title}
            visible={visible}
            onOk={() => doSetting()}
            okButtonProps={{disabled: isDisabled && isTimeSetting}}
            onCancel={() => cancelConfirm()}
            width='750px'
            confirmLoading={isLoading}
        >
             <Form>
                <FormItem>
                    <>
                        <HostSelect selectHosts={selectHosts} defaultHosts={hosts}/>
                        {isErrorHost &&<span className='errorHost'>请选择下发主机</span>}
                    </>
                </FormItem>
                {isTimeSetting && (
                    <FormItem>
                        <div>
                            <div className='settingNavTop'>
                                <div className='navTopLine'></div>
                                <div className='navTopTxt'>执行频率</div>
                            </div>
                            <div className='cronBox'>
                                <span className='labelTxt'>cron表达式：</span>
                                <Input 
                                    maxLength={500}
                                    placeholder="请输入cron表达式" 
                                    style={{ width: '599px'}} 
                                    onChange={changeCron}
                                    value={specValue}/>
                                {isError && (<span className='error'>{errorMsg}</span>)}
                            </div>
                        </div>
                     <FormItem>
                        <div className='runBox'>
                            <span className='labelTxt'>近10次运行时间：</span>
                            <TextArea 
                                value={timeList.join("\n")}
                                placeholder="若校验结果有误或为空，请核对 cron 表达式是否正确" 
                                style={{ width: '599px', height: 95}}/>
                        </div>
                     </FormItem>
                    </FormItem>
                )}
             </Form>
             {isTimeSetting && (
                <Button 
                    className='cronBtn'
                    type="primary" 
                    disabled={spec.length > 0 && spec.trim() == specValue.trim()} 
                    onClick={doCheck}>校验cron表达式</Button>)}
        </Modal>
    )
}

export default Form.create<any>()(TimeSetting)