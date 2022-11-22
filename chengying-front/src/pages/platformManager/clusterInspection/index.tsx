import React from 'react';
import moment from 'moment';
import JsPDF from 'jspdf';
const html2canvas = require('html2canvas');
import {
    Select,
    Button,
    Row,
    Col,
    Descriptions,
    message,
    DatePicker,
} from 'antd'
const { Option } = Select;

import './index.scss'

import Api from '@/services/clusterInspectionService'
import StatisticalSetting from './components/statisticalSetting'
import { form_Title_info, chartTable } from './constants'
import MoreCharts from './components/moreChart'
import MoreTable from './components/MoreTable'


import { clusterManagerService } from '@/services';

const { RangePicker } = DatePicker

interface IState {
    clusterOptionList: any;
    node_status: any;
    service_status: any;
    mysql_slave_status: any;
    modalShow: any;
    clusterTime: any;
    currentClusterId: any;
    isHasNameNode: boolean;
    configList: any[];
    bigData: any;
    pdfLoading: boolean;
    pdfTitle: string;
    clusterName: string;
    showTable: boolean;
}

class ClusterInspection extends React.Component {
    constructor(props) {
        super(props)

    }

    state: IState = {
        clusterOptionList: [],
        node_status: {},
        service_status: { },
        mysql_slave_status: '',
        currentClusterId: null,//选中时当前的节点id! 
        clusterTime: [],
        modalShow: false,
        isHasNameNode: false,
        bigData: {},
        configList: [],
        pdfLoading: false,
        pdfTitle: '立即下载',
        clusterName: '',
        showTable: false
    }

    // 获取表格数据
    

    componentDidMount() {
        const end = moment();
        const start = moment().subtract(6, 'days');
        this.setState({clusterTime: [start, end],  showTable: true}, () => {
            this.getClusterListFn()
        })
    }
    // 生成PDF文档
    generatePdf = () => {
        const { clusterTime, clusterName } = this.state;
        return new Promise(function (resolve, reject) {
            html2canvas(document.getElementById('inspection-report-body')).then(function (canvas) {
                var contentWidth = canvas.width;
                var contentHeight = canvas.height;
    
                // 一页pdf显示html页面生成的canvas高度;
                var pageHeight = contentWidth / 595.28 * 841.89;
                // 未生成pdf的html页面高度
                var leftHeight = contentHeight;
                // pdf页面偏移
                var position = 0;
                // a4纸的尺寸[595.28,841.89]，html页面生成的canvas在pdf中图片的宽高
                var imgWidth = 555.28;
                var imgHeight = 555.28 / contentWidth * contentHeight;
    
                var pageData = canvas.toDataURL('image/jpeg', 1.0);
    
                var pdf = new JsPDF('p', 'pt', 'a4');
                // 有两个高度需要区分，一个是html页面的实际高度，和生成pdf的页面高度(841.89)
                // 当内容未超过pdf一页显示的范围，无需分页
                if (leftHeight < pageHeight) {
                    pdf.addImage(pageData, 'JPEG', 20, 0, imgWidth, imgHeight);
                } else {
                    while (leftHeight > 0) {
                        pdf.addImage(pageData, 'JPEG', 20, position, imgWidth, imgHeight)
                        leftHeight -= pageHeight;
                        position -= 841.89;
                        // 避免添加空白页
                        if (leftHeight > 0) {
                            pdf.addPage();
                        }
                    }
                }
                pdf.save(`${clusterName}集群巡检(${clusterTime[0].format('YYYY-MM-DD')}~${clusterTime[1].format('YYYY-MM-DD')}).pdf`);
                setTimeout(() => {
                    resolve(true)
                }, 2500)
            });
        });
    }

    // 下载PDF
     exportPDF = () => {
         this.setState({
             pdfTitle: '报告生成中',
             pdfLoading: true
         })
         this.generatePdf().then((res) => {
             if (res) {
                this.setState({
                    pdfTitle: '报告已生成',
                    pdfLoading: false
                })
             }
         })
      }
    reloadData = () => {
        const end = moment();
        const start = moment().subtract(6, 'days');
        this.setState({
            clusterTime: [start, end], 
            showTable: false,
            currentClusterId: this.state.clusterOptionList[0].id}, () => {
             // 获取状态统计的数据
             this.getStateDataApi()
             // 获取表格的数据
             this.getBaseInfo()
        })
    }

    getClusterBigDataServer = async() => {
        const res = await Api.getClusterBigDataServer({cluster_id: this.state.currentClusterId})
        if (res.data.code == 0) {
            this.setState({bigData: res.data.data})
        } else {
            message.error(res.data.msg)
        }
    }

    getBaseInfo = async() => {
        const res = await Api.getClusterInspectionBaseInfoData({cluster_id: this.state.currentClusterId})
        if (res.data.code == 0) {
            this.setState({
                node_status: res.data.data.node_status,
                service_status: res.data.data.service_status,
                mysql_slave_status: res.data.data.mysql_slave_status,
                isHasNameNode: res.data.data.have_name_node == 1 ? true : false,
                showTable: true
            })
            if (res.data.data.have_name_node == 1) {
                this.getClusterBigDataServer()
            }
        } else {
            message.error(res.data.msg)
        }
    }
    /**
     * 同步数据的请求
     */
    // 
    getStateDataApi = async () => {
        const res = await Api.getApplicationService({cluster_id: this.state.currentClusterId})
        if (res.data.code == 0) {
            const moduleArr = [];
            const chartConfigArr = [];
            const moduleTypes = [];
            res.data.data.reduce((pre, cur, index, arr) => {
              const obj = {
                noData: false,
                ...cur,
                moduleName: moduleArr.includes(cur.module) ? '' : cur.module,
                moduleType: moduleTypes.includes(cur.type) ? '' : cur.type,
              };
              moduleArr.push(cur.module);
              moduleTypes.push(cur.type)
              chartConfigArr.push(obj);
              return cur;
            }, {});
            console.log(chartConfigArr)
            // 同步数据
            this.setState({
                configList: chartConfigArr
            })
        } else {
            message.error(res.data.msg)
        }
    }

    // 获取集群列表!
    getClusterListFn = () => {
        const params = {
            type: '',
            'sort-by': 'id',
            'sort-dir': 'desc',
            limit: 0,
            start: 0,
        };
        clusterManagerService.getClusterLists(params).then((res: any) => {
            if (res.data.code === 0) {
                this.setState({
                    clusterOptionList: res.data.data.clusters || [],
                    currentClusterId: res.data.data.clusters[0].id || '',
                    clusterName: res.data.data.clusters[0].name || ''
                }, () => {
                    // 获取状态统计的数据
                    this.getStateDataApi()
                    // 获取表格的数据
                    this.getBaseInfo()
                });

            } else {
                message.error(res.data.msg)
            }
        });
    }


    /**
     * 时间设置
     */
    clusterTime = (date?: any) => {
        this.setState({clusterTime: date}, () => {
            this.getStateDataApi()
        })
    }

    // 将时间转转化为时间戳
    clusterTimeFormat = (str) => {
        let initStr = str;
        let newStr = initStr.split('-').join('/');
        return (new Date(newStr).getTime())
    }

    /**
     * 
     *  子组件统计设置
     */
    statisticalSetting = (value?: any) => {

        this.setState({
            modalShow: value
        })
    }
    modalClose = () => {
        this.setState({
            modalShow: false
        })
    }
    // 处理select的数据
    handleSelectFn = (key: number, options: any) => {
        this.setState({
            currentClusterId: key,
            clusterName: options.props.children,
            showTable: false
        }, () => {
            this.getStateDataApi()
            this.getBaseInfo()
        })
    }

    delNoDataChart = (index: any) => {
        const arr = [...this.state.configList];
        arr[index].noData = true;
        this.setState({configList: arr})
      };

    render() {
        const { 
            modalShow,
            service_status,
            node_status, 
            mysql_slave_status, 
            clusterOptionList,
            isHasNameNode,
            configList,
            clusterTime,
            currentClusterId,
            bigData,
            pdfLoading,
            pdfTitle,
            showTable
        } = this.state;

        const disabledDate = (current) => {
            const range = 14 * 24 * 60 * 60 * 1000;
            return (
              (current && current.valueOf() > Date.now()) ||
              (Date.now() - range > current && current.valueOf())
            );
        };

        return (
            <div>
                <header className='cluster-header'>
                    <div className='button-setting'>
                        <Button type='primary' 
                            loading={pdfLoading} 
                            onClick={this.exportPDF}>
                            {pdfTitle}
                        </Button>
                        {/* 获取子组件值,并设置为true. */}
                        <Button onClick={() => this.statisticalSetting(true)}>
                            统计设置
                        </Button>
                        <Button icon='reload' onClick={this.reloadData}>
                            刷新
                        </Button>
                    </div>
                    <span className='spanLeft'>集群:</span>
                    <Select
                        defaultValue={currentClusterId}
                        showSearch
                        value={currentClusterId}
                        placeholder="请选择集群"
                        onSelect={this.handleSelectFn}
                    >
                        {
                            clusterOptionList.map((item) => {
                                return <Option value={item.id}>{item.name}</Option>
                            })
                        }
                    </Select>
                    <span className='spanLeft' style={{ marginLeft: 20 }}>时间:</span>
                    <RangePicker
                        disabledDate={disabledDate}
                        suffixIcon={<span className="emicon emicon-calendar" />}
                        separator="-"
                        value={clusterTime}
                        allowClear={false}
                        onChange={this.clusterTime}
                    />
                    {/* 统计设置里面的弹窗 */}
                    {/* 这里是form表单! */}
                    {modalShow && (
                        <StatisticalSetting visible={modalShow} ModalClose={this.modalClose} />
                    )}
                </header>
                <div className='clusterContent'>
                    <div className='inspection-report-body' id='inspection-report-body'>
                    <div className='formItem-title'> 指标汇总</div>
                    <div className='IndicatorsSummary'>
                        <Row>
                            <Col span={8}>
                                <Descriptions
                                    title="服务状态统计(最新)"
                                >
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-yunhangyichang" />
                                            </div>
                                            <div className='iconRight'>运行异常服务数: {service_status?.running_fail_num}</div>
                                        </div>
                                    </Descriptions.Item>
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-zhujidownji" />
                                            </div>
                                            <div className='iconRight'>主机down机服务数: {service_status?.host_down_num}</div>
                                        </div>
                                    </Descriptions.Item>
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-jiankangjiancha" />
                                            </div>
                                            <div className='iconRight'>健康检测异常服务数: {service_status?.healthy_check_error_num}</div>
                                        </div>
                                    </Descriptions.Item>
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-gaojingfuwushu" />
                                            </div>
                                            <div className='iconRight'>告警服务数: {service_status?.alerting_num}</div>
                                        </div>
                                    </Descriptions.Item>

                                </Descriptions>
                            </Col>
                            <Col span={8}>
                                <Descriptions
                                    title="节点最新统计(最新)"
                                >
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-jiedianzongshu" />
                                            </div>
                                            <div className='iconRight'>节点总数: {node_status?.total}</div>
                                        </div>
                                    </Descriptions.Item>
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-agentjiedianyichang" />
                                            </div>
                                            <div className='iconRight'>agent异常节点数: {node_status?.agent_error_num}</div>
                                        </div>
                                    </Descriptions.Item>
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-gaojingjiedian" />
                                            </div>
                                            <div className='iconRight'>告警节点数: {node_status?.alerting_num}</div>
                                        </div>
                                    </Descriptions.Item>
                                </Descriptions>
                            </Col>
                            <Col span={8}>
                                <Descriptions
                                    title="MySql主从同步状态(最新)"
                                >
                                    <Descriptions.Item span={3}>
                                        <div className='serviceBox'>
                                            <div className='iconLeft'>
                                                <i className="emicon emicon-mysqlzhucongtongbu" />
                                            </div>
                                            <div className='iconRight'>{mysql_slave_status == 1 ? '正常' : '异常'}</div>
                                        </div>
                                    </Descriptions.Item>
                                </Descriptions>
                            </Col>
                        </Row>
                    </div>
                    {/* 表格 */}
                    {
                        form_Title_info.map((value, index) => {
                            return (
                                <div className='tableTitleStyle' key='value.form_title'>
                                    <MoreTable
                                        showTableVisible={showTable}
                                        cluster_id={currentClusterId}
                                        listItem={value}
                                    />
                                </div>
                            )

                        })
                    }
                    {configList.map((item, index) => {
                        return (
                            <div className='chartList'>
                                {item.moduleType && item.moduleType == 3 && isHasNameNode && !item.noData &&  (
                                    <>
                                        <div className='formItem-title'>大数据服务运行</div>
                                        <div className='formItem-box'>
                                            <div className='formItem-box-item'>
                                                <div className='formItem-box-itemTop'> NameNode 内存（最新）</div>
                                                <div className='formItem-box-itemBottom'>{bigData?.name_node_mem}</div>
                                            </div>
                                            <div className='formItem-box-item'>
                                                <div className='formItem-box-itemTop'> DataNode 内存（最新）</div>
                                                <div className='formItem-box-itemBottom'>{bigData?.data_node_mem}</div>
                                            </div>
                                            <div className='formItem-box-item'>
                                                <div className='formItem-box-itemTop'>DataNode 存活数（最新）</div>
                                                <div className='formItem-box-itemBottom'>{bigData?.data_node_live_nums}</div>
                                            </div>
                                            <div className='formItem-box-item'>
                                                <div className='formItem-box-itemTop'> DataNode 死亡数（最新）</div>
                                                <div className='formItem-box-itemBottom'>{bigData?.data_node_dead_nums}</div>
                                            </div>
                                            <div className='formItem-box-item'>
                                                <div className='formItem-box-itemTop'> HDFS文件数（最新） </div>
                                                <div className='formItem-box-itemBottom'>{bigData?.hdfs_file_num}</div>
                                            </div>
                                        </div>
                                    </>
                                )}
                                {item.moduleType && item.moduleType == 2 && !item.noData && (<div className='formItem-title'>应用服务运行</div>)} 
                                {item.moduleType && item.moduleType == 4 && !item.noData && (<div className='formItem-title'>节点运行</div>)} 
                                {item.moduleType && item.moduleType !== 4 && !item.noData && (<div className='chartList-title'>各服务 Full GC 趋势</div>)}
                                {item.module && !item.noData && (
                                    <div className="chartList-sutitle">{item.module}</div>
                                )}
                                {!item.noData && (
                                    <>
                                        <MoreCharts
                                        config={item}
                                        time={clusterTime}
                                        key={index}
                                        index={index}
                                        delNoDataChart={this.delNoDataChart}
                                        />
                                        {chartTable.map((citem: any) => {
                                            return (
                                                <>
                                                {citem.insert == item?.module && (
                                                    <div className='tableTitleStyle' key='value.form_title'>
                                                    <MoreTable
                                                        showTableVisible={showTable}
                                                        cluster_id={currentClusterId}
                                                        listItem={citem}
                                                    />
                                                </div> 
                                                )}
                                                </>
                                            )
                                        })}
                                    </>
                                )}
                            </div>
                        )
                    })}
                
                </div>
                </div>
            </div>
        )
    }
}

export default ClusterInspection