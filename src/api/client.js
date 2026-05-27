import service from '@/utils/request'

export const listClient = (data) => {
  return service({
    url: '/admin/client/list',
    method: 'post',
    data
  })
}

export const setClientName = (data) => {
    return service({
      url: '/admin/client/setName',
      method: 'post',
      data
    })
}
export const setClientManager = (data) => {
    return service({
      url: '/admin/client/setManager',
      method: 'post',
      data
    })
}
export const remarkClient = (data) => {
    return service({
      url: '/admin/client/remark',
      method: 'post',
      data
    })
}

export const ddClient = (data) => {
    return service({
      url: '/admin/client/dd',
      method: 'post',
      data
    })
}
export const reviewClient = (data) => {
    return service({
      url: '/admin/client/review',
      method: 'post',
      data
    })
}
export const changeClientStatus = (data) => {
    return service({
      url: '/admin/client/changeStatus',
      method: 'post',
      data
    })
}
export const enhancedKYB = (data) => {
  return service({
    url: '/admin/client/kyb',
    method: 'post',
    data
  })
}

export const getDueDiligence = (params) => {
  return service({
    url: '/admin/client/dueDiligence/get',
    method: 'get',
    params
  })
}

export const adminLogin = (data) => {
  return service({
    url: '/admin/client/adminLogin',
    method: 'post',
    data
  })
}