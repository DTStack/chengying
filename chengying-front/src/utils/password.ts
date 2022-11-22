import JSEncrypt from 'jsencrypt';
import { hex2b64, b64tohex } from 'jsencrypt/lib/lib/jsbn/base64';
import { SM2 } from 'gm-crypto';

export function encryptStr(str, publicKey) {
  if (!publicKey) {
    return str;
  }
  const encrypt = new JSEncrypt();
  encrypt.setPublicKey(publicKey);
  return b64tohex(encrypt.encrypt(str) || '');
}

export function decryptStr(str) {
  return window.atob ? window.atob(str) : str;
}

// sm2加密
export function encryptSM(str, key) {
  const result = SM2.encrypt(str, key, {
    inputEncoding: 'utf8',
    outputEncoding: 'hex', // 支持 hex/base64 等格式
  });
  // 04 表示非压缩
  return '04' + result;
}

// sm2解密
export function decryptSM(str, key) {
  return SM2.decrypt(str, key, {
    inputEncoding: 'hex',
    outputEncoding: 'utf8', // 支持 hex/base64 等格式
  });
}
