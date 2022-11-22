export default {
    getRoleCodes: {
        msg: 'ok',
        code: 0,
        data: [
            'menu_cluster_manage',
            'sub_menu_cluster_manage',
            'cluster_view',
            'cluster_edit',
            'sub_menu_cluster_overview',
            'sub_menu_cluster_host',
            'sub_menu_cluster_image_store',
            'image_store_view',
            'image_store_edit',
            'menu_app_manage',
            'sub_menu_package_manage',
            'package_view',
            'package_upload_delete',
            'sub_menu_installed_app_manage',
            'installed_app_view',
            'menu_deploy_guide',
            'menu_product_overview',
            'menu_service',
            'service_view',
            'service_product_start_stop',
            'service_start_stop',
            'service_roll_restart',
            'service_config_edit',
            'service_config_distribute',
            'service_dashboard_view',
            'menu_product_host',
            'menu_product_diagnosis',
            'sub_menu_log_view',
            'log_view',
            'log_download',
            'sub_menu_event_diagnosis',
            'sub_menu_config_change',
            'menu_monitor',
            'sub_menu_dashboard',
            'sub_menu_alarm',
            'sub_menu_alarm_record',
            'alarm_record_view',
            'alarm_record_open_close',
            'sub_menu_alarm_channel',
            'alarm_channel_view',
            'alarm_channel_edit',
            'menu_user_manage',
            'sub_menu_user_manage',
            'user_view',
            'user_add',
            'user_edit',
            'user_delete',
            'user_able_disable',
            'user_reset_password',
            'sub_menu_role_manage',
            'sub_menu_user_info',
            'menu_security_audit'
        ]
    },
    getMembers: {
        msg: 'ok',
        code: 0,
        data: {
            count: 9,
            list: [
                {
                    id: 33,
                    username: 'weixing@dtstack.com',
                    password: '',
                    company: '',
                    full_name: 'weixing',
                    email: 'weixing@dtstack.com',
                    phone: '',
                    status: 0,
                    update_time: {
                        Time: '2020-09-10T11:48:42+08:00',
                        Valid: true
                    },
                    create_time: {
                        Time: '2020-09-10T11:48:21+08:00',
                        Valid: true
                    },
                    is_deleted: 0,
                    role_id: 3,
                    role_name: 'Cluster Reader',
                    UpdateTimeFormat: '2020-09-10 11:48:21'
                }, {
                    id: 29,
                    username: 'yuwan',
                    password: '',
                    company: '1231',
                    full_name: 'yuwan',
                    email: 'yuwan@dtstack.com',
                    phone: '13512312121',
                    status: 0,
                    update_time: {
                        Time: '2020-09-04T17:07:19+08:00',
                        Valid: true
                    },
                    create_time: {
                        Time: '2020-09-03T10:15:29+08:00',
                        Valid: true
                    },
                    is_deleted: 0,
                    role_id: 2,
                    role_name: 'Cluster Operator',
                    UpdateTimeFormat: '2020-09-03 10:15:29'
                }, {
                    id: 1,
                    username: 'admin@dtstack.com',
                    password: '',
                    company: 'dtstack',
                    full_name: 'admin',
                    email: 'admin@dtstack.com',
                    phone: '13345678901',
                    status: 0,
                    update_time: {
                        Time: '2020-08-28T10:14:30+08:00',
                        Valid: true
                    },
                    create_time: {
                        Time: '2020-06-03T11:06:56+08:00',
                        Valid: true
                    },
                    is_deleted: 0,
                    role_id: 1,
                    role_name: 'Administrator',
                    UpdateTimeFormat: '2020-06-03 11:06:56'
                }
            ]
        }
    },
    getRoleList: {
        msg: 'ok',
        code: 0,
        data: [
            {
                desc: '超级管理员，具备产品所有操作权限',
                id: 1,
                name: 'Administrator',
                update_time: '2020-06-03T11:06:56+08:00'
            },
            {
                desc: '集群操作人员，一般指运维人员，具有安装部署、集群运维、监控告警功能操作权限',
                id: 2,
                name: 'Cluster Operator',
                update_time: '2020-06-03T11:06:56+08:00'
            },
            {
                desc: '普通用户，只有集群的只读权限',
                id: 3,
                name: 'Cluster Reader',
                update_time: '2020-08-21T11:37:04+08:00'
            }
        ]
    },
    getAuthorityTree: {
        code: 0,
        msg: 'ok',
        data: {
            description: '超级管理员，具备产品所有操作权限',
            role_name: 'Administrator',
            permissions: {
                Permissions: [{
                    title: '应用管理',
                    code: 'menu_app_manage',
                    permission: 7,
                    children: [
                        {
                            title: '安装包管理',
                            code: 'sub_menu_package_manage',
                            permission: 7,
                            children: [
                                {
                                    title: '查看',
                                    code: 'package_view',
                                    permission: 7,
                                    children: [

                                    ],
                                    selected: true
                                },
                                {
                                    title: '产品包上传/删除',
                                    code: 'package_upload_delete',
                                    permission: 3,
                                    children: [

                                    ],
                                    selected: true
                                }
                            ],
                            selected: true
                        },
                        {
                            title: '已部署应用',
                            code: 'sub_menu_installed_app_manage',
                            permission: 7,
                            children: [
                                {
                                    title: '查看',
                                    code: 'installed_app_view',
                                    permission: 7,
                                    children: [

                                    ],
                                    selected: true
                                }
                            ],
                            selected: true
                        }
                    ],
                    selected: true
                }]
            }
        }
    },
    getLoginedUserInfo: {
        code: 0,
        data: {
            id: 1,
            username: 'test@dtstack.com',
            password: '',
            company: 'dtstack',
            full_name: 'test',
            email: 'test@dtstack.com',
            phone: '13311111111',
            status: 0,
            update_time: {
                Time: '2020-12-08T17:57:41+08:00',
                Valid: true
            },
            create_time: {
                Time: '2020-06-03T11:06:56+08:00',
                Valid: true
            },
            is_deleted: 0,
            role_id: 1,
            role_name: 'Administrator',
            UpdateTimeFormat: ''
        },
        msg: 'ok'
    },
    motifyUserInfo: {
        code: 0,
        msg: 'ok',
        data: true
    }
}
