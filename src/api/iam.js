import service from '@/utils/request'

// IAM 登录
export const iamLogin = (data) => {
  return service({
    url: '/client/iam/login',
    method: 'post',
    headers: { 'X-Auth-Mail': data.email, 'X-Auth-Password': data.password },
    data: data
  })
}

// 获取角色列表
export const getIamRolesList = () => {
  return service({
    url: '/client/iam/roles/list',
    method: 'get'
  })
}

// 获取IAM账号列表
export const getIamUserList = (params) => {
  return service({
    url: '/client/iam/list',
    method: 'get',
    params: params
  })
}

// 创建IAM账号
export const createIamUser = (data) => {
  return service({
    url: '/client/iam/create',
    method: 'post',
    data: data
  })
}

// 更新IAM账号
export const updateIamUser = (data) => {
  return service({
    url: '/client/iam/update',
    method: 'put',
    data: data
  })
}

// 删除IAM账号
export const deleteIamUser = (data) => {
  return service({
    url: '/client/iam/delete',
    method: 'delete',
    data: data
  })
}

// 禁用/启用IAM账号
export const toggleIamUserStatus = (data) => {
  return service({
    url: '/client/iam/toggleStatus',
    method: 'post',
    data: data
  })
}

// 重置IAM账号密码
export const resetIamUserPassword = (data) => {
  return service({
    url: '/client/iam/resetPassword',
    method: 'post',
    data: data
  })
}

