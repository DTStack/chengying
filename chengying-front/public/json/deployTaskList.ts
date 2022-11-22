export default {
    msg: 'ok',
    code: 0,
    data: {
        complete: 'deployed',
        count: 1,
        list: [
            {
                create_time: '2019-01-02 17:19:11',
                deploy_uuid: '296f1d67-5e3d-4ce0-99e6-c72bc7188c87',
                group: 'default',
                id: 87,
                instance_id: 49,
                ip: '172.16.10.100',
                product_name: 'DTinsight',
                product_version: '3.3.0',
                progress: 100,
                schema: '{"Version":"3.3.0","Instance":{"ConfigPaths":["conf/application.properties","conf/service.properties"],"Logs":["logs/*.log"],"Environment":{"HO_HEAP_SIZE":"512m","HO_MIN_HEAP_SIZE":"512m"},"HealthCheck":{"Shell":"./bin/health.sh 172.16.10.100 8087","Period":"20s","StartPeriod":"","Timeout":"","Retries":3},"Cmd":"./bin/base.sh","PostDeploy":"","PostUndeploy":"","PrometheusPort":"9510"},"Group":"default","DependsOn":["DTinsight_Gateway"],"Config":{"mail_from":{"Default":`${window.APPCONFIG.company}数据产品技术支持`,"Desc":"internal","Type":"internal","Value":`${window.APPCONFIG.company}数据产品技术支持test`},"mail_host":{"Default":"smtp.mxhichina.com","Desc":"internal","Type":"internal","Value":"smtp.mxhichina.com"},"mail_password":{"Default":"@@dtstack2018..","Desc":"internal","Type":"internal","Value":"@@dtstack2018.."},"mail_port":{"Default":"25","Desc":"internal","Type":"internal","Value":"25"},"mail_subject":{"Default":"RDOS","Desc":"internal","Type":"internal","Value":"RDOS"},"mail_username":{"Default":"data_support@dtstack.com","Desc":"internal","Type":"internal","Value":"data_support@dtstack.com"},"mysql_ip":{"Default":{"Host":["node1"],"IP":["172.16.10.26"],"NodeId":0,"SingleIndex":0},"Desc":"internal","Type":"internal","Value":{"Host":["node1"],"IP":["172.16.10.26"],"NodeId":0,"SingleIndex":0}},"mysql_password":{"Default":"abc123","Desc":"internal","Type":"internal","Value":"abc123"},"mysql_user":{"Default":"dtstack","Desc":"internal","Type":"internal","Value":"dtstack"},"rdos_api_gateway_ip":{"Default":{"Host":["node3"],"IP":["172.16.10.100"],"NodeId":0,"SingleIndex":0},"Desc":"internal","Type":"internal","Value":{"Host":["node3"],"IP":["172.16.10.100"],"NodeId":0,"SingleIndex":0}},"rdos_api_web_ip":{"Default":{"Host":["node3"],"IP":["172.16.10.100"],"NodeId":1,"SingleIndex":0},"Desc":"internal","Type":"internal","Value":{"Host":["node3"],"IP":["172.16.10.100"],"NodeId":1,"SingleIndex":0}},"rdos_tengine_url":{"Default":"insight.dtstack.com","Desc":"internal","Type":"internal","Value":"insight.dtstack.com"},"redis_db":{"Default":"1","Desc":"internal","Type":"internal","Value":"1"},"redis_host":{"Default":{"Host":["node1"],"IP":["172.16.10.26"],"NodeId":0,"SingleIndex":0},"Desc":"internal","Type":"internal","Value":{"Host":["node1"],"IP":["172.16.10.26"],"NodeId":0,"SingleIndex":0}},"redis_password":{"Default":"abc123","Desc":"internal","Type":"internal","Value":"abc123"},"uic_host":{"Default":{"Host":["node2"],"IP":["172.16.10.47"],"NodeId":0,"SingleIndex":0},"Desc":"internal","Type":"internal","Value":{"Host":["node2"],"IP":["172.16.10.47"],"NodeId":0,"SingleIndex":0}},"uic_port":{"Default":"82","Desc":"internal","Type":"internal","Value":"82"}},"BaseProduct":"","BaseService":"","BaseParsed":false}',
                service_name: 'DTinsighdst_API',
                service_version: '3.3.0',
                sid: 'acfec36e-ce74-4d60-9265-492a627166f2',
                status: 'health-checked',
                status_message: '',
                update_time: '2019-01-02 17:20:13'
            }
        ]
    }
};
