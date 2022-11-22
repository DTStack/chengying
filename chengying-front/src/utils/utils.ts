import { message } from 'antd';
import * as Cookies from 'js-cookie';
import { COOKIES } from '@/constants/const';

// tslint:disable
const utils = {
  // 公用的方法
  // loginOut: function () {
  //     Http.get(req_url.login_out,{}).then(ret => {
  //         Http.get(req_url.get_config,{}).then(data => {
  //             window.location.href = uicUrl + ":" + uicPort + redirectLink + window.location.href;
  //         })
  //         return ret;
  //     })
  // },
  /**
   * 获取页面宽度
   * @return {[type]} [description]
   */
  pageWidth: function () {
    return Math.max(
      document.documentElement.clientWidth,
      window.innerWidth || 0
    );
  },

  /**
   * 获取页面高度
   * @return {[type]} [description]
   */
  pageHeight: function () {
    return Math.max(
      document.documentElement.clientHeight,
      window.innerHeight || 0
    );
  },

  /**
   * 根据参数名获取URL数据
   * @param  {[type]} name [description]
   * @param  {[type]} url  [description]
   * @return {[type]}      [description]
   */
  getParameterByName: function (name: any, url: any) {
    if (!url) {
      url = window.location.href;
    }
    name = name.replace(/[[\]]/g, '\\$&');
    var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)');
    var results = regex.exec(url);
    if (!results) {
      return null;
    }
    if (!results[2]) {
      return '';
    }
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
  },

  /**
   * 获取图片的Base64格式
   * @param  {[type]}   img      [description]
   * @param  {Function} callback [description]
   * @return {[type]}            [description]
   */
  getBase64: function (img: any, callback: any) {
    const reader = new FileReader();
    reader.addEventListener('load', () => callback(reader.result));
    reader.readAsDataURL(img);
  },

  /**
   * 百分比转换
   * @param  {[type]} num       [description]
   * @param  {[type]} precision [description]
   * @return {[type]}           [description]
   */
  percent: function (num: any, precision: any) {
    if (!num || num === Infinity) {
      return 0 + '%';
    }
    if (num > 1) {
      num = 1;
    }
    precision = precision || 2;
    precision = Math.pow(10, precision);
    return Math.round(num * precision * 100) / precision + '%';
  },

  getCssText: function (object: any) {
    var str = '';
    for (var attr in object) {
      str += attr + ':' + object[attr] + ';';
    }
    return str;
  },

  formateDateTime: function (timestap: any) {
    if (!timestap) {
      return '';
    }
    const times = new Date(timestap);
    const year = times.getFullYear();
    let month = times.getMonth() + 1 + '';
    let day = times.getDate() + '';

    let hour = times.getHours() + '';
    let minutes = times.getMinutes() + '';
    let seconds = times.getSeconds() + '';

    if (month.toString().length === 1) {
      month = `0${month}`;
    }

    if (day.toString().length === 1) {
      day = `0${day}`;
    }

    if (hour.toString().length === 1) {
      hour = `0${hour}`;
    }

    if (minutes.toString().length === 1) {
      minutes = `0${minutes}`;
    }

    if (seconds.toString().length === 1) {
      seconds = `0${seconds}`;
    }
    return `${year}-${month}-${day} ${hour}:${minutes}:${seconds}`;
  },

  formateDate: function (timestap: any) {
    if (!timestap) {
      return '';
    }
    const times = new Date(timestap);
    const year = times.getFullYear();
    let month = times.getMonth() + 1 + '';
    let day = times.getDate() + '';
    if (month.toString().length === 1) {
      month = `0${month}`;
    }
    if (day.toString().length === 1) {
      day = `0${day}`;
    }
    return `${year}-${month}-${day}`;
  },

  trim: function (str: any) {
    return typeof str === 'string'
      ? str.replace(/^[\s\uFEFF\xA0]+|[\s\uFEFF\xA0]+$/g, '')
      : str;
  },

  messageFilter: function (
    data: any,
    callback: any,
    info?: any,
    hideSuccessMessage?: any
  ) {
    if (data.status.code === 1001) {
      // if(info){
      //     message.success(info);
      // }
      // message.success(data.status.msg);
      if (!hideSuccessMessage) {
        message.success('成功');
      }
      if (callback) {
        callback();
      }
    } else {
      // console.log(data.status.msg)
      // message.error(data.status.msg);
    }
  },

  // filter查询条件的构造
  getTableParam: (filter: any) => {
    const param = {};
    filter.forEach((item: any, i: any) => {
      const keyArr = item.split('-');
      param[keyArr[1]] = keyArr[2];
    });
    return param;
  },

  // 抑制告警
  setAlertControl: (data: any) => {
    const alarmSuppression: any = {};
    const acKey = data.alert_control;
    if (acKey === '1') {
      alarmSuppression.close = 1;
    } else {
      alarmSuppression.close = 0;
      if (data['alert_num_' + acKey]) {
        alarmSuppression.times_after_close = data['alert_num_' + acKey];
      }
      if (data['recover_num_' + acKey]) {
        alarmSuppression.close_after_long_time_reset_amount =
          data['recover_num_' + acKey];
        alarmSuppression.close_after_long_time_reset_unit =
          data['recover_unit_' + acKey];
      }
    }
    return alarmSuppression;
  },
  getAlertControl: (alertControl: any) => {
    const fields = {
      alert_control: '',
    };
    if (alertControl.close === 0) {
      if (alertControl.times_after_close !== 0) {
        fields.alert_control = '2';
        fields['alert_num_' + fields.alert_control] =
          alertControl.times_after_close;
      } else {
        fields.alert_control = '3';
      }
      if (alertControl.close_after_long_time_reset_amount !== '') {
        if (fields.alert_control === '2') {
          fields.alert_control = '4';
        }
        fields['alert_num_' + fields.alert_control] =
          alertControl.times_after_close;
        fields['recover_num_' + fields.alert_control] =
          alertControl.close_after_long_time_reset_amount;
        fields['recover_unit_' + fields.alert_control] =
          alertControl.close_after_long_time_reset_unit;
      }
    }
    return fields;
  },

  // 模糊匹配
  filterOption: (inputValue: any, option: any) => {
    if (option.props.children.indexOf(inputValue) !== -1) {
      return true;
    } else {
      return false;
    }
  },
  // 从url中分离参数
  getParamsFromUrl(url: any) {
    var obj = {};
    var keyvalue = [];
    var key = '';
    var value = '';
    var paraString = url.substring(url.indexOf('?') + 1, url.length).split('&');
    for (var i in paraString) {
      keyvalue = paraString[i].split('=');
      key = keyvalue[0];
      value = keyvalue[1];
      obj[key] = value;
    }
    return obj;
  },
  jsonToQuery(obj: any) {
    function cleanArray(actual: any) {
      const newArray = [];
      for (let i = 0; i < actual.length; i++) {
        if (actual[i]) {
          newArray.push(actual[i]);
        }
      }
      return newArray;
    }
    if (!obj) {
      return '';
    }
    return cleanArray(
      Object.keys(obj).map((key) => {
        if (obj[key] === undefined) {
          return '';
        }
        return encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]);
      })
    ).join('&');
  },
  // 判断Object是否为{}
  checkNullObj(obj: any) {
    return Object.keys(obj).length === 0;
  },

  // service中对于简单res的处理
  serviceCallback(res: any, successMsg?: string, callback?: Function) {
    if (res.data.code === 0) {
      message.success(successMsg);
      callback && callback();
    } else {
      message.error(res.data.msg);
    }
  },

  // 是否有权限操作
  noAuthorityToDO(authorityList: any, type: string) {
    if (!authorityList[type]) {
      message.error('权限不足，请联系管理员!');
      return true;
    }
    return false;
  },

  get k8sNamespace(): string {
    return Cookies.get(COOKIES.NAMESPACE);
  },

  setNaviKey: (top: string | undefined, sub: string | undefined): void => {
    sessionStorage.setItem('firstLevelNav', top);
    sessionStorage.setItem('siderLevelNav', sub);
  },

  // 转化为GB
  formatGBUnit: (r: string): number => {
    let result: number = 0;
    if (r.indexOf('MB') > -1) {
      result = parseFloat(r.replace('MB', '')) / 1024;
    } else if (r.indexOf('KB') > -1) {
      result = parseFloat(r.replace('KB', '')) / 1024 / 1024;
    } else if (r.indexOf('GB') > -1) {
      result = parseFloat(r.replace('GB', ''));
    } else if (r.indexOf('TB') > -1) {
      result = parseFloat(r.replace('TB', '')) * 1024;
    } else {
      result = 0;
    }
    return result;
  },
};

export default utils;
