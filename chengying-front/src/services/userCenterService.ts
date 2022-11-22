import apis from '@/constants/apis';
import * as http from '@/utils/http';

const { userCenter } = apis;

export default {
  login(params: any) {
    return http[userCenter.login.method](userCenter.login.url, params);
  },
  getMembers(params: any) {
    return http[userCenter.getMembers.method](
      userCenter.getMembers.url,
      params
    );
  },
  removeMember(params: any) {
    return http[userCenter.removeMember.method](
      userCenter.removeMember.url,
      params
    );
  },
  resetPassword(params: any) {
    return http[userCenter.resetPassword.method](
      userCenter.resetPassword.url,
      params
    );
  },
  regist(params: any) {
    return http[userCenter.regist.method](userCenter.regist.url, params);
  },
  // 个人信息修改
  motifyUserInfo(params: any) {
    return http[userCenter.motifyUserInfo.method](
      userCenter.motifyUserInfo.url,
      params
    );
  },
  // 成员管理 - 编辑信息
  modifyInfoByAdmin(params: any) {
    return http[userCenter.modifyInfoByAdmin.method](
      userCenter.modifyInfoByAdmin.url,
      params
    );
  },
  validEmail(params: any) {
    return http[userCenter.validEmail.method](
      userCenter.validEmail.url,
      params
    );
  },

  resetPasswordSelf(params: any) {
    return http[userCenter.resetPasswordSelf.method](
      userCenter.resetPasswordSelf.url,
      params
    );
  },
  logout(params?: any) {
    return http[userCenter.logout.method](userCenter.logout.url, params);
  },
  getValidCode() {
    return http[userCenter.getValidCode.method](userCenter.getValidCode.url);
  },
  checkValidCode(params: any) {
    return http[userCenter.checkValidCode.method](
      userCenter.checkValidCode.url,
      params
    );
  },
  taggleStatus(params: any, type: boolean) {
    if (type) {
      return this.enableUser(params);
    }
    return this.disableUser(params);
  },
  enableUser(params: any) {
    return http[userCenter.enableUser.method](
      userCenter.enableUser.url,
      params
    );
  },
  disableUser(params: any) {
    return http[userCenter.disableUser.method](
      userCenter.disableUser.url,
      params
    );
  },
  getLoginedUserInfo() {
    return http[userCenter.getLoginedUserInfo.method](
      userCenter.getLoginedUserInfo.url
    );
  },
  // 获取公钥
  getPublicKey() {
    return http[userCenter.getPublicKey.method](userCenter.getPublicKey.url);
  },
  // 获取角色列表
  getRoleList() {
    return http[userCenter.getRoleList.method](userCenter.getRoleList.url);
  },
  // 获取权限树
  getAuthorityTree(params: any) {
    return http[userCenter.getAuthorityTree(params).method](
      userCenter.getAuthorityTree(params).url
    );
  },
  // 获取权限code
  getRoleCodes() {
    return http[userCenter.getRoleCodes.method](userCenter.getRoleCodes.url);
  },
  // 生成部署文档
  generate (params: any) {
    return http[userCenter.generate.method](
      userCenter.generate.url,
      params
    );
  },
  // 下载信息部署文档
  downloadInfo() {
    return http[userCenter.downloadInfo.method](userCenter.downloadInfo.url);
  }
};
