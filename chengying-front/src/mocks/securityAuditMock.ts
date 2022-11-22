export default {
    getSafetyAudit: {
        code: 0,
        data: {
            count: 10,
            list: [{
                content: '',
                create_time: '2020-12-08 16:52:11',
                id: 678,
                ip: '192.168.106.218',
                module: '产品访问',
                operation: '进入EM',
                operator: 'admin@dtstack.com'
            }]
        },
        msg: 'ok'
    },
    getAuditModule: {
        code: 0,
        data: {
            count: 5,
            list: ['用户管理', '产品访问', '集群管理', '部署向导', '集群运维']
        },
        msg: 'ok'
    },
    getAuditOperation: {
        code: 0,
        data: {
            count: 5,
            list: ['创建账号', '禁用账号', '启用账号', '移除账号', '重置密码']
        },
        msg: 'ok'
    }
}
