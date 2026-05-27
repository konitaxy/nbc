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
  if (time !== null && time !== '' && time!==undefined) {
    var date = new Date(time)
    return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
  } else {
    return ''
  }
}

export const formatDateYYYYMMDD = (time) => {
  if (time !== null && time !== '' && time!==undefined) {
    var date = new Date(time)
    return formatTimeToStr(date, 'yyyy-MM-dd')
  } else {
    return ''
  }
}
export const formatDateMMDD = (time) => {
  if (time !== null && time !== '' && time!==undefined) {
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
  if (time !== null && time !== '' && time!==undefined) {
    var date = new Date(time)
    return formatTimeToStr(date, format)
  } else {
    return ''
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
