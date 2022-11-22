import React from 'react';
import {connect} from 'react-redux';
import { InputNumber, Form, Button, message, Spin } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { AppStoreTypes } from '@/stores';
import './style.scss';
const FormItem = Form.Item;

import { globalConfig } from '@/services';

interface Prop extends FormComponentProps{
    authorityList?: any;
}
interface State {  
    autolimitNumber:any;
    limitNumber: any;
    loading: boolean;
}
const mapStateToProps = (state: AppStoreTypes) => ({
    authorityList: state.UserCenterStore.authorityList,
 });
 @(connect(mapStateToProps) as any)
class GlobalConfig extends React.Component<Prop, State>{
    state: State = {
        autolimitNumber:null,
        limitNumber: null,
        loading: true
    }

    changeNumber = (value: any) => {
        this.setState({limitNumber: value})
    }

    changeAutoNumber = (value: any) => {
        this.setState({autolimitNumber:value})
    }

    doSave = () => {
        const { limitNumber,autolimitNumber} = this.state;
        globalConfig.setGlobalConfig({service_install_timeout_limit: limitNumber,auto_test_timeout_limit:autolimitNumber}).then(res => {
            const { data } = res;
            if (data.code !== 0) {
                message.error(data.msg);
            } else {
                message.success('操作成功');
            }
        })
    }

    getLimitNumber = async () => {
        const { data } = await globalConfig.getGlobalConfig();
        console.log(data)
        if (data.code !== 0) {
           message.error(data.msg);
        } else {
            this.setState({
                autolimitNumber: +data.data.auto_test_timeout_limit,
                limitNumber: +data.data.service_install_timeout_limit 
            }, () => {
                this.setState({loading: false})
            })
        }
    }
  



    componentDidMount () {
        this.getLimitNumber()
    }


    render() {


        console.log("这是from属性:", this.props.form);
        
        const { getFieldDecorator, getFieldsError } = this.props.form;
        const isDisable = getFieldsError()?.limitNumber || !this.props.authorityList?.sub_menu_configuration_global_config_edit
        const { limitNumber, autolimitNumber,loading } = this.state;
        const formLayout = {
            labelCol: {
              xs: { span: 2 },
              sm: { span: 2 },
            },
            wrapperCol: {
              xs: { span: 24 },
              sm: { span: 15 },
            },
          };
           
        return (
            <div className='globalConfig'>
                <Spin  spinning={loading}> 
                {!loading && 
                <Form>
                <FormItem label='组件包安装超时 ' {...formLayout}>
                {getFieldDecorator('limitNumber', {
                    initialValue: 'limitNumber',
                    rules: [
                        { required: true, message: '请输入时长' },
                        { pattern: /^[1-9]*[1-9][0-9]*$/, message: '请输入大于0的正整数' },
                    ],
                    })(
                        <div>
                        <InputNumber 
                            defaultValue={limitNumber}
                            step={1}
                            min={1}
                            onChange={this.changeNumber}
                        />  分钟
                        </div> 
                    )}
                </FormItem>   
                <FormItem label='自动化测试超时 ' {...formLayout}>
                {getFieldDecorator('autolimitNumber', {
                    initialValue: 'autolimitNumber',
                    rules: [
                        { required: true, message: '请输入时长' },
                        { pattern: /^[1-9]*[1-9][0-9]*$/, message: '请输入大于0的正整数' },
                    ],
                    })(
                        <div>
                        <InputNumber 
                            defaultValue={autolimitNumber}
                            step={1}
                            min={1}
                            onChange={this.changeAutoNumber}
                        />  分钟
                        </div> 
                    )}
                </FormItem> 
                
                <FormItem>
                    <div className='btnList'>
                        <Button
                          onClick={this.doSave}  
                          htmlType="submit" 
                          disabled={isDisable ? true : false}>保存</Button>
                    </div>
                </FormItem>
                </Form>
                }
                </Spin>
            </div>
        )
    }
}
export default Form.create<Prop>()(GlobalConfig);