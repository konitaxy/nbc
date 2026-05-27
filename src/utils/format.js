import { formatTimeToStr } from '@/utils/date'
import { getDict } from '@/utils/dictionary'

export const formatBoolean = (bool) => {
  if (bool !== null) {
    return bool ? '是' : '否'
  } else {
    return ''
  }
}
export const formatDate = (time) => {
  if (time !== null && time !== '') {
    var date = new Date(time)
    return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
  } else {
    return ''
  }
}

export const formatDateYYYYMMDD = (time) => {
  if (time !== null && time !== '') {
    var date = new Date(time)
    return formatTimeToStr(date, 'yyyy-MM-dd')
  } else {
    return ''
  }
}
export const formatDateMMDD = (time) => {
  if (time !== null && time !== '') {
    var date = new Date(time)
    return formatTimeToStr(date, 'MM-dd')
  } else {
    return ''
  }
}

export const formatDateFormat = (time,format) => {
  if(format == ''){
    format = 'yyyy-MM-dd'
  }
  if (time !== null && time !== '') {
    var date = new Date(time)
    return formatTimeToStr(date, format)
  } else {
    return ''
  }
}
export const formatTimeDifference = (dateStr) =>{
    // 解析输入的时间字符串（支持带时区的 ISO 8601 / RFC3339）
    const targetTime = new Date(dateStr);
    
    // 检查是否为有效时间
    if (isNaN(targetTime)) {
        return "无效时间格式";
    }

    const now = new Date();
    
    // 计算时间差（毫秒），取绝对值
    const diffMs = Math.abs(now - targetTime);
    // 转换为分钟、小时、天
    const minutes = Math.floor(diffMs / (1000 * 60));
    const hours = Math.floor(diffMs / (1000 * 60 * 60));
    const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    // 判断输出单位
    if (days >= 1) {
        return `${days}天`;
    } else if (hours >= 1) {
        return `${hours}小时`;
    } else {
        return `${minutes}分钟`;
    }
}
export const filterDict = (value, options) => {
  const rowLabel = options && options.filter(item => item.value === value)
  return rowLabel && rowLabel[0] && rowLabel[0].label
}

export const getDictFunc = async(type) => {
  const dicts = await getDict(type)
  return dicts
}

export const addYear = (day) =>{
let leapDate = new Date(day)
leapDate.setFullYear(leapDate.getFullYear() + 1)
return formatTimeToStr(leapDate, 'yyyy-MM-dd')
}


export const getCookie =(name) =>{
  const nameEQ = name + "=";
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) {
          return decodeURIComponent(c.substring(nameEQ.length, c.length));
      }
  }
  return null;
}

export const setCookie =(name, value, days)=> {
  let expires = "";
  if (days) {
      const date = new Date();
      date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
      expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + encodeURIComponent(value) + expires + "; path=/";
}

export const deleteCookie = (name)=> {
  document.cookie = name + "=; max-age=0; path=/";
}