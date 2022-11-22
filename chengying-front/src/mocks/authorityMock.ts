const authority = [
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
];

function getAuthorityList() {
    const obj = {};
    authority.forEach((item: string) => {
        obj[item] = true;
    })
    return obj;
}

export const authorityList = getAuthorityList();
