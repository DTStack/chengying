export default {
    searchDeployLog: {
        msg: 'ok',
        code: 0,
        data: {
            result: '*************************** 1.flinkx ***************************\n'
        }
    },
    unDeployService: {
        msg: 'ok',
        code: 0,
        data: {
            deploy_uuid: '01416169-9da6-4b54-a28d-bc21502d14ff'
        }
    },
    getUnDeployList: {
        msg: 'ok',
        code: 0,
        data: {
            complete: 'undeploying',
            count: 1,
            list: [
                {
                    create_time: '2020-11-28 16:40:32',
                    deploy_uuid: '01416169-9da6-4b54-a28d-bc21502d14ff',
                    group: 'default',
                    id: 3762,
                    instance_id: 3051,
                    ip: '172.16.100.172',
                    product_name: 'DTBase',
                    product_version: '2.0.9',
                    progress: 0,
                    schema: '{"ServiceDisplay":"","Version":"3.2.5-8","Instance":{"ConfigPaths":["conf/redis.conf","conf/sentinel.conf"],"Logs":["logs/*.log"],"HealthCheck":{"Shell":"./bin/health.sh 172.16.100.172 6379","Period":"20s","StartPeriod":"","Timeout":"","Retries":2},"RunUser":"","Cmd":"bin/start_redis.sh","HARoleCmd":"./bin/show_redis_role.sh 2\\u003e/dev/null","PostDeploy":"","PostUndeploy":"","PrometheusPort":"9121"},"Group":"default","Config":{"db":{"Default":"1","Desc":"internal","Type":"internal","Value":"1"},"host":{"Default":{"Host":["em2"],"IP":["172.16.100.172"],"NodeId":1,"SingleIndex":0},"Desc":"internal","Type":"internal","Value":{"Host":["em2"],"IP":["172.16.100.172"],"NodeId":1,"SingleIndex":0}},"password":{"Default":"DT@Stack#123","Desc":"internal","Type":"internal","Value":"DT@Stack#123"},"redis_port":{"Default":"26379","Desc":"internal","Type":"internal","Value":"26379"}},"BaseProduct":"","BaseProductVersion":"","BaseService":"","BaseParsed":false,"BaseAtrribute":"","ServiceAddr":{"Host":["em2"],"IP":["172.16.100.172"],"NodeId":1,"SingleIndex":0}}',
                    service_name: 'redis',
                    service_version: '3.2.5-8',
                    sid: '0de43f27-9afb-4931-969c-b940458d3c63',
                    status: 'uninstalling',
                    status_message: '',
                    update_time: '2020-11-28 16:40:32'
                }
            ]
        }
    }
}
