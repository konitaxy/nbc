// utils/validators.js
import {useLanguageStore} from '@/pinia/modules/language'

/**
 * 验证邮箱格式
 * @param {Object} rule - Element Plus 的校验规则对象
 * @param {string} value - 当前输入的值
 * @param {Function} callback - Element Plus 的校验回调函数
 */
export const validateEmail = (rule, value, callback) => {
    const lang = useLanguageStore()
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!value) {
      callback(new Error(lang.t('lang.validation.email_required')));
    } else if (!emailRegex.test(value)) {
      callback(new Error(lang.t('lang.validation.email_invalid')));
    } else {
      callback(); // 校验通过
    }
  };
  
  /**
   * 验证密码 (至少8位，包含大小写字母)
   * @param {Object} rule - Element Plus 的校验规则对象
   * @param {string} value - 当前输入的值
   * @param {Function} callback - Element Plus 的校验回调函数
   */
  export const validatePassword = (rule, value, callback) => {
    const lang = useLanguageStore()
    if (!value) {
      callback(new Error(lang.t('lang.validation.password_required')));
    } else if (value.length < 8) {
      callback(new Error(lang.t('lang.validation.password_length')));
    } else if (!/(?=.*[a-z])/.test(value)) {
      callback(new Error(lang.t('lang.validation.password_lowercase')));
    } else if (!/(?=.*[A-Z])/.test(value)) {
      callback(new Error(lang.t('lang.validation.password_uppercase')));
    } else {
      callback(); // 校验通过
    }
  };
  
  // 你可以继续添加其他通用校验规则
  // 例如：验证手机号
  export const validatePhone = (rule, value, callback) => {
    const phoneRegex = /^1[3-9]\d{9}$/;
    if (!value) {
      callback(new Error('请输入手机号'));
    } else if (!phoneRegex.test(value)) {
      callback(new Error('请输入正确的手机号'));
    } else {
      callback();
    }
  }

  export const validateVerifyCode = (rule, value, callback) => {
    const lang = useLanguageStore()
    if (!value) {
      callback(new Error(lang.t('lang.validation.verify_code_required')));
    } else if (value.length !== 6) {
      callback(new Error(lang.t('lang.validation.verify_code_length')));
    } else {
      callback();
    }
  }
  
  // 例如：验证确认密码 (需要传入要比较的密码字段名)
  
  // 导出一个包含所有通用规则的对象 (可选，方便一次性导入)
  export const commonValidators = {
    validateEmail,
    validatePassword,
    validatePhone,
    validateVerifyCode
  };