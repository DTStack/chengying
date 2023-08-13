#!/bin/sh

APP_HOME=/easy-agent-server

#########################################################
# Update Config
#########################################################

escape_spec_char() {
  local var_value=$1

  var_value="${var_value//\\/\\\\}"
#  var_value="${var_value//[$'\n']/}"
  var_value="${var_value//\//\\/}"
  var_value="${var_value//./\\.}"
  var_value="${var_value//\*/\\*}"
  var_value="${var_value//^/\\^}"
  var_value="${var_value//\$/\\\$}"
  var_value="${var_value//\&/\\\&}"
  var_value="${var_value//\[/\\[}"
  var_value="${var_value//\]/\\]}"

  echo $var_value
}


update_config_var() {
  local config_path=$1
  local var_name=$2
  local var_value=$3

  if [ ! -f "$config_path" ]; then
    echo "**** Configuration file '$config_path' does not exist"
    return
  fi

  # Escaping characters in parameter value
  var_value=$(escape_spec_char "$var_value")

  echo -n "** Updating '$config_path' parameter \"$var_name\": '$var_value'... "
  echo "updated"
  sed -i -e "s#$var_name#$var_value#g" "$config_path"
}

MATRIX_CONF=${APP_HOME}/example-config.yml


DB_PORT=${DB_PORT:-"3306"}


echo "********************"
echo "* DB_HOST: ${DB_HOST}"
echo "* DB_PORT: ${DB_PORT}"
echo "* DB_USER: ${DB_USER}"
echo "* DB_PWD: ${DB_PWD}"
echo "* AGENT_IP: ${AGENT_IP}"
echo "* MATRIX_IP: ${MATRIX_IP}"
echo "********************"

update_config_var $MATRIX_CONF "{DB_HOST}" "${DB_HOST}"
update_config_var $MATRIX_CONF "{DB_PORT}" "${DB_PORT}"
update_config_var $MATRIX_CONF "{DB_USER}" "${DB_USER}"
update_config_var $MATRIX_CONF "{DB_PWD}" "${DB_PWD}"
update_config_var $MATRIX_CONF "{AGENT_IP}" "${AGENT_IP}"
update_config_var $MATRIX_CONF "{MATRIX_IP}" "${MATRIX_IP}"

function update_item() {
  local para_name=$1
  local para_val=$2
  local conf_file=$3
  if [ "$para_val" != "" ]; then
    sed -ri -e "/^${para_name}=/c ${para_name}=${para_val}" $conf_file
    sed -ri -e "/^${para_name} =/c ${para_name}=${para_val}" $conf_file
    #echo "${para_name}=${para_val}"
  fi
}

function application_prop_conf() {
  conf_file=${APP_HOME}/conf/application.properties
  for i in $( set -o posix ; set |grep ^DT_APP_ |sort -rn ); do
    key=$(echo ${i} | awk -F'=' '{print $1}' |awk -F 'DT_APP_' '{print $2}')
    key_real=$(echo ${key} | sed -e 's/_/./g')
    val=$(echo ${i} | awk -F'=' '{print $2}')
    update_item $key_real $val $conf_file
  done
}
application_prop_conf;

function hive_prop_conf() {
  conf_file=${APP_HOME}/conf/hive.properties
  for i in $( set -o posix ; set |grep ^DT_HIVE_ |sort -rn ); do
    key=$(echo ${i} | awk -F'=' '{print $1}' |awk -F 'DT_HIVE_' '{print $2}')
    val=$(echo ${i} | awk -F'=' '{print $2}')
    update_item $key $val $conf_file
  done
}
hive_prop_conf;

function service_prop_conf() {
  conf_file=${APP_HOME}/conf/service.properties
  for i in $( set -o posix ; set |grep ^DT_SERVICE_ |sort -rn ); do
    key=$(echo ${i} | awk -F'=' '{print $1}' |awk -F 'DT_SERVICE_' '{print $2}')
    val=$(echo ${i} | awk -F'=' '{print $2}')
    update_item $key $val $conf_file
  done
}
service_prop_conf;

